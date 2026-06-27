package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/presets"
	"tracker/internal/ui"
	"tracker/pkg/table"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Управление экспортом",
	Long:  "Команды для управления экспортом задач, включая пресеты.",
}

var exportPresetCmd = &cobra.Command{
	Use:   "preset",
	Short: "Управление пресетами экспорта",
	Long: `Пресеты позволяют сохранять часто используемые конфигурации экспорта.

Примеры:
  tracker export preset save monthly --format xlsx --period "last month" --all-users
  tracker export --preset monthly
  tracker export preset list
  tracker export preset remove monthly`,
}

var exportPresetSaveCmd = &cobra.Command{
	Use:   "save [имя]",
	Short: "Сохранить пресет",
	Long: `Сохранить текущие параметры экспорта как пресет.

Все флаги, указанные после имени пресета, будут сохранены.

Примеры:
  tracker export preset save monthly --format xlsx --period "last month"
  tracker export preset save weekly --format csv --company COMP1 --all-users`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		preset := &presets.ExportPreset{
			Name: name,
		}

		if format, _ := cmd.Flags().GetString("format"); format != "" {
			preset.Format = format
		}

		if period, _ := cmd.Flags().GetString("period"); period != "" {
			preset.Period = period
		}

		if dateFrom, _ := cmd.Flags().GetString("date-from"); dateFrom != "" {
			preset.DateFrom = dateFrom
		}
		if dateTo, _ := cmd.Flags().GetString("date-to"); dateTo != "" {
			preset.DateTo = dateTo
		}

		if timezone, _ := cmd.Flags().GetString("timezone"); timezone != "" {
			preset.Timezone = timezone
		}

		if company, _ := cmd.Flags().GetString("company"); company != "" {
			preset.Company = company
		}
		if solution, _ := cmd.Flags().GetString("solution"); solution != "" {
			preset.Solution = solution
		}
		if assignee, _ := cmd.Flags().GetString("assignee"); assignee != "" {
			preset.Assignee = assignee
		}
		if search, _ := cmd.Flags().GetString("search"); search != "" {
			preset.Search = search
		}
		if ticket, _ := cmd.Flags().GetString("ticket"); ticket != "" {
			preset.Ticket = ticket
		}

		if openOnly, _ := cmd.Flags().GetBool("open-only"); openOnly {
			preset.OpenOnly = true
		}
		if closedOnly, _ := cmd.Flags().GetBool("closed-only"); closedOnly {
			preset.ClosedOnly = true
		}
		if pausedOnly, _ := cmd.Flags().GetBool("paused-only"); pausedOnly {
			preset.PausedOnly = true
		}
		if activeOnly, _ := cmd.Flags().GetBool("active-only"); activeOnly {
			preset.ActiveOnly = true
		}
		if allUsers, _ := cmd.Flags().GetBool("all-users"); allUsers {
			preset.AllUsers = true
		}

		if fields, _ := cmd.Flags().GetString("fields"); fields != "" {
			preset.Fields = strings.Split(fields, ",")
		}

		if desc, _ := cmd.Flags().GetString("description"); desc != "" {
			preset.Description = desc
		}

		if err := presets.Save(preset); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Пресет %s сохранён", ui.Bold(name)))
		return nil
	},
}

var exportPresetListCmd = &cobra.Command{
	Use:   "list",
	Short: "Показать список пресетов",
	RunE: func(cmd *cobra.Command, args []string) error {
		allPresets, err := presets.List()
		if err != nil {
			return err
		}

		if len(allPresets) == 0 {
			fmt.Println(ui.Warning("Пресеты не найдены."))
			fmt.Println(ui.Dim("Создайте пресет: tracker export preset save <имя> [флаги]"))
			return nil
		}

		fmt.Println()
		tbl := table.New("Имя", "Формат", "Период", "Описание")
		tbl.SetColumnWidths(map[int]int{0: 25, 1: 10, 2: 20, 3: 40})
		for _, p := range allPresets {
			desc := p.Description
			if desc == "" {
				desc = ui.Dim("—")
			}
			period := p.Period
			if period == "" {
				if p.DateFrom != "" {
					period = p.DateFrom
				} else {
					period = ui.Dim("—")
				}
			}
			tbl.AddRow(
				ui.Bold(p.Name),
				p.Format,
				period,
				desc,
			)
		}
		tbl.Render()
		fmt.Println()

		return nil
	},
}

var exportPresetShowCmd = &cobra.Command{
	Use:   "show [имя]",
	Short: "Показать содержимое пресета",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		preset, err := presets.Get(name)
		if err != nil {
			return err
		}

		fmt.Println()
		ui.Header(fmt.Sprintf("Пресет: %s", ui.CyanBold(preset.Name)))

		if preset.Description != "" {
			ui.Label("Описание", preset.Description)
		}

		ui.Label("Формат", ui.Cyan(preset.Format))

		if preset.Period != "" {
			ui.Label("Период", ui.Cyan(preset.Period))
		}
		if preset.DateFrom != "" {
			ui.Label("Дата от", ui.Cyan(preset.DateFrom))
		}
		if preset.DateTo != "" {
			ui.Label("Дата до", ui.Cyan(preset.DateTo))
		}
		if preset.Timezone != "" {
			ui.Label("Часовой пояс", preset.Timezone)
		}

		fmt.Println()
		ui.Header("Фильтры")

		if preset.Company != "" {
			ui.Label("Компания", ui.Cyan(preset.Company))
		}
		if preset.Solution != "" {
			ui.Label("Статус", preset.Solution)
		}
		if preset.Assignee != "" {
			ui.Label("Исполнитель", ui.Cyan(preset.Assignee))
		}
		if preset.Search != "" {
			ui.Label("Поиск", preset.Search)
		}
		if preset.Ticket != "" {
			ui.Label("Тикет", ui.Cyan(preset.Ticket))
		}

		if preset.OpenOnly {
			ui.Label("Только открытые", ui.StatusOK())
		}
		if preset.ClosedOnly {
			ui.Label("Только закрытые", ui.StatusOK())
		}
		if preset.PausedOnly {
			ui.Label("Только на паузе", ui.StatusOK())
		}
		if preset.ActiveOnly {
			ui.Label("Только активные", ui.StatusOK())
		}
		if preset.AllUsers {
			ui.Label("Все пользователи", ui.StatusOK())
		}

		if len(preset.Fields) > 0 {
			ui.Label("Поля", ui.Cyan(strings.Join(preset.Fields, ", ")))
		}

		fmt.Println()
		return nil
	},
}

var exportPresetRemoveCmd = &cobra.Command{
	Use:   "remove [имя]",
	Short: "Удалить пресет",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := presets.Delete(name); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Пресет %s удалён", ui.Bold(name)))
		return nil
	},
}

func init() {
	exportPresetSaveCmd.Flags().StringP("format", "f", "", "Формат экспорта (csv, json, xlsx)")
	exportPresetSaveCmd.Flags().String("period", "", "Относительный период (например: 'last month', 'this week')")
	exportPresetSaveCmd.Flags().String("date-from", "", "Дата начала (RFC3339 или относительная)")
	exportPresetSaveCmd.Flags().String("date-to", "", "Дата конца (RFC3339 или относительная)")
	exportPresetSaveCmd.Flags().String("timezone", "", "Часовой пояс")
	exportPresetSaveCmd.Flags().StringP("company", "q", "", "Компания")
	exportPresetSaveCmd.Flags().String("solution", "", "Статус решения")
	exportPresetSaveCmd.Flags().StringP("assignee", "a", "", "Исполнитель")
	exportPresetSaveCmd.Flags().StringP("search", "s", "", "Поиск")
	exportPresetSaveCmd.Flags().String("ticket", "", "Тикет")
	exportPresetSaveCmd.Flags().Bool("open-only", false, "Только открытые")
	exportPresetSaveCmd.Flags().Bool("closed-only", false, "Только закрытые")
	exportPresetSaveCmd.Flags().Bool("paused-only", false, "Только на паузе")
	exportPresetSaveCmd.Flags().Bool("active-only", false, "Только активные")
	exportPresetSaveCmd.Flags().Bool("all-users", false, "Все пользователи")
	exportPresetSaveCmd.Flags().String("fields", "", "Поля для экспорта (через запятую)")
	exportPresetSaveCmd.Flags().String("description", "", "Описание пресета")

	exportPresetCmd.AddCommand(exportPresetSaveCmd)
	exportPresetCmd.AddCommand(exportPresetListCmd)
	exportPresetCmd.AddCommand(exportPresetShowCmd)
	exportPresetCmd.AddCommand(exportPresetRemoveCmd)

	exportCmd.AddCommand(exportPresetCmd)
}
