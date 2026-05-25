package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/config"
	"github.com/vpsik/workspace-installer/pkg/detector"
	"github.com/vpsik/workspace-installer/pkg/scanner"
	"github.com/vpsik/workspace-installer/pkg/state"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Scan environment and show current service status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		fmt.Println("🔍 Scanning environment...")
		scanResult := scanner.Run()

		if !scanResult.DockerAvailable {
			fmt.Println("⚠ Docker is not available")
		} else {
			fmt.Printf("✓ Docker available\n")
			fmt.Printf("  Containers: %d\n", len(scanResult.Containers))
			fmt.Printf("  Networks: %d\n", len(scanResult.Networks))
		}

		fmt.Println("\n📋 Detecting services...")
		enabled := cfg.Services.EnabledList()
		detection := detector.Run(scanResult, enabled)
		svcState := state.Build(detection)

		for _, s := range svcState.Services {
			icon := "❌"
			if s.Status.String() == "installed" {
				icon = "✅"
			}
			fmt.Printf("  %s %s: %s\n", icon, s.Name, s.Details)
		}

		if len(scanResult.Errors) > 0 {
			fmt.Println("\n⚠ Errors:")
			for _, err := range scanResult.Errors {
				fmt.Fprintf(os.Stderr, "  %s\n", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
