# VPSIk Workspace — Implementation Taskfile

> Standard self-hosted workspace for any VPS.
> Architecture: isolated Docker network + Traefik + optional Cloudflare Tunnel.

---

## Phase 1: Project Foundation (Documentation & Structure)

### 1.1 Documentation
- [x] README.md — rewritten with architecture, one-command install, roadmap
- [x] TASKFILE.md — this file
- [ ] `.opencode/plans/architecture.md` — full architecture plan document
- [ ] `workspace.example.yaml` — update with all new services
- [ ] `api.example.yaml` — update with new endpoints

### 1.2 Directory Structure
- [ ] Create `/opt/workspace/` layout in installer config
- [ ] Create `traefik/` config directory structure
- [ ] Create `scripts/` directory with helper scripts

### 1.3 CI/CD
- [ ] `.github/workflows/release.yml` — goreleaser for multi-platform bins
- [ ] `.github/workflows/ci.yml` — update with new service tests

---

## Phase 2: Installer Core (workspace-installer)

### 2.1 Config System (`pkg/config/config.go`)
- [ ] Add new service fields: CodeServer, Plane, Outline, Mattermost, Cloudflare
- [ ] Add `BackupConfig` struct (repository, schedule, keep-policy)
- [ ] Add `NetworkConfig` struct (name, proxy-port, use-tunnel)
- [ ] Add `SystemConfig` struct (install-path)
- [ ] Update `EnabledList()` with new services

### 2.2 Docker Compose Generator (`pkg/docker/compose.go`)
- [ ] Add `Expose` field to `ServiceCompose`
- [ ] Add `TraefikHost` and `TraefikPort` fields
- [ ] Add `DependsOn` field
- [ ] Add `ExtraLabels` field
- [ ] Add service templates:
  - [ ] `code-server` (coder/code-server:latest, expose 8443)
  - [ ] `plane` (makeplane/plane-app:latest, expose 8080)
  - [ ] `outline` (outlinewiki/outline:latest, expose 3000)
  - [ ] `mattermost` (mattermost/mattermost-team-edition:latest, expose 8065)
  - [ ] `cloudflared` (cloudflare/cloudflared:latest)
  - [ ] `redis` (redis:7-alpine, expose 6379)
- [ ] Update existing templates: `ports` → `expose` + Traefik labels
- [ ] Change network name from `vpsik` to `workspace_net`
- [ ] Rewrite `GenerateComposeFile` to use YAML marshaling
- [ ] Update `GenerateEnvFile` to include Traefik hosts
- [ ] Add `GenerateTraefikLabels` helper function

### 2.3 Scanner (`pkg/scanner/scanner.go`)
- [ ] Add `SystemInfo` struct (OS, Arch, CPU, RAM, Disk)
- [ ] Add `CoolifyDetected` field with Coolify info
- [ ] Add system info detection (OS, RAM, Disk)
- [ ] Add Coolify-specific detection (container, port, version)
- [ ] Add Docker Compose version detection
- [ ] Expand port scanning with more ports

### 2.4 Detector (`pkg/detector/detector.go`)
- [ ] Add detectors for: code-server, plane, outline, mattermost, cloudflared, redis
- [ ] Add `detectDockerCompose` function (detect services in compose files)
- [ ] Improve `detectContainer` to check Docker labels

### 2.5 Commands

#### `cmd/init.go`
- [ ] Add `--auto` flag for non-interactive mode
- [ ] Add `--domain` flag
- [ ] Add `--services` flag (comma-separated)
- [ ] Add `--output` flag for custom output path

#### `cmd/install.go`
- [ ] Add `--yes` flag to skip confirmation prompt
- [ ] Add automatic Coolify detection and install prompt
- [ ] Improve post-deploy health checks

#### `cmd/doctor.go` (new)
- [ ] System requirements check (Docker, Compose, RAM, Disk)
- [ ] Port conflict detection
- [ ] Network existence check
- [ ] NTP sync check
- [ ] `/opt/workspace` permissions check
- [ ] `--fix` flag for auto-repair

#### `cmd/backup.go` (new)
- [ ] `--all` flag for full backup
- [ ] Service-specific backup strategies
- [ ] Restic integration
- [ ] Backup listing

#### `cmd/restore.go` (new)
- [ ] `--latest` flag
- [ ] `--snapshot` flag with ID
- [ ] `--list` flag to show available snapshots
- [ ] Service-specific restore strategies

### 2.6 New Packages

#### `pkg/backup/`
- [ ] `backup.go` — backup orchestrator
- [ ] `executor.go` — per-service backup execution
- [ ] `strategies.go` — backup strategies for each service
- [ ] `restic.go` — Restic wrapper

#### `pkg/restore/`
- [ ] `restore.go` — restore orchestrator

#### `pkg/system/`
- [ ] `check.go` — system health checks
- [ ] `info.go` — OS, CPU, RAM, Disk detection
- [ ] `docker.go` — Docker install helpers
- [ ] `fix.go` — auto-fix implementations

---

## Phase 3: Root Compose & Network

### 3.1 Root `docker-compose.yml`
- [ ] Change network to `workspace_net`
- [ ] Add Traefik labels for all services
- [ ] Use `expose` instead of `ports` for all internal services
- [ ] Add `cloudflared` service (optional, disabled by default)
- [ ] Add `redis` service for authentik
- [ ] Add `traefik/acme.json` volume for SSL certs
- [ ] Create `docker-compose.override.yml` for dev mode

### 3.2 `.env.example`
- [ ] Add Traefik host variables for all services
- [ ] Add Authentik secrets
- [ ] Add Cloudflare Tunnel token
- [ ] Add backup config variables
- [ ] Add new service URLs

---

## Phase 4: API Expansion (workspace-api)

### 4.1 Config (`internal/config/config.go`)
- [ ] Add CodeServer, Plane, Outline, Mattermost endpoints
- [ ] Update validation for new services

### 4.2 Clients
- [ ] `internal/client/codeserver.go` — status, health
- [ ] `internal/client/plane.go` — projects, issues
- [ ] `internal/client/outline.go` — collections, docs
- [ ] `internal/client/mattermost.go` — channels, posts, health

### 4.3 Handlers (`internal/handler/clients.go`, `proxy.go`, `status.go`)
- [ ] Add new clients to `Clients` struct
- [ ] Add proxy handlers for all new services
- [ ] Add status checks for new services
- [ ] Register new routes in `main.go`

---

## Phase 5: Dashboard Expansion (workspace-dashboard)

### 5.1 New Pages
- [ ] `app/services/codeserver/page.tsx` — Code-Server status & launch
- [ ] `app/services/plane/page.tsx` — Plane projects overview
- [ ] `app/services/outline/page.tsx` — Outline docs access
- [ ] `app/services/mattermost/page.tsx` — Mattermost channels

### 5.2 Component Updates
- [ ] Update `Sidebar.tsx` with new service links
- [ ] Update `ServiceCard.tsx` if needed
- [ ] Update `api.ts` with new endpoints

---

## Phase 6: Install Script & CI/CD

### 6.1 `install.sh`
- [ ] OS detection (Ubuntu, Debian, Fedora, Arch)
- [ ] Docker installation (if missing + `--fix`)
- [ ] Docker Compose plugin installation
- [ ] Binary download from GitHub Releases
- [ ] SHA256 verification
- [ ] Run `vpsik doctor --fix`
- [ ] Run `vpsik init --auto --domain $DOMAIN`
- [ ] Run `vpsik install --yes`
- [ ] Print summary

### 6.2 `.github/workflows/release.yml`
- [ ] Build for linux/amd64, linux/arm64, darwin/amd64
- [ ] Create GitHub Release on tag push
- [ ] Upload binaries as release assets

---

## Phase 7: Testing & Verification

### 7.1 Update Tests
- [ ] Fix `compose_test.go` for new network name (`workspace_net`)
- [ ] Fix `compose_test.go` for new service templates count
- [ ] Add tests for new commands (doctor, backup, restore)
- [ ] Add tests for new scanner features (system info, Coolify)

### 7.2 Integration Test
- [ ] `vpsik init --auto --domain test.local`
- [ ] `vpsik doctor`
- [ ] `vpsik plan`
- [ ] `go test ./...` passes
- [ ] `go build -o vpsik .` succeeds

---

## Progress Tracking

| Phase | Tasks | Status |
|-------|-------|--------|
| 1: Foundation | 6 | 🔄 In progress |
| 2: Installer Core | 30+ | ⏳ Pending |
| 3: Root Compose | 6 | ⏳ Pending |
| 4: API Expansion | 6 | ⏳ Pending |
| 5: Dashboard | 6 | ⏳ Pending |
| 6: Install Script | 8 | ⏳ Pending |
| 7: Testing | 6 | ⏳ Pending |

---

*Last updated: 2026-05-25*
