package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"tracker/internal/plugin"
	"tracker/internal/ui"

	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:     "plugin",
	Aliases: []string{"plugins"},
	Short:   "Управление плагинами",
	Long: `Управление плагинами для интеграции с внешними сервисами.

Примеры:
  tracker plugin list                    # Список установленных плагинов
  tracker plugin install jira /path/to/plugin  # Установить плагин
  tracker plugin enable jira             # Включить плагин
  tracker plugin disable jira            # Выключить плагин
  tracker plugin info jira               # Информация о плагине`,
}

var pluginListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls", "l"},
	Short:   "Показать список установленных плагинов",
	RunE: func(cmd *cobra.Command, args []string) error {
		registry, err := plugin.LoadPluginsRegistry()
		if err != nil {
			return fmt.Errorf("ошибка загрузки реестра: %w", err)
		}

		if len(registry.Plugins) == 0 {
			fmt.Println(ui.Info("Плагины не установлены"))
			fmt.Println("Используйте 'tracker plugin install <name> <path>' для установки")
			return nil
		}

		fmt.Println(ui.Bold("Установленные плагины:"))
		fmt.Println()

		for name, config := range registry.Plugins {
			status := ui.Success("✓ включён")
			if !config.Enabled {
				status = ui.Warning("✗ выключен")
			}

			fmt.Printf("  %s %s\n", ui.Bold(name), status)
			fmt.Printf("    Путь: %s\n", ui.Dim(config.Path))
			fmt.Printf("    Версия: %s\n", ui.Dim(config.Version))
			fmt.Println()
		}

		return nil
	},
}

var pluginInstallCmd = &cobra.Command{
	Use:     "install [имя] [путь]",
	Aliases: []string{"add"},
	Short:   "Установить плагин",
	Long: `Установить плагин из локального файла.

Примеры:
  tracker plugin install jira /usr/local/bin/tracker-jira-plugin
  tracker plugin install yandex ./plugins/tracker-yandex`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := args[1]

		absPath, err := filepath.Abs(path)
		if err != nil {
			return fmt.Errorf("ошибка получения абсолютного пути: %w", err)
		}

		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("файл не найден: %s", absPath)
		}

		registry, err := plugin.LoadPluginsRegistry()
		if err != nil {
			return fmt.Errorf("ошибка загрузки реестра: %w", err)
		}

		if _, exists := registry.Plugins[name]; exists {
			return fmt.Errorf("плагин %q уже установлен. Используйте 'tracker plugin uninstall %s' сначала", name, name)
		}

		config := &plugin.PluginConfig{
			Path:    absPath,
			Version: "unknown", // TODO: получать версию из плагина
			Enabled: true,
			Config:  make(map[string]string),
		}

		if err := plugin.SavePluginConfig(name, config); err != nil {
			return fmt.Errorf("ошибка сохранения конфигурации: %w", err)
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Плагин %s установлен", ui.Bold(name)))
		fmt.Printf("  Путь: %s\n", ui.Dim(absPath))
		fmt.Printf("  Статус: %s\n", ui.Success("включён"))

		return nil
	},
}

var pluginUninstallCmd = &cobra.Command{
	Use:     "uninstall [имя]",
	Aliases: []string{"remove", "rm", "del"},
	Short:   "Удалить плагин",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		registry, err := plugin.LoadPluginsRegistry()
		if err != nil {
			return fmt.Errorf("ошибка загрузки реестра: %w", err)
		}

		if _, exists := registry.Plugins[name]; !exists {
			return fmt.Errorf("плагин %q не установлен", name)
		}

		if err := plugin.DeletePluginConfig(name); err != nil {
			return fmt.Errorf("ошибка удаления конфигурации: %w", err)
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Плагин %s удалён", ui.Bold(name)))
		return nil
	},
}

var pluginEnableCmd = &cobra.Command{
	Use:   "enable [имя]",
	Short: "Включить плагин",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		config, err := plugin.GetPluginConfig(name)
		if err != nil {
			return err
		}

		config.Enabled = true
		if err := plugin.SavePluginConfig(name, config); err != nil {
			return fmt.Errorf("ошибка сохранения конфигурации: %w", err)
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Плагин %s включён", ui.Bold(name)))
		return nil
	},
}

var pluginDisableCmd = &cobra.Command{
	Use:   "disable [имя]",
	Short: "Выключить плагин",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		config, err := plugin.GetPluginConfig(name)
		if err != nil {
			return err
		}

		config.Enabled = false
		if err := plugin.SavePluginConfig(name, config); err != nil {
			return fmt.Errorf("ошибка сохранения конфигурации: %w", err)
		}

		fmt.Println(ui.Checkmark(), ui.Successf("Плагин %s выключен", ui.Bold(name)))
		return nil
	},
}

var pluginInfoCmd = &cobra.Command{
	Use:   "info [имя]",
	Short: "Показать информацию о плагине",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		config, err := plugin.GetPluginConfig(name)
		if err != nil {
			return err
		}

		fmt.Println(ui.Bold("Информация о плагине:"))
		fmt.Println()
		fmt.Printf("  Имя: %s\n", ui.Bold(name))
		fmt.Printf("  Путь: %s\n", config.Path)
		fmt.Printf("  Версия: %s\n", config.Version)

		status := ui.Success("включён")
		if !config.Enabled {
			status = ui.Warning("выключен")
		}
		fmt.Printf("  Статус: %s\n", status)

		if len(config.Config) > 0 {
			fmt.Println()
			fmt.Println(ui.Bold("  Конфигурация:"))
			for key, value := range config.Config {
				fmt.Printf("    %s: %s\n", key, value)
			}
		}

		return nil
	},
}

func init() {
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginInstallCmd)
	pluginCmd.AddCommand(pluginUninstallCmd)
	pluginCmd.AddCommand(pluginEnableCmd)
	pluginCmd.AddCommand(pluginDisableCmd)
	pluginCmd.AddCommand(pluginInfoCmd)
}
