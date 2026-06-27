package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/aliases"
	"tracker/internal/ui"
	"tracker/pkg/table"
)

var aliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Управление алиасами команд",
	Long: `Алиасы позволяют создавать короткие команды для часто используемых действий.

Примеры:
  tracker alias add ll "task list --today"
  tracker alias add w "task list --week"
  tracker alias add st "status"
  tracker ll   # выполнит: tracker task list --today`,
}

var aliasAddCmd = &cobra.Command{
	Use:   "add [имя] [команда]",
	Short: "Добавить новый алиас",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		value := strings.Join(args[1:], " ")

		if err := aliases.Add(name, value); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Алиас %s добавлен: %s",
			ui.Bold(name), ui.Cyan(value)))
		return nil
	},
}

var aliasListCmd = &cobra.Command{
	Use:   "list",
	Short: "Показать список алиасов",
	RunE: func(cmd *cobra.Command, args []string) error {
		allAliases, err := aliases.List()
		if err != nil {
			return err
		}

		if len(allAliases) == 0 {
			fmt.Println(ui.Warning("Алиасы не настроены."))
			fmt.Println(ui.Dim("Добавьте алиас: tracker alias add <имя> <команда>"))
			return nil
		}

		names := make([]string, 0, len(allAliases))
		for name := range allAliases {
			names = append(names, name)
		}
		sort.Strings(names)

		fmt.Println()
		tbl := table.New("Алиас", "Команда")
		tbl.SetColumnWidths(map[int]int{0: 20, 1: 60})
		for _, name := range names {
			tbl.AddRow(ui.Bold(name), ui.Cyan(allAliases[name]))
		}
		tbl.Render()
		fmt.Println()

		return nil
	},
}

var aliasRemoveCmd = &cobra.Command{
	Use:   "remove [имя]",
	Short: "Удалить алиас",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := aliases.Remove(name); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Алиас %s удалён", ui.Bold(name)))
		return nil
	},
}

func init() {
	aliasCmd.AddCommand(aliasAddCmd)
	aliasCmd.AddCommand(aliasListCmd)
	aliasCmd.AddCommand(aliasRemoveCmd)
}
