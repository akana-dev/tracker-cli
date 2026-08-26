package task

import (
	"fmt"
	"strings"
	"time"

	"tracker/internal/app"
	"tracker/internal/config"
	"tracker/internal/service"
	"tracker/internal/ui"
	"tracker/pkg/timeparse"

	"github.com/spf13/cobra"
)

var AddCmd = &cobra.Command{
	Use:     "add [название]",
	Aliases: []string{"new", "create"},
	Short:   "Создать новую задачу",
	Args:    cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")
		if err := service.ValidateTitle(title); err != nil {
			return err
		}

		start, _ := cmd.Flags().GetString("start")
		end, _ := cmd.Flags().GetString("end")
		company, _ := cmd.Flags().GetString("company")
		assignee, _ := cmd.Flags().GetString("assignee")
		solution, _ := cmd.Flags().GetString("solution")
		comment, _ := cmd.Flags().GetString("comment")
		tagNames, _ := cmd.Flags().GetStringSlice("tag")

		if err := service.ValidateComment(comment); err != nil {
			return err
		}
		if err := service.ValidateSolution(solution); err != nil {
			return err
		}

		if company == "" {
			if server, err := config.GetCurrentServer(); err == nil && server.DefaultCompany != "" {
				company = server.DefaultCompany
			}
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

		if len(tagNames) > 0 {
			var allTags []string
			for _, t := range tagNames {
				for _, part := range strings.Split(t, ",") {
					part = strings.TrimSpace(part)
					if part != "" {
						allTags = append(allTags, part)
					}
				}
			}
			if len(allTags) > 0 {
				tagIDs, err := resolveTagNamesToIDs(allTags)
				if err != nil {
					return err
				}
				payload["tag_ids"] = tagIDs
			}
		}

		c := app.GetClient()
		task, err := c.CreateTask(payload)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Создана задача %s: %s",
			ui.Ticket(task.Ticket), ui.Bold(task.Title)))
		fmt.Println(ui.Dimf("Исполнитель: %s", task.GetAssigneeDisplay()))
		return nil
	},
}

func init() {
	AddCmd.Flags().StringP("start", "s", "now", "Начало")
	AddCmd.Flags().StringP("end", "e", "", "Конец")
	AddCmd.Flags().StringP("company", "q", "", "Компания (по умолчанию — из конфига)")
	AddCmd.Flags().StringP("assignee", "a", "", "Исполнитель")
	AddCmd.Flags().StringP("solution", "S", "", "Статус")
	AddCmd.Flags().StringP("comment", "C", "", "Комментарий")
	AddCmd.Flags().StringSliceP("tag", "T", nil, "Теги задачи (через запятую: bug,urgent)")
}
