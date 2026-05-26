package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed VERSION
var versionFile string

var configPath string

var rootCmd = &cobra.Command{
	Use:   "workspace",
	Short: "WorkSpace OS Installer",
	Long: `WorkSpace OS — AI-native engineering workspace bootstrapper.

Detects existing services, plans what needs to be installed,
and reconciles the environment to the desired state.`,
	Version: strings.TrimSpace(versionFile),
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
