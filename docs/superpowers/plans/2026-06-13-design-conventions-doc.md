# Design Conventions Documentation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Create `design.md` + 8 directory `CLAUDE.md` files that document Gitea's design conventions as an AI assistant reference.

**Architecture:** Each file is independent after Task 1. Every rule must be backed by observable code patterns (grep-verifiable). Documentation is principles-first; no code dumps. Engineer explores → drafts → verifies → writes → (optional) commits per file.

**Tech Stack:** Markdown; grep/find/Read for exploration; Gitea Go/JS/CSS/templates as source material.

**Spec:** `docs/superpowers/specs/2026-06-13-design-conventions-doc-design.md`

**Commits:** User requested no commits during the active session. Each task ends with an OPTIONAL commit step — skip unless the user explicitly asks to commit.

---

## File Structure

| # | File | Responsibility | Target Size |
|---|------|----------------|-------------|
| 1 | `design.md` | Cross-cutting principles (8 topics) | <400 lines |
| 2 | `modules/CLAUDE.md` | Isolation rule, sub-package org | <150 lines |
| 3 | `models/CLAUDE.md` | XORM models, db.Session, migrations, unittest | <150 lines |
| 4 | `services/CLAUDE.md` | Orchestration, transaction boundary, async | <150 lines |
| 5 | `cmd/CLAUDE.md` | CLI structure, init order, context origin | <150 lines |
| 6 | `routers/web/CLAUDE.md` | Web handler patterns, middleware, render | <150 lines |
| 7 | `routers/api/v1/CLAUDE.md` | Swagger, binding, response, versioning | <150 lines |
| 8 | `web_src/CLAUDE.md` | JS org, Vue, CSS, Fomantic | <150 lines |
| 9 | `templates/CLAUDE.md` | Inheritance, helpers, i18n in templates | <150 lines |

Tasks can run in any order after Task 1, but the sequence below is recommended (each builds context useful for the next).

---

### Task 1: Root `design.md` — Cross-Cutting Principles

**Files:**
- Create: `design.md`
- Spec ref: Section "`design.md` Outline (Topic-First / Plan B)"

**Target:** 8 topic sections, ~300-400 lines total. Uniform format per rule: `**Rule**: statement` / `**Why**: rationale` / `**Frequency**: universal|common|occasional`.

- [ ] **Step 1: Explore architecture evidence**

Run from repo root:

```bash
# Verify modules/ isolation rule (spec says modules can't import models/services/routers)
grep -rln "code.gitea.io/gitea/models\b\|code.gitea.io/gitea/services\|code.gitea.io/gitea/routers" modules/ | head -20
```

Expected: a small number of files (currently known: `modules/metrics/collector.go`, `modules/eventsource/manager_run.go`, `modules/repository/env.go`). Investigate each — are they legitimate exceptions or violations? Document them in the "Exceptions" subsection.

```bash
# Confirm layered dependency direction
grep -rln "code.gitea.io/gitea/cmd" models/ modules/ services/ routers/ 2>/dev/null
# Expected: empty (no lower layer imports cmd/)
```

- [ ] **Step 2: Explore naming patterns**

```bash
# Package naming (lowercase single word expected)
find . -name "*.go" -not -path "*/vendor/*" -not -path "*/node_modules/*" -not -name "*_test.go" | head -100 | xargs grep "^package" | awk '{print $2}' | sort -u | head -30

# File naming (snake_case.go expected)
find models/ services/ modules/ -maxdepth 2 -name "*.go" -not -name "*_test.go" | head -30
```

Note: Gitea uses `snake_case.go` for filenames. Confirm and document.

- [ ] **Step 3: Explore error handling patterns**

```bash
# Wrapping pattern
grep -rn "fmt.Errorf" --include="*.go" services/ | grep "%w" | head -10
# Expected: fmt.Errorf("description: %w", err) pattern

# Sentinel error helpers (note: uses util.NewNotExistErrorf, not errors.New)
grep -rn "^var Err" --include="*.go" models/ | head -10
grep -rn "util.New.*Errorf\|util.New.*Error(" --include="*.go" models/ | head -10

# Error swallowing anti-pattern (should be rare)
grep -rn "_ = " --include="*.go" services/ models/ | grep -v "_test.go" | grep -v "defer" | head -10
```

Document the helper-based error definition style. Read `modules/util/error.go` to understand the error type hierarchy.

- [ ] **Step 4: Explore context propagation**

```bash
# Context as first param
grep -rn "func.*ctx context.Context" --include="*.go" models/ services/ | wc -l
grep -rn "func.*context.Context" --include="*.go" models/db/db.go 2>/dev/null | head -5

# DefaultContext for non-request flows
grep -rn "db.DefaultContext" --include="*.go" | head -5
```

- [ ] **Step 5: Explore logging patterns**

```bash
ls modules/log/*.go | head
grep -rn "log.Error\|log.Warn\|log.Info\|log.Debug" --include="*.go" services/ | grep -v "_test.go" | head -15
```

Read `modules/log/log.go` (or equivalent entry point) to identify the structured-fields API.

- [ ] **Step 6: Explore i18n usage**

```bash
grep -rn "\.Tr(" --include="*.go" routers/web/ | head -10
ls options/locale/ | head -10
```

Document: locale files in `options/locale/`, `Tr()` is called on the i18n context.

- [ ] **Step 7: Explore testing patterns**

```bash
grep -rln "code.gitea.io/gitea/models/unittest" --include="*_test.go" | wc -l
ls models/unittest/
cat models/unittest/testdb.go | head -50  # understand the setup API
grep -rn "\.Skip(" --include="*_test.go" | head -5  # check skip usage
```

- [ ] **Step 8: Explore security patterns**

```bash
# Permission middleware names
grep -rn "RequireSignIn\|MustRepo\|context.RepoAssigned\|perm.CheckUnit" routers/common/ routers/web/ 2>/dev/null | head -10

# SQL injection avoidance (XORM-only)
grep -rn "\.SQL(\|\.Raw(" --include="*.go" models/ | head -5
```

- [ ] **Step 9: Draft Section 0 (Document Positioning)**

Write 2-3 paragraphs covering:
- Audience: AI assistants (Claude, etc.) modifying Gitea code
- How to use this doc: read for cross-cutting; layer on directory `CLAUDE.md` for directory-specific rules
- Full index of directory `CLAUDE.md` files with one-line description each

- [ ] **Step 10: Draft Section 1 (Architecture Principles)**

Cover these rules (verify each against Step 1 evidence):
- Layered dependency: `cmd → routers → services → models → modules`
- `modules/` MUST NOT import from `models/`, `services/`, `routers/` (with documented exceptions from Step 1)
- No circular deps within a layer
- Services orchestrate; models do CRUD only; routers bind/validate/dispatch

For each rule, include any known violations found in Step 1 as "Exceptions" so AI doesn't enforce blindly.

- [ ] **Step 11: Draft Section 2 (Naming Conventions)**

Cover:
- Package names: lowercase, single word (confirm from Step 2)
- Exported types: PascalCase
- Receivers: short consistent abbreviations (e.g., `u *User`, not `user *User`)
- File naming: `snake_case.go`; test files: `xxx_test.go`
- DB column naming: snake_case via XORM tags

- [ ] **Step 12: Draft Section 3 (Error Handling)**

Cover (verify against Step 3 evidence):
- Wrap with `fmt.Errorf("description: %w", err)` — wrap at each layer with what-was-being-done context
- Sentinel errors via `util.NewNotExistErrorf`, `ErrXxx` package-level vars
- Return errors up; log only at boundary (router/CLI)
- Never swallow with `_ =`

- [ ] **Step 13: Draft Sections 4-8**

Use evidence from Steps 4-8. For each section:
- 3-5 rules
- Each with `**Rule**`, `**Why**`, `**Frequency**`
- Include `**Exceptions**` only when actually observed

- [ ] **Step 14: Verify every rule against grep evidence**

Re-read each rule. For each, run a grep that confirms the pattern exists in code. Rules with no evidence → mark as "aspirational" or remove.

- [ ] **Step 15: Write the file**

Use the Write tool. Verify size:

```bash
wc -l design.md
# Expected: < 400 lines
```

- [ ] **Step 16: Optional commit**

```bash
git add design.md
git commit -m "docs: add design.md with cross-cutting design conventions"
```

---

### Task 2: `modules/CLAUDE.md`

**Files:**
- Create: `modules/CLAUDE.md`
- Spec ref: "`modules/CLAUDE.md`" scope

**Target:** <150 lines. Opens with one-sentence statement of `modules/` responsibility.

- [ ] **Step 1: Explore modules/ structure and isolation evidence**

```bash
# Sub-package inventory
ls -d modules/*/ | head -30

# Documented isolation violations (from Task 1 Step 1)
grep -rln "code.gitea.io/gitea/models\b\|code.gitea.io/gitea/services\|code.gitea.io/gitea/routers" modules/

# For each violation file, read the imports section to understand WHY
head -30 modules/metrics/collector.go
head -30 modules/eventsource/manager_run.go
head -30 modules/repository/env.go
```

- [ ] **Step 2: Explore public API patterns**

```bash
# Sample several modules to identify API export conventions
cat modules/queue/queue.go 2>/dev/null | head -50
cat modules/sync/exclusive.go 2>/dev/null | head -30
```

Look for: how types are exported, whether internal types leak, naming of constructors.

- [ ] **Step 3: Explore test conventions in modules/**

```bash
# Confirm: modules tests should NOT use unittest (no DB)
grep -rln "models/unittest" modules/ | head
# Expected: empty (modules tests don't touch DB)

# Sample test file structure
cat modules/util/utility_test.go 2>/dev/null | head -40
```

- [ ] **Step 4: Draft rules**

Open with: "`modules/` contains shared libraries and utilities used across all layers. It MUST NOT depend on `models/`, `services/`, or `routers/`."

Cover these rules (from spec):
- Isolation rule: never import `models/`, `services/`, `routers/` — with the documented exceptions from Step 1 (investigate WHY each exception exists)
- Sub-package organization: single responsibility per sub-package; avoid catch-all `util/`
- Public API design: hide internal types, don't leak `context.Context` from internal APIs unnecessarily
- Test requirement: no DB dependency, runnable in isolation; if a test needs DB, it belongs in `models/` not `modules/`

Add 2-3 more rules based on Step 2 observations.

- [ ] **Step 5: Verify each rule**

Run greps to confirm. Any rule without evidence → remove or mark "aspirational".

- [ ] **Step 6: Write file**

```bash
wc -l modules/CLAUDE.md
# Expected: < 150 lines
```

- [ ] **Step 7: Optional commit**

```bash
git add modules/CLAUDE.md
git commit -m "docs(modules): add directory CLAUDE.md with isolation and org rules"
```

---

### Task 3: `models/CLAUDE.md`

**Files:**
- Create: `models/CLAUDE.md`
- Spec ref: "`models/CLAUDE.md`" scope

**Target:** <150 lines.

- [ ] **Step 1: Explore XORM model patterns**

```bash
# Sample model definitions
cat models/user/user.go 2>/dev/null | head -80
cat models/repo/repo.go 2>/dev/null | head -80

# Bean registration via init()
grep -rn "db.RegisterModel\|RegisterModel" --include="*.go" models/ | head -10

# XORM tag styles
grep -rn "xorm:" --include="*.go" models/user/user.go models/repo/repo.go | head -20
```

- [ ] **Step 2: Explore db.Session / DefaultContext usage**

```bash
cat models/db/engine.go 2>/dev/null | head -80
grep -rn "func.*context.Context.*Engine\|func.*Engine.*context.Context" --include="*.go" models/db/ | head -10
grep -rn "db.GetEngine\|db.DefaultContext" --include="*.go" models/ | head -10
```

- [ ] **Step 3: Explore migration rules**

```bash
ls models/migrations/ | head -20
ls models/migrations/ | tail -20
cat models/migrations/migrations.go 2>/dev/null | head -60

# Confirm: migrations are versioned sequentially
grep -rn "func.*up.*int64\|migrations\[.*\] = " models/migrations/migrations.go 2>/dev/null | head
```

- [ ] **Step 4: Explore unittest framework**

```bash
cat models/unittest/testdb.go | head -80
cat models/unittest/fixtures.go | head -60
ls models/unittest/fixtures/ 2>/dev/null | head
```

- [ ] **Step 5: Draft rules**

Open with: "`models/` is the data access layer: XORM models, DB sessions, migrations, and the unittest framework. No business logic."

Cover:
- XORM model definition: struct with `xorm:` tags, registered via `init()` calling `db.RegisterModel(new(Type))`, table name auto-derived or explicit
- `db.Session` / `db.DefaultContext`: acquire via `db.GetEngine(ctx)`, propagate ctx, never store engine in package global
- Migration rules: sequential version, immutable once released, up-only, no down migrations
- Unittest: every DB-backed test uses `unittest.InitTestDB` / `unittest.LoadFiles`, fixtures in `models/unittest/fixtures/`
- Boundary: data integrity + CRUD; no business logic; no HTTP; no template rendering
- Transactions: caller in `services/` opens, `models/` functions accept ctx+engine

Add 2-3 rules from observation.

- [ ] **Step 6: Verify each rule**

- [ ] **Step 7: Write file & verify size**

```bash
wc -l models/CLAUDE.md
```

- [ ] **Step 8: Optional commit**

```bash
git add models/CLAUDE.md
git commit -m "docs(models): add directory CLAUDE.md with XORM and migration rules"
```

---

### Task 4: `services/CLAUDE.md`

**Files:**
- Create: `services/CLAUDE.md`
- Spec ref: "`services/CLAUDE.md`" scope

**Target:** <150 lines.

- [ ] **Step 1: Explore orchestration patterns**

```bash
ls -d services/*/ | head -20
cat services/repository/repository.go 2>/dev/null | head -50
cat services/pull/pull.go 2>/dev/null | head -50
```

- [ ] **Step 2: Explore transaction boundary**

```bash
# Transaction wrapping at services layer
grep -rn "db.WithTx\|db.Tx\|sess.Begin\|sess.Commit" --include="*.go" services/ | head -15
```

- [ ] **Step 3: Explore async work**

```bash
# Queue usage (not raw goroutines)
grep -rn "queue.*Add\|queue.*Push" --include="*.go" services/ | head -10
# Raw goroutines (should be rare / justified)
grep -rn "^\s*go " --include="*.go" services/ | grep -v "_test.go" | head -10
```

- [ ] **Step 4: Explore webhook/notify triggers**

```bash
grep -rn "Notify\|webhook.\|NewNotifier" --include="*.go" services/ | head -10
```

- [ ] **Step 5: Draft rules**

Open with: "`services/` orchestrates business logic: calls `models/`, called by `routers/`. Owns transactions and async dispatch."

Cover:
- Orchestration: services call `models/`; routers call `services/`; never the reverse
- Transaction boundary: multi-model operations wrapped in `db.WithTx(func(ctx) {...})` here, not in routers
- Async work: use `queue` package, never raw `go` for fire-and-forget (raw `go` only for short-lived operations with explicit lifecycle)
- Webhooks/notifications: trigger from here, not from `models/`
- Boundary: do not render templates, do not bind HTTP forms, do not import `routers/`

- [ ] **Step 6: Verify each rule**

- [ ] **Step 7: Write file & verify size**

```bash
wc -l services/CLAUDE.md
```

- [ ] **Step 8: Optional commit**

```bash
git add services/CLAUDE.md
git commit -m "docs(services): add directory CLAUDE.md with orchestration rules"
```

---

### Task 5: `cmd/CLAUDE.md`

**Files:**
- Create: `cmd/CLAUDE.md`
- Spec ref: "`cmd/CLAUDE.md`" scope

**Target:** <150 lines.

- [ ] **Step 1: Explore CLI structure**

```bash
ls cmd/ | head -30
cat cmd/main.go 2>/dev/null | head -50  # or whatever the root is
```

- [ ] **Step 2: Explore command registration pattern**

```bash
# Find the cli library used (urfave/cli likely)
grep -rn "urfave/cli\|flag.\|cli.Command" cmd/ | head -10
cat cmd/web.go 2>/dev/null | head -80
```

- [ ] **Step 3: Explore init order**

```bash
grep -rn "InitWebInstalled\|routers.InitWeb\|InitDB\|InitLog" cmd/ routers/init.go | head -15
cat routers/init.go 2>/dev/null | head -100
```

Document the precise init sequence (config → DB → git → i18n → storage → services → SSH).

- [ ] **Step 4: Draft rules**

Open with: "`cmd/` holds CLI entry points. Each subcommand dispatches to `services/`; no business logic here."

Cover:
- Command registration: urfave/cli pattern with flags and subcommands
- Init order: aligns with `routers.InitWebInstalled()` — config first, services last
- Context origin: `ctx context.Context` originates here (often `context.Background()` or signal-aware ctx), propagates downward
- Boundary: dispatch to `services/` only; no direct DB/model manipulation

- [ ] **Step 5: Verify each rule**

- [ ] **Step 6: Write file & verify size**

```bash
wc -l cmd/CLAUDE.md
```

- [ ] **Step 7: Optional commit**

```bash
git add cmd/CLAUDE.md
git commit -m "docs(cmd): add directory CLAUDE.md with CLI and init rules"
```

---

### Task 6: `routers/web/CLAUDE.md`

**Files:**
- Create: `routers/web/CLAUDE.md`
- Spec ref: "`routers/web/CLAUDE.md`" scope

**Target:** <150 lines.

- [ ] **Step 1: Explore handler signatures**

```bash
ls routers/web/ | head -20
cat routers/web/repo/view.go 2>/dev/null | head -60  # sample handler
grep -rn "func.*context.Context" routers/web/repo/ | head -10
```

- [ ] **Step 2: Explore form binding**

```bash
grep -rn "ctx.Form\|c.Form\|ParseForm\|FormString" routers/web/ | head -15
```

- [ ] **Step 3: Explore middleware / perm checks**

```bash
ls routers/web/repo/setting/ 2>/dev/null | head
grep -rn "RequireSignInView\|MustRepo\|context.RepoRef\|RepoMustExist" routers/web/ routers/common/ | head -10
cat routers/web/web.go 2>/dev/null | head -80  # route registration
```

- [ ] **Step 4: Explore template render**

```bash
grep -rn "ctx.HTML\|ctx.Render\|\.HTML(" routers/web/ | head -15
grep -rn "ctx.Flash\|ctx.Data\[" routers/web/ | head -15
grep -rn "ctx.Redirect\|ctx.JSON" routers/web/ | head -10
```

- [ ] **Step 5: Draft rules**

Open with: "`routers/web/` handles browser-facing HTTP: bind form, check perms, call services, render template."

Cover:
- Handler signature: `func(c *context.Context)` (from `services/context`)
- Form binding via `ctx.FormString`/`ctx.FormInt` etc. at handler entry
- Permission middleware chain in route registration (`MustRepo`, `RequireSignInView`)
- Template render: `ctx.HTML(200, "template-name", data)`
- Flash messages: `ctx.Flash.Error(...)` then `ctx.Redirect(...)`
- Never call `models/` directly — go through `services/`
- Never write business logic — handler is a thin adapter

- [ ] **Step 6: Verify each rule**

- [ ] **Step 7: Write file & verify size**

```bash
wc -l routers/web/CLAUDE.md
```

- [ ] **Step 8: Optional commit**

```bash
git add routers/web/CLAUDE.md
git commit -m "docs(routers/web): add directory CLAUDE.md with handler patterns"
```

---

### Task 7: `routers/api/v1/CLAUDE.md`

**Files:**
- Create: `routers/api/v1/CLAUDE.md`
- Spec ref: "`routers/api/v1/CLAUDE.md`" scope

**Target:** <150 lines.

- [ ] **Step 1: Explore API handler structure**

```bash
ls routers/api/v1/ | head -20
cat routers/api/v1/repo/repo.go 2>/dev/null | head -80  # sample
grep -rn "func.*context.APIContext\|func.*APIContext" routers/api/v1/ | head -10
```

- [ ] **Step 2: Explore Swagger annotation**

```bash
grep -rn "// @Success\|// @Router\|// @Param" routers/api/v1/repo/ | head -15
cat routers/api/v1/swagger.go 2>/dev/null | head -50
ls .spectral.yaml 2>/dev/null
```

- [ ] **Step 3: Explore binding & response patterns**

```bash
grep -rn "ctx.Bind\|c.Bind\|ShouldBind" routers/api/v1/ | head -10
grep -rn "ctx.JSON\|ctx.Error\|APIError" routers/api/v1/ | head -15
grep -rn "ctx.Error\(.*Status\|ctx.NotFound\|ctx.ServerError" routers/api/v1/ | head -10
```

- [ ] **Step 4: Explore perm checks in API**

```bash
grep -rn "ctx.Repo.GitAccessMode\|RequireSignInView\|ctx.CheckPerm" routers/api/v1/ | head -10
```

- [ ] **Step 5: Draft rules**

Open with: "`routers/api/v1/` serves the REST API. Each handler has Swagger annotations, binds JSON/form, calls services, returns JSON."

Cover:
- Swagger annotation: required fields (`@Summary`, `@Router`, `@Param`, `@Success`, `@Failure`); run `make generate-swagger` after any API change
- Request binding via `ctx.Bind(...)` at handler entry; reject early on validation
- Response: `ctx.JSON(status, body)`; errors via `ctx.Error(status, msg)` or `ctx.NotFound()`
- Status code conventions: 400 (bad input), 401 (auth), 403 (perm), 404 (not found), 409 (conflict), 422 (validation), 500 (server)
- API version compatibility: add fields, don't remove/rename; v1 response shapes are contractual
- Perm checks at handler entry, never assume
- Never call `models/` directly — go through `services/`

- [ ] **Step 6: Verify each rule**

- [ ] **Step 7: Write file & verify size**

```bash
wc -l routers/api/v1/CLAUDE.md
```

- [ ] **Step 8: Optional commit**

```bash
git add routers/api/v1/CLAUDE.md
git commit -m "docs(routers/api): add directory CLAUDE.md with API and Swagger rules"
```

---

### Task 8: `web_src/CLAUDE.md`

**Files:**
- Create: `web_src/CLAUDE.md`
- Spec ref: "`web_src/CLAUDE.md`" scope

**Target:** <150 lines.

- [ ] **Step 1: Explore frontend structure**

```bash
ls web_src/
ls web_src/js/ | head -30
ls web_src/css/ | head -20
ls web_src/js/components/ 2>/dev/null | head  # Vue components
cat web_src/package.json 2>/dev/null | head -40
```

- [ ] **Step 2: Explore JS module patterns**

```bash
# Plain JS entry points
cat web_src/js/index.js 2>/dev/null | head -30  # or whatever the root is
ls web_src/js/features/ 2>/dev/null | head

# Vue component sample
find web_src -name "*.vue" | head -5
cat $(find web_src -name "*.vue" | head -1) | head -40
```

- [ ] **Step 3: Explore CSS conventions**

```bash
ls web_src/css/
cat web_src/css/_variables.css 2>/dev/null | head -20
grep -rn "GT-" web_src/css/ | head -10  # Gitea-specific class prefix?
```

- [ ] **Step 4: Explore async / CSRF patterns**

```bash
grep -rn "fetch(\|GET(\|POST(\|csrf" web_src/js/ | head -15
cat web_src/js/fetch.js 2>/dev/null | head -40  # or wherever the fetch wrapper is
```

- [ ] **Step 5: Draft rules**

Open with: "`web_src/` is the frontend source: plain JS modules, Vue components, CSS, Fomantic UI customizations."

Cover:
- JS organization: feature/page-based modules under `web_src/js/`; avoid scattered inline `<script>` in templates
- Vue components: single-file `.vue`; props/events/slots; use Vue only when interactivity warrants it (admin pages mostly Vue)
- CSS naming: use existing `GT-` prefix or BEM-ish; scope styles to component files where possible; avoid global pollution
- Fomantic UI: use bundled components; customization via overrides, not fork
- Async: use the project's fetch wrapper (handles CSRF + JSON); never raw `XMLHttpRequest`
- Build: webpack entries declared per page; check `web_src/index.js` and build config

- [ ] **Step 6: Verify each rule**

- [ ] **Step 7: Write file & verify size**

```bash
wc -l web_src/CLAUDE.md
```

- [ ] **Step 8: Optional commit**

```bash
git add web_src/CLAUDE.md
git commit -m "docs(web_src): add directory CLAUDE.md with frontend conventions"
```

---

### Task 9: `templates/CLAUDE.md`

**Files:**
- Create: `templates/CLAUDE.md`
- Spec ref: "`templates/CLAUDE.md`" scope

**Target:** <150 lines.

- [ ] **Step 1: Explore template structure**

```bash
ls templates/ | head -30
ls templates/base/ 2>/dev/null
ls templates/repo/ | head -15
cat templates/base/base.tmpl 2>/dev/null | head -50  # or whatever the root layout is
cat templates/base/head.tmpl 2>/dev/null | head -40
```

- [ ] **Step 2: Explore inheritance chain**

```bash
grep -rn "{{template\|{{define\|{{block" templates/base/ templates/repo/ | head -20
```

Document the inheritance: base → head/header/footer → page-specific.

- [ ] **Step 3: Explore helpers**

```bash
# Available template helpers
grep -rn "i18n.Tr\|\.Tr(\|AssetUrl\|\.Render\|\.Avatar\|\.DateFmt" templates/ | head -20
cat modules/templates/templates.go 2>/dev/null | head -60  # helper registration
ls modules/templates/ 2>/dev/null
```

- [ ] **Step 4: Explore i18n usage in templates**

```bash
grep -rn '{{.i18n.Tr\|{{.i18n\.Tr ' templates/ | head -15
```

- [ ] **Step 5: Draft rules**

Open with: "`templates/` holds Go HTML templates for the Web UI. Inherits base layout; uses helpers from `modules/templates/`; user-visible strings via i18n."

Cover:
- Inheritance: every page template starts with `{{template "base/head" .}}` and ends with `{{template "base/footer" .}}`
- Helpers: `.i18n.Tr "key"`, `.AssetUrl`, `.Render`, `.Avatar`, `.DateFmt` — see `modules/templates/` for full list
- i18n: all user-visible text via `{{.i18n.Tr "key"}}`; no hardcoded English
- Naming: template path mirrors route (e.g., `templates/repo/home.tmpl` for repo home page)
- Partials: shared fragments in `templates/base/` or per-feature subdirs; include via `{{template "name" .}}`

- [ ] **Step 6: Verify each rule**

- [ ] **Step 7: Write file & verify size**

```bash
wc -l templates/CLAUDE.md
```

- [ ] **Step 8: Optional commit**

```bash
git add templates/CLAUDE.md
git commit -m "docs(templates): add directory CLAUDE.md with template conventions"
```

---

### Task 10: Cross-File Consistency Check

**Files:**
- Read all: `design.md`, `*/CLAUDE.md` (8 files)

- [ ] **Step 1: Verify all 9 files exist**

```bash
ls design.md modules/CLAUDE.md models/CLAUDE.md services/CLAUDE.md cmd/CLAUDE.md routers/web/CLAUDE.md routers/api/v1/CLAUDE.md web_src/CLAUDE.md templates/CLAUDE.md
```

Expected: all 9 paths print without error.

- [ ] **Step 2: Check for rule conflicts**

Read `design.md` Section 1 (Architecture) and each `*/CLAUDE.md` opening statement. Confirm:
- Layered dependency direction stated consistently
- No directory file claims a responsibility that contradicts `design.md`
- "No `models/` import from routers" rule is consistent across `routers/web/CLAUDE.md`, `routers/api/v1/CLAUDE.md`, `services/CLAUDE.md`

- [ ] **Step 3: Verify rule evidence still holds**

Spot-check 5 random rules across all files. For each, run the grep that should confirm the pattern. If any rule now lacks evidence (e.g., code was refactored), update the doc.

- [ ] **Step 4: Update root `CLAUDE.md` index (optional)**

If the root `CLAUDE.md` should mention the new `design.md`, add a single line pointer under a "References" or "Further reading" section. Do NOT duplicate content — just point.

- [ ] **Step 5: Final commit (optional)**

```bash
git add CLAUDE.md
git commit -m "docs: link design.md from root CLAUDE.md"
```

---

## Self-Review Checklist (for the engineer executing this plan)

Before declaring complete:

- [ ] All 9 files created, each under its target line count
- [ ] Every rule in every file has been verified against grep evidence during that task
- [ ] No rule conflicts between `design.md` and directory `CLAUDE.md` files
- [ ] `modules/CLAUDE.md` documents the actual isolation violations (don't pretend they don't exist)
- [ ] `design.md` Section 0 includes a full index of the 8 directory `CLAUDE.md` files
- [ ] Root `CLAUDE.md` either untouched or only has a one-line pointer added
