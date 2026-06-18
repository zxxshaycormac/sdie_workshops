# CLAUDE.md

Guidance for Claude Code (claude.ai/code) and other AI assistants working in this repository.

## SDIE Workflow Map

This repository follows the SDIE methodology (Specification → Design → Implementation → Evaluation), a maintenance framework for brownfield software like Gitea. The four phases are not linear — they overlap, validate each other, and feed back.

- **(S)pec** — Define what is being changed and why, via OpenSpec.
- **(D)esign** — Translate the spec into architecture, contracts, and tasks.
- **(I)mplementation** — Write code and tests against the design.
- **(E)val** — Validate the change satisfies the original need (Validation) and the spec (Verification).

Work top-down: read (S)pec before implementing; consult (D)esign before writing code; verify in (E)val before claiming Done. Sections below are organized by phase.

---

## (S)pecification — Specify + Business Context

The baseline feature spec lives at `openspec/specs/spec.md` — an index into 9 MECE category files (~80 features, ~85KB). It is the source of truth for what Gitea v1.22.x does today.

### When Starting Work

1. Identify the category your change touches (e.g. CI/CD & Automation, Collaboration).
2. Read the feature entry to find UI routes, API endpoints, `app.ini` config keys, and constraints.
3. Note surrounding features for cross-cutting impact.

### When Proposing a Change

Use OpenSpec delta specs (ADDED / MODIFIED / REMOVED markers). Reference the baseline: "Change feature X in category Y from behavior A to behavior B." Examples live in `openspec/changes/`.

### Acceptance Criteria

Define BDD-style acceptance scenarios for any non-trivial change. They become the contract for (E)val.

---

## (D)esign — Architecture + Planning

### Cross-cutting Conventions

@design.md

Architecture & layering, naming, error handling, context propagation, logging, i18n, testing, security. Apply these to every change before writing code.

### Layered Structure

The codebase follows a strict layered architecture. Dependencies flow downward:

```
cmd/       → CLI entry points, app initialization
routers/   → HTTP handlers (routers/web/ for UI, routers/api/v1/ for REST API)
services/  → Business logic (the orchestration layer)
models/    → Data access, XORM models, DB migrations
modules/   → Shared libraries and utilities (no dependency on models/services)
```

**Key rule**: `modules/` MUST NOT import from `models/`, `services/`, or `routers/`. Known exceptions are listed in `design.md` §1 — do not extend them.

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

### Module-Level CLAUDE.md

Most directories ship their own `CLAUDE.md` layered on top of `design.md`. Read the one in the directory you're editing before writing code.

### Planning

For non-trivial changes, write an implementation plan before coding. Use plan mode or the Superpowers `writing-plans` skill.

---

## (I)mplementation — Development + Testing

### Build & Run

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

### Lint & Format

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

### Code Generation

```bash
make generate                 # All Go generate targets
make generate-swagger         # Regenerate Swagger spec
make svg                      # Regenerate SVG assets
```

### Conventions

- PR title format: `<scope>: <description>`
- Go tests use `models/unittest/` framework for DB-backed test setup
- API tests go in `routers/api/v1/*_test.go`
- Integration tests in `tests/integration/`
- Swagger annotations on API handlers — run `make generate-swagger` after API changes
- This is Gitea 1.22.x branch — backport-focused, avoid unnecessary refactors

### Local Rules Per Directory

Read the `CLAUDE.md` of the directory you're editing before writing code. Local entry points, common pitfalls, and patterns live there.

### TDD Discipline

Follow Red → Green → Refactor. The Superpowers `test-driven-development` skill enforces this.

---

## (E)valuation — Validation + Verification

**Validation** (Are you building the right thing?) — does the change satisfy the original need stated in the spec / proposal?
**Verification** (Are you building it right?) — does the change conform to the spec, contracts, and conventions?

### Tests

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

### Contract Verification

- API handlers carry Swagger annotations — run `make generate-swagger` after any API change and commit the regenerated spec.
- If you changed DB schema, add a migration in `models/migrations/` and update fixtures in `models/fixtures/`.

### Definition of Done

A change is not Done until all of the following hold:

- [ ] Relevant spec category in `openspec/specs/` was consulted
- [ ] `design.md` principles honored (layering, naming, error handling, context, etc.)
- [ ] New code covered by tests (unit / integration / e2e as appropriate)
- [ ] `make lint` passes
- [ ] `make generate-swagger` regenerated if API changed
- [ ] Migration + fixtures updated if DB schema changed
- [ ] Acceptance criteria from (S)pec verified

### Cross-stage Validation

Before claiming Done, sanity-check across phases:

- Does the code match the design contract? **(D × I)**
- Does the code match the spec? **(S × I)**
- Does the test result match the acceptance criteria? **(S × E)**

---

## SDIE Principles

These principles govern the methodology; the SDIE framework document is the source of truth.

- **MECE** — Tasks atomic and non-overlapping; spec exhausts boundary cases. Prevents AI hallucination on ambiguous inputs.
- **Token economics** — Token is cheap; spend it on multi-angle verification, exhaustive analysis, and parallel exploration. Rework is expensive.
- **Harness Engineering** — Hooks, plan mode, and permission modes constrain the agent at each phase. The Steering Loop: when a failure pattern recurs, improve the harness rather than retrying.
- **Cross-stage validation** — Quality lives at the boundaries between phases, not inside any single one (see Cross-stage Validation above).
- **Real-time discarding** — When an approach fails 2-3 times in a row, abort and restart with a smaller scope. Context pollution compounds failures.
