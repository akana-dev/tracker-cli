package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/service"
	"tracker/internal/ui"
	"tracker/pkg/table"
	"tracker/pkg/timeparse"
)

var EditCmd = &cobra.Command{
	Use:   "edit [тикет]",
	Short: "Редактировать задачу",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		payload := map[string]interface{}{}
		changes := []string{}

		if title, _ := cmd.Flags().GetString("title"); title != "" {
			if err := service.ValidateTitle(title); err != nil {
				return err
			}
			payload["title"] = title
			changes = append(changes, "название")
		}

		if start, _ := cmd.Flags().GetString("start"); start != "" {
			startTime, err := timeparse.Parse(start)
			if err != nil {
				return fmt.Errorf("ошибка в start: %w", err)
			}
			payload["start_time"] = startTime.UTC().Format(time.RFC3339)
			changes = append(changes, fmt.Sprintf("начало=%s", startTime.Local().Format("15:04")))
		}

		if end, _ := cmd.Flags().GetString("end"); end != "" {
			endTime, err := timeparse.Parse(end)
			if err != nil {
				return fmt.Errorf("ошибка в end: %w", err)
			}
			payload["end_time"] = endTime.UTC().Format(time.RFC3339)
			changes = append(changes, fmt.Sprintf("конец=%s", endTime.Local().Format("15:04")))
		}

		if company, _ := cmd.Flags().GetString("company"); company != "" {
			payload["company_name"] = company
			changes = append(changes, "компания")
		}

		if assignee, _ := cmd.Flags().GetString("assignee"); assignee != "" {
			payload["assignee_username"] = assignee
			changes = append(changes, "исполнитель")
		}

		if solution, _ := cmd.Flags().GetString("solution"); solution != "" {
			if err := service.ValidateSolution(solution); err != nil {
				return err
			}
			payload["solution"] = solution
			changes = append(changes, "статус")
		}

		if comment, _ := cmd.Flags().GetString("comment"); comment != "" {
			if err := service.ValidateComment(comment); err != nil {
				return err
			}
			payload["comment"] = comment
			changes = append(changes, "комментарий")
		}

		if cmd.Flags().Changed("tag") {
			tagNames, _ := cmd.Flags().GetStringSlice("tag")
			tagIDs, err := resolveTagNamesToIDs(tagNames)
			if err != nil {
				return err
			}
			payload["tag_ids"] = tagIDs
			if len(tagNames) == 0 {
				changes = append(changes, "теги очищены")
			} else {
				changes = append(changes, fmt.Sprintf("теги=[%s]", strings.Join(tagNames, ", ")))
			}
		}

		if len(payload) == 0 {
			fmt.Println(ui.Warning("Нет изменений для сохранения"))
			return nil
		}

		updatedTask, err := client.UpdateTask(task.ID, payload)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Задача %s обновлена (%s)",
			ui.Ticket(ticket), strings.Join(changes, ", ")))

		fmt.Println()
		tbl := table.New("Тикет", "Дата", "Сессии", "Часы", "Задача", "Исполнитель", "Статус")
		tbl.SetColumnWidths(map[int]int{
			0: 10, 1: 12, 2: 34, 3: 6, 4: 45, 5: 25, 6: 22,
		})

		tbl.AddRow(
			ui.Ticket(updatedTask.Ticket),
			updatedTask.StartTime.Local().Format("02.01.2006"),
			service.FormatSessions(*updatedTask),
			fmt.Sprintf("%.1f", service.CalculateTaskHours(*updatedTask)),
			service.FormatTaskCell(*updatedTask),
			ui.Cyan(updatedTask.GetAssigneeDisplay()),
			service.FormatStatus(*updatedTask),
		)
		tbl.Render()

		return nil
	},
}

func init() {
	EditCmd.Flags().StringP("title", "t", "", "Новое название")
	EditCmd.Flags().StringP("start", "s", "", "Новое время начала")
	EditCmd.Flags().StringP("end", "e", "", "Новое время окончания")
	EditCmd.Flags().StringP("company", "q", "", "Новая компания")
	EditCmd.Flags().StringP("assignee", "a", "", "Новый исполнитель")
	EditCmd.Flags().StringP("solution", "S", "", "Новый статус")
	EditCmd.Flags().StringP("comment", "C", "", "Новый комментарий")
	EditCmd.Flags().StringSliceP("tag", "T", nil, "Новые теги задачи (полная замена)")
}
