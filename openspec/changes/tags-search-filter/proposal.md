## Why

The repository Tags page (`/{owner}/{repo}/tags`) has no way to search or filter tags by name. For repositories with hundreds or thousands of tags, users cannot quickly locate a specific tag. The sibling Branches and Commits pages already provide server-side keyword filtering on the same `/branches` and `/commits/{branch}/search` routes; Tags is the only listing page in the code browser without parity, breaking the established UX consistency.

## What Changes

Add server-side keyword filtering to the Tags page, matching the Branches page pattern exactly:

- The `TagsList` handler reads a `q` query param and applies it as a substring filter on `tag_name`.
- The `FindReleasesOptions` type gains an opt-in `Keyword` field; when set, `ToConds()` adds `builder.Like{"tag_name", opts.Keyword}` — the same condition the Branches page applies to `name`. xorm's `Like` auto-wraps the value as `%kw%`. Case sensitivity follows each database's default collation (CI on SQLite/MySQL-default, CS on PostgreSQL) — identical to Branches behavior.
- The Tags template renders a search form using the existing `shared/search/combo` partial, pre-filled with the current keyword.
- Pagination links preserve `q` via `pager.SetDefaultParams(ctx)`, the same mechanism Branches uses.
- No client-side JavaScript; the form is a plain GET submit.
- RSS/Atom feed endpoints (`.rss`, `.atom`) and the JSON `/tags/list` selector endpoint are unchanged.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `releases` (PKG-04): adds requirements governing tag-name filtering on the standalone Tags listing page. Existing release workflow requirements are unchanged; the new requirements attach only to the `/{owner}/{repo}/tags` browsing surface.

## Impact

- **Code:**
  - `models/repo/release.go` — `FindReleasesOptions` gets a `Keyword string` field; `ToConds()` adds `builder.Like{"tag_name", opts.Keyword}` when set (matching `models/git/branch_list.go:102-104`).
  - `routers/web/repo/release.go` — `TagsList` reads `q`, threads it through `opts.Keyword`, exposes `ctx.Data["Keyword"]`, and replaces the pagination's total-count source from the middleware-set `NumTags` to a filtered `db.Count[repo_model.Release](ctx, opts)` so page count reflects the search result (Branches does the same with `branchesCount`).
  - `templates/repo/tag/list.tmpl` — pulls the `<h4>` and a new search-segment out of the `{{if .Releases}}` gate, inserts a `shared/search/combo` form in its own `ui attached segment`, and adds an `{{else}}` branch rendering `search.no_results` for empty search results.
  - `options/locale/locale_en-US.ini` — adds a `tag_kind` entry under the `[search]` section (sister to `branch_kind` / `commit_kind`).
- **APIs:** No REST API change. Swagger unaffected.
- **DB schema:** No migration. The LIKE runs against the existing `tag_name` column on `release`.
- **Tests:** New unit test in `models/repo/release_test.go`; new integration test in `tests/integration/repo_tag_test.go` (or extension of existing).
- **Performance:** A leading-wildcard LIKE (`%kw%`) cannot use the `tag_name` index. Acceptable for browsing-scale traffic; mirrors the Branches page's identical trade-off. If a repo ever hits pathological scale, a future change can add a full-text index — out of scope here.
