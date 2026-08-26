package cli

import (
	"tracker/internal/cli/task"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func copyCommand(src *cobra.Command) *cobra.Command {
	dst := &cobra.Command{
		Use:     src.Use,
		Short:   src.Short,
		Long:    src.Long,
		Example: src.Example,
		Args:    src.Args,
		RunE:    src.RunE,
		Run:     src.Run,
	}

	if len(src.Aliases) > 0 {
		dst.Aliases = make([]string, len(src.Aliases))
		copy(dst.Aliases, src.Aliases)
	}

	src.Flags().VisitAll(func(f *pflag.Flag) {
		dst.Flags().AddFlag(f)
	})

	for _, sub := range src.Commands() {
		dst.AddCommand(copyCommand(sub))
	}

	return dst
}

var (
	showCmd        *cobra.Command
	closeCmd       *cobra.Command
	pauseCmd       *cobra.Command
	resumeCmd      *cobra.Command
	assignCmd      *cobra.Command
	editCmd        *cobra.Command
	deleteCmd      *cobra.Command
	addCmd         *cobra.Command
	lsCmd          *cobra.Command
	useCmd         *cobra.Command
	exportTasksCmd *cobra.Command
)

func init() {
	showCmd = copyCommand(task.ViewCmd)
	showCmd.Use = "show [тикет]"
	showCmd.Aliases = []string{"view", "info"}

	closeCmd = copyCommand(task.CloseCmd)
	closeCmd.Use = "close [тикет]"
	closeCmd.Aliases = []string{"done", "finish"}

	pauseCmd = copyCommand(task.PauseCmd)
	pauseCmd.Use = "pause [тикет]"
	pauseCmd.Aliases = []string{"hold"}

	resumeCmd = copyCommand(task.ResumeCmd)
	resumeCmd.Use = "resume [тикет]"
	resumeCmd.Aliases = []string{"continue"}

	assignCmd = copyCommand(task.AssignCmd)
	assignCmd.Use = "assign [тикет] [пользователь]"
	assignCmd.Aliases = []string{"a"}

	editCmd = copyCommand(task.EditCmd)
	editCmd.Use = "edit [тикет]"
	editCmd.Aliases = []string{"update", "modify"}

	deleteCmd = copyCommand(task.DeleteCmd)
	deleteCmd.Use = "delete [тикет]"
	deleteCmd.Aliases = []string{"del", "rm", "remove"}

	addCmd = copyCommand(task.AddCmd)
	addCmd.Use = "add [название...]"
	addCmd.Aliases = []string{"new", "create"}

	lsCmd = copyCommand(task.ListCmd)
	lsCmd.Use = "ls"
	lsCmd.Aliases = []string{"list", "l"}

	useCmd = copyCommand(task.FromCmd)
	useCmd.Use = "use [имя_шаблона]"
	useCmd.Aliases = []string{"from", "template"}

	exportTasksCmd = copyCommand(task.ExportCmd)
	exportTasksCmd.Use = "export [формат]"
	exportTasksCmd.Aliases = []string{"exp"}

	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(closeCmd)
	rootCmd.AddCommand(pauseCmd)
	rootCmd.AddCommand(resumeCmd)
	rootCmd.AddCommand(assignCmd)
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(lsCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(exportTasksCmd)
}
