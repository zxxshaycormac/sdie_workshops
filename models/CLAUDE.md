# `models/` — Design Conventions

`models/` is the data access layer: XORM models, DB sessions, migrations, and the unittest framework. No business logic. This file goes deeper than `design.md` on what it is like to work inside `models/`; the cross-cutting rules (architecture, naming, errors, context, logging, i18n, testing, security) remain authoritative and are not repeated here.

---

## 1. XORM Model Definition

**Rule**: Every persisted struct lives in a `models/<domain>/` package, has one `ID int64 \`xorm:"pk autoincr"\`` (or documented composite key), uses `xorm:` tags for column constraints (`NOT NULL`, `DEFAULT`, `INDEX`, `UNIQUE`, `created`/`updated`), and registers itself at package init via `db.RegisterModel(new(Type))`. Time fields use `timeutil.TimeStamp` with `created`/`updated` tags so XORM maintains them automatically.
**Why**: `RegisterModel` adds the bean to the slice that `SyncAllTables` walks to keep the schema in sync; forgetting it means the table is silently never created. The `xorm:` tags are the source of truth for column shape — there is no separate schema file in production. `timeutil.TimeStamp` keeps wall-clock values as `int64` seconds across every database backend.
**Frequency**: universal. ~110 call sites across ~96 files of `db.RegisterModel` in `models/`.
**Exceptions**: Migrations define their own beans (`models/migrations/v1_*/v*.go`) for legacy column shapes and intentionally do NOT call `RegisterModel` — they are one-shot scripts that run against a fresh engine. `models/system` and a few `TableName()` overrides (`models/auth/oauth2.go`, `models/organization/org.go`) rename the table; the struct still follows the registration rule.

---

## 2. Engine Acquisition via Context

**Rule**: Inside `models/`, acquire the XORM engine for the current request by calling `db.GetEngine(ctx)` at the top of the function. Never reference the package-global `db.x` from outside `models/db/`; never store an engine or session in a package-level variable or struct field that outlives the request.
**Why**: The engine returned by `GetEngine(ctx)` is the caller's transaction if the context carries one (`db.Context` implements `Engined`), or a context-bound view of the default engine otherwise. Bypassing it breaks transaction propagation and cancellation, and storing it leaks a connection across requests.
**Frequency**: universal. Thousands of `db.GetEngine(ctx)` call sites; the global `x` is read directly only inside `models/db/*.go` itself.
**Exceptions**: `models/db/engine.go` declares `var x *xorm.Engine` and uses it internally; `models/unittest/testdb.go` uses `context.Background()` once at bootstrap to seed `SetDefaultEngine`. Background goroutines that lose the request context fall back to `db.DefaultContext` (see `models/org_team.go:438`), with a `// FIXME` acknowledging the loss — match that pattern, do not invent a new one.

---

## 3. Transaction Boundary

**Rule**: Transactions are opened at the orchestration layer. The preferred pattern is `db.WithTx(ctx, func(ctx context.Context) error { ... })`, which auto-commits on nil and rolls back on error; legacy code uses `ctx, committer, err := db.TxContext(ctx); defer committer.Close()` and must call `committer.Commit()` before returning on the success path. `models/` functions that need a transaction accept `ctx` and let the caller decide whether to wrap; they do NOT begin their own transactions for single-statement reads.
**Why**: `WithTx`/`TxContext` are transaction-aware — if the parent context already holds a session, they reuse it as a `halfCommitter` so nested calls do not produce nested transactions. Beginning a transaction inside a single-row read doubles the round-trip and obscures the actual call shape.
**Frequency**: common. `db.WithTx` appears across `models/org_team.go`, `models/repo_transfer.go`, `models/project/board.go`, `models/user/badge.go`, etc.; `TxContext` is the older idiom in `models/org.go`, `models/repo.go`, `models/webhook/webhook.go`.
**Exceptions**: Migrations call `sess.Begin()` / `sess.Commit()` directly on `*xorm.Engine` (see `models/migrations/v1_9/v85.go`) because they run before `models/db` is fully wired; do not copy that style outside `models/migrations/`.

---

## 4. Migrations

**Rule**: A migration is a function `func(x *xorm.Engine) error` in `models/migrations/v1_<release>/v<N>.go`, registered at the bottom of the `migrations` slice in `models/migrations/migrations.go` via `NewMigration("<desc>", v1_xx.Fn)`. Migrations are append-only, sequential, immutable once released, and forward-only — there is no `Down`. The integer `N` (currently in the v280–v298 range) is the next free number; `minDBVersion` (70, Gitea 1.5.3) is the floor below which `Migrate` refuses to run.
**Why**: The runner in `Migrate(x)` slices `migrations[v-minDBVersion:]` and writes the new `Version` row after each step; renumbering a released migration breaks every upgraded database. Downgrading is intentionally fatal (`log.Fatal("Migration Error: ...")` in `migrations.go:675`) to protect user data.
**Frequency**: universal. 18 version directories `v1_6` through `v1_23`; the slice is a flat history with `// vX -> vY` comments separating releases.
**Exceptions**: `noopMigration` is used three times (`models/migrations/migrations.go:368`, `:405`, `:572`) to retire a botched migration while keeping the version number stable — copy this only if you must neutralise an already-released migration. Destructive operations DO exist inside migrations (`base.DropTableColumns`, `x.DropTables`) — destructive is allowed inside a migration, never in a `models/` runtime function.

---

## 5. Unittest Framework

**Rule**: Every test that touches the DB lives under `models/` or `tests/integration/` and goes through `models/unittest`. The package's `TestMain` calls `unittest.MainTest(m)`, which initialises a SQLite DB, syncs all registered tables, and loads fixtures. Each `TestXxx` calls `unittest.PrepareTestDatabase()` as its first line to reset fixtures to a known state.
**Why**: `unittest.MainTest` reads YAML fixtures from `models/fixtures/` (one file per table — `user.yml`, `repository.yml`, `access.yml`, 72 files total) and `PrepareTestDatabase` reloads them between tests so they are hermetic. Bypassing the framework means writing fixtures by hand and breaking parallelism.
**Frequency**: universal across DB-backed tests; `unittest.PrepareTestDatabase` is the established entry point.
**Exceptions**: Pure-logic tests with no DB dependency (a parser, a validator) skip `unittest` entirely — see `design.md` Section 7. Tests in `modules/` that DO need the DB are documented as isolation violations in `modules/CLAUDE.md` Section 1; do not add new ones.

---

## 6. Boundary — Data Integrity, Not Orchestration

**Rule**: Functions in `models/` perform CRUD, schema-touching queries, and the data-integrity checks that the DB itself cannot enforce (uniqueness preflight, reserved-name checks, foreign-key bookkeeping). They do NOT make policy decisions, send notifications, render templates, or call into `services/` or `routers/`. HTTP handling, templating, and business orchestration belong to upper layers.
**Why**: Keeping `models/` narrow lets it be tested with fixtures alone and reused by every caller — the CLI, the web router, the API router, the webhook receiver — without dragging in HTTP-shaped dependencies. A `models/` function that imports `routers/` creates a cycle; one that imports `services/` inverts the architecture.
**Frequency**: universal as a rule.
**Exceptions**: Older packages (`models/issues`, `models/repo`, `models/org_team.go`) carry orchestration-shaped helpers (e.g. `models.AddRepository` (in `models/org_team.go`) writes the team-repo row AND watches the repo AND emits a team add event) because they predate `services/`. Match the surrounding file when extending one of these; place new orchestration in `services/`. `models/repo/repo.go` imports `html/template` solely to use `template.HTML` as a column type for safe HTML payloads — it does not render.

---

## 7. Additional Conventions

**Rule**: Sentinel errors in `models/` use typed structs (`ErrXxxNotExist`, `ErrXxxAlreadyExist`) with an `Unwrap()` method that returns one of the shared sentinels in `modules/util` (`util.ErrNotExist`, `util.ErrPermissionDenied`, etc.). Provide an `IsErrXxx(err error) bool` companion. Prefer the constructors in `modules/util/error.go` (`util.NewNotExistErrorf`, `util.NewAlreadyExistErrorf`, `util.NewPermissionDeniedErrorf`, `util.NewInvalidArgumentErrorf`) over raw `errors.New`.
**Why**: Typed errors carry context (the missing ID, the conflicting name) at the data layer while still unwrapping to a sentinel the router can map to an HTTP status with `errors.Is`.
**Frequency**: common; the dominant pattern in `models/user`, `models/repo`, `models/issues`, `models/organization`.
**Exceptions**: A handful of older errors (`models/git/protected_branch.go: ErrBranchIsProtected = errors.New(...)`) predate the convention. Match local style in those files.

**Rule**: Build queries with the XORM session DSL (`Where`, `And`, `In`, `ID`, `OrderBy`, `Join`, `Count`, `Find`, `Get`, `Iterate`) or the `xorm.io/builder` package for composable conditions. Never interpolate user-controlled values into a SQL string.
**Why**: The DSL parameterises; concatenation is the SQL injection vector (see `design.md` Section 8).
**Frequency**: universal for new code; `builder.Eq`, `builder.In`, `builder.Like` appear across every `models/` subpackage.
**Exceptions**: `sess.SQL("...", constants)` exists in legacy code and is acceptable ONLY when every interpolated value is a constant or internal ID (the dominant pattern in `models/migrations/`). Audit any new `SQL(...)` call.

**Rule**: A `models/<domain>/` package may import other `models/` packages through the aliased `xxx_model` form (`user_model`, `repo_model`) but must avoid cycles. Cross-domain reads prefer the typed accessor (`user_model.GetUserByID`) over a hand-rolled join.
**Why**: Reuse keeps consistency invariants in one place; a cycle would break the build because Go forbids it.
**Frequency**: universal.
**Exceptions**: none — when two packages need each other, factor the shared type into a third lower package (as `models/db` already does for engine primitives).
