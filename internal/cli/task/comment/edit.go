package comment

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/service"
	"tracker/internal/ui"
)

var EditCmd = &cobra.Command{
	Use:   "edit [тикет] [id комментария] [новый текст]",
	Short: "Редактировать свой комментарий",
	Long: `Редактировать свой комментарий. Только автор может редактировать.

Примеры:
  tracker comment edit NTC-7 42 "Новый текст"
  tracker comment edit NTC-7 42 --editor
  tracker comment edit NTC-7 42 --file new_content.md`,
	Args: cobra.MinimumNArgs(2),
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

		editor, _ := cmd.Flags().GetBool("editor")
		filePath, _ := cmd.Flags().GetString("file")

		var content string

		switch {
		case editor:
			comments, err := client.ListComments(task.ID, 500, 0)
			if err != nil {
				return err
			}

			var currentContent string
			for _, c := range comments {
				if c.ID == commentID {
					currentContent = c.Content
					break
				}
			}

			if currentContent == "" {
				return fmt.Errorf("комментарий #%d не найден", commentID)
			}

			content, err = readCommentFromEditor(currentContent)
		case filePath != "":
			content, err = readCommentFromFile(filePath)
		default:
			if len(args) < 3 {
				return fmt.Errorf("укажите новый текст или используйте --editor/--file")
			}
			content = strings.Join(args[2:], " ")
		}

		if err != nil {
			return err
		}

		content = strings.TrimSpace(content)
		if content == "" {
			return fmt.Errorf("комментарий не может быть пустым")
		}

		if err := service.ValidateComment(content); err != nil {
			return err
		}

		_, err = client.UpdateComment(task.ID, commentID, content)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Комментарий #%d обновлён", commentID))
		return nil
	},
}

func init() {
	EditCmd.Flags().BoolP("editor", "e", false, "Открыть редактор")
	EditCmd.Flags().StringP("file", "f", "", "Прочитать из файла")
}
