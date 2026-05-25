package handler

import (
	"encoding/json"
	"net/http"

	"github.com/vpsik/workspace-api/internal/client"
	"github.com/vpsik/workspace-api/internal/config"
)

type ProxyHandler struct {
	gitea   *client.GiteaClient
	coolify *client.CoolifyClient
	ollama  *client.OllamaClient
}

func NewProxyHandler(cfg *config.APIConfig) *ProxyHandler {
	return &ProxyHandler{
		gitea:   client.NewGiteaClient(cfg.Services.Gitea.URL, cfg.Services.Gitea.Token),
		coolify: client.NewCoolifyClient(cfg.Services.Coolify.URL, cfg.Services.Coolify.Token),
		ollama:  client.NewOllamaClient(cfg.Services.Ollama.URL),
	}
}

func (h *ProxyHandler) GiteaRepos(w http.ResponseWriter, r *http.Request) {
	repos, err := h.gitea.ListRepos(r.Context())
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repos)
}

func (h *ProxyHandler) GiteaIssues(w http.ResponseWriter, r *http.Request) {
	issues, err := h.gitea.ListIssues(r.Context())
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(issues)
}

func (h *ProxyHandler) CoolifyServers(w http.ResponseWriter, r *http.Request) {
	servers, err := h.coolify.ListServers(r.Context())
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(servers)
}

func (h *ProxyHandler) CoolifyDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := h.coolify.ListDeployments(r.Context())
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployments)
}

func (h *ProxyHandler) OllamaModels(w http.ResponseWriter, r *http.Request) {
	models, err := h.ollama.ListModels(r.Context())
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models)
}

type chatRequest struct {
	Model    string            `json:"model"`
	Messages []client.ChatMessage `json:"messages"`
}

func (h *ProxyHandler) OllamaChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	reply, err := h.ollama.Chat(r.Context(), req.Model, req.Messages)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"reply": reply})
}
