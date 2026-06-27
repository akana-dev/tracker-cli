package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
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

	url := fmt.Sprintf(githubAPIURL, owner, repo)

	reqCtx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("User-Agent", "tracker-cli/"+currentVersion)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка сети: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API вернул статус %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("ошибка парсинга ответа: %w", err)
	}

	if release.Draft || (!includePreRelease && release.Prerelease) {
		return &CheckResult{
			CurrentVersion: currentVersion,
			HasUpdate:      false,
		}, nil
	}

	latestVersion := release.TagName
	if len(latestVersion) > 0 && latestVersion[0] == 'v' {
		latestVersion = latestVersion[1:]
	}

	hasUpdate := latestVersion != currentVersion
	downloadURL := findAssetForPlatform(release.Assets)

	return &CheckResult{
		HasUpdate:      hasUpdate,
		LatestVersion:  latestVersion,
		CurrentVersion: currentVersion,
		ReleaseURL:     release.HtmlURL,
		Changelog:      release.Body,
		DownloadURL:    downloadURL,
	}, nil
}

func findAssetForPlatform(assets []GitHubAsset) string {
	targetOS := runtime.GOOS
	targetArch := runtime.GOARCH

	var ext string
	if targetOS == "windows" {
		ext = ".exe"
	}

	targetName := fmt.Sprintf("tracker-%s-%s%s", targetOS, targetArch, ext)

	for _, asset := range assets {
		if asset.Name == targetName {
			return asset.BrowserDownloadURL
		}
	}

	return ""
}
