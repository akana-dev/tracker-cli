package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func Card(title string, content string) string {
	t := GetTheme()
	if !ShouldUseColors() {
		return fmt.Sprintf("┌─ %s ─────────────────────\n%s\n└────────────────────────────", title, content)
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(t.Primary))

	cardStyle := t.styles.card.Copy().
		BorderForeground(lipgloss.Color(t.CardBorder))

	titleBar := lipgloss.JoinHorizontal(lipgloss.Left,
		titleStyle.Render(title),
	)

	return cardStyle.Render(titleBar + "\n" + content)
}

func Badge(text string) string {
	t := GetTheme()
	if !ShouldUseColors() {
		return fmt.Sprintf("[%s]", text)
	}
	return t.styles.badge.Render(text)
}

func BadgeWithColor(text, bgColor, fgColor string) string {
	if !ShouldUseColors() {
		return fmt.Sprintf("[%s]", text)
	}
	style := lipgloss.NewStyle().
		Padding(0, 1).
		Background(lipgloss.Color(bgColor)).
		Foreground(lipgloss.Color(fgColor))
	return style.Render(text)
}

func Divider(width int) string {
	t := GetTheme()
	if !ShouldUseColors() {
		return strings.Repeat("-", width)
	}
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Muted))
	return style.Render(strings.Repeat("─", width))
}

func SectionHeader(title string) string {
	t := GetTheme()
	if !ShouldUseColors() {
		return fmt.Sprintf("\n%s\n%s", title, strings.Repeat("-", len(title)))
	}
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(t.Primary)).
		MarginBottom(1)
	return style.Render(title)
}

func KeyValue(key, value string, keyWidth int) string {
	t := GetTheme()
	if !ShouldUseColors() {
		return fmt.Sprintf("  %-*s %s", keyWidth, key+":", value)
	}
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(t.Muted)).
		Width(keyWidth)
	return fmt.Sprintf("  %s %s", keyStyle.Render(key+":"), value)
}

func StatusBadge(status string, active bool) string {
	t := GetTheme()
	if !ShouldUseColors() {
		if active {
			return fmt.Sprintf("[%s]", status)
		}
		return status
	}

	var bgColor, fgColor string
	if active {
		bgColor = t.Success
		fgColor = "#000000"
	} else {
		bgColor = t.Muted
		fgColor = "#FFFFFF"
	}

	return BadgeWithColor(status, bgColor, fgColor)
}

func ProgressBar(current, total int, width int) string {
	if total == 0 {
		return strings.Repeat(" ", width)
	}

	t := GetTheme()
	filled := (current * width) / total
	empty := width - filled

	if !ShouldUseColors() {
		return fmt.Sprintf("[%s%s] %d/%d",
			strings.Repeat("=", filled),
			strings.Repeat(" ", empty),
			current, total)
	}

	filledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Success))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted))

	return fmt.Sprintf("%s%s %d/%d",
		filledStyle.Render(strings.Repeat("█", filled)),
		emptyStyle.Render(strings.Repeat("░", empty)),
		current, total)
}

func Indent(text string, spaces int) string {
	prefix := strings.Repeat(" ", spaces)
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if line != "" {
			lines[i] = prefix + line
		}
	}
	return strings.Join(lines, "\n")
}

func Box(text string) string {
	t := GetTheme()
	if !ShouldUseColors() {
		lines := strings.Split(text, "\n")
		maxLen := 0
		for _, line := range lines {
			if len(line) > maxLen {
				maxLen = len(line)
			}
		}
		var result strings.Builder
		result.WriteString("┌")
		result.WriteString(strings.Repeat("─", maxLen+2))
		result.WriteString("┐\n")
		for _, line := range lines {
			result.WriteString("│ ")
			result.WriteString(line)
			result.WriteString(strings.Repeat(" ", maxLen-len(line)))
			result.WriteString(" │\n")
		}
		result.WriteString("└")
		result.WriteString(strings.Repeat("─", maxLen+2))
		result.WriteString("┘")
		return result.String()
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(t.CardBorder)).
		Padding(0, 1)
	return style.Render(text)
}
