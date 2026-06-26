package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"tracker/internal/models"
)

var (
	ConfigDir   string
	ServersFile string
)

var (
	configCache *models.ServersConfig
	configMu    sync.RWMutex
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	ConfigDir = filepath.Join(home, ".tracker")
	ServersFile = filepath.Join(ConfigDir, "servers.json")
}

func invalidateCache() {
	configMu.Lock()
	configCache = nil
	configMu.Unlock()
}

func LoadServersConfig() (*models.ServersConfig, error) {
	configMu.RLock()
	if configCache != nil {
		cached := configCache
		configMu.RUnlock()
		return cached, nil
	}
	configMu.RUnlock()

	data, err := os.ReadFile(ServersFile)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := &models.ServersConfig{Servers: make(map[string]*models.Server)}
			setCache(cfg)
			return cfg, nil
		}
		return nil, err
	}

	var config models.ServersConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if config.Servers == nil {
		config.Servers = make(map[string]*models.Server)
	}

	setCache(&config)
	return &config, nil
}

func setCache(cfg *models.ServersConfig) {
	configMu.Lock()
	configCache = cfg
	configMu.Unlock()
}

func SaveServersConfig(config *models.ServersConfig) error {
	if err := os.MkdirAll(ConfigDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(ServersFile, data, 0600); err != nil {
		return err
	}

	setCache(config)
	return nil
}

func GetCurrentServer() (*models.Server, error) {
	config, err := LoadServersConfig()
	if err != nil {
		return nil, err
	}

	if config.Current == "" {
		return nil, fmt.Errorf("текущий сервер не выбран. Выполните: tracker server add")
	}

	server, exists := config.Servers[config.Current]
	if !exists {
		return nil, fmt.Errorf("сервер '%s' не найден", config.Current)
	}

	return server, nil
}

func SetCurrentServer(name string) error {
	config, err := LoadServersConfig()
	if err != nil {
		return err
	}

	if _, exists := config.Servers[name]; !exists {
		return fmt.Errorf("сервер '%s' не найден", name)
	}

	config.Current = name
	return SaveServersConfig(config)
}

func AddServer(name, apiURL string) error {
	config, err := LoadServersConfig()
	if err != nil {
		return err
	}

	server := &models.Server{
		Name:        name,
		APIURL:      apiURL,
		AuthMethods: []string{"password", "ad"},
	}

	config.Servers[name] = server

	if config.Current == "" {
		config.Current = name
	}

	return SaveServersConfig(config)
}

func RemoveServer(name string) error {
	config, err := LoadServersConfig()
	if err != nil {
		return err
	}

	if _, exists := config.Servers[name]; !exists {
		return fmt.Errorf("сервер '%s' не найден", name)
	}

	delete(config.Servers, name)

	if config.Current == name {
		config.Current = ""
		if len(config.Servers) > 0 {
			for n := range config.Servers {
				config.Current = n
				break
			}
		}
	}

	return SaveServersConfig(config)
}

type ServerInfo struct {
	Name      string
	APIURL    string
	IsCurrent bool
	HasToken  bool
	UserRole  string
}

func ListServers() ([]ServerInfo, error) {
	config, err := LoadServersConfig()
	if err != nil {
		return nil, err
	}

	var result []ServerInfo
	for name, server := range config.Servers {
		result = append(result, ServerInfo{
			Name:      name,
			APIURL:    server.APIURL,
			IsCurrent: name == config.Current,
			HasToken:  server.Token != "",
			UserRole:  server.UserRole,
		})
	}

	return result, nil
}

func GetAPIURL() string {
	server, err := GetCurrentServer()
	if err != nil {
		return "http://localhost:8000"
	}
	return server.APIURL
}

func LoadToken() string {
	server, err := GetCurrentServer()
	if err != nil {
		return ""
	}
	return server.Token
}

func SaveToken(token string) error {
	config, err := LoadServersConfig()
	if err != nil {
		return err
	}

	if config.Current == "" {
		return fmt.Errorf("текущий сервер не выбран")
	}

	server, exists := config.Servers[config.Current]
	if !exists {
		return fmt.Errorf("сервер не найден")
	}

	server.Token = token
	return SaveServersConfig(config)
}

func SaveUserRole(role string) error {
	config, err := LoadServersConfig()
	if err != nil {
		return err
	}

	if config.Current == "" {
		return nil
	}

	server, exists := config.Servers[config.Current]
	if !exists {
		return nil
	}

	server.UserRole = role
	return SaveServersConfig(config)
}

func GetUserRole() string {
	server, err := GetCurrentServer()
	if err != nil {
		return "user"
	}
	if server.UserRole == "" {
		return "user"
	}
	return server.UserRole
}

func GetADDomain() string {
	server, err := GetCurrentServer()
	if err != nil {
		return ""
	}
	return server.ADDomain
}

func SaveConfig(server *models.Server) error {
	config, err := LoadServersConfig()
	if err != nil {
		return err
	}

	if config.Current == "" {
		return fmt.Errorf("текущий сервер не выбран")
	}

	config.Servers[config.Current] = server
	return SaveServersConfig(config)
}
