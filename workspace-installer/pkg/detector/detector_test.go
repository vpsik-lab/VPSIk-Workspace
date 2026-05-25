package detector

import (
	"testing"

	"github.com/vpsik/workspace-installer/pkg/scanner"
)

func TestStatusString(t *testing.T) {
	if StatusInstalled.String() != "installed" {
		t.Errorf("expected installed, got %s", StatusInstalled.String())
	}
	if StatusMissing.String() != "missing" {
		t.Errorf("expected missing, got %s", StatusMissing.String())
	}
}

func TestDetectContainer_ByName(t *testing.T) {
	scan := &scanner.ScanResult{
		Containers: []scanner.ContainerInfo{
			{Name: "gitea-server", Image: "gitea/gitea:1.21", Status: "running"},
		},
		UsedPorts: []int{3000},
	}

	info := detectContainer(scan, "gitea", "gitea/gitea", 3000, "http://localhost:3000")
	if info.Status != StatusInstalled {
		t.Errorf("expected installed, got %s", info.Status)
	}
	if info.Container != "gitea-server" {
		t.Errorf("expected container gitea-server, got %s", info.Container)
	}
}

func TestDetectContainer_ByImage(t *testing.T) {
	scan := &scanner.ScanResult{
		Containers: []scanner.ContainerInfo{
			{Name: "my-gitea", Image: "gitea/gitea:latest", Status: "running"},
		},
	}

	info := detectContainer(scan, "gitea", "gitea/gitea", 3000, "")
	if info.Status != StatusInstalled {
		t.Errorf("expected installed by image match, got %s", info.Status)
	}
}

func TestDetectContainer_ByPort(t *testing.T) {
	scan := &scanner.ScanResult{
		Containers: []scanner.ContainerInfo{
			{Name: "other-service", Image: "nginx", Status: "running"},
		},
		UsedPorts: []int{9000},
	}

	info := detectContainer(scan, "authentik", "authentik", 9000, "")
	if info.Status != StatusInstalled {
		t.Errorf("expected installed by port, got %s", info.Status)
	}
	if info.Port != 9000 {
		t.Errorf("expected port 9000, got %d", info.Port)
	}
}

func TestDetectContainer_NotFound(t *testing.T) {
	scan := &scanner.ScanResult{
		Containers: []scanner.ContainerInfo{
			{Name: "nginx", Image: "nginx", Status: "running"},
		},
	}

	info := detectContainer(scan, "gitea", "gitea/gitea", 3000, "")
	if info.Status != StatusMissing {
		t.Errorf("expected missing, got %s", info.Status)
	}
}

func TestDetectBinary_NotFound(t *testing.T) {
	info := detectBinary("nonexistent-binary-xyz-123")
	if info.Status != StatusMissing {
		t.Errorf("expected missing, got %s", info.Status)
	}
}

func TestRun_EnabledServices(t *testing.T) {
	scan := &scanner.ScanResult{
		Containers: []scanner.ContainerInfo{
			{Name: "gitea", Image: "gitea/gitea:latest", Status: "running"},
		},
		UsedPorts: []int{3000},
	}

	result := Run(scan, []string{"gitea"})
	if len(result.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(result.Services))
	}
	if result.Services[0].Name != "gitea" {
		t.Errorf("expected gitea, got %s", result.Services[0].Name)
	}
}

func TestRun_SkipDisabled(t *testing.T) {
	scan := &scanner.ScanResult{}
	result := Run(scan, []string{})
	if len(result.Services) != 0 {
		t.Errorf("expected 0 services, got %d", len(result.Services))
	}
}

func TestRun_UnknownService(t *testing.T) {
	scan := &scanner.ScanResult{}
	result := Run(scan, []string{"unknown-service"})
	if len(result.Services) != 0 {
		t.Errorf("expected 0 for unknown service, got %d", len(result.Services))
	}
}
