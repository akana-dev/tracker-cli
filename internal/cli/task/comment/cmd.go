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
  tracker comment list TEST-7
  tracker comment add TEST-7 "Текст комментария"
  tracker comment add TEST-7 --editor
  tracker comment edit TEST-7 42 "Новый текст"
  tracker comment delete TEST-7 42
  tracker comment watch TEST-7`,
}

func init() {
	CommentCmd.AddCommand(ListCmd)
	CommentCmd.AddCommand(AddCmd)
	CommentCmd.AddCommand(EditCmd)
	CommentCmd.AddCommand(DeleteCmd)
	CommentCmd.AddCommand(WatchCmd)
}
