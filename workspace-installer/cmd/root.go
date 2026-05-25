package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configPath string

var rootCmd = &cobra.Command{
	Use:   "vpsik",
	Short: "VPSIk Workspace Installer",
	Long: `VPSIk Workspace — AI-native engineering workspace bootstrapper.

Detects existing services, plans what needs to be installed,
and reconciles the environment to the desired state.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "workspace.yaml", "Path to workspace config file")
}
