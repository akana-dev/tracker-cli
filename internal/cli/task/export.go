package task

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/export"
	"tracker/internal/presets"
	"tracker/internal/ui"
)

var ExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Экспорт задач в файл",
	Long: `Экспорт задач в файл различных форматов с гибкой фильтрацией.

Поддерживаемые форматы: csv, json, xlsx

Примеры:
  tracker task export --format csv --output tasks.csv
  tracker task export --preset monthly
  tracker task export --format xlsx --period "last month"
  tracker task export --format csv --preview --today
  tracker task export --interactive
  tracker task export --format csv --fields ticket,title,hours --today`,
	RunE: func(cmd *cobra.Command, args []string) error {
		interactive, _ := cmd.Flags().GetBool("interactive")
		if interactive {
			return export.RunInteractive()
		}

		presetName, _ := cmd.Flags().GetString("preset")
		var preset *presets.ExportPreset
		if presetName != "" {
			var err error
			preset, err = presets.Get(presetName)
			if err != nil {
				return err
			}
			fmt.Println(ui.Dimf("Используется пресет: %s", ui.Bold(presetName)))
		}

		format, _ := cmd.Flags().GetString("format")
		output, _ := cmd.Flags().GetString("output")
		filename, _ := cmd.Flags().GetString("filename")
		timezone, _ := cmd.Flags().GetString("timezone")
		period, _ := cmd.Flags().GetString("period")
		dateFrom, _ := cmd.Flags().GetString("date-from")
		dateTo, _ := cmd.Flags().GetString("date-to")
		solution, _ := cmd.Flags().GetString("solution")
		company, _ := cmd.Flags().GetString("company")
		search, _ := cmd.Flags().GetString("search")
		assignee, _ := cmd.Flags().GetString("assignee")
		ticket, _ := cmd.Flags().GetString("ticket")

		openOnly, _ := cmd.Flags().GetBool("open-only")
		closedOnly, _ := cmd.Flags().GetBool("closed-only")
		pausedOnly, _ := cmd.Flags().GetBool("paused-only")
		activeOnly, _ := cmd.Flags().GetBool("active-only")
		allUsers, _ := cmd.Flags().GetBool("all-users")

		fieldsStr, _ := cmd.Flags().GetString("fields")
		fieldsPreset, _ := cmd.Flags().GetString("fields-preset")
		preview, _ := cmd.Flags().GetBool("preview")

		if today, _ := cmd.Flags().GetBool("today"); today {
			dateFrom = "today"
		}
		if week, _ := cmd.Flags().GetBool("week"); week {
			dateFrom = "last 7 days"
		}
		if month, _ := cmd.Flags().GetBool("month"); month {
			dateFrom = "last 30 days"
		}

		if preset != nil {
			if !cmd.Flags().Changed("format") && preset.Format != "" {
				format = preset.Format
			}
			if !cmd.Flags().Changed("period") && preset.Period != "" {
				period = preset.Period
			}
			if !cmd.Flags().Changed("date-from") && preset.DateFrom != "" {
				dateFrom = preset.DateFrom
			}
			if !cmd.Flags().Changed("date-to") && preset.DateTo != "" {
				dateTo = preset.DateTo
			}
			if !cmd.Flags().Changed("timezone") && preset.Timezone != "" {
				timezone = preset.Timezone
			}
			if !cmd.Flags().Changed("company") && preset.Company != "" {
				company = preset.Company
			}
			if !cmd.Flags().Changed("solution") && preset.Solution != "" {
				solution = preset.Solution
			}
			if !cmd.Flags().Changed("assignee") && preset.Assignee != "" {
				assignee = preset.Assignee
			}
			if !cmd.Flags().Changed("search") && preset.Search != "" {
				search = preset.Search
			}
			if !cmd.Flags().Changed("ticket") && preset.Ticket != "" {
				ticket = preset.Ticket
			}
			if !cmd.Flags().Changed("open-only") && preset.OpenOnly {
				openOnly = true
			}
			if !cmd.Flags().Changed("closed-only") && preset.ClosedOnly {
				closedOnly = true
			}
			if !cmd.Flags().Changed("paused-only") && preset.PausedOnly {
				pausedOnly = true
			}
			if !cmd.Flags().Changed("active-only") && preset.ActiveOnly {
				activeOnly = true
			}
			if !cmd.Flags().Changed("all-users") && preset.AllUsers {
				allUsers = true
			}
			if !cmd.Flags().Changed("fields") && len(preset.Fields) > 0 {
				fieldsStr = strings.Join(preset.Fields, ",")
			}
		}

		if format == "" {
			format = "csv"
		}
		if timezone == "" {
			timezone = "Europe/Moscow"
		}

		resolvedFrom, resolvedTo, err := export.ResolveDates(period, dateFrom, dateTo)
		if err != nil {
			return err
		}

		fields := export.ResolveFields(fieldsStr, fieldsPreset)

		params := map[string]string{
			"format":   format,
			"timezone": timezone,
		}
		if filename != "" {
			params["filename"] = filename
		}
		if resolvedFrom != "" {
			params["date_from"] = resolvedFrom
		}
		if resolvedTo != "" {
			params["date_to"] = resolvedTo
		}
		if solution != "" {
			params["solution"] = solution
		}
		if company != "" {
			params["company"] = company
		}
		if search != "" {
			params["search"] = search
		}
		if assignee != "" {
			params["assignee"] = assignee
		}
		if ticket != "" {
			params["ticket"] = ticket
		}
		if len(fields) > 0 {
			params["fields"] = strings.Join(fields, ",")
		}

		if openOnly {
			params["open_only"] = "true"
		}
		if closedOnly {
			params["closed_only"] = "true"
		}
		if pausedOnly {
			params["paused_only"] = "true"
		}
		if activeOnly {
			params["active_only"] = "true"
		}
		if allUsers {
			params["all_users"] = "true"
		}

		if preview {
			return export.RunPreview(params, format)
		}

		data, apiFilename, err := client.ExportTasks(params)
		if err != nil {
			return err
		}

		if output == "" {
			if filename != "" {
				output = filename
			} else {
				output = apiFilename
			}
		}

		if err := os.WriteFile(output, data, 0644); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Экспортировано в %s", ui.Bold(output)))
		summary := buildExportSummary(params, resolvedFrom, resolvedTo)
		if summary != "" {
			fmt.Println(ui.Dim(summary))
		}

		return nil
	},
}

func buildExportSummary(params map[string]string, dateFrom, dateTo string) string {
	var parts []string

	if dateFrom != "" {
		parts = append(parts, fmt.Sprintf("Период: %s", dateFrom))
		if dateTo != "" {
			parts[len(parts)-1] += fmt.Sprintf(" — %s", dateTo)
		}
	}
	if v, ok := params["company"]; ok && v != "" {
		parts = append(parts, fmt.Sprintf("Компания: %s", v))
	}
	if v, ok := params["solution"]; ok && v != "" {
		parts = append(parts, fmt.Sprintf("Статус: %s", v))
	}
	if v, ok := params["assignee"]; ok && v != "" {
		parts = append(parts, fmt.Sprintf("Исполнитель: %s", v))
	}
	if v, ok := params["all_users"]; ok && v == "true" {
		parts = append(parts, "Все пользователи")
	}

	if len(parts) == 0 {
		return ""
	}
	return "  " + strings.Join(parts, " | ")
}

func init() {
	ExportCmd.Flags().StringP("format", "f", "", "Формат экспорта: csv, json, xlsx")
	ExportCmd.Flags().StringP("output", "o", "", "Имя выходного файла")
	ExportCmd.Flags().String("filename", "", "Имя файла для скачивания")
	ExportCmd.Flags().String("timezone", "", "Часовой пояс (по умолчанию: Europe/Moscow)")

	ExportCmd.Flags().String("period", "", "Относительный период (today, week, month, 'last 7 days')")
	ExportCmd.Flags().String("date-from", "", "Дата начала (RFC3339 или относительная)")
	ExportCmd.Flags().String("date-to", "", "Дата конца (RFC3339 или относительная)")
	ExportCmd.Flags().BoolP("today", "t", false, "Только сегодня")
	ExportCmd.Flags().BoolP("week", "w", false, "За неделю")
	ExportCmd.Flags().BoolP("month", "m", false, "За месяц")

	ExportCmd.Flags().String("solution", "", "Фильтр по статусу")
	ExportCmd.Flags().StringP("company", "q", "", "Фильтр по компании")
	ExportCmd.Flags().StringP("search", "s", "", "Поиск")
	ExportCmd.Flags().StringP("assignee", "a", "", "Фильтр по исполнителю")
	ExportCmd.Flags().String("ticket", "", "Поиск по тикету")

	ExportCmd.Flags().Bool("open-only", false, "Только открытые")
	ExportCmd.Flags().Bool("closed-only", false, "Только закрытые")
	ExportCmd.Flags().Bool("paused-only", false, "Только на паузе")
	ExportCmd.Flags().Bool("active-only", false, "Только активные")
	ExportCmd.Flags().Bool("all-users", false, "Все пользователи")

	ExportCmd.Flags().String("fields", "", "Поля через запятую")
	ExportCmd.Flags().String("fields-preset", "", "Пресет колонок: minimal, standard, full")

	ExportCmd.Flags().BoolP("interactive", "i", false, "Интерактивный режим")
	ExportCmd.Flags().Bool("preview", false, "Preview режим")
	ExportCmd.Flags().String("preset", "", "Имя пресета")
}
