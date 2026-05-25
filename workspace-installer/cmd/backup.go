package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/config"
)

var (
	backupAll     bool
	backupDryRun  bool
	backupList    bool
)

var backupCmd = &cobra.Command{
	Use:   "backup [service...]",
	Short: "Create backups of services",
	Long: `Backup service data using restic.

Examples:
  vpsik backup --all              Backup all services
  vpsik backup gitea postgres     Backup specific services
  vpsik backup --list             List available snapshots
  vpsik backup --dry-run          Simulate backup
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		repo := cfg.Backup.Repository
		if repo == "" {
			repo = "/opt/workspace/backups"
		}

		if backupList {
			return listSnapshots(repo)
		}

		if backupDryRun {
			fmt.Println("🔍 Dry-run: would backup the following services:")
		} else {
			fmt.Println("💾 Starting backup...")
		}

		ensureResticRepo(repo, cfg)

		var services []string
		if backupAll {
			services = cfg.Services.EnabledList()
		} else if len(args) > 0 {
			services = args
		} else {
			return fmt.Errorf("specify services, --all, or --list")
		}

		for _, svc := range services {
			if backupDryRun {
				fmt.Printf("  Would backup: %s\n", svc)
				continue
			}
			fmt.Printf("  Backing up %s...", svc)
			if err := backupService(svc, repo); err != nil {
				fmt.Printf(" ❌ %v\n", err)
			} else {
				fmt.Println(" ✅")
			}
		}

		if !backupDryRun {
			fmt.Println("\n✅ Backup complete")
		}
		return nil
	},
}

func backupService(name, repo string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	path, ok := backupPaths[name]
	if !ok {
		return fmt.Errorf("no backup strategy for %s", name)
	}

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/data", repo),
		"-v", fmt.Sprintf("%s:/source:ro", path),
		"restic/restic:latest",
		"-r", "/data",
		"backup", "/source",
		"--tag", name,
		"--hostname", "vpsik",
	)
	cmd.Env = os.Environ()
	pw := os.Getenv("RESTIC_PASSWORD")
	if pw != "" {
		cmd.Env = append(cmd.Env, "RESTIC_PASSWORD="+pw)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("backup failed: %s", string(output))
	}
	return nil
}

func ensureResticRepo(repo string, cfg *config.WorkspaceConfig) {
	if _, err := os.Stat(filepath.Join(repo, "config")); os.IsNotExist(err) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
			"-v", fmt.Sprintf("%s:/data", repo),
			"restic/restic:latest",
			"-r", "/data", "init",
		)
		password := os.Getenv("RESTIC_PASSWORD")
		if password == "" {
			password = "vpsik"
		}
		cmd.Env = append(os.Environ(), "RESTIC_PASSWORD="+password)
		cmd.Run()
	}
}

func listSnapshots(repo string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/data", repo),
		"restic/restic:latest",
		"-r", "/data", "snapshots",
	)
	password := os.Getenv("RESTIC_PASSWORD")
	if password == "" {
		password = "vpsik"
	}
	cmd.Env = append(os.Environ(), "RESTIC_PASSWORD="+password)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

var backupPaths = map[string]string{
	"gitea":      "/var/lib/docker/volumes/gitea-data/_data",
	"postgres":   "/var/lib/docker/volumes/postgres-data/_data",
	"ollama":     "/var/lib/docker/volumes/ollama-data/_data",
	"grafana":    "/var/lib/docker/volumes/grafana-data/_data",
	"authentik":  "/var/lib/docker/volumes/authentik-media/_data",
	"opencode":   "/var/lib/docker/volumes/opencode-data/_data",
	"openwebui":  "/var/lib/docker/volumes/openwebui-data/_data",
	"code-server": "/var/lib/docker/volumes/codeserver-data/_data",
	"outline":    "/var/lib/docker/volumes/outline-data/_data",
	"mattermost": "/var/lib/docker/volumes/mattermost-data/_data",
	"restic":     "/var/lib/docker/volumes/restic-data/_data",
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.Flags().BoolVar(&backupAll, "all", false, "Backup all services")
	backupCmd.Flags().BoolVar(&backupDryRun, "dry-run", false, "Simulate backup")
	backupCmd.Flags().BoolVar(&backupList, "list", false, "List available snapshots")
}
