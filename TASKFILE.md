# VPSIk WorkSpace — Implementation Taskfile

> Standard self-hosted workspace for any VPS.
> Architecture: isolated Docker network + Traefik + optional Cloudflare Tunnel.

---

## Open Source Edition — ✅ Complete

The open source edition is fully implemented and functional.

### Implemented Features

- [x] Docker scanner & detector
- [x] Service templates (19 services)
- [x] Plan engine & state management
- [x] Traefik reverse proxy with SSL
- [x] API proxy server (Go)
- [x] Dashboard (Next.js 14)
- [x] `install.sh` bootstrap script
- [x] `vpsik init --auto` & `vpsik install --yes`
- [x] `vpsik doctor --fix` system health
- [x] `vpsik backup` & `vpsik restore` (Restic)
- [x] `vpsik upgrade` & `vpsik uninstall`
- [x] All service integrations (Authentik, Gitea, Ollama, OpenWebUI, Code-Server, Mattermost, Outline, Plane, Coolify, Grafana, Prometheus)
- [x] Network isolation via `workspace_net`
- [x] Env file generation with auto-generated secrets
- [x] API config generation with bcrypt password hashing
- [x] Post-deploy health checks
- [x] State persistence (JSON)

### Commands

```
vpsik init      Generate configuration (interactive or --auto)
vpsik plan      Show installation plan
vpsik install   Deploy services (--yes for silent, --dry-run)
vpsik status    Show service health
vpsik doctor    System checks (--fix for auto-repair)
vpsik upgrade   Pull latest images & recreate
vpsik uninstall Remove all services (--volumes, --network)
vpsik backup    Create backups (--all, --dry-run, --list)
vpsik restore   Restore from snapshots (--latest, --snapshot)
```

---

## VPSIk Pro Edition — 🔄 Planned

Team-focused edition with AI-powered features.

- [ ] Multi-user authentication & team roles
- [ ] AI CTO — full project analysis & architecture review
- [ ] Code review automation with AI reports
- [ ] Shared projects & real-time collaboration
- [ ] Notifications (WhatsApp, Telegram, Slack, Notion)
- [ ] Team permissions & access control

---

## WorkSpace SaaS Edition — 🔄 Planned

Fully-hosted cloud edition managed by VPSIk.

- [ ] Everything in Pro
- [ ] Direct GitHub integration
- [ ] Multi-workspace support
- [ ] Managed hosting (no VPS required)
- [ ] High-availability infrastructure

---

## Progress Tracking

| Edition | Status |
|---------|--------|
| Open Source | ✅ Complete |
| VPSIk Pro | 🔄 Planned |
| WorkSpace SaaS | 🔄 Planned |

---

*Last updated: 2026-05-26*
