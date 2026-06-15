# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Feature Specs

Baseline feature specification lives at `openspec/specs/spec.md` (index into 9 category files). Consult it before implementing or modifying a feature to find the relevant routes, API endpoints, config keys, and constraints.

## Build & Run

```bash
make build                    # Build everything (frontend + backend)
make backend                  # Backend only
make frontend                 # Frontend only
make watch                    # Dev hot-reload (both)
make watch-backend            # Dev hot-reload (backend only)

# Build with SQLite support (common for local dev)
TAGS="bindata sqlite sqlite_unlock_notify" make build

# Run the server
./gitea web
GITEA_RUN_MODE=dev ./gitea web
```

## Tests

```bash
make test                     # All tests
make test-backend             # Backend unit tests
make test-frontend            # Frontend tests

# Integration tests (require running DB)
make test-sqlite              # SQLite integration tests
make test-mysql               # MySQL integration tests
make test-pgsql               # PostgreSQL integration tests

# Run a single Go test
go test -run TestFunctionName ./path/to/package
make test-sqlite#TestFunctionName   # Single integration test

# E2E tests (Playwright)
make test-e2e-sqlite
```

## Lint & Format

```bash
make lint                     # Lint everything
make lint-fix                 # Auto-fix lint issues
make fmt                      # Format Go code
make tidy                     # go mod tidy
make checks                   # Consistency checks

# Individual linters
make lint-go                  # Go lint
make lint-js                  # JS lint
make lint-css                 # CSS lint
make lint-templates           # Go HTML templates
```

## Code Generation

```bash
make generate                 # All Go generate targets
make generate-swagger         # Regenerate Swagger spec
make svg                      # Regenerate SVG assets
```

## Architecture

### Layered Structure

The codebase follows a strict layered architecture. Dependencies flow downward:

```
cmd/       → CLI entry points, app initialization
routers/   → HTTP handlers (routers/web/ for UI, routers/api/v1/ for REST API)
services/  → Business logic (the orchestration layer)
models/    → Data access, XORM models, DB migrations
modules/   → Shared libraries and utilities (no dependency on models/services)
```

**Key rule**: `modules/` must not import from `models/`, `services/`, or `routers/`.

### Router

- Uses Chi router via `modules/web`
- Routes registered in `routers/web/` (UI) and `routers/api/v1/` (API)
- `routers/init.go` → `NormalRoutes()` sets up all route groups
- Middleware chain in `routers/common/`

### Database (XORM)

- ORM: XORM with engine in `models/db/engine.go`
- Models register via `init()` functions calling `RegisterModel()`
- DB sessions obtained through `db.DefaultContext` or explicit context
- Migrations in `models/migrations/` — version-numbered, sequential

### Templates & Frontend

- Go HTML templates in `templates/` (`.tmpl` files)
- Static assets embedded via bindata (`go generate`)
- Frontend source in `web_src/` — JS, CSS, Vue components, Fomantic UI
- Custom assets in `custom/public/` override built-in ones at runtime
- Template helpers in `modules/templates/`

### Configuration

- INI-based config loaded by `modules/setting/`
- Environment variables can override settings
- `LoadSettings()` at startup, `LoadSettingsForInstall()` for install wizard

### Initialization Sequence

`cmd/web.go` → `routers/init.go:InitWebInstalled()`:
1. Configuration loading
2. Database engine init
3. Git repository setup
4. Translation/i18n
5. Storage backends
6. Services (indexer, mirror, webhook, etc.)
7. SSH server

## Conventions

- PR title format: `<scope>: <description>`
- Go tests use `models/unittest/` framework for DB-backed test setup
- API tests go in `routers/api/v1/*_test.go`
- Integration tests in `tests/integration/`
- Swagger annotations on API handlers — run `make generate-swagger` after API changes
- This is Gitea 1.22.x branch — backport-focused, avoid unnecessary refactors

## Design Conventions

@design.md

cross-cutting principles: architecture & layering, naming, error handling, context propagation, logging, i18n, testing, security
