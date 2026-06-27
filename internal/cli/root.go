package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"tracker/internal/cli/task"
	"tracker/internal/config"
	"tracker/internal/installer"
)

var ErrHelp = errors.New("help requested")

var (
	installFlag   bool
	uninstallFlag bool
)

var rootCmd = &cobra.Command{
	Use:           "tracker",
	Short:         "Трекер времени задач",
	Long:          "Трекер времени задач с поддержкой нескольких серверов",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if isPublicCommand(cmd) {
			return nil
		}

		token := config.LoadToken()
		if token == "" {
			return fmt.Errorf("не авторизованы. Выполните: tracker login")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if installFlag {
			return installer.Install()
		}
		if uninstallFlag {
			return installer.Uninstall()
		}

		return cmd.Help()
	},
}

func isPublicCommand(cmd *cobra.Command) bool {
	publicCmds := map[string]bool{
		"tracker": true,
		"login":   true, "register": true, "configure": true,
		"server": true, "help": true, "completion": true,
		"alias": true, "tag": true, "template": true, "config": true,
		"export": true,
	}

	if publicCmds[cmd.Name()] {
		return true
	}

	current := cmd.Parent()
	for current != nil {
		if publicCmds[current.Name()] {
			return true
		}
		current = current.Parent()
	}

	return false
}

func Execute() error {
	SetupHelp(rootCmd)
	setupSubCommandsHelp(rootCmd)
	return rootCmd.Execute()
}

func setupSubCommandsHelp(cmd *cobra.Command) {
	for _, sub := range cmd.Commands() {
		SetupHelp(sub)
		setupSubCommandsHelp(sub)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&installFlag, "install", false, "Установить tracker в систему (добавить в PATH)")
	rootCmd.Flags().BoolVar(&uninstallFlag, "uninstall", false, "Удалить tracker из системы")

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(meCmd)
	rootCmd.AddCommand(registerCmd)

	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(configCmd)

	rootCmd.AddCommand(task.Cmd)
	rootCmd.AddCommand(companyCmd)
	rootCmd.AddCommand(exportCmd)

	rootCmd.AddCommand(aliasCmd)
	rootCmd.AddCommand(tagCmd)
	rootCmd.AddCommand(templateCmd)
}
