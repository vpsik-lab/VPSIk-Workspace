package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/cliui"
	"github.com/vpsik/workspace-installer/pkg/config"
	"github.com/vpsik/workspace-installer/pkg/detector"
	"github.com/vpsik/workspace-installer/pkg/plan"
	"github.com/vpsik/workspace-installer/pkg/scanner"
	"github.com/vpsik/workspace-installer/pkg/state"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Scan environment and show installation plan",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		spin := cliui.NewSpinner("Scanning environment...")
		spin.Start()
		scanResult := scanner.Run()
		spin.Stop()

		spin = cliui.NewSpinner("Detecting services...")
		spin.Start()
		enabled := cfg.Services.EnabledList()
		detection := detector.Run(scanResult, enabled)
		svcState := state.Build(detection)
		installPlan := plan.Build(svcState, enabled)
		spin.Stop()

		fmt.Println()
		fmt.Println(cliui.Header(installPlan.Summary()))

		tbl := cliui.NewTable([]cliui.Column{
			{Header: "Action", Width: 10},
			{Header: "Service", Width: 16},
			{Header: "Reason"},
		})
		for _, item := range installPlan.Items {
			action := cliui.Success("skip")
			if item.Action == plan.ActionInstall {
				action = cliui.Warning("install")
			}
			tbl.AddRow(action, item.Service, item.Reason)
		}
		tbl.Print()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
