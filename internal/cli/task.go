package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/config"
	"tracker/internal/service"
	"tracker/internal/ui"
	"tracker/pkg/table"
	"tracker/pkg/timeparse"
)

const defaultPageSize = 20

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
			"company_name": company,
			"start_time":   startTime.UTC().Format(time.RFC3339),
		}

		if end != "" {
			endTime, err := timeparse.Parse(end)
			if err != nil {
				return fmt.Errorf("ошибка в end: %w", err)
			}
			payload["end_time"] = endTime.UTC().Format(time.RFC3339)
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

		fmt.Println(ui.Checkmark(), ui.Successf("Создана задача %s: %s",
			ui.Ticket(task.Ticket), ui.Bold(task.Title)))
		fmt.Println(ui.Dimf("Исполнитель: %s", task.GetAssigneeDisplay()))
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

		all, _ := cmd.Flags().GetBool("all")
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		if cmd.Flags().Changed("page") && cmd.Flags().Changed("offset") {
			return fmt.Errorf("нельзя использовать --page и --offset одновременно. Используйте что-то одно")
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
				limit = defaultPageSize
			}
			if cmd.Flags().Changed("page") && page > 1 {
				offset = (page - 1) * limit
			}
		}

		resp, err := client.ListTasks(params, limit, offset)
		if err != nil {
			return err
		}

		tasks := resp.Tasks

		if len(tasks) == 0 {
			fmt.Println(ui.Warning("Задачи не найдены."))
			return nil
		}

		stats := service.CalculateTasksStats(tasks)

		fmt.Println()

		headerParts := []string{
			ui.Bold(fmt.Sprintf("Найдено: %d", resp.Total)),
		}

		if limit > 0 && resp.Total > 0 {
			currentPage := resp.CurrentPage()
			totalPages := resp.Pages()
			headerParts = append(headerParts,
				fmt.Sprintf("Страница: %s", ui.Cyan(fmt.Sprintf("%d из %d", currentPage, totalPages))))

			startIdx := resp.Offset + 1
			endIdx := resp.Offset + len(tasks)
			headerParts = append(headerParts,
				fmt.Sprintf("Показано: %s", ui.Dim(fmt.Sprintf("%d-%d", startIdx, endIdx))))
		}

		headerParts = append(headerParts,
			fmt.Sprintf("В работе: %s", ui.Green(fmt.Sprintf("%d", stats.Active))),
			fmt.Sprintf("На паузе: %s", ui.Yellow(fmt.Sprintf("%d", stats.Paused))),
			fmt.Sprintf("Закрыто: %s", ui.Dim(fmt.Sprintf("%d", stats.Closed))),
			fmt.Sprintf("Время: %s", ui.Cyan(fmt.Sprintf("%.1f ч.", stats.TotalHours))),
		)

		fmt.Printf("%s %s\n", ui.Bold("Задачи:"), strings.Join(headerParts, " | "))
		fmt.Println()

		tbl := table.New("Тикет", "Дата", "Сессии", "Часы", "Задача", "Исполнитель", "Статус")
		tbl.SetColumnWidths(map[int]int{
			0: 10,
			1: 12,
			2: 34,
			3: 6,
			4: 45,
			5: 25,
			6: 22,
		})

		for _, t := range tasks {
			tbl.AddRow(
				ui.Ticket(t.Ticket),
				t.StartTime.Local().Format("02.01.2006"),
				service.FormatSessions(t),
				fmt.Sprintf("%.1f", service.CalculateTaskHours(t)),
				service.FormatTaskCell(t),
				ui.Cyan(t.GetAssigneeDisplay()),
				service.FormatStatus(t),
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

var taskViewCmd = &cobra.Command{
	Use:   "view [тикет]",
	Short: "Подробная информация о задаче",
	Long:  "Показать полную информацию о задаче с сессиями, комментарием и метаданными",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		server, _ := config.GetCurrentServer()
		serverName := "—"
		if server != nil {
			serverName = server.Name
		}

		fmt.Println()

		statusStr := service.FormatStatus(*task)
		fmt.Printf("  %s  %s\n", ui.CyanBold(task.Ticket), statusStr)
		fmt.Println()

		ui.Header("Основная информация")
		ui.Label("Название", ui.Bold(task.Title))
		ui.Label("Компания", ui.Cyan(task.CompanyName))
		ui.Label("Сервер", ui.Dim(serverName))
		ui.Label("Создатель", ui.Cyan(task.GetOwnerDisplay()))

		if task.IsAssignedToSomeone() {
			ui.Label("Исполнитель", ui.Cyan(task.GetAssigneeDisplay()))
		} else {
			ui.Label("Исполнитель", ui.Cyan(task.GetAssigneeDisplay())+ui.Dim(" (создатель)"))
		}

		fmt.Println()

		ui.Header("Время")
		ui.Label("Начало", task.StartTime.Local().Format("02.01.2006 15:04"))

		if task.IsClosed() {
			ui.Label("Окончание", task.EndTime.Local().Format("02.01.2006 15:04"))
			duration := task.EndTime.Sub(task.StartTime.Time)
			ui.Label("Длительность", service.FormatDuration(duration))
		} else {
			ui.Label("Окончание", ui.Warning("не закрыта"))
		}

		if task.IsPaused() {
			ui.Label("На паузе с", ui.Warning(task.PausedAt.Local().Format("02.01.2006 15:04")))
		}

		totalHours := service.CalculateTaskHours(*task)
		ui.Label("Отработано", ui.Cyan(fmt.Sprintf("%.1f ч.", totalHours)))

		fmt.Println()

		ui.Header("Статус и описание")

		solution := "—"
		if task.Solution != nil && *task.Solution != "" {
			solution = *task.Solution
		}
		ui.Label("Решение", statusStr+" "+solution)

		if task.Comment != nil && *task.Comment != "" {
			ui.Label("Комментарий", "")
			service.PrintIndented(*task.Comment, "    ")
		} else {
			ui.Label("Комментарий", ui.Dim("—"))
		}

		fmt.Println()

		ui.Header(fmt.Sprintf("Сессии (%d)", len(task.Sessions)))

		if len(task.Sessions) == 0 {
			fmt.Println("    " + ui.Dim("Нет сессий"))
		} else {
			for i, s := range task.Sessions {
				sessionNum := i + 1
				startLocal := s.StartTime.Time.Local()
				startStr := startLocal.Format("02.01.2006 15:04")

				fmt.Printf("    %s ", ui.Dim(fmt.Sprintf("#%d", sessionNum)))

				if s.EndTime != nil && !s.EndTime.IsZero() {
					endStr := service.FormatEndTime(s.StartTime.Time, s.EndTime.Time)
					duration := s.EndTime.Time.UTC().Sub(s.StartTime.Time.UTC())

					fmt.Printf("%s — %s  %s\n",
						startStr,
						endStr,
						ui.Cyan(fmt.Sprintf("(%s)", service.FormatDuration(duration))),
					)
				} else {
					if task.IsPaused() {
						pauseDuration := task.PausedAt.Time.UTC().Sub(s.StartTime.Time.UTC())
						fmt.Printf("%s — %s\n",
							startStr,
							ui.Paused(fmt.Sprintf("на паузе (%s)", service.FormatDuration(pauseDuration))),
						)
					} else {
						elapsed := time.Now().Sub(startLocal)
						fmt.Printf("%s — %s\n",
							startStr,
							ui.InProgress(fmt.Sprintf("в работе (%s)", service.FormatDuration(elapsed))),
						)
					}
				}
			}
		}

		fmt.Println()

		ui.Header("Права доступа")
		if task.CanEdit {
			ui.Label("Редактирование", ui.StatusOK())
		} else {
			ui.Label("Редактирование", ui.StatusNo())
		}
		if task.CanDelete {
			ui.Label("Удаление", ui.StatusOK())
		} else {
			ui.Label("Удаление", ui.StatusNo())
		}

		fmt.Println()

		fmt.Println(ui.Dim("Команды для работы с задачей:"))
		fmt.Printf("  %s  %s\n", ui.Cyan("edit"), ui.Dim("Редактировать задачу"))
		fmt.Printf("  %s  %s\n", ui.Cyan("pause"), ui.Dim("Поставить на паузу"))
		fmt.Printf("  %s  %s\n", ui.Cyan("resume"), ui.Dim("Возобновить"))
		fmt.Printf("  %s  %s\n", ui.Cyan("close"), ui.Dim("Закрыть задачу"))
		fmt.Printf("  %s  %s\n", ui.Cyan("assign"), ui.Dim("Назначить исполнителя"))
		fmt.Printf("  %s  %s\n", ui.Cyan("delete"), ui.Dim("Удалить задачу"))
		fmt.Println()

		return nil
	},
}

var taskEditCmd = &cobra.Command{
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
			payload["solution"] = solution
			changes = append(changes, "статус")
		}

		if comment, _ := cmd.Flags().GetString("comment"); comment != "" {
			payload["comment"] = comment
			changes = append(changes, "комментарий")
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

var taskCloseCmd = &cobra.Command{
	Use:   "close [тикет]",
	Short: "Закрыть задачу",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		solution, _ := cmd.Flags().GetString("solution")

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		payload := map[string]interface{}{"solution": solution}
		if _, err := client.UpdateTask(task.ID, payload); err != nil {
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

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		if _, err := client.PauseTask(task.ID); err != nil {
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

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		if _, err := client.ResumeTask(task.ID); err != nil {
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

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		payload := map[string]interface{}{"assignee_username": username}
		if _, err := client.UpdateTask(task.ID, payload); err != nil {
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

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		if err := client.DeleteTask(task.ID); err != nil {
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

	taskListCmd.Flags().Bool("all", false, "Показать все задачи (без пагинации)")
	taskListCmd.Flags().Int("page", 1, "Номер страницы (конвертируется в offset на клиенте)")
	taskListCmd.Flags().Int("limit", defaultPageSize, "Количество задач на странице")
	taskListCmd.Flags().Int("offset", 0, "Смещение от начала (альтернатива --page)")

	taskEditCmd.Flags().StringP("title", "t", "", "Новое название")
	taskEditCmd.Flags().StringP("start", "s", "", "Новое время начала")
	taskEditCmd.Flags().StringP("end", "e", "", "Новое время окончания")
	taskEditCmd.Flags().StringP("company", "q", "", "Новая компания")
	taskEditCmd.Flags().StringP("assignee", "a", "", "Новый исполнитель")
	taskEditCmd.Flags().String("solution", "", "Новый статус")
	taskEditCmd.Flags().StringP("comment", "C", "", "Новый комментарий")

	taskCloseCmd.Flags().String("solution", "Решено", "Статус решения")

	taskExportCmd.Flags().StringP("format", "f", "csv", "Формат: csv, json, xlsx")
	taskExportCmd.Flags().StringP("output", "o", "", "Имя файла")
	taskExportCmd.Flags().BoolP("today", "t", false, "Только сегодня")
	taskExportCmd.Flags().BoolP("week", "w", false, "За неделю")
	taskExportCmd.Flags().BoolP("month", "m", false, "За месяц")

	taskCmd.AddCommand(taskAddCmd)
	taskCmd.AddCommand(taskListCmd)
	taskCmd.AddCommand(taskViewCmd)
	taskCmd.AddCommand(taskEditCmd)
	taskCmd.AddCommand(taskCloseCmd)
	taskCmd.AddCommand(taskPauseCmd)
	taskCmd.AddCommand(taskResumeCmd)
	taskCmd.AddCommand(taskAssignCmd)
	taskCmd.AddCommand(taskDeleteCmd)
	taskCmd.AddCommand(taskExportCmd)
}
