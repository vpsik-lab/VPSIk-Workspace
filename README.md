# VPSIk WorkSpace

<p align="center">
  <b>AI-Native Engineering Workspace — One-Command Setup for Any VPS.</b>
</p>

<p align="center">
  <a href="https://github.com/vpsik-lab/VPSIk-Workspace">GitHub</a> •
  <a href="https://vpsik.com">Website</a> •
  <a href="mailto:opensource@vpsik.com">Contact</a>
</p>

<p align="center">
  <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT">
  <img src="https://img.shields.io/badge/Go-1.22-blue" alt="Go">
  <img src="https://img.shields.io/badge/Next.js-14-black" alt="Next.js">
  <img src="https://img.shields.io/badge/status-stable-green" alt="Status">
</p>

```bash
curl -fsSL https://raw.githubusercontent.com/vpsik-lab/VPSIk-Workspace/main/install.sh | bash
```

---

## Editions

| Feature | Open Source 🆓 | VPSIk Pro 💼 | WorkSpace SaaS ☁️ |
|---------|:-------------:|:-----------:|:----------------:|
| **Users** | Single | Team (multi-user) | Unlimited |
| **Hosting** | Self-hosted | Self-hosted | Hosted by VPSIk |
| **Price** | Free | License | Subscription |
| One-command install | ✅ | ✅ | — |
| Git hosting (Gitea) | ✅ | ✅ | ✅ |
| Local AI (Ollama) | ✅ | ✅ | ✅ |
| AI Chat (OpenWebUI) | ✅ | ✅ | ✅ |
| IDE in browser (Code-Server) | ✅ | ✅ | ✅ |
| Team communication (Mattermost) | ✅ | ✅ | ✅ |
| Knowledge base (Outline) | ✅ | ✅ | ✅ |
| Project management (Plane) | ✅ | ✅ | ✅ |
| App deployment (Coolify) | ✅ | ✅ | ✅ |
| SSO / Identity (Authentik) | ✅ | ✅ | ✅ |
| Monitoring (Grafana + Prometheus) | ✅ | ✅ | ✅ |
| Backup & restore (Restic) | ✅ | ✅ | ✅ |
| AI CTO — full project analysis | — | ✅ | ✅ |
| Team roles & permissions | — | ✅ | ✅ |
| Shared projects | — | ✅ | ✅ |
| Code review automation | — | ✅ | ✅ |
| AI development reports | — | ✅ | ✅ |
| Notifications (WhatsApp / Telegram / Slack / Notion) | — | ✅ | ✅ |
| Direct GitHub integration | — | — | ✅ |
| Multi-workspace | — | — | ✅ |
| Fully managed (no VPS needed) | — | — | ✅ |

---

## VPSIk Pro

Turn your workspace into a full-scale engineering company.

- **AI CTO** — AI analyzes your project as a virtual CTO: architecture reviews, tech debt detection, dependency analysis, performance recommendations.
- **Team Roles** — Owner, Admin, Developer, Reviewer with granular permissions.
- **Shared Projects** — Collaborate in real-time. Assign tasks, track progress, review code.
- **Code Review Automation** — Automated PR reviews with detailed reports and suggested fixes.
- **Development Reports** — AI-generated weekly/monthly reports covering progress, bottlenecks, and next steps.
- **Notifications** — Receive alerts and summaries via WhatsApp, Telegram, Slack, and Notion.

---

## WorkSpace SaaS

The fully-hosted cloud experience. No server, no setup.

- **Everything in Pro** — All Pro features included.
- **Direct GitHub Integration** — Sync repositories, manage issues, trigger workflows directly.
- **Multi-Workspace** — Separate workspaces for different teams or projects.
- **Notifications** — Full integration across WhatsApp, Telegram, Slack, and Notion.
- **Managed Infrastructure** — High-availability, automatic backups, SSL, monitoring — handled by VPSIk.
- **No VPS Required** — Sign up at [vpsik.com](https://vpsik.com) and start immediately.

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

## Prerequisites

| Requirement | Minimum | Recommended |
|------------|---------|-------------|
| RAM | 4 GB | 8 GB |
| Disk | 20 GB | 50 GB (SSD) |
| Docker | 24+ | Latest |
| OS | Linux (any) | Ubuntu 22.04+ / Debian 12 |

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

## VPSIk

**Website:** [vpsik.com](https://vpsik.com)
**GitHub:** [github.com/vpsik-lab](https://github.com/vpsik-lab)
**Email:** opensource@vpsik.com

Built by **Youssef Soliman** — [github.com/Ymasoli](https://github.com/Ymasoli)

---

## License

MIT
