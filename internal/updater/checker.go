package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	githubAPIURL   = "https://api.github.com/repos/skymoore/vibe-zsh/releases/latest"
	checkInterval  = 7 * 24 * time.Hour
	requestTimeout = 10 * time.Second
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func shouldCheckForUpdates(state *State) bool {
	if !state.CheckEnabled {
		return false
	}

	if os.Getenv("VIBE_AUTO_UPDATE") == "false" {
		return false
	}

	interval := checkInterval
	if envInterval := os.Getenv("VIBE_UPDATE_CHECK_INTERVAL"); envInterval != "" {
		if d, err := time.ParseDuration(envInterval); err == nil {
			interval = d
		}
	}

	return time.Since(state.LastCheck) >= interval
}

func checkLatestVersion() (string, error) {
	client := &http.Client{
		Timeout: requestTimeout,
	}

	req, err := http.NewRequest("GET", githubAPIURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func compareVersions(current, latest string) bool {
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")
	return latest > current && current != "dev"
}

func CheckForUpdatesBackground(currentVersion string) {
	state, err := LoadState()
	if err != nil {
		return
	}

	if !shouldCheckForUpdates(state) {
		return
	}

	latestVersion, err := checkLatestVersion()
	if err != nil {
		return
	}

	state.LastCheck = time.Now()
	state.CurrentVersion = currentVersion

	if compareVersions(currentVersion, latestVersion) {
		if state.AvailableVersion != latestVersion {
			state.AvailableVersion = latestVersion
			state.NotificationShown = false
		}
	} else {
		state.AvailableVersion = ""
		state.NotificationShown = false
	}

	SaveState(state)
}

func ShowUpdateNotification(currentVersion string) {
	state, err := LoadState()
	if err != nil {
		return
	}

	if state.AvailableVersion == "" || state.NotificationShown {
		return
	}

	if !compareVersions(currentVersion, state.AvailableVersion) {
		return
	}

	fmt.Fprintf(os.Stderr, "\n⚠️  vibe %s available (current: %s)\n", state.AvailableVersion, currentVersion)
	fmt.Fprintf(os.Stderr, "   Run: vibe --update\n\n")

	state.NotificationShown = true
	SaveState(state)
}
