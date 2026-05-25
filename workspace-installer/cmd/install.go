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

		fmt.Println("🔍 Scanning environment...")
		scanResult := scanner.Run()

		if !scanResult.DockerAvailable {
			fmt.Println("❌ Docker is required but not available")
			return fmt.Errorf("docker not available")
		}

		networkName := cfg.NetworkName()

		if scanResult.CoolifyDetected != nil {
			fmt.Printf("✓ Coolify detected: %s (port %d)\n",
				scanResult.CoolifyDetected.Container,
				scanResult.CoolifyDetected.Port)

			if scanResult.CoolifyDetected.HasProxy {
				fmt.Println("  ⚠ Coolify is using ports 80/443 — will use internal networking")
			}
		} else {
			fmt.Println("ℹ Coolify not detected")
			if !autoYes {
				fmt.Print("  Install Coolify? [y/N]: ")
				reader := bufio.NewReader(os.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))
				if response == "y" || response == "yes" {
					cfg.Services.Coolify = true
				}
			}
		}

		fmt.Println("\n📋 Detecting services...")
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

		if !installPlan.HasChanges() {
			fmt.Println("\n✅ All services are already installed. Nothing to do.")
			return nil
		}

		if dryRun {
			fmt.Println("\n🔍 Dry-run mode — no changes applied.")
			return nil
		}

		if !autoYes {
			fmt.Print("\n⚠ Apply this plan? (yes/no): ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(response)
			if response != "yes" {
				fmt.Println("Installation cancelled.")
				return nil
			}
		}

		var toInstall []string
		for _, item := range installPlan.Items {
			if item.Action == plan.ActionInstall {
				toInstall = append(toInstall, item.Service)
			}
		}

		if err := docker.EnsureNetwork(networkName); err != nil {
			return fmt.Errorf("ensure network: %w", err)
		}
		fmt.Printf("✓ Network '%s' ready\n", networkName)

		composeDir := cfg.InstallPath()
		if composeDir == "" {
			composeDir = filepath.Dir(configPath)
		}
		composePath := filepath.Join(composeDir, "compose", "docker-compose.yml")

		envPath := filepath.Join(filepath.Dir(composePath), ".env")
		if err := docker.GenerateEnvFile(toInstall, cfg.Workspace.Domain, envPath); err != nil {
			return fmt.Errorf("generate env: %w", err)
		}
		fmt.Printf("✓ Generated .env at %s\n", envPath)

		if err := docker.GenerateComposeFile(toInstall, networkName, cfg.Workspace.Domain, composePath); err != nil {
			return fmt.Errorf("generate compose: %w", err)
		}
		fmt.Printf("✓ Generated docker-compose at %s\n", composePath)

		apiConfigPath := filepath.Join(composeDir, "api.yaml")
		adminPassword, err := docker.GenerateAPIConfig(toInstall, cfg.Workspace.Domain, apiConfigPath)
		if err != nil {
			return fmt.Errorf("generate api config: %w", err)
		}
		fmt.Printf("✓ Generated API config at %s\n", apiConfigPath)
		fmt.Printf("  Admin password: %s\n", adminPassword)

		fmt.Println("\n📦 Pulling images...")
		if err := docker.PullImages(toInstall); err != nil {
			return fmt.Errorf("pull images: %w", err)
		}

		fmt.Println("\n🚀 Deploying services...")
		if err := docker.Deploy(composePath); err != nil {
			return fmt.Errorf("deploy: %w", err)
		}

		fmt.Println("\n⏳ Checking service health...")
		for _, name := range toInstall {
			fmt.Printf("  %s...", name)
			healthy := waitForService(name, 60*time.Second)
			if healthy {
				fmt.Println(" ✅")
			} else {
				fmt.Println(" ⚠ not responding yet")
			}
		}

		svcState.Save(filepath.Join(composeDir, "workspace-state.json"))

		fmt.Println("\n✅ Installation complete!")
		fmt.Printf("   Domain: %s\n", cfg.Workspace.Domain)
		fmt.Println("   Services deployed:")
		for _, name := range toInstall {
			fmt.Printf("     - %s\n", name)
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
