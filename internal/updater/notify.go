package updater

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"tracker/internal/ui"
)

const (
	DefaultNotifyInterval = 24 * time.Hour
)

var (
	done     = make(chan struct{})
	doneOnce sync.Once

	notified = false
	notifyMu sync.Mutex
)

func NotifyAboutUpdate(result *CheckResult) {
	notifyMu.Lock()
	if notified {
		notifyMu.Unlock()
		return
	}
	notified = true
	notifyMu.Unlock()

	if result == nil || !result.HasUpdate {
		return
	}

	if isDevVersion() {
		return
	}

	cache := LoadCache()
	if !cache.ShouldNotify(DefaultNotifyInterval) {
		return
	}

	fmt.Printf("\n%s → %s (%s)\n",
		ui.Dim(fmt.Sprintf("Доступна новая версия: %s", formatVersion(result.CurrentVersion))),
		ui.Success(formatVersion(result.LatestVersion)),
		ui.Dim("выполните: tracker update"),
	)

	cache.MarkNotified()
	SaveCache(cache)
}

func ShowFullUpdateInfo(result *CheckResult) {
	if result == nil || !result.HasUpdate {
		return
	}

	fmt.Println()
	fmt.Printf("Доступна новая версия: %s → %s\n",
		ui.Dim(formatVersion(result.CurrentVersion)),
		ui.Success(formatVersion(result.LatestVersion)))
	fmt.Println()

	if result.Changelog != "" {
		fmt.Println(ui.Bold("Что нового:"))
		lines := strings.Split(result.Changelog, "\n")
		shown := 0
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			fmt.Printf("  %s %s\n", ui.Bullet(), line)
			shown++
			if shown >= 10 {
				fmt.Println(ui.Dim("  ... и другие изменения"))
				break
			}
		}
		fmt.Println()
	}

	fmt.Printf("  Для обновления выполните: %s\n", ui.Cyan("tracker update"))
	fmt.Printf("  Или скачайте вручную: %s\n", ui.Dim(result.ReleaseURL))
	fmt.Println()
}

func formatVersion(v string) string {
	if v == "dev" || v == "" {
		return "dev"
	}
	if strings.HasPrefix(v, "v") {
		return v
	}
	return "v" + v
}

func CheckAndNotify(owner, repo string, checkInterval time.Duration) {
	defer doneOnce.Do(func() { close(done) })

	if os.Getenv("TRACKER_NO_UPDATE_CHECK") != "" {
		return
	}

	cache := LoadCache()

	if !cache.IsExpired(checkInterval) && cache.HasUpdate {
		NotifyAboutUpdate(&CheckResult{
			HasUpdate:      cache.HasUpdate,
			LatestVersion:  cache.LatestVersion,
			CurrentVersion: getCurrentVersion(),
			ReleaseURL:     cache.ReleaseURL,
			Changelog:      cache.Changelog,
		})
		return
	}

	ctx := context.Background()
	result, err := CheckForUpdate(ctx, owner, repo, false)
	if err != nil {
		return
	}

	cache = &CheckCache{
		LastCheck:     time.Now(),
		LastNotified:  cache.LastNotified,
		LatestVersion: result.LatestVersion,
		ReleaseURL:    result.ReleaseURL,
		Changelog:     result.Changelog,
		HasUpdate:     result.HasUpdate,
	}
	SaveCache(cache)

	if result.HasUpdate {
		NotifyAboutUpdate(result)
	}
}

func WaitForCheck(timeout int) {
	select {
	case <-done:
	case <-time.After(time.Duration(timeout) * time.Second):
	}
}

func ResetNotification() {
	notifyMu.Lock()
	notified = false
	notifyMu.Unlock()
}
