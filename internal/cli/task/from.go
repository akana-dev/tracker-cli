package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/config"
	"tracker/internal/service"
	"tracker/internal/tags"
	"tracker/internal/templates"
	"tracker/internal/ui"
	"tracker/pkg/timeparse"
)

var FromCmd = &cobra.Command{
	Use:   "from [шаблон]",
	Short: "Создать задачу из шаблона",
	Long: `Создать новую задачу на основе шаблона.

Все параметры шаблона можно переопределить через флаги.

Примеры:
  tracker task from daily-standup
  tracker task from daily-standup --title "Планёрка по проекту X"
  tracker task from daily-standup --start 10:00
  tracker task from daily-standup --tag urgent`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]

		tmpl, err := templates.Load(templateName)
		if err != nil {
			return err
		}

		title := tmpl.Title
		if t, _ := cmd.Flags().GetString("title"); t != "" {
			title = t
		}

		company := tmpl.Company
		if c, _ := cmd.Flags().GetString("company"); c != "" {
			company = c
		}
		if company == "" {
			if server, err := config.GetCurrentServer(); err == nil && server.DefaultCompany != "" {
				company = server.DefaultCompany
			}
		}

		assignee := tmpl.Assignee
		if a, _ := cmd.Flags().GetString("assignee"); a != "" {
			assignee = a
		}

		solution := tmpl.Solution
		if s, _ := cmd.Flags().GetString("solution"); s != "" {
			solution = s
		}

		comment := tmpl.Comment
		if c, _ := cmd.Flags().GetString("comment"); c != "" {
			comment = c
		}

		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")

		tagsFromTemplate := tmpl.Tags
		tagsFromFlag, _ := cmd.Flags().GetStringSlice("tag")
		allTags := append(tagsFromTemplate, tagsFromFlag...)

		if err := service.ValidateTitle(title); err != nil {
			return err
		}
		if err := service.ValidateComment(comment); err != nil {
			return err
		}
		if err := service.ValidateSolution(solution); err != nil {
			return err
		}

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

		if len(allTags) > 0 {
			if err := tags.Set(task.Ticket, allTags); err != nil {
				fmt.Println(ui.Warningf("Задача создана, но не удалось применить теги: %v", err))
			}
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Создана задача %s: %s",
			ui.Ticket(task.Ticket), ui.Bold(task.Title)))
		if len(allTags) > 0 {
			fmt.Println(ui.Dimf("Теги: %s", strings.Join(allTags, ", ")))
		}
		fmt.Println(ui.Dimf("Исполнитель: %s", task.GetAssigneeDisplay()))
		return nil
	},
}

func init() {
	FromCmd.Flags().StringP("title", "t", "", "Переопределить название")
	FromCmd.Flags().StringP("start", "s", "now", "Начало")
	FromCmd.Flags().StringP("end", "e", "", "Конец")
	FromCmd.Flags().StringP("company", "q", "", "Переопределить компанию")
	FromCmd.Flags().StringP("assignee", "a", "", "Переопределить исполнителя")
	FromCmd.Flags().String("solution", "", "Переопределить статус")
	FromCmd.Flags().StringP("comment", "C", "", "Переопределить комментарий")
	FromCmd.Flags().StringSlice("tag", nil, "Дополнительные теги")
}
