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
	"github.com/vpsik/workspace-installer/pkg/cliui"
	"github.com/vpsik/workspace-installer/pkg/scanner"
)

var autoFix bool

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system health and fix issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(cliui.Header("WorkSpace OS — System Health Check"))

		issues := 0
		fixed := 0
		scan := scanner.Run()

		tbl := cliui.NewTable([]cliui.Column{
			{Header: "Check", Width: 24},
			{Header: "Status"},
			{Header: "Details"},
		})

		// 1. OS
		tbl.AddRow("OS", cliui.Success("OK"), fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))

		// 2. Docker
		if scan.DockerAvailable {
			tbl.AddRow("Docker", cliui.Success("OK"), scan.System.DockerVersion)
		} else {
			tbl.AddRow("Docker", cliui.Error("FAIL"), "not installed")
			issues++
			if autoFix && fixDocker() {
				tbl.AddRow("  → Fix", cliui.Success("FIXED"), "Docker installed")
				fixed++
			}
		}

		// 3. Docker Compose
		if scan.ComposeAvailable {
			tbl.AddRow("Docker Compose", cliui.Success("OK"), scan.System.ComposeVersion)
		} else {
			tbl.AddRow("Docker Compose", cliui.Error("FAIL"), "not installed")
			issues++
			if autoFix && fixDockerCompose() {
				tbl.AddRow("  → Fix", cliui.Success("FIXED"), "Docker Compose installed")
				fixed++
			}
		}

		// 4. RAM
		if scan.System.RAMMB > 0 {
			if scan.System.RAMMB < 2048 {
				tbl.AddRow("RAM", cliui.Warning("WARN"), fmt.Sprintf("%d MB (min 2048)", scan.System.RAMMB))
				issues++
			} else {
				tbl.AddRow("RAM", cliui.Success("OK"), fmt.Sprintf("%d MB", scan.System.RAMMB))
			}
		}

		// 5. Disk
		if scan.System.DiskFreeMB > 0 {
			if scan.System.DiskFreeMB < 10240 {
				tbl.AddRow("Disk", cliui.Warning("WARN"), fmt.Sprintf("%d MB free (min 10240)", scan.System.DiskFreeMB))
				issues++
			} else {
				tbl.AddRow("Disk", cliui.Success("OK"), fmt.Sprintf("%d MB free", scan.System.DiskFreeMB))
			}
		}

		// 6. Port conflicts
		conflicts := 0
		for port, desc := range map[int]string{80: "HTTP", 443: "HTTPS", 3000: "Dashboard", 8000: "Coolify"} {
			for _, up := range scan.UsedPorts {
				if up == port {
					tbl.AddRow("Port "+fmt.Sprint(port), cliui.Warning("WARN"), desc+" in use")
					conflicts++
					break
				}
			}
		}
		if conflicts == 0 {
			tbl.AddRow("Port Conflicts", cliui.Success("OK"), "none detected")
		}

		// 7. Coolify
		if scan.CoolifyDetected != nil {
			ci := scan.CoolifyDetected
			status := "stopped"
			if ci.Running {
				status = "running"
			}
			tbl.AddRow("Coolify", cliui.Success("OK"), fmt.Sprintf("%s (port %d, %s)", ci.Container, ci.Port, status))
		} else {
			tbl.AddRow("Coolify", cliui.DimText("INFO"), "not detected")
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
			tbl.AddRow("Network", cliui.Success("OK"), "workspace_net exists")
		} else {
			tbl.AddRow("Network", cliui.DimText("INFO"), "will be created on install")
		}

		// 9. NTP
		if checkNTP() {
			tbl.AddRow("NTP", cliui.Success("OK"), "time synchronized")
		} else {
			tbl.AddRow("NTP", cliui.Warning("WARN"), "time not synchronized")
			issues++
			if autoFix && fixNTP() {
				tbl.AddRow("  → Fix", cliui.Success("FIXED"), "NTP enabled")
				fixed++
			}
		}

		tbl.Print()

		if issues == 0 {
			fmt.Println(cliui.Success("\n  All checks passed — system is ready!"))
		} else {
			msg := fmt.Sprintf("Found %d issues", issues)
			if fixed > 0 {
				msg += fmt.Sprintf(" (%d fixed)", fixed)
			}
			if !autoFix {
				msg += " — run with --fix to auto-repair"
			}
			fmt.Println(cliui.Warning("\n  " + msg))
		}

		return nil
	},
}

func fixDocker() bool {
	fmt.Println(cliui.Warning("     Installing Docker from get.docker.com..."))
	fmt.Println(cliui.DimText("     This runs a script from the internet. Verify at https://get.docker.com"))
	fmt.Print(cliui.Highlight("     Continue? (yes/no): "))
	var resp string
	fmt.Scanln(&resp)
	if resp != "yes" {
		fmt.Println("     Skipped.")
		return false
	}
	cmd := exec.Command("sh", "-c",
		"curl -fsSL https://get.docker.com | sh")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run() == nil
}

func fixDockerCompose() bool {
	fmt.Println(cliui.Warning("     Installing Docker Compose plugin..."))
	fmt.Println(cliui.DimText("     Downloads binary from GitHub releases (no checksum verification)"))
	fmt.Print(cliui.Highlight("     Continue? (yes/no): "))
	var resp string
	fmt.Scanln(&resp)
	if resp != "yes" {
		fmt.Println("     Skipped.")
		return false
	}
	cmd := exec.Command("sh", "-c",
		"DOCKER_CONFIG=${DOCKER_CONFIG:-$HOME/.docker}; "+
			"mkdir -p $DOCKER_CONFIG/cli-plugins; "+
			"curl -fsSL https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m) "+
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
