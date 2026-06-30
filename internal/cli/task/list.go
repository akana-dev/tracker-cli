package task

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/models"
	"tracker/internal/service"
	"tracker/internal/tags"
	"tracker/internal/ui"
	"tracker/pkg/table"
	"tracker/pkg/timeparse"
)

var ListCmd = &cobra.Command{
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

		tagFilter, _ := cmd.Flags().GetStringSlice("tag")

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

		resp, err := client.ListTasks(params, limit, offset)
		if err != nil {
			return err
		}

		tasks := resp.Tasks

		if len(tagFilter) > 0 {
			ticketsWithTags, err := tags.FilterTicketsByTag(tagFilter)
			if err != nil {
				return err
			}
			tagSet := make(map[string]bool, len(ticketsWithTags))
			for _, t := range ticketsWithTags {
				tagSet[t] = true
			}

			var filtered []models.Task
			for _, t := range tasks {
				if tagSet[t.Ticket] {
					filtered = append(filtered, t)
				}
			}
			tasks = filtered
		}

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
			0: 10, 1: 12, 2: 34, 3: 6, 4: 45, 5: 25, 6: 22,
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

func init() {
	ListCmd.Flags().BoolP("today", "t", false, "Только сегодня")
	ListCmd.Flags().BoolP("week", "w", false, "За неделю")
	ListCmd.Flags().BoolP("month", "m", false, "За месяц")
	ListCmd.Flags().StringP("company", "q", "", "Фильтр по компании")
	ListCmd.Flags().StringP("solution", "S", "", "Фильтр по статусу")
	ListCmd.Flags().StringP("assignee", "a", "", "Фильтр по исполнителю")
	ListCmd.Flags().StringP("search", "s", "", "Поиск")
	ListCmd.Flags().StringSlice("tag", nil, "Фильтр по тегам")

	ListCmd.Flags().BoolP("all", "A", false, "Показать все задачи")
	ListCmd.Flags().IntP("page", "p", 1, "Номер страницы")
	ListCmd.Flags().IntP("limit", "l", service.DefaultPageSize, "Количество задач на странице")
	ListCmd.Flags().IntP("offset", "o", 0, "Смещение от начала")
}
