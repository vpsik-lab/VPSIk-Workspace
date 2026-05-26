package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/cliui"
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

		spin := cliui.NewSpinner("Scanning environment...")
		spin.Start()
		scanResult := scanner.Run()
		spin.Stop()

		if !scanResult.DockerAvailable {
			fmt.Println(cliui.Warning("  Docker is not available"))
		} else {
			fmt.Println(cliui.Success("  Docker available"))
			fmt.Printf("  Containers: %d | Networks: %d\n", len(scanResult.Containers), len(scanResult.Networks))
		}

		spin = cliui.NewSpinner("Detecting services...")
		spin.Start()
		enabled := cfg.Services.EnabledList()
		detection := detector.Run(scanResult, enabled)
		svcState := state.Build(detection)
		spin.Stop()

		tbl := cliui.NewTable([]cliui.Column{
			{Header: "Service", Width: 16},
			{Header: "Status", Width: 12},
			{Header: "Details"},
		})
		for _, s := range svcState.Services {
			status := cliui.Error("✗ missing")
			if s.Status == detector.StatusInstalled {
				status = cliui.Success("✓ installed")
			}
			tbl.AddRow(s.Name, status, s.Details)
		}
		tbl.Print()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
