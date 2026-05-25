package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/vpsik/workspace-api/internal/client"
)

type ProxyHandler struct {
	clients *Clients
}

func NewProxyHandler(clients *Clients) *ProxyHandler {
	return &ProxyHandler{clients: clients}
}

func (h *ProxyHandler) GiteaRepos(w http.ResponseWriter, r *http.Request) {
	repos, err := h.clients.Gitea.ListRepos(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, repos)
}

func (h *ProxyHandler) GiteaIssues(w http.ResponseWriter, r *http.Request) {
	issues, err := h.clients.Gitea.ListIssues(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, issues)
}

func (h *ProxyHandler) CoolifyServers(w http.ResponseWriter, r *http.Request) {
	servers, err := h.clients.Coolify.ListServers(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, servers)
}

func (h *ProxyHandler) CoolifyDeployments(w http.ResponseWriter, r *http.Request) {
	deployments, err := h.clients.Coolify.ListDeployments(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, deployments)
}

func (h *ProxyHandler) CoolifyProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.clients.Coolify.ListProjects(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, projects)
}

func (h *ProxyHandler) CoolifyEnvironments(w http.ResponseWriter, r *http.Request) {
	projectUUID := r.PathValue("project")
	if projectUUID == "" {
		jsonError(w, "project UUID required", http.StatusBadRequest)
		return
	}
	envs, err := h.clients.Coolify.ListEnvironments(r.Context(), projectUUID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, envs)
}

func (h *ProxyHandler) CoolifyApplications(w http.ResponseWriter, r *http.Request) {
	projectUUID := r.PathValue("project")
	envName := r.PathValue("env")
	if projectUUID == "" || envName == "" {
		jsonError(w, "project UUID and env name required", http.StatusBadRequest)
		return
	}
	apps, err := h.clients.Coolify.ListApplications(r.Context(), projectUUID, envName)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, apps)
}

func (h *ProxyHandler) CoolifyDeployResource(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ResourceUUID string `json:"resource_uuid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.clients.Coolify.Deploy(r.Context(), req.ResourceUUID); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, map[string]string{"status": "deploying"})
}

func (h *ProxyHandler) CoolifyRestartResource(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ResourceUUID string `json:"resource_uuid"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.clients.Coolify.Restart(r.Context(), req.ResourceUUID); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, map[string]string{"status": "restarting"})
}

func (h *ProxyHandler) CoolifyDeploymentLogs(w http.ResponseWriter, r *http.Request) {
	deploymentID := r.PathValue("id")
	if deploymentID == "" {
		jsonError(w, "deployment ID required", http.StatusBadRequest)
		return
	}
	logs, err := h.clients.Coolify.GetDeploymentLogs(r.Context(), deploymentID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, map[string]string{"logs": logs})
}

type envVarsRequest struct {
	ProjectUUID string            `json:"project_uuid"`
	EnvName     string            `json:"env_name"`
	AppUUID     string            `json:"app_uuid"`
	EnvVars     map[string]string `json:"env_vars,omitempty"`
}

func (h *ProxyHandler) CoolifyGetEnvVars(w http.ResponseWriter, r *http.Request) {
	var req envVarsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	envVars, err := h.clients.Coolify.GetEnvVars(r.Context(), req.ProjectUUID, req.EnvName, req.AppUUID)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, envVars)
}

func (h *ProxyHandler) CoolifyUpdateEnvVars(w http.ResponseWriter, r *http.Request) {
	var req envVarsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.clients.Coolify.UpdateEnvVars(r.Context(), req.ProjectUUID, req.EnvName, req.AppUUID, req.EnvVars); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, map[string]string{"status": "updated"})
}

func (h *ProxyHandler) GiteaCreateWebhook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Repo    string `json:"repo"`
		URL     string `json:"url"`
		Secret  string `json:"secret,omitempty"`
		Events  []string `json:"events"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	hook, err := h.clients.Gitea.CreateWebhook(r.Context(), req.Repo, req.URL, req.Secret, req.Events)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, hook)
}

func (h *ProxyHandler) GiteaListWebhooks(w http.ResponseWriter, r *http.Request) {
	repo := r.PathValue("repo")
	if repo == "" {
		jsonError(w, "repo name required", http.StatusBadRequest)
		return
	}
	hooks, err := h.clients.Gitea.ListWebhooks(r.Context(), repo)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, hooks)
}

func (h *ProxyHandler) OllamaModels(w http.ResponseWriter, r *http.Request) {
	models, err := h.clients.Ollama.ListModels(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, models)
}

type chatRequest struct {
	Model    string              `json:"model"`
	Messages []client.ChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

func (h *ProxyHandler) OllamaChat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Stream {
		h.handleStreamingChat(w, r, req)
		return
	}

	reply, err := h.clients.Ollama.Chat(r.Context(), req.Model, req.Messages)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, map[string]string{"reply": reply})
}

func (h *ProxyHandler) handleStreamingChat(w http.ResponseWriter, r *http.Request, req chatRequest) {
	body, err := h.clients.Ollama.ChatStream(r.Context(), req.Model, req.Messages)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer body.Close()

	flusher, ok := w.(http.Flusher)
	if !ok {
		jsonError(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var chunk client.OllamaChatResponse
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			continue
		}

		data, _ := json.Marshal(map[string]interface{}{
			"content": chunk.Message.Content,
			"done":    chunk.Done,
		})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()

		if chunk.Done {
			break
		}
	}
}

func (h *ProxyHandler) OllamaPullModel(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Model string `json:"model"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if req.Model == "" {
		jsonError(w, "model name is required", http.StatusBadRequest)
		return
	}

	if err := h.clients.Ollama.PullModel(r.Context(), req.Model); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, map[string]string{"status": "pulling", "model": req.Model})
}

func (h *ProxyHandler) OllamaDeleteModel(w http.ResponseWriter, r *http.Request) {
	model := strings.TrimPrefix(r.URL.Path, "/api/ollama/models/")
	if model == "" {
		jsonError(w, "model name is required", http.StatusBadRequest)
		return
	}

	if err := h.clients.Ollama.DeleteModel(r.Context(), model); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, map[string]string{"status": "deleted", "model": model})
}

type taskRequest struct {
	Model    string `json:"model"`
	Task     string `json:"task"`
	Content  string `json:"content"`
}

func (h *ProxyHandler) OllamaTask(w http.ResponseWriter, r *http.Request) {
	var req taskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	prompts := map[string]string{
		"explain":   "Explain the following in simple terms:\n\n",
		"summarize": "Summarize the following concisely:\n\n",
		"review":    "Review this code for bugs, issues, and improvements:\n\n",
		"expand":    "Expand on the following with more detail:\n\n",
	}

	prompt, ok := prompts[req.Task]
	if !ok {
		jsonError(w, fmt.Sprintf("unknown task: %s", req.Task), http.StatusBadRequest)
		return
	}

	messages := []client.ChatMessage{
		{Role: "user", Content: prompt + req.Content},
	}

	if req.Model == "" {
		req.Model = "llama3.2"
	}

	reply, err := h.clients.Ollama.Chat(r.Context(), req.Model, messages)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, map[string]string{"reply": reply, "task": req.Task})
}

func (h *ProxyHandler) OpenCodeChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Message  string `json:"message"`
		Context  string `json:"context,omitempty"`
		RepoPath string `json:"repo_path,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	reply, err := h.clients.OpenCode.Chat(r.Context(), req.Message, req.Context, req.RepoPath)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}

	writeJSON(w, map[string]string{"reply": reply})
}

// ─── Restic / Backup ─────────────────────────────────────────────

func (h *ProxyHandler) ResticSnapshots(w http.ResponseWriter, r *http.Request) {
	snapshots, err := h.clients.Restic.ListSnapshots(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, snapshots)
}

func (h *ProxyHandler) ResticBackup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Paths []string `json:"paths"`
		Tags  []string `json:"tags,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if len(req.Paths) == 0 {
		jsonError(w, "at least one path required", http.StatusBadRequest)
		return
	}

	stats, err := h.clients.Restic.Backup(r.Context(), req.Paths, req.Tags)
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, stats)
}

func (h *ProxyHandler) ResticRestore(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SnapshotID string `json:"snapshot_id"`
		Target     string `json:"target"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.SnapshotID == "" || req.Target == "" {
		jsonError(w, "snapshot_id and target required", http.StatusBadRequest)
		return
	}

	if err := h.clients.Restic.Restore(r.Context(), req.SnapshotID, req.Target); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, map[string]string{"status": "restoring"})
}

func (h *ProxyHandler) ResticForget(w http.ResponseWriter, r *http.Request) {
	var req struct {
		KeepLast int      `json:"keep_last"`
		Tags     []string `json:"tags,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if req.KeepLast <= 0 {
		req.KeepLast = 7
	}

	if err := h.clients.Restic.Forget(r.Context(), req.KeepLast, req.Tags); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, map[string]string{"status": "forget completed"})
}

func (h *ProxyHandler) ResticCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.clients.Restic.Check(r.Context()); err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, map[string]string{"status": "ok"})
}

func (h *ProxyHandler) ResticStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.clients.Restic.Stats(r.Context())
	if err != nil {
		jsonError(w, err.Error(), http.StatusBadGateway)
		return
	}
	writeJSON(w, stats)
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
