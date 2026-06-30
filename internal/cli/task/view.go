package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"tracker/internal/cli/task/comment"
	"tracker/internal/client"
	"tracker/internal/config"
	"tracker/internal/service"
	"tracker/internal/tags"
	"tracker/internal/ui"
)

var ViewCmd = &cobra.Command{
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

		taskTags, _ := tags.Get(ticket)
		if len(taskTags) > 0 {
			ui.Label("Теги", ui.Cyan(strings.Join(taskTags, ", ")))
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

		noComments, _ := cmd.Flags().GetBool("no-comments")
		if !noComments {
			fmt.Println()
			ui.Header(fmt.Sprintf("Комментарии (%d)", len(task.Comments)))
			fmt.Println()

			if len(task.Comments) == 0 {
				fmt.Println("    " + ui.Dim("Нет комментариев"))
			} else {
				fmt.Print(comment.FormatCommentForView(task.Comments))
			}
		}

		fmt.Println()

		fmt.Println(ui.Dim("Команды для работы с задачей:"))
		fmt.Printf("  %s  %s\n", ui.Cyan("edit"), ui.Dim("Редактировать задачу"))
		fmt.Printf("  %s  %s\n", ui.Cyan("pause"), ui.Dim("Поставить на паузу"))
		fmt.Printf("  %s  %s\n", ui.Cyan("resume"), ui.Dim("Возобновить"))
		fmt.Printf("  %s  %s\n", ui.Cyan("close"), ui.Dim("Закрыть задачу"))
		fmt.Printf("  %s  %s\n", ui.Cyan("assign"), ui.Dim("Назначить исполнителя"))
		fmt.Printf("  %s  %s\n", ui.Cyan("delete"), ui.Dim("Удалить задачу"))
		fmt.Printf("  %s  %s\n", ui.Cyan("comment list"), ui.Dim("Показать комментарии"))
		fmt.Printf("  %s  %s\n", ui.Cyan("comment add"), ui.Dim("Добавить комментарий"))
		fmt.Println()

		return nil
	},
}

func init() {
	ViewCmd.Flags().BoolP("no-comments", "N", false, "Не показывать комментарии")
}
