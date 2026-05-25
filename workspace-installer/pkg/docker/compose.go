package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

var serviceTemplates = map[string]*ServiceCompose{
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
		Ports: []string{"3000:8080"},
		Volumes: []string{
			"openwebui-data:/app/backend/data",
		},
		EnvVars: map[string]string{
			"OLLAMA_BASE_URL": "http://ollama:11434",
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
		if tpl, ok := serviceTemplates[name]; ok {
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
