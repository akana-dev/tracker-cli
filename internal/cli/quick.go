package cli

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/state"
	"tracker/internal/ui"
)

var (
	startStopPrevious bool
)

var StartCmd = &cobra.Command{
	Use:   "start [тикет]",
	Short: "Начать работу над задачей (предыдущая ставится на паузу)",
	Long: `Начать работу над задачей. По умолчанию предыдущая активная задача 
автоматически ставится на паузу. С флагом --stop предыдущая задача будет остановлена.`,
	Aliases: []string{"s"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])

		_ = state.CleanupOldPausedTasks()

		currentState, err := state.Load()
		if err != nil {
			return fmt.Errorf("ошибка загрузки состояния: %w", err)
		}

		if currentState != nil && currentState.CurrentTicket != "" {
			prevTicket := currentState.CurrentTicket

			if prevTicket == ticket {
				fmt.Printf("%s Задача %s уже активна\n",
					ui.Info(), ui.Ticket(ticket))
				return nil
			}

			prevTask, err := client.GetTaskByTicket(prevTicket)
			if err != nil {
				return fmt.Errorf("предыдущая задача %s не найдена: %w", prevTicket, err)
			}

			duration := time.Since(currentState.StartedAt)

			if startStopPrevious {
				payload := map[string]interface{}{
					"paused_at": time.Now().UTC().Format(time.RFC3339),
				}
				if _, err := client.PauseTask(prevTask.ID, payload); err != nil {
					return fmt.Errorf("ошибка остановки %s: %w", prevTicket, err)
				}
				fmt.Printf("%s Остановлена работа над %s (время: %s)\n",
					ui.Warning(), ui.Ticket(prevTicket), formatDuration(duration))
			} else {
				payload := map[string]interface{}{
					"paused_at": time.Now().UTC().Format(time.RFC3339),
				}
				if _, err := client.PauseTask(prevTask.ID, payload); err != nil {
					return fmt.Errorf("ошибка паузы %s: %w", prevTicket, err)
				}

				currentState.AddPausedTask(state.PausedTask{
					Ticket:    prevTicket,
					TaskID:    prevTask.ID,
					PausedAt:  time.Now(),
					StartedAt: currentState.StartedAt,
					Duration:  int64(duration.Seconds()),
				})

				fmt.Printf("%s Приостановлена %s (время: %s)\n",
					ui.Warning(), ui.Ticket(prevTicket), formatDuration(duration))
			}
		} else {
			currentState = &state.State{
				PausedTasks: []state.PausedTask{},
			}
		}

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		if task.IsClosed() {
			return fmt.Errorf("задача %s уже закрыта", ui.Ticket(ticket))
		}

		if !task.IsPaused() {
			currentState.RemovePausedTask(ticket)
			currentState.CurrentTicket = ticket
			currentState.StartedAt = time.Now()
			if err := state.Save(currentState); err != nil {
				return fmt.Errorf("ошибка сохранения состояния: %w", err)
			}

			fmt.Printf("%s Задача %s уже активна\n",
				ui.Info(), ui.Ticket(ticket))
			return nil
		}

		payload := map[string]interface{}{
			"resumed_at": time.Now().UTC().Format(time.RFC3339),
		}
		if _, err := client.ResumeTask(task.ID, payload); err != nil {
			return fmt.Errorf("ошибка возобновления задачи: %w", err)
		}

		currentState.RemovePausedTask(ticket)

		currentState.CurrentTicket = ticket
		currentState.StartedAt = time.Now()
		if err := state.Save(currentState); err != nil {
			return fmt.Errorf("ошибка сохранения состояния: %w", err)
		}

		fmt.Printf("%s Начата работа над %s\n",
			ui.Checkmark(), ui.Ticket(ticket))
		return nil
	},
}

var StopCmd = &cobra.Command{
	Use:     "stop",
	Short:   "Остановить работу над текущей задачей",
	Aliases: []string{"end"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentState, err := state.Load()
		if err != nil {
			return fmt.Errorf("ошибка загрузки состояния: %w", err)
		}

		if currentState == nil || currentState.CurrentTicket == "" {
			fmt.Println(ui.Warning(), "Нет активной задачи")
			return nil
		}

		ticket := currentState.CurrentTicket
		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		if task.IsClosed() {
			currentState.CurrentTicket = ""
			currentState.StartedAt = time.Time{}
			if err := state.Save(currentState); err != nil {
				return fmt.Errorf("ошибка сохранения состояния: %w", err)
			}
			fmt.Printf("%s Задача %s уже закрыта\n",
				ui.Warning(), ui.Ticket(ticket))
			return nil
		}

		if task.IsPaused() {
			currentState.CurrentTicket = ""
			currentState.StartedAt = time.Time{}
			if err := state.Save(currentState); err != nil {
				return fmt.Errorf("ошибка сохранения состояния: %w", err)
			}
			fmt.Printf("%s Задача %s уже на паузе\n",
				ui.Warning(), ui.Ticket(ticket))
			return nil
		}

		duration := time.Since(currentState.StartedAt)

		payload := map[string]interface{}{
			"paused_at": time.Now().UTC().Format(time.RFC3339),
		}
		if _, err := client.PauseTask(task.ID, payload); err != nil {
			return fmt.Errorf("ошибка остановки задачи: %w", err)
		}

		currentState.CurrentTicket = ""
		currentState.StartedAt = time.Time{}
		if err := state.Save(currentState); err != nil {
			return fmt.Errorf("ошибка сохранения состояния: %w", err)
		}

		fmt.Printf("%s Остановлена работа над %s (время: %s)\n",
			ui.Checkmark(), ui.Ticket(ticket), formatDuration(duration))
		return nil
	},
}

var StatusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Показать статус текущей задачи",
	Aliases: []string{"st", "stat"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentState, err := state.Load()
		if err != nil {
			return fmt.Errorf("ошибка загрузки состояния: %w", err)
		}

		if currentState == nil || currentState.CurrentTicket == "" {
			fmt.Println(ui.Info(), "Нет активной задачи")
			fmt.Println("Используйте 'tracker start <тикет>' для начала работы")
			return nil
		}

		ticket := currentState.CurrentTicket
		duration := time.Since(currentState.StartedAt)

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			fmt.Printf("%s Задача %s не найдена на сервере. Очищаю состояние...\n",
				ui.Warning(), ui.Ticket(ticket))
			currentState.CurrentTicket = ""
			currentState.StartedAt = time.Time{}
			_ = state.Save(currentState)
			return nil
		}

		hours := int(duration.Hours())
		minutes := int(duration.Minutes()) % 60
		seconds := int(duration.Seconds()) % 60

		fmt.Printf("%s %s\n", ui.Info(), ui.Bold("Текущая задача:"))
		fmt.Printf("  Тикет:      %s\n", ui.Ticket(ticket))
		fmt.Printf("  Заголовок:  %s\n", task.Title)
		fmt.Printf("  Начата:     %s\n",
			currentState.StartedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Длительность: %02d:%02d:%02d\n",
			hours, minutes, seconds)

		if len(currentState.PausedTasks) > 0 {
			fmt.Printf("\n  %s На паузе: %d задач (см. 'tracker paused')\n",
				ui.Info(), len(currentState.PausedTasks))
		}

		return nil
	},
}

var PausedCmd = &cobra.Command{
	Use:     "paused",
	Short:   "Показать список приостановленных задач",
	Aliases: []string{"p"},
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		_ = state.CleanupOldPausedTasks()

		currentState, err := state.Load()
		if err != nil {
			return fmt.Errorf("ошибка загрузки состояния: %w", err)
		}

		if currentState == nil || len(currentState.PausedTasks) == 0 {
			fmt.Println(ui.Info(), "Нет приостановленных задач")
			return nil
		}

		sort.Slice(currentState.PausedTasks, func(i, j int) bool {
			return currentState.PausedTasks[i].PausedAt.After(
				currentState.PausedTasks[j].PausedAt)
		})

		fmt.Printf("%s %s\n", ui.Info(), ui.Bold("Приостановленные задачи:"))
		fmt.Println()

		for _, task := range currentState.PausedTasks {
			pausedAgo := time.Since(task.PausedAt)
			duration := time.Duration(task.Duration) * time.Second

			fmt.Printf("  %s  %s  (на паузе %s, отработано %s)\n",
				ui.Ticket(task.Ticket),
				formatDuration(duration),
				formatRelativeTime(pausedAgo),
				formatDuration(duration))
		}

		fmt.Println()
		fmt.Printf("Используйте 'tracker start <тикет>' для возврата к задаче\n")
		return nil
	},
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func formatRelativeTime(d time.Duration) string {
	if d < time.Minute {
		return "только что"
	}
	if d < time.Hour {
		mins := int(d.Minutes())
		return fmt.Sprintf("%d мин назад", mins)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		return fmt.Sprintf("%d ч назад", hours)
	}
	days := int(d.Hours() / 24)
	return fmt.Sprintf("%d дн назад", days)
}

func init() {
	StartCmd.Flags().BoolVar(&startStopPrevious, "stop", false,
		"Остановить предыдущую задачу вместо паузы")
}
