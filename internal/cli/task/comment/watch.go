package comment

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/ui"
)

var WatchCmd = &cobra.Command{
	Use:   "watch [тикет]",
	Short: "Следить за новыми комментариями",
	Long: `Периодически опрашивать API и показывать новые комментарии в реальном времени.

Примеры:
  tracker comment watch NTC-7
  tracker comment watch NTC-7 --interval 10s
  tracker comment watch NTC-7 --interval 1m

Для остановки — Ctrl+C`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		interval, _ := cmd.Flags().GetDuration("interval")

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		fmt.Println(ui.Dimf("Слежение за комментариями задачи %s", ui.CyanBold(ticket)))
		fmt.Println(ui.Dimf("Интервал: %s | Для остановки — Ctrl+C", interval))
		fmt.Println()

		comments, err := client.ListComments(task.ID, 500, 0)
		if err != nil {
			return err
		}

		lastID := 0
		for _, c := range comments {
			if c.ID > lastID {
				lastID = c.ID
			}
		}

		fmt.Println(ui.Dimf("Известно комментариев: %d", len(comments)))
		fmt.Println(ui.Dim("Ожидание новых комментариев..."))
		fmt.Println()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigCh
			fmt.Println()
			fmt.Println(ui.Warning("\nОстановка слежения..."))
			cancel()
		}()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return nil
			case <-ticker.C:
				newComments, err := client.ListComments(task.ID, 500, 0)
				if err != nil {
					fmt.Println(ui.Warningf("Ошибка опроса: %v", err))
					continue
				}

				for _, c := range newComments {
					if c.ID > lastID {
						fmt.Println("─────────────────────────────────────────────────────────")
						fmt.Printf("%s ", ui.Success("[NEW]"))
						fmt.Print(RenderComment(c, 0))
						lastID = c.ID
					}
				}
			}
		}
	},
}

func init() {
	WatchCmd.Flags().DurationP("interval", "n", 30*time.Second, "Интервал опроса (например: 10s, 1m)")
}
