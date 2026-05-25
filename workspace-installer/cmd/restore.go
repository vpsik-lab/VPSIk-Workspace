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
	restoreLatest  bool
	restoreList    bool
	restoreSnapshot string
)

var restoreCmd = &cobra.Command{
	Use:   "restore <service>",
	Short: "Restore a service from backup",
	Long: `Restore service data from restic snapshots.

Examples:
  vpsik restore gitea --latest          Restore latest gitea backup
  vpsik restore postgres --snapshot ID  Restore specific snapshot
  vpsik restore --list                  List available snapshots
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

		if restoreList {
			return listSnapshots(repo)
		}

		if len(args) < 1 {
			return fmt.Errorf("service name required")
		}

		service := args[0]

		path, ok := backupPaths[service]
		if !ok {
			return fmt.Errorf("no restore strategy for %s", service)
		}

		if restoreLatest {
			fmt.Printf("Restoring latest snapshot for %s...\n", service)
			return restoreLatestSnapshot(service, path, repo)
		}

		if restoreSnapshot != "" {
			fmt.Printf("Restoring snapshot %s for %s...\n", restoreSnapshot, service)
			return restoreSpecificSnapshot(restoreSnapshot, path, repo)
		}

		return fmt.Errorf("specify --latest or --snapshot")
	},
}

func restoreLatestSnapshot(service, target, repo string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	parent := filepath.Dir(target)

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/data", repo),
		"-v", fmt.Sprintf("%s:/restore", parent),
		"restic/restic:latest",
		"-r", "/data",
		"restore", "latest",
		"--tag", service,
		"--target", "/restore",
		"--host", "vpsik",
	)
	cmd.Env = os.Environ()
	if pw := os.Getenv("RESTIC_PASSWORD"); pw != "" {
		cmd.Env = append(cmd.Env, "RESTIC_PASSWORD="+pw)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}

	return nil
}

func restoreSpecificSnapshot(snapshotID, target, repo string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/data", repo),
		"-v", fmt.Sprintf("%s:/restore", target),
		"restic/restic:latest",
		"-r", "/data",
		"restore", snapshotID,
		"--target", "/restore",
	)
	cmd.Env = os.Environ()
	if pw := os.Getenv("RESTIC_PASSWORD"); pw != "" {
		cmd.Env = append(cmd.Env, "RESTIC_PASSWORD="+pw)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("restore failed: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.Flags().BoolVar(&restoreLatest, "latest", false, "Restore latest snapshot")
	restoreCmd.Flags().BoolVar(&restoreList, "list", false, "List available snapshots")
	restoreCmd.Flags().StringVar(&restoreSnapshot, "snapshot", "", "Snapshot ID to restore")
}
