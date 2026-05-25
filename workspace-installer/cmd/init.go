package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a workspace configuration file interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("🚀 VPSIk Workspace — Interactive Setup")
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
			{"restic", "Restic (Backup)"},
		}

		services := make(map[string]bool)
		for _, svc := range allServices {
			fmt.Printf("Install %s? [Y/n]: ", svc.label)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			services[svc.name] = answer == "" || answer == "y" || answer == "yes"
		}

		cfg := struct {
			Workspace struct {
				Domain string `yaml:"domain"`
			} `yaml:"workspace"`
			Services map[string]bool `yaml:"services"`
		}{}
		cfg.Workspace.Domain = domain
		cfg.Services = services

		data, err := yaml.Marshal(&cfg)
		if err != nil {
			return fmt.Errorf("marshal config: %w", err)
		}

		outputPath := configPath
		if outputPath == "workspace.yaml" {
			if _, err := os.Stat(outputPath); err == nil {
				fmt.Printf("\n⚠ %s already exists. Overwrite? [y/N]: ", outputPath)
				answer, _ := reader.ReadString('\n')
				answer = strings.TrimSpace(strings.ToLower(answer))
				if answer != "y" && answer != "yes" {
					fmt.Println("Init cancelled.")
					return nil
				}
			}
		}

		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			return fmt.Errorf("write config: %w", err)
		}

		fmt.Printf("\n✅ Configuration written to %s\n", outputPath)
		fmt.Println("\nRun 'vpsik plan' to see what needs to be installed.")
		fmt.Println("Run 'vpsik install' to apply.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
