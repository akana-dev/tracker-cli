package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

// Цветовая палитра
var (
	PrimaryColor = lipgloss.Color("#7D56F4")
	SuccessColor = lipgloss.Color("#04B575")
	WarningColor = lipgloss.Color("#FFA500")
	ErrorColor   = lipgloss.Color("#FF4444")
	CyanColor    = lipgloss.Color("#67B7A1")
	YellowColor  = lipgloss.Color("#FFD700")
	DimColor     = lipgloss.Color("#6272A4")
	InfoColor    = lipgloss.Color("#8BE9FD")
)

// TermWidth возвращает ширину терминала
func TermWidth() int {
	width, _, err := term.GetSize(0)
	if err != nil || width <= 0 {
		return 120
	}
	return width
}

// TermHeight возвращает высоту терминала
func TermHeight() int {
	_, height, err := term.GetSize(0)
	if err != nil || height <= 0 {
		return 40
	}
	return height
}

// Базовые стили
var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginBottom(1)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor)

	InfoStyle = lipgloss.NewStyle().
			Foreground(InfoColor)

	DimStyle = lipgloss.NewStyle().
			Foreground(DimColor)

	BoldStyle = lipgloss.NewStyle().Bold(true)

	TicketStyle = lipgloss.NewStyle().
			Foreground(CyanColor).
			Bold(true).
			Padding(0, 1)

	AdminRoleStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	ManagerRoleStyle = lipgloss.NewStyle().
				Foreground(YellowColor).
				Bold(true)

	UserRoleStyle = lipgloss.NewStyle().
			Foreground(SuccessColor)

	InProgressStyle = lipgloss.NewStyle().
			Foreground(SuccessColor)

	PausedStyle = lipgloss.NewStyle().
			Foreground(WarningColor)

	ClosedStyle = lipgloss.NewStyle().
			Foreground(DimColor)
)

// BoxStyle возвращает стиль с рамкой, адаптированный под ширину терминала
func BoxStyle() lipgloss.Style {
	width := TermWidth() - 4 // отступы
	if width < 40 {
		width = 40
	}
	if width > 80 {
		width = 80
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor).
		Padding(1, 2).
		Width(width).
		MarginBottom(1)
}

// Функции форматирования
func Success(text string) string {
	return SuccessStyle.Render(text)
}

func Error(text string) string {
	return ErrorStyle.Render(text)
}

func Warning(text string) string {
	return WarningStyle.Render(text)
}

func Info(text string) string {
	return InfoStyle.Render(text)
}

func Bold(text string) string {
	return BoldStyle.Render(text)
}

func Dim(text string) string {
	return DimStyle.Render(text)
}

func Cyan(text string) string {
	return lipgloss.NewStyle().Foreground(CyanColor).Render(text)
}

func Green(text string) string {
	return lipgloss.NewStyle().Foreground(SuccessColor).Render(text)
}

func Red(text string) string {
	return lipgloss.NewStyle().Foreground(ErrorColor).Render(text)
}

func Yellow(text string) string {
	return lipgloss.NewStyle().Foreground(YellowColor).Render(text)
}

// Функции форматирования с форматированием строк
func Successf(format string, a ...interface{}) string {
	return Success(fmt.Sprintf(format, a...))
}

func Errorf(format string, a ...interface{}) string {
	return Error(fmt.Sprintf(format, a...))
}

func Warningf(format string, a ...interface{}) string {
	return Warning(fmt.Sprintf(format, a...))
}

func Infof(format string, a ...interface{}) string {
	return Info(fmt.Sprintf(format, a...))
}

func Boldf(format string, a ...interface{}) string {
	return Bold(fmt.Sprintf(format, a...))
}

func Dimf(format string, a ...interface{}) string {
	return Dim(fmt.Sprintf(format, a...))
}

func Cyanf(format string, a ...interface{}) string {
	return Cyan(fmt.Sprintf(format, a...))
}

// Header выводит заголовок секции
func Header(title string) {
	fmt.Println(TitleStyle.Render(title))
}

// Label выводит метку со значением
func Label(label, value string) {
	labelWidth := 12
	fmt.Printf("  %s  %s\n",
		DimStyle.Width(labelWidth).Render(label+":"),
		value,
	)
}

// Checkmark возвращает зелёную галочку
func Checkmark() string {
	return Success("✓")
}

// Cross возвращает красный крестик
func Cross() string {
	return Error("✗")
}

// Bullet возвращает маркер списка
func Bullet() string {
	return Dim("•")
}

// Ticket возвращает тикет задачи с цветом
func Ticket(ticket string) string {
	return TicketStyle.Render(ticket)
}

// StatusOK возвращает "да" зелёным
func StatusOK() string {
	return Success("да")
}

// StatusNo возвращает "нет" красным
func StatusNo() string {
	return Error("нет")
}

// RoleColor возвращает роль с цветом
func RoleColor(role string) string {
	switch role {
	case "admin":
		return AdminRoleStyle.Render(role)
	case "manager":
		return ManagerRoleStyle.Render(role)
	case "user":
		return UserRoleStyle.Render(role)
	default:
		return role
	}
}

// InProgress возвращает текст зелёным (в работе)
func InProgress(text string) string {
	return InProgressStyle.Render(text)
}

// Paused возвращает текст жёлтым (на паузе)
func Paused(text string) string {
	return PausedStyle.Render(text)
}

// Closed возвращает текст тусклым (закрыто)
func Closed(text string) string {
	return ClosedStyle.Render(text)
}

// Box оборачивает текст в рамку с адаптацией под ширину терминала
func Box(text string) string {
	return BoxStyle().Render(text)
}

// KeyValue выводит пару ключ-значение
func KeyValue(key, value string) {
	fmt.Printf("  %s %s\n",
		DimStyle.Width(15).Render(key+":"),
		value,
	)
}

// Truncate обрезает строку до указанной длины
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return s
	}
	visLen := visibleLength(s)
	if visLen <= maxLen {
		return s
	}
	// Обрезаем с учётом ANSI escape кодов
	result := ""
	count := 0
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			result += string(r)
			continue
		}
		if inEscape {
			result += string(r)
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		if count >= maxLen-1 {
			result += "…"
			break
		}
		result += string(r)
		count++
	}
	return result
}

// visibleLength возвращает видимую длину строки (без ANSI)
func visibleLength(s string) int {
	length := 0
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		length++
	}
	return length
}

// Divider возвращает горизонтальный разделитель
func Divider() string {
	width := TermWidth() - 4
	if width < 20 {
		width = 20
	}
	return Dim(strings.Repeat("─", width))
}
