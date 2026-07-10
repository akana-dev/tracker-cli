package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	githubAPIURL         = "https://api.github.com/repos/%s/%s/releases/latest"
	DefaultCheckInterval = 24 * time.Hour
	requestTimeout       = 5 * time.Second
)

type GitHubRelease struct {
	TagName    string        `json:"tag_name"`
	Name       string        `json:"name"`
	Body       string        `json:"body"`
	HtmlURL    string        `json:"html_url"`
	Prerelease bool          `json:"prerelease"`
	Draft      bool          `json:"draft"`
	Assets     []GitHubAsset `json:"assets"`
}

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type CheckResult struct {
	HasUpdate      bool
	LatestVersion  string
	CurrentVersion string
	ReleaseURL     string
	Changelog      string
	DownloadURL    string
}

func CheckForUpdate(ctx context.Context, owner, repo string, includePreRelease bool) (*CheckResult, error) {
	if isDevVersion() {
		return &CheckResult{
			CurrentVersion: "dev",
			HasUpdate:      false,
		}, nil
	}

	currentVersion := getCurrentVersion()

	reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	var release GitHubRelease
	var err error

	if includePreRelease {
		release, err = getLatestReleaseIncludingPreRelease(reqCtx, owner, repo)
	} else {
		release, err = getLatestStableRelease(reqCtx, owner, repo)
	}

	if err != nil {
		return nil, err
	}

	if release.Draft {
		return &CheckResult{
			CurrentVersion: currentVersion,
			HasUpdate:      false,
		}, nil
	}

	latestVersion := release.TagName

	hasUpdate := compareSemverWithSuffix(latestVersion, currentVersion) > 0

	downloadURL, err := findAssetForPlatform(release.Assets)
	if err != nil {
		downloadURL = ""
	}

	return &CheckResult{
		HasUpdate:      hasUpdate,
		LatestVersion:  strings.TrimPrefix(latestVersion, "v"),
		CurrentVersion: currentVersion,
		ReleaseURL:     release.HtmlURL,
		Changelog:      release.Body,
		DownloadURL:    downloadURL,
	}, nil
}

func getLatestStableRelease(ctx context.Context, owner, repo string) (GitHubRelease, error) {
	url := fmt.Sprintf(githubAPIURL, owner, repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return GitHubRelease{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("User-Agent", "tracker-cli/"+getCurrentVersion())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return GitHubRelease{}, fmt.Errorf("ошибка сети: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GitHubRelease{}, fmt.Errorf("GitHub API вернул статус %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return GitHubRelease{}, fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	return release, nil
}

func getLatestReleaseIncludingPreRelease(ctx context.Context, owner, repo string) (GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=10", owner, repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return GitHubRelease{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("User-Agent", "tracker-cli/"+getCurrentVersion())
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return GitHubRelease{}, fmt.Errorf("ошибка сети: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return GitHubRelease{}, fmt.Errorf("GitHub API вернул статус %d", resp.StatusCode)
	}

	var releases []GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return GitHubRelease{}, fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	if len(releases) == 0 {
		return GitHubRelease{}, fmt.Errorf("релизы не найдены")
	}

	latest := releases[0]
	for _, release := range releases[1:] {
		if isMoreRecent(release, latest) {
			latest = release
		}
	}

	return latest, nil
}

func findAssetForPlatform(assets []GitHubAsset) (string, error) {
	targetOS := runtime.GOOS
	targetArch := runtime.GOARCH

	var ext string
	if targetOS == "windows" {
		ext = ".exe"
	}

	targetName := fmt.Sprintf("tracker-%s-%s%s", targetOS, targetArch, ext)

	for _, asset := range assets {
		if asset.Name == targetName {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("бинарник для %s/%s не найден в релизе (ожидался: %s)", targetOS, targetArch, targetName)
}

func isMoreRecent(release, current GitHubRelease) bool {
	v1 := strings.TrimPrefix(release.TagName, "v")
	v2 := strings.TrimPrefix(current.TagName, "v")

	return compareSemverWithSuffix(v1, v2) > 0
}

func compareSemver(v1, v2 string) int {
	return compareSemverWithSuffix(v1, v2)
}

func compareSemverWithSuffix(v1, v2 string) int {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	base1, suffix1 := splitVersionAndSuffix(v1)
	base2, suffix2 := splitVersionAndSuffix(v2)

	parts1 := strings.Split(base1, ".")
	parts2 := strings.Split(base2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var p1, p2 int
		if i < len(parts1) {
			p1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			p2, _ = strconv.Atoi(parts2[i])
		}

		if p1 < p2 {
			return -1
		}
		if p1 > p2 {
			return 1
		}
	}

	if suffix1 == "" && suffix2 == "" {
		return 0
	}
	if suffix1 == "" && suffix2 != "" {
		return 1
	}
	if suffix1 != "" && suffix2 == "" {
		return -1
	}

	if suffix1 < suffix2 {
		return -1
	}
	if suffix1 > suffix2 {
		return 1
	}
	return 0
}

func splitVersionAndSuffix(version string) (string, string) {
	parts := strings.SplitN(version, "-", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
