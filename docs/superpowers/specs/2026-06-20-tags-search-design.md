# Tags Page Keyword Search — Design

**Date:** 2026-06-20
**Author:** brainstorming session (via `superpowers:brainstorming`)
**Change ID:** `tags-search-filter` (OpenSpec scaffold at `openspec/changes/tags-search-filter/`)
**Branch target:** Gitea 1.22.x (backport-focused)

---

## 1. Problem

The repository Tags page (`/{owner}/{repo}/tags`) has no way to filter tags by name. On repositories with hundreds or thousands of tags, users cannot quickly find a specific release. The sibling **Branches** and **Commits** pages both already provide keyword search, creating an inconsistent experience.

## 2. Goal

Users can search tags by name on `/{owner}/{repo}/tags`, with UX consistent with the Branches page (`/{owner}/{repo}/branches?q=...`).

## 3. Non-Goals

- Search on the Releases page (`/{owner}/{repo}/releases`) — shares handler and data layer with Tags, but explicitly out of scope. Data-layer changes are written so Releases can opt in later without further refactor.
- Filtering the RSS/Atom tag feeds (`/tags.rss`, `/tags.atom`) — feed readers expect the full list.
- Changes to `GetTagList` (`/tags/list`, the JSON dropdown endpoint).
- New JS modules, new packages, new routes, DB migrations.

## 4. Approach

**Chosen: server-side keyword filtering** (Approach A from the brainstorm).

Replicate the Branches search flow exactly: `?q=<kw>` query param → server-side `LIKE` filter at the data layer → pre-filled search input via `shared/search/combo` → pager carries `q` across pagination links.

Rejected alternatives:
- **Client-side filtering** — breaks pagination (page only renders 10–30 tags today); over-fetches on large repos.
- **API + fetch-render** — loses URL state (search not shareable/bookmarkable); worse SEO; inconsistent with Branches (which is a full page reload).

## 5. Component Changes (4 files)

### 5.1 `models/repo/release.go` — data-layer keyword filter

Add a `Keyword string` field to `FindReleasesOptions` and a `builder.Like("tag_name", opts.Keyword)` clause to `ToConds()`:

```go
type FindReleasesOptions struct {
    db.ListOptions
    RepoID        int64
    IncludeDrafts bool
    IncludeTags   bool
    IsPreRelease  optional.Option[bool]
    IsDraft       optional.Option[bool]
    TagNames      []string
    Keyword       string  // NEW — substring match on tag_name
    HasSha1       optional.Option[bool]
}

func (opts FindReleasesOptions) ToConds() builder.Cond {
    var cond builder.Cond = builder.Eq{"repo_id": opts.RepoID}
    // ...existing clauses unchanged...
    if opts.Keyword != "" {
        cond = cond.And(builder.Like("tag_name", opts.Keyword))
    }
    return cond
}
```

This mirrors `git_model.FindBranchOptions.Keyword` → `builder.Like("name", ...)`. No new index — substring `LIKE '%kw%'` cannot use a b-tree index, and the result set is bounded by `RepoID` (already indexed). Acceptable for v1.22.x backport scope; revisit with a trigram/pg_trgm index only if production metrics show a real bottleneck.

Every existing caller of `FindReleasesOptions.ToConds()` (e.g. `GetTagNamesByRepoID`, releases list, single release lookup) continues to behave identically because `Keyword` defaults to `""` and the new clause is skipped when empty.

### 5.2 `routers/web/repo/release.go` — `TagsList` handler

Read `q`, pass it to the options, set `ctx.Data["Keyword"]`, and compute a filtered count for pagination when the keyword is non-empty:

```go
func TagsList(ctx *context.Context) {
    // ...existing ctx.Data setup unchanged...

    kw := ctx.FormString("q")

    listOptions := db.ListOptions{
        Page:     ctx.FormInt("page"),
        PageSize: ctx.FormInt("limit"),
    }
    if listOptions.PageSize == 0 {
        listOptions.PageSize = setting.Repository.Release.DefaultPagingNum
    }
    if listOptions.PageSize > setting.API.MaxResponseItems {
        listOptions.PageSize = setting.API.MaxResponseItems
    }

    opts := repo_model.FindReleasesOptions{
        ListOptions:   listOptions,
        IncludeDrafts: true,
        IncludeTags:   true,
        HasSha1:       optional.Some(true),
        RepoID:        ctx.Repo.Repository.ID,
        Keyword:       kw,
    }

    releases, err := db.Find[repo_model.Release](ctx, opts)
    if err != nil {
        ctx.ServerError("GetReleasesByRepoID", err)
        return
    }
    ctx.Data["Releases"] = releases

    // CHANGED: filtered count when searching; middleware-set NumTags otherwise.
    var total int64
    if kw != "" {
        total, err = db.Count[repo_model.Release](ctx, opts)
        if err != nil {
            ctx.ServerError("CountReleases", err)
            return
        }
    } else {
        total = ctx.Data["NumTags"].(int64)
    }
    pager := context.NewPagination(int(total), opts.PageSize, opts.Page, 5)
    pager.SetDefaultParams(ctx)  // carries q + page across pagination links
    ctx.Data["Page"] = pager
    ctx.Data["Keyword"] = kw

    ctx.Data["PageIsViewCode"] = !ctx.Repo.Repository.UnitEnabled(ctx, unit.TypeReleases)
    ctx.HTML(http.StatusOK, tplTagsList)
}
```

**Why two count paths:** the middleware-set `NumTags` (`services/context/repo.go:517`) counts *all* tags for the repo header badge. Under search, that count is wrong for pagination. The branch is cheap — one `COUNT(*) WHERE repo_id = ? AND tag_name LIKE ?`.

### 5.3 `templates/repo/tag/list.tmpl` — search input + empty state

Insert the search form above the table (mirroring `templates/repo/branch/list.tmpl:76-80`):

```html
<div class="ui attached segment">
    <form class="ignore-dirty" method="get">
        {{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}
    </form>
</div>
```

Convert the existing `{{if .Releases}}` block so a search with no hits renders an explicit empty-state message instead of an empty gap:

```html
{{if .Releases}}
    {{/* existing table */}}
{{else if .Keyword}}
    <div class="ui attached segment center">
        {{ctx.Locale.Tr "repo.release.tags.no_match" .Keyword}}
    </div>
{{end}}
```

### 5.4 `options/locale/locale_en-US.ini` — two new keys

```ini
[search]
tag_kind = Search tags…
[repo.release]
tags.no_match = No tags match "%s"
```

`search.branch_kind` already exists as the sibling placeholder — `tag_kind` slots in next to it. Per `design.md` §6 we touch `locale_en-US.ini` only; other locales are translated by the community.

## 6. Data Flow

```
GET /{owner}/{repo}/tags?q=v1.2&page=2
  │
  ▼
[Chi] → context.RepoAssignment → context.RepoRefByType(Tag)
  │
  ▼
repo.TagsList(ctx)
  │   reads kw := ctx.FormString("q")           ← "" when absent
  │   builds FindReleasesOptions{Keyword: kw}
  │
  ├──► db.Find[Release](ctx, opts)
  │      └─ ToConds() → builder.Like("tag_name", kw)   ← parameterised
  │
  ├──► db.Count[Release](ctx, opts)             ← only when kw != ""
  │
  ▼
ctx.Data["Releases", "Keyword", "Page"]
ctx.HTML(tplTagsList)
  │
  ▼
templates/repo/tag/list.tmpl
  ├─ search input pre-filled with {{.Keyword}}
  ├─ matching tags table
  └─ pager.SetDefaultParams carries ?q=v1.2&page=N across links
```

**State transitions:**
- Empty `q` → behaves exactly like today (unfiltered, uses middleware's `NumTags`).
- Non-empty `q`, matches exist → table shows matches; pagination reflects filtered count.
- Non-empty `q`, no matches → empty-state message; pagination shows 0 pages.
- User clears the box and submits → `q=` is sent; handler treats empty as "no filter"; full list returns.

## 7. Error Handling

| Failure | Where | Response |
|---|---|---|
| `db.Find` fails | `TagsList` | `ctx.ServerError("GetReleasesByRepoID", err)` — same pattern the handler already uses. 5xx page rendered by `services/context`. |
| `db.Count` fails (search path only) | `TagsList` | `ctx.ServerError("CountReleases", err)`. New log message string; identical handling. |
| Repo empty / archived | Middleware (`MustBeNotEmpty`, `RepoMustNotBeArchived`) | Already gated before `TagsList` runs — no change. |
| Permission denied | `reqRepoCodeReader` middleware | Already gated — no change. |
| Invalid `page` (e.g. `?page=abc`) | `ctx.FormInt("page")` | Returns `0`; `NewPagination` clamps to page 1. Pre-existing behavior. |
| Keyword with SQL special chars (`%`, `_`, `\`) | `builder.Like` | XORM parameterises and escapes the value; a user-supplied `%` is a literal percent in the search string, not a wildcard. Not a SQL-injection vector. |

**Explicitly not added:** length cap on `q` (Branches has none either); client-side debounce (form is submit-on-Enter, same as Branches); client-side validation (single text field).

## 8. Testing

### 8.1 Data-layer unit test — `models/repo/release_test.go`

Extend the existing test file. New test covers the `Keyword` clause directly:

```go
func TestFindReleasesOptions_KeywordFilter(t *testing.T) {
    unittest.PrepareTestDatabase(t)

    cases := []struct {
        kw             string
        expectTagNames []string
    }{
        {"v1",           []string{"v1.0", "v1.1"}},      // substring
        {"nonexistent",  nil},
    }
    // assert returned tag names against fixtures for each kw
    //
    // Note on case sensitivity: SQL `LIKE` is case-insensitive on MySQL/SQLite
    // and case-sensitive on PostgreSQL by default. Branches search has the
    // same DB-dependent behaviour. Do NOT assert a mixed-case query in this
    // test — the result would differ between `make test-sqlite` (CI) and a
    // local Postgres run. If case-insensitive matching is required across
    // backends in the future, that is a separate change (ILIKE / LOWER()).
}
```

Uses `unittest.PrepareTestDatabase` + `models/fixtures/release.yml` (already populated with tag rows).

### 8.2 Integration test — `tests/integration/repo_tag_test.go`

A new file under the integration framework:

```go
func TestTagsListSearch(t *testing.T) {
    defer tests.PrepareTestEnv(t)()
    // GET /{owner}/{repo}/tags?q=<kw>
    // assert: 200 status
    // assert: response body contains only matching tag names
    // assert: response body does NOT contain non-matching tag names
    // assert: search input value="kw" present in HTML
    // assert: pagination link carries q=kw
}
```

Three sub-cases: `q` matches subset; `q` matches nothing (empty-state message rendered); `q` absent (today's behavior preserved — regression guard).

### 8.3 E2E (Playwright) — `tests/e2e/`

Add only if a comparable Branches-search E2E exists to mirror. **Not in the Definition of Done unless an existing branches-search E2E establishes the precedent.** Check during implementation.

### 8.4 Definition of Done

- [ ] `FindReleasesOptions.Keyword` covered by `models/repo/release_test.go`
- [ ] `TagsList` search flow covered by `tests/integration/repo_tag_test.go`
- [ ] `make test-backend` passes
- [ ] `make lint` passes (especially `make lint-templates` for the `.tmpl` edit)
- [ ] Manually verified in browser: empty `q`, matching `q`, non-matching `q`, pagination carries `q` to page 2
- [ ] Locale keys added to `locale_en-US.ini` only — no machine translation
- [ ] Acceptance scenarios below pass

## 9. Acceptance Scenarios (BDD)

```
Scenario: Substring search returns matching tags
  Given a repository with tags "v1.0", "v1.1", "v2.0"
  When the user visits /tags?q=v1
  Then the page shows "v1.0" and "v1.1"
  And  does not show "v2.0"
  And  the search input is pre-filled with "v1"
  And  pagination links carry q=v1

Scenario: No matches shows guided empty state
  Given a repository with tags "v1.0", "v1.1"
  When the user visits /tags?q=zzz
  Then the page shows "No tags match 'zzz'"
  And  no pagination links are rendered

Scenario: Absent q preserves existing behavior
  Given a repository with tags "v1.0", "v1.1", "v2.0"
  When the user visits /tags
  Then all tags are listed
  And  pagination reflects the unfiltered total
```

## 10. Cross-Stage Sanity (per CLAUDE.md SDIE)

- **(S × I):** code matches spec — handler reads `q`, data layer filters on `tag_name`, template renders the search input. ✓
- **(D × I):** code matches design conventions — SQL built via `builder.Like` (§8 of `design.md`), no inline `<script>` (templates §7), locale keys via `ctx.Locale.Tr` (§6). ✓
- **(S × E):** acceptance scenarios test exactly what the spec promises. ✓

## 11. Open Items / Follow-Ups

- Migrate this design into the OpenSpec change at `openspec/changes/tags-search-filter/` as `proposal.md` + delta spec files before implementation, so the change directory is no longer an empty scaffold. (Can happen in the writing-plans phase.)
- If a Branches-search E2E exists, mirror it for Tags. Otherwise skip — integration test covers the contract.
