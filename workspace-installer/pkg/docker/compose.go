package docker

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v3"
)

type ServiceCompose struct {
	Name        string
	Image       string
	BuildCtx    string
	Ports       []string
	Expose      []string
	Volumes     []string
	EnvVars     map[string]string
	Network     string
	TraefikHost string
	TraefikPort string
	DependsOn   []string
	Replicas    int
	ExtraLabels map[string]string
	Command     []string
}

var ServiceTemplates = map[string]*ServiceCompose{
	"traefik": {
		Name:  "traefik",
		Image: "traefik:v3.0",
		Ports: []string{"80:80", "443:443"},
		Volumes: []string{
			"/var/run/docker.sock:/var/run/docker.sock:ro",
			"traefik-data:/etc/traefik",
		},
		Network: "workspace_net",
	},
	"postgres": {
		Name:   "postgres",
		Image:  "postgres:16-alpine",
		Expose: []string{"5432"},
		Volumes: []string{
			"postgres-data:/var/lib/postgresql/data",
		},
		EnvVars: map[string]string{
			"POSTGRES_USER":     "${POSTGRES_USER:-vpsik}",
			"POSTGRES_PASSWORD": "${POSTGRES_PASSWORD}",
			"POSTGRES_DB":       "${POSTGRES_DB:-vpsik}",
		},
		Network: "workspace_net",
	},
	"redis": {
		Name:    "redis",
		Image:   "redis:7-alpine",
		Expose:  []string{"6379"},
		Volumes: []string{"redis-data:/data"},
		EnvVars: map[string]string{
			"REDIS_PASSWORD": "${REDIS_PASSWORD}",
		},
		Network:  "workspace_net",
	},
	"authentik": {
		Name:        "authentik-server",
		Image:       "ghcr.io/goauthentik/server:latest",
		Expose:      []string{"9000"},
		Network:     "workspace_net",
		TraefikHost: "auth.${VPSIK_DOMAIN}",
		TraefikPort: "9000",
		DependsOn:   []string{"postgres", "redis"},
		EnvVars: map[string]string{
			"AUTHENTIK_SECRET_KEY":     "${AUTHENTIK_SECRET_KEY}",
			"AUTHENTIK_BOOT_PASSWORD":  "${AUTHENTIK_BOOT_PASSWORD}",
			"AUTHENTIK_REDIS__HOST":    "redis",
			"AUTHENTIK_POSTGRESQL__HOST": "postgres",
			"AUTHENTIK_POSTGRESQL__USER": "${POSTGRES_USER:-vpsik}",
			"AUTHENTIK_POSTGRESQL__PASSWORD": "${POSTGRES_PASSWORD:-vpsik}",
			"AUTHENTIK_POSTGRESQL__NAME": "${POSTGRES_DB:-vpsik}",
		},
	},
	"gitea": {
		Name:        "gitea",
		Image:       "gitea/gitea:latest",
		Expose:      []string{"3000"},
		Volumes:     []string{"gitea-data:/data"},
		Network:     "workspace_net",
		TraefikHost: "git.${VPSIK_DOMAIN}",
		TraefikPort: "3000",
		EnvVars: map[string]string{
			"USER_UID":                  "1000",
			"USER_GID":                  "1000",
			"GITEA__server__DOMAIN":     "git.${VPSIK_DOMAIN}",
			"GITEA__server__ROOT_URL":   "https://git.${VPSIK_DOMAIN}",
			"GITEA__database__DB_TYPE":  "postgres",
			"GITEA__database__HOST":     "postgres:5432",
			"GITEA__database__NAME":     "${POSTGRES_DB:-vpsik}",
			"GITEA__database__USER":     "${POSTGRES_USER:-vpsik}",
			"GITEA__database__PASSWD":   "${POSTGRES_PASSWORD:-vpsik}",
		},
		DependsOn: []string{"postgres"},
	},
	"coolify": {
		Name:        "coolify",
		Image:       "coollabsio/coolify:latest",
		Expose:      []string{"3000"},
		Volumes:     []string{"coolify-data:/data", "/var/run/docker.sock:/var/run/docker.sock"},
		Network:     "workspace_net",
		TraefikHost: "coolify.${VPSIK_DOMAIN}",
		TraefikPort: "3000",
		EnvVars: map[string]string{
			"COOLIFY_APP_ID":    "${COOLIFY_APP_ID}",
			"COOLIFY_SECRET_KEY": "${COOLIFY_SECRET_KEY}",
		},
	},
	"ollama": {
		Name:        "ollama",
		Image:       "ollama/ollama:latest",
		Expose:      []string{"11434"},
		Volumes:     []string{"ollama-data:/root/.ollama"},
		Network:     "workspace_net",
		TraefikHost: "ollama.${VPSIK_DOMAIN}",
		TraefikPort: "11434",
	},
	"opencode": {
		Name:        "opencode",
		Image:       "opencodeai/opencode:latest",
		Expose:      []string{"30081"},
		Volumes:     []string{"opencode-data:/data"},
		Network:     "workspace_net",
		TraefikHost: "opencode.${VPSIK_DOMAIN}",
		TraefikPort: "30081",
	},
	"openwebui": {
		Name:        "openwebui",
		Image:       "ghcr.io/open-webui/open-webui:main",
		Expose:      []string{"8080"},
		Volumes:     []string{"openwebui-data:/app/backend/data"},
		Network:     "workspace_net",
		TraefikHost: "chat.${VPSIK_DOMAIN}",
		TraefikPort: "8080",
		EnvVars: map[string]string{
			"OLLAMA_BASE_URL": "http://ollama:11434",
		},
	},
	"code-server": {
		Name:        "codeserver",
		Image:       "codercom/code-server:latest",
		Expose:      []string{"8443"},
		Volumes:     []string{"codeserver-data:/home/coder"},
		Network:     "workspace_net",
		TraefikHost: "code.${VPSIK_DOMAIN}",
		TraefikPort: "8443",
		EnvVars: map[string]string{
			"PASSWORD": "${CODESERVER_PASSWORD}",
			"SUDO_PASSWORD": "${CODESERVER_PASSWORD}",
		},
	},
	"plane": {
		Name:        "plane",
		Image:       "makeplane/plane-app:latest",
		Expose:      []string{"8080"},
		Network:     "workspace_net",
		TraefikHost: "plane.${VPSIK_DOMAIN}",
		TraefikPort: "8080",
		DependsOn:   []string{"postgres", "redis"},
	},
	"outline": {
		Name:        "outline",
		Image:       "outlinewiki/outline:latest",
		Expose:      []string{"3000"},
		Volumes:     []string{"outline-data:/var/lib/outline/data"},
		Network:     "workspace_net",
		TraefikHost: "docs.${VPSIK_DOMAIN}",
		TraefikPort: "3000",
		DependsOn:   []string{"postgres", "redis"},
		EnvVars: map[string]string{
			"SECRET_KEY":     "${OUTLINE_SECRET_KEY}",
			"PGSSLMODE":      "disable",
			"DATABASE_URL":   "postgres://${POSTGRES_USER:-vpsik}:${POSTGRES_PASSWORD:-vpsik}@postgres:5432/${POSTGRES_DB:-vpsik}",
			"REDIS_URL":      "redis://:${REDIS_PASSWORD:-vpsik}@redis:6379",
		},
	},
	"mattermost": {
		Name:        "mattermost",
		Image:       "mattermost/mattermost-team-edition:latest",
		Expose:      []string{"8065"},
		Volumes:     []string{"mattermost-data:/mattermost/data"},
		Network:     "workspace_net",
		TraefikHost: "mattermost.${VPSIK_DOMAIN}",
		TraefikPort: "8065",
		DependsOn:   []string{"postgres"},
		EnvVars: map[string]string{
			"MM_USERNAME":   "mmuser",
			"MM_PASSWORD":   "${POSTGRES_PASSWORD:-vpsik}",
			"MM_DBNAME":     "${POSTGRES_DB:-vpsik}",
			"MM_SQLSETTINGS_DATASOURCE": "postgres://mmuser:${POSTGRES_PASSWORD:-vpsik}@postgres:5432/${POSTGRES_DB:-vpsik}?sslmode=disable&connect_timeout=10",
		},
	},
	"grafana": {
		Name:        "grafana",
		Image:       "grafana/grafana:latest",
		Expose:      []string{"3000"},
		Volumes:     []string{"grafana-data:/var/lib/grafana"},
		Network:     "workspace_net",
		TraefikHost: "metrics.${VPSIK_DOMAIN}",
		TraefikPort: "3000",
		EnvVars: map[string]string{
			"GF_SECURITY_ADMIN_USER":     "${GRAFANA_USER:-admin}",
			"GF_SECURITY_ADMIN_PASSWORD": "${GRAFANA_PASSWORD}",
		},
	},
	"prometheus": {
		Name:    "prometheus",
		Image:   "prom/prometheus:latest",
		Expose:  []string{"9090"},
		Volumes: []string{"prometheus-data:/prometheus"},
		Network: "workspace_net",
	},
	"restic": {
		Name:    "restic",
		Image:   "restic/restic:latest",
		Volumes: []string{"restic-data:/data"},
		Network: "workspace_net",
	},
	"dashboard": {
		Name:        "dashboard",
		Image:       "vpsik-dashboard:latest",
		Expose:      []string{"3000"},
		Ports:       []string{"3000:3000"},
		Network:     "workspace_net",
		TraefikHost: "${VPSIK_DOMAIN}",
		TraefikPort: "3000",
		DependsOn:   []string{"api"},
	},
	"api": {
		Name:        "api",
		Image:       "vpsik-api:latest",
		Expose:      []string{"8081"},
		Ports:       []string{"8081:8081"},
		Volumes:     []string{"/opt/workspace/api.yaml:/etc/vpsik/api.yaml:ro"},
		Network:     "workspace_net",
		TraefikHost: "api.${VPSIK_DOMAIN}",
		TraefikPort: "8081",
		Command:     []string{"/etc/vpsik/api.yaml"},
	},
	"cloudflare": {
		Name:  "cloudflared",
		Image: "cloudflare/cloudflared:latest",
		EnvVars: map[string]string{
			"TUNNEL_TOKEN": "${CF_TUNNEL_TOKEN}",
		},
		Network: "workspace_net",
		ExtraLabels: map[string]string{
			"com.docker.compose.managed": "vpsik",
		},
	},
}

type composeService struct {
	Image       string            `yaml:"image,omitempty"`
	ContainerName string          `yaml:"container_name,omitempty"`
	Restart     string            `yaml:"restart,omitempty"`
	Ports       []string          `yaml:"ports,omitempty"`
	Expose      []string          `yaml:"expose,omitempty"`
	Volumes     []string          `yaml:"volumes,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Networks    []string          `yaml:"networks,omitempty"`
	DependsOn   []string          `yaml:"depends_on,omitempty"`
	Command     []string          `yaml:"command,omitempty"`
}

type composeFile struct {
	Version  string                    `yaml:"version,omitempty"`
	Networks map[string]networkDef     `yaml:"networks"`
	Volumes  map[string]map[string]interface{} `yaml:"volumes"`
	Services map[string]composeService `yaml:"services"`
}

type networkDef struct {
	Name       string `yaml:"name,omitempty"`
	External   bool   `yaml:"external,omitempty"`
	Driver     string `yaml:"driver,omitempty"`
}

func GenerateTraefikLabels(svc *ServiceCompose) map[string]string {
	labels := make(map[string]string)
	name := svc.Name

	labels["traefik.enable"] = "true"

	if svc.TraefikHost != "" {
		labels[fmt.Sprintf("traefik.http.routers.%s.rule", name)] = fmt.Sprintf("Host(`%s`)", svc.TraefikHost)
		labels[fmt.Sprintf("traefik.http.routers.%s.entrypoints", name)] = "websecure"
	}

	if svc.TraefikPort != "" {
		labels[fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", name)] = svc.TraefikPort
	}

	for k, v := range svc.ExtraLabels {
		labels[k] = v
	}

	return labels
}

func GenerateComposeFile(services []string, network string, domain string, outputPath string) error {
	volumes := make(map[string]bool)
	svcMap := make(map[string]composeService)
	svcSet := make(map[string]bool, len(services))
	for _, s := range services {
		svcSet[s] = true
	}

	for _, name := range services {
		tpl, ok := ServiceTemplates[name]
		if !ok {
			continue
		}

		svc := *tpl
		svc.EnvVars = make(map[string]string, len(tpl.EnvVars))
		for k, v := range tpl.EnvVars {
			svc.EnvVars[k] = v
		}

		var deps []string
		for _, d := range svc.DependsOn {
			if svcSet[d] {
				deps = append(deps, d)
			}
		}

		cf := composeService{
			Image:         svc.Image,
			ContainerName: svc.Name,
			Restart:       "unless-stopped",
			Expose:        svc.Expose,
			Ports:         svc.Ports,
			Volumes:       svc.Volumes,
			Environment:   svc.EnvVars,
			Networks:      []string{network},
			DependsOn:     deps,
			Command:       svc.Command,
		}

		labels := GenerateTraefikLabels(&svc)
		if len(labels) > 0 {
			resolvedLabels := make(map[string]string)
			for k, v := range labels {
				resolvedLabels[k] = os.Expand(v, func(key string) string {
					if val := os.Getenv(key); val != "" {
						return val
					}
					return "${" + key + "}"
				})
			}
			cf.Labels = resolvedLabels
		}

		svcMap[svc.Name] = cf

		for _, vol := range tpl.Volumes {
			if vol == "" {
				continue
			}
			vname := vol
			if idx := strings.IndexByte(vol, ':'); idx >= 0 {
				vname = vol[:idx]
			}
			if vname != "" && vname[0] != '/' && vname[0] != '.' && vname[0] != '~' {
				volumes[vname] = true
			}
		}
	}

	volMap := make(map[string]map[string]interface{})
	var volKeys []string
	for v := range volumes {
		volKeys = append(volKeys, v)
	}
	sort.Strings(volKeys)
	for _, v := range volKeys {
		volMap[v] = map[string]interface{}{}
	}

	compose := composeFile{
		Networks: map[string]networkDef{
			network: {Name: network, External: true},
		},
		Volumes:  volMap,
		Services: svcMap,
	}

	data, err := yaml.Marshal(compose)
	if err != nil {
		return fmt.Errorf("marshal compose: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	return os.WriteFile(outputPath, data, 0644)
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}

func randomHex(length int) string {
	const charset = "0123456789abcdef"
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
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
				if v == "" || (strings.HasPrefix(v, "${") && strings.HasSuffix(v, "}") && !strings.Contains(v, ":-")) {
					switch k {
					case "POSTGRES_PASSWORD":
						envVars[k] = randomString(24)
					case "REDIS_PASSWORD":
						envVars[k] = randomString(24)
					case "AUTHENTIK_SECRET_KEY":
						envVars[k] = randomHex(64)
					case "AUTHENTIK_BOOT_PASSWORD":
						envVars[k] = randomString(24)
					case "GRAFANA_PASSWORD":
						envVars[k] = randomString(24)
					case "CODESERVER_PASSWORD":
						envVars[k] = randomString(24)
					case "RESTIC_PASSWORD":
						envVars[k] = randomString(24)
					default:
						envVars[k] = v
					}
				} else {
					envVars[k] = v
				}
			}
		}
	}

	var lines []string
	for k, v := range envVars {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(lines)

	content := "# VPSIk Workspace — Auto-generated environment\n"
	content += "# Generated for domain: " + domain + "\n\n"
	content += strings.Join(lines, "\n") + "\n"

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	return os.WriteFile(outputPath, []byte(content), 0644)
}

func GenerateAPIConfig(services []string, domain string, outputPath string) (string, error) {
	jwtSecret := randomHex(64)
	adminPassword := randomString(16)
	adminHash, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	type apiUser struct {
		Username     string `yaml:"username"`
		PasswordHash string `yaml:"password_hash"`
	}

	type apiAuth struct {
		JWTSecret string    `yaml:"jwt_secret"`
		Users     []apiUser `yaml:"users"`
	}

	type svcEndpoint struct {
		URL   string `yaml:"url"`
		Token string `yaml:"token,omitempty"`
	}

	type resticCfg struct {
		Binary   string `yaml:"binary"`
		RepoURL  string `yaml:"repo_url"`
		Password string `yaml:"password"`
	}

	type apiServices struct {
		Gitea      svcEndpoint `yaml:"gitea"`
		Coolify    svcEndpoint `yaml:"coolify"`
		Ollama     svcEndpoint `yaml:"ollama"`
		OpenCode   svcEndpoint `yaml:"opencode"`
		Restic     resticCfg   `yaml:"restic"`
		CodeServer svcEndpoint `yaml:"code-server"`
		Plane      svcEndpoint `yaml:"plane,omitempty"`
		Outline    svcEndpoint `yaml:"outline,omitempty"`
		Mattermost svcEndpoint `yaml:"mattermost,omitempty"`
	}

	svcs := apiServices{
		Gitea:      svcEndpoint{URL: "http://gitea:3000"},
		Coolify:    svcEndpoint{URL: "http://coolify:3000"},
		Ollama:     svcEndpoint{URL: "http://ollama:11434"},
		OpenCode:   svcEndpoint{URL: "", Token: ""},
		Restic:     resticCfg{Binary: "restic", RepoURL: "", Password: ""},
		CodeServer: svcEndpoint{URL: "http://codeserver:8443"},
	}

	cfg := struct {
		Server   map[string]interface{} `yaml:"server"`
		Auth     apiAuth                `yaml:"auth"`
		Services apiServices            `yaml:"services"`
	}{
		Server: map[string]interface{}{
			"port": 8081,
			"host": "0.0.0.0",
		},
		Auth: apiAuth{
			JWTSecret: jwtSecret,
			Users: []apiUser{
				{Username: "admin", PasswordHash: string(adminHash)},
			},
		},
		Services: svcs,
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return "", fmt.Errorf("marshal api config: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", fmt.Errorf("create output dir: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return "", fmt.Errorf("write api config: %w", err)
	}

	return adminPassword, nil
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

		// Skip pull for locally-built images
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		check := exec.CommandContext(ctx, "docker", "image", "inspect", tpl.Image)
		if err := check.Run(); err == nil {
			cancel()
			fmt.Printf("  Using local image %s\n", tpl.Image)
			continue
		}
		cancel()

		fmt.Printf("  Pulling %s...\n", tpl.Image)
		ctx, cancel = context.WithTimeout(context.Background(), 300*time.Second)
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
