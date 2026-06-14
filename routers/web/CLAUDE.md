# `routers/web/` — Design Conventions

`routers/web/` contains browser-facing HTTP handlers that bind request input, enforce auth via middleware, render HTML templates, and dispatch side-effecting work to `services/`. This file goes deeper than `design.md` on the router-specific rules; the cross-cutting rules (architecture, naming, errors, context, logging, i18n, testing, security) remain authoritative and are not repeated here.

---

## 1. Handler Signature — `func Xxx(ctx *context.Context)`

**Rule**: Every exported handler takes a single `ctx *context.Context` parameter, where `context` is aliased as `code.gitea.io/gitea/services/context`. Do not use the stdlib `context.Context` or invent a router-local type.
**Why**: `services/context.Context` wraps `*http.Request`, `http.ResponseWriter`, the session, the flash store, `ctx.Data` (template data map), `ctx.Repo`, `ctx.Doer`, and `ctx.Org` in one place. Using it uniformly means every handler has the same entry contract and the middleware chain (`common.Sessioner`, `context.Contexter`, `webAuth`, `context.RepoAssignment`) can populate it before dispatch.
**Frequency**: universal. ~640 handler functions across `routers/web/` (non-test) match `func ... *context.Context)` (verified `grep -rn "func.*\*context.Context)" routers/web/ --include="*.go" | grep -v "_test.go" | wc -l` = 668 including non-handler helpers; ~560 are strictly exported handlers); sample at `routers/web/repo/issue_dependency.go:16` (`AddDependency`), `routers/web/repo/contributors.go:20` (`Contributors`), `routers/web/admin/admin.go` (`Dashboard`, etc.). ~150 non-test `.go` files in `routers/web/` (verified count: 149).
**Exceptions**: Helper functions inside a handler file may take additional typed parameters (e.g. `findReadmeFileInEntries` at `routers/web/repo/view.go:79`) but the route-registered entry point is always `*context.Context`.

---

## 2. Form Binding — `web.Bind` at Route Registration, `web.GetForm(ctx)` Inside the Handler

**Rule**: Structured POST forms are declared as `services/forms/*.go` types and bound at route registration via `web.Bind(forms.XxxForm{})` (defined at `modules/web/route.go:18`). Inside the handler the bound form is retrieved with `form := web.GetForm(ctx).(*forms.XxxForm)`. For ad-hoc query/body params use `ctx.FormString`, `ctx.FormInt`, `ctx.FormBool`, `ctx.FormTrim`.
**Why**: `web.Bind` runs the Fomantic-aware binder as middleware before the handler, so validation errors are surfaced through `ctx.HasError()` consistently; `GetForm` is the single retrieval point that preserves the typed pointer. Declaring the form struct separately keeps the binding rules reusable across routes and testable in isolation.
**Frequency**: common. 128 `web.Bind(forms...)` call sites in `routers/web/web.go`; 114 `GetForm(ctx)` retrievals in handlers. Example route+handler pair: `routers/web/web.go:1192` registers `web.Bind(forms.CreateIssueForm{})` for `NewIssuePost`; the handler at `routers/web/repo/issue.go:1206` opens with `form := web.GetForm(ctx).(*forms.CreateIssueForm)`. Scalar param reads via `ctx.FormString` appear ~80 times (`routers/web/auth/auth.go:127,150,444`).
**Exceptions**: GET handlers that take only path parameters (`/repos/{id}`) and route-only flags do not need a form struct.

---

## 3. Permission Middleware — Declare at the Route, Never Check Inline

**Rule**: Authentication and authorisation are attached as middleware in `routers/web/web.go` via `verifyAuthWithOptions(&common.VerifyOptions{...})` (defined at `routers/web/web.go:138`) and the `context.RequireRepoAdmin()` / `context.RequireRepoReader(unit.TypeXxx)` / `context.RequireRepoReaderOr(...)` / `context.RepoAssignment` / `context.RepoRef()` / `context.UserAssignmentWeb()` / `context.OrgAssignment()` middleware factories. The established aliases are `reqSignIn`, `ignSignIn`, `ignSignInAndCsrf`, `adminReq`, `reqRepoAdmin`, `reqRepoCodeReader`, `reqRepoIssuesOrPullsReader` (declared at `routers/web/web.go:295-302, 681, 802-813`).
**Why**: Attaching the check at route registration covers every handler in the group and makes the auth contract inspectable in one file. Inline `if !ctx.IsSigned` checks get forgotten when routes are added; the middleware cannot. Coarse access is enforced by middleware; fine-grained per-object checks (collaborator mode on one repo) still happen in the handler after the middleware passes.
**Frequency**: universal. `verifyAuthWithOptions` produces the five `req*`/`ign*` aliases used on every route group; `context.RepoAssignment` is attached to every repo-scoped route (`routers/web/web.go:1042, 1131, 1137, 1153, 1175, 1185`). The factories live in `services/context/` (`permission.go:16,58,84`, `repo.go:402,765,884`, `user.go:15`, `org.go:276`).
**Exceptions**: Handlers that need object-level permission the middleware cannot know about (e.g. "can this user edit this specific comment?") do the check inline after middleware passes; this is correct layered enforcement, not a violation of the rule.

---

## 4. Template Rendering — `ctx.HTML(http.StatusOK, tplXxx)` with Package-Level `base.TplName` Constants

**Rule**: Render HTML with `ctx.HTML(status, name)` (defined at `services/context/context_response.go:69`). The template name is a package-level `const tplXxx base.TplName = "path/under/templates"`. Pass data to the template via `ctx.Data["Key"] = value` before rendering.
**Why**: Naming templates as constants prevents typos at the call site (the compiler checks the type even though the string is dynamic) and co-locates the list of templates a package renders at the top of the file. `ctx.HTML` also wires `TemplateName`/`TemplateLoadTimes` into `ctx.Data` in dev mode and renders the `status/500` fallback if the template fails (`services/context/context_response.go:73-92`).
**Frequency**: universal. 271 `ctx.HTML(...)` call sites in `routers/web/`; every package declares its `tplXxx` constants at file head (e.g. `routers/web/admin/admin.go:32-41`, `routers/web/repo/view.go:62-69`, `routers/web/user/home.go:43-48`). 2,436 `ctx.Data[...]` assignments carry render context.
**Exceptions**: JSON-shaped responses from a web handler (progressive-enhancement endpoints, fetch handlers) use `ctx.JSONError`/`ctx.JSONRedirect` (`routers/web/repo/issue.go:1238,1278`) or `ctx.JSON` — see Section 7.

---

## 5. Flash Messages and Redirects — Flash Then Redirect, Never Both Render

**Rule**: On a successful write, set `ctx.Flash.Success(ctx.Tr("..."))` (or `.Error`, `.Info`, `.Warning` — methods at `modules/web/middleware/flash.go:44,56,62`) and then call `ctx.Redirect(setting.AppSubURL + "/...")`. On a validation failure that must re-render the form, use `ctx.RenderWithErr(msg, tpl, form)` (`services/context/context_response.go:117`) which assigns the form back to `ctx.Data` and sets an error flash in one call.
**Why**: Flash storage is session-scoped; setting it before `ctx.Redirect` carries the message across the redirect so the next page renders it once. Rendering a template directly after setting a flash would consume it on the same request, losing it across the redirect boundary. `RenderWithErr` is the established idiom for "binding failed, redraw the form with values and error" — it pairs `middleware.AssignForm` + `Flash.Error(..., true)` + `ctx.HTML` so the user sees their input preserved.
**Frequency**: common. 390 `ctx.Flash.*` call sites and 331 `ctx.Redirect(...)` call sites across `routers/web/`. Example success path: `routers/web/auth/auth.go:528` (`Flash.Success` then redirect). Failure path: `routers/web/repo/issue.go:1237-1239` uses `ctx.HasError()` then `ctx.JSONError(ctx.GetErrMsg())`.
**Exceptions**: Endpoints invoked via fetch/XHR that expect JSON, not a redirect, use `ctx.JSONError`/`ctx.JSONRedirect` instead — the flash is for human browser flows.

---

## 6. Error Handling — `ctx.ServerError` for 5xx, `ctx.NotFound` for 4xx-Not-Found, Never `log.Error` Then Return

**Rule**: Surface errors through `ctx.ServerError(logMsg, err)` (`services/context/context_response.go:158`) or `ctx.NotFound(logMsg, err)` (`context_response.go:126`); both log at the right level and render the appropriate status page. Do not call `log.Error` and return silently — the helper does the logging and the response. For "expected" failures that are not errors (e.g. feature disabled), use `ctx.NotFound("MustEnableIssues", nil)` with a descriptive log message and nil error.
**Why**: `ctx.ServerError` renders the `status/500` template in production and dumps the stack/error in dev; `ctx.NotFound` honours the `Accept` header to return either the 404 HTML page or a plain `Not found.` body for non-HTML clients (`context_response.go:138-148`). Calling `log.Error` and returning leaves the response half-written or empty.
**Frequency**: universal. 1,048 `ctx.ServerError(...)` call sites in `routers/web/`; `ctx.NotFound` is used for both missing resources and feature-gating (`routers/web/repo/issue.go:118,132`).
**Exceptions**: Validation errors that should re-render a form use `ctx.RenderWithErr` (Section 5) or `ctx.JSONError` (Section 7) rather than a 5xx.

---

## 7. Boundary — Call `services/`, Not `models/` for Writes; Read-Through to `models/` Is Legacy

**Rule**: New handlers orchestrate side effects through `services/` packages (e.g. `issue_service`, `repo_service`, `pull_service`, `notify_service`). They MUST NOT call write-path `models/` functions directly and MUST NOT invoke `db.GetEngine`/`db.WithTx`/`db.TxContext`. Reading reference data from `models/` (e.g. `repo_model.GetRepoAssignees`, enumeration queries) is the established pattern; transactional multi-step writes go through a service.
**Why**: The same workflow runs from the web UI, the API, and the CLI; placing it in `services/` gives one call site that fires webhooks, notifications, mail, and indexing consistently. If a handler opens the transaction, every call site has to remember to commit, and the side-effect fan-out (`services/notify`) gets skipped on the CLI path.
**Frequency**: aspirational for writes, common for reads. 144 files in `routers/web/` import at least one `services/` package; 132 files import at least one `models/` package (overwhelmingly read-side helpers like `repo_model`, `user_model`, `unit`, `perm`). Verified direct engine access: `routers/web/healthcheck/check.go:95` calls `db.GetEngine(ctx).Ping()` (acceptable — it is the healthcheck). Dropdown-population helpers walk reference data (`routers/web/repo/issue.go:556,564` call `db.Find[issues_model.Milestone]` to populate the milestone dropdown on the New Issue form).
**Exceptions**: Reference-data lookups in handlers (assignees, milestones, labels, units, permission enums) are accepted reads; matching existing file style when extending them is fine. New write paths must go through `services/`.

---

## 8. Route Registration — One File (`web.go`), Middleware Order Matters

**Rule**: All routes are registered in `routers/web/web.go` inside `registerRoutes(m *web.Route)` (`routers/web/web.go:298`). The middleware chain is applied left-to-right; order is significant: auth (`reqSignIn`/`ignSignIn`) → assignment (`context.RepoAssignment`/`UserAssignmentWeb`/`OrgAssignment`) → permission (`reqRepoAdmin`/`reqRepoCodeReader`) → ref resolution (`context.RepoRef()`) → binder (`web.Bind(...)`) → handler.
**Why**: `context.RepoAssignment` must run before any `reqRepoXxx` check because the permission check reads `ctx.Repo.Permission`; `context.RepoRef()` must run after `RepoAssignment` because it walks `ctx.Repo.Repository`'s git refs. Misordering produces nil-pointer panics or skipped permission checks.
**Frequency**: universal. The middleware aliases are declared once (`routers/web/web.go:295-302, 681, 802-813`) and reused across every route group; the canonical ordering appears at `routers/web/web.go:1131,1137,1152,1191-1193`.
**Exceptions**: Top-of-chain middleware that runs once per request (`common.Sessioner`, `context.Contexter`, `webAuth`, `user.GetNotificationCount`, `repo.GetActiveStopwatch`, `goGet`) is appended to the global `mid` slice in `Routes()` (`routers/web/web.go:270-286`) rather than repeated per route.

---

## 9. Localisation — Every User-Visible String via `ctx.Tr`

**Rule**: All strings the user will see (flash messages, template titles, JSON error messages) come from locale keys via `ctx.Tr("key.name", args...)`. Template-name constants (`tplXxx base.TplName`) are programmatic identifiers and are not localised.
**Why**: `design.md` Section 6 establishes this as a cross-cutting rule; `routers/web/` is the densest source of user-visible strings because every HTML page and every flash is rendered here. Hardcoded English forces every translator to chase code changes.
**Frequency**: universal. `ctx.Tr(...)` appears in nearly every handler that sets a flash or `ctx.Data["Title"]` (e.g. `routers/web/admin/admin.go:134`, `routers/web/auth/auth.go:528`).
**Exceptions**: Log messages passed to `ctx.ServerError`/`ctx.NotFound` are operator-facing and stay English; HTTP header values are not localised.
