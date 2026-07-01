package task

import (
	"fmt"
	"strings"
	"tracker/internal/client"
	"tracker/internal/ui"

	"github.com/spf13/cobra"
)

var BulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "Массовые операции над задачами",
}

var bulkCloseCmd = &cobra.Command{
	Use:   "close [тикет1] [тикет2] ...",
	Short: "Массовое закрытие задач",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var taskIDs []int
		var errors []string

		for _, ticket := range args {
			task, err := client.GetTaskByTicket(strings.ToUpper(ticket))
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", ticket, err))
				continue
			}
			taskIDs = append(taskIDs, task.ID)
		}

		if len(taskIDs) == 0 {
			return fmt.Errorf("не найдено ни одной задачи для закрытия")
		}

		result, err := client.BulkCloseTasks(taskIDs)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Обработано задач: %d", result.Total))
		if result.Succeeded > 0 {
			fmt.Println(ui.Success, ui.Successf("  ✓ Успешно закрыто: %d", result.Succeeded))
		}
		if result.Failed > 0 {
			fmt.Println(ui.Warning, ui.Warningf("  ⚠ Ошибок: %d", result.Failed))
			for _, r := range result.Results {
				if r.Status == "error" || r.Status == "skipped" {
					fmt.Printf("    - Задача #%d: %s\n", r.TaskID, r.Detail)
				}
			}
		}

		if len(errors) > 0 {
			fmt.Println(ui.Warning, ui.Warning("Ошибки поиска задач:"))
			for _, e := range errors {
				fmt.Printf("  - %s\n", e)
			}
		}

		return nil
	},
}

var bulkAssignCmd = &cobra.Command{
	Use:   "assign [username] [тикет1] [тикет2] ...",
	Short: "Массовое назначение исполнителя",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		assignee := args[0]
		tickets := args[1:]

		var taskIDs []int
		var errors []string

		for _, ticket := range tickets {
			task, err := client.GetTaskByTicket(strings.ToUpper(ticket))
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", ticket, err))
				continue
			}
			taskIDs = append(taskIDs, task.ID)
		}

		if len(taskIDs) == 0 {
			return fmt.Errorf("не найдено ни одной задачи для назначения")
		}

		result, err := client.BulkAssignTasks(taskIDs, assignee)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Обработано задач: %d", result.Total))
		if result.Succeeded > 0 {
			fmt.Println(ui.Success, ui.Successf("  ✓ Успешно назначено: %d", result.Succeeded))
		}
		if result.Failed > 0 {
			fmt.Println(ui.Warning, ui.Warningf("  ⚠ Ошибок: %d", result.Failed))
			for _, r := range result.Results {
				if r.Status == "error" || r.Status == "skipped" {
					fmt.Printf("    - Задача #%d: %s\n", r.TaskID, r.Detail)
				}
			}
		}

		if len(errors) > 0 {
			fmt.Println(ui.Warning, ui.Warning("Ошибки поиска задач:"))
			for _, e := range errors {
				fmt.Printf("  - %s\n", e)
			}
		}

		return nil
	},
}

var bulkDeleteCmd = &cobra.Command{
	Use:   "delete [тикет1] [тикет2] ...",
	Short: "Массовое удаление задач",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Printf(ui.Warning("Вы действительно хотите удалить %d задач? [y/N]: "), len(args))
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" {
				return nil
			}
		}

		var taskIDs []int
		var errors []string

		for _, ticket := range args {
			task, err := client.GetTaskByTicket(strings.ToUpper(ticket))
			if err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", ticket, err))
				continue
			}
			taskIDs = append(taskIDs, task.ID)
		}

		if len(taskIDs) == 0 {
			return fmt.Errorf("не найдено ни одной задачи для удаления")
		}

		result, err := client.BulkDeleteTasks(taskIDs)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Обработано задач: %d", result.Total))
		if result.Succeeded > 0 {
			fmt.Println(ui.Success, ui.Successf("  ✓ Успешно удалено: %d", result.Succeeded))
		}
		if result.Failed > 0 {
			fmt.Println(ui.Warning, ui.Warningf("  ⚠ Ошибок: %d", result.Failed))
			for _, r := range result.Results {
				if r.Status == "error" || r.Status == "skipped" {
					fmt.Printf("    - Задача #%d: %s\n", r.TaskID, r.Detail)
				}
			}
		}

		if len(errors) > 0 {
			fmt.Println(ui.Warning, ui.Warning("Ошибки поиска задач:"))
			for _, e := range errors {
				fmt.Printf("  - %s\n", e)
			}
		}

		return nil
	},
}

func init() {
	bulkDeleteCmd.Flags().BoolP("force", "f", false, "Пропустить подтверждение")

	BulkCmd.AddCommand(bulkCloseCmd)
	BulkCmd.AddCommand(bulkAssignCmd)
	BulkCmd.AddCommand(bulkDeleteCmd)
}
