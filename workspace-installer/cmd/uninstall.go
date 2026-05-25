package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/config"
	"github.com/vpsik/workspace-installer/pkg/docker"
)

var (
	removeVolumes bool
	removeNetwork bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove all deployed services",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		composeDir := filepath.Dir(configPath)
		composePath := filepath.Join(composeDir, "docker-compose.generated.yml")

		if _, err := os.Stat(composePath); os.IsNotExist(err) {
			fmt.Println("No generated compose file found. Nothing to uninstall.")
			return nil
		}

		fmt.Print("⚠ This will stop and remove all services. Continue? (yes/no): ")
		var response string
		fmt.Scanln(&response)
		if response != "yes" {
			fmt.Println("Uninstall cancelled.")
			return nil
		}

		fmt.Println("\n🛑 Stopping and removing services...")
		if err := docker.Down(composePath); err != nil {
			return fmt.Errorf("docker compose down: %w", err)
		}

		if removeVolumes {
			fmt.Println("🗑 Removing volumes...")
			if err := docker.RemoveVolumes(composePath); err != nil {
				return fmt.Errorf("remove volumes: %w", err)
			}
		}

		if removeNetwork {
			fmt.Println("🌐 Removing network 'workspace_net'...")
			if err := docker.RemoveNetwork("workspace_net"); err != nil {
				return fmt.Errorf("remove network: %w", err)
			}
		}

		if err := os.Remove(composePath); err == nil {
			fmt.Println("🧹 Removed generated compose file.")
		}

		fmt.Println("\n✅ Uninstall complete.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolVar(&removeVolumes, "volumes", false, "Remove persistent volumes")
	uninstallCmd.Flags().BoolVar(&removeNetwork, "network", false, "Remove the workspace network")
}
