# `modules/` — Design Conventions

`modules/` contains shared libraries and utilities used across all layers. It MUST NOT depend on `models/`, `services/`, or `routers/` (with a small set of documented exceptions). This file goes deeper than `design.md` Section 1 on what it is like to work inside `modules/`; the cross-cutting rules (architecture, naming, errors, context, logging, i18n, testing, security) remain authoritative and are not repeated here.

---

## 1. Isolation and the Known Violations

**Rule**: New code in `modules/` imports nothing from `models/`, `services/`, or `routers/`. Existing files that do are documented as accepted technical debt and are NOT precedents.
**Why**: `modules/` is imported by everything else. An upward import anywhere creates a cycle risk for every transitive consumer. The rule is the default; violations are warnings.
**Frequency**: universal as a rule; ~60 files violate it (mostly concentrated in a few clusters).
**Exceptions**: Each violating cluster below has a structural reason. When extending one of these files, match local style; when starting fresh, treat the cluster as a smell, not a pattern.

- **`modules/repository/`** (`env.go`, `commits.go`, `create.go`, `branch.go`, `delete.go`, `fork.go`, `repo.go`, `collaborator.go`, `init.go`) imports `models/repo` and `models/user`. **Why**: these files build the `GITEA_*` environment variables git hooks consume, and the hook contract requires concrete repo/user IDs. The cluster is effectively a service-layer helper misfiled under `modules/`; prefer `services/repository/` for new git-hook glue.
- **`modules/metrics/collector.go`** imports `models/activities`, `models/db`. **Why**: the Prometheus `Collector` interface needs raw row counts and the collector runs server-side after boot; refactoring would mean duplicating count queries. Acceptable but isolated to one file.
- **`modules/eventsource/manager_run.go`** imports `models/activities`, `models/issues`, `services/convert`. **Why**: the live notification poller fans DB rows out to SSE clients; the conversion to wire types belongs in `services/convert`. Long-standing debt.
- **`modules/auth/webauthn/webauthn.go`** imports `models/auth`, `models/db`, `models/user`. **Why**: it implements the third-party `webauthn.User` interface, whose methods must return concrete user/credential structs. The dependency is forced by the upstream library shape.
- **`modules/system/db.go`** imports `models/system`. **Why**: `DBStore` persists app-state items via the model layer; the file is the storage adapter for an in-process key/value API.
- **`modules/private/serv.go`** imports `models/asymkey`, `models/perm`, `models/user`. **Why**: `modules/private/` is a client for Gitea's own internal HTTP API (see `newInternalRequest` in `internal.go`); the model types appear only as wire-format response structs. Most of `private/` is clean — only the SSH-key path leaks model types.
- **`modules/actions/task_state.go`, `modules/actions/log.go`** import `models/actions`. **Why**: pure render helpers that walk `actions_model.ActionTask`/`ActionTaskStep`. Could move to `services/actions` or `templates/` helpers.
- **`modules/templates/util_render.go`, `util_misc.go`, `util_avatar.go`, `helper.go`** import many `models/*` packages. `helper.go` additionally imports `services/gitdiff` and `services/webtheme` — the only place in `modules/` besides `eventsource/manager_run.go` that crosses to `services/`. **Why**: template helpers are by nature presentation glue between models and HTML; this is the single largest structural concession. New template helpers that touch models follow this cluster; pure-string helpers stay in `modules/templates/util_*.go` without model imports.
- **`modules/gitgraph/graph_models.go`** imports `models/asymkey`, `models/db`, `models/git`, `models/repo`, `models/user`. **Why**: commit-graph rendering joins commit SHAs to GPG signatures and repo metadata. The model-touching logic is segregated in one `*_models.go` file so the rest of the package stays clean.
- **`modules/ssh/ssh.go`** imports `models/asymkey`. **Why**: the built-in SSH server authenticates against the public-key table. The dependency is contained to the SSH auth path.
- **`modules/session/db.go`, `modules/badge/badge.go`, `modules/activitypub/user_settings.go`, `modules/indexer/issues/db/db.go`, `modules/indexer/issues/internal/model.go`, `modules/indexer/issues/dboptions.go`, `modules/indexer/stats/db.go`, `modules/indexer/stats/indexer.go`, `modules/indexer/code/*.go`** — additional files that touch `models/db` for XORM sessions or read model types for indexing. These are persistence/indexing adapters and follow the same "match local style, don't spread the pattern" rule.
- **Tests only**: `modules/markup/html_test.go`, `modules/markup/markdown/markdown_test.go`, and similar tests in `modules/repository/`, `modules/system/`, `modules/templates/`, `modules/activitypub/`, `modules/indexer/` import `models/unittest` (see Section 4 for the full list). Tests are exempt from the production-import rule.

---

## 2. Sub-package Organization

**Rule**: Each `modules/<name>/` package owns one responsibility. Avoid a catch-all `util/`; if a new utility does not fit an existing sub-package, prefer a new narrow package.
**Why**: A focused package is easy to grep, easy to test in isolation, and easy to import without pulling unrelated symbols into the call site. `modules/util/` is the historical exception — it is split across many small files (`error.go`, `path.go`, `string.go`, `paginate.go`, `keypair.go`, `sanitize.go`, etc.) so each concern remains independently addressable; new utilities should follow this file-per-concern split rather than dumping into `util.go`.
**Frequency**: common.
**Exceptions**: `modules/util/` is the only sanctioned multi-concern package; do not create another. Sub-packages of `modules/indexer/` (e.g. `issues/`, `code/`, `stats/`) use an `internal/` child to hold the public `Indexer` interface plus a backend-specific sub-folder per storage (`bleve/`, `elasticsearch/`, `meilisearch/`, `db/`) — copy this layout when adding a new indexer type.

---

## 3. Public API Design

**Rule**: Hide internal types behind small exported interfaces and helper constructors. Prefer returning `*ConcreteType` from a `NewXxx` constructor over exposing package-level state; when package-level state is unavoidable (e.g. `setting`, `log`), expose only functions, never the variables.
**Why**: `modules/` is the lowest layer; its API is the contract every upper layer programs against. Stable, narrow exports keep refactorings local.
**Frequency**: universal.
**Exceptions**: `modules/log`, `modules/setting`, and `modules/cache` intentionally expose stateful singletons (`log.Trace`, `setting.AppURL`, `cache.GetCache`) because they model process-wide configuration; this is the established pattern for those three packages and should not spread elsewhere.

**Rule**: Internal helpers that build HTTP requests or perform transport work take `ctx context.Context` and keep their signatures unexported (`newInternalRequest`, `requestJSONResp`). Exported wrappers (`ServNoCommand`, `HookPreReceive`, `GenerateActionsRunnerToken`) take `ctx` first and hide the wire-format details from callers.
**Why**: Call sites in `cmd/` and `routers/` should not need to know about internal tokens, URLs, or response-shape structs.
**Frequency**: common in `modules/private/` and `modules/httplib/`.
**Exceptions**: none.

---

## 4. Tests

**Rule**: Tests under `modules/` run with no database. A test that needs a DB belongs in `models/` or `tests/integration/`, not here. Pure-logic tests live next to the file (`foo_test.go`), use `github.com/stretchr/testify/assert`, and may import `modules/setting` or `modules/test` for fixture-style helpers but never `models/unittest`.
**Why**: `modules/` is the dependency-free base layer; reaching for the DB unittest harness here would mean the layer is no longer runnable in isolation, and would create a build-time cycle (`models/unittest` transitively imports `models/...`).
**Frequency**: universal as a rule.
**Exceptions**: The 15 files under `modules/` that import `models/unittest` are tests of code that itself violates the isolation rule (`modules/repository/*_test.go`, `modules/markup/html_test.go`, `modules/system/appstate_test.go`, `modules/templates/util_render_test.go`, `modules/activitypub/*_test.go`, `modules/indexer/**/*_test.go`). These tests inherit the violation from the package under test; do not add a new `unittest` import to test code that does not already depend on the DB.

---

## 5. Additional Conventions

**Rule**: Keep `context.Context` as the first parameter of exported functions that do I/O, but do NOT thread it through pure helpers (formatters, validators, parsers, math). Never store a request-scoped `context.Context` in a struct that outlives the request.
**Why**: Pure helpers stay pure so they can be reused without a context; storing contexts in long-lived structs leaks cancellation deadlines into background state.
**Frequency**: universal.
**Exceptions**: structs whose explicit job is to carry a request scope (`AvatarUtils{ctx context.Context}` in `modules/templates/util_avatar.go`) hold the context for the duration of one render pass; this is a scoped adapter, not a leaked context.

**Rule**: Reusable building blocks (queues, process manager, graceful shutdown, LRU caches, optional values) live in dedicated packages and expose a stable interface that upper layers compose — never inline the mechanism.
**Why**: `modules/queue`, `modules/process`, `modules/graceful`, `modules/regexplru`, `modules/optional` are depended on by long-running services; their interfaces must change slowly.
**Frequency**: universal.
**Exceptions**: none.

**Rule**: When a package needs a backend abstraction, define the interface in an `internal/` child package and put each concrete backend in its own sibling package, with a top-level package that selects and wires them.
**Why**: This is the `modules/indexer/{code,issues,stats}/` shape: `internal/` holds the `Indexer` interface and shared data types, `bleve/`/`elasticsearch/`/`meilisearch/`/`db/` are swappable backends, and the parent package picks one based on `setting.*`. It keeps backend imports off the call sites and lets `go build` prune unused drivers.
**Frequency**: common for new indexers and storage adapters.
**Exceptions**: `modules/storage/` uses a similar but flatter shape (one package per backend, no `internal/`); follow whichever the surrounding package already uses.
