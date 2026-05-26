package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/config"
	"gopkg.in/yaml.v3"
)

var (
	autoMode     bool
	domainFlag   string
	servicesFlag string
	outputFlag   string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a workspace configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		if autoMode {
			return runAutoInit()
		}
		return runInteractiveInit()
	},
}

func runAutoInit() error {
	domain := domainFlag
	if domain == "" {
		domain = "workspace.vpsik.com"
	}

	enabledServices := config.DefaultConfig().Services.EnabledList()
	if servicesFlag != "" {
		enabledServices = parseServices(servicesFlag)
	}

	svcMap := make(map[string]bool)
	for _, s := range enabledServices {
		svcMap[s] = true
	}

	cfg := config.DefaultConfig()
	cfg.Workspace.Domain = domain
	cfg.Services.Traefik = svcMap["traefik"]
	cfg.Services.Postgres = svcMap["postgres"]
	cfg.Services.Redis = contains(enabledServices, "redis")
	cfg.Services.Authentik = contains(enabledServices, "authentik")
	cfg.Services.Gitea = contains(enabledServices, "gitea")
	cfg.Services.Ollama = contains(enabledServices, "ollama")
	cfg.Services.OpenWebUI = contains(enabledServices, "openwebui")
	cfg.Services.CodeServer = contains(enabledServices, "code-server")
	cfg.Services.Restic = contains(enabledServices, "restic")
	cfg.Services.Grafana = contains(enabledServices, "grafana")
	cfg.Services.Prometheus = contains(enabledServices, "prometheus")
	cfg.Services.Coolify = contains(enabledServices, "coolify")
	cfg.Services.Mattermost = contains(enabledServices, "mattermost")
	cfg.Services.Outline = contains(enabledServices, "outline")
	cfg.Services.Plane = contains(enabledServices, "plane")
	cfg.Services.OpenCode = contains(enabledServices, "opencode")
	cfg.Services.Cloudflare = contains(enabledServices, "cloudflare")
	cfg.Services.Api = contains(enabledServices, "api")
	cfg.Services.Dashboard = contains(enabledServices, "dashboard")

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	outputPath := outputFlag
	if outputPath == "" {
		outputPath = configPath
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("create output dir: %w", err)
		}
	}

	if info, err := os.Stat(outputPath); err == nil && info.IsDir() {
		if err := os.Remove(outputPath); err != nil {
			return fmt.Errorf("remove existing directory at output path: %w", err)
		}
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Printf("✅ Configuration written to %s\n", outputPath)
	fmt.Printf("   Domain: %s\n", domain)
	fmt.Printf("   Services: %s\n", strings.Join(enabledServices, ", "))
	fmt.Println("\nRun 'workspace doctor' to check system health.")
	fmt.Println("Run 'workspace plan' to see what needs to be installed.")
	fmt.Println("Run 'workspace install --yes' to deploy.")

	return nil
}

func runInteractiveInit() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🚀 WorkSpace OS — Interactive Setup")
	fmt.Println(strings.Repeat("─", 40))

	fmt.Print("Domain (e.g., workspace.vpsik.com) [workspace.vpsik.com]: ")
	domain, _ := reader.ReadString('\n')
	domain = strings.TrimSpace(domain)
	if domain == "" {
		domain = "workspace.vpsik.com"
	}

	allServices := []struct {
		name  string
		label string
	}{
		{"authentik", "Authentik (SSO / Identity Provider)"},
		{"gitea", "Gitea (Git Service)"},
		{"coolify", "Coolify (App Deployment)"},
		{"ollama", "Ollama (Local LLM)"},
		{"opencode", "OpenCode.ai (AI Coding)"},
		{"openwebui", "Open WebUI (AI Chat UI)"},
		{"code-server", "Code-Server (VS Code in Browser)"},
		{"mattermost", "Mattermost (Team Chat)"},
		{"outline", "Outline (Knowledge Base)"},
		{"plane", "Plane (Project Management)"},
		{"api", "Workspace API"},
		{"dashboard", "Workspace Dashboard"},
		{"restic", "Restic (Backup)"},
		{"grafana", "Grafana (Monitoring)"},
		{"cloudflare", "Cloudflare Tunnel"},
	}

	services := make(map[string]bool)
	for _, svc := range allServices {
		fmt.Printf("Install %s? [Y/n]: ", svc.label)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		services[svc.name] = answer == "" || answer == "y" || answer == "yes"
	}

	cfg := config.DefaultConfig()
	cfg.Workspace.Domain = domain
	cfg.Services.Authentik = services["authentik"]
	cfg.Services.Gitea = services["gitea"]
	cfg.Services.Coolify = services["coolify"]
	cfg.Services.Ollama = services["ollama"]
	cfg.Services.OpenCode = services["opencode"]
	cfg.Services.OpenWebUI = services["openwebui"]
	cfg.Services.CodeServer = services["code-server"]
	cfg.Services.Mattermost = services["mattermost"]
	cfg.Services.Outline = services["outline"]
	cfg.Services.Plane = services["plane"]
	cfg.Services.Restic = services["restic"]
	cfg.Services.Grafana = services["grafana"]
	cfg.Services.Cloudflare = services["cloudflare"]
	cfg.Services.Api = services["api"]
	cfg.Services.Dashboard = services["dashboard"]

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	outputPath := outputFlag
	if outputPath == "" {
		outputPath = configPath
	}
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("\n⚠ %s already exists. Overwrite? [y/N]: ", outputPath)
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			fmt.Println("Init cancelled.")
			return nil
		}
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	if info, err := os.Stat(outputPath); err == nil && info.IsDir() {
		if err := os.Remove(outputPath); err != nil {
			return fmt.Errorf("remove existing directory at output path: %w", err)
		}
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	fmt.Printf("\n✅ Configuration written to %s\n", outputPath)
	fmt.Println("\nRun 'workspace plan' to see what needs to be installed.")
	fmt.Println("Run 'workspace install' to apply.")

	return nil
}

func parseServices(input string) []string {
	if input == "" {
		return config.DefaultConfig().Services.EnabledList()
	}
	var result []string
	for _, s := range strings.Split(input, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&autoMode, "auto", false, "Non-interactive mode")
	initCmd.Flags().StringVar(&domainFlag, "domain", "", "Workspace domain")
	initCmd.Flags().StringVar(&servicesFlag, "services", "", "Comma-separated services")
	initCmd.Flags().StringVar(&outputFlag, "output", "", "Output path for config file")
}
