package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"tracker/internal/config"
	"tracker/internal/ui"
)

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Настроить подключение к API",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetCurrentServer()
		if err != nil {
			if err := config.AddServer("default", "http://localhost:8000"); err != nil {
				return fmt.Errorf("не удалось создать сервер: %w", err)
			}
			server, err = config.GetCurrentServer()
			if err != nil {
				return err
			}
			fmt.Println(ui.Warningf("Создан дефолтный сервер: %s", server.APIURL))
			fmt.Println()
		}

		fmt.Println()
		ui.Header("Текущие настройки:")
		ui.Label("Сервер", ui.Bold(server.Name))
		ui.Label("API URL", ui.Cyan(server.APIURL))
		ui.Label("Методы", ui.Cyan(strings.Join(server.AuthMethods, ", ")))

		adDomain := server.ADDomain
		if adDomain == "" {
			adDomain = "—"
		} else {
			adDomain = ui.Cyan(adDomain)
		}
		ui.Label("AD домен", adDomain)
		ui.Label("Роль", ui.RoleColor(server.UserRole))
		fmt.Println()

		newURL := readLineWithDefault("API URL (Enter — оставить)", server.APIURL)
		server.APIURL = newURL

		newMethods := readLineWithDefault("Методы через запятую (Enter — оставить)", strings.Join(server.AuthMethods, ", "))
		if newMethods != strings.Join(server.AuthMethods, ", ") {
			methods := strings.Split(newMethods, ",")
			for i, m := range methods {
				methods[i] = strings.TrimSpace(m)
			}
			server.AuthMethods = methods
		}

		if contains(server.AuthMethods, "ad") {
			newDomain := readLineWithDefault("AD домен (Enter — оставить)", server.ADDomain)
			server.ADDomain = newDomain
		}

		if err := config.SaveConfig(server); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(ui.Checkmark(), ui.Success("Настройки сохранены"))
		return nil
	},
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
