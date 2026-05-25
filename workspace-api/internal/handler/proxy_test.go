package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vpsik/workspace-api/internal/config"
)

func TestJSONHelpers(t *testing.T) {
	w := httptest.NewRecorder()
	writeJSON(w, map[string]string{"status": "ok"})

	resp := w.Result()
	defer resp.Body.Close()

	if resp.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json, got %s", resp.Header.Get("Content-Type"))
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result["status"] != "ok" {
		t.Errorf("expected ok, got %s", result["status"])
	}
}

func TestJSONError(t *testing.T) {
	w := httptest.NewRecorder()
	jsonError(w, "something went wrong", http.StatusBadRequest)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if result["error"] != "something went wrong" {
		t.Errorf("expected error message, got %s", result["error"])
	}
}

func TestOllamaChat_InvalidBody(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Ollama: config.ServiceEndpoint{URL: "http://ollama:11434"},
		},
	}
	clients := NewClients(cfg)
	h := NewProxyHandler(clients)

	req := httptest.NewRequest("POST", "/api/ollama/chat", strings.NewReader("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.OllamaChat(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestOllamaChat_EmptyBody(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Ollama: config.ServiceEndpoint{URL: "http://ollama:11434"},
		},
	}
	clients := NewClients(cfg)
	h := NewProxyHandler(clients)

	req := httptest.NewRequest("POST", "/api/ollama/chat", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.OllamaChat(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("expected 502, got %d", resp.StatusCode)
	}
}

func TestOllamaPullModel_EmptyModel(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Ollama: config.ServiceEndpoint{URL: "http://ollama:11434"},
		},
	}
	clients := NewClients(cfg)
	h := NewProxyHandler(clients)

	req := httptest.NewRequest("POST", "/api/ollama/pull", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.OllamaPullModel(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for empty model, got %d", resp.StatusCode)
	}
}

func TestOllamaTask_UnknownTask(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Ollama: config.ServiceEndpoint{URL: "http://ollama:11434"},
		},
	}
	clients := NewClients(cfg)
	h := NewProxyHandler(clients)

	body := `{"model":"llama3.2","task":"invalid-task","content":"test"}`
	req := httptest.NewRequest("POST", "/api/ollama/task", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.OllamaTask(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400 for unknown task, got %d", resp.StatusCode)
	}
}

func TestOpenCodeChat_NotConfigured(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Ollama: config.ServiceEndpoint{URL: "http://ollama:11434"},
		},
	}
	clients := NewClients(cfg)
	h := NewProxyHandler(clients)

	body := `{"message":"hello"}`
	req := httptest.NewRequest("POST", "/api/opencode/chat", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.OpenCodeChat(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("expected 502, got %d", resp.StatusCode)
	}
}

func TestResticBackup_NoPaths(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Restic: config.ResticConfig{Binary: "restic", RepoURL: "/repo", Password: "pw"},
		},
	}
	clients := NewClients(cfg)
	h := NewProxyHandler(clients)

	req := httptest.NewRequest("POST", "/api/restic/backup", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.ResticBackup(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestResticRestore_MissingFields(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Restic: config.ResticConfig{Binary: "restic", RepoURL: "/repo", Password: "pw"},
		},
	}
	clients := NewClients(cfg)
	h := NewProxyHandler(clients)

	tests := []struct {
		name string
		body string
	}{
		{"empty object", `{}`},
		{"missing target", `{"snapshot_id":"abc"}`},
		{"missing snapshot", `{"target":"/tmp"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/restic/restore", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h.ResticRestore(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", resp.StatusCode)
			}
		})
	}
}

func TestGiteaListRepos_NoClient(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Gitea: config.ServiceEndpoint{},
		},
	}
	clients := NewClients(cfg)
	h := NewProxyHandler(clients)

	req := httptest.NewRequest("GET", "/api/gitea/repos", nil)
	w := httptest.NewRecorder()
	h.GiteaRepos(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadGateway {
		t.Errorf("expected 502, got %d", resp.StatusCode)
	}
}
