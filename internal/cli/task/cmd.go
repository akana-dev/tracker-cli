package task

import (
	"github.com/spf13/cobra"

	"tracker/internal/cli/task/comment"
)

var Cmd = &cobra.Command{
	Use:   "task",
	Short: "Управление задачами",
}

func init() {
	Cmd.AddCommand(AddCmd)
	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(ViewCmd)
	Cmd.AddCommand(EditCmd)
	Cmd.AddCommand(CloseCmd)
	Cmd.AddCommand(PauseCmd)
	Cmd.AddCommand(ResumeCmd)
	Cmd.AddCommand(AssignCmd)
	Cmd.AddCommand(DeleteCmd)
	Cmd.AddCommand(ExportCmd)
	Cmd.AddCommand(FromCmd)
	Cmd.AddCommand(BulkCmd)
	Cmd.AddCommand(comment.CommentCmd)
}
