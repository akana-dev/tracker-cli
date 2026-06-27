package comment

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/ui"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [тикет] [id комментария]",
	Short: "Удалить комментарий",
	Long: `Удалить комментарий. Может удалить автор или admin.

Примеры:
  tracker comment delete NTC-7 42`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		commentID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("неверный ID комментария: %s", args[1])
		}

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		if err := client.DeleteComment(task.ID, commentID); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Комментарий #%d удалён", commentID))
		return nil
	},
}
