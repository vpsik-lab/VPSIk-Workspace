package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load("/tmp/nonexistent-config.yaml")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
	if cfg != nil {
		t.Fatal("expected nil config on error")
	}
}

func TestLoad_ValidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "api.yaml")
	content := []byte(`
server:
  port: 8081
  host: "127.0.0.1"
auth:
  jwt_secret: test-secret
  users:
    - username: admin
      password_hash: hash
services:
  gitea:
    url: http://gitea:3000
  coolify:
    url: http://coolify:3000
  ollama:
    url: http://ollama:11434
`)

	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 8081 {
		t.Errorf("expected port 8081, got %d", cfg.Server.Port)
	}
	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("expected host 127.0.0.1, got %s", cfg.Server.Host)
	}
	if cfg.Auth.JWTSecret != "test-secret" {
		t.Errorf("expected test-secret, got %s", cfg.Auth.JWTSecret)
	}
	if len(cfg.Auth.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(cfg.Auth.Users))
	}
	if cfg.Auth.Users[0].Username != "admin" {
		t.Errorf("expected admin, got %s", cfg.Auth.Users[0].Username)
	}
}

func TestLoad_DefaultPort(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "minimal.yaml")
	content := []byte("auth:\n  jwt_secret: x\n  users:\n    - username: a\n      password_hash: b\nservices:\n  gitea:\n    url: http://gitea:3000\n  coolify:\n    url: http://coolify:3000\n  ollama:\n    url: http://ollama:11434\n")

	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Server.Port != 8080 {
		t.Errorf("expected default port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("expected default host 0.0.0.0, got %s", cfg.Server.Host)
	}
}

func TestLoad_ResticConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "with-restic.yaml")
	content := []byte(`
auth:
  jwt_secret: x
  users:
    - username: admin
      password_hash: hash
services:
  gitea: {url: http://gitea:3000}
  coolify: {url: http://coolify:3000}
  ollama: {url: http://ollama:11434}
  restic:
    binary: /usr/local/bin/restic
    repo_url: /data/backup
    password: secret123
`)

	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.Services.Restic.Binary != "/usr/local/bin/restic" {
		t.Errorf("expected restic binary, got %s", cfg.Services.Restic.Binary)
	}
	if cfg.Services.Restic.RepoURL != "/data/backup" {
		t.Errorf("expected /data/backup, got %s", cfg.Services.Restic.RepoURL)
	}
	if cfg.Services.Restic.Password != "secret123" {
		t.Errorf("expected secret123, got %s", cfg.Services.Restic.Password)
	}
}

func TestFindUser(t *testing.T) {
	cfg := &APIConfig{
		Auth: AuthConfig{
			Users: []User{
				{Username: "admin", PasswordHash: "hash1"},
				{Username: "dev", PasswordHash: "hash2"},
			},
		},
	}

	u := cfg.FindUser("admin")
	if u == nil {
		t.Fatal("expected to find admin")
	}
	if u.PasswordHash != "hash1" {
		t.Errorf("expected hash1, got %s", u.PasswordHash)
	}

	u = cfg.FindUser("nonexistent")
	if u != nil {
		t.Fatal("expected nil for nonexistent user")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.yaml")
	os.WriteFile(path, []byte("{invalid: yaml: broken"), 0644)

	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestLoad_EmptyFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.yaml")
	os.WriteFile(path, []byte(""), 0644)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
}
