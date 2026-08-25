package task

import (
	"fmt"
	"strings"
	"time"

	"tracker/internal/app"
	"tracker/internal/cli/task/comment"
	"tracker/internal/config"
	"tracker/internal/service"
	"tracker/internal/ui"

	"github.com/spf13/cobra"
)

var ViewCmd = &cobra.Command{
	Use:   "view [тикет]",
	Short: "Подробная информация о задаче",
	Long:  "Показать полную информацию о задаче с сессиями, комментарием и метаданными",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		c := app.GetClient()
		task, err := c.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		server, _ := config.GetCurrentServer()
		serverName := "—"
		if server != nil {
			serverName = server.Name
		}

		fmt.Println()
		fmt.Println(ui.SectionHeader(fmt.Sprintf("%s  %s", ui.CyanBold(task.Ticket), service.FormatStatus(*task))))
		fmt.Println()

		fmt.Println(ui.SectionHeader("Основная информация"))
		fmt.Println(ui.KeyValue("Название", ui.Bold(task.Title), 16))
		fmt.Println(ui.KeyValue("Компания", ui.Cyan(task.CompanyName), 16))
		fmt.Println(ui.KeyValue("Сервер", ui.Dim(serverName), 16))
		fmt.Println(ui.KeyValue("Создатель", ui.Cyan(task.GetOwnerDisplay()), 16))

		if task.IsAssignedToSomeone() {
			fmt.Println(ui.KeyValue("Исполнитель", ui.Cyan(task.GetAssigneeDisplay()), 16))
		} else {
			fmt.Println(ui.KeyValue("Исполнитель", ui.Cyan(task.GetAssigneeDisplay())+ui.Dim(" (создатель)"), 16))
		}

		if len(task.Tags) > 0 {
			tagInfos := make([]ui.TagInfo, 0, len(task.Tags))
			for _, tag := range task.Tags {
				tagInfos = append(tagInfos, ui.TagInfo{
					Name:  tag.Name,
					Color: tag.Color,
				})
			}
			fmt.Println(ui.KeyValue("Теги", ui.TagsDisplay(tagInfos), 16))
		}

		fmt.Println()

		fmt.Println(ui.SectionHeader("Время"))
		fmt.Println(ui.KeyValue("Начало", task.StartTime.Local().Format("02.01.2006 15:04"), 16))

		if task.IsClosed() {
			fmt.Println(ui.KeyValue("Окончание", task.EndTime.Local().Format("02.01.2006 15:04"), 16))
			duration := task.EndTime.Sub(task.StartTime.Time)
			fmt.Println(ui.KeyValue("Длительность", service.FormatDuration(duration), 16))
		} else {
			fmt.Println(ui.KeyValue("Окончание", ui.Warning("не закрыта"), 16))
		}

		if task.IsPaused() {
			fmt.Println(ui.KeyValue("На паузе с", ui.Warning(task.PausedAt.Local().Format("02.01.2006 15:04")), 16))
		}

		totalHours := service.CalculateTaskHours(*task)
		fmt.Println(ui.KeyValue("Отработано", ui.Cyan(fmt.Sprintf("%.1f ч.", totalHours)), 16))
		fmt.Println()

		fmt.Println(ui.SectionHeader("Статус и описание"))
		solution := "—"
		if task.Solution != nil && *task.Solution != "" {
			solution = *task.Solution
		}
		fmt.Println(ui.KeyValue("Решение", service.FormatStatus(*task)+" "+solution, 16))

		if task.Comment != nil && *task.Comment != "" {
			fmt.Println(ui.KeyValue("Комментарий", "", 16))
			service.PrintIndented(*task.Comment, "    ")
		} else {
			fmt.Println(ui.KeyValue("Комментарий", ui.Dim("—"), 16))
		}

		fmt.Println()

		fmt.Println(ui.SectionHeader(fmt.Sprintf("Сессии (%d)", len(task.Sessions))))
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
						elapsed := time.Since(startLocal)
						fmt.Printf("%s — %s\n",
							startStr,
							ui.InProgress(fmt.Sprintf("в работе (%s)", service.FormatDuration(elapsed))),
						)
					}
				}
			}
		}

		fmt.Println()

		fmt.Println(ui.SectionHeader("Права доступа"))
		if task.CanEdit {
			fmt.Println(ui.KeyValue("Редактирование", ui.StatusOK(), 16))
		} else {
			fmt.Println(ui.KeyValue("Редактирование", ui.StatusNo(), 16))
		}
		if task.CanDelete {
			fmt.Println(ui.KeyValue("Удаление", ui.StatusOK(), 16))
		} else {
			fmt.Println(ui.KeyValue("Удаление", ui.StatusNo(), 16))
		}

		noComments, _ := cmd.Flags().GetBool("no-comments")
		if !noComments {
			fmt.Println()
			fmt.Println(ui.SectionHeader(fmt.Sprintf("Комментарии (%d)", len(task.Comments))))
			fmt.Println()
			if len(task.Comments) == 0 {
				fmt.Println("    " + ui.Dim("Нет комментариев"))
			} else {
				fmt.Print(comment.FormatCommentForView(task.Comments))
			}
		}

		fmt.Println()
		fmt.Println(ui.Divider(80))
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
