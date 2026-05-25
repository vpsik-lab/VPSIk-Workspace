package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ResticClient struct {
	binaryPath string
	repoURL    string
	password   string
}

type ResticSnapshot struct {
	ID        string            `json:"id"`
	Time      string            `json:"time"`
	Hostname  string            `json:"hostname"`
	Tags      []string          `json:"tags"`
	Paths     []string          `json:"paths"`
	ShortID   string            `json:"short_id"`
}

type ResticBackupStats struct {
	FilesNew        int    `json:"files_new"`
	FilesChanged    int    `json:"files_changed"`
	FilesUnmodified int    `json:"files_unmodified"`
	DirNew          int    `json:"dir_new"`
	DirChanged      int    `json:"dir_changed"`
	TotalBytes      int64  `json:"total_bytes"`
}

func NewResticClient(binaryPath, repoURL, password string) *ResticClient {
	if binaryPath == "" {
		binaryPath = "restic"
	}
	return &ResticClient{
		binaryPath: binaryPath,
		repoURL:    repoURL,
		password:   password,
	}
}

func (c *ResticClient) env() []string {
	return append(os.Environ(),
		fmt.Sprintf("RESTIC_REPOSITORY=%s", c.repoURL),
		fmt.Sprintf("RESTIC_PASSWORD=%s", c.password),
	)
}

func (c *ResticClient) run(ctx context.Context, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, c.binaryPath, args...)
	cmd.Env = c.env()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("restic %s: %s", strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func (c *ResticClient) CheckHealth(ctx context.Context) error {
	_, err := c.run(ctx, "version")
	return err
}

func (c *ResticClient) Backup(ctx context.Context, paths []string, tags []string) (*ResticBackupStats, error) {
	args := []string{"backup", "--json"}
	for _, tag := range tags {
		args = append(args, "--tag", tag)
	}
	args = append(args, "--")
	for _, path := range paths {
		args = append(args, path)
	}

	out, err := c.run(ctx, args...)
	if err != nil {
		return nil, err
	}

	var stats ResticBackupStats
	if err := json.Unmarshal([]byte(out), &stats); err != nil {
		return nil, fmt.Errorf("parse backup output: %w", err)
	}
	return &stats, nil
}

func (c *ResticClient) ListSnapshots(ctx context.Context) ([]ResticSnapshot, error) {
	out, err := c.run(ctx, "snapshots", "--json")
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(out) == "" || strings.TrimSpace(out) == "null" {
		return []ResticSnapshot{}, nil
	}

	var snapshots []ResticSnapshot
	if err := json.Unmarshal([]byte(out), &snapshots); err != nil {
		return nil, fmt.Errorf("parse snapshots: %w", err)
	}
	return snapshots, nil
}

func (c *ResticClient) Restore(ctx context.Context, snapshotID, target string) error {
	_, err := c.run(ctx, "restore", snapshotID, "--target", target)
	return err
}

func (c *ResticClient) Forget(ctx context.Context, keepLast int, tags []string) error {
	args := []string{"forget", "--keep-last", fmt.Sprintf("%d", keepLast), "--prune"}
	for _, tag := range tags {
		args = append(args, "--tag", tag)
	}

	_, err := c.run(ctx, args...)
	return err
}

func (c *ResticClient) Check(ctx context.Context) error {
	_, err := c.run(ctx, "check")
	return err
}

func (c *ResticClient) Stats(ctx context.Context) (map[string]interface{}, error) {
	out, err := c.run(ctx, "stats", "--json")
	if err != nil {
		return nil, err
	}

	var stats map[string]interface{}
	if err := json.Unmarshal([]byte(out), &stats); err != nil {
		return nil, fmt.Errorf("parse stats: %w", err)
	}
	return stats, nil
}
