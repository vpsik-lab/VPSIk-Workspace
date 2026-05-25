# VPSIk Workspace — Architecture Plan

## Vision

A **one-command, self-contained engineering workspace** that runs on any VPS with Docker.
Standard, portable, and isolated from other services (Coolify, existing stacks).

---

## Core Principles

1. **Single command install**: `curl ... install.sh | bash`
2. **Zero port conflicts**: All services internal via `workspace_net`, only Traefik exposes 80/443 (or Cloudflare Tunnel)
3. **Works anywhere**: Ubuntu, Debian, Fedora, Arch — any VPS with Docker
4. **Coolify-friendly**: Detects Coolify, avoids its ports, communicates via API
5. **Portable**: `/opt/workspace/` contains everything — backup, move, restore
6. **Incremental**: Start minimal, expand later

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Internet                                 │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │  Cloudflare Tunnel   │  ← Optional: avoids port 80/443
              │  (cloudflared)       │    conflicts with Coolify
              └──────────┬───────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │  Traefik v3          │  ← Reverse proxy, SSL (80/443)
              │  (vpsik-traefik)     │
              └──────────┬───────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │   workspace_net      │  ← Isolated Docker bridge network
              │   (172.x.x.x/16)     │
              └──────────────────────┘
                         │
         ┌───────────────┼───────────────────┐
         ▼               ▼                   ▼
  ┌────────────┐  ┌────────────┐   ┌────────────────┐
  │ Identity   │  │ Core Dev   │   │ Collaboration  │
  │ ─────────  │  │ ─────────  │   │ ─────────────  │
  │ Authentik  │  │ Gitea      │   │ Mattermost     │
  │ PostgreSQL │  │ Ollama     │   │ Outline        │
  │ Redis      │  │ OpenWebUI  │   │ Plane          │
  │            │  │ Code-Server│   │                │
  └────────────┘  └────────────┘   └────────────────┘
         │               │                  │
         └───────────────┼──────────────────┘
                         ▼
              ┌──────────────────────┐
              │  Management Layer    │
              │  ─────────────────── │
              │  Dashboard (Next.js) │
              │  API Proxy (Go)      │
              │  CLI (Go — vpsik)    │
              │  Restic (Backup)     │
              │  Grafana/Prometheus  │
              └──────────────────────┘
```

---

## Component Details

### 1. CLI — `vpsik` (workspace-installer)

The single source of truth. The CLI:
- Scans the server (OS, Docker, resources)
- Detects running services (containers, ports, APIs)
- Generates `docker-compose.yml` from templates
- Prevents port conflicts (detects in-use ports, adjusts)
- Deploys, upgrades, backs up, restores

**Commands:**
```
init      → Generate config (interactive or --auto)
plan      → Show what needs to change
install   → Apply the plan (with --yes for silent)
status    → Show current service states
doctor    → Health check + auto-fix
backup    → Create restic backups
restore   → Restore from restic snapshots
upgrade   → Pull latest images + recreate
uninstall → Tear down everything
```

### 2. Docker Compose Layout

Each service uses `expose:` (internal only) instead of `ports:`.
Traefik handles all external routing via labels.

```yaml
services:
  gitea:
    image: gitea/gitea:latest
    expose: ["3000"]                        # ← Internal only
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.gitea.rule=Host(`git.DOMAIN`)"
      - "traefik.http.services.gitea.loadbalancer.server.port=3000"
    networks: [workspace_net]
```

### 3. Network Isolation

- **Network**: `workspace_net` (custom bridge, not `vpsik`)
- **No ports exposed** except Traefik (80/443) and optionally Cloudflare Tunnel
- **No conflict** with Coolify's own Docker network or ports

### 4. Traefik Configuration

```yaml
# Dynamic config for each service:
- "traefik.enable=true"
- "traefik.http.routers.{service}.rule=Host(`{subdomain}.{DOMAIN}`)"
- "traefik.http.routers.{service}.entrypoints=websecure"
- "traefik.http.services.{service}.loadbalancer.server.port={port}"
```

### 5. Backup Strategy (Restic)

Each service has a specific backup strategy:

| Service | Method | Frequency |
|---------|--------|-----------|
| PostgreSQL | `pg_dump` → restic | Daily 03:00 |
| Gitea | `gitea dump` → restic | Daily 04:00 |
| Authentik | DB dump + media | Daily 05:00 |
| Ollama | restic volume | Weekly |
| Code-Server | restic /home/coder | Daily |
| Configs | restic /opt/workspace/configs | Daily |

### 6. Coolify Integration

- **No direct integration**: Coolify runs in its own space
- **Detection**: `vpsik` detects Coolify containers and ports
- **Conflict avoidance**: If Coolify uses 80/443, VPSIk uses Cloudflare Tunnel
- **Remote management**: Via Coolify API endpoints in workspace-api
- **Auto-install**: If Coolify not found, `vpsik install` offers to deploy it

---

## Service Templates (Complete)

| Service | Image | Expose | Traefik Host | Depends On |
|---------|-------|--------|-------------|------------|
| traefik | traefik:v3.0 | 80,443 | — | — |
| postgres | postgres:16-alpine | 5432 | — | — |
| redis | redis:7-alpine | 6379 | — | — |
| authentik-server | goauthentik/server | 9000 | auth.DOMAIN | postgres, redis |
| authentik-worker | goauthentik/server | — | — | authentik-server |
| gitea | gitea/gitea | 3000 | git.DOMAIN | postgres |
| ollama | ollama/ollama | 11434 | ollama.DOMAIN | — |
| openwebui | open-webui/open-webui | 8080 | chat.DOMAIN | ollama |
| code-server | coder/code-server | 8443 | code.DOMAIN | — |
| mattermost | mattermost/mattermost | 8065 | mattermost.DOMAIN | postgres |
| outline | outlinewiki/outline | 3000 | docs.DOMAIN | postgres, redis |
| plane | makeplane/plane | 8080 | plane.DOMAIN | postgres, redis |
| coolify | coollabsio/coolify | 3000 | coolify.DOMAIN | — |
| grafana | grafana/grafana | 3000 | metrics.DOMAIN | prometheus |
| prometheus | prom/prometheus | 9090 | — | — |
| restic | restic/restic | — | — | — |
| workspace-api | (build) | 8081 | api.DOMAIN | postgres |
| workspace-dashboard | (build) | 3000 | workspace.DOMAIN | api |
| cloudflared | cloudflare/cloudflared | — | — | — |

---

## Installation Flow

```
User runs: curl ... install.sh | bash

install.sh:
  ├── Detect OS → install Docker if needed
  ├── Install Docker Compose plugin
  ├── Download vpsik binary from GitHub Releases
  ├── Create /opt/workspace/ structure
  ├── vpsik doctor --fix
  │     ├── Check Docker running
  │     ├── Check ports (80, 443, etc.)
  │     ├── Check RAM > 2GB, Disk > 10GB
  │     ├── Check NTP sync
  │     └── Auto-fix if possible
  ├── vpsik init --auto --domain $DOMAIN
  │     └── Generate /opt/workspace/compose/docker-compose.yml
  ├── vpsik install --yes
  │     ├── Create network workspace_net
  │     ├── Pull images
  │     ├── Deploy containers
  │     ├── Health checks
  │     └── Save state
  └── Print success message with URLs
```

---

## File System Layout

```
/opt/workspace/
├── compose/
│   ├── docker-compose.yml
│   └── docker-compose.override.yml
├── configs/
│   ├── traefik/
│   │   ├── traefik.yml
│   │   ├── config.yml
│   │   └── acme.json         (mode 600)
│   ├── cloudflared/
│   │   └── config.yml
│   ├── prometheus/
│   │   └── prometheus.yml
│   └── grafana/
│       └── dashboards/
│           └── default.json
├── data/
│   ├── postgres/
│   ├── gitea/
│   ├── authentik/
│   │   ├── media/
│   │   └── certs/
│   ├── ollama/
│   │   └── models/
│   ├── codeserver/
│   ├── mattermost/
│   ├── outline/
│   ├── plane/
│   └── restic/
├── backups/
│   └── (restic repository)
├── scripts/
│   ├── backup.sh
│   ├── restore.sh
│   └── healthcheck.sh
└── .env
```

---

## Environment Variables

All configuration via `.env` file:

```env
# Domain
VPSIK_DOMAIN=workspace.vpsik.com

# Traefik subdomains (auto-generated from domain)
TRAEFIK_AUTH_HOST=auth.${VPSIK_DOMAIN}
TRAEFIK_GIT_HOST=git.${VPSIK_DOMAIN}
TRAEFIK_OLLAMA_HOST=ollama.${VPSIK_DOMAIN}
TRAEFIK_CHAT_HOST=chat.${VPSIK_DOMAIN}
TRAEFIK_CODE_HOST=code.${VPSIK_DOMAIN}
TRAEFIK_DOCS_HOST=docs.${VPSIK_DOMAIN}
TRAEFIK_PLANE_HOST=plane.${VPSIK_DOMAIN}
TRAEFIK_MATTERMOST_HOST=mattermost.${VPSIK_DOMAIN}
TRAEFIK_METRICS_HOST=metrics.${VPSIK_DOMAIN}
TRAEFIK_COOLIFY_HOST=coolify.${VPSIK_DOMAIN}
TRAEFIK_API_HOST=api.${VPSIK_DOMAIN}

# Database
POSTGRES_USER=vpsik
POSTGRES_PASSWORD=<generate>
POSTGRES_DB=vpsik

# Authentik
AUTHENTIK_SECRET_KEY=<generate>
AUTHENTIK_BOOT_PASSWORD=<generate>

# Services
GITEA_URL=http://gitea:3000
COOLIFY_URL=http://coolify:3000
OLLAMA_URL=http://ollama:11434
CODESERVER_URL=http://codeserver:8443
PLANE_URL=http://plane:8080
OUTLINE_URL=http://outline:3000
MATTERMOST_URL=http://mattermost:8065

# Cloudflare
CF_TUNNEL_TOKEN=

# Backup
RESTIC_REPOSITORY=/opt/workspace/backups
RESTIC_PASSWORD=<generate>
```

---

## Port Conflict Resolution

| Scenario | Resolution |
|----------|-----------|
| Coolify on 80/443 | Enable Cloudflare Tunnel, Traefik on 8080/8443 |
| Port already in use | Scan changes port mapping |
| Multiple Docker networks | Use isolated `workspace_net` |
| Firewall blocks ports | `vpsik doctor --fix` opens them |

---

## Development Workflow

```bash
# 1. Clone
git clone https://github.com/vpsik-lab/VPSIk-Workspace.git
cd VPSIk-Workspace

# 2. Build installer
cd workspace-installer && go build -o vpsik .

# 3. Init for local dev
./vpsik init --auto --domain dev.local

# 4. Build & run API
cd ../workspace-api && go build -o workspace-api .

# 5. Build & run dashboard
cd ../workspace-dashboard && npm install && npm run dev
```

---

## Release Process

1. Tag: `git tag v0.1.0 && git push origin v0.1.0`
2. GitHub Actions builds binaries for 3 platforms
3. Creates GitHub Release with assets
4. `install.sh` points to latest release
5. User runs one-command install
