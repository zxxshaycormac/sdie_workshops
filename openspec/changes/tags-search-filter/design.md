## Context

The repository Tags page (`routers/web/repo/release.go:204` `TagsList`) is the only one of the three "ref-browsing" pages (Branches, Commits, Tags) without a search/filter box. Branches added one via `q := ctx.FormString("q")` plumbed through `git_model.FindBranchOptions.Keyword` into `builder.Like{"name", opts.Keyword}` in `models/git/branch_list.go:102-104`. The Tags page handler already paginates through the `release` table (`db.Find[repo_model.Release]` with `IncludeTags: true`) and the `release` table already has a `lower_tag_name` column (`models/repo/release.go:76`) — both prerequisites for a DB-side LIKE filter are in place.

A JSON sibling endpoint `GET /{org}/{repo}/tags/list` (`routers/web/repo/repo.go:716` `GetTagList`) returns the full tag-name list (no pagination) for autocomplete-style consumers. It uses `repo_model.GetTagNamesByRepoID`, a function shared by 7 other call sites (compare view, issue autocomplete, RSS/Atom feeds).

The Branches page is the canonical reference for both UX (the `shared/search/combo` template, `pager.SetDefaultParams(ctx)` propagating `q` across pages) and backend pattern (per-options-struct `Keyword` field + DB LIKE).

## Goals / Non-Goals

**Goals:**
- Add a `q` query parameter to the Tags HTML page that filters tags by case-insensitive substring match on the tag name.
- Add the same parameter to the `/tags/list` JSON endpoint.
- Mirror the Branches UX end-to-end: same search partial, same `q` URL parameter, same pagination continuity.
- No schema migrations, no signature changes to shared functions, no new frontend JS.

**Non-Goals:**
- Filtering by release title or release notes (tag name only).
- Fuzzy/exact search toggle (issues page has `combo_fuzzy`; Tags uses plain `combo`).
- Filtering the RSS/Atom feed endpoints (feeds are designed for subscription, not interactive filtering).
- Changes to the public REST API at `/api/v1/repos/{owner}/{repo}/tags` (out of scope; tracked separately if needed).
- Client-side instant-filter JS (plain form submit matches Branches).

## Decisions

### Decision 1 — Filter at DB layer via new `Keyword` field on `FindReleasesOptions`

Add `Keyword string` to `FindReleasesOptions` and a clause in `toConds()`:

```go
if opts.Keyword != "" {
    cond = cond.And(builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)})
}
```

**Why over alternatives:**
- *Filter inline in handler* — duplicates the WHERE-clause logic across `TagsList` and any future caller; bypasses the typed `db.Find[repo_model.Release]` API; forces the handler to drop to engine-level.
- *Filter in Go after `db.Find`* — breaks pagination: `numTags` no longer matches the filtered result count, page sizes become wrong, large repos OOM. Non-starter for the same reason Branches doesn't do this.
- *Pre-split keyword into `TagNames`* — `TagNames` does exact-match `builder.In`, not substring.

The options-struct pattern is already established on this struct (`TagNames`, `IsPreRelease`, `IsDraft`, `HasSha1` all follow the same shape) and on the sibling `FindBranchOptions.Keyword`.

### Decision 2 — Use `lower_tag_name` + `strings.ToLower(kw)` for case-insensitive matching

The `release` table has both `tag_name` and `lower_tag_name` columns; the latter is already populated and indexed.

**Why over `builder.Like{"tag_name", kw}`** — relies on DB default collation, which is fragile: MySQL `utf8mb4_general_ci` is case-insensitive, PostgreSQL `COLLATE` varies, SQLite depends on `LIKE` vs `GLOB`. The Branches implementation has this fragility on the `name` column; Tags will not.

### Decision 3 — JSON endpoint filters in Go after fetch, not via shared-function signature change

`GetTagList` (`routers/web/repo/repo.go:716`) calls `repo_model.GetTagNamesByRepoID`, which is shared by 7 call sites and returns the full unpaginated list. The JSON handler filters in Go:

```go
if q := ctx.FormTrim("q"); q != "" {
    lower := strings.ToLower(q)
    filtered := tags[:0]
    for _, t := range tags {
        if strings.Contains(strings.ToLower(t), lower) {
            filtered = append(filtered, t)
        }
    }
    tags = filtered
}
```

**Why over alternatives:**
- *Add `keyword` parameter to `GetTagNamesByRepoID`* — forces all 7 callers (compare view, issue autocomplete, RSS/Atom feeds) to pass `""` for no behavior change. Noisy diff, real risk of someone passing the wrong thing.
- *Add a new model-layer function `GetTagNamesByRepoIDAndKeyword`* — duplicates the existing query for one caller; not worth it.
- *Use the DB-side filter from Decision 1* — the endpoint is unpaginated and returns a slice of names, not releases; running a second query just to filter is more work than filtering the slice we already have.

Semantically equivalent to the HTML path: case-insensitive substring match. Reuses the slice's backing array (no allocation when no matches).

### Decision 4 — Reuse the `shared/search/combo` template partial

The Tags template adds the same form Branches uses (`templates/repo/branch/list.tmpl:76-80`):

```html
<div class="ui attached segment">
    <form class="ignore-dirty" method="get">
        {{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}
    </form>
</div>
```

`shared/search/combo` already renders the `name="q"` input + submit button. No new JS needed; plain form-on-submit matches Branches.

### Decision 5 — Do NOT escape SQL LIKE wildcards (`%`, `_`)

User input containing `%` or `_` will be interpreted as wildcards by the DB. This is a pre-existing condition on the Branches search (`models/git/branch_list.go:103` is also unescaped) and is treated as a low-risk self-DoS rather than a security issue.

**Why not escape:** Branches sets the precedent, escaping adds complexity (XORM's `builder.Like` doesn't easily support an `ESCAPE` clause), and the user can only affect their own search results.

## Risks / Trade-offs

- **[Risk] SQL LIKE wildcard characters in user input match more than intended** → Mitigation: documented as a known limitation matching Branches; users can self-recover by refining their query.
- **[Risk] Performance on repos with very many tags** → Mitigation: `lower_tag_name` is already indexed; LIKE with a leading-wildcard pattern (`%foo%`) cannot use the index, but the existing pagination caps the result set per page. No regression vs. the unfiltered query for the same page size.
- **[Risk] Empty keyword changes behavior for other `FindReleasesOptions` callers** → Mitigation: the `if opts.Keyword != ""` guard makes empty keyword a strict no-op. All existing callers (releases page, RSS/Atom, API) construct `FindReleasesOptions` without setting `Keyword`, so they see zero behavior change.
- **[Trade-off] JSON endpoint uses in-memory filter, HTML endpoint uses DB filter** → Acceptable because the two endpoints have different shapes (paginated vs full-list) and the JSON path avoids touching a shared function with 7 callers. Semantically equivalent.

## Migration Plan

No migrations. The change is purely additive at the HTTP and Go layer; the `lower_tag_name` column already exists. Deployment is a standard rolling restart.

**Rollback:** Revert the commit; existing queries without `q` continue to work because the field was already optional.

## Open Questions

None — all decisions resolved during the design conversation. The Branches page provides a complete reference for every choice.
