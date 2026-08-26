package task

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"tracker/internal/app"
	"tracker/internal/input"
	"tracker/internal/service"
	"tracker/internal/ui"
	"tracker/pkg/timeparse"
)

var CloseCmd = &cobra.Command{
	Use:     "close [тикет]",
	Aliases: []string{"done", "finish"},
	Short:   "Закрыть задачу",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		solution, _ := cmd.Flags().GetString("solution")

		if err := service.ValidateSolution(solution); err != nil {
			return err
		}

		c := app.GetClient()
		task, err := c.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		payload := map[string]interface{}{"solution": solution}
		if _, err := c.UpdateTask(task.ID, payload); err != nil {
			return err
		}

		fmt.Println(ui.SuccessAction("Задача закрыта", ticket))
		fmt.Println(ui.Dimf("  Статус: %s", solution))
		return nil
	},
}

var PauseCmd = &cobra.Command{
	Use:     "pause [тикет]",
	Aliases: []string{"hold"},
	Short:   "Поставить задачу на паузу",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		atStr, _ := cmd.Flags().GetString("at")

		c := app.GetClient()
		task, err := c.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		var payload map[string]interface{}
		if atStr != "" {
			pauseTime, err := timeparse.Parse(atStr)
			if err != nil {
				return fmt.Errorf("ошибка в at: %w", err)
			}
			if pauseTime.After(time.Now()) {
				return fmt.Errorf("время паузы не может быть в будущем")
			}
			payload = map[string]interface{}{
				"paused_at": pauseTime.UTC().Format(time.RFC3339),
			}
		}

		if _, err := c.PauseTask(task.ID, payload); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Warningf("Задача %s поставлена на паузу",
			ui.Ticket(ticket)))
		return nil
	},
}

var ResumeCmd = &cobra.Command{
	Use:     "resume [тикет]",
	Aliases: []string{"continue"},
	Short:   "Возобновить задачу",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		startStr, _ := cmd.Flags().GetString("start")

		c := app.GetClient()
		task, err := c.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		var payload map[string]interface{}
		if startStr != "" {
			startTime, err := timeparse.Parse(startStr)
			if err != nil {
				return fmt.Errorf("ошибка в start: %w", err)
			}
			payload = map[string]interface{}{
				"resumed_at": startTime.UTC().Format(time.RFC3339),
			}
		}

		if _, err := c.ResumeTask(task.ID, payload); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Задача %s возобновлена",
			ui.Ticket(ticket)))
		return nil
	},
}

var AssignCmd = &cobra.Command{
	Use:   "assign [тикет] [пользователь]",
	Short: "Назначить задачу исполнителю",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		username := args[1]

		if err := service.ValidateUsername(username); err != nil {
			return err
		}

		c := app.GetClient()
		task, err := c.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		payload := map[string]interface{}{"assignee_username": username}
		if _, err := c.UpdateTask(task.ID, payload); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Задача %s назначена на %s",
			ui.Ticket(ticket), ui.Bold(username)))
		return nil
	},
}

var DeleteCmd = &cobra.Command{
	Use:     "delete [тикет]",
	Aliases: []string{"del", "rm", "remove"},
	Short:   "Удалить задачу",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		force, _ := cmd.Flags().GetBool("force")

		c := app.GetClient()
		task, err := c.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		if !task.CanDelete {
			return fmt.Errorf("нет прав на удаление задачи %s", ui.Ticket(ticket))
		}

		if !force {
			if !term.IsTerminal(int(os.Stdin.Fd())) {
				return fmt.Errorf("удаление задачи требует интерактивного терминала или флага --force")
			}

			confirmed := input.ReadBool(fmt.Sprintf("Удалить задачу %s?", ui.Ticket(ticket)), false)
			if !confirmed {
				fmt.Println(ui.Dim("Отменено"))
				return nil
			}
		}

		if err := c.DeleteTask(task.ID); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Задача %s удалена",
			ui.Ticket(ticket)))
		return nil
	},
}

func init() {
	CloseCmd.Flags().StringP("solution", "s", "Решено", "Статус решения")
	PauseCmd.Flags().StringP("at", "t", "", "Время паузы (по умолчанию — текущее)")
	ResumeCmd.Flags().StringP("start", "s", "", "Время возобновления")
	DeleteCmd.Flags().BoolP("force", "f", false, "Не запрашивать подтверждение")
}
