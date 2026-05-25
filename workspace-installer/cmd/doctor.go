package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/vpsik/workspace-installer/pkg/scanner"
)

var autoFix bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system health and fix issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🏥 VPSIk Workspace — System Health Check")
		fmt.Println(strings.Repeat("─", 40))

		issues := 0
		fixed := 0

		scan := scanner.Run()

		// 1. OS
		fmt.Printf("  OS: %s/%s\n", runtime.GOOS, runtime.GOARCH)

		// 2. Docker
		if scan.DockerAvailable {
			fmt.Printf("  ✅ Docker: %s\n", scan.System.DockerVersion)
		} else {
			fmt.Println("  ❌ Docker: not installed")
			issues++
			if autoFix {
				if fixDocker() {
					fmt.Println("     ✅ Fixed: Docker installed")
					fixed++
				} else {
					fmt.Println("     ❌ Could not install Docker automatically")
				}
			}
		}

		// 3. Docker Compose
		if scan.ComposeAvailable {
			fmt.Printf("  ✅ Docker Compose: %s\n", scan.System.ComposeVersion)
		} else {
			fmt.Println("  ❌ Docker Compose: not installed")
			issues++
			if autoFix {
				if fixDockerCompose() {
					fmt.Println("     ✅ Fixed: Docker Compose installed")
					fixed++
				}
			}
		}

		// 4. RAM
		if scan.System.RAMMB > 0 {
			if scan.System.RAMMB < 2048 {
				fmt.Printf("  ⚠ RAM: %d MB (minimum 2048 MB)\n", scan.System.RAMMB)
				issues++
			} else {
				fmt.Printf("  ✅ RAM: %d MB\n", scan.System.RAMMB)
			}
		}

		// 5. Disk
		if scan.System.DiskFreeMB > 0 {
			if scan.System.DiskFreeMB < 10240 {
				fmt.Printf("  ⚠ Disk: %d MB free (minimum 10240 MB)\n", scan.System.DiskFreeMB)
				issues++
			} else {
				fmt.Printf("  ✅ Disk: %d MB free\n", scan.System.DiskFreeMB)
			}
		}

		// 6. Port conflicts
		criticalPorts := map[int]string{
			80:   "HTTP (Coolify/Traefik)",
			443:  "HTTPS (Coolify/Traefik)",
			3000: "Dashboard/Gitea",
			8000: "Coolify",
		}
		conflicts := 0
		for port, desc := range criticalPorts {
			for _, up := range scan.UsedPorts {
				if up == port {
					fmt.Printf("  ⚠ Port %d in use: %s\n", port, desc)
					conflicts++
					break
				}
			}
		}
		if conflicts == 0 {
			fmt.Println("  ✅ No port conflicts detected")
		}

		// 7. Coolify
		if scan.CoolifyDetected != nil {
			ci := scan.CoolifyDetected
			status := "stopped"
			if ci.Running {
				status = "running"
			}
			fmt.Printf("  ✅ Coolify: %s (port %d, %s)\n", ci.Container, ci.Port, status)
			if ci.HasProxy {
				fmt.Println("     ⚠ Using ports 80/443 — use Cloudflare Tunnel for VPSIk")
			}
		} else {
			fmt.Println("  ℹ Coolify: not detected")
		}

		// 8. workspace_net
		netFound := false
		for _, n := range scan.Networks {
			if n == "workspace_net" {
				netFound = true
				break
			}
		}
		if netFound {
			fmt.Println("  ✅ Network: workspace_net exists")
		} else {
			fmt.Println("  ℹ Network: workspace_net will be created on install")
		}

		// 9. NTP
		if checkNTP() {
			fmt.Println("  ✅ NTP: time synchronized")
		} else {
			fmt.Println("  ⚠ NTP: time not synchronized")
			issues++
			if autoFix {
				if fixNTP() {
					fmt.Println("     ✅ Fixed: NTP enabled")
					fixed++
				}
			}
		}

		fmt.Println(strings.Repeat("─", 40))
		if issues == 0 {
			fmt.Println("✅ All checks passed — system is ready!")
		} else {
			fmt.Printf("⚠ Found %d issues", issues)
			if fixed > 0 {
				fmt.Printf(" (%d fixed)", fixed)
			}
			if !autoFix {
				fmt.Print(" — run with --fix to auto-repair")
			}
			fmt.Println()
		}

		return nil
	},
}

func fixDocker() bool {
	fmt.Println("     Installing Docker...")
	cmd := exec.Command("sh", "-c",
		"curl -fsSL https://get.docker.com | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}

func fixDockerCompose() bool {
	fmt.Println("     Installing Docker Compose plugin...")
	cmd := exec.Command("sh", "-c",
		"DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}; "+
			"mkdir -p $DOCKER_CONFIG/cli-plugins; "+
			"curl -SL https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m) "+
			"-o $DOCKER_CONFIG/cli-plugins/docker-compose; "+
			"chmod +x $DOCKER_CONFIG/cli-plugins/docker-compose")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}

func checkNTP() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "timedatectl", "show", "-p", "NTPSynchronized", "--value")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "yes"
}

func fixNTP() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "timedatectl", "set-ntp", "true")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}

func init() {
	rootCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVar(&autoFix, "fix", false, "Auto-fix detected issues")
}
