## Why

The repository Tags page (`/{owner}/{repo}/tags`) lacks keyword filtering, while the sibling Branches and Commits pages both support a `q` parameter for name-based search. On repositories with many tags, users have no way to quickly locate a specific tag — they must page through the entire list manually. This creates an inconsistent UX across the three ref-listing pages and degrades usability on release-heavy projects.

## What Changes

- Add a `q` query parameter to the Tags page that filters tags by name.
- Filtering is case-insensitive (matches the existing `lower_tag_name` column), so `V1.2` finds `v1.2.3`.
- Render a `shared/search/combo` form above the tag list, identical in shape to the Branches page search input.
- Add `search.tag_kind` locale key (`Search tags...`) to `locale_en-US.ini`; other locales are translated by the community.
- Fix a latent pagination bug exposed by this work: `TagsList` previously used an unfiltered tag count for the pager widget. Switching to `db.FindAndCount` makes the page count reflect the filtered result set, so paginating a search no longer shows empty pages.

No REST API change. No DB migration. No new JavaScript. The `Releases` page handler (`/{owner}/{repo}/releases`) is untouched.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `code-management`: Adds a new requirement under ref-listing for the Tags page to support keyword filtering, matching the existing Branches page behavior. Also fixes the Tags-page pagination count to honor the active filter.

## Impact

**Code paths:**
- `routers/web/repo/release.go` (`TagsList` handler): reads `q`, passes keyword into `FindReleasesOptions`, exposes it to the template via `ctx.Data["Keyword"]`, and uses `db.FindAndCount` for the filtered pager total. Also adds `page` normalization (`page <= 0 → 1`) because `db.FindAndCount` does not normalize `page == 0` the way `db.Find` does — a latent footgun shared with `routers/web/repo/branch.go:50-52`.
- `models/repo/release.go`: `FindReleasesOptions` gains a `Keyword string` field; `ToConds()` gains a `builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)}` clause. Backward compatible — empty keyword skips the clause, so the 12+ existing callers of `FindReleasesOptions` are unaffected.
- `templates/repo/tag/list.tmpl`: inserts the search form between the header partial and the `{{if .Releases}}` block, so the form renders even on empty search results.
- `options/locale/locale_en-US.ini`: adds `tag_kind = Search tags...` under `[search]`.

**APIs:** None affected. The REST endpoint `GET /api/v1/repos/{owner}/{repo}/tags` is unchanged.

**Dependencies:** None added. No DB migration required — the `lower_tag_name` column already exists on the `release` table and is maintained by the model.

**Performance:** `LIKE '%keyword%'` on `lower_tag_name` is a full table scan; the column has no index. This matches the Branches page baseline (which also lacks an index on `name`). Adding an index is explicitly out of scope for this change.

**Tests:** Adds 5 integration subtests under `TestTagsSearch` covering hit/miss/empty/case-insensitive/pagination-with-keyword. Existing tag/release tests continue to pass.
