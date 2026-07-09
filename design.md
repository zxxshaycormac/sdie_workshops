# Gitea Design Conventions

## 0. Document Positioning

This document is written for **AI assistants** (Claude, other LLM agents) that modify Gitea code. Its purpose is to capture the **cross-cutting conventions** that the codebase expects every contributor to follow, so an assistant reading this can produce changes that fit alongside existing code without needing to reverse-engineer the rules from thousands of files.

**How to use this document.** Read it once for the principles. Before editing code in a specific directory, layer on that directory's `CLAUDE.md` for directory-specific rules (entry points, common pitfalls, local style). Treat the rules below as defaults; a directory `CLAUDE.md` may add stricter rules but should not contradict what is stated here. If a contradiction appears, trust the directory file for that directory and surface the discrepancy.

**Index of directory `CLAUDE.md` files.** Each of the following provides rules specific to its slice of the codebase; this document is the umbrella they inherit from:

- `modules/CLAUDE.md` — shared libraries; the strictest isolation rules apply here.
- `models/CLAUDE.md` — XORM models, DB sessions, migrations, fixtures.
- `services/CLAUDE.md` — business-logic orchestration layer.
- `cmd/CLAUDE.md` — CLI entry points and application bootstrap.
- `routers/web/CLAUDE.md` — browser-facing HTTP handlers (HTML responses).
- `routers/api/v1/CLAUDE.md` — REST API handlers (JSON responses).
- `web_src/CLAUDE.md` — frontend source (JS, CSS, Vue, Fomantic UI).
- `templates/CLAUDE.md` — Go HTML templates (`.tmpl`).

---

## 1. Architecture Principles

The codebase follows a strict layered dependency graph:

```
cmd → routers → services → models → modules
```

Dependencies flow downward only. Within a layer, packages may depend on sibling packages but must avoid cycles.

**Rule**: `modules/` MUST NOT import from `models/`, `services/`, or `routers/`. It is the foundational shared library and must remain dependency-free of upper layers.
**Why**: Lower layers cannot depend on upper layers without creating cycles; `modules/` is imported by everything, so any upward dependency would taint the entire codebase.
**Frequency**: universal (intended); violated in a small number of known files.
**Exceptions**: The following `modules/` files do import from `models/` or `services/` and are accepted technical debt — do not copy these patterns in new code:
- `modules/metrics/collector.go` — Prometheus collector reads row counts directly from `models/activities` and `models/db`.
- `modules/eventsource/manager_run.go` — internal event-source manager uses `models/activities`, `models/issues`, and `services/convert`.
- `modules/repository/*.go` — environment/push helpers (`env.go`, `commits.go`, `create.go`, `branch.go`, `delete.go`, `fork.go`, `repo.go`) depend on `models/repo` and `models/user` because they build git-hook environment variables.
- `modules/auth/webauthn/webauthn.go` — uses `models/auth`, `models/db`, `models/user` to implement the webauthn.User interface.
- `modules/system/db.go` — implements app-state persistence via `models/system`.
- `modules/private/serv.go` — internal SSH serv helper reads `models/asymkey`, `models/perm`, `models/user`.
- `modules/actions/task_state.go`, `modules/actions/log.go` — render `models/actions` types.
- `modules/templates/util_render.go`, `modules/templates/util_misc.go` — template helpers touch `models` for rendering.
- `modules/gitgraph/graph_models.go`, `modules/ssh/ssh.go`, `modules/markup/html_test.go` (test only).

When adding code to `modules/`, treat the absence of upward imports as the rule and the files above as warnings, not precedents.

**Rule**: Services orchestrate; models do CRUD only; routers bind, validate, and dispatch.
**Why**: Keeping layers narrow in responsibility prevents business logic from leaking into HTTP handlers and prevents HTTP-shaped concerns from leaking into data access.
**Frequency**: common (uniform in new code; legacy code occasionally blurs the line, especially older `models/` packages that contain some business logic).
**Exceptions**: Some `models/` packages (notably `models/issues`, `models/repo`) contain logic that arguably belongs in services. Match the surrounding file when extending existing code, but prefer placing new orchestration in `services/`.

**Rule**: No circular imports within a layer.
**Why**: Go forbids import cycles; the package layout must form a DAG.
**Frequency**: universal.
**Exceptions**: none observed.

---

## 2. Naming Conventions

**Rule**: Package names are lowercase, single words, no underscores or hyphens.
**Why**: Go convention; matches the import-path final segment and keeps call sites readable (`db.Get`, `repo_model.GetByID`).
**Frequency**: universal.
**Exceptions**: When a package name collides with a stdlib or common name, Gitea uses a `xxx_model` suffix (e.g. `repo_model`, `user_model`, `asymkey_model`, `activities_model`, `issues_model`). These appear most often in `models/` subpackages that are imported from sites that also use the bare `repo`/`user` words.

**Rule**: Exported types, functions, and constants use PascalCase (`Repository`, `CreateRepository`, `DefaultContext`). Unexported identifiers use camelCase (`prepareRepoCommit`).
**Why**: Go visibility rules; the case of the first letter encodes export status.
**Frequency**: universal.
**Exceptions**: none.

**Rule**: Method receivers use short abbreviations of the type name, typically one to three lowercase letters.
**Why**: Brevity at call sites; consistent idiom across the file.
**Frequency**: universal.
**Exceptions**: A small number of files use longer receiver names (`upload *Upload`, `cfg *UnitConfig`). Follow the existing receiver in the file you are editing; for a new type pick a short abbreviation.

**Rule**: File names use `snake_case.go`. Test files use `xxx_test.go` next to the file they test.
**Why**: Matches Go ecosystem tooling expectations (test discovery, build constraints).
**Frequency**: universal.
**Exceptions**: none observed.

**Rule**: Database columns are named in `snake_case`, declared via XORM struct tags (`xorm:"pk autoincr"`, `xorm:"UNIQUE(s) INDEX NOT NULL"`, `xorm:"-"` for non-persisted fields).
**Why**: XORM auto-maps field names to columns; explicit tags disambiguate cases the auto-mapper would get wrong and document intent.
**Frequency**: universal across `models/` structs.
**Exceptions**: none.

---

## 3. Error Handling

**Rule**: Wrap errors once per architectural layer transition (e.g., when returning from `models/` to `services/`, or `services/` to `routers/`) with `fmt.Errorf("<description>: %w", err)`. The wrapping description should explain what that layer was doing. Do not re-wrap within the same layer unless adding new context.
**Why**: Wrapping preserves the error chain so callers can `errors.Is` / `errors.As` while still producing a readable stack-like message at the boundary.
**Frequency**: universal in new code; pervasive in `services/` and `routers/`.
**Exceptions**: none. The `%v` form is discouraged for cross-layer wrapping because it breaks `errors.Is`.

**Rule**: Define sentinel errors as package-level `var ErrXxx` variables. Use the constructors in `modules/util/error.go` (`util.NewNotExistErrorf`, `util.NewAlreadyExistErrorf`, `util.NewPermissionDeniedErrorf`, `util.NewInvalidArgumentErrorf`) so the error unwraps to one of the four shared sentinels (`util.ErrNotExist`, `util.ErrAlreadyExist`, `util.ErrPermissionDenied`, `util.ErrInvalidArgument`).
**Why**: Callers across layers test for these sentinels with `errors.Is(err, util.ErrNotExist)` to choose HTTP status codes or branch logic; the constructor wraps the sentinel in a `SilentWrap` so the user-facing message can differ from the classification.
**Frequency**: common in new code; the established pattern in `models/`.
**Exceptions**: A handful of older model errors still use raw `errors.New` (e.g. `models/git/protected_branch.go: ErrBranchIsProtected = errors.New("branch is protected")`) and a few use `db.ErrNotExist{Resource: "..."}` directly (`models/git/lfs.go`). When adding new sentinel errors prefer the `util.New*Errorf` constructors; when extending code that already uses raw `errors.New`, match local style.

**Rule**: Return errors up to the caller. Log them only at the boundary (router handler, CLI command, long-running service loop).
**Why**: Logging deep in the call stack produces duplicate log lines for the same error and hides the failure from callers that may want to react to it.
**Frequency**: common; boundary logging is the established pattern (`routers/web/repo/render.go: log.Error(...)`, `services/webhook/notifier.go: log.Error(...)` inside a service loop).
**Exceptions**: Long-running loops (`services/webhook/deliver.go`, queue consumers) log errors because they have no caller to return to.

**Rule**: Never swallow an error with `_ =`. If a function returns an error you genuinely cannot act on, document why in a comment.
**Why**: Silent error drops hide bugs; the next reader cannot tell whether the drop was intentional.
**Frequency**: universal as a rule; violated in practice.
**Exceptions**: The dominant acceptable pattern is ignoring the return of I/O primitives whose errors are known to be irrelevant in context (e.g. `_, _ = sb.Write(lineBytes)` on a `*strings.Builder`, `_ = rd.UnreadByte()`, `_ = reader.Close()` in a defer). These appear throughout `services/gitdiff/gitdiff.go` and are idiomatic. Truly inappropriate drops — where the error would matter — do exist (`services/uinotification/notify.go: _ = ns.issueQueue.Push(opts)`) and are technical debt, not precedent.

---

## 4. Context Propagation

**Rule**: `ctx context.Context` is the first parameter of every function that performs I/O (DB, HTTP, git, filesystem, queue). The same rule applies to methods — `ctx` precedes the receiver-named parameters.
**Why**: Cancellation, deadlines, tracing, and the XORM session (`xorm.Session.Context(ctx)`) all key off the propagated context; threading it through every I/O call lets upper layers cancel work cleanly.
**Frequency**: universal. ~2,260 functions across `models/` and `services/` accept `ctx context.Context` as a leading parameter; `modules/` has ~324 more.
**Exceptions**: Pure functions (formatters, parsers, validators) take no context. Constructors that only allocate a struct take no context.

**Rule**: For non-request flows (background workers, queue handlers, CLI subcommands without an explicit request scope) use `db.DefaultContext` as the root context.
**Why**: Background work still needs a context that owns a DB session and respects process-level cancellation; `db.DefaultContext` is the project-wide root.
**Frequency**: common in `cmd/`, queue consumers, and tests.
**Exceptions**: Do not use `db.DefaultContext` inside an HTTP handler — use the request's context. Do not store the request context in a struct that outlives the request.

**Rule**: Timeouts and cancellation originate at the boundary (router middleware, CLI command), not inside `models/` or `services/`.
**Why**: The boundary knows the operational budget; lower layers should honour whatever context they are given.
**Frequency**: common.
**Exceptions**: Some long-running git operations wrap the incoming context with `process.GetManager().AddTypedContext(...)` to register with the process manager. This is a Gitea-specific extension of the boundary rule, not a violation.

---

## 5. Logging

The logging API lives in `modules/log` and exposes package-level functions: `Trace`, `Debug`, `Info`, `Warn`, `Error`, `ErrorWithSkip`, `Critical`, `Fatal`. All take a printf-style format string and variadic args.

**Rule**: Pick the level by what the operator needs to do with the line.
- `Trace` — extremely chatty per-item progress (e.g. `log.Trace("Task[%d] has already been delivered", task.ID)`).
- `Debug` — one-line diagnostics useful when reproducing a problem (e.g. `log.Debug("ParsePatch(%d, %d, %d, ..., %s)", ...)`).
- `Info` — noteworthy normal operation the operator may want to see (e.g. `log.Info("Branch %q doesn't match branch filter %q, skipping", branch, w.BranchFilter)`).
- `Warn` — recoverable, unexpected, or degraded behaviour (e.g. `log.Warn("GetHookTaskByID[%d] warn: %v", taskID, err)`).
- `Error` — operation failed; surfaced in logs but execution continues (e.g. `log.Error("Unable to deliver webhook task[%d]: %v", task.ID, err)`).
- `Critical` / `Fatal` — only at process-init boundaries; `Fatal` exits.
**Why**: Operators filter logs by level; misclassification either floods production or hides real failures.
**Frequency**: universal.
**Exceptions**: none.

Note: the `modules/log` package API is printf-style only; no structured key/value field helpers (`WithField`, `WithValues`) exist at the package level, so format fields inline.

**Rule**: Never log PII, passwords, tokens, session IDs, or full request bodies. Log opaque identifiers (user ID, repo ID, task ID) instead.
**Why**: Logs are commonly shipped to shared aggregation systems; secrets in logs are a security incident.
**Frequency**: universal as a rule.
**Exceptions**: none acceptable. If you find existing code that logs a secret, treat it as a bug to fix, not a pattern to copy.

---

## 6. Internationalization (i18n)

Locale files live in `options/locale/` as `locale_<lang>.ini` (40+ languages; `locale_en-US.ini` is the source of truth).

**Rule**: Every user-visible string must come from a locale key. In Go code use `ctx.Tr("key.name", args...)` (the `context.Context` of a handler exposes `Tr`). In templates use `{{ctx.Locale.Tr "key.name"}}` (or `{{.locale.Tr ...}}` where a locale is passed into a partial).
**Why**: Gitea ships in many languages; hardcoded English forces every translator to chase code changes and breaks non-English UIs.
**Frequency**: common; the established pattern in `routers/web/` and `routers/api/v1/` (API uses `Tr` for human-readable error strings returned alongside machine-readable codes).
**Exceptions**: Programmatic strings never shown to users (internal log lines, error wrapping descriptions, HTTP header values) are not localized — `fmt.Errorf("GetUserByID: %w", err)` stays English.

**Rule**: Add new keys to `options/locale/locale_en-US.ini` first. Other locale files are translated by the community; do not machine-translate them.
**Why**: English is the canonical source; the translation workflow diffs against it.
**Frequency**: universal.
**Exceptions**: none.

**Rule**: Do not concatenate translated fragments. Use placeholder substitution (`ctx.Tr("user.form.name_reserved", newName)`) so translators can reorder the sentence.
**Why**: Word order differs across languages; concatenation produces ungrammatical translations.
**Frequency**: universal.
**Exceptions**: none.

---

## 7. Testing Principles

**Rule**: Any test that touches the database uses the `models/unittest/` framework. The test entry point is `unittest.InitXORMWithFixture` / `unittest.LoadFixtures` (typically invoked from a `TestMain` or a shared `init`).
**Why**: Tests share one SQLite (or configured) DB; `unittest` resets fixtures between tests so they are hermetic and parallelisable.
**Frequency**: universal — 332 `*_test.go` files import `code.gitea.io/gitea/models/unittest`.
**Exceptions**: Pure-logic unit tests with no DB dependency do not need `unittest`.

**Rule**: Put pure-logic tests next to the file they test (`foo_test.go`). Put HTTP-shaped integration tests in `tests/integration/`. Put API-contract tests next to the handler (`routers/api/v1/<area>_test.go`).
**Why**: Mirrors Go ecosystem expectations and Gitea's `make test*` targets.
**Frequency**: universal.
**Exceptions**: none.

**Rule**: Do not commit `t.Skip()` without a tracking issue and a comment naming it. Conditional skips for environment capabilities (e.g. `t.Skip("Test skipped for non-Minio-storage.")`) are acceptable.
**Why**: Unconditional skips silently rot; conditional skips tied to a documented capability are the established pattern.
**Frequency**: common. Do not add bare `t.Skip()` for unrelated work.
**Exceptions**: Observed acceptable uses — conditional-environment skips in `tests/integration/` for LDAP, Minio storage, and external markup renderers (each tied to an absent optional capability, not to pending work).

---

## 8. Security General Principles

**Rule**: A trust boundary is the first function that receives untrusted data — typically the same place Section 4 calls "the boundary" (router handler, CLI command), plus additionally webhook receivers and migration importers. Validate input at the trust boundary. Lower layers may assume data from a *single, trusted* caller (e.g., a router that already validated) has been validated; when a function receives data from an external/trusted-less source (webhook, migration, queue), it becomes a trust boundary itself and must re-validate per Rule 3.
**Why**: Centralising validation prevents the same input from being checked inconsistently across call sites; it also keeps lower layers focused.
**Frequency**: common. Routers perform binding and validation via the request context before invoking services.
**Exceptions**: When a service is callable from multiple boundaries with different validation rules, the service re-validates defensively.

**Rule**: Authorisation goes through middleware and the request context, not ad-hoc checks inside handlers.
**Why**: A single `verifyAuthWithOptions(&VerifyOptions{SignInRequired: ...})` declaration on a route group covers every handler in that group; scattered checks get forgotten when routes are added.
**Frequency**: universal. See `routers/web/web.go` (`ignSignIn := verifyAuthWithOptions(...)`) and the `common.VerifyOptions` type.
**Exceptions**: Some handlers do additional per-object permission checks (e.g. collaborator mode on a specific repo) after the middleware passes — this is correct, the middleware handles the coarse check, the handler the fine one.

**Rule**: Never trust user input. Re-validate at every trust boundary (as defined in Rule 1 above) — e.g. when a service receives data from a webhook payload or a migration source, even though the source is "internal".
**Why**: Defence in depth; a misconfigured upstream or a forged webhook should not corrupt the DB.
**Frequency**: common.
**Exceptions**: none.

**Rule**: Build SQL through XORM session methods (`Where`, `And`, `In`, `ID`, `OrderBy`, etc.) or the `builder` package. Never concatenate user-controlled values into a SQL string.
**Why**: String concatenation is the canonical SQL injection vector; XORM and `builder` parameterise.
**Frequency**: universal for new code.
**Exceptions**: Raw `sess.SQL("...")` exists in `models/` (e.g. `models/org_team.go`, `models/repo.go`) and is heavily used in `models/migrations/`. These are acceptable only when the SQL string contains no user-controlled values (constants, internal IDs, schema introspection). When introducing a new `SQL(...)` call, audit every interpolated value and prefer XORM/builder when any value is dynamic.

**Rule**: Permission checks use the `perm` package (`models/perm`) with explicit `AccessMode` values; do not invent ad-hoc role checks.
**Why**: A single source of truth for what "read", "write", "admin" mean keeps permission decisions consistent across the codebase.
**Frequency**: universal.
**Exceptions**: none.

---

## 9. UI Behavior Conventions

This section governs cross-cutting UI behavior — the contracts every page should honor regardless of feature. It does not prescribe visual design (colors, spacing, typography); those belong in a future style guide. It covers only behaviors that are testable and observable.

**Rule**: Every button or link that triggers an async operation shall enter a disabled + spinner state for the duration of the request and shall not accept a second click until the response arrives.
**Why**: Prevents double-submit race conditions (duplicate comments, duplicate merges, duplicate deletes) and gives the user visible feedback that the click was received.
**Frequency**: universal for forms, action buttons, and async links.
**Exceptions**: Pure-navigation links need no spinner; button-as-link patterns that don't trigger a request can stay enabled.

**Rule**: Operations expected to take longer than 200ms shall display a loading indicator (skeleton, spinner, or `is-loading` class) at the location where the result will appear, not at a distant page chrome.
**Why**: Below 200ms the operation feels instant; above 200ms the user needs feedback. Locating the indicator at the result site lets the user keep reading surrounding content.
**Frequency**: common for list views, dashboard widgets, modal data fetches.
**Exceptions**: Pre-cached or pre-rendered data needs no loading state.

**Rule**: Errors from form submissions shall be displayed inline next to the offending field when the error is field-level, and at the top of the form when the error is form-level. Errors from page-level operations shall render a dedicated status template (`templates/status/*.tmpl`) rather than a partial page.
**Why**: Inline errors let the user fix the problem without losing context. Page-level errors via dedicated templates ensure consistent styling and prevent half-rendered UI.
**Frequency**: universal.
**Exceptions**: Background operations (webhook delivery, queue items, cron failures) surface errors via the admin notices panel (ADM-13), not the user UI.

**Rule**: Empty lists and first-use states shall render a guided empty state (the `templates/repo/empty.tmpl` pattern) explaining what is missing and what action the user can take — never a bare "no results".
**Why**: A bare "no results" leaves the user uncertain whether the system is broken, the data is missing, or the filter is wrong. A guided empty state explains the situation and offers a next step.
**Frequency**: common for new repositories, first-time visits, filtered lists.
**Exceptions**: Search results with no matches may use a simpler "no results for query" message since the user already knows what they searched for.

**Rule**: Transient success and non-blocking error feedback shall use the toast system (`web_src/js/modules/toast.js`) rather than flash redirects or `alert()` dialogs. Toasts shall auto-dismiss after a configurable timeout and shall stack without overlapping.
**Why**: Toasts let the user continue working without dismissing a modal. Auto-dismiss prevents toast graveyards. Stacking prevents overlap when multiple toasts fire in quick succession (e.g. multi-file upload).
**Frequency**: common for save success, copy-to-clipboard confirmation, async failure.
**Exceptions**: Blocking errors that require user acknowledgment (e.g. merge conflict, validation failure on submit) shall use a modal or inline error, not a toast.

**Rule**: Destructive or irreversible operations (delete repository, delete user, transfer repository, force-push to a protected branch) shall require confirmation via a modal (`web_src/js/features/comp/ConfirmModal.js`) and shall require the user to type the entity name when the operation is irreversible.
**Why**: Confirmation modals prevent accidental clicks on destructive actions. Typed confirmation for irreversible actions acts as an intent filter — the friction is the feature.
**Frequency**: common for delete/transfer operations.
**Exceptions**: Reversible destructive operations (delete a draft comment, unstar a repo, unpin an issue) can use a simple confirm modal without typed input.

**Rule**: Toggles that change per-user state (reactions, stars, watches, subscription mode) shall update the UI optimistically before the server confirms and shall revert on failure. Operations that mutate shared state (merge PR, close issue, change visibility, transfer ownership) shall wait for server confirmation before updating the UI.
**Why**: Per-user toggles are idempotent and almost always succeed; making the user wait feels slow. Shared-state mutations need server truth because other users see the result.
**Frequency**: common.
**Exceptions**: none.

**Rule**: Keyboard shortcuts shall use the established defaults (`/` for site search, `?` for shortcut help, `c` for compose/create in scoped contexts) and shall be discoverable via the help overlay triggered by `?`.
**Why**: Consistent defaults let power users transfer muscle memory across pages. Discoverability via `?` ensures new users find the shortcuts without documentation.
**Frequency**: common for high-traffic pages.
**Exceptions**: Single-purpose admin pages may omit shortcuts; pages with text inputs shall suppress single-key shortcuts while the input has focus.

**Rule**: Each semantically-distinct value of a displayed enum shall render with a distinct visual treatment — icon, color, or label — across every presentation surface that displays it. A branch that "accepts" a value but maps it to the same glyph as another value is a latent bug; the collision is especially likely to surface when a downstream fix newly makes a previously-collapsed value reachable.
**Why**: Users act on status at a glance; two meanings sharing one glyph are indistinguishable. Where rendering is duplicated across surfaces (e.g. paired server and client templates), distinctness must hold in every copy — a correct value can still collide in a surface that was never updated.
**Frequency**: universal for any displayed enum (status badges, commit states, run/job outcomes, etc.).
**Exceptions**: values genuinely synonymous to the user may share a glyph (e.g. an `unknown` fallback), as a deliberate, commented choice — never as an accidental fallthrough.
