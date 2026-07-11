## Context

The repository Tags page is a server-rendered web page. The current request path is:

```text
GET /{username}/{reponame}/tags
  -> repository assignment and code-read authorization in routers/web/web.go
  -> repo.TagsList in routers/web/repo/release.go
  -> db.Find[repo_model.Release] with repo_model.FindReleasesOptions
  -> templates/repo/tag/list.tmpl
```

`TagsList` treats tag-backed rows as release records for the repository, including drafts when they point to a real Git object (`sha1` is non-empty). Repository assignment separately calculates `NumTags`, which is rendered in repository navigation. The handler currently reuses that repository-wide count for pagination.

Branches already implements a suitable browser interaction: a GET form using `q`, the shared search combo, a model-level substring condition, `Keyword` template data, and pagination that preserves the keyword. Commits also has search UI, but its separate `/search` path and Git commit-query grammar do not fit a simple tag-name filter.

Tag records contain both `tag_name` and `lower_tag_name`, but the normalized field is not a complete invariant. In particular, `services/release.CreateNewTag` can insert a tag-backed release without setting `lower_tag_name` when the Git tag already exists. The change therefore treats `tag_name` as authoritative and uses the repository's cross-database case-insensitive query helper. The change crosses the web handler, model query options, template, localization source, and integration tests, so the behavior requires a coordinated design.

## Goals / Non-Goals

**Goals:**

- Add a visually consistent, server-rendered tag-name search control to the Tags page.
- Provide case-insensitive substring matching with consistent behavior across supported databases.
- Preserve the keyword through result rendering and pagination.
- Keep repository-wide tag totals distinct from filtered result totals.
- Preserve authorization, ordering, actions, feeds, and the unfiltered page behavior.

**Non-Goals:**

- Searching releases, notes, commits, SHAs, API responses, selector dropdowns, or feeds.
- Introducing client-side filtering, JavaScript, CSS, a new HTTP route, a database migration, or an index.
- Changing how tags are created, synchronized, ordered, authorized, or deleted.
- Editing generated assets, template bindata, or generated Swagger output.

## Decisions

### 1. Use `GET /tags?q=<keyword>` on the existing route

`templates/repo/tag/list.tmpl` will submit a GET form to the current page. `TagsList` will read a trimmed `q` value and expose it as `Keyword`. Empty and whitespace-only values will follow the existing unfiltered path.

This matches the Branches page, produces bookmarkable URLs, and allows `Pagination.SetDefaultParams` to retain `q`. No route, form binding type, API contract, or JavaScript initializer is needed.

**Alternative considered:** Add `/tags/search`, following Commits. Rejected because Commits needs a distinct Git search grammar and handler branch, while tag filtering is a normal variation of the list query.

### 2. Add an explicit tag-name keyword to `FindReleasesOptions`

`models/repo/release.go` will add a narrowly named option such as `TagNameKeyword`. When non-empty, `ToConds` will add `db.BuildCaseInsensitiveLike("tag_name", keyword)`. Existing exact `TagNames` behavior remains unchanged.

The helper applies a portable case-insensitive substring comparison for SQLite, MySQL, PostgreSQL, and MSSQL. Using the authoritative column prevents valid rows with an empty or stale `lower_tag_name` from disappearing, keeps persistence logic in the model layer, and ensures only the tag name participates in matching.

**Alternatives considered:**

- Filter `tag_name` with a raw `builder.Like`, as Branches does for branch names. Rejected because the supported databases can apply different case-sensitivity rules.
- Filter `lower_tag_name` with a lowercased keyword. Rejected because `CreateNewTag` can leave that field empty when registering an existing Git tag.
- Filter the loaded page in Go or JavaScript. Rejected because filtering after pagination produces incomplete results and incorrect totals.
- Query Git tags directly. Rejected because this page intentionally uses tag-backed release rows to preserve release metadata and existing actions.

### 3. Query the page and filtered count together

`TagsList` will normalize the requested page to at least `1` and use `db.FindAndCount[repo_model.Release]` with the same options for both rows and total matches. Page normalization is required because the generic helper applies its limit only when `Page >= 1`. The returned filtered count will initialize the result paginator. The existing `ctx.Data["NumTags"]` remains untouched and continues to represent the repository-wide count in `repo/sub_menu.tmpl` and `repo/release_tag_header.tmpl`.

`Keyword` will be set before calling `pager.SetDefaultParams(ctx)`, causing pagination links to retain `q`. When a caller explicitly supplies a positive `limit`, pagination will also retain the bounded page size so following a link cannot change the result window. Existing ordering from `FindReleasesOptions.ToOrders` remains unchanged.

**Alternative considered:** Reuse `NumTags` for pagination. Rejected because a filtered list could expose empty or duplicate pages based on the unfiltered total.

### 4. Reuse the shared search control and add a filtered-empty state

`templates/repo/tag/list.tmpl` will reuse `shared/search/combo`, following `templates/repo/branch/list.tmpl`. The search section must remain rendered whenever a repository has tags or a keyword is active, even when the filtered result slice is empty. A no-match branch will render the existing `search.no_results` message instead of hiding the entire tag section.

The template will continue to render the existing tag table and action permission checks for matches. No custom CSS or browser code is required. The English locale source will add `search.tag_kind = Search tags...`; other languages remain managed by Crowdin under the repository's localization policy.

**Alternative considered:** Reuse `repo.find_tag`. Rejected because it is worded for selector-style lookup rather than the consistent `Search <kind>...` placeholders used by Branches and Commits.

### 5. Verify model semantics and the rendered workflow

Model-focused tests will cover the new release query condition where practical. `tests/integration/release_test.go` will extend `TestViewTagsList` or add adjacent cases for unfiltered compatibility, partial matching, case-insensitive matching, keyword echo, and no results. A pagination case will prove that the matching count and `q` are used when results span pages.

The focused verification scope is:

- `./tools/harness/verify.sh go ./models/repo/...`
- `make 'test-sqlite#TestViewTagsList'` or the exact adjacent Tags-search test name
- `./tools/harness/verify.sh spec`
- `openspec validate --all --strict --no-interactive`

Multi-database integration suites are broader than the initial change, but the model test should exercise the repository's default SQLite configuration and `BuildCaseInsensitiveLike` is the existing portability boundary for supported databases. Full backend, frontend, and E2E suites are not required unless implementation uncovers shared behavior beyond these control points.

## Risks / Trade-offs

- [Risk] Existing or edge-case rows can have an empty or stale `lower_tag_name`. -> Query authoritative `tag_name` through `db.BuildCaseInsensitiveLike`; do not add a migration for this feature.
- [Risk] A leading-wildcard substring query cannot use a conventional index efficiently on very large tag sets. -> Keep the scope to one repository, retain page limits, and avoid adding an index that would not accelerate `%keyword%`; measure before proposing a specialized search index.
- [Risk] Search UI disappears after a zero-result query because the current template wraps the entire list section in `if .Releases`. -> Move the search container outside the result-only condition and add an explicit filtered-empty branch.
- [Risk] Using the repository-wide `NumTags` for result pagination exposes incorrect pages. -> Use the filtered count returned from the exact result query while leaving `NumTags` unchanged for navigation.
- [Trade-off] Only English source text is added in this repository. -> Follow the established Crowdin workflow; do not hand-edit translated locale files that will be overwritten.

## Migration Plan

No schema or data migration is planned. Deployment consists of the model option, handler query, template, locale source, and tests. Existing URLs without `q` retain their current behavior.

Rollback removes those source changes; stored data and public routes require no rollback. Generated bindata or assets must not be edited as part of either deployment or rollback.

## Open Questions

None. The incomplete `lower_tag_name` path was resolved in the design by querying authoritative `tag_name`; no migration or repair scope is required.
