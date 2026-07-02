# Tags Page Search — Design

**Date:** 2026-07-02
**Status:** Approved
**Scope:** Add keyword search/filter to the repository Tags page (`/{org}/{repo}/tags`), consistent with the existing Branches and Commits search UX.

## Problem

The repository Tags page has no search or filter. When a repo has many tags, users cannot quickly find a specific one. The Branches page (`/{org}/{repo}/branches`) and Commits page already support keyword search via the `q` query parameter. The Tags page lacks this, creating an inconsistent experience.

## Solution

Replicate the Branches search pattern exactly: a plain HTML GET form that submits `?q=keyword`, filtered server-side by a DB-level `LIKE` condition, rendered through the shared `shared/search/combo` template partial.

## Design Decision: Case-Insensitive Matching

The `release` model has a `lower_tag_name` column (lowercased copy of `tag_name`). Search filters on `lower_tag_name` using a lowercased keyword, giving case-insensitive matching: `v1` finds `V1.0.0`. This is more forgiving than the Branches page (which uses case-sensitive `LIKE` on `name`) and is trivial to implement since the column already exists.

## Architecture & Data Flow

```
User types in search box
  → browser native GET form submits ?q=keyword (no JavaScript)
  → TagsList handler reads ctx.FormString("q")
  → sets opts.Keyword on FindReleasesOptions
  → db.FindAndCount applies builder.Like{"lower_tag_name", lower(kw)} in ToConds()
  → filtered, paginated results returned with accurate filtered count
  → ctx.Data["Keyword"] = kw (search box shows current term)
  → pager.SetDefaultParams(ctx) carries q across pagination links
```

No new routes, no JavaScript, no DB migration. The `GET /tags` route already accepts arbitrary query params.

## Component Changes

### 1. `models/repo/release.go` — Add keyword to options + filter

Add `Keyword string` field to `FindReleasesOptions` (struct at line 229-238).

In `ToConds()` (line 240-266), add after the existing conditions:

```go
if opts.Keyword != "" {
    cond = cond.And(builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)})
}
```

This mirrors `FindBranchOptions.Keyword` + the `builder.Like` idiom in `models/git/branch_list.go:88-106`. Requires adding `"strings"` and `"xorm.io/builder"` imports if not already present (builder is already used in this file).

### 2. `routers/web/repo/release.go` — Read `q`, apply filter, fix pagination count

In `TagsList` (lines 204-251):

- Add `kw := ctx.FormString("q")` (mirrors `branch.go:55`)
- Set `Keyword: kw` on the `FindReleasesOptions` struct
- Switch from `db.Find[repo_model.Release](ctx, opts)` to `db.FindAndCount[repo_model.Release](ctx, opts)` — returns `(releases, filteredCount, err)`
- Use `filteredCount` as the pagination total (when no keyword is set, this equals the total tag count, so behavior is unchanged in the no-search case; the old `ctx.Data["NumTags"]` read is no longer needed)
- Set `ctx.Data["Keyword"] = kw` (mirrors `branch.go:79`)
- `pager.SetDefaultParams(ctx)` already carries `q` across page links — no change needed

**Why the pagination fix matters:** `ctx.Data["NumTags"]` is set by middleware as the total count of all tags in the repo (no keyword). Using it for pagination when filtering would show wrong page counts (e.g., "page 1 of 10" when only 2 tags match). Switching to `db.FindAndCount` gives an accurate count in both cases — with keyword (filtered) and without (total). Branches solves this via `LoadBranches` returning the filtered count; we solve it via `db.FindAndCount`.

### 3. `templates/repo/tag/list.tmpl` — Insert search form

Insert between the `<h4>` header (lines 8-12) and the table div (line 14), matching `branch/list.tmpl:76-80`:

```go
<div class="ui attached segment">
    <form class="ignore-dirty" method="get">
        {{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}
    </form>
</div>
```

This is a native HTML GET form — no JavaScript. The `shared/search/combo` partial renders the `name="q"` input (defined in `templates/shared/search/input.tmpl`) plus the search button. The `ignore-dirty` class prevents the global unsaved-changes guard from intercepting navigation.

### 4. `options/locale/locale_en-US.ini` — Add placeholder key

Add `search.tag_kind = Search tags` (mirrors existing `search.branch_kind` and `search.commit_kind`).

## Error Handling

No new error paths. The keyword filter is a pure DB `LIKE` condition. If `db.FindAndCount` fails, the existing `ctx.ServerError("...", err)` pattern applies unchanged. An empty keyword (`?q=`) is treated as "no filter" — `builder.Like` is skipped when `Keyword == ""`, identical to how branches behave.

## Testing

### Unit test — `models/repo/release_test.go`

Add a test verifying that `FindReleasesOptions{Keyword: "..."}`:
- Returns only releases whose `lower_tag_name` matches the lowercased keyword
- Returns an empty set when the keyword matches nothing
- Returns all matching releases (regression: no keyword = existing behavior)

Use existing fixtures in `models/fixtures/release.yml`. Assert via `db.Find[repo_model.Release]` with the keyword set.

### Integration test — `tests/integration/`

Add a test that:
- Loads the tags page with `?q=<known-tag-substring>` and asserts only matching tags appear in the HTML response
- Loads with `?q=<nonexistent>` and asserts the table body is empty but the page renders (200 OK)
- Loads with no `q` and asserts all fixture tags appear (regression check)
- Verifies the search input box is present and its `value` attribute reflects the current `q`

No frontend/E2E test needed — the search is a native GET form with no JavaScript, so HTML integration tests are sufficient.

## Definition of Done

- [ ] `make lint` passes
- [ ] `make test-backend` passes (including new unit + integration tests)
- [ ] Search box renders on the tags page, visually consistent with branches
- [ ] Typing a keyword and submitting filters the tag list
- [ ] Pagination preserves the `q` param across pages
- [ ] Case-insensitive matching works (`v1` finds `V1.0.0`)
- [ ] Empty keyword shows all tags (no regression)

## Out of Scope

- Adding a tag count display in the header (branches has one; tags does not — not adding to stay focused)
- Searching tag commit messages or release notes (only tag names, matching branches behavior)
- API endpoint changes (the `/tags/list` JSON API is a separate handler, unaffected)
