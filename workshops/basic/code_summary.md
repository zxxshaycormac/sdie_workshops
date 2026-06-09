# Gitea Project Code Summary

## Overall Architecture

A **Go backend** + **JavaScript/CSS frontend** self-hosted Git service, following a classic layered architecture:

```
cmd/       -> Entry points (CLI commands)
routers/   -> HTTP handlers (API + Web UI)
services/  -> Business logic layer
models/    -> Data access & DB migrations
modules/   -> Shared utilities & libraries
```

---

## Backend (Go) — ~301,500 LOC

| Directory | Go LOC | Files | Role |
|-----------|--------|-------|------|
| **modules/** | 82,185 | 857 | Shared libraries (git ops, markup, settings, indexer...) |
| **routers/** | 64,068 | 381 | HTTP handlers |
| **models/** | 56,172 | 573 | Data models, DB migrations, fixtures |
| **services/** | 58,762 | 405 | Business logic |
| **tests/** | 31,554 | 207 | Integration & e2e tests |
| **cmd/** | 7,016 | 43 | CLI entry points |
| **build/** | 1,185 | 11 | Build tooling |
| **contrib/** | 562 | 3 | Contributions |
| **root files** | 52 | 2 | `main.go`, `build.go` |

### Top sub-packages by size

| Package | Sub-package | Go LOC | Purpose |
|---------|------------|--------|---------|
| modules | `git` | 12,653 | Git operations |
| modules | `emoji` | 3,598 | Emoji data |
| modules | `markup` | 7,044 | Markup rendering |
| modules | `setting` | 6,846 | Configuration |
| modules | `packages` | 6,132 | Package registry |
| modules | `indexer` | 4,883 | Search indexing |
| modules | `charset` | 2,335 | Charset detection |
| modules | `util` | 2,545 | General utilities |
| modules | `structs` | 2,445 | API type definitions |
| modules | `templates` | 2,440 | Template processing |
| modules | `log` | 1,739 | Logging |
| modules | `queue` | 1,764 | Job queue |
| routers | `web` | 33,122 | Web UI routes |
| routers | `api` | 27,727 | REST API routes |
| models | `migrations` | 9,604 | DB schema migrations |
| models | `issues` | 11,547 | Issue tracking models |
| models | `repo` | 5,293 | Repository models |
| models | `user` | 3,251 | User models |
| models | `asymkey` | 3,127 | Asymmetric key models |
| services | `webhook` | 5,842 | Webhook dispatch |
| services | `repository` | 7,409 | Repository operations |
| services | `pull` | 4,125 | Pull request logic |
| services | `auth` | 3,890 | Authentication |
| services | `context` | 3,341 | Request context |
| services | `migrations` | 8,084 | External repo migration |
| services | `mailer` | 2,456 | Email sending |
| services | `packages` | 2,425 | Package management |
| services | `gitdiff` | 2,600 | Diff processing |

---

## Frontend — ~48,400 LOC

| Directory | LOC | Files | Details |
|-----------|-----|-------|---------|
| `web_src/js` | 12,517 | 157 | JavaScript application code |
| `web_src/css` | 12,437 | 81 | CSS stylesheets |
| `web_src/fomantic` | 22,916 | 6 | Fomantic UI (Semantic UI fork) customizations |
| `web_src/svg` | 488 | 57 | SVG icons |
| Vue components | 3,364 | 16 | `.vue` single-file components |

---

## Other Notable Areas

| Directory | LOC | Files | Content |
|-----------|-----|-------|---------|
| **options/** | 72,975 | 28 | INI locale/translation files (i18n) |
| **docs/** | 20,503 | 204 | Markdown documentation |
| **tests/integration/** | 31,169 | 203 | Go integration tests |
| **templates/** | — | — | Go HTML templates (not counted by cloc) |
| **assets/** | 1,229 | 3 | JSON assets |
| **docker/** | 309 | 14 | Dockerfile & scripts |
| **public/** | 388 | 383 | Static SVG assets |
| **snap/** | 75 | 1 | Snap packaging |

---

## Key Takeaways

- **~301K Go LOC** makes up the core backend, well-layered (routers -> services -> models -> modules)
- **`modules/git`** (12.6K LOC) is the largest single package — the Git integration is the heart of the project
- **`models/migrations`** (9.6K LOC across 263 files) shows extensive DB schema evolution
- **`routers/api`** (27.7K LOC) and **`routers/web`** (33.1K LOC) are roughly 1:2, indicating a large web UI + REST API surface
- **Frontend** is ~48K LOC with heavy use of Fomantic UI; relatively modest custom JS/CSS
- **72K LOC of translations** in `options/` shows strong i18n support
- **Test code** (~31.5K LOC integration tests) is focused on integration rather than unit testing
