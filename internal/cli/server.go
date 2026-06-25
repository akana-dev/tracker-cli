package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"tracker/internal/config"
	"tracker/internal/ui"
	"tracker/pkg/table"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Управление серверами",
}

var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "Показать список серверов",
	RunE: func(cmd *cobra.Command, args []string) error {
		servers, err := config.ListServers()
		if err != nil {
			return err
		}

		if len(servers) == 0 {
			fmt.Println(ui.Warning("Серверы не настроены. Добавьте сервер: tracker server add"))
			return nil
		}

		fmt.Println()
		tbl := table.New("Имя", "URL", "Статус", "Роль")
		for _, s := range servers {
			status := ui.Dim("○")
			if s.IsCurrent {
				status = ui.Success("● текущий")
			}
			auth := ui.Error("✗")
			if s.HasToken {
				auth = ui.Success("✓")
			}
			role := s.UserRole
			if role == "" {
				role = "—"
			} else {
				role = ui.RoleColor(role)
			}

			tbl.AddRow(
				ui.Bold(s.Name),
				s.APIURL,
				fmt.Sprintf("%s %s", status, auth),
				role,
			)
		}
		tbl.Render()
		fmt.Println()

		return nil
	},
}

var serverAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Добавить новый сервер",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		url, _ := cmd.Flags().GetString("url")

		if name == "" || url == "" {
			return fmt.Errorf("укажите --name и --url")
		}

		if err := config.AddServer(name, url); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Сервер %s добавлен", ui.Bold(name)))
		fmt.Println(ui.Dimf("URL: %s", url))
		return nil
	},
}

var serverRemoveCmd = &cobra.Command{
	Use:   "remove [имя]",
	Short: "Удалить сервер",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := config.RemoveServer(name); err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Сервер %s удалён", ui.Bold(name)))
		return nil
	},
}

var serverUseCmd = &cobra.Command{
	Use:   "use [имя]",
	Short: "Переключиться на другой сервер",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := config.SetCurrentServer(name); err != nil {
			return err
		}

		server, err := config.GetCurrentServer()
		if err != nil {
			return err
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Переключено на сервер %s", ui.Bold(name)))
		fmt.Println(ui.Dimf("URL: %s", server.APIURL))

		if server.Token == "" {
			fmt.Println(ui.Warning("Требуется авторизация: tracker login"))
		}
		return nil
	},
}

var serverCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Показать текущий сервер",
	RunE: func(cmd *cobra.Command, args []string) error {
		server, err := config.GetCurrentServer()
		if err != nil {
			return err
		}

		fmt.Println()
		ui.Header(fmt.Sprintf("Текущий сервер: %s", ui.Bold(server.Name)))
		ui.Label("URL", server.APIURL)
		ui.Label("Авторизован", func() string {
			if server.Token != "" {
				return ui.StatusOK()
			}
			return ui.StatusNo()
		}())
		ui.Label("Роль", ui.RoleColor(server.UserRole))
		fmt.Println()
		return nil
	},
}

func init() {
	serverAddCmd.Flags().StringP("name", "n", "", "Имя сервера")
	serverAddCmd.Flags().StringP("url", "u", "", "URL API")

	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverAddCmd)
	serverCmd.AddCommand(serverRemoveCmd)
	serverCmd.AddCommand(serverUseCmd)
	serverCmd.AddCommand(serverCurrentCmd)
}
