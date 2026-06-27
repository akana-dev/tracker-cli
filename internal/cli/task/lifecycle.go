package task

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/service"
	"tracker/internal/ui"
)

var CloseCmd = &cobra.Command{
	Use:   "close [тикет]",
	Short: "Закрыть задачу",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		solution, _ := cmd.Flags().GetString("solution")

		if err := service.ValidateSolution(solution); err != nil {
			return err
		}

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

var PauseCmd = &cobra.Command{
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

var ResumeCmd = &cobra.Command{
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

var DeleteCmd = &cobra.Command{
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

func init() {
	CloseCmd.Flags().String("solution", "Решено", "Статус решения")
}
