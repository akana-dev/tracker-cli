package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/tags"
	"tracker/internal/ui"
	"tracker/pkg/table"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Управление тегами задач",
	Long: `Теги позволяют классифицировать задачи по произвольным категориям.

Примеры:
  tracker tag add NTC-7 backend refactoring
  tracker tag remove NTC-7 backend
  tracker tag list
  tracker task list --tag backend`,
}

var tagAddCmd = &cobra.Command{
	Use:   "add [тикет] [теги...]",
	Short: "Добавить теги к задаче",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		tagList := args[1:]

		if err := tags.Add(ticket, tagList); err != nil {
			return err
		}

		allTags, err := tags.Get(ticket)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Теги задачи %s: %s",
			ui.Ticket(ticket), ui.Cyan(strings.Join(allTags, ", "))))
		return nil
	},
}

var tagRemoveCmd = &cobra.Command{
	Use:   "remove [тикет] [теги...]",
	Short: "Удалить теги у задачи",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])
		tagList := args[1:]

		if err := tags.Remove(ticket, tagList); err != nil {
			return err
		}

		allTags, err := tags.Get(ticket)
		if err != nil {
			return err
		}

		if len(allTags) == 0 {
			fmt.Println(ui.Checkmark(), ui.Successf("Все теги удалены у задачи %s", ui.Ticket(ticket)))
		} else {
			fmt.Println(ui.Checkmark(), ui.Successf("Теги задачи %s: %s",
				ui.Ticket(ticket), ui.Cyan(strings.Join(allTags, ", "))))
		}
		return nil
	},
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "Показать все теги",
	RunE: func(cmd *cobra.Command, args []string) error {
		allTags, err := tags.ListAll()
		if err != nil {
			return err
		}

		if len(allTags) == 0 {
			fmt.Println(ui.Warning("Теги не найдены."))
			fmt.Println(ui.Dim("Добавьте теги: tracker tag add <тикет> <тег>"))
			return nil
		}

		ticketsByTag, err := tags.ListByTicket()
		if err != nil {
			return err
		}

		tagCount := make(map[string]int)
		for _, taskTags := range ticketsByTag {
			for _, tag := range taskTags {
				tagCount[tag]++
			}
		}

		fmt.Println()
		tbl := table.New("Тег", "Использований")
		tbl.SetColumnWidths(map[int]int{0: 30, 1: 15})
		for _, tag := range allTags {
			tbl.AddRow(
				ui.Cyan(tag),
				fmt.Sprintf("%d", tagCount[tag]),
			)
		}
		tbl.Render()
		fmt.Println()

		return nil
	},
}

var tagShowCmd = &cobra.Command{
	Use:   "show [тикет]",
	Short: "Показать теги задачи",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])

		taskTags, err := tags.Get(ticket)
		if err != nil {
			return err
		}

		if len(taskTags) == 0 {
			fmt.Println(ui.Warningf("У задачи %s нет тегов", ui.Ticket(ticket)))
			return nil
		}

		fmt.Println()
		ui.Header(fmt.Sprintf("Теги задачи %s", ui.CyanBold(ticket)))
		for _, tag := range taskTags {
			fmt.Printf("  %s %s\n", ui.Bullet(), ui.Cyan(tag))
		}
		fmt.Println()

		return nil
	},
}

func init() {
	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagRemoveCmd)
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagShowCmd)
}
