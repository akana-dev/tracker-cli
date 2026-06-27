package comment

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/service"
	"tracker/internal/ui"
)

var ListCmd = &cobra.Command{
	Use:   "list [тикет]",
	Short: "Показать список комментариев задачи",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		all, _ := cmd.Flags().GetBool("all")
		page, _ := cmd.Flags().GetInt("page")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		if cmd.Flags().Changed("page") && cmd.Flags().Changed("offset") {
			return fmt.Errorf("нельзя использовать --page и --offset одновременно")
		}
		if page < 1 {
			return fmt.Errorf("--page должен быть >= 1")
		}
		if offset < 0 {
			return fmt.Errorf("--offset должен быть >= 0")
		}
		if limit < 0 {
			return fmt.Errorf("--limit должен быть >= 0")
		}

		if all {
			limit = 0
			offset = 0
		} else {
			if !cmd.Flags().Changed("limit") {
				limit = service.DefaultPageSize
			}
			if cmd.Flags().Changed("page") && page > 1 {
				offset = (page - 1) * limit
			}
		}

		comments, err := client.ListComments(task.ID, limit, offset)
		if err != nil {
			return err
		}

		if len(comments) == 0 {
			fmt.Println(ui.Warning("Комментариев нет."))
			return nil
		}

		fmt.Println()
		ui.Header(fmt.Sprintf("Комментарии задачи %s", ui.CyanBold(ticket)))
		fmt.Println()

		if limit > 0 {
			startIdx := offset + 1
			endIdx := offset + len(comments)
			fmt.Printf("  %s\n",
				ui.Dim(fmt.Sprintf("Показано: %d-%d", startIdx, endIdx)))

			if len(comments) == limit {
				nextPage := page + 1
				fmt.Printf("  %s\n",
					ui.Dim(fmt.Sprintf("Возможно есть ещё. Используйте %s",
						ui.Cyan(fmt.Sprintf("--page %d", nextPage)))))
			}
			fmt.Println()
		}

		for i, c := range comments {
			fmt.Println("─────────────────────────────────────────────────────────")
			fmt.Print(RenderComment(c, offset+i+1))
		}
		fmt.Println("─────────────────────────────────────────────────────────")
		fmt.Println()

		return nil
	},
}

func init() {
	ListCmd.Flags().Bool("all", false, "Показать все комментарии (без пагинации)")
	ListCmd.Flags().Int("page", 1, "Номер страницы")
	ListCmd.Flags().Int("limit", service.DefaultPageSize, "Количество комментариев на странице")
	ListCmd.Flags().Int("offset", 0, "Смещение от начала")
}
