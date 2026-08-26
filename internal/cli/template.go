package cli

import (
	"fmt"
	"os"
	"strconv"
	"tracker/internal/client"
	"tracker/internal/input"
	"tracker/internal/ui"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var templateCmd = &cobra.Command{
	Use:     "template",
	Aliases: []string{"templates", "tmpl"},
	Short:   "Управление шаблонами задач",
}

var templateAddCmd = &cobra.Command{
	Use:     "add [имя]",
	Aliases: []string{"new", "create"},
	Short:   "Создать шаблон",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		company, _ := cmd.Flags().GetString("company")
		solution, _ := cmd.Flags().GetString("solution")
		isPublic, _ := cmd.Flags().GetBool("public")

		if title == "" {
			return fmt.Errorf("обязательный параметр: --title")
		}

		template, err := client.CreateTemplate(name, title, description, company, solution, isPublic)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Шаблон создан: %s (ID: %d)", template.Name, template.ID))
		return nil
	},
}

var templateListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "Показать все шаблоны",
	RunE: func(cmd *cobra.Command, args []string) error {
		includeAll, _ := cmd.Flags().GetBool("all")
		templates, err := client.ListTemplates(includeAll)
		if err != nil {
			return err
		}

		if len(templates) == 0 {
			fmt.Println(ui.Info("Шаблоны не найдены"))
			return nil
		}

		for _, t := range templates {
			publicStr := ""
			if t.IsPublic {
				publicStr = " (публичный)"
			}
			fmt.Printf("  %d. %s - %s%s\n", t.ID, t.Name, t.Title, publicStr)
		}

		return nil
	},
}

var templateShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Показать детали шаблона",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("некорректный ID: %s", args[0])
		}

		template, err := client.GetTemplateByID(id)
		if err != nil {
			return fmt.Errorf("шаблон #%d не найден: %w", id, err)
		}

		fmt.Println(ui.Bold("ID:"), template.ID)
		fmt.Println(ui.Bold("Имя:"), template.Name)
		fmt.Println(ui.Bold("Заголовок:"), template.Title)
		if template.Description != "" {
			fmt.Println(ui.Bold("Описание:"), template.Description)
		}
		if template.CompanyName != "" {
			fmt.Println(ui.Bold("Компания:"), template.CompanyName)
		}
		if template.DefaultSolution != "" {
			fmt.Println(ui.Bold("Статус:"), template.DefaultSolution)
		}
		fmt.Println(ui.Bold("Публичный:"), template.IsPublic)
		fmt.Println(ui.Bold("Владелец:"), template.OwnerUsername)

		return nil
	},
}

var templateUseCmd = &cobra.Command{
	Use:   "use [id]",
	Short: "Создать задачу из шаблона",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("некорректный ID: %s", args[0])
		}

		task, err := client.UseTemplate(id)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Задача создана: %s (%s)", task.Ticket, task.Title))
		return nil
	},
}

var templateDeleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Aliases: []string{"del", "rm", "remove"},
	Short:   "Удалить шаблон",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("некорректный ID: %s", args[0])
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				return fmt.Errorf("удаление требует интерактивного терминала или флага --force")
			}

			confirmed := input.ReadBool("Удалить шаблон?", false)
			if !confirmed {
				fmt.Println(ui.Dim("Отменено"))
				return nil
			}
		}

		if err := client.DeleteTemplate(id); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Success("Шаблон удалён"))
		return nil
	},
}

func init() {
	templateAddCmd.Flags().StringP("title", "t", "", "Заголовок задачи (обязательный)")
	templateAddCmd.Flags().StringP("description", "d", "", "Описание")
	templateAddCmd.Flags().StringP("company", "c", "", "Компания")
	templateAddCmd.Flags().StringP("solution", "s", "", "Статус по умолчанию")
	templateAddCmd.Flags().BoolP("public", "p", false, "Публичный шаблон")

	templateListCmd.Flags().BoolP("all", "a", false, "Показать все шаблоны (только для admin)")

	templateDeleteCmd.Flags().BoolP("force", "f", false, "Пропустить подтверждение")

	templateCmd.AddCommand(templateAddCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateUseCmd)
	templateCmd.AddCommand(templateDeleteCmd)
}
