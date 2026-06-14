# Design Conventions Documentation — Design Spec

**Date**: 2026-06-13
**Branch**: training/base_line
**Status**: Approved (pending implementation)
**Source**: `/superpowers:brainstorming` session

## Goal

Discover and document Gitea's design conventions and standards from existing code, as an **AI assistant reference document**. The output tells Claude (or any AI assistant) what conventions to follow when modifying Gitea code, beyond what `CLAUDE.md` already covers.

## Background

The Gitea codebase (~301K Go LOC + ~48K frontend LOC) has implicit design conventions spread across `cmd/`, `routers/`, `services/`, `models/`, `modules/`, `web_src/`, and `templates/`. The root `CLAUDE.md` captures build/test commands and a high-level architecture summary, but not the deeper conventions an AI assistant needs when modifying code: error wrapping patterns, XORM session propagation, Swagger annotation requirements, i18n Tr() calls, template helper signatures, etc.

These are currently discoverable only by reading adjacent code — slow and inconsistent. This spec defines a documentation structure that captures them in a way that is:

- **Progressive**: AI receives relevant rules contextually based on working directory
- **Principle-first**: Rules and rationale, not extensive code examples (AI can read code itself)
- **Discoverable**: Clear division between cross-cutting and directory-specific principles

## File Structure

### Deliverables (9 files)

1. `design.md` (repo root) — main document, topic-organized, cross-cutting principles
2. `modules/CLAUDE.md`
3. `models/CLAUDE.md`
4. `services/CLAUDE.md`
5. `cmd/CLAUDE.md`
6. `routers/web/CLAUDE.md`
7. `routers/api/v1/CLAUDE.md`
8. `web_src/CLAUDE.md`
9. `templates/CLAUDE.md`

### Division of Concerns

| Document | Scope |
|----------|-------|
| `design.md` | Principles that apply regardless of which directory AI is working in |
| `<dir>/CLAUDE.md` | Principles specific to working inside that directory |

Rationale: Claude Code auto-loads `CLAUDE.md` from the current working directory upward. By placing directory-specific rules in each directory's `CLAUDE.md`, AI receives progressively more specific guidance as it navigates deeper. Cross-cutting rules in root `design.md` are referenced from anywhere.

## `design.md` Outline (Topic-First / Plan B)

8 sections. Uniform internal structure per section:

```
## N. Topic Name

**Principle**: <one-sentence statement>
**Why**: <1-2 sentence rationale>
**Exceptions**: <if any>
```

Sections:

0. **Document positioning** — audience (AI assistants), how to use this doc with the directory `CLAUDE.md` index
1. **Architecture principles** — layering (`cmd → routers → services → models → modules`), `modules/` isolation rule (no imports from `models/`/`services/`/`routers/`), no circular deps, dependency flow downward only
2. **Naming conventions** — package / type / function / file / DB column naming rules
3. **Error handling** — wrapping with `fmt.Errorf("...: %w", err)`, sentinel errors (`ErrNotFound` style), return-vs-log boundary, never swallow errors
4. **Context propagation** — `context.Context` as first param on all I/O functions, cancellation propagation, `DefaultContext` usage, timeouts at boundary
5. **Logging** — level selection (Debug/Info/Warn/Error), structured fields, no sensitive data (tokens, passwords)
6. **Internationalization (i18n)** — all user-visible strings via `Tr()`, no hardcoded English in user-facing paths
7. **Testing principles** — `unittest` framework for DB-backed tests, integration vs unit test placement, never skip / `.Skip()` without tracking issue
8. **Security general principles** — input validation at boundary, permission checks (`PermChecker`), never trust user input

## Directory `CLAUDE.md` Scopes

Each file opens with one sentence stating the directory's responsibility, then a list of rules. Each rule: one sentence + optional rationale.

### `modules/CLAUDE.md`
- Isolation rule: never import `models/`, `services/`, `routers/`
- Sub-package organization: single responsibility per sub-package, avoid catch-all `util/`
- Public API design: hide internal types, don't leak `context.Context` from internal APIs
- Test requirement: no DB dependency, runnable in isolation

### `models/CLAUDE.md`
- XORM model definition rules: field tags, `Bean` registration via `init()`, table naming
- `db.Session` / `DefaultContext` acquisition and propagation patterns
- Migration rules: sequential version numbers, published migrations are immutable, up-only (no down migrations in production)
- `unittest` framework usage: fixtures, test DB preparation, `TestMain` setup
- Boundary: data integrity + CRUD only, no business logic
- Transactions are controlled by callers in `services/`, not within `models/`

### `services/CLAUDE.md`
- Orchestration role: call `models/`, called by `routers/`
- Transaction boundary: wrap multi-model operations here, do not leak to `routers/`
- Async work via `queue` package, never raw `go` keyword for fire-and-forget
- Webhook / notification trigger points belong here, not in `models/`
- Boundary: do not render templates, do not bind HTTP forms

### `cmd/CLAUDE.md`
- Command registration structure (urfave/cli flags and subcommand organization)
- Initialization order: aligns with `routers.InitWebInstalled()` sequence
- `context.Context` originates here and propagates downward
- No business logic — only dispatch to `services/`

### `routers/web/CLAUDE.md`
- Handler signature: uses `*context.Context` from `services/context`
- Form binding and validation entry points
- Permission middleware (`MustRepo`, `RequireSignInView`, etc.)
- Template rendering convention (`ctx.Render`, data passing)
- Flash message and redirect patterns
- Never call `models/` directly — go through `services/`

### `routers/api/v1/CLAUDE.md`
- Swagger annotation: required fields; `make generate-swagger` mandatory after API changes
- Request binding (`Bind`/`Decode`), uniform error format
- Response patterns (`JSONError`, status code conventions: 400/403/404/409/422/500)
- API version compatibility: do not break 1.x response shapes; add fields, don't remove/rename
- Permission checks at handler entry
- Never call `models/` directly — go through `services/`

### `web_src/CLAUDE.md`
- JS organization: modular per page/feature, avoid scattered jQuery
- Vue component conventions: props/events/slots, single-file components, when to use Vue vs plain JS
- CSS naming and scoping: avoid global pollution, BEM-ish conventions actually used
- Fomantic UI: use existing components, customize via overrides not forks
- Async requests: `fetch` + CSRF token handling
- Build entry points: webpack config locations

### `templates/CLAUDE.md`
- Template inheritance chain: `base` → `layout` → `page`
- Helper invocations: `.i18n.Tr`, `.Render`, `.AssetUrl`, etc.
- User-visible strings must go through `i18n`
- Template naming and organization (by feature directory)
- Partial template reuse conventions

## Writing Principles (How to Extract During Implementation)

When filling in these documents:

1. **Discover, don't invent** — every rule must be backed by observable code patterns. If only 1-2 examples exist, mark as "emerging pattern" not a rule.
2. **Cite frequency** — note whether a pattern is universal, common, or occasional. This tells AI how strictly to enforce.
3. **Note known violations** — if a rule has exceptions in the codebase, mention them so AI doesn't enforce blindly.
4. **Avoid redundancy with root `CLAUDE.md`** — root covers build/test/architecture high-level. `design.md` goes deeper on conventions; directory `CLAUDE.md` goes deeper on directory-specific patterns.
5. **No code dumps** — principles only. AI reads code itself; the doc tells AI what to look for and what to follow.
6. **Use exploration tools** — `grep` / Explore agent to verify a pattern is wide before elevating to a rule. Do not extrapolate from a single file.

## Out of Scope

- Refactoring existing code to match conventions (documentation only)
- Frontend deep style guide (covered at high level only; not a Fomantic tutorial)
- API specification (Swagger is authoritative; doc references but does not duplicate)
- Database schema documentation (covered by `models/` + migrations directly)
- Contributor workflow / PR process (covered by `CONTRIBUTING.md`)
- Rewriting or extending the root `CLAUDE.md` (root stays as-is unless rule conflicts surface)

## Success Criteria

- AI assistant reading `design.md` + a directory's `CLAUDE.md` can correctly identify conventions for code in that directory without further exploration
- Each rule is backed by observable code patterns (discoverable via grep)
- No rule conflicts between `design.md` and directory `CLAUDE.md` (cross-cutting vs specific clearly separated)
- Documents are concise: `design.md` under ~400 lines; each directory `CLAUDE.md` under ~150 lines
- Every deliverable file committed to git
