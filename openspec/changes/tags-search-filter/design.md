## Context

Gitea's repository browser exposes three listing surfaces that share visual layout and user expectations:

- `/{owner}/{repo}/branches` — list of branches, **server-side keyword filter** on branch name via `q`.
- `/{owner}/{repo}/commits/{branch}` — list of commits, **server-side keyword filter** via `q` and `all` toggle (`SearchCommits`).
- `/{owner}/{repo}/tags` — list of standalone tags (`Release` rows with `IsTag=true`), **no filter**.

The Branches implementation is the closest analog: handler reads `ctx.FormString("q")`, service passes it as `Keyword` into `git_model.FindBranchOptions`, the underlying query does a LIKE, the template renders `{{template "shared/search/combo" ...}}`, and pagination preserves `q` via `pager.SetDefaultParams(ctx)`. No client-side JS is involved.

The Tags page already paginates via `db.Find[repo_model.Release]` with `FindReleasesOptions{IncludeTags: true, HasSha1: optional.Some(true)}` at `routers/web/repo/release.go:204-251`, page size 10. Adding a keyword filter is a direct extension of this query.

## Goals / Non-Goals

**Goals:**

- A user can type a substring into a search box on `/tags` and see only tags whose `tag_name` contains that substring.
- Behavior, URL shape (`?q=...`), template partial, and pagination semantics match the Branches page exactly — a user who knows one knows the other.
- Empty query degrades cleanly to today's behavior (full paginated list).
- Total-tag count (`NumTags` in the header) continues to reflect the repo total, not the filtered subset, mirroring Branches.

**Non-Goals:**

- No filtering on commit subject, tag annotation message, or any field other than `tag_name`.
- No client-side / live filtering JS. Server-side only.
- No changes to the RSS/Atom feed endpoints or the JSON `/tags/list` selector endpoint.
- No full-text index, no schema migration, no new columns.
- No changes to the Releases page (`/releases`) — that surface has its own list semantics.

## Decisions

**D1 — Filter lives in `FindReleasesOptions`, not in the handler.**
The Branches pattern puts the keyword inside the model-level `FindOptions` so the same query does filter + paginate + count in one shot. Replicating that here means `db.Find` and `db.Count` automatically agree on the filtered set; a handler-side post-filter would force two passes and break pagination correctness.

**D2 — `Keyword` is opt-in.**
`FindReleasesOptions` is shared with the Releases page and several other call sites. Adding `Keyword string` with a `if opts.Keyword != "" { ... }` guard inside `FindReleases` means every existing caller continues to behave identically. No behavioral risk to releases.

**D3 — LIKE pattern is `%kw%` (substring, case-insensitive via column collation).**
Mirrors Branches' behavior exactly. Accepts the leading-wildcard index limitation in exchange for matching user expectations ("find any tag containing 'v1.2'"). Special chars `%` and `_` in user input are interpreted as LIKE wildcards — this is consistent with Branches and deemed acceptable.

**D4 — Empty `q` returns the unfiltered list; no redirect.**
Branches' list handler also accepts empty `q` gracefully (only Commits' *search* sub-route redirects on empty). The Tags list stays put — the search box is always present, the list is always the same endpoint.

**D5 — Template uses `shared/search/combo`, placed in its own `ui attached segment`.**
The same partial Branches uses (`templates/repo/branch/list.tmpl:77-79`), in the same surrounding `<div class="ui attached segment">` wrapper, between the `<h4>` title and the table. One visual pattern across both listing pages.

The current Tags template gates the entire `<h4>` + table inside `{{if .Releases}}`. Restructuring is required: pull the `<h4>` and the new search-segment out of the `{{if}}`, and add an `{{else}}` branch that renders `search.no_results` so an empty search result still shows the search box (letting the user refine the query) instead of a blank page.

**D6 — Locale key `search.tag_kind`.**
Sister keys exist: `search.branch_kind`, `search.commit_kind`. Adding `search.tag_kind = Search tag name` matches the family.

**D7 — `NumTags` (header tab count) stays unfiltered; pagination count is filtered.**
Two distinct counts are at play:
- The header tab shows `{{.NumTags}}` (`templates/repo/release_tag_header.tmpl:10`), set by middleware as the repo total. Stays unfiltered — it answers "how many tags does this repo have".
- The pagination control's total comes from `context.NewPagination(total, ...)` at `release.go:245`. Today `total` is `int(NumTags)`, which is correct only because there is no filter. Once `Keyword` is set, this must become the filtered count via `db.Count[repo_model.Release](ctx, opts)`, exactly as Branches uses `branchesCount` (`routers/web/repo/branch.go:84`). Otherwise pagination shows wrong page count after a search.

## Risks / Trade-offs

**Leading-wildcard LIKE bypasses the `tag_name` index.**
On a repo with tens of thousands of tags, the query falls back to a scan. This is the same trade-off Branches makes today and is acceptable for browsing-scale traffic. If a real-world hot repo proves this wrong, a follow-up change can add a trigram/full-text index — explicitly deferred.

**LIKE wildcard characters in user input.**
A user searching for `100%` will match any tag containing `100` followed by anything. Documented behavior, matches Branches. Not worth escaping (would diverge from the established pattern and surprise users who copy a search between pages).

**`FindReleasesOptions` is a shared type.**
Adding a field is low-risk because the guard makes it opt-in, but any future caller of `FindReleases` inherits a new (default-empty) knob. Acceptable — the field is self-documenting and the default is a no-op.

**No client-side JS means a full page reload per search.**
Consistent with Branches and Commits. Modern Gitea is moving toward more interactive surfaces, but introducing JS here would diverge from the sibling pages and expand scope. Defer until/unless the broader code-browser gets a Vue rewrite.
