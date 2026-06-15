# Repository Tags Page Search Filter — Design Spec

**Date**: 2026-06-15
**Branch**: training/base_line
**Status**: Approved (pending implementation)
**Source**: `/superpowers:brainstorming` session

## Goal

Add keyword search to the repository Tags page (`/{owner}/{repo}/tags`) so users can find a specific tag by name when a repository has many tags. The UX matches the existing Branches page search pattern.

## Background

The Tags page currently supports only `page` and `limit` query parameters — there is no way to filter by name. The Branches page (`routers/web/repo/branch.go`) and Commits page (`routers/web/repo/commit.go`) both already support a `q` query parameter for name-based filtering. The Tags page is the inconsistent outlier.

The Tags page handler `TagsList` (`routers/web/repo/release.go:204-251`) calls `models/repo/release.go:FindReleases` via `FindReleasesOptions`. The underlying `release` table already has a `lower_tag_name` column (maintained by the model on every tag insert/update) but no index.

## Approach

Follow the established Branches search pattern exactly, with one improvement: use the existing `lower_tag_name` column for case-insensitive matching instead of Branches' case-sensitive `builder.Like{"name", ...}`.

Rejected alternatives:
- **Client-side JS filtering** — incompatible with server-side pagination; only filters the current page.
- **Shared search helper abstraction** — only two call sites with meaningfully different shapes; premature.

## Architecture & Data Flow

```
Browser GET /{owner}/{repo}/tags?q=v1.2&page=2
        │
        ▼
routers/web/repo/release.go: TagsList(ctx)
        │  reads q := ctx.FormString("q")
        │  builds FindReleasesOptions{..., Keyword: q}
        ▼
models/repo/release.go: FindReleases(ctx, opts)
        │  ToConds() adds: builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)}
        ▼
XORM → SELECT ... FROM release WHERE ... AND lower_tag_name LIKE '%v1.2%'
        │
        ▼
Template tag/list.tmpl renders list + search form (shared/search/combo)
        │
        ▼
Pagination via SetDefaultParams(ctx) carries q through page links
```

Key points:
- **Unidirectional**: browser → handler → model → SQL → template. No JS, no AJAX, no new API.
- **Reuses existing column**: `release.lower_tag_name` is already maintained by the model; LIKE on this column gives case-insensitive matching with no `LOWER()` function call (consistent behavior across PostgreSQL/MySQL/SQLite).
- **Pagination + search coexist**: `SetDefaultParams(ctx)` automatically threads `q` through pagination links, so paginating search results does not drop the keyword.
- **Empty results**: template falls back to the existing "no tags" empty state; no new branch needed.

The single difference from the Branches pattern: Branches uses `builder.Like{"name", ...}` (case-sensitive); Tags uses the `lower_tag_name` column (case-insensitive). This is a free upgrade because the column already exists.

## Code Changes (5 files, ~30 lines)

### 1. `models/repo/release.go`

Add a `Keyword string` field to the `FindReleasesOptions` struct (lines 229-238).

In `ToConds()` (lines 240-266), append a new clause after the existing conditions (after the `HasSha1` block at lines 258-264, before the `return cond` at line 265):

```go
if opts.Keyword != "" {
    cond = cond.And(builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)})
}
```

The `"strings"` package is already imported (line 14); no import change needed.

### 2. `routers/web/repo/release.go`

In `TagsList` (lines 204-251), add immediately after the `listOptions` block closes at line 224 (before the `opts := repo_model.FindReleasesOptions{` literal at line 226):

```go
keyword := ctx.FormString("q")
```

In the `opts` literal (lines 226-234), add a `Keyword: keyword,` field (e.g., after the `RepoID` line at 233).

Right after `ctx.Data["Releases"] = releases` (line 242), add:

```go
ctx.Data["Keyword"] = keyword
```

Pagination at lines 244-247 already calls `pager.SetDefaultParams(ctx)`, which automatically carries `q` through page links — no change needed there.

### 3. `templates/repo/tag/list.tmpl`

The current template opens with `release_tag_header` at line 6 and wraps the tag table in `{{if .Releases}}` at line 7. Insert the search form between lines 6 and 7 so it renders even when the search returns zero results (otherwise users could not retry with a different keyword). Pattern copied from `templates/repo/branch/list.tmpl:76-80`:

```html
<div class="ui attached segment">
    <form class="ignore-dirty" method="get">
        {{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}
    </form>
</div>
```

### 4. `options/locale/locale_en-US.ini`

The `[search]` section starts at line 162. `branch_kind` is at line 178 and `commit_kind` at line 179. Add at line 180:

```ini
tag_kind = Search tags...
```

Other 40+ locale files are translated by the community; do not machine-translate.

### 5. `tests/integration/repo_tag_test.go` (append to existing file)

The file already exists. Add two integration tests:
- **Hit**: GET `/{owner}/{repo}/tags?q=<known-tag-name-fragment>` returns HTTP 200 and the response body contains the matching tag row.
- **Miss**: GET `/{owner}/{repo}/tags?q=<nonexistent-string>` returns HTTP 200 and the response body does not contain any tag row.

Use the existing fixture repository (e.g. `repo1`, which already has tags like `v1.1`). No new fixtures needed.

## Files Explicitly NOT Touched

- `routers/web/repo/release.go` Releases-related handlers — scope is Tags only.
- `routers/api/v1/repo/tag.go` — REST API unchanged.
- `services/repository/` — no new service-layer code; `TagsList` calls the model directly.
- DB migrations — `lower_tag_name` column already exists.

## Edge Cases

| Scenario | Behavior |
|---|---|
| `q` empty or missing | `Keyword == ""` → `ToConds()` skips LIKE → identical to current behavior (backward compatible) |
| `q` contains SQL wildcards (`%`, `_`) | XORM's `builder.Like` escapes them as parameter values; searching `100%` does not match all tags |
| `q` contains whitespace | `FormString` preserves whitespace, passed through to LIKE (matches Branches behavior; no trim) |
| Mixed-case query (e.g. `V1.2`) | Handler passes raw keyword to model; model applies `strings.ToLower` before LIKE on `lower_tag_name` → `V1.2` and `v1.2` both match `v1.2.3` |
| Search results exceed `limit` | Pagination applies; `SetDefaultParams(ctx)` carries `q` through page links |
| Zero search results | Template renders existing "no tags" empty state; no new branch |
| Releases page shares `FindReleasesOptions` | Releases handler does not read `q`, so `Keyword` is always empty string → `ToConds()` adds no condition → Releases behavior unchanged |

## Security

- **SQL injection**: `builder.Like` is parameterized; user input is never concatenated into the SQL string. Compliant with `design.md` §8.
- **XSS**: keyword renders via `{{.Keyword}}` in the template; Go's `html/template` auto-escapes. The `shared/search/combo` `Value` field goes through the standard escape path.
- **Authorization**: search reuses the existing `RepoAssignment` + `MustBeRepo` middleware chain on `TagsList`; no permission check is bypassed.

## Risks & Trade-offs

1. **`lower_tag_name` has no index**: `LIKE '%keyword%'` is a full table scan on repositories with many tags (10k+). The Branches table has the same characteristic, so this matches the project's accepted baseline. Adding an index is **out of scope** for this change — it is an optimization, not a feature. If performance becomes an issue in the future, add `INDEX lower_tag_name` via a separate migration.

2. **i18n coverage**: only `locale_en-US.ini` is updated. Community translation workflow handles the other locales per `design.md` §6.

3. **Test fixture dependency**: integration test relies on the fixture repository having tags. The default `repo1` fixture already has `v1.1` etc., so no new fixtures are required.

## Out of Scope (YAGNI)

- AJAX instant search (GET form submit is enough; matches Branches).
- Search-result highlighting (Branches does not have it).
- Search by commit message / author (out of scope for "search tag names").
- REST API changes (user confirmed Tags-only scope).
- DB index / migration (see Risks #1).
- Refactoring of shared search helpers.
