package cli

import (
	"fmt"

	"tracker/internal/app"
	"tracker/internal/config"
	"tracker/internal/plugin"
	"tracker/internal/ui"
	"tracker/pkg/table"

	"github.com/spf13/cobra"
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
		fmt.Println(ui.SectionHeader("Серверы"))
		fmt.Println()
		fmt.Printf("  %s\n", ui.Dim(fmt.Sprintf("Найдено: %d", len(servers))))
		fmt.Println()
		fmt.Println(ui.Divider(70))
		fmt.Println()

		tbl := table.New("Имя", "URL", "Статус", "Роль", "Плагин")
		tbl.SetColumnWidths(map[int]int{0: 15, 1: 35, 2: 20, 3: 12, 4: 15})

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
				role = ui.Dim("—")
			} else {
				role = ui.RoleColor(role)
			}
			pluginName := ui.Dim("—")
			if s.Plugin != "" {
				pluginName = ui.Cyan(s.Plugin)
			}

			tbl.AddRow(
				ui.Bold(s.Name),
				s.APIURL,
				fmt.Sprintf("%s %s", status, auth),
				role,
				pluginName,
			)
		}

		tbl.Render()

		fmt.Println()
		fmt.Println(ui.Divider(70))
		fmt.Println()
		fmt.Printf("  %s %s\n", ui.Dim("Переключиться:"), ui.Cyan("tracker server use <имя>"))
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
		if err := app.RefreshProvider(); err != nil {
			fmt.Printf("Предупреждение: не удалось обновить клиент: %v\n", err)
		}

		server, err := config.GetCurrentServer()
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(ui.Checkmark(), ui.Successf("Переключено на сервер %s", ui.Bold(name)))
		fmt.Println()
		ui.Label("URL", ui.Cyan(server.APIURL))
		if server.Plugin != "" {
			ui.Label("Плагин", ui.Cyan(server.Plugin))
		}
		if server.Token == "" {
			fmt.Println()
			fmt.Println(ui.Warning("Требуется авторизация: tracker login"))
		}
		fmt.Println()
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
		fmt.Println(ui.SectionHeader(fmt.Sprintf("Текущий сервер: %s", ui.Bold(server.Name))))
		fmt.Println()

		ui.Label("URL", ui.Cyan(server.APIURL))
		ui.Label("Авторизован", func() string {
			if server.Token != "" {
				return ui.StatusOK()
			}
			return ui.StatusNo()
		}())
		ui.Label("Роль", ui.RoleColor(server.UserRole))
		if server.Plugin != "" {
			ui.Label("Плагин", ui.Cyan(server.Plugin))
		}

		fmt.Println()
		fmt.Println(ui.Divider(70))
		fmt.Println()
		return nil
	},
}

var serverSetPluginCmd = &cobra.Command{
	Use:   "set-plugin [имя_сервера] [имя_плагина]",
	Short: "Привязать плагин к серверу",
	Long: `Привязать плагин к серверу. Все запросы к этому серверу будут перенаправлены через плагин.
Примеры:
tracker server set-plugin work jira
tracker server set-plugin personal yandex`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]
		pluginName := args[1]

		cfg, err := config.LoadServersConfig()
		if err != nil {
			return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
		}

		server, exists := cfg.Servers[serverName]
		if !exists {
			return fmt.Errorf("сервер %q не найден", serverName)
		}

		pluginConfig, err := plugin.GetPluginConfig(pluginName)
		if err != nil {
			return fmt.Errorf("плагин %q не установлен: %w", pluginName, err)
		}

		if !pluginConfig.Enabled {
			return fmt.Errorf("плагин %q выключен. Используйте 'tracker plugin enable %s'", pluginName, pluginName)
		}

		server.Plugin = pluginName
		if err := config.SaveServersConfig(cfg); err != nil {
			return fmt.Errorf("ошибка сохранения конфигурации: %w", err)
		}

		if serverName == cfg.Current {
			if err := app.RefreshProvider(); err != nil {
				fmt.Printf("Предупреждение: не удалось загрузить плагин: %v\n", err)
			}
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Плагин %s привязан к серверу %s",
			ui.Bold(pluginName), ui.Bold(serverName)))
		fmt.Println(ui.Dim("Все запросы к этому серверу теперь будут проходить через плагин"))
		return nil
	},
}

var serverUnsetPluginCmd = &cobra.Command{
	Use:   "unset-plugin [имя_сервера]",
	Short: "Убрать плагин с сервера",
	Long: `Убрать привязку плагина к серверу. Запросы будут перенаправлены на нативный API.
Примеры:
tracker server unset-plugin work`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serverName := args[0]

		cfg, err := config.LoadServersConfig()
		if err != nil {
			return fmt.Errorf("ошибка загрузки конфигурации: %w", err)
		}

		server, exists := cfg.Servers[serverName]
		if !exists {
			return fmt.Errorf("сервер %q не найден", serverName)
		}

		if server.Plugin == "" {
			fmt.Println(ui.Warning("Сервер не использует плагин"))
			return nil
		}

		oldPlugin := server.Plugin
		server.Plugin = ""

		if err := config.SaveServersConfig(cfg); err != nil {
			return fmt.Errorf("ошибка сохранения конфигурации: %w", err)
		}

		if serverName == cfg.Current {
			if err := app.RefreshProvider(); err != nil {
				fmt.Printf("Предупреждение: не удалось обновить клиент: %v\n", err)
			}
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Плагин %s отвязан от сервера %s",
			ui.Bold(oldPlugin), ui.Bold(serverName)))
		fmt.Println(ui.Dim("Запросы теперь будут идти через нативный API"))
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
	serverCmd.AddCommand(serverSetPluginCmd)
	serverCmd.AddCommand(serverUnsetPluginCmd)
}
