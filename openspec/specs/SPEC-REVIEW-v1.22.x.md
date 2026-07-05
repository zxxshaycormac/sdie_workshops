# OpenSpec Baseline Review — Gitea v1.22.x vs Codebase

**Method:** 9 parallel agents, one per MECE domain, each verified its `spec.md` against the actual Gitea v1.22.x source with `file:line` evidence. This document synthesizes their findings. Detailed per-domain reports are preserved in the agent transcripts.

---

## Executive Summary

The baseline is **structurally sound** (the 9-domain decomposition maps cleanly onto the codebase) and **roughly 75–85% accurate** on documented behaviors. However, every domain contains at least a few **factual inaccuracies** and several contain **fabricated features** (describing behavior the code does not implement). The two most common defect classes:

1. **Invented features** — the spec describes capabilities that don't exist in v1.22.x (e.g. `workflow_dispatch`, GPG keyserver import, Jupyter/PDF rendering, API rate limiting, time estimates on issues, org-level default repo permission, Meilisearch code search, organization sitemaps).
2. **Wrong/misnamed config keys** — the Configuration Reference sections are the least reliable part of the spec across all domains (fabricated keys, wrong names, wrong sections, deprecated keys presented as current).

| # | Domain | Accuracy | Critical Issues |
|---|--------|----------|-----------------|
| 01 | Identity & Access | ~70% | 2FA uses MD5 (not PBKDF2); GPG keyserver import doesn't exist; OpenID misclassified as auth source type; heatmap is 15-min buckets not UTC days |
| 02 | Collaboration | ~85% | Time-tracking fabricates a "description" field + a "time estimate" feature; content history stores full text not hashes; `.github/CODEOWNERS` not consulted; project columns auto-move cards instead of blocking |
| 03 | Organization | ~80% | "Badges" is a user-level feature, not org-level; invite expiration, default repo permission, org fork permissions don't exist; missing "Releases" unit type |
| 04 | Code Management | ~75% | Repo visibility model misstated (limited is user/org attr, not repo); `[git]` config block has 7+ wrong/fabricated keys; missing entire `[mirror]` and `[migrations]` sections |
| 05 | CI/CD & Automation | ~75% | `workflow_dispatch`/`repository_dispatch` unsupported; `[task]` section deprecated; head SHA not persisted on auto-merge record; merge styles list incomplete |
| 06 | Packages & Artifacts | **~95%** | 21 package types verified exactly; minor: attachment default is 2MB not 4MB, missing scope returns 401 not 403, `actions_artifacts` plural not singular |
| 07 | Administration & Ops | ~70% | Health endpoint is `/api/healthz` not `/-/healthz`; pprof runs on separate `localhost:6060`; notice types are Repository/Task not warning/error/info; CLI list wrong; ~17 doctor checks undocumented; `--yes` flag doesn't exist |
| 08 | Search & Discovery | ~65% | Meilisearch is issue-indexer-only (no code search); org sitemaps don't exist; config keys misnamed; fuzzy/exact modes belong to code search not repo search; license filter doesn't exist |
| 09 | Integration & Extension | ~70% | Fabricated Jupyter/PDF renderers + API rate limiting (429); themes undercounted (5 not 3, default is `gitea-auto`); OAuth2 providers undercounted (16 not 3); OAuth2 *provider* side undocumented; `RENDER_CONTENT_MODE` value is `iframe` not `image` |

---

## Cross-Cutting Defect Patterns

### Pattern 1: Fabricated Features (highest severity)
These describe behavior the code does NOT implement. Each will mislead any delta spec or implementer.

| Domain | Spec Claim | Reality |
|--------|-----------|---------|
| 05 CI/CD | `workflow_dispatch`, `repository_dispatch` triggers | Unsupported (docs explicitly say "ignored by Gitea Actions") |
| 02 Collab | TrackedTime has a `description` field; issue time estimates | No such field/feature exists |
| 01 Identity | GPG keyserver import by long ID | No fetch logic; only pasted armored keys |
| 09 Integration | Jupyter `.ipynb` rendering, PDF viewer | No renderers registered |
| 09 Integration | API rate limiting → HTTP 429 | No rate-limit module |
| 09 Integration | SSE connection cap → 429; camo bypass list | Neither exists |
| 03 Org | Configurable invite expiration; org default repo permission; org fork permissions | None exist |
| 08 Search | Meilisearch for code search; organization sitemaps; license filter | None exist |
| 07 Admin | `restore-cert` CLI command; `--yes` doctor flag | Neither exists |
| 01 Identity | "session storage `redis-cluster`" | Not a supported provider |

### Pattern 2: Config Reference Unreliable
The Configuration Reference is the **weakest section in every domain**. Recurring problems:
- **Fabricated keys** (e.g. `CONCURRENT_CREATION_LIMIT`, `DEFAULT_APPLICATION_VISIBILITY`, `MAX_LIFE_TIME`, `CLIENT_SECRET_LENGTH`, `VERBOSE_PUSH_TIMEOUT`, `RENDERED_MAX_FILE_SIZE`).
- **Wrong names** (e.g. `AUTO_WATCH_REPOS` → `AUTO_WATCH_NEW_REPOS`; `MAX_LIFE_TIME` → `SESSION_LIFE_TIME`; `EXPLORE_DEFAULT_SORT` → `EXPLORE_PAGING_DEFAULT_SORT`; `DISABLE_DIFF_DIFF_TEXT` → `DISABLE_DIFF_HIGHLIGHT`).
- **Wrong sections** (e.g. CSRF cookie is in `[security]` not `[session]`; `ALLOW_LOCAL_NETWORKS` is in `[migrations]` not `[mirror]`).
- **Deprecated keys presented as current** (e.g. `[task]` deprecated since v1.19 for `[queue.task]`; `[git.reflog]` remapped to `[git.config]`).
- **Entire sections missing** (`[mirror]`, `[migrations]`, `[webhook]`, `[cron.<name>]`, `[oauth2_client]`, `[federation]`, `[queue]`, `[metrics]`).

**Recommendation:** Regenerate every Configuration Reference block directly from `modules/setting/*.go` structs rather than hand-authoring.

### Pattern 3: Mischaracterized / Mis-scoped Features
- **Badges** (03 Org) — documented as org-level; actually a user-account feature.
- **OpenID** (01 Identity) — documented as an `auth.Type` source; actually a dedicated consumer flow.
- **Repo visibility** (04 Code) — "limited" is a user/org attribute; repos are public/private only.
- **Fuzzy/exact search modes** (08 Search) — belong to code search, not repo search.
- **Admin blocking** (01/03) — block *creation* is not denied for admins; only enforcement is skipped.

### Pattern 4: Under-counted Inventories
Several spec sections give incomplete lists where the code is exhaustive:
- **OAuth2 client providers**: spec implies 3 ("GitHub, GitLab, Google"); code has **16**.
- **Themes**: spec says 3; code ships **5** (incl. 2 colorblind variants).
- **Markup renderers**: spec misses Console and Asciicast (`.cast`).
- **Merge styles**: spec lists 3; code defines **7**.
- **Doctor checks**: spec lists ~6 categories; code has **~27**.
- **Activity periods**: spec lists 3; code supports **7**.
- **Webhook events**: spec's list is partial vs the full enum.
- **CLI subcommands**: spec invents 3, omits 6 real ones.
- **Dump archive formats**: spec lists 8; code has **9** (omits `tar.sz`).

---

## Per-Domain Highlights

### 01 — Identity & Access
- **2FA secret encryption is MD5-derived**, not PBKDF2 (`models/auth/twofactor.go:95-98`). Flag as a possible security weakness.
- **Heatmap aggregates in 15-minute buckets**, not UTC days (`user_heatmap.go:41`), to permit client-side timezone shifting.
- Reserved-usernames list is stale (missing `.well-known`, `avatars`, `commits`, `manifest.json`, `milestones`, `repo-avatars`).
- Undocumented config: `[oauth2_client]` (account linking, auto-registration), `[federation]`, `REGISTER_MANUAL_CONFIRM`.

### 02 — Collaboration
- **Time-tracking (COLL-09) needs a rewrite** — half its requirements describe non-existent features (description, estimates).
- **Content history stores full text**, not hashes; capped at 20 revisions.
- CODEOWNERS searched in `CODEOWNERS`, `docs/CODEOWNERS`, `.gitea/CODEOWNERS` only (no `.github/`).
- Project column deletion auto-moves cards to default column (doesn't block).
- Auto-merge does NOT verify a check/approval rule exists before scheduling.

### 03 — Organization
- **Badges section should be removed/relocated** — it's a user-account feature (`models/user/badge.go`).
- Three documented capabilities don't exist: invite expiration, default repo permission, org fork permissions.
- Missing "Releases" from the unit-permission list (10 types, not 9).
- Undocumented real features: org Actions secrets/runners, package cleanup rules + Cargo index, per-unit team access modes.

### 04 — Code Management
- **Visibility model is wrong** — repos are `IsPrivate bool`; three-level visibility is a user/org attribute.
- `[git]` config block: 4 wrong names, 2 fabricated keys, 1 fabricated feature key.
- Missing entire `[mirror]` and `[migrations]` sections.
- Missing features: `.bundle` archive format, Gitea as migration source, wiki-as-git-repo, tag browsing, repo adoption, `DISABLE_STARS`.

### 05 — CI/CD & Automation
- **Remove `workflow_dispatch`/`repository_dispatch`** — most user-misleading inaccuracy.
- `[task]` section deprecated; replace with `[queue.<name>]` keys.
- Head SHA tracked in queue message, not `pull_auto_merge` table; stale entries silently skipped, not "aborted with notification."
- Missing: `pull_request_target`, AGit auto-merge, fork-PR approval state machine, `GITEA_`/`GITHUB_`/`CI` name-prefix forbiddance, `SKIP_WORKFLOW_STRINGS`.

### 06 — Packages & Artifacts (most accurate domain)
- **21 package types verified exactly** (`models/packages/package.go:33-53`).
- Fix attachment default 4MB → **2MB** (`attachment.go:16`).
- Missing read:package scope returns **401**, not 403.
- Storage type name `actions_artifacts` (plural).
- Undocumented: public-only token scope, multi-hasher (MD5+SHA1+SHA256+SHA512), `_catalog` endpoint, `LFS_HTTP_AUTH_EXPIRY`.

### 07 — Administration & Ops
- Health endpoint: `/-/healthz` → **`/api/healthz`**.
- No `/api/v1/metrics` — only `/metrics`.
- **pprof runs on a separate `localhost:6060` server**, not main routes. Config key is `[server].ENABLE_PPROF`.
- "smtp" log mode doesn't exist (only console/file/conn).
- Notice types are `Repository`/`Task`; notices are **deleted**, not "resolved."
- `--yes` doctor flag doesn't exist; dump flags are `--skip-lfs-data` etc. (with `-data` suffix).
- CLI: `restore-cert` doesn't exist; `convert`/`checks` are `doctor` subcommands; missing `embedded`, `dump-repository`, `restore-repository`, `actions`, `cert`, `docs`.

### 08 — Search & Discovery
- **Meilisearch is issue-indexer-only** — code search is Bleve + Elasticsearch.
- **Organization sitemaps don't exist.**
- Config keys wrong: `EXPLORE_DEFAULT_SORT` → `EXPLORE_PAGING_DEFAULT_SORT`; `EXPLORE_DISABLE_USERS_PAGE` → `[service].DISABLE_USERS_PAGE` (and it **redirects**, doesn't 404).
- Fuzzy algorithm is **Levenshtein**, not Damerau-Levenshtein.
- Autocomplete (`/user/search_candidates`) is single-mode — no mention/assignee distinction, no write-access scoping, no 2-char minimum.
- Issue indexer has an undocumented **`db`** backend.

### 09 — Integration & Extension
- **Remove fabricated renderers** (Jupyter, PDF) and **API rate limiting (429)**.
- **5 themes** (not 3); default is `gitea-auto` (not `gitea` light).
- **16 OAuth2 client providers** (not 3) — list them all.
- **OAuth2 Provider side** (Gitea as authorization server) is entirely undocumented despite being a major feature.
- Config: `RENDER_CONTENT_MODE` value is `iframe` not `image`; `[ssh.minimum_key_sizes]` defaults wrong (ed25519=256, rsa=3071, no DSA).
- Add Asciicast + Console renderers, Mermaid cap, Federation/ActivityPub.

---

## Top Recommendations (Priority Order)

1. **Delete all fabricated features** (Pattern 1 table) — these are the highest-risk items for any delta spec author.
2. **Regenerate every Configuration Reference from `modules/setting/*.go`** — automate this; hand-authored config blocks are unreliable across all 9 domains.
3. **Fix the most user-misleading claims**: 2FA MD5 encryption, `workflow_dispatch`, repo visibility model, health endpoint path, Meilisearch code search.
4. **Complete the inventories** (OAuth2 providers, themes, renderers, merge styles, doctor checks, CLI commands, webhook events).
5. **Re-scope mischaracterized features** (Badges → user domain; OpenID → consumer flow; fuzzy/exact → code search).
6. **Add the OAuth2 Provider** as a first-class feature (currently absent from domain 09 despite being a core capability).
7. **Add missing config sections**: `[mirror]`, `[migrations]`, `[webhook]`, `[cron.<name>]`, `[oauth2_client]`, `[federation]`, `[queue]`, `[metrics]`.

---

*Generated by parallel codebase verification. Per-domain detailed reports with full `file:line` evidence are available in the agent transcripts.*
