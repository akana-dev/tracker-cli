package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v3"

	"tracker/internal/models"
	"tracker/internal/templates"
	"tracker/internal/ui"
	"tracker/pkg/table"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Управление шаблонами задач",
	Long: `Шаблоны позволяют быстро создавать задачи с предопределёнными параметрами.

Шаблоны хранятся в ~/.tracker/templates/<name>.yaml в формате YAML.

Примеры:
  tracker template add daily-standup
  tracker template list
  tracker template show daily-standup
  tracker task from daily-standup
  tracker task from daily-standup --title "Новое название"`,
}

var templateAddCmd = &cobra.Command{
	Use:   "add [имя]",
	Short: "Создать новый шаблон",
	Long: `Создать новый шаблон. Откроет редактор ($EDITOR или nano) с примером YAML.

После сохранения файл будет проверен и добавлен в список шаблонов.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		if templates.Exists(name) {
			return fmt.Errorf("шаблон %q уже существует", name)
		}

		tmpFile, err := os.CreateTemp("", "tracker-template-*.yaml")
		if err != nil {
			return fmt.Errorf("не удалось создать временный файл: %w", err)
		}
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		if _, err := tmpFile.WriteString(templates.GenerateExample()); err != nil {
			tmpFile.Close()
			return fmt.Errorf("не удалось записать пример: %w", err)
		}
		tmpFile.Close()

		content, err := os.ReadFile(tmpPath)
		if err != nil {
			return err
		}
		content = []byte(strings.Replace(string(content), "name: example", "name: "+name, 1))
		if err := os.WriteFile(tmpPath, content, 0600); err != nil {
			return err
		}

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano"
		}

		fmt.Println(ui.Dimf("Открываю редактор %s...", editor))
		fmt.Println(ui.Dim("Заполните шаблон и сохраните файл. Для отмены закройте редактор без изменений."))
		fmt.Println()

		editorCmd := exec.Command(editor, tmpPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		if err := editorCmd.Run(); err != nil {
			return fmt.Errorf("ошибка редактора: %w", err)
		}

		resultContent, err := os.ReadFile(tmpPath)
		if err != nil {
			return err
		}

		tmpl, err := templates.ParseYAML(resultContent)
		if err != nil {
			return fmt.Errorf("ошибка парсинга YAML: %w", err)
		}

		if strings.TrimSpace(tmpl.Title) == "" {
			return fmt.Errorf("название задачи (title) обязательно")
		}

		tmpl.Name = name
		if err := templates.Save(tmpl); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Шаблон %s создан", ui.Bold(name)))
		return nil
	},
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "Показать список шаблонов",
	RunE: func(cmd *cobra.Command, args []string) error {
		allTemplates, err := templates.List()
		if err != nil {
			return err
		}

		if len(allTemplates) == 0 {
			fmt.Println(ui.Warning("Шаблоны не найдены."))
			fmt.Println(ui.Dim("Создайте шаблон: tracker template add <имя>"))
			return nil
		}

		fmt.Println()
		tbl := table.New("Имя", "Название", "Компания", "Теги")
		tbl.SetColumnWidths(map[int]int{0: 20, 1: 40, 2: 15, 3: 30})
		for _, tmpl := range allTemplates {
			tagsStr := "—"
			if len(tmpl.Tags) > 0 {
				tagsStr = strings.Join(tmpl.Tags, ", ")
			}
			company := tmpl.Company
			if company == "" {
				company = ui.Dim("—")
			}
			tbl.AddRow(
				ui.Bold(tmpl.Name),
				tmpl.Title,
				company,
				ui.Cyan(tagsStr),
			)
		}
		tbl.Render()
		fmt.Println()

		return nil
	},
}

var templateShowCmd = &cobra.Command{
	Use:   "show [имя]",
	Short: "Показать содержимое шаблона",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		tmpl, err := templates.Load(name)
		if err != nil {
			return err
		}

		fmt.Println()
		ui.Header(fmt.Sprintf("Шаблон: %s", ui.CyanBold(tmpl.Name)))
		ui.Label("Название", ui.Bold(tmpl.Title))
		if tmpl.Company != "" {
			ui.Label("Компания", ui.Cyan(tmpl.Company))
		}
		if tmpl.Assignee != "" {
			ui.Label("Исполнитель", ui.Cyan(tmpl.Assignee))
		}
		if tmpl.Solution != "" {
			ui.Label("Решение", tmpl.Solution)
		}
		if len(tmpl.Tags) > 0 {
			ui.Label("Теги", ui.Cyan(strings.Join(tmpl.Tags, ", ")))
		}
		if tmpl.Comment != "" {
			ui.Label("Комментарий", "")
			fmt.Println("    " + tmpl.Comment)
		}
		fmt.Println()

		return nil
	},
}

var templateRemoveCmd = &cobra.Command{
	Use:   "remove [имя]",
	Short: "Удалить шаблон",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := templates.Delete(name); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Шаблон %s удалён", ui.Bold(name)))
		return nil
	},
}

var templateEditCmd = &cobra.Command{
	Use:   "edit [имя]",
	Short: "Редактировать шаблон",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		tmpl, err := templates.Load(name)
		if err != nil {
			return err
		}

		tmpFile, err := os.CreateTemp("", "tracker-template-*.yaml")
		if err != nil {
			return fmt.Errorf("не удалось создать временный файл: %w", err)
		}
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		content, err := templates.MarshalYAML(tmpl)
		if err != nil {
			tmpFile.Close()
			return err
		}

		if _, err := tmpFile.Write(content); err != nil {
			tmpFile.Close()
			return err
		}
		tmpFile.Close()

		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "nano"
		}

		editorCmd := exec.Command(editor, tmpPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr
		if err := editorCmd.Run(); err != nil {
			return fmt.Errorf("ошибка редактора: %w", err)
		}

		resultContent, err := os.ReadFile(tmpPath)
		if err != nil {
			return err
		}

		updated, err := templates.ParseYAML(resultContent)
		if err != nil {
			return fmt.Errorf("ошибка парсинга YAML: %w", err)
		}

		if strings.TrimSpace(updated.Title) == "" {
			return fmt.Errorf("название задачи (title) обязательно")
		}

		updated.Name = name
		if err := templates.Save(updated); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Шаблон %s обновлён", ui.Bold(name)))
		return nil
	},
}

func init() {
	templateCmd.AddCommand(templateAddCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateRemoveCmd)
	templateCmd.AddCommand(templateEditCmd)
}

func ParseYAML(data []byte) (*models.Template, error) {
	var tmpl models.Template
	if err := yaml.Unmarshal(data, &tmpl); err != nil {
		return nil, err
	}
	return &tmpl, nil
}

func MarshalYAML(tmpl *models.Template) ([]byte, error) {
	return yaml.Marshal(tmpl)
}
