## Context

The Tags page (`routers/web/repo/release.go:TagsList`) and Releases page share the `FindReleasesOptions` filter struct and the underlying `release` table. The Tags handler currently supports only `page` and `limit` query parameters — no keyword filtering. The Branches (`routers/web/repo/branch.go:Branches`) and Commits (`routers/web/repo/commit.go:SearchCommits`) pages both support a `q` query parameter for name-based filtering, making Tags the inconsistent outlier across the three ref-listing pages.

The `release` table already maintains a `lower_tag_name` column on every tag insert/update (`models/repo/release.go:154,188`), so case-insensitive name matching needs no `LOWER()` SQL function call and behaves identically across SQLite/MySQL/PostgreSQL.

A latent pagination bug also exists in the current Tags handler: it pulls the pager's total count from `ctx.Data["NumTags"]`, which is computed once in `services/context/repo.go:517-522` (during repo-assignment middleware) **without** any keyword filter. So even if a keyword were applied to the page query, the pager widget would still show the unfiltered total. The Branches page avoids this by using `db.FindAndCount` (`services/repository/branch.go:78`), which returns both the filtered page and the filtered total in one call.

Constraints:
- Backport-focused branch (v1.22.x) — avoid unnecessary refactors.
- No DB migrations, no new dependencies, no new JavaScript.
- Match the existing Branches search UX so users get a consistent experience.

## Goals / Non-Goals

**Goals:**
- Users can filter the Tags page by tag name via a `q` query parameter.
- Matching is case-insensitive substring on the tag name.
- The Tags page search UX matches the Branches page (same `shared/search/combo` partial, same `q` parameter, same `SetDefaultParams` pagination carry-through).
- The pager widget reflects the filtered count, so paginating a search does not produce empty pages.
- No regression to the Releases page or any other caller of `FindReleasesOptions`.

**Non-Goals:**
- Searching by commit message, author, or date — tag name only.
- AJAX instant search — GET form submit is enough.
- Search-result highlighting.
- REST API changes (`GET /api/v1/repos/{owner}/{repo}/tags`).
- Adding an index on `lower_tag_name` (deferred optimization).
- Refactoring the shared search infrastructure.

## Decisions

### Decision 1: Reuse the Branches search pattern verbatim, with one improvement

**Choice:** Mirror `routers/web/repo/branch.go:55-79` end-to-end: read `q` via `ctx.FormString`, pass keyword into the options struct, set `ctx.Data["Keyword"]`, render `shared/search/combo` in the template.

**Rationale:** Project convention is consistency across sibling pages. Branches already established the pattern; deviating would create new cognitive load for no benefit.

**Alternatives considered:**
- **Client-side JS filtering.** Rejected: incompatible with server-side pagination. The whole point of the feature is to find a tag on a repo with many tags, where the user is on page 1 of N. JS filtering only filters the current page.
- **Shared search-helper abstraction.** Rejected: only two call sites with meaningfully different shapes (branches use `git_model.FindBranchOptions`, tags use `repo_model.FindReleasesOptions`). Premature abstraction.

### Decision 2: Filter on `lower_tag_name` for case-insensitive matching (improvement over Branches)

**Choice:** Add `Keyword string` to `FindReleasesOptions` and append `builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)}` in `ToConds()`.

**Rationale:** The `lower_tag_name` column already exists and is maintained by the model on every write. Branches uses `builder.Like{"name", ...}` (case-sensitive). Using `lower_tag_name` is a free upgrade — case-insensitive matching with no `LOWER()` call and consistent behavior across DB backends. The pattern is established elsewhere in the codebase (`models/packages/*.go`, `models/organization/org.go`).

**Alternatives considered:**
- `builder.Like{"tag_name", opts.Keyword}` — case-sensitive, matches Branches exactly but is a worse UX. Rejected.
- Custom ILIKE for PostgreSQL — diverges across DB backends. Rejected.

### Decision 3: Switch `TagsList` from `db.Find` to `db.FindAndCount`

**Choice:** Replace `releases, err := db.Find[repo_model.Release](ctx, opts)` with `releases, total, err := db.FindAndCount[repo_model.Release](ctx, opts)`, and feed `total` (filtered) to `context.NewPagination` instead of `ctx.Data["NumTags"]` (unfiltered).

**Rationale:** `db.Find` discards the count. Reusing the middleware-computed `NumTags` would silently break pagination whenever a keyword narrows the result set. `db.FindAndCount` is the established Branches pattern and returns both values in one DB round-trip (the COUNT is issued alongside the SELECT).

**Side effect — page normalization:** `db.FindAndCount` (at `models/db/list.go:188-215`) only applies `sess.Limit(...)` when `opts.Page >= 1`, whereas `db.Find` normalizes `page == 0 → 1` internally. Without handler-side normalization, a request with no `page` param (the common case) would silently return all matching rows, ignoring `limit`. The handler now normalizes: `if listOptions.Page <= 0 { listOptions.Page = 1 }`, matching `routers/web/repo/branch.go:50-52`.

**Alternatives considered:**
- Fix `db.FindAndCount` to normalize page==0. Rejected: it may be intentional for callers wanting no-pagination semantics; modifying the shared helper risks the ~hundreds of existing callers.
- Compute a separate filtered count via `db.Count`. Rejected: two queries instead of one.

### Decision 4: Keep `ctx.Data["NumTags"]` untouched

**Choice:** Leave the middleware-computed unfiltered `NumTags` in place. It is consumed by the repo header to display total tag count regardless of any active filter.

**Rationale:** `NumTags` is repo-level metadata, not page-level result count. Reusing it for pagination was the bug; deleting it would break the header.

### Decision 5: Template form placement — outside `{{if .Releases}}`

**Choice:** Insert the search form between `{{template "repo/release_tag_header" .}}` and `{{if .Releases}}` in `templates/repo/tag/list.tmpl`.

**Rationale:** If the form were inside the `{{if}}` block, a zero-result search would hide the form, so users could not retry with a different keyword. Placing it outside guarantees the search box is always visible.

## Risks / Trade-offs

**`lower_tag_name` has no index → full table scan on `LIKE '%keyword%'`** → Acceptable. Matches Branches baseline. If a future performance investigation flags this, add `INDEX lower_tag_name` via a separate migration. Explicitly out of scope here.

**i18n coverage gap on day 1** → Only `locale_en-US.ini` ships the new `tag_kind` key. → Acceptable. Community translation workflow handles the 40+ other locales; `design.md` §6 establishes this as the canonical pattern.

**`db.FindAndCount` issues an extra COUNT query on every Tags page render** → Negligible: a single indexed count on a per-repo filter. The previous code already issued an equivalent count in middleware (`services/context/repo.go:517`). Net new queries: zero on the no-keyword path (the middleware count is still needed for the header); one on the keyword path (the filtered count).

**Test fixture size limits regression coverage** → The default fixture repo (`user2/repo1`) has only 3 visible tags. The `PaginationWithKeyword` subtest works around this by using `q=v1.0&limit=1` to force the filtered-vs-unfiltered page-count divergence (`1` vs `3`) into the on/off threshold for the pagination widget.

**No integration test for AJAX-style live search** → N/A: there is no AJAX search in this design.
