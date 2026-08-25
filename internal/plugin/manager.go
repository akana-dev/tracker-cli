package plugin

import (
	"fmt"
	"os/exec"
	"sync"

	"github.com/hashicorp/go-plugin"
)

type PluginManager struct {
	mu            sync.RWMutex
	activePlugins map[string]*plugin.Client
	plugins       map[string]TrackerPlugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		activePlugins: make(map[string]*plugin.Client),
		plugins:       make(map[string]TrackerPlugin),
	}
}

func (m *PluginManager) LoadPlugin(name string) (TrackerPlugin, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p, exists := m.plugins[name]; exists {
		return p, nil
	}

	config, err := GetPluginConfig(name)
	if err != nil {
		return nil, fmt.Errorf("плагин %q не найден: %w", name, err)
	}

	if !config.Enabled {
		return nil, fmt.Errorf("плагин %q отключён", name)
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"tracker": &TrackerPluginPlugin{},
		},
		Cmd: exec.Command(config.Path),
	})

	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("ошибка запуска плагина %q: %w", name, err)
	}

	raw, err := rpcClient.Dispense("tracker")
	if err != nil {
		client.Kill()
		return nil, fmt.Errorf("ошибка получения плагина %q: %w", name, err)
	}

	pluginImpl := raw.(TrackerPlugin)

	if len(config.Config) > 0 {
		if err := pluginImpl.Configure(config.Config); err != nil {
			client.Kill()
			return nil, fmt.Errorf("ошибка конфигурации плагина %q: %w", name, err)
		}
	}

	m.activePlugins[name] = client
	m.plugins[name] = pluginImpl

	return pluginImpl, nil
}

func (m *PluginManager) UnloadPlugin(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.activePlugins[name]
	if !exists {
		return nil
	}

	client.Kill()

	delete(m.activePlugins, name)
	delete(m.plugins, name)

	return nil
}

func (m *PluginManager) GetPlugin(name string) (TrackerPlugin, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.plugins[name]
	return p, exists
}

func (m *PluginManager) ListLoadedPlugins() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.plugins))
	for name := range m.plugins {
		names = append(names, name)
	}
	return names
}

func (m *PluginManager) KillAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, client := range m.activePlugins {
		client.Kill()
		delete(m.activePlugins, name)
		delete(m.plugins, name)
	}
}
