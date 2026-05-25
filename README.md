# VPSIk Workspace

**AI-native engineering workspace — One-command self-hosted stack for any VPS.**

```bash
curl -fsSL https://raw.githubusercontent.com/vpsik-lab/VPSIk-Workspace/main/install.sh | bash
```

---

## Architecture

```
Internet
   ↓
Cloudflare Tunnel (optional)  ← avoids port conflicts with Coolify
   ↓
Traefik (reverse proxy, SSL)
   ↓
workspace_net (isolated Docker network)
   ↓
┌──────────────────────────────────────────────────────┐
│  Dashboard  (workspace-dashboard / Next.js 14)      │
│  API Proxy  (workspace-api / Go)                    │
│  CLI        (workspace-installer / Go — "vpsik")    │
├──────────────────────────────────────────────────────┤
│  Authentik  (SSO / Identity Provider)               │
│  Gitea      (Git hosting)                           │
│  Ollama     (Local LLM runtime)                     │
│  OpenWebUI  (AI chat interface)                     │
│  Code-Server (VS Code in browser)                   │
│  Mattermost (Team communication)                    │
│  Outline    (Knowledge base / docs)                 │
│  Plane      (Project management)                    │
│  Coolify    (App deployment platform)               │
├──────────────────────────────────────────────────────┤
│  PostgreSQL (Central database)                      │
│  Redis      (Cache / message broker)                │
│  Restic     (Backup & restore)                      │
│  Grafana    (Metrics dashboards)                    │
│  Prometheus (Time-series metrics)                   │
└──────────────────────────────────────────────────────┘
```

**Key design decisions:**
- **Isolated**: `workspace_net` Docker network — no port conflicts
- **Portable**: Single `docker compose up -d` on any server
- **External**: Runs outside Coolify, communicates via API only
- **Self-contained**: All data in `/opt/workspace/` — easy backup & migration

---

## One-Command Install

```bash
# Minimal (defaults to workspace.vpsik.com):
curl -fsSL https://raw.githubusercontent.com/vpsik-lab/VPSIk-Workspace/main/install.sh | bash

# With custom domain:
curl -fsSL ...install.sh | bash -s -- --domain workspace.myvps.com

# Dry-run (see what would happen):
curl -fsSL ...install.sh | bash -s -- --dry-run
```

**What it does:**
1. Checks system (OS, Docker, RAM, Disk)
2. Installs dependencies if needed (Docker, Compose)
3. Creates `/opt/workspace/` directory structure
4. Downloads `vpsik` CLI binary from GitHub Releases
5. Runs `vpsik doctor --fix` (auto-fix issues)
6. Runs `vpsik init --auto` (generates config)
7. Runs `vpsik install --yes` (deploys all services)
8. Prints dashboard URLs

---

## Services

| Service | Internal URL | Traefik Host | Description |
|---------|-------------|--------------|-------------|
| Dashboard | `dashboard:3000` | `workspace.DOMAIN` | Management UI |
| API | `api:8081` | `api.workspace.DOMAIN` | Backend proxy |
| Authentik | `authentik:9000` | `auth.workspace.DOMAIN` | SSO / Identity |
| Gitea | `gitea:3000` | `git.workspace.DOMAIN` | Git hosting |
| Ollama | `ollama:11434` | `ollama.workspace.DOMAIN` | LLM runtime |
| OpenWebUI | `openwebui:8080` | `chat.workspace.DOMAIN` | AI chat UI |
| Code-Server | `codeserver:8443` | `code.workspace.DOMAIN` | VS Code in browser |
| Mattermost | `mattermost:8065` | `mattermost.workspace.DOMAIN` | Team chat |
| Outline | `outline:3000` | `docs.workspace.DOMAIN` | Knowledge base |
| Plane | `plane:8080` | `plane.workspace.DOMAIN` | Project management |
| Coolify | `coolify:3000` | `coolify.workspace.DOMAIN` | App deployment |
| Grafana | `grafana:3000` | `metrics.workspace.DOMAIN` | Dashboards |
| Prometheus | `prometheus:9090` | — | Metrics storage |
| PostgreSQL | `postgres:5432` | — | Database |
| Restic | — | — | Backup engine |

---

## CLI Reference (`vpsik`)

```
Usage:  vpsik [command] [flags]

Commands:
  init          Generate workspace configuration
    --auto        Non-interactive mode
    --domain      Workspace domain (default: workspace.vpsik.com)
    --services    Comma-separated services (default: all)

  plan          Scan environment and show installation plan

  install       Install missing services
    --yes         Skip confirmation prompt
    --dry-run     Show plan without installing

  status        Show current service status

  doctor        Check system health
    --fix         Auto-fix detected issues

  backup        Create backups
    --all         Backup all services
    --dry-run     Simulate backup

  restore       Restore from backup
    --latest      Restore latest snapshot
    --snapshot    Restore specific snapshot ID

  upgrade       Pull latest images and recreate services

  uninstall     Remove all deployed services
    --volumes     Remove persistent volumes
    --network     Remove Docker network

Global Flags:
  --config, -c  Config file path (default: workspace.yaml)
```

---

## Directory Structure

```
/opt/workspace/
├── compose/
│   ├── docker-compose.yml          # Generated compose file
│   └── docker-compose.override.yml # Dev overrides (opens ports)
├── configs/
│   ├── traefik/
│   │   ├── traefik.yml
│   │   ├── config.yml
│   │   └── acme.json
│   ├── cloudflared/
│   │   └── config.yml
│   ├── prometheus/
│   │   └── prometheus.yml
│   └── grafana/
│       └── dashboards/
├── data/
│   ├── postgres/
│   ├── gitea/
│   ├── authentik/
│   ├── ollama/
│   ├── codeserver/
│   └── restic/
├── backups/
│   └── (restic snapshots)
├── scripts/
│   ├── backup.sh
│   ├── restore.sh
│   └── healthcheck.sh
└── .env
```

---

## Roadmap

### Phase 1: Foundation  ﻿✅
- [x] Docker scanner & detector
- [x] Service templates (7 services)
- [x] Plan engine & state management
- [x] Traefik reverse proxy
- [x] API proxy server (Go)
- [x] Dashboard (Next.js 14)

### Phase 2: One-Command Install
- [ ] `install.sh` bootstrap script
- [ ] GitHub release workflow
- [ ] `vpsik init --auto`
- [ ] `vpsik install --yes`
- [ ] `vpsik doctor --fix`

### Phase 3: Backup & Restore
- [ ] `vpsik backup` command
- [ ] `vpsik restore` command
- [ ] Restic integration
- [ ] Service-specific backup strategies
- [ ] Scheduled backups

### Phase 4: Service Expansion
- [ ] Authentik SSO integration
- [ ] Code-Server
- [ ] Mattermost
- [ ] Outline
- [ ] Plane

### Phase 5: Monitoring & Observability
- [ ] Grafana dashboards per service
- [ ] Prometheus alerts
- [ ] Service health endpoints
- [ ] Log aggregation

### Phase 6: Enterprise
- [ ] Multi-server support
- [ ] Team management via Authentik
- [ ] Audit logging
- [ ] SLA monitoring

---

## Development

```bash
# Installer
cd workspace-installer && go build -o vpsik .
./vpsik init --auto --domain dev.local

# API
cd workspace-api && go build -o workspace-api .

# Dashboard
cd workspace-dashboard && npm install && npm run dev

# Tests
cd workspace-installer && go test ./...
cd workspace-api && go test ./...
cd workspace-dashboard && npm run build
```

---

## Configuration

See `workspace.example.yaml` for the full config reference.
Environment variables override YAML settings at runtime.

---

## License

MIT
