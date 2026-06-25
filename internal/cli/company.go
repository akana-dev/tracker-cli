package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/config"
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
		companies, err := client.ListCompanies()
		if err != nil {
			return err
		}

		if len(companies) == 0 {
			fmt.Println(ui.Warning("Компании не найдены."))
			return nil
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

	companyCmd.AddCommand(companyListCmd)
	companyCmd.AddCommand(companyAddCmd)
	companyCmd.AddCommand(companyDeleteCmd)
}
