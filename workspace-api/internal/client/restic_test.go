package client

import (
	"encoding/json"
	"testing"
)

func TestNewResticClient_Defaults(t *testing.T) {
	c := NewResticClient("", "/repo", "pass")
	if c.binaryPath != "restic" {
		t.Errorf("expected default restic, got %s", c.binaryPath)
	}
	if c.repoURL != "/repo" {
		t.Errorf("expected /repo, got %s", c.repoURL)
	}
	if c.password != "pass" {
		t.Errorf("expected pass, got %s", c.password)
	}
}

func TestNewResticClient_CustomBinary(t *testing.T) {
	c := NewResticClient("/usr/local/bin/restic", "/repo", "pass")
	if c.binaryPath != "/usr/local/bin/restic" {
		t.Errorf("expected custom path, got %s", c.binaryPath)
	}
}

func TestNewResticClient_EmptyBinary(t *testing.T) {
	c := NewResticClient("", "/repo", "pass")
	if c.binaryPath != "restic" {
		t.Errorf("expected restic, got %s", c.binaryPath)
	}
}

func TestResticSnapshot_JSON(t *testing.T) {
	jsonData := `[
		{"id": "abc123", "short_id": "abc", "time": "2024-01-15T10:00:00Z", "hostname": "server1", "tags": ["daily", "prod"], "paths": ["/data", "/etc"]},
		{"id": "def456", "short_id": "def", "time": "2024-01-16T10:00:00Z", "hostname": "server1", "tags": null, "paths": ["/data"]}
	]`

	expectedShortIDs := []string{"abc", "def"}
	expectedHostnames := []string{"server1", "server1"}
	expectedTagCounts := []int{2, 0}
	expectedPathCounts := []int{2, 1}

	snapshots, err := parseSnapshots(jsonData)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if len(snapshots) != 2 {
		t.Fatalf("expected 2 snapshots, got %d", len(snapshots))
	}

	for i, snap := range snapshots {
		if snap.ShortID != expectedShortIDs[i] {
			t.Errorf("snapshot[%d]: expected short_id %s, got %s", i, expectedShortIDs[i], snap.ShortID)
		}
		if snap.Hostname != expectedHostnames[i] {
			t.Errorf("snapshot[%d]: expected hostname %s, got %s", i, expectedHostnames[i], snap.Hostname)
		}
		if len(snap.Tags) != expectedTagCounts[i] {
			t.Errorf("snapshot[%d]: expected %d tags, got %d", i, expectedTagCounts[i], len(snap.Tags))
		}
		if len(snap.Paths) != expectedPathCounts[i] {
			t.Errorf("snapshot[%d]: expected %d paths, got %d", i, expectedPathCounts[i], len(snap.Paths))
		}
	}
}

func TestResticSnapshot_EmptyJSON(t *testing.T) {
	snapshots, err := parseSnapshots("[]")
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if len(snapshots) != 0 {
		t.Errorf("expected 0 snapshots, got %d", len(snapshots))
	}
}

func TestResticSnapshot_InvalidJSON(t *testing.T) {
	_, err := parseSnapshots("{invalid}")
	if err == nil {
		t.Fatal("expected parse error for invalid JSON")
	}
}

func TestResticBackupStats_JSON(t *testing.T) {
	jsonData := `{"files_new": 10, "files_changed": 3, "files_unmodified": 100, "dir_new": 2, "dir_changed": 1, "total_bytes": 1048576}`

	stats, err := parseBackupStats(jsonData)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	if stats.FilesNew != 10 {
		t.Errorf("expected 10 new files, got %d", stats.FilesNew)
	}
	if stats.FilesChanged != 3 {
		t.Errorf("expected 3 changed files, got %d", stats.FilesChanged)
	}
	if stats.TotalBytes != 1048576 {
		t.Errorf("expected 1048576 bytes, got %d", stats.TotalBytes)
	}
}

func parseSnapshots(data string) ([]ResticSnapshot, error) {
	var snapshots []ResticSnapshot
	if err := json.Unmarshal([]byte(data), &snapshots); err != nil {
		return nil, err
	}
	return snapshots, nil
}

func parseBackupStats(data string) (*ResticBackupStats, error) {
	var stats ResticBackupStats
	if err := json.Unmarshal([]byte(data), &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}
