package task

import (
	"fmt"
	"os"
	"strings"

	"tracker/internal/app"
	"tracker/internal/input"
	"tracker/internal/models"
	"tracker/internal/ui"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const maxBulkTasks = 100

var BulkCmd = &cobra.Command{
	Use:     "bulk",
	Aliases: []string{"mass"},
	Short:   "Массовые операции над задачами",
}

var bulkCloseCmd = &cobra.Command{
	Use:   "close [тикет1] [тикет2] ...",
	Short: "Массовое закрытие задач",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > maxBulkTasks {
			return fmt.Errorf("слишком много задач: %d (максимум %d)", len(args), maxBulkTasks)
		}

		c := app.GetClient()
		var taskIDs []int
		var errors []string

		err := ui.WithSpinner(fmt.Sprintf("Поиск %d задач", len(args)), func() error {
			for _, ticket := range args {
				task, err := c.GetTaskByTicket(strings.ToUpper(ticket))
				if err != nil {
					errors = append(errors, fmt.Sprintf("%s: %v", ticket, err))
					continue
				}
				taskIDs = append(taskIDs, task.ID)
			}
			return nil
		})
		if err != nil {
			return err
		}

		if len(taskIDs) == 0 {
			return fmt.Errorf("не найдено ни одной задачи для закрытия")
		}

		result, err := ui.WithSpinnerResult("Закрытие задач", func() (*models.BulkResponse, error) {
			return c.BulkCloseTasks(taskIDs)
		})
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(ui.Checkmark(), ui.Successf("Обработано задач: %d", result.Total))
		if result.Succeeded > 0 {
			fmt.Println(ui.Successf("  Успешно закрыто: %d", result.Succeeded))
		}
		if result.Failed > 0 {
			fmt.Println(ui.Warningf("  Ошибок: %d", result.Failed))
			for _, r := range result.Results {
				if r.Status == "error" || r.Status == "skipped" {
					fmt.Printf("    - Задача #%d: %s\n", r.TaskID, r.Detail)
				}
			}
		}
		if len(errors) > 0 {
			fmt.Println(ui.Warning("Ошибки поиска задач:"))
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

		if len(tickets) > maxBulkTasks {
			return fmt.Errorf("слишком много задач: %d (максимум %d)", len(tickets), maxBulkTasks)
		}

		c := app.GetClient()
		var taskIDs []int
		var errors []string

		err := ui.WithSpinner(fmt.Sprintf("Поиск %d задач", len(tickets)), func() error {
			for _, ticket := range tickets {
				task, err := c.GetTaskByTicket(strings.ToUpper(ticket))
				if err != nil {
					errors = append(errors, fmt.Sprintf("%s: %v", ticket, err))
					continue
				}
				taskIDs = append(taskIDs, task.ID)
			}
			return nil
		})
		if err != nil {
			return err
		}

		if len(taskIDs) == 0 {
			return fmt.Errorf("не найдено ни одной задачи для назначения")
		}

		result, err := ui.WithSpinnerResult(
			fmt.Sprintf("Назначение на %s", ui.Bold(assignee)),
			func() (*models.BulkResponse, error) {
				return c.BulkAssignTasks(taskIDs, assignee)
			},
		)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(ui.Checkmark(), ui.Successf("Обработано задач: %d", result.Total))
		if result.Succeeded > 0 {
			fmt.Println(ui.Successf("  Успешно назначено: %d", result.Succeeded))
		}
		if result.Failed > 0 {
			fmt.Println(ui.Warningf("  Ошибок: %d", result.Failed))
			for _, r := range result.Results {
				if r.Status == "error" || r.Status == "skipped" {
					fmt.Printf("    - Задача #%d: %s\n", r.TaskID, r.Detail)
				}
			}
		}
		if len(errors) > 0 {
			fmt.Println(ui.Warning("Ошибки поиска задач:"))
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
		if len(args) > maxBulkTasks {
			return fmt.Errorf("слишком много задач: %d (максимум %d)", len(args), maxBulkTasks)
		}

		force, _ := cmd.Flags().GetBool("force")

		c := app.GetClient()
		var taskIDs []int
		var taskInfos []struct {
			Ticket string
			Title  string
			ID     int
		}
		var errors []string

		err := ui.WithSpinner(fmt.Sprintf("Поиск %d задач", len(args)), func() error {
			for _, ticket := range args {
				task, err := c.GetTaskByTicket(strings.ToUpper(ticket))
				if err != nil {
					errors = append(errors, fmt.Sprintf("%s: %v", ticket, err))
					continue
				}
				if !task.CanDelete {
					errors = append(errors, fmt.Sprintf("%s: нет прав на удаление", ticket))
					continue
				}
				taskIDs = append(taskIDs, task.ID)
				taskInfos = append(taskInfos, struct {
					Ticket string
					Title  string
					ID     int
				}{
					Ticket: task.Ticket,
					Title:  task.Title,
					ID:     task.ID,
				})
			}
			return nil
		})
		if err != nil {
			return err
		}

		if len(taskIDs) == 0 {
			return fmt.Errorf("не найдено ни одной задачи для удаления")
		}

		if !force && ui.IsInteractive() {
			fmt.Println()
			ui.Header("Подтверждение удаления")
			fmt.Println()

			for i, info := range taskInfos {
				fmt.Printf("  %s %s  %s\n",
					ui.Dim(fmt.Sprintf("%d.", i+1)),
					ui.Ticket(info.Ticket),
					ui.Bold(info.Title),
				)
			}
			fmt.Println()

			var confirmed bool
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewConfirm().
						Title(fmt.Sprintf("Удалить %d задач?", len(taskInfos))).
						Description("Это действие нельзя отменить").
						Affirmative("Да, удалить").
						Negative("Отмена").
						Value(&confirmed),
				),
			)

			if err := form.Run(); err != nil || !confirmed {
				fmt.Println(ui.Dim("Отменено"))
				return nil
			}
		} else if !force {
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				return fmt.Errorf("массовое удаление требует интерактивного терминала или флага --force")
			}
			prompt := fmt.Sprintf("Вы действительно хотите удалить %d задач?", len(taskInfos))
			confirmed := input.ReadBool(prompt, false)
			if !confirmed {
				fmt.Println(ui.Dim("Отменено"))
				return nil
			}
		}

		result, err := ui.WithSpinnerResult("Удаление задач", func() (*models.BulkResponse, error) {
			return c.BulkDeleteTasks(taskIDs)
		})
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(ui.Checkmark(), ui.Successf("Обработано задач: %d", result.Total))
		if result.Succeeded > 0 {
			fmt.Println(ui.Successf("  Успешно удалено: %d", result.Succeeded))
		}
		if result.Failed > 0 {
			fmt.Println(ui.Warningf("  Ошибок: %d", result.Failed))
			for _, r := range result.Results {
				if r.Status == "error" || r.Status == "skipped" {
					fmt.Printf("    - Задача #%d: %s\n", r.TaskID, r.Detail)
				}
			}
		}
		if len(errors) > 0 {
			fmt.Println(ui.Warning("Ошибки поиска задач:"))
			for _, e := range errors {
				fmt.Printf("  - %s\n", e)
			}
		}
		fmt.Println()
		return nil
	},
}

func init() {
	bulkDeleteCmd.Flags().BoolP("force", "f", false, "Пропустить подтверждение")

	BulkCmd.AddCommand(bulkCloseCmd)
	BulkCmd.AddCommand(bulkAssignCmd)
	BulkCmd.AddCommand(bulkDeleteCmd)
}
