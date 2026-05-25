package detector

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/vpsik/workspace-installer/pkg/scanner"
)

type Status int

const (
	StatusMissing   Status = iota
	StatusInstalled Status = iota
)

func (s Status) String() string {
	switch s {
	case StatusInstalled:
		return "installed"
	case StatusMissing:
		return "missing"
	default:
		return "unknown"
	}
}

type ServiceInfo struct {
	Name        string
	Status      Status
	Details     string
	Container   string
	Port        int
	APIEndpoint string
}

type Result struct {
	Services []ServiceInfo
}

func Run(scan *scanner.ScanResult, enabledServices []string) *Result {
	detectors := map[string]func() *ServiceInfo{
		"authentik": func() *ServiceInfo { return detectContainer(scan, "authentik", "authentik", 9000, "http://localhost:9000") },
		"gitea":     func() *ServiceInfo { return detectContainer(scan, "gitea", "gitea/gitea", 3000, "http://localhost:3000/api/healthz") },
		"coolify":   func() *ServiceInfo { return detectContainer(scan, "coolify", "coollabsio/coolify", 8000, "http://localhost:8000/api/v1/health") },
		"ollama":    func() *ServiceInfo { return detectPort("ollama", "ollama/ollama", 11434, "http://localhost:11434/api/tags") },
		"opencode":  func() *ServiceInfo { return detectContainer(scan, "opencode", "opencode", 30081, "") },
		"openwebui": func() *ServiceInfo { return detectContainer(scan, "openwebui", "openwebui", 3000, "http://localhost:3000") },
		"restic":    func() *ServiceInfo { return detectBinary("restic") },
	}

	result := &Result{}
	for _, name := range enabledServices {
		if detect, ok := detectors[name]; ok {
			info := detect()
			result.Services = append(result.Services, *info)
		}
	}

	return result
}

func detectContainer(scan *scanner.ScanResult, name, imageHint string, port int, apiURL string) *ServiceInfo {
	info := &ServiceInfo{Name: name, Status: StatusMissing}

	for _, c := range scan.Containers {
		if strings.Contains(strings.ToLower(c.Name), name) ||
			strings.Contains(strings.ToLower(c.Image), imageHint) {
			info.Status = StatusInstalled
			info.Container = c.Name
			info.Details = fmt.Sprintf("Container %s (%s) — %s", c.Name, c.Image, c.Status)
			return info
		}
	}

	for _, p := range scan.UsedPorts {
		if p == port {
			info.Status = StatusInstalled
			info.Port = port
			info.Details = fmt.Sprintf("Port %d is in use", port)
			return info
		}
	}

	if apiURL != "" {
		if CheckAPI(apiURL) {
			info.Status = StatusInstalled
			info.APIEndpoint = apiURL
			info.Details = fmt.Sprintf("API reachable at %s", apiURL)
			return info
		}
	}

	info.Details = "Not detected"
	return info
}

func detectPort(name, imageHint string, port int, apiURL string) *ServiceInfo {
	info := &ServiceInfo{Name: name, Status: StatusMissing}

	if CheckPort(port) {
		info.Status = StatusInstalled
		info.Port = port
		info.Details = fmt.Sprintf("Port %d is in use", port)
		return info
	}

	if apiURL != "" && CheckAPI(apiURL) {
		info.Status = StatusInstalled
		info.APIEndpoint = apiURL
		info.Details = fmt.Sprintf("API reachable at %s", apiURL)
		return info
	}

	info.Details = "Not detected"
	return info
}

func detectBinary(name string) *ServiceInfo {
	info := &ServiceInfo{Name: name, Status: StatusMissing}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := exec.CommandContext(ctx, "which", name).Run(); err == nil {
		info.Status = StatusInstalled
		info.Details = fmt.Sprintf("Binary found: %s", name)
		return info
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel2()

	if out, err := exec.CommandContext(ctx2, name, "version").Output(); err == nil {
		info.Status = StatusInstalled
		info.Details = fmt.Sprintf("Binary found: %s (%s)", name, strings.TrimSpace(string(out)))
		return info
	}

	info.Details = "Not detected"
	return info
}

func CheckAPI(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode < 500
}

func CheckPort(port int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "ps",
		"--filter", fmt.Sprintf("publish=%d", port),
		"--format", "{{.Names}}")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) != ""
}
