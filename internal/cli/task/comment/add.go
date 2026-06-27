package comment

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/client"
	"tracker/internal/input"
	"tracker/internal/service"
	"tracker/internal/ui"
)

var AddCmd = &cobra.Command{
	Use:   "add [тикет] [текст]",
	Short: "Добавить комментарий к задаче",
	Long: `Добавить комментарий к задаче. Поддерживает Markdown и упоминания @username.

Примеры:
  tracker comment add NTC-7 "Текст комментария"
  tracker comment add NTC-7 --editor
  tracker comment add NTC-7 --file comment.md
  tracker comment add NTC-7 --interactive`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticket := strings.ToUpper(args[0])

		task, err := client.GetTaskByTicket(ticket)
		if err != nil {
			return fmt.Errorf("тикет %s не найден: %w", ticket, err)
		}

		interactive, _ := cmd.Flags().GetBool("interactive")
		editor, _ := cmd.Flags().GetBool("editor")
		filePath, _ := cmd.Flags().GetString("file")

		var content string

		switch {
		case interactive:
			content, err = readInteractiveComment()
		case editor:
			content, err = readCommentFromEditor("")
		case filePath != "":
			content, err = readCommentFromFile(filePath)
		default:
			if len(args) < 2 {
				return fmt.Errorf("укажите текст комментария или используйте --editor/--file/--interactive")
			}
			content = strings.Join(args[1:], " ")
		}

		if err != nil {
			return err
		}

		content = strings.TrimSpace(content)
		if content == "" {
			return fmt.Errorf("комментарий не может быть пустым")
		}

		if err := service.ValidateComment(content); err != nil {
			return err
		}

		comment, err := client.CreateComment(task.ID, content)
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Комментарий #%d добавлен к задаче %s",
			comment.ID, ui.Ticket(ticket)))
		return nil
	},
}

func readInteractiveComment() (string, error) {
	fmt.Println()
	ui.Header("Интерактивный ввод комментария")
	fmt.Println(ui.Dim("Поддерживается Markdown. Упоминания: @username"))
	fmt.Println(ui.Dim("Введите комментарий. Для завершения — пустая строка."))
	fmt.Println()

	var lines []string
	for {
		fmt.Print("> ")
		line := input.ReadLine()

		if line == "" && len(lines) > 0 {
			break
		}

		lines = append(lines, line)
	}

	return strings.Join(lines, "\n"), nil
}

func readCommentFromEditor(initialContent string) (string, error) {
	tmpFile, err := os.CreateTemp("", "tracker-comment-*.md")
	if err != nil {
		return "", fmt.Errorf("не удалось создать временный файл: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	template := `# Напишите комментарий к задаче
# Поддерживается Markdown: **жирный**, *курсив*, ` + "`код`" + `, списки
# Упоминания: @username
# Строки, начинающиеся с #, будут удалены
# Для отмены оставьте файл пустым

`
	if initialContent != "" {
		template += initialContent
	}

	if _, err := tmpFile.WriteString(template); err != nil {
		tmpFile.Close()
		return "", err
	}
	tmpFile.Close()

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nano"
	}

	fmt.Println(ui.Dimf("Открываю редактор %s...", editor))

	editorCmd := exec.Command(editor, tmpPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr
	if err := editorCmd.Run(); err != nil {
		return "", fmt.Errorf("ошибка редактора: %w", err)
	}

	content, err := os.ReadFile(tmpPath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(content), "\n")
	var result []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			result = append(result, line)
		}
	}

	return strings.TrimSpace(strings.Join(result, "\n")), nil
}

func readCommentFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать файл: %w", err)
	}
	return string(data), nil
}

func init() {
	AddCmd.Flags().BoolP("interactive", "i", false, "Интерактивный режим ввода")
	AddCmd.Flags().BoolP("editor", "e", false, "Открыть редактор для ввода")
	AddCmd.Flags().StringP("file", "f", "", "Прочитать комментарий из файла")
}
