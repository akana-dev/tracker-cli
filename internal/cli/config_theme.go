package cli

import (
	"fmt"
	"strings"

	"tracker/internal/ui"

	"github.com/spf13/cobra"
)

var configThemeCmd = &cobra.Command{
	Use:   "theme [название]",
	Short: "Установить или показать тему оформления",
	Long: `Управление темой оформления tracker.
Доступные темы: dark, light, plain

Примеры:
tracker config theme            # Показать текущую тему
tracker config theme dark       # Установить темную тему
tracker config theme light      # Установить светлую тему
tracker config theme plain      # Отключить цвета (для скриптов)`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			// Показать текущую тему
			theme := ui.GetTheme()
			fmt.Println()
			ui.Header("Текущая тема")
			ui.Label("Название", ui.Bold(theme.Name))
			fmt.Println()

			// Показать превью
			fmt.Println("  Превью:")
			fmt.Printf("    %s Успешная операция\n", ui.Checkmark())
			fmt.Printf("    %s Предупреждение\n", ui.Warning("!"))
			fmt.Printf("    %s Ошибка\n", ui.Cross())
			fmt.Printf("    Тикет: %s\n", ui.Ticket("TEST-123"))
			fmt.Printf("    Бейдж: %s\n", ui.Badge("active"))
			fmt.Println()

			fmt.Println("  Доступные темы:")
			for _, name := range ui.ListThemes() {
				marker := "  "
				if name == theme.Name {
					marker = ui.Success("▸ ")
				}
				fmt.Printf("  %s%s\n", marker, name)
			}
			fmt.Println()
			return nil
		}

		themeName := strings.ToLower(args[0])
		theme, err := ui.SetThemeByName(themeName)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Тема установлена: %s", ui.Bold(theme.Name)))
		fmt.Println(ui.Dim("Изменения применятся сразу для новых команд"))
		return nil
	},
}

func init() {
	configCmd.AddCommand(configThemeCmd)
}
