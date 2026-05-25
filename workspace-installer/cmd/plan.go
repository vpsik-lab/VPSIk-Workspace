package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
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

		fmt.Println("🔍 Scanning environment...")
		scanResult := scanner.Run()

		fmt.Println("📋 Detecting services...")
		enabled := cfg.Services.EnabledList()
		detection := detector.Run(scanResult, enabled)
		svcState := state.Build(detection)

		installPlan := plan.Build(svcState, enabled)

		fmt.Printf("\n📊 Plan: %s\n\n", installPlan.Summary())

		for _, item := range installPlan.Items {
			switch item.Action {
			case plan.ActionInstall:
				fmt.Printf("  🔧 %s — %s\n", item.Service, item.Reason)
			case plan.ActionSkip:
				fmt.Printf("  ✅ %s — %s\n", item.Service, item.Reason)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(planCmd)
}
