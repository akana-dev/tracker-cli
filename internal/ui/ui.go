package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	Success = color.New(color.FgGreen).SprintFunc()
	Error   = color.New(color.FgRed).SprintFunc()
	Warning = color.New(color.FgYellow).SprintFunc()
	Info    = color.New(color.FgCyan).SprintFunc()
	Bold    = color.New(color.Bold).SprintFunc()
	Dim     = color.New(color.Faint).SprintFunc()
	Cyan    = color.New(color.FgCyan).SprintFunc()
	Green   = color.New(color.FgGreen).SprintFunc()
	Red     = color.New(color.FgRed).SprintFunc()
	Yellow  = color.New(color.FgYellow).SprintFunc()
	Blue    = color.New(color.FgBlue).SprintFunc()
	Magenta = color.New(color.FgMagenta).SprintFunc()

	SuccessBold = color.New(color.FgGreen, color.Bold).SprintFunc()
	ErrorBold   = color.New(color.FgRed, color.Bold).SprintFunc()
	WarningBold = color.New(color.FgYellow, color.Bold).SprintFunc()
	InfoBold    = color.New(color.FgCyan, color.Bold).SprintFunc()
	CyanBold    = color.New(color.FgCyan, color.Bold).SprintFunc()
	RedBold     = color.New(color.FgRed, color.Bold).SprintFunc()
	GreenBold   = color.New(color.FgGreen, color.Bold).SprintFunc()
	YellowBold  = color.New(color.FgYellow, color.Bold).SprintFunc()
	MagentaBold = color.New(color.FgMagenta, color.Bold).SprintFunc()
)

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

func CyanBoldf(format string, a ...interface{}) string {
	return CyanBold(fmt.Sprintf(format, a...))
}

func Header(title string) {
	fmt.Println(InfoBold(title))
}

func Label(label, value string) {
	fmt.Printf("  %s  %s\n", Dim(label+":"), value)
}

func Checkmark() string {
	return Success("✓")
}

func Cross() string {
	return Error("✗")
}

func Bullet() string {
	return Dim("•")
}

func Ticket(ticket string) string {
	return CyanBold(ticket)
}

func StatusOK() string {
	return Success("да")
}

func StatusNo() string {
	return Error("нет")
}

func RoleColor(role string) string {
	switch role {
	case "admin":
		return RedBold(role)
	case "manager":
		return YellowBold(role)
	case "user":
		return Green(role)
	default:
		return role
	}
}

func InProgress(text string) string {
	return Green(text)
}

func Paused(text string) string {
	return Warning(text)
}

func Closed(text string) string {
	return Dim(text)
}

func TagWithColor(name, hexColor string) string {
	if hexColor == "" {
		return name
	}

	r, g, b, ok := parseHexColor(hexColor)
	if !ok {
		return name
	}

	return fmt.Sprintf("\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, name)
}

func parseHexColor(hex string) (r, g, b uint8, ok bool) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0, 0, 0, false
	}

	rVal, err := strconv.ParseUint(hex[0:2], 16, 8)
	if err != nil {
		return 0, 0, 0, false
	}
	gVal, err := strconv.ParseUint(hex[2:4], 16, 8)
	if err != nil {
		return 0, 0, 0, false
	}
	bVal, err := strconv.ParseUint(hex[4:6], 16, 8)
	if err != nil {
		return 0, 0, 0, false
	}

	return uint8(rVal), uint8(gVal), uint8(bVal), true
}

func TagsDisplay(tags []TagInfo) string {
	if len(tags) == 0 {
		return ""
	}

	parts := make([]string, 0, len(tags))
	for _, t := range tags {
		parts = append(parts, TagWithColor(t.Name, t.Color))
	}

	return strings.Join(parts, ", ")
}

type TagInfo struct {
	Name  string
	Color string
}
