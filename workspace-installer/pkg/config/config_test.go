package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Workspace.Domain != "workspace.vpsik.com" {
		t.Errorf("expected default domain workspace.vpsik.com, got %s", cfg.Workspace.Domain)
	}
	if !cfg.Services.Gitea {
		t.Error("expected gitea to be enabled by default")
	}
	if !cfg.Services.Traefik {
		t.Error("expected traefik to be enabled by default")
	}
}

func TestEnabledList(t *testing.T) {
	tests := []struct {
		name     string
		services Services
		expected int
	}{
		{"all enabled", Services{Authentik: true, Gitea: true, Coolify: true, Ollama: true, OpenCode: true, OpenWebUI: true, Restic: true, Traefik: true, Postgres: true, Grafana: true, Prometheus: true, CodeServer: true, Redis: true}, 13},
		{"none enabled", Services{}, 0},
		{"only gitea", Services{Gitea: true}, 1},
		{"traefik+postgres", Services{Traefik: true, Postgres: true}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enabled := tt.services.EnabledList()
			if len(enabled) != tt.expected {
				t.Errorf("expected %d enabled, got %d: %v", tt.expected, len(enabled), enabled)
			}
		})
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := Load("/tmp/nonexistent-file-12345.yaml")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestLoad_ValidFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "workspace.yaml")
	content := []byte("workspace:\n  domain: test.example.com\nservices:\n  gitea: true\n  ollama: false\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Workspace.Domain != "test.example.com" {
		t.Errorf("expected test.example.com, got %s", cfg.Workspace.Domain)
	}

	if !cfg.Services.Gitea {
		t.Error("expected gitea to be true")
	}
	if cfg.Services.Ollama {
		t.Error("expected ollama to be false")
	}
}

func TestLoad_DefaultsMerged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "minimal.yaml")
	content := []byte("workspace:\n  domain: test.com\nservices:\n  gitea: true\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !cfg.Services.Authentik {
		t.Error("expected authentik to remain true from defaults")
	}
}

func TestEnabledList_Order(t *testing.T) {
	s := Services{Gitea: true, Ollama: true, Authentik: true}
	list := s.EnabledList()

	expected := []string{"authentik", "gitea", "ollama"}
	for i, name := range expected {
		if i >= len(list) {
			t.Fatalf("expected %s at index %d, got out of bounds", name, i)
		}
		if list[i] != name {
			t.Errorf("expected %s at index %d, got %s", name, i, list[i])
		}
	}
}
