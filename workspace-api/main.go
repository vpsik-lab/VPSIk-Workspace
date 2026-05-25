package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vpsik/workspace-api/internal/config"
	"github.com/vpsik/workspace-api/internal/handler"
	"github.com/vpsik/workspace-api/internal/middleware"
)

func main() {
	configPath := "api.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := validateConfig(cfg); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	clients := handler.NewClients(cfg)
	authHandler := handler.NewAuthHandler(cfg)
	statusHandler := handler.NewStatusHandler(clients)
	proxyHandler := handler.NewProxyHandler(clients)

	protected := middleware.AuthMiddleware(cfg.Auth.JWTSecret)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.Handle("GET /api/auth/verify", protected(http.HandlerFunc(authHandler.Verify)))
	mux.Handle("GET /api/status", protected(http.HandlerFunc(statusHandler.Status)))

	mux.Handle("GET /api/gitea/repos", protected(http.HandlerFunc(proxyHandler.GiteaRepos)))
	mux.Handle("GET /api/gitea/issues", protected(http.HandlerFunc(proxyHandler.GiteaIssues)))

	mux.Handle("GET /api/coolify/servers", protected(http.HandlerFunc(proxyHandler.CoolifyServers)))
	mux.Handle("GET /api/coolify/deployments", protected(http.HandlerFunc(proxyHandler.CoolifyDeployments)))
	mux.Handle("GET /api/coolify/projects", protected(http.HandlerFunc(proxyHandler.CoolifyProjects)))
	mux.Handle("GET /api/coolify/projects/{project}/environments", protected(http.HandlerFunc(proxyHandler.CoolifyEnvironments)))
	mux.Handle("GET /api/coolify/projects/{project}/environments/{env}/applications", protected(http.HandlerFunc(proxyHandler.CoolifyApplications)))
	mux.Handle("POST /api/coolify/deploy", protected(http.HandlerFunc(proxyHandler.CoolifyDeployResource)))
	mux.Handle("POST /api/coolify/restart", protected(http.HandlerFunc(proxyHandler.CoolifyRestartResource)))
	mux.Handle("GET /api/coolify/deployments/{id}/logs", protected(http.HandlerFunc(proxyHandler.CoolifyDeploymentLogs)))
	mux.Handle("POST /api/coolify/env/get", protected(http.HandlerFunc(proxyHandler.CoolifyGetEnvVars)))
	mux.Handle("POST /api/coolify/env/update", protected(http.HandlerFunc(proxyHandler.CoolifyUpdateEnvVars)))

	mux.Handle("POST /api/gitea/webhooks", protected(http.HandlerFunc(proxyHandler.GiteaCreateWebhook)))
	mux.Handle("GET /api/gitea/repos/{repo}/webhooks", protected(http.HandlerFunc(proxyHandler.GiteaListWebhooks)))

	mux.Handle("GET /api/ollama/models", protected(http.HandlerFunc(proxyHandler.OllamaModels)))
	mux.Handle("POST /api/ollama/chat", protected(http.HandlerFunc(proxyHandler.OllamaChat)))
	mux.Handle("POST /api/ollama/pull", protected(http.HandlerFunc(proxyHandler.OllamaPullModel)))
	mux.Handle("DELETE /api/ollama/models/{name}", protected(http.HandlerFunc(proxyHandler.OllamaDeleteModel)))
	mux.Handle("POST /api/ollama/task", protected(http.HandlerFunc(proxyHandler.OllamaTask)))

	mux.Handle("POST /api/opencode/chat", protected(http.HandlerFunc(proxyHandler.OpenCodeChat)))

	mux.Handle("GET /api/restic/snapshots", protected(http.HandlerFunc(proxyHandler.ResticSnapshots)))
	mux.Handle("POST /api/restic/backup", protected(http.HandlerFunc(proxyHandler.ResticBackup)))
	mux.Handle("POST /api/restic/restore", protected(http.HandlerFunc(proxyHandler.ResticRestore)))
	mux.Handle("POST /api/restic/forget", protected(http.HandlerFunc(proxyHandler.ResticForget)))
	mux.Handle("GET /api/restic/check", protected(http.HandlerFunc(proxyHandler.ResticCheck)))
	mux.Handle("GET /api/restic/stats", protected(http.HandlerFunc(proxyHandler.ResticStats)))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	h := middleware.CORS(middleware.Logging(mux))

	server := &http.Server{
		Addr:         addr,
		Handler:      h,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Workspace API listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func validateConfig(cfg *config.APIConfig) error {
	if cfg.Auth.JWTSecret == "" || cfg.Auth.JWTSecret == "change-me" {
		return fmt.Errorf("jwt_secret must be set to a secure value")
	}
	if len(cfg.Auth.Users) == 0 {
		return fmt.Errorf("at least one user must be configured")
	}
	if cfg.Services.Gitea.URL == "" {
		return fmt.Errorf("gitea URL is required")
	}
	if cfg.Services.Coolify.URL == "" {
		return fmt.Errorf("coolify URL is required")
	}
	if cfg.Services.Ollama.URL == "" {
		return fmt.Errorf("ollama URL is required")
	}
	if cfg.Services.OpenCode.URL != "" && cfg.Services.OpenCode.Token == "" {
		return fmt.Errorf("opencode token is required when opencode URL is set")
	}
	if cfg.Services.Restic.RepoURL != "" && cfg.Services.Restic.Password == "" {
		return fmt.Errorf("restic password is required when restic repo URL is set")
	}
	return nil
}
