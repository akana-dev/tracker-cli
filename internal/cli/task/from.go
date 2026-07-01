package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/config"
	"tracker/internal/service"
	"tracker/internal/templates"
	"tracker/internal/ui"
	"tracker/pkg/timeparse"
)

var FromCmd = &cobra.Command{
	Use:   "from [имя_шаблона]",
	Short: "Создать задачу из шаблона",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		templateName := args[0]

		tmpl, err := templates.Load(templateName)
		if err != nil {
			return fmt.Errorf("шаблон '%s' не найден: %w", templateName, err)
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

		if err := service.ValidateTitle(title); err != nil {
			return err
		}
		if err := service.ValidateComment(comment); err != nil {
			return err
		}
		if err := service.ValidateSolution(solution); err != nil {
			return err
		}

		startStr, _ := cmd.Flags().GetString("start")
		startTime, err := timeparse.Parse(startStr)
		if err != nil {
			return fmt.Errorf("ошибка в start: %w", err)
		}

		payload := map[string]interface{}{
			"title":        title,
			"company_name": company,
			"start_time":   startTime.UTC().Format(time.RFC3339),
		}

		if endStr, _ := cmd.Flags().GetString("end"); endStr != "" {
			endTime, err := timeparse.Parse(endStr)
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

		tagNames, _ := cmd.Flags().GetStringSlice("tag")
		if len(tagNames) > 0 {
			tagIDs, err := resolveTagNamesToIDs(tagNames)
			if err != nil {
				return err
			}
			payload["tag_ids"] = tagIDs
		}

		task, err := client.CreateTask(payload)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Создана задача из шаблона '%s': %s",
			templateName, ui.Ticket(task.Ticket)))
		fmt.Println(ui.Dimf("Название: %s", task.Title))
		fmt.Println(ui.Dimf("Исполнитель: %s", task.GetAssigneeDisplay()))

		if len(tagNames) > 0 {
			fmt.Println(ui.Dimf("Теги: %s", strings.Join(tagNames, ", ")))
		}

		return nil
	},
}

func init() {
	FromCmd.Flags().StringP("title", "t", "", "Переопределить название")
	FromCmd.Flags().StringP("start", "s", "now", "Начало")
	FromCmd.Flags().StringP("end", "e", "", "Конец")
	FromCmd.Flags().StringP("company", "q", "", "Переопределить компанию")
	FromCmd.Flags().StringP("assignee", "a", "", "Переопределить исполнителя")
	FromCmd.Flags().StringP("solution", "S", "", "Переопределить статус")
	FromCmd.Flags().StringP("comment", "C", "", "Переопределить комментарий")
	FromCmd.Flags().StringSliceP("tag", "T", nil, "Дополнительные теги")
}
