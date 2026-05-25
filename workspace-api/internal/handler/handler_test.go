package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vpsik/workspace-api/internal/config"
)

func TestStatusHandler_WithRestic(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Restic: config.ResticConfig{
				Binary:   "restic",
				RepoURL:  "/tmp/test-repo",
				Password: "test",
			},
		},
	}

	clients := NewClients(cfg)
	h := NewStatusHandler(clients)

	req := httptest.NewRequest("GET", "/api/status", nil)
	w := httptest.NewRecorder()
	h.Status(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var sr StatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if len(sr.Services) == 0 {
		t.Fatal("expected at least one service")
	}

	if sr.Timestamp == "" {
		t.Error("expected timestamp")
	}

	foundRestic := false
	for _, s := range sr.Services {
		if s.Name == "restic" {
			foundRestic = true
			break
		}
	}
	if !foundRestic {
		t.Error("expected restic in services")
	}
}

func TestStatusHandler_TimestampFormat(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Gitea:   config.ServiceEndpoint{URL: "http://gitea:3000"},
			Coolify: config.ServiceEndpoint{URL: "http://coolify:3000"},
			Ollama:  config.ServiceEndpoint{URL: "http://ollama:11434"},
		},
	}

	clients := NewClients(cfg)
	h := NewStatusHandler(clients)

	req := httptest.NewRequest("GET", "/api/status", nil)
	w := httptest.NewRecorder()
	h.Status(w, req)

	var sr StatusResponse
	resp := w.Result()
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&sr)

	_, err := time.Parse(time.RFC3339, sr.Timestamp)
	if err != nil {
		t.Errorf("expected RFC3339 timestamp, got %q: %v", sr.Timestamp, err)
	}
}

func TestNewClients(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Gitea:   config.ServiceEndpoint{URL: "http://gitea:3000", Token: "abc"},
			Coolify: config.ServiceEndpoint{URL: "http://coolify:3000", Token: "def"},
			Ollama:  config.ServiceEndpoint{URL: "http://ollama:11434"},
			OpenCode: config.ServiceEndpoint{URL: "http://opencode:8080", Token: "ghi"},
			Restic:  config.ResticConfig{Binary: "restic", RepoURL: "/repo", Password: "pass"},
		},
	}

	clients := NewClients(cfg)
	if clients == nil {
		t.Fatal("expected non-nil clients")
	}

	if clients.Gitea == nil {
		t.Error("expected gitea client")
	}
	if clients.Coolify == nil {
		t.Error("expected coolify client")
	}
	if clients.Ollama == nil {
		t.Error("expected ollama client")
	}
	if clients.OpenCode == nil {
		t.Error("expected opencode client")
	}
	if clients.Restic == nil {
		t.Error("expected restic client")
	}
}

func TestNewClients_NilOpenCode(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Gitea:   config.ServiceEndpoint{URL: "http://gitea:3000"},
			Coolify: config.ServiceEndpoint{URL: "http://coolify:3000"},
			Ollama:  config.ServiceEndpoint{URL: "http://ollama:11434"},
			CodeServer: config.ServiceEndpoint{},
			Plane:      config.ServiceEndpoint{},
			Outline:    config.ServiceEndpoint{},
			Mattermost: config.ServiceEndpoint{},
		},
	}

	clients := NewClients(cfg)
	if clients.OpenCode != nil {
		t.Error("expected nil opencode client when URL is empty")
	}
	if clients.CodeServer == nil {
		t.Error("expected codeserver client even with empty URL")
	}
}
