package app

import (
	"fmt"
	"sync"

	"tracker/internal/client"
	"tracker/internal/config"
	"tracker/internal/plugin"
)

var (
	provider     client.Client
	pluginMgr    *plugin.PluginManager
	providerOnce sync.Once
	providerMu   sync.RWMutex
)

func InitProvider() error {
	var initErr error
	providerOnce.Do(func() {
		pluginMgr = plugin.NewPluginManager()

		server, err := config.GetCurrentServer()
		if err != nil {
			initErr = err
			return
		}

		if server != nil && server.Plugin != "" {
			p, err := pluginMgr.LoadPlugin(server.Plugin)
			if err != nil {
				fmt.Printf("Предупреждение: не удалось загрузить плагин %q: %v\n", server.Plugin, err)
				fmt.Printf("Используется нативный клиент\n")
				provider = client.NewNativeClient()
				return
			}
			router := plugin.NewRouter(client.NewNativeClient())
			router.AttachPlugin(p)
			provider = router
		} else {
			provider = client.NewNativeClient()
		}
	})
	return initErr
}

func GetClient() client.Client {
	providerMu.RLock()
	defer providerMu.RUnlock()

	if provider == nil {
		return client.NewNativeClient()
	}
	return provider
}

func RefreshProvider() error {
	providerMu.Lock()
	defer providerMu.Unlock()

	if pluginMgr != nil {
		pluginMgr.KillAll()
	}

	providerOnce = sync.Once{}
	provider = nil

	return InitProvider()
}

func Cleanup() {
	providerMu.Lock()
	defer providerMu.Unlock()

	if pluginMgr != nil {
		pluginMgr.KillAll()
	}
}
