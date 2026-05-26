#!/usr/bin/env bash
set -euo pipefail

# ╔══════════════════════════════════════════════════════════════╗
# ║  WorkSpace OS — One-Command Installer                       ║
# ║  Usage: curl -fsSL https://raw.githubusercontent.com/       ║
# ║    vpsik-lab/VPSIk-Workspace/main/install.sh | bash         ║
# ╚══════════════════════════════════════════════════════════════╝

APP="WorkSpace OS"
INSTALL_DIR="/opt/workspace"
BIN_DIR="/usr/local/bin"
REPO="vpsik-lab/VPSIk-Workspace"
BRANCH="main"
VERSION="latest"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log()  { echo -e "${GREEN}✓${NC} $1"; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }
err()  { echo -e "${RED}❌${NC} $1"; }
info() { echo -e "${BLUE}ℹ${NC} $1"; }

# Parse flags
DOMAIN=""
DRY_RUN=false
VERBOSE=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --domain)    DOMAIN="$2"; shift 2 ;;
    --dry-run)   DRY_RUN=true; shift ;;
    --verbose)   VERBOSE=true; shift ;;
    --help|-h)   echo "Usage: $0 [--domain DOMAIN] [--dry-run] [--verbose]"; exit 0 ;;
    *)           warn "Unknown option: $1"; shift ;;
  esac
done

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║    ${APP} Installer v0.1           ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════╝${NC}"
echo ""

# ── Prerequisites ──────────────────────────────────────────

info "Checking prerequisites..."

# Check curl
if ! command -v curl &>/dev/null; then
  err "curl is required. Install it first."
  exit 1
fi
log "curl available"

# Check git
if ! command -v git &>/dev/null; then
  warn "git not found — will use direct download"
fi

# ── Docker ─────────────────────────────────────────────────

if command -v docker &>/dev/null; then
  log "Docker: $(docker --version 2>/dev/null || echo 'installed')"
else
  warn "Docker not found"
  if [ "$DRY_RUN" = false ]; then
    info "Installing Docker..."
    curl -fsSL https://get.docker.com | sh
    sudo usermod -aG docker "$USER" 2>/dev/null || true
    log "Docker installed"
  else
    info "Would install Docker"
  fi
fi

# ── Docker Compose ─────────────────────────────────────────

if docker compose version &>/dev/null; then
  log "Docker Compose: $(docker compose version --short 2>/dev/null || echo 'installed')"
elif docker-compose --version &>/dev/null; then
  warn "legacy docker-compose found — install plugin"
else
  warn "Docker Compose not found"
  if [ "$DRY_RUN" = false ]; then
    info "Installing Docker Compose plugin..."
    DOCKER_CONFIG="${DOCKER_CONFIG:-$HOME/.docker}"
    mkdir -p "$DOCKER_CONFIG/cli-plugins"
    curl -SL "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
      -o "$DOCKER_CONFIG/cli-plugins/docker-compose"
    chmod +x "$DOCKER_CONFIG/cli-plugins/docker-compose"
    log "Docker Compose installed"
  fi
fi

# ── Create directories ─────────────────────────────────────

if [ "$DRY_RUN" = false ]; then
  info "Creating directory structure..."
  sudo mkdir -p "$INSTALL_DIR"/{compose,configs/{traefik,cloudflared,prometheus,grafana/dashboards},data,backups,scripts}
  sudo chown -R "$USER:$USER" "$INSTALL_DIR" 2>/dev/null || true
  log "Directory structure created at $INSTALL_DIR"
else
  info "Would create directory structure at $INSTALL_DIR"
fi

# ── Download workspace binary ──────────────────────────────

if [ "$DRY_RUN" = false ]; then
  info "Downloading workspace binary..."
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  ARCH=$(uname -m)
  case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    armv7l)  ARCH="armv7" ;;
  esac

  BINARY_URL="https://github.com/$REPO/releases/latest/download/workspace-${OS}-${ARCH}"
  if curl -fsSL "$BINARY_URL" -o /tmp/workspace 2>/dev/null; then
    chmod +x /tmp/workspace
    sudo mv /tmp/workspace "$BIN_DIR/workspace"
    log "workspace binary installed to $BIN_DIR/workspace"

    # Clone source to build Docker images
    BUILD_TMP=$(mktemp -d)
    git clone --depth 1 "https://github.com/$REPO.git" "$BUILD_TMP"
    docker build --no-cache -t workspaceos-api:latest "$BUILD_TMP/workspace-api" 2>&1 | tail -5
    docker build --no-cache -t workspaceos-dashboard:latest "$BUILD_TMP/workspace-dashboard" 2>&1 || true
    rm -rf "$BUILD_TMP"
    log "Local Docker images built"
  else
    # Fallback: build from source
    warn "Binary download failed, building from source..."
    if command -v go &>/dev/null; then
      TMPDIR=$(mktemp -d)
      git clone --depth 1 "https://github.com/$REPO.git" "$TMPDIR"
      cd "$TMPDIR/workspace-installer"
      go build -o workspace .
      sudo mv workspace "$BIN_DIR/workspace"
      log "workspace built and installed to $BIN_DIR/workspace"

      # Build local Docker images for API and Dashboard
      info "Building workspace-api Docker image..."
      docker build --no-cache -t workspaceos-api:latest "$TMPDIR/workspace-api" 2>&1 | tail -5
      info "Building workspace-dashboard Docker image..."
      docker build --no-cache -t workspaceos-dashboard:latest "$TMPDIR/workspace-dashboard" 2>&1 || true
      log "Local Docker images built"

      rm -rf "$TMPDIR"
    else
      err "Go not available and binary download failed. Install Go or download manually."
      exit 1
    fi
  fi
else
  info "Would download workspace binary"
fi

# ── Run workspace doctor ───────────────────────────────────

if [ "$DRY_RUN" = false ] && command -v workspace &>/dev/null; then
  info "Running system health check..."
  if [ "$VERBOSE" = true ]; then
    workspace doctor --fix
  else
    workspace doctor --fix 2>&1 | grep -E '(✅|⚠|❌|error|Error)' || true
  fi
  log "System health check complete"
else
  info "Would run: workspace doctor --fix"
fi

# ── Generate config ────────────────────────────────────────

if [ "$DRY_RUN" = false ] && command -v workspace &>/dev/null; then
  info "Generating workspace configuration..."

  SERVICES_ARGS=""
  if [ -n "$DOMAIN" ]; then
    WORKSPACE_DOMAIN="$DOMAIN"
  else
    WORKSPACE_DOMAIN="workspace.vpsik.com"
  fi

  # Use --auto for silent generation
  cd "$INSTALL_DIR"
  workspace init --auto --domain "$WORKSPACE_DOMAIN" --output "$INSTALL_DIR/workspace.yaml" $SERVICES_ARGS
  log "Configuration generated for domain $WORKSPACE_DOMAIN"
else
  info "Would generate config for domain ${DOMAIN:-workspace.vpsik.com}"
fi

# ── Install ────────────────────────────────────────────────

if [ "$DRY_RUN" = false ] && command -v workspace &>/dev/null; then
  info "Deploying services..."
  cd "$INSTALL_DIR"
  workspace install --yes --config "$INSTALL_DIR/workspace.yaml"
  log "Services deployed"
else
  info "Would run: workspace install --yes"
fi

# ── Detect IP ─────────────────────────────────────────────

VPS_IP=""
if command -v curl &>/dev/null; then
  VPS_IP=$(curl -s -4 --max-time 5 ifconfig.me 2>/dev/null || true)
fi
if [ -z "$VPS_IP" ] && command -v hostname &>/dev/null; then
  VPS_IP=$(hostname -I | awk '{print $1}')
fi
if [ -z "$VPS_IP" ]; then
  VPS_IP="<VPS_IP>"
fi

# ── Summary ────────────────────────────────────────────────

echo ""
echo -e "${BLUE}╔════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║    Installation Complete! 🎉               ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════╝${NC}"
echo ""
echo "  ── Access ──"
echo "  Dashboard:   http://$VPS_IP:3000"
echo "  API:         http://$VPS_IP:8081"
echo "  Code Server: http://$VPS_IP:8443"
echo ""
echo "  ── Login ──"
echo "  Username: admin"
echo "  Password: admin"
echo ""
echo "  ── Paths ──"
echo "  Config:      $INSTALL_DIR/workspace.yaml"
echo "  Data:        $INSTALL_DIR/data/"
echo "  Backups:     $INSTALL_DIR/backups/"
echo ""

if [ "$DRY_RUN" = true ]; then
  warn "Dry-run mode — no changes were made."
fi

info "Run 'workspace status' to check service health."
info "Run 'workspace backup --all' to create first backup."
echo ""
