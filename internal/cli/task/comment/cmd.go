package comment

import (
	"github.com/spf13/cobra"
)

var CommentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Управление комментариями задач",
	Long: `Команды для работы с комментариями к задачам.

Поддерживается Markdown и упоминания @username.

Примеры:
  tracker comment list NTC-7
  tracker comment add NTC-7 "Текст комментария"
  tracker comment add NTC-7 --editor
  tracker comment edit NTC-7 42 "Новый текст"
  tracker comment delete NTC-7 42
  tracker comment watch NTC-7`,
}

func init() {
	CommentCmd.AddCommand(ListCmd)
	CommentCmd.AddCommand(AddCmd)
	CommentCmd.AddCommand(EditCmd)
	CommentCmd.AddCommand(DeleteCmd)
	CommentCmd.AddCommand(WatchCmd)
}
