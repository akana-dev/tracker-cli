package cli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

type colorSet struct {
	cyan    func(...interface{}) string
	green   func(...interface{}) string
	yellow  func(...interface{}) string
	dim     func(...interface{}) string
	bold    func(...interface{}) string
	red     func(...interface{}) string
	magenta func(...interface{}) string
}

func newColorSet() *colorSet {
	return &colorSet{
		cyan:    color.New(color.FgCyan, color.Bold).SprintFunc(),
		green:   color.New(color.FgGreen).SprintFunc(),
		yellow:  color.New(color.FgYellow).SprintFunc(),
		dim:     color.New(color.Faint).SprintFunc(),
		bold:    color.New(color.Bold).SprintFunc(),
		red:     color.New(color.FgRed, color.Bold).SprintFunc(),
		magenta: color.New(color.FgMagenta).SprintFunc(),
	}
}

type commandGroup struct {
	title string
	names []string
}

var groupsConfig = map[string][]commandGroup{
	"tracker": {
		{"Авторизация:", []string{"login", "logout", "me", "register"}},
		{"Конфигурация:", []string{"configure", "server"}},
		{"Работа с данными:", []string{"task", "company"}},
		{"Администрирование:", []string{"users", "role"}},
	},
	"task": {
		{"Создание:", []string{"add"}},
		{"Просмотр:", []string{"list", "view", "export"}},
		{"Редактирование:", []string{"edit", "assign", "close"}},
		{"Управление состоянием:", []string{"pause", "resume", "delete"}},
	},
}

func SetupHelp(cmd *cobra.Command) {
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		printHelp(c)
	})
}

func printHelp(cmd *cobra.Command) {
	colors := newColorSet()

	fmt.Println()

	if cmd.Long != "" {
		fmt.Println(colors.bold(cmd.Long))
	} else if cmd.Short != "" {
		fmt.Println(colors.bold(cmd.Short))
	}
	fmt.Println()

	fmt.Printf("%s %s\n", colors.green("Использование:"), colors.cyan(cmd.Use))
	fmt.Println()

	if cmd.HasAvailableSubCommands() {
		printCommands(cmd, colors)
	}

	if cmd.HasAvailableFlags() {
		fmt.Println(colors.bold("Флаги:"))
		fmt.Print(cmd.Flags().FlagUsages())
		fmt.Println()
	}

	if cmd.HasAvailableInheritedFlags() && cmd != cmd.Root() {
		fmt.Println(colors.bold("Глобальные флаги:"))
		fmt.Print(cmd.InheritedFlags().FlagUsages())
		fmt.Println()
	}

	fmt.Printf("%s '%s' для подробной информации о команде.\n",
		colors.dim("Используйте"),
		colors.cyan(fmt.Sprintf("%s [command] --help", cmd.CommandPath())),
	)
	fmt.Println()
}

func printCommands(cmd *cobra.Command, colors *colorSet) {
	cmds := cmd.Commands()

	if groups, ok := groupsConfig[cmd.Name()]; ok {
		printGroupedCommands(cmds, groups, colors)
	} else {
		printSimpleCommands(cmds, colors)
	}
}

func printGroupedCommands(cmds []*cobra.Command, groups []commandGroup, colors *colorSet) {
	fmt.Println(colors.bold("Команды:"))

	knownNames := make(map[string]bool)
	for _, g := range groups {
		for _, name := range g.names {
			knownNames[name] = true
		}
	}

	for _, group := range groups {
		filtered := filterCommands(cmds, group.names)
		if len(filtered) == 0 {
			continue
		}
		printCommandGroup(group.title, filtered, colors.yellow, colors.cyan)
	}

	var others []*cobra.Command
	for _, c := range cmds {
		if !c.Hidden && !knownNames[c.Name()] {
			others = append(others, c)
		}
	}
	if len(others) > 0 {
		printCommandGroup("Другие команды:", others, colors.dim, colors.cyan)
	}
}

func printSimpleCommands(cmds []*cobra.Command, colors *colorSet) {
	var visible []*cobra.Command
	for _, c := range cmds {
		if !c.Hidden {
			visible = append(visible, c)
		}
	}
	if len(visible) == 0 {
		return
	}

	fmt.Println(colors.bold("Команды:"))
	printCommandList(visible, colors.cyan)
}

func printCommandGroup(title string, cmds []*cobra.Command, titleColor, cmdColor func(...interface{}) string) {
	fmt.Printf("  %s\n", titleColor(title))
	printCommandList(cmds, cmdColor)
}

func printCommandList(cmds []*cobra.Command, cmdColor func(...interface{}) string) {
	maxLen := 0
	for _, c := range cmds {
		if len(c.Name()) > maxLen {
			maxLen = len(c.Name())
		}
	}

	for _, c := range cmds {
		padding := strings.Repeat(" ", maxLen-len(c.Name())+4)
		fmt.Printf("    %s%s%s\n", cmdColor(c.Name()), padding, c.Short)
	}
	fmt.Println()
}

func filterCommands(cmds []*cobra.Command, names []string) []*cobra.Command {
	nameSet := make(map[string]bool, len(names))
	for _, n := range names {
		nameSet[n] = true
	}

	var filtered []*cobra.Command
	for _, c := range cmds {
		if nameSet[c.Name()] && !c.Hidden {
			filtered = append(filtered, c)
		}
	}
	return filtered
}
