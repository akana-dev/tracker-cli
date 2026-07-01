package cli

import (
	"fmt"
	"strconv"
	"strings"
	"tracker/internal/client"
	"tracker/internal/ui"

	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Управление тегами задач",
}

var tagAddCmd = &cobra.Command{
	Use:   "add [имя]",
	Short: "Создать тег",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		color, _ := cmd.Flags().GetString("color")

		if color != "" && !isValidHexColor(color) {
			return fmt.Errorf("некорректный цвет: %s (формат: #RRGGBB)", color)
		}

		tag, err := client.CreateTag(name, color)
		if err != nil {
			return err
		}

		colorStr := ""
		if tag.Color != "" {
			colorStr = fmt.Sprintf(" [%s]", tag.Color)
		}
		fmt.Println(ui.Checkmark(), ui.Successf("Тег создан: %s%s", tag.Name, colorStr))
		return nil
	},
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "Показать все теги",
	RunE: func(cmd *cobra.Command, args []string) error {
		search, _ := cmd.Flags().GetString("search")
		tags, err := client.ListTags(search)
		if err != nil {
			return err
		}

		if len(tags) == 0 {
			fmt.Println(ui.Info("Теги не найдены"))
			return nil
		}

		for _, tag := range tags {
			nameDisplay := ui.TagWithColor(tag.Name, tag.Color)

			colorStr := ""
			if tag.Color != "" {
				colorStr = ui.Dimf(" (%s)", tag.Color)
			}
			fmt.Printf("  %d. %s%s\n", tag.ID, nameDisplay, colorStr)
		}

		return nil
	},
}

var tagUpdateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Обновить тег",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("некорректный ID: %s", args[0])
		}

		name, _ := cmd.Flags().GetString("name")
		color, _ := cmd.Flags().GetString("color")

		if color != "" && !isValidHexColor(color) {
			return fmt.Errorf("некорректный цвет: %s (формат: #RRGGBB)", color)
		}

		payload := map[string]interface{}{}
		if name != "" {
			payload["name"] = name
		}
		if color != "" {
			payload["color"] = color
		}

		if len(payload) == 0 {
			return fmt.Errorf("укажите хотя бы --name или --color")
		}

		tag, err := client.UpdateTag(id, payload)
		if err != nil {
			return err
		}

		colorStr := ""
		if tag.Color != "" {
			colorStr = fmt.Sprintf(" [%s]", tag.Color)
		}
		fmt.Println(ui.Checkmark(), ui.Successf("Тег обновлён: %s%s", tag.Name, colorStr))
		return nil
	},
}

var tagDeleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Удалить тег",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("некорректный ID: %s", args[0])
		}

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			fmt.Print(ui.Warning("Удалить тег? Тег будет удалён из всех задач. [y/N]: "))
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" {
				return nil
			}
		}

		if err := client.DeleteTag(id); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Success("Тег удалён"))
		return nil
	},
}

func isValidHexColor(color string) bool {
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	for _, c := range color[1:] {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}

func init() {
	tagAddCmd.Flags().StringP("color", "c", "", "Цвет тега (формат: #RRGGBB)")
	tagListCmd.Flags().StringP("search", "s", "", "Поиск по имени тега")
	tagUpdateCmd.Flags().StringP("name", "n", "", "Новое имя тега")
	tagUpdateCmd.Flags().StringP("color", "c", "", "Новый цвет тега")
	tagDeleteCmd.Flags().BoolP("force", "f", false, "Пропустить подтверждение")

	tagCmd.AddCommand(tagAddCmd)
	tagCmd.AddCommand(tagListCmd)
	tagCmd.AddCommand(tagUpdateCmd)
	tagCmd.AddCommand(tagDeleteCmd)
}
