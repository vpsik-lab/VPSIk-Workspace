package scanner

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

type ContainerInfo struct {
	Name   string
	Image  string
	Status string
	Ports  []string
}

type ScanResult struct {
	DockerAvailable bool
	Containers      []ContainerInfo
	Networks        []string
	UsedPorts       []int
	Errors          []string
}

func Run() *ScanResult {
	result := &ScanResult{}

	dockerAvailable, err := checkDocker()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("docker check: %v", err))
	}
	result.DockerAvailable = dockerAvailable

	if !dockerAvailable {
		return result
	}

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

	return result
}

func checkDocker() (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", "info", "--format", "{{.ServerVersion}}")
	out, err := cmd.Output()
	if err != nil {
		return false, nil
	}

	return strings.TrimSpace(string(out)) != "", nil
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
		80, 443, 3000, 4000, 5000, 5432, 6379,
		8000, 8080, 8443, 9000, 9090, 11434, 30081,
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
