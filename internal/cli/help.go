package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const reset = "\033[0m"

func SetupHelp(cmd *cobra.Command) {
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		printHelp(c)
	})
}

func printHelp(cmd *cobra.Command) {
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	dim := color.New(color.Faint).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	fmt.Println()

	if cmd.Long != "" {
		fmt.Println(bold(cmd.Long))
	} else if cmd.Short != "" {
		fmt.Println(bold(cmd.Short))
	}
	fmt.Println()

	fmt.Printf("%s %s\n", green("Использование:"), cyan(cmd.Use))
	fmt.Println()

	if cmd.HasAvailableSubCommands() {
		cmds := cmd.Commands()

		authCmds := []string{"login", "logout", "me", "register"}
		configCmds := []string{"configure", "server"}
		workCmds := []string{"task", "company"}
		adminCmds := []string{"users", "role"}

		printCmdGroup := func(title string, names []string, titleColor func(...interface{}) string) {
			var filtered []*cobra.Command
			for _, c := range cmds {
				for _, name := range names {
					if c.Name() == name {
						filtered = append(filtered, c)
						break
					}
				}
			}
			if len(filtered) == 0 {
				return
			}

			fmt.Printf("  %s\n", titleColor(title))
			maxLen := 0
			for _, c := range filtered {
				if len(c.Name()) > maxLen {
					maxLen = len(c.Name())
				}
			}
			for _, c := range filtered {
				padding := strings.Repeat(" ", maxLen-len(c.Name())+4)
				fmt.Printf("    %s%s%s\n", cyan(c.Name()), padding, c.Short)
			}
			fmt.Println()
		}

		fmt.Println(bold("Команды:"))
		printCmdGroup("Авторизация:", authCmds, yellow)
		printCmdGroup("Конфигурация:", configCmds, yellow)
		printCmdGroup("Работа с данными:", workCmds, yellow)
		printCmdGroup("Администрирование:", adminCmds, yellow)

		var other []*cobra.Command
		for _, c := range cmds {
			isKnown := false
			for _, name := range append(append(append(authCmds, configCmds...), workCmds...), adminCmds...) {
				if c.Name() == name {
					isKnown = true
					break
				}
			}
			if !isKnown && !c.Hidden {
				other = append(other, c)
			}
		}
		if len(other) > 0 {
			fmt.Printf("  %s\n", dim("Другие команды:"))
			maxLen := 0
			for _, oc := range other {
				if len(oc.Name()) > maxLen {
					maxLen = len(oc.Name())
				}
			}
			for _, c := range other {
				padding := strings.Repeat(" ", maxLen-len(c.Name())+4)
				fmt.Printf("    %s%s%s\n", cyan(c.Name()), padding, c.Short)
			}
			fmt.Println()
		}
	}

	if cmd.HasAvailableFlags() {
		fmt.Println(bold("Флаги:"))
		fmt.Print(cmd.Flags().FlagUsages())
		fmt.Println()
	}

	if cmd.HasAvailableInheritedFlags() && cmd != cmd.Root() {
		fmt.Println(bold("Глобальные флаги:"))
		fmt.Print(cmd.InheritedFlags().FlagUsages())
		fmt.Println()
	}

	fmt.Printf("%s '%s' для подробной информации о команде.\n",
		dim("Используйте"),
		cyan(fmt.Sprintf("%s [command] --help", cmd.CommandPath())),
	)
	fmt.Println()
}
