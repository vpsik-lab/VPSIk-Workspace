# VPSIk Workspace

AI-native engineering workspace — intelligent installer, unified dashboard, AI layer, and deployment platform.

## Architecture

```
workspace-installer/     Go CLI — detect → plan → apply (idempotent)
workspace-api/           Go HTTP API server — proxy + auth + status
workspace-dashboard/     Next.js 14 dashboard — management UI
workspace-agents/        (future) AI agents
```

## Quick Start

```bash
# 1. Configure
cp workspace.example.yaml workspace.yaml
vim workspace.yaml

# 2. Bootstrap
cd workspace-installer && go build -o vpsik .
./vpsik init
./vpsik plan
./vpsik apply

# 3. Start API
cd workspace-api && go build -o workspace-api .
./workspace-api api.yaml

# 4. Start Dashboard
cd workspace-dashboard && npm run dev
```

## Docker

```bash
# Build and run all services
docker compose up -d

# Services:
#   Dashboard:     http://localhost:3000
#   API:           http://localhost:8081
#   Grafana:       http://localhost:3002
#   Prometheus:    http://localhost:9090
```

## Services

| Service | Port | Description |
|---------|------|-------------|
| Workspace Dashboard | 3000 | Management UI |
| API Server | 8081 | Backend proxy |
| Grafana | 3002 | Metrics dashboards |
| Prometheus | 9090 | Time-series metrics |
| Open WebUI | 3001 | LLM chat UI |
| Gitea | 3000 | Git hosting |
| Coolify | 8000 | Deployment platform |
| Ollama | 11434 | LLM runtime |

## Phases

1. **Infrastructure Bootstrap** — installer, compose templates, state
2. **Core Workspace** — API, auth, dashboard pages
3. **AI Workspace** — Ollama chat, OpenCode.ai, AI tasks
4. **Deployment Platform** — Coolify, Gitea webhooks
5. **Backup & Recovery** — Restic snapshots, restore, prune
6. **Monitoring** — Grafana, Prometheus, health checks
7. **CI/CD** — GitHub Actions, Docker builds

## Development

```bash
# Installer tests
cd workspace-installer && go test ./...

# API build
cd workspace-api && go build -o workspace-api .

# Dashboard build
cd workspace-dashboard && npm run build
```

## Configuration

See `workspace.example.yaml` for installer config and `api.example.yaml` for API config. Environment variables override YAML settings at runtime.
