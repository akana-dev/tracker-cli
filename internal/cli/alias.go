package cli

import (
	"fmt"
	"sort"

	"tracker/internal/aliases"
	"tracker/internal/ui"
	"tracker/pkg/table"

	"github.com/spf13/cobra"
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
		value := fmt.Sprintf("%s", args[1])
		if len(args) > 2 {
			for i := 2; i < len(args); i++ {
				value += " " + args[i]
			}
		}

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
			fmt.Println()
			fmt.Println(ui.Warning("Алиасы не настроены."))
			fmt.Println(ui.Dim("Добавьте алиас: tracker alias add <имя> <команда>"))
			fmt.Println()
			return nil
		}

		names := make([]string, 0, len(allAliases))
		for name := range allAliases {
			names = append(names, name)
		}
		sort.Strings(names)

		fmt.Println()
		fmt.Println(ui.SectionHeader("Алиасы"))
		fmt.Println()
		fmt.Printf("  %s\n", ui.Dim(fmt.Sprintf("Найдено: %d", len(allAliases))))
		fmt.Println()
		fmt.Println(ui.Divider(70))
		fmt.Println()

		tbl := table.New("Алиас", "Команда")
		tbl.SetColumnWidths(map[int]int{0: 20, 1: 60})

		for _, name := range names {
			tbl.AddRow(ui.Bold(name), ui.Cyan(allAliases[name]))
		}

		tbl.Render()

		fmt.Println()
		fmt.Println(ui.Divider(70))
		fmt.Println()
		fmt.Printf("  %s %s\n", ui.Dim("Удалить:"), ui.Cyan("tracker alias remove <имя>"))
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
