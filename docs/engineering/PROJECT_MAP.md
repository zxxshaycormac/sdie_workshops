# Project Map

This map is an orientation aid, not a substitute for tracing the requested
behavior. Paths and flows below were verified against this workspace when the
engineering harness was established.

## System Shape

Gitea is a Go 1.22 monolith with server-rendered Go templates and a JavaScript
frontend that includes Vue 3 components. XORM supports SQLite, MySQL,
PostgreSQL, and MSSQL. The repository also includes package registries, Git
protocol handling, Actions, webhooks, queues, cron work, and multiple test
layers.

Approximate size at harness creation:

| Surface | Files |
| --- | ---: |
| Go | 2,486 |
| Go tests | 721 |
| Routers | 381 |
| Services | 405 |
| Models | 573 |
| Modules | 861 |
| JavaScript and Vue | 155 |
| Templates | 486 |

## Startup And Route Assembly

```text
main.go
  -> cmd.NewMainApp / cmd.RunMainApp
  -> cmd.CmdWeb / runWeb
  -> routers.InitWebInstalled
       -> settings, storage, cache, database, models, services
       -> indexers, queues, webhooks, Actions, SSH, cron
  -> routers.NormalRoutes
       -> /             routers/web
       -> /api/v1       routers/api/v1
       -> /api/internal routers/private
       -> /api/packages and /v2 when packages are enabled
       -> /api/actions* when Actions are enabled
```

Primary control points:

- `main.go`: process entry and build metadata.
- `cmd/main.go`: command registration, global flags, and configuration setup.
- `cmd/web.go`: install versus normal startup and HTTP serving.
- `routers/init.go`: subsystem initialization order and top-level mounts.
- `routers/web/web.go`: browser route tree.
- `routers/api/v1/api.go`: public API v1 route tree and authorization groups.
- `routers/private/internal.go`: internal route tree.

Initialization order matters. A component used during startup must be ready
before dependents and must participate in graceful shutdown where applicable.

## Layer Responsibilities

| Layer | Responsibility | Typical evidence |
| --- | --- | --- |
| `routers/` | HTTP routing, middleware, binding, status and response handling | Route registration and handler tests |
| `services/context/`, `services/forms/` | Request context and web form contracts | Binding and validation tests |
| `services/` | Use-case orchestration, transactions across domains, notifications | Service tests and caller traces |
| `models/` | Persistence models, queries, transactions, migrations | Model tests and database variants |
| `modules/` | Reusable infrastructure, Git operations, settings, public structs | Package tests and integration callers |
| `templates/` | Server-rendered UI and Swagger template output | Handler data, snapshots, or E2E |
| `web_src/js/` | Browser behavior and Vue components | Vitest and E2E tests |
| `web_src/css/` | Source styles and themes | Style lint and visual inspection |
| `tests/` | Integration and Playwright-based end-to-end behavior | Database-specific Make targets |

The layer boundaries are conventions, not a guarantee. Confirm the surrounding
package pattern before introducing a new dependency direction.

## Common Change Traces

### Public API

1. Find the mount in `routers/api/v1/api.go`.
2. Read the concrete handler under `routers/api/v1/`.
3. Trace service calls and model or module queries.
4. Inspect input and output contracts under `modules/structs/`.
5. Inspect conversion code, commonly under `services/convert/`.
6. Update Swagger annotations and referenced types.
7. Add handler, service, or contract tests at the narrowest real boundary.
8. For generated Swagger changes, run the generator; do not hand-edit output.

Compatibility defaults are documented in `CONTRIBUTING.md`: preserve existing
fields and GitHub-compatible semantics unless the change has an explicit reason
to diverge.

### Browser Page

1. Find the route in `routers/web/web.go`.
2. Read its middleware and handler under `routers/web/`.
3. Inspect `services/forms/` for submitted data.
4. Trace services and models that build or mutate state.
5. Find the template under `templates/` and the data keys it consumes.
6. Find initialization in `web_src/js/index.js` and the related feature or Vue
   component.
7. Check the related CSS source, Vitest test, and E2E flow.

### Persistence

1. Locate the model package and its `db.RegisterModel` use.
2. Trace query construction, context use, and transaction ownership.
3. Check migrations under `models/migrations/` for schema changes.
4. Check SQLite, MySQL, PostgreSQL, and MSSQL assumptions.
5. Trace converters and API or template consumers before changing field shape.
6. Use model tests first, then relevant migration or integration suites.

### Background Work

Trace the full state machine rather than only the worker function:

```text
producer -> payload/storage -> queue or cron registration -> consumer
         -> retry/timeout/cancel -> final persistence -> notification
```

Relevant areas include `modules/queue/`, `services/task/`, `services/cron/`,
`services/actions/`, `services/webhook/`, `services/mirror/`, and their models.

### Git And Repository Operations

Distinguish database metadata from filesystem Git state. Repository workflows
can cross `services/repository/`, `models/repo/`, `modules/git/`,
`modules/gitrepo/`, hooks under `routers/private/`, and notifications. Check
failure cleanup and filesystem side effects explicitly.

## Contracts And Cross-Layer Synchronization

- Public JSON and option types: `modules/structs/`.
- API route and operation annotations: `routers/api/v1/`.
- Swagger type references: `routers/api/v1/swagger/` and
  `templates/swagger/v1_json.tmpl`.
- Web form contracts: `services/forms/`.
- Configuration defaults and loading: `modules/setting/` and `custom/conf` only
  as runtime input.
- Localized user text: `options/locale/`.
- Templates and browser initialization: `templates/` plus `web_src/`.

When one member of a contract set changes, search for every serializer,
converter, template key, API annotation, and test that consumes it.

## Generated Outputs

Do not edit generated outputs to make a test pass. Change the source and run the
repository generator.

| Generated output | Source or generator |
| --- | --- |
| `public/assets/` | `web_src/`, webpack, frontend Make targets |
| `modules/*/bindata.go` and `.hash` | templates/options/public source and build generators |
| `templates/swagger/v1_json.tmpl` | Swagger annotations and `make generate-swagger` |
| `*.pb.go` | Adjacent `.proto` and protobuf generation |
| `assets/*.json`, option bundles | `build/` and Make generation targets |

Inspect the Makefile target before running a generator because some checks use
Git diffs, which are unavailable in this workspace until Git metadata exists.

## Runtime-Only And Sensitive Paths

The following are not project source:

- `custom/conf/`: local server configuration and secrets.
- `data/`: SQLite data, repositories, sessions, artifacts, keys, and storage.
- `log/`: runtime logs.
- `gitea`: built server binary.
- `node_modules/`, `.venv/`, `vendor/`: installed dependencies.
- `.make_evidence/`, coverage files, and test binaries: build/test evidence.

Only enter these paths for a task explicitly scoped to local runtime diagnosis.
Never paste credentials or tokens into OpenSpec artifacts or completion reports.

## Test Locations

- Go unit/package tests: adjacent `*_test.go` files.
- Frontend unit tests: `web_src/**/*.test.js` and related Vitest patterns.
- Integration tests: `tests/integration/`.
- End-to-end tests: `tests/e2e/*.test.e2e.js` driven by Playwright.
- Code-format tests and generators: `build/`.

Use `docs/engineering/VERIFICATION.md` to choose commands before editing.

## Git Metadata

This workspace normally has Git metadata restored. Use Git status, diffs, and
history when they are available, but do not assume they exist in every exported
or unpacked copy of the project. If `.git` is missing, state that limitation
before relying on blame, upstream delta, clean-tree checks, or Git-based
generator verification.
