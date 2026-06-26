package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

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
	red := color.New(color.FgRed, color.Bold).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

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

		switch cmd.Name() {
		case "tracker":
			printRootHelp(cmds, bold, cyan, green, yellow, red, magenta, dim)
		case "task":
			printTaskHelp(cmds, bold, cyan, green, yellow, dim)
		case "company":
			printCompanyHelp(cmds, bold, cyan, dim)
		case "server":
			printServerHelp(cmds, bold, cyan, dim)
		default:
			printDefaultHelp(cmds, bold, cyan, dim)
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

func printRootHelp(cmds []*cobra.Command, bold, cyan, green, yellow, red, magenta, dim func(...interface{}) string) {
	authCmds := []string{"login", "logout", "me", "register"}
	configCmds := []string{"configure", "server"}
	workCmds := []string{"task", "company"}
	adminCmds := []string{"users", "role"}

	printCmdGroup := func(title string, names []string, titleColor func(...interface{}) string) {
		var filtered []*cobra.Command
		for _, c := range cmds {
			for _, name := range names {
				if c.Name() == name && !c.Hidden {
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
	knownNames := append(append(append(authCmds, configCmds...), workCmds...), adminCmds...)
	for _, c := range cmds {
		isKnown := false
		for _, name := range knownNames {
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

func printTaskHelp(cmds []*cobra.Command, bold, cyan, green, yellow, dim func(...interface{}) string) {
	createCmds := []string{"add"}
	viewCmds := []string{"list", "view", "export"}
	editCmds := []string{"edit", "assign", "close"}
	lifecycleCmds := []string{"pause", "resume", "delete"}

	printCmdGroup := func(title string, names []string, titleColor func(...interface{}) string) {
		var filtered []*cobra.Command
		for _, c := range cmds {
			for _, name := range names {
				if c.Name() == name && !c.Hidden {
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
	printCmdGroup("Создание:", createCmds, yellow)
	printCmdGroup("Просмотр:", viewCmds, yellow)
	printCmdGroup("Редактирование:", editCmds, yellow)
	printCmdGroup("Управление состоянием:", lifecycleCmds, yellow)

	// Остальные
	var other []*cobra.Command
	knownNames := append(append(append(createCmds, viewCmds...), editCmds...), lifecycleCmds...)
	for _, c := range cmds {
		isKnown := false
		for _, name := range knownNames {
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
		for _, c := range other {
			fmt.Printf("    %s    %s\n", cyan(c.Name()), c.Short)
		}
		fmt.Println()
	}
}

func printCompanyHelp(cmds []*cobra.Command, bold, cyan, dim func(...interface{}) string) {
	fmt.Println(bold("Команды:"))
	maxLen := 0
	for _, c := range cmds {
		if !c.Hidden && len(c.Name()) > maxLen {
			maxLen = len(c.Name())
		}
	}
	for _, c := range cmds {
		if c.Hidden {
			continue
		}
		padding := strings.Repeat(" ", maxLen-len(c.Name())+4)
		fmt.Printf("  %s%s%s\n", cyan(c.Name()), padding, c.Short)
	}
	fmt.Println()
}

func printServerHelp(cmds []*cobra.Command, bold, cyan, dim func(...interface{}) string) {
	fmt.Println(bold("Команды:"))
	maxLen := 0
	for _, c := range cmds {
		if !c.Hidden && len(c.Name()) > maxLen {
			maxLen = len(c.Name())
		}
	}
	for _, c := range cmds {
		if c.Hidden {
			continue
		}
		padding := strings.Repeat(" ", maxLen-len(c.Name())+4)
		fmt.Printf("  %s%s%s\n", cyan(c.Name()), padding, c.Short)
	}
	fmt.Println()
}

func printDefaultHelp(cmds []*cobra.Command, bold, cyan, dim func(...interface{}) string) {
	if len(cmds) == 0 {
		return
	}
	fmt.Println(bold("Команды:"))
	maxLen := 0
	for _, c := range cmds {
		if !c.Hidden && len(c.Name()) > maxLen {
			maxLen = len(c.Name())
		}
	}
	for _, c := range cmds {
		if c.Hidden {
			continue
		}
		padding := strings.Repeat(" ", maxLen-len(c.Name())+4)
		fmt.Printf("  %s%s%s\n", cyan(c.Name()), padding, c.Short)
	}
	fmt.Println()
}
