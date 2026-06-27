package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/cli/task"
	"tracker/internal/config"
	"tracker/internal/installer"
	"tracker/internal/updater"
	"tracker/internal/version"
)

const (
	githubOwner = "akana-dev"
	githubRepo  = "tracker-cli"
)

var ErrHelp = errors.New("help requested")

var silentCommands = map[string]bool{
	"login":     true,
	"register":  true,
	"configure": true,

	// Служебные команды
	"help":       true,
	"completion": true,
	"update":     true,

	"install":   true,
	"uninstall": true,
}

var (
	installFlag   bool
	uninstallFlag bool
	versionFlag   bool
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
		if versionFlag {
			fmt.Println(version.String())
			return nil
		}

		return cmd.Help()
	},
}

func shouldSkipUpdateCheck() bool {
	if len(os.Args) < 2 {
		return true
	}

	for _, arg := range os.Args[1:] {
		if arg == "--help" || arg == "-h" || arg == "--version" || arg == "-v" {
			return true
		}

		if strings.HasPrefix(arg, "-") {
			continue
		}

		if silentCommands[arg] {
			return true
		}

		break
	}

	return false
}

func isPublicCommand(cmd *cobra.Command) bool {
	publicCmds := map[string]bool{
		"tracker": true,
		"login":   true, "register": true, "configure": true,
		"server": true, "help": true, "completion": true,
		"alias": true, "tag": true, "template": true, "config": true,
		"export": true, "update": true,
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
	if !shouldSkipUpdateCheck() {
		go updater.CheckAndNotify(githubOwner, githubRepo, updater.DefaultCheckInterval)
	}

	SetupHelp(rootCmd)
	setupSubCommandsHelp(rootCmd)

	err := rootCmd.Execute()

	updater.WaitForCheck(1)

	return err
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
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Показать версию")

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
	rootCmd.AddCommand(updateCmd)
}
