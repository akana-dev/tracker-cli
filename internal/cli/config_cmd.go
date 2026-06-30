package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/config"
	"tracker/internal/ui"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Управление конфигурацией",
	Long: `Управление локальной конфигурацией трекера.

Доступные настройки:
  default-company  — компания по умолчанию для новых задач

Примеры:
  tracker config default-company COMP1
  tracker config default-company --clear
  tracker config show`,
}

var configDefaultCompanyCmd = &cobra.Command{
	Use:   "default-company [название]",
	Short: "Установить компанию по умолчанию",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		clear, _ := cmd.Flags().GetBool("clear")

		server, err := config.GetCurrentServer()
		if err != nil {
			return err
		}

		if clear {
			server.DefaultCompany = ""
			if err := config.SaveConfig(server); err != nil {
				return err
			}
			fmt.Println(ui.Checkmark(), ui.Success("Компания по умолчанию сброшена"))
			return nil
		}

		if len(args) == 0 {
			if server.DefaultCompany == "" {
				fmt.Println(ui.Dim("Компания по умолчанию не установлена"))
			} else {
				ui.Label("Компания по умолчанию", ui.Cyan(server.DefaultCompany))
			}
			return nil
		}

		name := strings.ToUpper(args[0])
		server.DefaultCompany = name
		if err := config.SaveConfig(server); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Компания по умолчанию: %s", ui.Bold(name)))
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Показать текущую конфигурацию",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetCurrentServer()
		if err != nil {
			return err
		}

		fmt.Println()
		ui.Header("Локальная конфигурация")
		ui.Label("Сервер", ui.Bold(server.Name))
		ui.Label("API URL", ui.Cyan(server.APIURL))

		defaultCompany := server.DefaultCompany
		if defaultCompany == "" {
			defaultCompany = ui.Dim("не установлена")
		} else {
			defaultCompany = ui.Cyan(defaultCompany)
		}
		ui.Label("Компания по умолчанию", defaultCompany)

		fmt.Println()
		return nil
	},
}

func init() {
	configDefaultCompanyCmd.Flags().BoolP("clear", "c", false, "Сбросить компанию по умолчанию")

	configCmd.AddCommand(configDefaultCompanyCmd)
	configCmd.AddCommand(configShowCmd)
}
