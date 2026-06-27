package updater

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"tracker/internal/config"
)

type CheckCache struct {
	LastCheck     time.Time `json:"last_check"`
	LastNotified  time.Time `json:"last_notified,omitempty"`
	LatestVersion string    `json:"latest_version"`
	ReleaseURL    string    `json:"release_url"`
	Changelog     string    `json:"changelog"`
	HasUpdate     bool      `json:"has_update"`
	Error         string    `json:"error,omitempty"`
}

var cacheFile = filepath.Join(config.ConfigDir, "update-check.json")

func LoadCache() *CheckCache {
	data, err := os.ReadFile(cacheFile)
	if err != nil {
		return &CheckCache{}
	}

	var cache CheckCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return &CheckCache{}
	}

	return &cache
}

func SaveCache(cache *CheckCache) error {
	if err := os.MkdirAll(config.ConfigDir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, data, 0600)
}

func (c *CheckCache) IsExpired(interval time.Duration) bool {
	return time.Since(c.LastCheck) > interval
}

func (c *CheckCache) ShouldNotify(notifyInterval time.Duration) bool {
	return time.Since(c.LastNotified) > notifyInterval
}

func (c *CheckCache) MarkNotified() {
	c.LastNotified = time.Now()
}
