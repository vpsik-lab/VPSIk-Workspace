package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/config"
	"github.com/vpsik/workspace-installer/pkg/detector"
	"github.com/vpsik/workspace-installer/pkg/docker"
	"github.com/vpsik/workspace-installer/pkg/plan"
	"github.com/vpsik/workspace-installer/pkg/scanner"
	"github.com/vpsik/workspace-installer/pkg/state"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Pull latest images and recreate services",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		fmt.Println("🔍 Scanning environment...")
		scanResult := scanner.Run()

		if !scanResult.DockerAvailable {
			return fmt.Errorf("docker not available")
		}

		enabled := cfg.Services.EnabledList()

		fmt.Println("📋 Detecting installed services...")
		detection := detector.Run(scanResult, enabled)
		svcState := state.Build(detection)
		installPlan := plan.Build(svcState, enabled)

		if installPlan.HasChanges() {
			fmt.Println("\n⚠ Some services are not installed. Upgrade will only update installed ones.")
		}

		fmt.Println("\n📦 Pulling latest images...")
		if err := docker.PullImages(enabled); err != nil {
			return fmt.Errorf("pull images: %w", err)
		}

		fmt.Println("🔄 Recreating services...")
		if err := docker.RecreateServices(configPath, cfg.Workspace.Domain); err != nil {
			return fmt.Errorf("recreate services: %w", err)
		}

		composeDir := cfg.InstallPath()
		if composeDir == "" {
			composeDir = filepath.Dir(configPath)
		}
		svcState.Save(filepath.Join(composeDir, "workspace-state.json"))

		fmt.Println("\n✅ Upgrade complete!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
