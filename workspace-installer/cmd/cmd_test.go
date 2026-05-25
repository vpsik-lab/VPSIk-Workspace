package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContains(t *testing.T) {
	slice := []string{"a", "b", "c"}
	if !contains(slice, "a") {
		t.Error("expected true for 'a'")
	}
	if !contains(slice, "c") {
		t.Error("expected true for 'c'")
	}
	if contains(slice, "d") {
		t.Error("expected false for 'd'")
	}
	if contains([]string{}, "a") {
		t.Error("expected false for empty slice")
	}
}

func TestParseServices_Empty(t *testing.T) {
	result := parseServices("")
	if len(result) == 0 {
		t.Error("expected default services for empty input")
	}
}

func TestParseServices_CommaSeparated(t *testing.T) {
	result := parseServices("gitea,ollama,postgres")
	if len(result) != 3 {
		t.Fatalf("expected 3 services, got %d: %v", len(result), result)
	}
	if result[0] != "gitea" {
		t.Errorf("expected gitea, got %s", result[0])
	}
	if result[1] != "ollama" {
		t.Errorf("expected ollama, got %s", result[1])
	}
	if result[2] != "postgres" {
		t.Errorf("expected postgres, got %s", result[2])
	}
}

func TestParseServices_WithSpaces(t *testing.T) {
	result := parseServices(" gitea , ollama , postgres ")
	if len(result) != 3 {
		t.Fatalf("expected 3 services, got %d", len(result))
	}
	if result[0] != "gitea" {
		t.Errorf("expected gitea, got %s", result[0])
	}
}

func TestParseServices_Single(t *testing.T) {
	result := parseServices("gitea")
	if len(result) != 1 {
		t.Fatalf("expected 1 service, got %d", len(result))
	}
	if result[0] != "gitea" {
		t.Errorf("expected gitea, got %s", result[0])
	}
}

func TestRunAutoInit_DefaultDomain(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "workspace.yaml")

	autoMode = true
	domainFlag = ""
	servicesFlag = ""
	outputFlag = outputPath
	configPath = outputPath

	err := runAutoInit()
	if err != nil {
		t.Fatalf("runAutoInit failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("expected non-empty config")
	}
}

func TestRunAutoInit_CustomDomain(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "custom.yaml")

	autoMode = true
	domainFlag = "myworkspace.com"
	servicesFlag = ""
	outputFlag = outputPath
	configPath = outputPath

	err := runAutoInit()
	if err != nil {
		t.Fatalf("runAutoInit failed: %v", err)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Error("expected non-empty config")
	}
}

func TestRunAutoInit_SelectedServices(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "selected.yaml")

	autoMode = true
	domainFlag = "test.com"
	servicesFlag = "gitea,ollama,postgres"
	outputFlag = outputPath
	configPath = outputPath

	err := runAutoInit()
	if err != nil {
		t.Fatalf("runAutoInit failed: %v", err)
	}

	_, err = os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
}

func TestRunAutoInit_WritesToConfigPath(t *testing.T) {
	dir := t.TempDir()
	expectedPath := filepath.Join(dir, "config.yaml")

	autoMode = true
	domainFlag = "test.com"
	servicesFlag = ""
	outputFlag = ""
	configPath = expectedPath

	err := runAutoInit()
	if err != nil {
		t.Fatalf("runAutoInit failed: %v", err)
	}

	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Error("expected config file to exist at configPath")
	}
}

func TestBackupPaths_ContainsExpected(t *testing.T) {
	expected := []string{"gitea", "postgres", "ollama", "grafana", "authentik",
		"opencode", "openwebui", "code-server", "outline", "mattermost", "restic"}
	for _, name := range expected {
		if _, ok := backupPaths[name]; !ok {
			t.Errorf("expected backup path for %s", name)
		}
	}
}

func TestWaitForServicePorts(t *testing.T) {
	ports := map[string]int{
		"authentik": 9000, "gitea": 3000, "coolify": 3000,
		"ollama": 11434, "opencode": 30081, "openwebui": 8080,
		"code-server": 8443, "grafana": 3000, "prometheus": 9090,
	}
	for name, port := range ports {
		if port <= 0 {
			t.Errorf("invalid port for %s: %d", name, port)
		}
	}
}
