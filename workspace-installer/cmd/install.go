package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/cliui"
	"github.com/vpsik/workspace-installer/pkg/config"
	"github.com/vpsik/workspace-installer/pkg/detector"
	"github.com/vpsik/workspace-installer/pkg/docker"
	"github.com/vpsik/workspace-installer/pkg/plan"
	"github.com/vpsik/workspace-installer/pkg/scanner"
	"github.com/vpsik/workspace-installer/pkg/state"
)

var (
	dryRun  bool
	autoYes bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install missing services (requires confirmation)",
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
			return fmt.Errorf(cliui.Error("Docker is required but not available"))
		}

		networkName := cfg.NetworkName()

		if scanResult.CoolifyDetected != nil {
			ci := scanResult.CoolifyDetected
			fmt.Println(cliui.Success("  Coolify detected: %s (port %d)", ci.Container, ci.Port))
			if ci.HasProxy {
				fmt.Println(cliui.Warning("  Coolify is using ports 80/443 — will use internal networking"))
			}
		} else {
			fmt.Println(cliui.DimText("  Coolify not detected"))
			if !autoYes {
				fmt.Print(cliui.Highlight("  Install Coolify? [y/N]: "))
				reader := bufio.NewReader(os.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))
				if response == "y" || response == "yes" {
					cfg.Services.Coolify = true
				}
			}
		}

		spin = cliui.NewSpinner("Detecting services...")
		spin.Start()
		enabled := cfg.Services.EnabledList()
		detection := detector.Run(scanResult, enabled)
		svcState := state.Build(detection)
		installPlan := plan.Build(svcState, enabled)
		spin.Stop()

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

		if !installPlan.HasChanges() {
			fmt.Println(cliui.Success("\n  All services already installed. Nothing to do."))
			return nil
		}

		if dryRun {
			fmt.Println(cliui.DimText("\n  Dry-run mode — no changes applied."))
			return nil
		}

		if !autoYes {
			fmt.Print(cliui.Warning("\n  Apply this plan? (yes/no): "))
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(response)
			if response != "yes" {
				fmt.Println(cliui.DimText("  Installation cancelled."))
				return nil
			}
		}

		var toInstall []string
		for _, item := range installPlan.Items {
			if item.Action == plan.ActionInstall {
				toInstall = append(toInstall, item.Service)
			}
		}

		spin = cliui.NewSpinner("Setting up network...")
		spin.Start()
		if err := docker.EnsureNetwork(networkName); err != nil {
			spin.StopError(err)
			return err
		}
		spin.Stop()

		composeDir := cfg.InstallPath()
		if composeDir == "" {
			composeDir = filepath.Dir(configPath)
		}
		composePath := filepath.Join(composeDir, "compose", "docker-compose.yml")
		envPath := filepath.Join(filepath.Dir(composePath), ".env")

		spin = cliui.NewSpinner("Generating configuration files...")
		spin.Start()
		if err := docker.GenerateEnvFile(toInstall, cfg.Workspace.Domain, envPath); err != nil {
			spin.StopError(err)
			return err
		}
		if err := docker.GenerateComposeFile(toInstall, networkName, cfg.Workspace.Domain, composePath); err != nil {
			spin.StopError(err)
			return err
		}
		apiConfigPath := filepath.Join(composeDir, "api.yaml")
		adminPassword, err := docker.GenerateAPIConfig(toInstall, cfg.Workspace.Domain, apiConfigPath)
		if err != nil {
			spin.StopError(err)
			return err
		}
		spin.Stop()
		fmt.Println(cliui.Label("API Config", apiConfigPath))
		fmt.Println(cliui.Label("Admin Password", adminPassword))

		spin = cliui.NewSpinner("Pulling Docker images...")
		spin.Start()
		if err := docker.PullImages(toInstall); err != nil {
			spin.StopError(err)
			return err
		}
		spin.Stop()

		spin = cliui.NewSpinner("Deploying services...")
		spin.Start()
		if err := docker.Deploy(composePath); err != nil {
			spin.StopError(err)
			return err
		}
		spin.Stop()

		fmt.Println(cliui.DimText("\n  Checking service health..."))
		for _, name := range toInstall {
			spin = cliui.NewSpinner(fmt.Sprintf("Waiting for %s...", name))
			spin.Start()
			healthy := waitForService(name, 60*time.Second)
			if healthy {
				spin.Stop()
			} else {
				spin.Stop()
				fmt.Println(cliui.Warning("  %s not responding yet", name))
			}
		}

		svcState.Save(filepath.Join(composeDir, "workspace-state.json"))

		fmt.Println(cliui.Success("\n  Installation complete!"))
		fmt.Println(cliui.Label("Domain", cfg.Workspace.Domain))
		for _, name := range toInstall {
			fmt.Println(cliui.Success("  - %s", name))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show plan without making changes")
	installCmd.Flags().BoolVar(&autoYes, "yes", false, "Skip confirmation prompt")
}

func waitForService(name string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	containerName := name
	if tpl, ok := docker.ServiceTemplates[name]; ok {
		containerName = tpl.Name
	}

	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		cmd := exec.CommandContext(ctx, "docker", "ps",
			"--filter", fmt.Sprintf("name=%s", containerName),
			"--filter", "status=running",
			"--format", "{{.Names}}")
		out, err := cmd.Output()
		cancel()
		if err == nil && strings.TrimSpace(string(out)) != "" {
			return true
		}
		time.Sleep(3 * time.Second)
	}
	return false
}
