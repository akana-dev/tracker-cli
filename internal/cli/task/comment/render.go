package comment

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/glamour"

	"tracker/internal/models"
	"tracker/internal/ui"
)

var mentionRegex = regexp.MustCompile(`@(\w+)`)

func RenderComment(c models.Comment, index int) string {
	var sb strings.Builder

	sb.WriteString(renderCommentHeader(c, index))
	sb.WriteString("\n")

	sb.WriteString(renderCommentContent(c.Content))
	sb.WriteString("\n")

	sb.WriteString(renderCommentActions(c))
	sb.WriteString("\n")

	return sb.String()
}

func renderCommentHeader(c models.Comment, index int) string {
	user := c.User.GetDisplayName()
	timeStr := c.CreatedAt.Local().Format("02.01.2006 15:04")

	edited := ""
	if c.IsEdited {
		edited = " " + ui.Dim("(ред.)")
	}

	return fmt.Sprintf("%s %s • %s%s",
		ui.Dim(fmt.Sprintf("#%d", index)),
		ui.Bold(user),
		ui.Dim(timeStr),
		edited,
	)
}

func renderCommentContent(content string) string {
	renderer, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return "    " + content
	}

	rendered, err := renderer.Render(content)
	if err != nil {
		return "    " + content
	}

	rendered = mentionRegex.ReplaceAllStringFunc(rendered, func(match string) string {
		return ui.Cyan(match)
	})

	lines := strings.Split(rendered, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = "    " + line
		}
	}

	return strings.Join(lines, "\n")
}

func renderCommentActions(c models.Comment) string {
	var actions []string

	if c.CanEdit {
		actions = append(actions, ui.Dim("[можно редактировать]"))
	}
	if c.CanDelete {
		actions = append(actions, ui.Dim("[можно удалить]"))
	}

	if len(actions) == 0 {
		return ""
	}

	return "    " + strings.Join(actions, " ")
}

func RenderCommentsList(comments []models.Comment, taskTicket string) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(ui.InfoBold(fmt.Sprintf("Комментарии задачи %s (%d):", taskTicket, len(comments))))
	sb.WriteString("\n")

	for i, c := range comments {
		sb.WriteString("─────────────────────────────────────────────────────────\n")
		sb.WriteString(RenderComment(c, i+1))
	}

	sb.WriteString("─────────────────────────────────────────────────────────\n")
	sb.WriteString("\n")

	return sb.String()
}

func FormatCommentForView(comments []models.Comment) string {
	if len(comments) == 0 {
		return ui.Dim("Нет комментариев")
	}

	var sb strings.Builder

	for i, c := range comments {
		user := c.User.Username
		timeStr := c.CreatedAt.Local().Format("02.01 15:04")

		firstLine := strings.Split(c.Content, "\n")[0]
		if len(firstLine) > 60 {
			firstLine = firstLine[:57] + "..."
		}

		edited := ""
		if c.IsEdited {
			edited = " (ред.)"
		}

		sb.WriteString(fmt.Sprintf("  %s %s • %s%s\n",
			ui.Dim(fmt.Sprintf("#%d", i+1)),
			ui.Cyan(user),
			ui.Dim(timeStr),
			edited,
		))
		sb.WriteString(fmt.Sprintf("      %s\n", firstLine))
	}

	return sb.String()
}
