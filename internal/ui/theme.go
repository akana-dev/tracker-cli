package ui

import (
	"encoding/json"
	"os"
	"path/filepath"

	"tracker/internal/config"

	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name string `json:"name"`

	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Success   string `json:"success"`
	Warning   string `json:"warning"`
	Error     string `json:"error"`
	Info      string `json:"info"`
	Muted     string `json:"muted"`

	BadgeBg    string `json:"badge_bg"`
	BadgeFg    string `json:"badge_fg"`
	CardBorder string `json:"card_border"`

	styles *themeStyles
}

type themeStyles struct {
	primary   lipgloss.Style
	secondary lipgloss.Style
	success   lipgloss.Style
	warning   lipgloss.Style
	error     lipgloss.Style
	info      lipgloss.Style
	muted     lipgloss.Style
	bold      lipgloss.Style
	ticket    lipgloss.Style
	badge     lipgloss.Style
	card      lipgloss.Style
}

var (
	currentTheme *Theme
	themeFile    = filepath.Join(config.ConfigDir, "theme.json")
)

func InitTheme() {
	if !ShouldUseColors() {
		currentTheme = plainTheme()
		return
	}

	if data, err := os.ReadFile(themeFile); err == nil {
		var t Theme
		if json.Unmarshal(data, &t) == nil && t.Name != "" {
			currentTheme = &t
			currentTheme.compile()
			return
		}
	}

	currentTheme = defaultTheme()
}

func GetTheme() *Theme {
	if currentTheme == nil {
		InitTheme()
	}
	return currentTheme
}

func (t *Theme) compile() {
	if !ShouldUseColors() {
		t.styles = &themeStyles{}
		return
	}

	t.styles = &themeStyles{
		primary:   lipgloss.NewStyle().Foreground(lipgloss.Color(t.Primary)).Bold(true),
		secondary: lipgloss.NewStyle().Foreground(lipgloss.Color(t.Secondary)),
		success:   lipgloss.NewStyle().Foreground(lipgloss.Color(t.Success)),
		warning:   lipgloss.NewStyle().Foreground(lipgloss.Color(t.Warning)),
		error:     lipgloss.NewStyle().Foreground(lipgloss.Color(t.Error)),
		info:      lipgloss.NewStyle().Foreground(lipgloss.Color(t.Info)),
		muted:     lipgloss.NewStyle().Foreground(lipgloss.Color(t.Muted)).Faint(true),
		bold:      lipgloss.NewStyle().Bold(true),
		ticket:    lipgloss.NewStyle().Foreground(lipgloss.Color(t.Primary)).Bold(true),
		badge: lipgloss.NewStyle().
			Padding(0, 1).
			Background(lipgloss.Color(t.BadgeBg)).
			Foreground(lipgloss.Color(t.BadgeFg)),
		card: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(t.CardBorder)).
			Padding(1, 2),
	}
}

func defaultTheme() *Theme {
	t := &Theme{
		Name:       "dark",
		Primary:    "#7AAFFF",
		Secondary:  "#A0A0A0",
		Success:    "#7FFF7F",
		Warning:    "#FFD77F",
		Error:      "#FF7F7F",
		Info:       "#7FDFFF",
		Muted:      "#808080",
		BadgeBg:    "#3A3A5A",
		BadgeFg:    "#FFFFFF",
		CardBorder: "#5A5A8A",
	}
	t.compile()
	return t
}

func lightTheme() *Theme {
	t := &Theme{
		Name:       "light",
		Primary:    "#0066CC",
		Secondary:  "#666666",
		Success:    "#008800",
		Warning:    "#CC8800",
		Error:      "#CC0000",
		Info:       "#0088CC",
		Muted:      "#999999",
		BadgeBg:    "#E0E0F0",
		BadgeFg:    "#000000",
		CardBorder: "#B0B0D0",
	}
	t.compile()
	return t
}

func plainTheme() *Theme {
	t := &Theme{
		Name: "plain",
	}
	t.compile()
	return t
}

func SaveTheme(t *Theme) error {
	if err := os.MkdirAll(config.ConfigDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(themeFile, data, 0600)
}

func SetThemeByName(name string) (*Theme, error) {
	var t *Theme
	switch name {
	case "dark":
		t = defaultTheme()
	case "light":
		t = lightTheme()
	case "plain":
		t = plainTheme()
	default:
		return nil, ErrUnknownTheme
	}

	if err := SaveTheme(t); err != nil {
		return nil, err
	}

	currentTheme = t
	return t, nil
}

func ListThemes() []string {
	return []string{"dark", "light", "plain"}
}

var ErrUnknownTheme = &ThemeError{"неизвестная тема"}

type ThemeError struct {
	msg string
}

func (e *ThemeError) Error() string {
	return e.msg
}
