package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/vpsik/workspace-api/internal/auth"
	"github.com/vpsik/workspace-api/internal/config"
	"github.com/vpsik/workspace-api/internal/middleware"
	"golang.org/x/crypto/bcrypt"
)

func mustHash(password string) string {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(h)
}

func TestAuthHandler_Login_Success(t *testing.T) {
	cfg := &config.APIConfig{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret-key-32-chars-min!!",
			Users: []config.User{
				{Username: "admin", PasswordHash: mustHash("password123")},
			},
		},
	}

	h := NewAuthHandler(cfg)

	body := `{"username":"admin","password":"password123"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var lr loginResponse
	if err := json.NewDecoder(resp.Body).Decode(&lr); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if lr.Token == "" {
		t.Error("expected token")
	}
	if lr.Username != "admin" {
		t.Errorf("expected admin, got %s", lr.Username)
	}

	cookies := resp.Cookies()
	foundCookie := false
	for _, c := range cookies {
		if c.Name == "workspace_token" {
			foundCookie = true
			if !c.HttpOnly {
				t.Error("expected HttpOnly cookie")
			}
			if c.Value == "" {
				t.Error("expected non-empty cookie value")
			}
			break
		}
	}
	if !foundCookie {
		t.Error("expected vpsik_token cookie")
	}
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	cfg := &config.APIConfig{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret-key-32-chars-min!!",
			Users: []config.User{
				{Username: "admin", PasswordHash: mustHash("password123")},
			},
		},
	}

	h := NewAuthHandler(cfg)

	tests := []struct {
		name     string
		body     string
		expected int
	}{
		{"wrong password", `{"username":"admin","password":"wrong"}`, http.StatusUnauthorized},
		{"wrong user", `{"username":"nobody","password":"password123"}`, http.StatusUnauthorized},
		{"invalid json", `not-json`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.Login(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, resp.StatusCode)
			}
		})
	}
}

func TestAuthHandler_Login_WrongMethod(t *testing.T) {
	cfg := &config.APIConfig{
		Auth: config.AuthConfig{
			JWTSecret: "test-secret",
			Users:     []config.User{{Username: "admin", PasswordHash: "hash"}},
		},
	}

	h := NewAuthHandler(cfg)
	req := httptest.NewRequest("GET", "/api/auth/login", nil)
	w := httptest.NewRecorder()

	h.Login(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	cfg := &config.APIConfig{Auth: config.AuthConfig{JWTSecret: "test"}}
	h := NewAuthHandler(cfg)

	req := httptest.NewRequest("POST", "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	cookies := resp.Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "workspace_token" && c.MaxAge < 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected expired vpsik_token cookie")
	}
}

func TestAuthHandler_Verify(t *testing.T) {
	cfg := &config.APIConfig{Auth: config.AuthConfig{JWTSecret: "test-secret"}}
	h := NewAuthHandler(cfg)

	req := httptest.NewRequest("GET", "/api/auth/verify", nil)
	ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.Verify(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if result["valid"] != true {
		t.Error("expected valid=true")
	}
	if result["username"] != "testuser" {
		t.Errorf("expected testuser, got %v", result["username"])
	}
}

func TestNewClients_AllServices(t *testing.T) {
	cfg := &config.APIConfig{
		Services: config.ServicesConfig{
			Gitea:      config.ServiceEndpoint{URL: "http://gitea:3000", Token: "tok"},
			Coolify:    config.ServiceEndpoint{URL: "http://coolify:3000", Token: "tok"},
			Ollama:     config.ServiceEndpoint{URL: "http://ollama:11434"},
			OpenCode:   config.ServiceEndpoint{URL: "http://opencode:8080", Token: "tok"},
			Restic:     config.ResticConfig{Binary: "restic", RepoURL: "/repo", Password: "pw"},
			CodeServer: config.ServiceEndpoint{URL: "http://codeserver:8443"},
			Plane:      config.ServiceEndpoint{URL: "http://plane:8080"},
			Outline:    config.ServiceEndpoint{URL: "http://outline:3000", Token: "tok"},
			Mattermost: config.ServiceEndpoint{URL: "http://mattermost:8065", Token: "tok"},
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
	if clients.CodeServer == nil {
		t.Error("expected codeserver client")
	}
	if clients.Plane == nil {
		t.Error("expected plane client")
	}
	if clients.Outline == nil {
		t.Error("expected outline client")
	}
	if clients.Mattermost == nil {
		t.Error("expected mattermost client")
	}
}
