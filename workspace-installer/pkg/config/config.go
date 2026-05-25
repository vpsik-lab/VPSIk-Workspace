package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type WorkspaceConfig struct {
	Workspace Workspace `yaml:"workspace"`
	Services  Services  `yaml:"services"`
}

type Workspace struct {
	Domain string `yaml:"domain"`
}

type Services struct {
	Authentik bool `yaml:"authentik"`
	Gitea     bool `yaml:"gitea"`
	Coolify   bool `yaml:"coolify"`
	Ollama    bool `yaml:"ollama"`
	OpenCode  bool `yaml:"opencode"`
	OpenWebUI bool `yaml:"openwebui"`
	Restic    bool `yaml:"restic"`
}

func (s Services) EnabledList() []string {
	var enabled []string
	if s.Authentik {
		enabled = append(enabled, "authentik")
	}
	if s.Gitea {
		enabled = append(enabled, "gitea")
	}
	if s.Coolify {
		enabled = append(enabled, "coolify")
	}
	if s.Ollama {
		enabled = append(enabled, "ollama")
	}
	if s.OpenCode {
		enabled = append(enabled, "opencode")
	}
	if s.OpenWebUI {
		enabled = append(enabled, "openwebui")
	}
	if s.Restic {
		enabled = append(enabled, "restic")
	}
	return enabled
}

func DefaultConfig() *WorkspaceConfig {
	return &WorkspaceConfig{
		Workspace: Workspace{
			Domain: "workspace.vpsik.com",
		},
		Services: Services{
			Authentik: true,
			Gitea:     true,
			Coolify:   true,
			Ollama:    true,
			OpenCode:  true,
			OpenWebUI: true,
			Restic:    true,
		},
	}
}

func Load(path string) (*WorkspaceConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file %s not found", path)
		}
		return nil, fmt.Errorf("read config: %w", err)
	}

	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	return cfg, nil
}
