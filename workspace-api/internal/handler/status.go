package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func httpCheck(url string) func(context.Context) error {
	return func(ctx context.Context) error {
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			return fmt.Errorf("unexpected status: %d", resp.StatusCode)
		}
		return nil
	}
}

type StatusHandler struct {
	clients *Clients
}

type ServiceStatus struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type StatusResponse struct {
	Services  []ServiceStatus `json:"services"`
	Timestamp string          `json:"timestamp"`
}

func NewStatusHandler(clients *Clients) *StatusHandler {
	return &StatusHandler{clients: clients}
}

func (h *StatusHandler) Status(w http.ResponseWriter, r *http.Request) {
	svcChecks := []struct {
		Name   string
		Check  func(context.Context) error
	}{
		{"gitea", h.clients.Gitea.CheckHealth},
		{"coolify", h.clients.Coolify.CheckHealth},
		{"ollama", h.clients.Ollama.CheckHealth},
		{"restic", h.clients.Restic.CheckHealth},
		{"code-server", h.clients.CodeServer.CheckHealth},
		{"plane", h.clients.Plane.CheckHealth},
		{"outline", h.clients.Outline.CheckHealth},
		{"mattermost", h.clients.Mattermost.CheckHealth},
	}

	if h.clients.OpenCode != nil {
		svcChecks = append(svcChecks, struct {
			Name   string
			Check  func(context.Context) error
		}{"opencode", h.clients.OpenCode.CheckHealth})
	}

	svcChecks = append(svcChecks,
		struct {
			Name  string
			Check func(context.Context) error
		}{"grafana", httpCheck("http://grafana:3000/api/health")},
		struct {
			Name  string
			Check func(context.Context) error
		}{"prometheus", httpCheck("http://prometheus:9090/-/healthy")},
	)

	services := svcChecks

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
