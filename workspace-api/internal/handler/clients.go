package handler

import (
	"github.com/vpsik/workspace-api/internal/client"
	"github.com/vpsik/workspace-api/internal/config"
)

type Clients struct {
	Gitea    *client.GiteaClient
	Coolify  *client.CoolifyClient
	Ollama   *client.OllamaClient
	OpenCode *client.OpenCodeClient
	Restic   *client.ResticClient
}

func NewClients(cfg *config.APIConfig) *Clients {
	return &Clients{
		Gitea:    client.NewGiteaClient(cfg.Services.Gitea.URL, cfg.Services.Gitea.Token),
		Coolify:  client.NewCoolifyClient(cfg.Services.Coolify.URL, cfg.Services.Coolify.Token),
		Ollama:   client.NewOllamaClient(cfg.Services.Ollama.URL),
		OpenCode: client.NewOpenCodeClient(cfg.Services.OpenCode.URL, cfg.Services.OpenCode.Token),
		Restic:   client.NewResticClient(cfg.Services.Restic.Binary, cfg.Services.Restic.RepoURL, cfg.Services.Restic.Password),
	}
}
