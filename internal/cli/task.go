package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/models"
	"tracker/internal/ui"
	"tracker/pkg/table"
	"tracker/pkg/timeparse"
)

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Управление задачами",
}

var taskAddCmd = &cobra.Command{
	Use:   "add [название]",
	Short: "Создать новую задачу",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")
		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		company, _ := cmd.Flags().GetString("company")
		assignee, _ := cmd.Flags().GetString("assignee")
		solution, _ := cmd.Flags().GetString("solution")
		comment, _ := cmd.Flags().GetString("comment")

		startTime, err := timeparse.Parse(start)
		if err != nil {
			return fmt.Errorf("ошибка в start: %w", err)
		}

		payload := map[string]interface{}{
			"title":        title,
			"start_time":   startTime.Format("2006-01-02T15:04:05"),
			"company_name": company,
		}

		if end != "" {
			endTime, err := timeparse.Parse(end)
			if err != nil {
				return fmt.Errorf("ошибка в end: %w", err)
			}
			payload["end_time"] = endTime.Format("2006-01-02T15:04:05")
		}
		if assignee != "" {
			payload["assignee_username"] = assignee
		}
		if solution != "" {
			payload["solution"] = solution
		}
		if comment != "" {
			payload["comment"] = comment
		}

		task, err := client.CreateTask(payload)
		if err != nil {
			return err
		}

		assigneeName := task.OwnerUsername
		if task.AssigneeUsername != nil {
			assigneeName = *task.AssigneeUsername
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Создана задача %s: %s",
			ui.Ticket(task.Ticket), ui.Bold(task.Title)))
		fmt.Println(ui.Dimf("Исполнитель: %s", assigneeName))
		return nil
	},
}

var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "Показать задачи",
	RunE: func(cmd *cobra.Command, args []string) error {
		params := map[string]string{}

		if today, _ := cmd.Flags().GetBool("today"); today {
			params["date_from"] = timeparse.TodayStart()
		}
		if week, _ := cmd.Flags().GetBool("week"); week {
			params["date_from"] = timeparse.WeekAgo()
		}
		if month, _ := cmd.Flags().GetBool("month"); month {
			params["date_from"] = timeparse.MonthAgo()
		}
		if company, _ := cmd.Flags().GetString("company"); company != "" {
			params["company"] = company
		}
		if solution, _ := cmd.Flags().GetString("solution"); solution != "" {
			params["solution"] = solution
		}
		if assignee, _ := cmd.Flags().GetString("assignee"); assignee != "" {
			params["assignee"] = assignee
		}
		if search, _ := cmd.Flags().GetString("search"); search != "" {
			params["search"] = search
		}

		tasks, err := client.ListTasks(params)
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			fmt.Println(ui.Warning("Задачи не найдены."))
			return nil
		}

		activeCount := 0
		pausedCount := 0
		closedCount := 0
		totalHours := 0.0
		for _, t := range tasks {
			totalHours += t.TotalHours
			if t.EndTime != nil && !t.EndTime.IsZero() {
				closedCount++
			} else if t.PausedAt != nil && !t.PausedAt.IsZero() {
				pausedCount++
			} else {
				activeCount++
			}
		}

		fmt.Println()
		fmt.Printf("%s Найдено: %s | В работе: %s | На паузе: %s | Закрыто: %s | Время: %s\n",
			ui.Bold("Задачи:"),
			ui.Bold(fmt.Sprintf("%d", len(tasks))),
			ui.Green(fmt.Sprintf("%d", activeCount)),
			ui.Yellow(fmt.Sprintf("%d", pausedCount)),
			ui.Dim(fmt.Sprintf("%d", closedCount)),
			ui.Cyan(fmt.Sprintf("%.1f ч.", totalHours)),
		)
		fmt.Println()

		tbl := table.New("Тикет", "Дата", "Сессии", "Часы", "Задача", "Исполнитель", "Статус")
		for _, t := range tasks {
			assignee := t.OwnerUsername
			if t.AssigneeUsername != nil {
				assignee = *t.AssigneeUsername
			}

			tbl.AddRow(
				ui.Ticket(t.Ticket),
				t.StartTime.Format("02.01.2006"),
				formatSessions(t),
				fmt.Sprintf("%.1f", t.TotalHours),
				ui.Bold(t.Title),
				ui.Cyan(assignee),
				formatStatus(t),
			)
		}
		tbl.Render()

		return nil
	},
}

var taskCloseCmd = &cobra.Command{
	Use:   "close [тикет]",
	Short: "Закрыть задачу",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		solution, _ := cmd.Flags().GetString("solution")

		taskID, err := findTaskID(ticket)
		if err != nil {
			return err
		}

		payload := map[string]interface{}{"solution": solution}
		if _, err := client.UpdateTask(taskID, payload); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Задача %s закрыта со статусом %s",
			ui.Ticket(ticket), ui.Bold(solution)))
		return nil
	},
}

var taskPauseCmd = &cobra.Command{
	Use:   "pause [тикет]",
	Short: "Поставить задачу на паузу",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		taskID, err := findTaskID(ticket)
		if err != nil {
			return err
		}

		if _, err := client.PauseTask(taskID); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Warningf("Задача %s поставлена на паузу",
			ui.Ticket(ticket)))
		return nil
	},
}

var taskResumeCmd = &cobra.Command{
	Use:   "resume [тикет]",
	Short: "Возобновить задачу",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		taskID, err := findTaskID(ticket)
		if err != nil {
			return err
		}

		if _, err := client.ResumeTask(taskID); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Задача %s возобновлена",
			ui.Ticket(ticket)))
		return nil
	},
}

var taskAssignCmd = &cobra.Command{
	Use:   "assign [тикет] [пользователь]",
	Short: "Назначить задачу исполнителю",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		username := args[1]

		taskID, err := findTaskID(ticket)
		if err != nil {
			return err
		}

		payload := map[string]interface{}{"assignee_username": username}
		if _, err := client.UpdateTask(taskID, payload); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Задача %s назначена на %s",
			ui.Ticket(ticket), ui.Bold(username)))
		return nil
	},
}

var taskDeleteCmd = &cobra.Command{
	Use:   "delete [тикет]",
	Short: "Удалить задачу",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		taskID, err := findTaskID(ticket)
		if err != nil {
			return err
		}

		if err := client.DeleteTask(taskID); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Задача %s удалена",
			ui.Ticket(ticket)))
		return nil
	},
}

var taskExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Экспорт задач в файл",
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		output, _ := cmd.Flags().GetString("output")

		params := map[string]string{}
		if today, _ := cmd.Flags().GetBool("today"); today {
			params["date_from"] = timeparse.TodayStart()
		}
		if week, _ := cmd.Flags().GetBool("week"); week {
			params["date_from"] = timeparse.WeekAgo()
		}
		if month, _ := cmd.Flags().GetBool("month"); month {
			params["date_from"] = timeparse.MonthAgo()
		}

		data, filename, err := client.ExportTasks(format, params)
		if err != nil {
			return err
		}

		if output == "" {
			output = filename
		}

		if err := os.WriteFile(output, data, 0644); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Экспортировано в %s", ui.Bold(output)))
		return nil
	},
}

func init() {
	taskAddCmd.Flags().StringP("start", "s", "now", "Начало")
	taskAddCmd.Flags().StringP("end", "e", "", "Конец")
	taskAddCmd.Flags().StringP("company", "q", "", "Компания")
	taskAddCmd.Flags().StringP("assignee", "a", "", "Исполнитель")
	taskAddCmd.Flags().String("solution", "", "Статус")
	taskAddCmd.Flags().StringP("comment", "C", "", "Комментарий")

	taskListCmd.Flags().BoolP("today", "t", false, "Только сегодня")
	taskListCmd.Flags().BoolP("week", "w", false, "За неделю")
	taskListCmd.Flags().BoolP("month", "m", false, "За месяц")
	taskListCmd.Flags().StringP("company", "q", "", "Фильтр по компании")
	taskListCmd.Flags().String("solution", "", "Фильтр по статусу")
	taskListCmd.Flags().StringP("assignee", "a", "", "Фильтр по исполнителю")
	taskListCmd.Flags().StringP("search", "s", "", "Поиск")

	taskCloseCmd.Flags().String("solution", "Решено", "Статус решения")

	taskExportCmd.Flags().StringP("format", "f", "csv", "Формат: csv, json, xlsx")
	taskExportCmd.Flags().StringP("output", "o", "", "Имя файла")
	taskExportCmd.Flags().BoolP("today", "t", false, "Только сегодня")
	taskExportCmd.Flags().BoolP("week", "w", false, "За неделю")
	taskExportCmd.Flags().BoolP("month", "m", false, "За месяц")

	taskCmd.AddCommand(taskAddCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskCloseCmd)
	taskCmd.AddCommand(taskPauseCmd)
	taskCmd.AddCommand(taskResumeCmd)
	taskCmd.AddCommand(taskAssignCmd)
	taskCmd.AddCommand(taskDeleteCmd)
	taskCmd.AddCommand(taskExportCmd)
}

func findTaskID(ticket string) (int, error) {
	tasks, err := client.ListTasks(map[string]string{"limit": "1000"})
	if err != nil {
		return 0, err
	}

	for _, t := range tasks {
		if t.Ticket == ticket {
			return t.ID, nil
		}
	}

	return 0, fmt.Errorf("тикет %s не найден", ticket)
}

func formatSessions(t models.Task) string {
	if len(t.Sessions) == 0 {
		return ui.Dim("—")
	}

	var lines []string
	for _, s := range t.Sessions {
		startStr := s.StartTime.Format("02.01 15:04")

		if s.EndTime != nil && !s.EndTime.IsZero() {
			endStr := s.EndTime.Format("15:04")
			lines = append(lines, fmt.Sprintf("%s - %s (%.1fч)",
				startStr, endStr, s.DurationHours))
		} else {
			if t.PausedAt != nil && !t.PausedAt.IsZero() {
				pauseDuration := t.PausedAt.Sub(s.StartTime.Time).Hours()
				lines = append(lines, ui.Paused(
					fmt.Sprintf("%s - на паузе (%.1fч)", startStr, pauseDuration)))
			} else {
				hoursWorking := time.Since(s.StartTime.Time).Hours()
				if hoursWorking > 8 {
					lines = append(lines, ui.Error(
						fmt.Sprintf("%s - в работе %.1fч", startStr, hoursWorking)))
				} else if hoursWorking > 4 {
					lines = append(lines, ui.Warning(
						fmt.Sprintf("%s - в работе %.1fч", startStr, hoursWorking)))
				} else {
					lines = append(lines, ui.InProgress(
						fmt.Sprintf("%s - в работе %.1fч", startStr, hoursWorking)))
				}
			}
		}
	}

	return strings.Join(lines, "\n")
}

func formatStatus(t models.Task) string {
	solution := "—"
	if t.Solution != nil {
		solution = *t.Solution
	}

	if t.PausedAt != nil && !t.PausedAt.IsZero() {
		return ui.Paused(fmt.Sprintf("%s (на паузе)", solution))
	}
	if t.EndTime != nil && !t.EndTime.IsZero() {
		return ui.Closed(solution)
	}
	return ui.InProgress(solution)
}
