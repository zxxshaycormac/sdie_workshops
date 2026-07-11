## Why

The repository Tags page has no way to narrow a long tag list by name, even though the neighboring Branches and Commits pages expose search controls. Users need a consistent way to find a tag without paging through the entire list.

## What Changes

- Add a keyword search control to the repository Tags page, using the existing shared search UI pattern.
- Filter tags by a case-insensitive substring of the tag name through the existing `GET /{owner}/{repo}/tags` route and `q` query parameter.
- Preserve the keyword while paging and base pagination on the number of matching tags.
- Keep the search control available when a query has no matches and show a clear no-results state.
- Preserve the repository's total tag count in navigation; search result counts do not replace it.

### Non-goals

- Do not add tag search to the public API, tag selector dropdowns, RSS/Atom feeds, Releases page, or repository-wide search.
- Do not search release titles, notes, commit messages, or commit SHAs.
- Do not add client-side filtering, new JavaScript, generated assets, database columns, or migrations.
- Do not change tag visibility, authorization, ordering, creation, deletion, or download behavior.

## Capabilities

### New Capabilities

- `repository-tag-search`: Search and page repository tags by a tag-name keyword on the Tags page.

### Modified Capabilities

None. The current main specs only describe the engineering harness and do not define repository tag-list behavior.

## Impact

### Observed control points

- Web routing: `routers/web/web.go` already maps `GET /{username}/{reponame}/tags` to `repo.TagsList`; no route addition is required.
- Handler and pagination: `routers/web/repo/release.go` builds `FindReleasesOptions`, loads tag-backed release records, and currently paginates with the repository-wide `NumTags` value.
- Persistence query: `models/repo/release.go` defines `FindReleasesOptions`; tag names are stored in `tag_name`, while `lower_tag_name` is populated by most but not all creation paths.
- Template: `templates/repo/tag/list.tmpl` renders the tag list but has no search or filtered-empty state.
- Localization: the shared search component requires an English tag-search placeholder in `options/locale/locale_en-US.ini`; non-English translations remain managed through Crowdin.
- Tests: `tests/integration/release_test.go` already covers the rendered Tags page.

### Verified implementation constraint

- `CreateNewTag` can insert a tag-backed release with an empty `lower_tag_name` when the Git tag already exists, so search must use the authoritative `tag_name` field with the repository's cross-database case-insensitive query helper.

### Compatibility surface

The existing route, permissions, default unfiltered response, tag ordering, feed URLs, and repository-wide tag count remain compatible. The only new public browser behavior is the optional `q` query parameter and filtered rendered result set.
