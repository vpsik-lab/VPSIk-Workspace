package scanner

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ContainerInfo struct {
	Name   string
	Image  string
	Status string
	Ports  []string
}

type SystemInfo struct {
	OS             string
	Architecture   string
	CPU            int
	RAMMB          int64
	DiskFreeMB     int64
	DockerVersion  string
	ComposeVersion string
}

type CoolifyInfo struct {
	Installed bool
	Running   bool
	Container string
	Port      int
	Version   string
	HasProxy  bool
}

type ScanResult struct {
	DockerAvailable  bool
	DockerRunning    bool
	ComposeAvailable bool
	System           SystemInfo
	Containers       []ContainerInfo
	Networks         []string
	UsedPorts        []int
	CoolifyDetected  *CoolifyInfo
	Errors           []string
}

func Run() *ScanResult {
	result := &ScanResult{}

	detectSystem(&result.System)
	collectErrors(result)

	dockerAvailable, err := checkDocker()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("docker check: %v", err))
	}
	result.DockerAvailable = dockerAvailable
	result.DockerRunning = dockerAvailable

	if !dockerAvailable {
		return result
	}

	ver, _ := getDockerVersion()
	result.System.DockerVersion = ver

	cv, _ := getComposeVersion()
	result.System.ComposeVersion = cv
	result.ComposeAvailable = cv != ""

	if containers, err := listContainers(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("list containers: %v", err))
	} else {
		result.Containers = containers
	}

	if networks, err := listNetworks(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("list networks: %v", err))
	} else {
		result.Networks = networks
	}

	if ports, err := scanPorts(); err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("port scan: %v", err))
	} else {
		result.UsedPorts = ports
	}

	result.CoolifyDetected = detectCoolify(result.Containers, result.UsedPorts)

	return result
}

func detectSystem(info *SystemInfo) {
	info.OS = runtime.GOOS
	info.Architecture = runtime.GOARCH
	info.CPU = runtime.NumCPU()

	if mem, err := getSystemMemoryMB(); err == nil {
		info.RAMMB = mem
	}

	if disk, err := getDiskFreeMB("/"); err == nil {
		info.DiskFreeMB = disk
	}
}

func detectCoolify(containers []ContainerInfo, ports []int) *CoolifyInfo {
	ci := &CoolifyInfo{}

	for _, c := range containers {
		name := strings.ToLower(c.Name)
		image := strings.ToLower(c.Image)
		if strings.Contains(name, "coolify") || strings.Contains(image, "coolify") {
			ci.Installed = true
			ci.Running = strings.Contains(c.Status, "Up") || strings.Contains(c.Status, "running")
			ci.Container = c.Name
			for _, p := range c.Ports {
				if strings.Contains(p, "8000") || strings.Contains(p, "3000") {
					fmt.Sscanf(p, "%d", &ci.Port)
				}
			}
			break
		}
	}

	if !ci.Installed {
		for _, p := range ports {
			if p == 8000 || p == 3000 {
				ci.Installed = true
				ci.Port = p
				break
			}
		}
	}

	for _, p := range ports {
		if p == 80 || p == 443 {
			ci.HasProxy = true
		}
	}

	if ci.Installed {
		return ci
	}
	return nil
}

func checkDocker() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "info", "--format", "{{.ServerVersion}}")
	out, err := cmd.Output()
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(string(out)) != "", nil
}

func getDockerVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "version", "--format", "{{.Server.Version}}")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getComposeVersion() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "compose", "version", "--short")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func getSystemMemoryMB() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", "free -m | awk '/^Mem:/ {print $2}'")
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	var mem int64
	fmt.Sscanf(string(out), "%d", &mem)
	return mem, nil
}

func getDiskFreeMB(path string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "df", "-m", path)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		return 0, fmt.Errorf("no output from df")
	}
	fields := strings.Fields(lines[1])
	if len(fields) < 4 {
		return 0, fmt.Errorf("unexpected df output")
	}
	var free int64
	fmt.Sscanf(fields[3], "%d", &free)
	return free, nil
}

func listContainers() ([]ContainerInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "ps", "-a",
		"--format", "{{.Names}}\t{{.Image}}\t{{.Status}}\t{{.Ports}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var containers []ContainerInfo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "\t", 4)
		info := ContainerInfo{}
		if len(parts) > 0 {
			info.Name = parts[0]
		}
		if len(parts) > 1 {
			info.Image = parts[1]
		}
		if len(parts) > 2 {
			info.Status = parts[2]
		}
		if len(parts) > 3 && parts[3] != "" {
			info.Ports = strings.Fields(parts[3])
		}
		containers = append(containers, info)
	}
	return containers, nil
}

func listNetworks() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "network", "ls", "--format", "{{.Name}}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	networks := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(networks) == 1 && networks[0] == "" {
		return []string{}, nil
	}
	return networks, nil
}

func scanPorts() ([]int, error) {
	commonPorts := []int{
		22, 80, 443, 3000, 3001, 3002, 4000, 5000, 5432, 6379,
		8000, 8001, 8065, 8080, 8081, 8443, 9000, 9090, 11434, 30081,
	}
	var used []int
	for _, port := range commonPorts {
		if isPortUsed(port) {
			used = append(used, port)
		}
	}
	return used, nil
}

func isPortUsed(port int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func collectErrors(result *ScanResult) {
	if result.System.RAMMB > 0 && result.System.RAMMB < 2048 {
		result.Errors = append(result.Errors, fmt.Sprintf("low memory: %d MB (minimum 2048 MB)", result.System.RAMMB))
	}
	if result.System.DiskFreeMB > 0 && result.System.DiskFreeMB < 10240 {
		result.Errors = append(result.Errors, fmt.Sprintf("low disk space: %d MB (minimum 10240 MB)", result.System.DiskFreeMB))
	}
}
