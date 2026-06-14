# `cmd/` — Design Conventions

`cmd/` is the entry-point layer: it parses CLI arguments with `urfave/cli/v2`, builds the root `context.Context`, initialises config (and, for the `web` subcommand, the full `routers.InitWebInstalled` chain), then dispatches to `services/` or `routers/`. This file goes deeper than `design.md` on the CLI-specific rules; the cross-cutting rules (architecture, naming, errors, context, logging, i18n, testing, security) remain authoritative and are not repeated here.

---

## 1. Command Registration — `urfave/cli/v2` Package Globals

**Rule**: Every subcommand is a package-level `var CmdXxx = &cli.Command{...}` with `Name`, `Usage`, optional `Description`, optional `Before: PrepareConsoleLoggerLevel(...)`, `Action: runXxx`, optional `Flags []cli.Flag`, and optional `Subcommands`. The `Action` is a `func(*cli.Context) error` named `runXxx`. Commands with nested children (e.g. `CmdAdmin`'s `auth`, `regenerate`, `user` groups) set `Subcommands` and leave `Action` unset; leaf commands set `Action`.
**Why**: `urfave/cli` discovers commands from the `app.Commands` slice registered in `cmd/main.go:NewMainApp`; the package-global form lets the registration list reference them by name (`CmdWeb`, `CmdServ`, ...) without init-order surprises. Naming the action `runXxx` matches the convention across every command file.
**Frequency**: universal. 17 top-level commands registered in `cmd/main.go:130-152` (14 in `subCmdWithConfig` at lines 130-146, 3 in `subCmdStandalone` at lines 149-152); every leaf command follows the `runXxx` pattern (`runWeb` at `cmd/web.go:228`, `runServ` at `cmd/serv.go:127`, `runDoctorCheck` at `cmd/doctor.go:162`, `runCreateUser` at `cmd/admin_user_create.go:71`, `runDeleteUser` at `cmd/admin_user_delete.go:44`).
**Exceptions**: `cmdHelp()` at `cmd/main.go:19` returns `*cli.Command` (constructed per-call) so it can be appended as the `help` subcommand of every parent command by `prepareSubcommandWithConfig`.

---

## 2. Init Order — Config First, Services Last

**Rule**: CLI subcommands do NOT re-implement the init chain. For one-shot commands the action calls `initDB(ctx)` (declared at `cmd/cmd.go:60`) — which runs `setting.MustInstalled`, `setting.LoadDBSetting`, `setting.InitSQLLoggersForCli`, and `db.InitEngine(ctx)` — and only then dispatches to a service. For the long-running `web` command the action delegates the full boot sequence to `routers.InitWebInstalled(ctx)` (defined at `routers/init.go:114`), whose order is: `git.InitFull` → i18n → `setting.LoadSettings` → `storage.Init` → mailer/cache/feed/notification/archiver → markup → `common.InitDBEngine` → `system.Init` → `oauth2.Init` → `release_service.Init` → `models.Init` / `authmodel.Init` / `repo_service.Init` → indexer/mirror/webhook/pull/automerge/task/migrations → `ssh.Init` → `auth.Init` → `actions_service.Init` → `cron.NewContext`. `InitWebInstallPage` (`routers/init.go:107`) is the lighter variant used before the install lock is set.
**Why**: Services assume their queue is registered before any handler pushes to it; the DB engine must exist before any model call. Letting `routers/init.go` own the order keeps it inspectable in one place and prevents a CLI action from forgetting a step. `initDB` is the deliberately-truncated subset for commands that only need the engine (admin user edits, dump, doctor checks).
**Frequency**: universal. `routers.InitWebInstalled(graceful.GetManager().HammerContext())` is invoked exactly once at `cmd/web.go:194`; `installSignals` + `initDB` together appear in 27 of the 32 command files that declare an `Action:` (the dominant pattern in `cmd/admin_*`, `cmd/dump.go`, `cmd/keys.go`, `cmd/hook.go`, `cmd/mailer.go`, `cmd/doctor*.go`).
**Exceptions**: `migrate` and `migrate_storage` call `db.InitEngineWithMigration(context.Background(), migrations.Migrate)` (`cmd/migrate.go:39`, `cmd/migrate_storage.go:198`) instead of `initDB`, because they must run pending migrations before normal use. Standalone commands (`CmdCert`, `CmdGenerate`, `CmdDocs`) skip init entirely — they are registered in `subCmdStandalone` (`cmd/main.go:149-152`) precisely because they need neither config nor DB.

---

## 3. Context Origin — `cmd/` Owns the Root `context.Context`

**Rule**: The `context.Context` that flows through every lower layer originates in `cmd/`. The two established roots are: (a) `installSignals()` (`cmd/cmd.go:76`) returns `(ctx, cancel)` built from `context.WithCancel(context.Background())` and wires `SIGINT`/`SIGTERM` via `signal.Notify` to `cancel()`; (b) `cmd/web.go:235` constructs `managerCtx, cancel := context.WithCancel(context.Background())` and hands it to `graceful.InitManager(managerCtx)` at `cmd/web.go:236` for the web server's lifecycle. The `cli.Context` itself is bridged via `c.Context` (used at `cmd/admin_user_create.go:95`); when a one-shot command needs cancellation it creates its own `installSignals()` pair and passes that `ctx` downward.
**Why**: There is no HTTP request to provide a context for CLI flows, so the cancellation root must be constructed explicitly. `installSignals()` ensures Ctrl-C propagates as `ctx.Done()` into DB sessions and git subprocesses; `graceful.InitManager` extends the same idea to the long-running web process so HTTP handlers and queues observe a single shutdown signal. Lower layers must not invent their own root — see `design.md` Section 4.
**Frequency**: universal. `installSignals()` is invoked at the top of nearly every one-shot `runXxx` action (`runServ` at `cmd/serv.go:128`, `runDeleteUser` at `cmd/admin_user_delete.go:49`, `runCreateUser` via the same idiom at `cmd/admin_user_create.go:100`); `context.Background()` appears as a root only at `cmd/dump_repo.go:91,180`, `cmd/serv.go:70`, `cmd/web.go:221,235`, and inside `installSignals` itself.
**Exceptions**: Tests use `context.Background()` directly (`cmd/hook_test.go:18`, `cmd/migrate_storage_test.go:56`). Long-running goroutines spawned by the web action (e.g. `servePprof` at `cmd/web.go:219`) call `process.GetManager().AddTypedContext(context.Background(), ...)` (at `cmd/web.go:221`) to register with the process manager rather than tie the pprof server to request lifetimes.

---

## 4. Boundary — Dispatch to `services/`, Not Direct DB/model Code

**Rule**: New CLI actions orchestrate via `services/` (e.g. `user_service.DeleteUser` at `cmd/admin_user_delete.go:80`, `user_service.UpdateAuth` at `cmd/admin_user_change_password.go:65`). They must not reach into `models/` for write operations, must not call `db.GetEngine`, and must not hold a `*xorm.Engine`. Validation, arg parsing, and confirming (`argsSet` at `cmd/cmd.go:28`, `confirm` at `cmd/cmd.go:42`) belong here; the workflow belongs in a service.
**Why**: The same workflow runs from the web UI, the API, and the CLI; placing it in `services/` gives one call site for the orchestration and lets the CLI stay concerned with flags and I/O. Holding an engine across a CLI invocation also risks leaking a connection that the long-running process will never reclaim.
**Frequency**: common, but not yet universal — this is the aspirational rule. Verified direct violations: `cmd/doctor.go:133` passes a `func(x *xorm.Engine) error` into `db.InitEngineWithMigration` for schema recreation; `cmd/migrate_storage.go` calls `db.Iterate(...)` directly on eight model types (Attachment, LFSMetaObject, User, Repository, RepoArchiver, PackageBlob, ActionTask, ActionArtifact — call sites at lines 100, 107, 114, 124, 134, 142, 150, 166); `cmd/admin_user_create.go:128` calls `db.IsTableNotEmpty(&user_model.User{})`; `cmd/admin.go:161` calls `db.Count[repo_model.Release](...)`; `cmd/admin_auth.go:67` calls `db.Find[auth_model.Source](ctx, ...)`; `cmd/dump.go:220` calls `db.DumpDatabase(...)`. These are technical debt, not precedent.
**Exceptions**: When a CLI command is the only caller of an operation (schema doctoring, full-DB dump, storage migration iterating every row), reaching for the `db.` helper is pragmatic and accepted — the alternative would be a service function with a single CLI caller. For new code, prefer extracting the helper into a service and dispatching from the command; only fall back to direct `db.` calls when the operation is inherently CLI-shaped (whole-DB iteration, raw engine access for migrations).

---

## 5. Action Helpers — Shared Utilities in `cmd/cmd.go`

**Rule**: Reusable CLI helpers live in `cmd/cmd.go` and are referenced by every command that needs them: `argsSet(c, names...)` at `cmd/cmd.go:28` enforces required flags; `confirm()` at `cmd/cmd.go:42` prompts y/n; `initDB(ctx)` at `cmd/cmd.go:60` boots the engine; `installSignals()` at `cmd/cmd.go:76` builds the cancellable root context; `setupConsoleLogger(level, colorize, out)` at `cmd/cmd.go:98` swaps the console writer (panics if `out` is not `os.Stdout`/`os.Stderr` — see `cmd/cmd.go:99`); `PrepareConsoleLoggerLevel(defaultLevel)` at `cmd/cmd.go:123` returns a `Before` hook that honours the global `--quiet`/`--verbose`/`--debug` flags. Per-command `setup()` shims exist for legacy commands (e.g. `cmd/serv.go` calls `setup(ctx, c.Bool("debug"))`).
**Why**: Duplicating the signal-handler + engine-init dance in every action would invite drift; centralising it in `cmd/cmd.go` keeps the cancellation contract consistent. The `Before` hook form lets each command declare its own default log level (`web` uses `INFO` at `cmd/web.go:40`, `serv` uses `FATAL` at `cmd/serv.go:47` because any stdout noise breaks the git protocol) while still honouring user overrides.
**Frequency**: universal. `installSignals` + `initDB` is the opening stanza of 27 of the 32 command files containing an `Action:` field; `PrepareConsoleLoggerLevel` is referenced from `cmd/cmd.go`, `cmd/keys.go`, `cmd/hook.go`, `cmd/serv.go`, `cmd/main.go`, `cmd/web.go`.
**Exceptions**: The `web` action bypasses `initDB` because `routers.InitWebInstalled` does the equivalent work as part of the full boot chain. Standalone commands (`cert`, `generate`, `docs`) bypass both helpers because they neither read config nor touch the DB.

---

## 6. Logging on stdout — Forbidden for Git-Protocol Commands

**Rule**: Any command that runs inside an SSH or git remote-helper pipe (`serv`, `hook`, `keys`) sets `Before: PrepareConsoleLoggerLevel(log.FATAL)` so no log line ever reaches stdout. Commands that interact with a human at a terminal (`web`, `admin`, `doctor`) default to `log.INFO`.
**Why**: The comment at `cmd/cmd.go:121` states it directly: "some sub-commands (for git/ssh protocol) shouldn't output any log to stdout. Any log appears in git stdout pipe will break the git protocol, eg: client can't push and hangs forever." Git parses stdout byte-for-byte as pack/protocol data; an interleaved `log.Info` corrupts the stream silently.
**Frequency**: universal for git-pipe commands. `cmd/serv.go:47` and `cmd/hook.go` (via the same `PrepareConsoleLoggerLevel(log.FATAL)` `Before`) enforce the FATAL ceiling; `cmd/keys.go` likewise. Operators who need diagnostics from these commands must use `--debug` (which routes through the level override inside `PrepareConsoleLoggerLevel`) or set a log file.
**Exceptions**: Even in FATAL-ceiling commands, fatal errors that exit the process may print to stderr (git treats stderr as user-visible progress, not protocol); the rule is about stdout specifically. The `web` command is exempt because it never speaks the git wire protocol directly — git operations go through `serv`/`hook`.

---

## 7. Flag Inheritance — Global Flags Prepended, Short Flags Reserved

**Rule**: `appGlobalFlags()` (`cmd/main.go:50`) declares `-C/--custom-path`, `-c/--config`, `-w/--work-path` as globally shared. `prepareSubcommandWithConfig` (`cmd/main.go:77`) prepends these to every subcommand's flags and wraps its `Action` with `prepareWorkPathAndCustomConf`, which walks `ctx.Lineage()` to find the first ancestor that set each flag. Short single-letter aliases are reserved at the app level — subcommands must not reuse `-C`, `-c`, or `-w`.
**Why**: urfave/cli's per-level `Before` runs once per nesting depth, so doing work-path resolution in `Before` would fire multiple times; wrapping `Action` instead guarantees a single init. Reserving the short flags globally prevents the ambiguity where `gitea admin user create -c X` could mean either config or a subcommand-defined `-c`.
**Frequency**: universal. The wrap is applied to every entry in `subCmdWithConfig` via the loop at `cmd/main.go:162-164`; `app.HideHelp = true` (`cmd/main.go:160`) forces the custom `cmdHelp()` so the DEFAULT CONFIGURATION block (work path, custom path, config file) is printed alongside help.
**Exceptions**: Standalone commands (`cert`, `generate`, `docs`) intentionally skip the wrap (they are appended at `cmd/main.go:166` without `prepareSubcommandWithConfig`) because they do not read the config file and must work even before installation.
