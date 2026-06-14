## Why

The repository Tags page (`/{org}/{repo}/tags`) lacks a search/filter box, while the sibling Branches and Commits pages both already support keyword filtering. On repositories with many tags (release-managed projects, mirrors of large upstreams), users have no way to locate a specific tag without scanning the full paginated list. This brings the Tags page to parity with its siblings and removes a UX inconsistency.

## What Changes

- Add a search box to the Tags page (`templates/repo/tag/list.tmpl`) that submits a `q` query parameter via GET, mirroring the Branches page search UX.
- Accept `q` in the `TagsList` HTML handler (`routers/web/repo/release.go`) and filter results server-side via a new `Keyword` field on `FindReleasesOptions`.
- Accept `q` in the `GetTagList` JSON handler (`routers/web/repo/repo.go`) and filter the returned tag-name list in Go (the endpoint is unpaginated and the underlying function is shared by 7 other callers).
- Filter is case-insensitive substring match on the tag name, matching Branches semantics.
- RSS/Atom feed endpoints (`/tags.rss`, `/tags.atom`) are intentionally left unfiltered — feeds are designed for subscription, not interactive filtering.
- New locale key `search.tag_kind` added to `options/locale/locale_en-US.ini`.

## Capabilities

### New Capabilities

_None._

### Modified Capabilities

- `code-management`: Adds a new requirement to the "Code Review / Browsing" section specifying that users can filter the Tags page by keyword, with pagination that preserves the filter.

## Impact

**Code paths affected:**
- `models/repo/release.go` — add `Keyword` field to `FindReleasesOptions` and a `builder.Like{"lower_tag_name", ...}` clause in `toConds()`. No schema change (the `lower_tag_name` column already exists and is indexed).
- `routers/web/repo/release.go` — `TagsList` (line 204-251) reads `q`, passes it through `FindReleasesOptions.Keyword`, sets `ctx.Data["Keyword"]`.
- `routers/web/repo/repo.go` — `GetTagList` (line 716-725) reads `q` and filters the returned slice in Go.
- `templates/repo/tag/list.tmpl` — add a `shared/search/combo` form above the tag table.
- `options/locale/locale_en-US.ini` — add `search.tag_kind` key.

**APIs affected:**
- Web (browser-facing): `GET /{username}/{reponame}/tags` gains optional `q` query parameter. No Swagger regen needed (web routes are not Swagger-annotated).
- Web (JSON): `GET /{username}/{reponame}/tags/list` gains optional `q` query parameter. Also not Swagger-annotated (this endpoint serves the in-page autocomplete, not the public REST API at `/api/v1/`).
- Public REST API (`routers/api/v1/`): **not affected.** No changes to `/api/v1/repos/{owner}/{repo}/tags`.

**Dependencies / systems:** None. No new modules, no schema migrations, no new frontend JS (uses the existing `shared/search/combo` partial).

**Tests:**
- Unit test in `models/repo/release_test.go` for the new `Keyword` filter behavior.
- Integration tests in `tests/integration/repo_tag_test.go` for both HTML and JSON endpoints with `?q=...`.
