package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/config"
	"tracker/internal/service"
	"tracker/internal/ui"
	"tracker/pkg/table"
)

var companyCmd = &cobra.Command{
	Use:   "company",
	Short: "Управление компаниями",
}

var companyListCmd = &cobra.Command{
	Use:   "list",
	Short: "Показать список компаний",
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		if cmd.Flags().Changed("page") && cmd.Flags().Changed("offset") {
			return fmt.Errorf("нельзя использовать --page и --offset одновременно")
		}
		if page < 1 {
			return fmt.Errorf("--page должен быть >= 1")
		}
		if offset < 0 {
			return fmt.Errorf("--offset должен быть >= 0")
		}
		if limit < 0 {
			return fmt.Errorf("--limit должен быть >= 0")
		}

		if all {
			limit = 0
			offset = 0
		} else {
			if !cmd.Flags().Changed("limit") {
				limit = service.DefaultPageSize
			}
			if cmd.Flags().Changed("page") && page > 1 {
				offset = (page - 1) * limit
			}
		}

		resp, err := client.ListCompanies(limit, offset)
		if err != nil {
			return err
		}

		companies := resp.Companies

		if len(companies) == 0 {
			fmt.Println(ui.Warning("Компании не найдены."))
			return nil
		}

		fmt.Println()

		if limit > 0 && resp.Total > 0 {
			currentPage := resp.CurrentPage()
			totalPages := resp.Pages()
			startIdx := resp.Offset + 1
			endIdx := resp.Offset + len(companies)

			fmt.Printf("%s Найдено: %s | Страница: %s | Показано: %s\n",
				ui.Bold("Компании:"),
				ui.Bold(fmt.Sprintf("%d", resp.Total)),
				ui.Cyan(fmt.Sprintf("%d из %d", currentPage, totalPages)),
				ui.Dim(fmt.Sprintf("%d-%d", startIdx, endIdx)))
		} else {
			fmt.Printf("%s Найдено: %s\n",
				ui.Bold("Компании:"),
				ui.Bold(fmt.Sprintf("%d", resp.Total)))
		}
		fmt.Println()

		tbl := table.New("ID", "Название", "Описание", "След. номер")
		for _, c := range companies {
			desc := "—"
			if c.Description != nil {
				desc = *c.Description
			}
			tbl.AddRow(
				fmt.Sprintf("%d", c.ID),
				ui.Bold(c.Name),
				desc,
				fmt.Sprintf("%d", c.NextTaskNumber),
			)
		}
		tbl.Render()

		if limit > 0 && resp.HasNext() {
			fmt.Println()
			currentPage := resp.CurrentPage()
			nextPage := currentPage + 1
			fmt.Println(ui.Dimf("Следующая страница: %s | Показать все: %s",
				ui.Cyan(fmt.Sprintf("--page %d", nextPage)),
				ui.Cyan("--all")))
		}
		fmt.Println()

		return nil
	},
}

var companyAddCmd = &cobra.Command{
	Use:   "add [название]",
	Short: "Добавить новую компанию (только admin)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.GetUserRole() != "admin" {
			return fmt.Errorf("команда доступна только администраторам")
		}

		name := strings.ToUpper(args[0])
		description, _ := cmd.Flags().GetString("description")

		if err := service.ValidateCompanyName(name); err != nil {
			return err
		}
		if err := service.ValidateCompanyDescription(description); err != nil {
			return err
		}

		company, err := client.CreateCompany(name, description)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Компания %s добавлена (ID: %d)",
			ui.Bold(company.Name), company.ID))
		return nil
	},
}

var companyDeleteCmd = &cobra.Command{
	Use:   "delete [название]",
	Short: "Удалить компанию (только admin)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.GetUserRole() != "admin" {
			return fmt.Errorf("команда доступна только администраторам")
		}

		name := strings.ToUpper(args[0])
		if err := client.DeleteCompany(name); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Компания %s удалена", ui.Bold(name)))
		return nil
	},
}

func init() {
	companyAddCmd.Flags().StringP("description", "d", "", "Описание")

	companyListCmd.Flags().BoolP("all", "a", false, "Показать все компании (без пагинации)")
	companyListCmd.Flags().IntP("page", "p", 1, "Номер страницы")
	companyListCmd.Flags().IntP("limit", "l", service.DefaultPageSize, "Количество компаний на странице")
	companyListCmd.Flags().IntP("offset", "o", 0, "Смещение от начала")

	companyCmd.AddCommand(companyListCmd)
	companyCmd.AddCommand(companyAddCmd)
	companyCmd.AddCommand(companyDeleteCmd)
}
