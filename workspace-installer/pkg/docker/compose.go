package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type ServiceCompose struct {
	Name     string
	Image    string
	Ports    []string
	Volumes  []string
	EnvVars  map[string]string
	Network  string
}

var ServiceTemplates = map[string]*ServiceCompose{
	"authentik": {
		Name:  "authentik",
		Image: "ghcr.io/goauthentik/server:latest",
		Ports: []string{"9000:9000"},
		Volumes: []string{
			"authentik-media:/media",
			"authentik-certs:/certs",
		},
		EnvVars: map[string]string{
			"AUTHENTIK_SECRET_KEY": "${AUTHENTIK_SECRET_KEY}",
			"AUTHENTIK_BOOT_PASSWORD": "${AUTHENTIK_BOOT_PASSWORD}",
		},
		Network: "vpsik",
	},
	"gitea": {
		Name:  "gitea",
		Image: "gitea/gitea:latest",
		Ports: []string{"3000:3000", "22:22"},
		Volumes: []string{
			"gitea-data:/data",
		},
		EnvVars: map[string]string{
			"USER_UID": "1000",
			"USER_GID": "1000",
			"GITEA__server__DOMAIN": "${VPSIK_DOMAIN}",
		},
		Network: "vpsik",
	},
	"coolify": {
		Name:  "coolify",
		Image: "coollabsio/coolify:latest",
		Ports: []string{"8000:8000"},
		Volumes: []string{
			"coolify-data:/data",
			"/var/run/docker.sock:/var/run/docker.sock",
		},
		EnvVars: map[string]string{
			"COOLIFY_APP_ID": "${COOLIFY_APP_ID}",
			"COOLIFY_SECRET_KEY": "${COOLIFY_SECRET_KEY}",
		},
		Network: "vpsik",
	},
	"ollama": {
		Name:  "ollama",
		Image: "ollama/ollama:latest",
		Ports: []string{"11434:11434"},
		Volumes: []string{
			"ollama-data:/root/.ollama",
		},
		Network: "vpsik",
	},
	"opencode": {
		Name:  "opencode",
		Image: "opencodeai/opencode:latest",
		Ports: []string{"30081:30081"},
		Volumes: []string{
			"opencode-data:/data",
		},
		Network: "vpsik",
	},
	"openwebui": {
		Name:  "openwebui",
		Image: "ghcr.io/open-webui/open-webui:main",
		Ports: []string{"3001:8080"},
		Volumes: []string{
			"openwebui-data:/app/backend/data",
		},
		EnvVars: map[string]string{
			"OLLAMA_BASE_URL": "http://ollama:11434",
		},
		Network: "vpsik",
	},
	"postgres": {
		Name:  "postgres",
		Image: "postgres:16-alpine",
		Ports: []string{"5432:5432"},
		Volumes: []string{
			"postgres-data:/var/lib/postgresql/data",
		},
		EnvVars: map[string]string{
			"POSTGRES_USER":     "${POSTGRES_USER:-vpsik}",
			"POSTGRES_PASSWORD": "${POSTGRES_PASSWORD:-vpsik}",
			"POSTGRES_DB":       "${POSTGRES_DB:-vpsik}",
		},
		Network: "vpsik",
	},
	"traefik": {
		Name:  "traefik",
		Image: "traefik:v3.0",
		Ports: []string{"80:80", "443:443"},
		Volumes: []string{
			"/var/run/docker.sock:/var/run/docker.sock",
			"traefik-data:/etc/traefik",
		},
		Network: "vpsik",
	},
	"grafana": {
		Name:  "grafana",
		Image: "grafana/grafana:latest",
		Ports: []string{"3002:3000"},
		Volumes: []string{
			"grafana-data:/var/lib/grafana",
		},
		EnvVars: map[string]string{
			"GF_SECURITY_ADMIN_USER":     "${GRAFANA_USER:-admin}",
			"GF_SECURITY_ADMIN_PASSWORD": "${GRAFANA_PASSWORD:-admin}",
		},
		Network: "vpsik",
	},
	"prometheus": {
		Name:  "prometheus",
		Image: "prom/prometheus:latest",
		Ports: []string{"9090:9090"},
		Volumes: []string{
			"prometheus-data:/prometheus",
		},
		Network: "vpsik",
	},
	"restic": {
		Name:  "restic",
		Image: "restic/restic:latest",
		Volumes: []string{
			"restic-data:/data",
		},
		Network: "vpsik",
	},
}

const composeTemplate = `version: "3.8"

networks:
  {{.Network}}:
    external: true
    name: {{.Network}}

volumes:
{{range .Volumes}}
  {{.}}:
{{end}}

services:
{{range .Services}}
  {{.Name}}:
    image: {{.Image}}
    container_name: {{.Name}}
    restart: unless-stopped
    networks:
      - {{.Network}}
{{if .Ports}}
    ports:
{{range .Ports}}
      - "{{.}}"
{{end}}
{{end}}
{{if .Volumes}}
    volumes:
{{range .Volumes}}
      - "{{.}}"
{{end}}
{{end}}
{{if .EnvVars}}
    environment:
{{range $key, $val := .EnvVars}}
      {{$key}}: {{$val}}
{{end}}
{{end}}
{{end}}
`

type composeData struct {
	Network  string
	Volumes  []string
	Services []*ServiceCompose
}

func GenerateComposeFile(services []string, domain string, outputPath string) error {
	volumes := make(map[string]bool)
	var svcDefs []*ServiceCompose

	for _, name := range services {
		if tpl, ok := ServiceTemplates[name]; ok {
			svc := *tpl
			svc.EnvVars["VPSIK_DOMAIN"] = domain
			svcDefs = append(svcDefs, &svc)
			for _, vol := range tpl.Volumes {
				if vol[0] != '/' {
					volumes[vol] = true
				}
			}
		}
	}

	var volList []string
	for v := range volumes {
		volList = append(volList, v)
	}

	data := composeData{
		Network:  "vpsik",
		Volumes:  volList,
		Services: svcDefs,
	}

	tmpl, err := template.New("compose").Parse(composeTemplate)
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create compose file: %w", err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

func Deploy(composePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	absPath, err := filepath.Abs(composePath)
	if err != nil {
		return fmt.Errorf("abs path: %w", err)
	}

	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", absPath, "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose up: %w", err)
	}

	return nil
}

func EnsureNetwork(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "network", "inspect", name)
	if err := cmd.Run(); err == nil {
		return nil
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	cmd = exec.CommandContext(ctx2, "docker", "network", "create", name)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("create network %s: %w", name, err)
	}

	return nil
}

func Down(composePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	absPath, err := filepath.Abs(composePath)
	if err != nil {
		return fmt.Errorf("abs path: %w", err)
	}

	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", absPath, "down")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RemoveVolumes(composePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	absPath, err := filepath.Abs(composePath)
	if err != nil {
		return fmt.Errorf("abs path: %w", err)
	}

	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", absPath, "down", "-v")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func RemoveNetwork(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "network", "rm", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func PullImages(services []string) error {
	seen := make(map[string]bool)
	for _, name := range services {
		tpl, ok := ServiceTemplates[name]
		if !ok {
			continue
		}
		if seen[tpl.Image] {
			continue
		}
		seen[tpl.Image] = true

		fmt.Printf("  Pulling %s...\n", tpl.Image)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		cmd := exec.CommandContext(ctx, "docker", "pull", tpl.Image)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		cancel()
		if err != nil {
			return fmt.Errorf("pull %s: %w", tpl.Image, err)
		}
	}
	return nil
}

func RecreateServices(configPath, domain string) error {
	composeDir := filepath.Dir(configPath)
	composePath := filepath.Join(composeDir, "docker-compose.generated.yml")

	absPath, err := filepath.Abs(composePath)
	if err != nil {
		return fmt.Errorf("abs path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("no generated compose file found at %s", absPath)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "compose", "-f", absPath, "up", "-d", "--force-recreate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func GenerateEnvFile(services []string, domain string, outputPath string) error {
	envVars := make(map[string]string)
	envVars["VPSIK_DOMAIN"] = domain

	for _, name := range services {
		tpl, ok := ServiceTemplates[name]
		if !ok {
			continue
		}
		for k, v := range tpl.EnvVars {
			if _, exists := envVars[k]; !exists {
				envVars[k] = v
			}
		}
	}

	var sb strings.Builder
	for k, v := range envVars {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return os.WriteFile(outputPath, []byte(sb.String()), 0644)
}
