package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	Success = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.success.Render(fmt.Sprint(a...))
		}
		return color.New(color.FgGreen).Sprint(a...)
	}

	Error = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.error.Render(fmt.Sprint(a...))
		}
		return color.New(color.FgRed).Sprint(a...)
	}

	Warning = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.warning.Render(fmt.Sprint(a...))
		}
		return color.New(color.FgYellow).Sprint(a...)
	}

	Info = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.info.Render(fmt.Sprint(a...))
		}
		return color.New(color.FgCyan).Sprint(a...)
	}

	Bold = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.bold.Render(fmt.Sprint(a...))
		}
		return color.New(color.Bold).Sprint(a...)
	}

	Dim = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.muted.Render(fmt.Sprint(a...))
		}
		return color.New(color.Faint).Sprint(a...)
	}

	Cyan = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.info.Render(fmt.Sprint(a...))
		}
		return color.New(color.FgCyan).Sprint(a...)
	}

	Green  = Success
	Red    = Error
	Yellow = Warning
	Blue   = Info

	Magenta = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		return color.New(color.FgMagenta).Sprint(a...)
	}

	SuccessBold = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.success.Copy().Bold(true).Render(fmt.Sprint(a...))
		}
		return color.New(color.FgGreen, color.Bold).Sprint(a...)
	}

	ErrorBold = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.error.Copy().Bold(true).Render(fmt.Sprint(a...))
		}
		return color.New(color.FgRed, color.Bold).Sprint(a...)
	}

	WarningBold = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.warning.Copy().Bold(true).Render(fmt.Sprint(a...))
		}
		return color.New(color.FgYellow, color.Bold).Sprint(a...)
	}

	InfoBold = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.info.Copy().Bold(true).Render(fmt.Sprint(a...))
		}
		return color.New(color.FgCyan, color.Bold).Sprint(a...)
	}

	CyanBold = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		t := GetTheme()
		if t != nil && t.styles != nil {
			return t.styles.info.Copy().Bold(true).Render(fmt.Sprint(a...))
		}
		return color.New(color.FgCyan, color.Bold).Sprint(a...)
	}

	RedBold    = ErrorBold
	GreenBold  = SuccessBold
	YellowBold = WarningBold

	MagentaBold = func(a ...interface{}) string {
		if !ShouldUseColors() {
			return fmt.Sprint(a...)
		}
		return color.New(color.FgMagenta, color.Bold).Sprint(a...)
	}
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
	if !ShouldUseColors() {
		fmt.Println(title)
		fmt.Println(strings.Repeat("-", len(title)))
		return
	}
	fmt.Println(SectionHeader(title))
}

func Label(label, value string) {
	if !ShouldUseColors() {
		fmt.Printf("  %s: %s\n", label, value)
		return
	}
	fmt.Println(KeyValue(label, value, 16))
}

func Checkmark() string {
	if !ShouldUseColors() {
		return "[OK]"
	}
	return Success("✓")
}

func Cross() string {
	if !ShouldUseColors() {
		return "[ERR]"
	}
	return Error("✗")
}

func Bullet() string {
	if !ShouldUseColors() {
		return "-"
	}
	return Dim("•")
}

func Ticket(ticket string) string {
	if !ShouldUseColors() {
		return ticket
	}
	t := GetTheme()
	if t != nil && t.styles != nil {
		return t.styles.ticket.Render(ticket)
	}
	return CyanBold(ticket)
}

func StatusOK() string {
	if !ShouldUseColors() {
		return "да"
	}
	return Success("да")
}

func StatusNo() string {
	if !ShouldUseColors() {
		return "нет"
	}
	return Error("нет")
}

func RoleColor(role string) string {
	if !ShouldUseColors() {
		return role
	}
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
	if !ShouldUseColors() {
		return text
	}
	return Green(text)
}

func Paused(text string) string {
	if !ShouldUseColors() {
		return text
	}
	return Warning(text)
}

func Closed(text string) string {
	if !ShouldUseColors() {
		return text
	}
	return Dim(text)
}

func TagWithColor(name, hexColor string) string {
	if !ShouldUseColors() || hexColor == "" {
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
