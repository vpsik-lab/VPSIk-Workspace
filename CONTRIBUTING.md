# Contributing to WorkSpace OS

Thank you for your interest! We welcome contributions of all forms.

## Getting Started

1. Fork the repository
2. Clone your fork
3. Run `./install.sh --dry-run` to see what the installer does
4. Run `cd workspace-installer && go test ./...` to run tests

## Development Workflow

```bash
# Installer (Go)
cd workspace-installer
go build -o workspace .
./workspace init --auto --domain dev.local

# API (Go)
cd workspace-api
go build -o workspace-api .
./workspace-api test-api.yaml

# Dashboard (Next.js)
cd workspace-dashboard
npm install
npm run dev
```

## Code Style

- **Go**: Run `go vet ./...` and `golangci-lint run` before committing
- **Next.js**: Run `npm run lint` before committing
- Follow existing code patterns and conventions
- No commented-out code; no debug logs in production paths

## Pull Request Process

1. Create a feature branch from `main`
2. Write tests for new code
3. Ensure all tests pass: `go test ./...`
4. Update documentation if needed
5. Open a PR with a clear title and description

## Commit Messages

Use conventional commits:
- `fix:` — bug fix
- `feat:` — new feature
- `test:` — test additions/changes
- `docs:` — documentation
- `refactor:` — code restructuring
- `sec:` — security fix

## Questions?

Open a GitHub Discussion or issue.
