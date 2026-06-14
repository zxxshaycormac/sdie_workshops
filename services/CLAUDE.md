# `services/` — Design Conventions

`services/` is the orchestration layer: it composes one or more `models/` operations into business workflows, fires side effects (webhooks, notifications, mail, indexing), and exposes `Init()` entry points called from `routers/init.go`. Routers call into `services/`; `services/` never calls back up. This file goes deeper than `design.md` on the orchestration-specific rules; the cross-cutting rules (architecture, naming, errors, context, logging, i18n, testing, security) remain authoritative and are not repeated here.

---

## 1. Orchestration Role — Downward Only

**Rule**: Packages under `services/` call into `models/` (and the `xxx_model` aliases) and into sibling `services/` packages. They MUST NOT import `code.gitea.io/gitea/routers/`. A new orchestration function lives here even when the underlying data writes happen in a `models/` package.
**Why**: `services/` is the layer the router depends on; if `services/` depended back on the router, every refactor of HTTP routing would force a cascade through business logic. Placing orchestration here keeps `models/` CRUD-shaped and router-friendly.
**Frequency**: universal as a rule. ~316 files under `services/` import at least one `models/` package; `grep -rln "code.gitea.io/gitea/routers" services/` returns empty.
**Exceptions**: none in production code. `services/context/` is the context type that `routers/web/` dereferences — it lives in `services/` for historical reasons but is consumed exclusively by routers; do not add new business logic to it.

---

## 2. Transaction Boundary

**Rule**: Multi-model write operations are wrapped here, not in `routers/` and not in `models/`. The preferred idiom is `db.WithTx(ctx, func(ctx context.Context) error { ... })`, which auto-commits on nil and rolls back on error; legacy code in the same package may use `ctx, committer, err := db.TxContext(ctx); defer committer.Close(); ... committer.Commit()`. Match whichever idiom the surrounding file already uses.
**Why**: A service workflow (create repo + init git + seed labels + watch) must atomically succeed or roll back; if the router opened the transaction, every call site would have to remember to. `WithTx`/`TxContext` are transaction-aware — nested calls reuse the parent session as a `halfCommitter`, so a service calling another service inside a transaction does not double-begin.
**Frequency**: common. `db.WithTx` appears 30 times across 23 files (e.g. `services/repository/create.go:245`, `services/repository/branch.go:307,515,603`, `services/repository/fork.go:139,222`, `services/pull/merge.go:478`, `services/attachment/attachment.go:27`); `db.TxContext` appears 27 times across 19 files (e.g. `services/repository/delete.go:38,418`, `services/repository/transfer.go:104,343,435`, `services/org/org.go:23`, `services/user/user.go:63,213`). `WithTx` is the newer, dominant pattern — prefer it for new code.
**Exceptions**: Single-statement reads never open a transaction. Migrations are exempt (they run before `models/db` is wired).

---

## 3. Async Work — Queues Over Goroutines

**Rule**: Long-running or fire-and-forget work goes through `modules/queue` (`queue.CreateSimpleQueue` / `queue.CreateUniqueQueue`), with the queue created in a package-level `Init()` and started via `graceful.GetManager().RunWithCancel(queue)`. Raw `go func()` is reserved for short-lived, bounded operations whose lifetime is implicitly tied to the caller — pipeline stages wired with `sync.WaitGroup` (`services/pull/lfs.go:43-58`), fan-out for contributor-graph generation (`services/repository/contributors_graph.go:95`), or one-shot git invocations.
**Why**: Queues survive process restarts (when backed by LevelDB/Redis), bound concurrency, and report their backlog; the graceful manager wires them into shutdown so in-flight work drains rather than dies mid-statement. A bare `go someServiceFunc()` in a request path leaks work when the process stops.
**Frequency**: common. 12 `queue.CreateSimpleQueue`/`CreateUniqueQueue` call sites across `services/webhook/deliver.go:310`, `services/repository/push.go:47`, `services/repository/branch.go:572`, `services/repository/archiver/archiver.go:284`, `services/pull/check.go:390`, `services/automerge/automerge.go:36`, `services/mailer/mailer.go:410`, `services/release/tag.go:45`, `services/task/task.go:42`, `services/mirror/queue.go:43`, `services/actions/init.go:19`, `services/uinotification/notify.go:46`; 16 `graceful.GetManager().RunWith*` launches wrap those queues (some queues launch multiple workers, hence more launches than queues). Of 41 raw `go` statements outside tests, the queue launch form accounts for the majority of "background" work; the remaining raw `go` calls are bounded pipelines or explicitly-flagged TODOs.
**Exceptions**: `services/notify/notify.go:22` starts each notifier with a bare `go notifier.Run()` because the notifiers self-manage their queues internally (e.g. `services/uinotification/notify.go:63` wraps its own queue in `RunWithCancel`). Four sites launch `go AddTestPullRequestTask(...)` — `services/repository/push.go:169`, `services/pull/merge.go:177`, `services/pull/update.go:40,79` — and these are technical debt: the inline `// FIXME: graceful: AddTestPullRequestTask needs to become a queue!` at `services/pull/pull.go:330` explicitly acknowledges the gap. Do not copy this pattern.

---

## 4. Side Effects — Notify, Webhook, Mail From Services, Not Models

**Rule**: Notifications, webhooks, and email originate in `services/`. The `services/notify` package exposes a `Notifier` interface plus `RegisterNotifier`; orchestration code calls `notify_service.NewIssue(...)`, `notify_service.MergePullRequest(...)`, etc., and the registered notifiers fan out the side effect. Webhook payloads are prepared via `services/webhook.PrepareWebhooks` (`services/webhook/webhook.go:186`).
**Why**: Models should remain CRUD-only; if a model write directly triggered HTTP callbacks, the data layer would drag in network and template dependencies. Centralising the trigger in `services/` lets the same workflow run from the web UI, the API, the CLI, and migrations with the same fan-out.
**Frequency**: universal. Registered notifiers: `services/webhook/notifier.go:26`, `services/uinotification/notify.go:36`, plus notifiers in `services/mailer`, `services/feed`, `services/actions`, `services/indexer`, `services/mirror`, `services/automerge`. Each notifier's methods are the only producers of `PrepareWebhooks` calls and mail sends.
**Exceptions**: none. New event types add a method to the `Notifier` interface (`services/notify/notifier.go:18`), a `NullNotifier` no-op (`services/notify/null.go`), and a forwarding wrapper in `services/notify/notify.go`; then call the wrapper from the service that emits the event.

---

## 5. Boundary — No HTTP, No Templates, No Forms

**Rule**: `services/` does not render HTML, does not bind HTTP form parameters, and does not import `routers/`. HTTP-shaped response shaping (JSON serialisation, status codes) belongs in `routers/api/v1/`; the only `services/` package that produces wire types is `services/convert/` (model-to-API structs used by both the API router and notifiers).
**Why**: Keeping `services/` HTTP-free means a CLI subcommand or a migration importer can invoke the same workflow without dragging in the request/response stack. Binding forms in `services/` would couple the data shape to Fomantic form fields and prevent reuse.
**Frequency**: universal as a rule. `grep -rln "code.gitea.io/gitea/routers" services/` returns empty. The only `ctx.HTML`/`c.HTML` call sites in `services/` are `services/context/context_response.go` (the request-context helper used by routers) and `services/auth/sspi.go:92` (which has an inline `// FIXME: it doesn't look good to render the page here, why not redirect?` acknowledging the violation). Form-binding structs live in `services/forms/` but are imported and bound by `routers/web/`, not by other services.
**Exceptions**: `services/forms/` defines the HTTP-binding structs (e.g. `AuthenticationForm` in `services/forms/auth_form.go`) but they are data-only; binding happens in routers. `services/context/` and `services/contexttest/` define the request-context type used by routers and are not business logic — match local style when extending, do not add orchestration there.

---

## 6. Initialisation — `Init()` Called From `routers/init.go`

**Rule**: Every service that owns a queue, a background loop, or a one-time bootstrap exposes `func Init(...) error` (or `func Init()` for void bootstraps) and is wired into `routers/init.go` via `mustInit(...)` / `mustInitCtx(ctx, ...)`. The `Init` function creates the queue, registers notifiers, and launches the graceful consumer; it does NOT do per-request work.
**Why**: Centralising bootstrap in `routers/init.go` makes the startup order inspectable and ensures every queue is created before any handler can push to it. Skipping the registration leaves a nil queue that crashes on first push.
**Frequency**: common. `services/webhook/deliver.go:292` (`Init() error`), `services/repository/repository.go:98` (`Init(ctx) error`), `services/repository/archiver/archiver.go:271` (`Init(ctx) error`), `services/auth/auth.go:26` (`Init()`), `services/automerge/automerge.go:33` (`Init() error`), `services/migrations/migrate.go:485` (`Init() error`); invoked from the `mustInit(...)` / `mustInitCtx(...)` block in `routers/init.go` (~lines 129-173).
**Exceptions**: Notifier-only packages expose `NewNotifier()` + `RegisterNotifier` from inside their own `Init()` (e.g. `services/uinotification/notify.go:35-39`) rather than exposing queue creation publicly.

---

## 7. Cross-Service Composition

**Rule**: One service may call another; the import uses the `_service` alias form (`notify_service`, `repo_service`, `issue_service`, `pull_service`) to avoid colliding with the package's own receiver names. Calls flow from a higher-level orchestrator (e.g. `services/migrations/gitea_uploader.go` calling `repo_service.CreateRepositoryDirectly`) toward a lower-level one; avoid mutual recursion.
**Why**: The alias makes the call site readable (`notify_service.NewIssue(...)` is unambiguous in a file that also has a local `notify` variable) and the one-way direction keeps the dependency graph acyclic.
**Frequency**: common. `services/migrations/gitea_uploader.go:105,123` calls `repo_service`; `services/webhook/notifier.go:22` and `services/uinotification/notify.go:18` alias `notify_service`; the `_service` suffix is the established convention across every cross-service import.
**Exceptions**: none. When two services genuinely need each other, factor the shared helper into a third, lower service package.

---

## 8. Convert — The Single Wire-Format Serializer

**Rule**: Transformations from `models/` types to API-shaped structs (`api.Issue`, `api.PullRequest`, `api.Comment`, `api.Repository`) live in `services/convert/`. Both `routers/api/v1/` and service-layer notifiers (notably `services/webhook/notifier.go`) call into `convert.ToXxx(ctx, ...)` rather than re-shaping inline.
**Why**: Both the REST API and the webhook payload path must serialise the same model graph identically; centralising the conversion in `services/` (instead of `routers/api/v1/`) lets webhooks reuse it without importing the router layer.
**Frequency**: common. 77 files import `code.gitea.io/gitea/services/convert`; the consumers split roughly between `routers/api/v1/` (JSON responses) and `services/webhook/` (payload construction).
**Exceptions**: none. New API payload shapes start as a `ToXxx` function in `services/convert/` and are then surfaced by both the API handler and any webhook event that needs the same shape.
