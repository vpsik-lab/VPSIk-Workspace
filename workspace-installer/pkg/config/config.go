package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type WorkspaceConfig struct {
	Workspace Workspace     `yaml:"workspace"`
	Services  Services      `yaml:"services"`
	Backup    *BackupConfig   `yaml:"backup,omitempty"`
	Network   *NetworkConfig  `yaml:"network,omitempty"`
	System    *SystemConfig   `yaml:"system,omitempty"`
}

type Workspace struct {
	Domain string `yaml:"domain"`
}

type Services struct {
	Authentik  bool `yaml:"authentik"`
	Gitea      bool `yaml:"gitea"`
	Coolify    bool `yaml:"coolify"`
	Ollama     bool `yaml:"ollama"`
	OpenCode   bool `yaml:"opencode"`
	OpenWebUI  bool `yaml:"openwebui"`
	Restic     bool `yaml:"restic"`
	Traefik    bool `yaml:"traefik"`
	Postgres   bool `yaml:"postgres"`
	Grafana    bool `yaml:"grafana"`
	Prometheus bool `yaml:"prometheus"`
	CodeServer bool `yaml:"code-server"`
	Plane      bool `yaml:"plane"`
	Outline    bool `yaml:"outline"`
	Mattermost bool `yaml:"mattermost"`
	Cloudflare bool `yaml:"cloudflare"`
	Redis      bool `yaml:"redis"`
	Api        bool `yaml:"api"`
	Dashboard  bool `yaml:"dashboard"`
}

type BackupConfig struct {
	Repository string `yaml:"repository"`
	Schedule   string `yaml:"schedule"`
	KeepPolicy string `yaml:"keep-policy"`
}

type NetworkConfig struct {
	Name      string `yaml:"name"`
	ProxyPort int    `yaml:"proxy-port"`
	UseTunnel bool   `yaml:"use-tunnel"`
}

type SystemConfig struct {
	InstallPath string `yaml:"install-path"`
}

func (s Services) EnabledList() []string {
	var enabled []string
	if s.Authentik  { enabled = append(enabled, "authentik") }
	if s.Gitea      { enabled = append(enabled, "gitea") }
	if s.Coolify    { enabled = append(enabled, "coolify") }
	if s.Ollama     { enabled = append(enabled, "ollama") }
	if s.OpenCode   { enabled = append(enabled, "opencode") }
	if s.OpenWebUI  { enabled = append(enabled, "openwebui") }
	if s.Restic     { enabled = append(enabled, "restic") }
	if s.Traefik    { enabled = append(enabled, "traefik") }
	if s.Postgres   { enabled = append(enabled, "postgres") }
	if s.Grafana    { enabled = append(enabled, "grafana") }
	if s.Prometheus { enabled = append(enabled, "prometheus") }
	if s.CodeServer { enabled = append(enabled, "code-server") }
	if s.Plane      { enabled = append(enabled, "plane") }
	if s.Outline    { enabled = append(enabled, "outline") }
	if s.Mattermost { enabled = append(enabled, "mattermost") }
	if s.Cloudflare { enabled = append(enabled, "cloudflare") }
	if s.Redis      { enabled = append(enabled, "redis") }
	if s.Api        { enabled = append(enabled, "api") }
	if s.Dashboard  { enabled = append(enabled, "dashboard") }
	return enabled
}

func DefaultConfig() *WorkspaceConfig {
	return &WorkspaceConfig{
		Workspace: Workspace{
			Domain: "workspace.vpsik.com",
		},
		Services: Services{
			Traefik:    true,
			Postgres:   true,
			Redis:      true,
			Authentik:  true,
			Gitea:      true,
			Ollama:     true,
			OpenWebUI:  true,
			CodeServer: true,
			Mattermost: false,
			Outline:    false,
			Plane:      false,
			Coolify:    true,
			OpenCode:   false,
			Restic:     true,
			Grafana:    true,
			Prometheus: true,
			Cloudflare: false,
			Api:        true,
			Dashboard:  true,
		},
		System: &SystemConfig{
			InstallPath: "/opt/workspace",
		},
		Network: &NetworkConfig{
			Name:      "workspace_net",
			ProxyPort: 443,
			UseTunnel: false,
		},
		Backup: &BackupConfig{
			Repository: "/opt/workspace/backups",
			Schedule:   "0 3 * * *",
			KeepPolicy: "7d 4w 12m",
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

func (c *WorkspaceConfig) NetworkName() string {
	if c.Network != nil && c.Network.Name != "" {
		return c.Network.Name
	}
	return "workspace_net"
}

func (c *WorkspaceConfig) InstallPath() string {
	if c.System != nil && c.System.InstallPath != "" {
		return c.System.InstallPath
	}
	return "/opt/workspace"
}

func ConfigFilename(path string) string {
	base := strings.TrimSuffix(path, ".yaml")
	base = strings.TrimSuffix(base, ".yml")
	return base + ".yaml"
}
