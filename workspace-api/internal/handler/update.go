package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

var currentVersion = "v0.1.0"

type UpdateCheckResponse struct {
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version"`
	UpdateAvailable bool   `json:"update_available"`
	ReleaseURL      string `json:"release_url"`
}

func CheckUpdate(w http.ResponseWriter, r *http.Request) {
	latest, err := fetchLatestVersion()
	latestVersion := latest
	if err != nil {
		latestVersion = currentVersion
	}

	updateAvailable := false
	if latestVersion != currentVersion && err == nil {
		updateAvailable = true
	}

	resp := UpdateCheckResponse{
		CurrentVersion:  currentVersion,
		LatestVersion:   latestVersion,
		UpdateAvailable: updateAvailable,
		ReleaseURL:      fmt.Sprintf("https://github.com/vpsik-lab/VPSIk-Workspace/releases/tag/%s", latestVersion),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func fetchLatestVersion() (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/vpsik-lab/VPSIk-Workspace/releases/latest")
	if err != nil {
		return "", fmt.Errorf("github api: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if release.TagName == "" {
		return "", fmt.Errorf("no releases found")
	}

	return release.TagName, nil
}

func init() {
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range buildInfo.Settings {
			if setting.Key == "vcs.tag" && setting.Value != "" {
				currentVersion = strings.TrimPrefix(setting.Value, "v")
				if !strings.HasPrefix(currentVersion, "v") {
					currentVersion = "v" + currentVersion
				}
			}
		}
	}
}
