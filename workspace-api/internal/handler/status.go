package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/vpsik/workspace-api/internal/client"
	"github.com/vpsik/workspace-api/internal/config"
)

type StatusHandler struct {
	gitea   *client.GiteaClient
	coolify *client.CoolifyClient
	ollama  *client.OllamaClient
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type StatusResponse struct {
	Services []ServiceStatus `json:"services"`
	Timestamp string         `json:"timestamp"`
}

func NewStatusHandler(cfg *config.APIConfig) *StatusHandler {
	return &StatusHandler{
		gitea:   client.NewGiteaClient(cfg.Services.Gitea.URL, cfg.Services.Gitea.Token),
		coolify: client.NewCoolifyClient(cfg.Services.Coolify.URL, cfg.Services.Coolify.Token),
		ollama:  client.NewOllamaClient(cfg.Services.Ollama.URL),
	}
}

func (h *StatusHandler) Status(w http.ResponseWriter, r *http.Request) {
	services := []struct {
		Name   string
		Check  func(context.Context) error
	}{
		{"gitea", h.gitea.CheckHealth},
		{"coolify", h.coolify.CheckHealth},
		{"ollama", h.ollama.CheckHealth},
	}

	results := make([]ServiceStatus, len(services))
	var wg sync.WaitGroup
	wg.Add(len(services))

	for i, svc := range services {
		i, svc := i, svc
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()

			status := ServiceStatus{Name: svc.Name, Status: "healthy"}
			if err := svc.Check(ctx); err != nil {
				status.Status = "unhealthy"
				status.Error = err.Error()
			}
			results[i] = status
		}()
	}

	wg.Wait()

	resp := StatusResponse{
		Services:  results,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
