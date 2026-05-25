package cmd

import (
	"bufio"
	"fmt"
	"os"
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

var dryRun bool

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

		if !installPlan.HasChanges() {
			fmt.Println("\n✅ All services are already installed. Nothing to do.")
			return nil
		}

		if dryRun {
			fmt.Println("\n🔍 Dry-run mode — no changes applied.")
			return nil
		}

		fmt.Print("\n⚠ Apply this plan? (yes/no): ")
		reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

		if response != "yes" {
			fmt.Println("Installation cancelled.")
			return nil
		}

		var toInstall []string
		for _, item := range installPlan.Items {
			if item.Action == plan.ActionInstall {
				toInstall = append(toInstall, item.Service)
			}
		}

		if err := docker.EnsureNetwork("vpsik"); err != nil {
			return fmt.Errorf("ensure network: %w", err)
		}
		fmt.Println("✓ Network 'vpsik' ready")

		composeDir := filepath.Dir(configPath)
		composePath := filepath.Join(composeDir, "docker-compose.generated.yml")

		envPath := filepath.Join(composeDir, ".env")
		if err := docker.GenerateEnvFile(toInstall, cfg.Workspace.Domain, envPath); err != nil {
			return fmt.Errorf("generate env: %w", err)
		}
		fmt.Printf("✓ Generated .env at %s\n", envPath)

		if err := docker.GenerateComposeFile(toInstall, cfg.Workspace.Domain, composePath); err != nil {
			return fmt.Errorf("generate compose: %w", err)
		}
		fmt.Printf("✓ Generated docker-compose at %s\n", composePath)

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
}

func waitForService(name string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	servicePorts := map[string]int{
		"authentik": 9000,
		"gitea":     3000,
		"coolify":   8000,
		"ollama":    11434,
		"opencode":  30081,
		"openwebui": 3001,
	}
	serviceAPIs := map[string]string{
		"authentik": "http://localhost:9000",
		"gitea":     "http://localhost:3000/api/healthz",
		"coolify":   "http://localhost:8000/api/v1/health",
		"ollama":    "http://localhost:11434/api/tags",
		"openwebui": "http://localhost:3001",
	}

	for time.Now().Before(deadline) {
		if port, ok := servicePorts[name]; ok {
			if detector.CheckPort(port) {
				return true
			}
		}
		if apiURL, ok := serviceAPIs[name]; ok {
			if detector.CheckAPI(apiURL) {
				return true
			}
		}
		time.Sleep(2 * time.Second)
	}
	return false
}
