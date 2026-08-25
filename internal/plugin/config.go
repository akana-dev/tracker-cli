package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type PluginConfig struct {
	Path    string            `json:"path"`
	Version string            `json:"version"`
	Enabled bool              `json:"enabled"`
	Config  map[string]string `json:"config"`
}

type PluginsRegistry struct {
	Plugins map[string]*PluginConfig `json:"plugins"`
}

func GetPluginsConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("не удалось получить домашнюю директорию: %w", err)
	}

	configDir := filepath.Join(homeDir, ".tracker")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("не удалось создать директорию %s: %w", configDir, err)
	}

	return filepath.Join(configDir, "plugins.json"), nil
}

func LoadPluginsRegistry() (*PluginsRegistry, error) {
	configPath, err := GetPluginsConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &PluginsRegistry{Plugins: make(map[string]*PluginConfig)}, nil
		}
		return nil, fmt.Errorf("ошибка чтения plugins.json: %w", err)
	}

	var registry PluginsRegistry
	if err := json.Unmarshal(data, &registry); err != nil {
		return nil, fmt.Errorf("ошибка парсинга plugins.json: %w", err)
	}

	if registry.Plugins == nil {
		registry.Plugins = make(map[string]*PluginConfig)
	}

	return &registry, nil
}

func SavePluginsRegistry(registry *PluginsRegistry) error {
	configPath, err := GetPluginsConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(registry, "", "  ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("ошибка записи plugins.json: %w", err)
	}

	return nil
}

func GetPluginConfig(name string) (*PluginConfig, error) {
	registry, err := LoadPluginsRegistry()
	if err != nil {
		return nil, err
	}

	plugin, exists := registry.Plugins[name]
	if !exists {
		return nil, fmt.Errorf("плагин %q не установлен", name)
	}

	return plugin, nil
}

func SavePluginConfig(name string, config *PluginConfig) error {
	registry, err := LoadPluginsRegistry()
	if err != nil {
		return err
	}

	registry.Plugins[name] = config
	return SavePluginsRegistry(registry)
}

func DeletePluginConfig(name string) error {
	registry, err := LoadPluginsRegistry()
	if err != nil {
		return err
	}

	delete(registry.Plugins, name)
	return SavePluginsRegistry(registry)
}
