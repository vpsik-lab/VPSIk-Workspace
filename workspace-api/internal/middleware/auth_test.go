package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vpsik/workspace-api/internal/auth"
)

func TestExtractToken_FromHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer test-token-123")

	token := extractToken(req)
	if token != "test-token-123" {
		t.Errorf("expected test-token-123, got %s", token)
	}
}

func TestExtractToken_FromCookie(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "workspace_token", Value: "cookie-token"})

	token := extractToken(req)
	if token != "cookie-token" {
		t.Errorf("expected cookie-token, got %s", token)
	}
}

func TestExtractToken_HeaderTakesPrecedence(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer header-token")
	req.AddCookie(&http.Cookie{Name: "workspace_token", Value: "cookie-token"})

	token := extractToken(req)
	if token != "header-token" {
		t.Errorf("expected header-token, got %s", token)
	}
}

func TestExtractToken_Empty(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	token := extractToken(req)
	if token != "" {
		t.Errorf("expected empty, got %s", token)
	}
}

func TestExtractToken_InvalidHeaderFormat(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "InvalidScheme token")

	token := extractToken(req)
	if token != "" {
		t.Errorf("expected empty for invalid format, got %s", token)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	jwtSecret := "test-secret"
	tokenStr, err := auth.GenerateToken(jwtSecret, "testuser")
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	handler := AuthMiddleware(jwtSecret)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.Context().Value(UsernameKey)
		if username != "testuser" {
			t.Errorf("expected testuser, got %v", username)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	handler := AuthMiddleware("secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	handler := AuthMiddleware("secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}
