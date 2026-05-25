package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

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

	authHandler := handler.NewAuthHandler(cfg)
	statusHandler := handler.NewStatusHandler(cfg)
	proxyHandler := handler.NewProxyHandler(cfg)

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

	mux.Handle("GET /api/ollama/models", protected(http.HandlerFunc(proxyHandler.OllamaModels)))
	mux.Handle("POST /api/ollama/chat", protected(http.HandlerFunc(proxyHandler.OllamaChat)))

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)

	handler := middleware.CORS(mux)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
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
}
