# Phase 1 — Infrastructure Bootstrap Improvements

## التحليل الكامل لما تم وما ينقص

### ✅ موجود حالياً
- CLI tool `vpsik` مع أوامر: `status`, `plan`, `install`
- Config system (YAML) مع 7 خدمات
- Environment Scanner (Docker, containers, networks, ports)
- Service Detector لـ 7 خدمات
- State Engine + Plan Engine (Diff desired vs current)
- Docker Compose Generator + Deploy

### ❌ التحسينات المطلوبة

---

## المهام المطلوب تنفيذها

### 1. إضافة `vpsik init` — ملف: `cmd/init.go`
```go
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
```

### 2. إضافة `vpsik uninstall` — ملف: `cmd/uninstall.go`
```go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/config"
	"github.com/vpsik/workspace-installer/pkg/docker"
)

var (
	removeVolumes bool
	removeNetwork bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Remove all deployed services",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		composeDir := filepath.Dir(configPath)
		composePath := filepath.Join(composeDir, "docker-compose.generated.yml")

		if _, err := os.Stat(composePath); os.IsNotExist(err) {
			fmt.Println("No generated compose file found. Nothing to uninstall.")
			return nil
		}

		fmt.Print("⚠ This will stop and remove all services. Continue? (yes/no): ")
		var response string
		fmt.Scanln(&response)
		if response != "yes" {
			fmt.Println("Uninstall cancelled.")
			return nil
		}

		fmt.Println("\n🛑 Stopping and removing services...")
		if err := docker.Down(composePath); err != nil {
			return fmt.Errorf("docker compose down: %w", err)
		}

		if removeVolumes {
			fmt.Println("🗑 Removing volumes...")
			if err := docker.RemoveVolumes(composePath); err != nil {
				return fmt.Errorf("remove volumes: %w", err)
			}
		}

		if removeNetwork {
			fmt.Println("🌐 Removing network 'vpsik'...")
			if err := docker.RemoveNetwork("vpsik"); err != nil {
				return fmt.Errorf("remove network: %w", err)
			}
		}

		fmt.Println("\n✅ Uninstall complete.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolVar(&removeVolumes, "volumes", false, "Remove persistent volumes")
	uninstallCmd.Flags().BoolVar(&removeNetwork, "network", false, "Remove the vpsik network")
}
```

### 3. إضافة `vpsik upgrade` — ملف: `cmd/upgrade.go`
```go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/config"
	"github.com/vpsik/workspace-installer/pkg/detector"
	"github.com/vpsik/workspace-installer/pkg/docker"
	"github.com/vpsik/workspace-installer/pkg/plan"
	"github.com/vpsik/workspace-installer/pkg/scanner"
	"github.com/vpsik/workspace-installer/pkg/state"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Pull latest images and recreate services",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configPath)
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}

		fmt.Println("🔍 Scanning environment...")
		scanResult := scanner.Run()

		if !scanResult.DockerAvailable {
			return fmt.Errorf("docker not available")
		}

		enabled := cfg.Services.EnabledList()
		detection := detector.Run(scanResult, enabled)
		svcState := state.Build(detection)
		installPlan := plan.Build(svcState, enabled)

		if !installPlan.HasChanges() {
			fmt.Println("All services already installed. Proceeding with upgrade...")
		}

		fmt.Println("\n📦 Pulling latest images...")
		if err := docker.PullImages(enabled); err != nil {
			return fmt.Errorf("pull images: %w", err)
		}

		fmt.Println("🔄 Recreating services...")
		if err := docker.RecreateServices(configPath, cfg.Workspace.Domain); err != nil {
			return fmt.Errorf("recreate services: %w", err)
		}

		fmt.Println("\n✅ Upgrade complete!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}
```

### 4. تحديث `pkg/docker/compose.go` — الإضافات المطلوبة

**أ. إضافة دوال جديدة: `Down`, `RemoveVolumes`, `RemoveNetwork`, `PullImages`, `RecreateServices`**

أضف بعد الـ function `EnsureNetwork`:

```go
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
	var images []string
	for _, name := range services {
		if tpl, ok := serviceTemplates[name]; ok {
			images = append(images, tpl.Image)
		}
	}

	for _, image := range images {
		fmt.Printf("  Pulling %s...\n", image)
		ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
		cmd := exec.CommandContext(ctx, "docker", "pull", image)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		cancel()
		if err != nil {
			return fmt.Errorf("pull %s: %w", image, err)
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
```

**ب. حل conflict port 3000 — تغيير OpenWebUI من `3000:8080` إلى `3001:8080`**

في `serviceTemplates` غيّر:
```go
"openwebui": {
	...
	Ports: []string{"3001:8080"},  // تغيير 3000 → 3001
	...
}
```

**ج. إضافة PostgreSQL containers**

أضف خدمة postgres في `serviceTemplates`:
```go
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
```

**د. إضافة Traefik reverse proxy مع SSL**

أضف خدمة traefik:
```go
"traefik": {
	Name:  "traefik",
	Image: "traefik:v3.0",
	Ports: []string{"80:80", "443:443"},
	Volumes: []string{
		"/var/run/docker.sock:/var/run/docker.sock",
		"traefik-data:/etc/traefik",
	},
	EnvVars: map[string]string{
		"CF_DNS_API_TOKEN": "${CF_DNS_API_TOKEN}",
	},
	Network: "vpsik",
},
```

أضف متغير `EnableTraefik bool` في `config.go` في `Services`.

### 5. Post-deploy Health Checks — تعديل `cmd/install.go`

بعد `docker.Deploy(composePath)`، أضف:
```go
fmt.Println("\n⏳ Waiting for services to be healthy...")
for _, name := range toInstall {
	fmt.Printf("  Checking %s...", name)
	healthy := waitForService(name, 60*time.Second)
	if healthy {
		fmt.Println(" ✅")
	} else {
		fmt.Println(" ⚠ timeout")
	}
}
```

ودالة `waitForService` تستخدم نفس منطق الكشف في `detector.go`.

### 6. إضافة `.env` generation — في `pkg/docker/compose.go`

أضف:
```go
func GenerateEnvFile(services []string, domain string, outputPath string) error {
	envVars := make(map[string]string)
	envVars["VPSIK_DOMAIN"] = domain

	for _, name := range services {
		if tpl, ok := serviceTemplates[name]; ok {
			for k, v := range tpl.EnvVars {
				if _, exists := envVars[k]; !exists {
					envVars[k] = v
				}
			}
		}
	}

	var sb strings.Builder
	for k, v := range envVars {
		sb.WriteString(fmt.Sprintf("%s=%s\n", k, v))
	}

	return os.WriteFile(outputPath, []byte(sb.String()), 0644)
}
```

### 7. State Serialization — في `pkg/state/state.go`

أضف:
```go
import (
	"encoding/json"
	"os"
)

func (s *State) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read state: %w", err)
	}
	var s State
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("unmarshal state: %w", err)
	}
	return &s, nil
}
```

### 8. إضافة Traefik إلى Config — في `pkg/config/config.go`

أضف حقل جديد:
```go
type Services struct {
	Authentik bool `yaml:"authentik"`
	Gitea     bool `yaml:"gitea"`
	Coolify   bool `yaml:"coolify"`
	Ollama    bool `yaml:"ollama"`
	OpenCode  bool `yaml:"opencode"`
	OpenWebUI bool `yaml:"openwebui"`
	Restic    bool `yaml:"restic"`
	Traefik   bool `yaml:"traefik"`     // جديد
	Postgres  bool `yaml:"postgres"`    // جديد
}
```

---

## ترتيب التنفيذ

1. إنشاء `cmd/init.go` — أسهل وأكثر طلب من المستخدم
2. إنشاء `cmd/uninstall.go` + إضافة دوال Docker
3. إنشاء `cmd/upgrade.go` + إضافة دوال Docker
4. تحديث `pkg/docker/compose.go`: Fix port 3000 → 3001
5. تحديث `pkg/docker/compose.go`: إضافة PostgreSQL + Traefik templates
6. تحديث `pkg/config/config.go`: إضافة Traefik, Postgres
7. إضافة State Serialization في `pkg/state/state.go`
8. إضافة .env generation
9. إضافة post-deploy health checks في install.go
10. بناء واختبار: `go build -o vpsik .`
