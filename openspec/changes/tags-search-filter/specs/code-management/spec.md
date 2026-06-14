## ADDED Requirements

### 1. Tag Search

**User Story:** As a developer browsing a repository with many tags, I want to search tags by name so that I can quickly locate a specific release without scanning the full paginated list.

#### Ubiquitous Requirements (Tag Search Constraints)

- **TAG-01-001:** `The system shall provide a search input on the repository tags page that accepts a tag-name filter query.`
- **TAG-01-002:** `The system shall evaluate the tag search query as a case-insensitive substring match against the tag name.`
- **TAG-01-003:** `The system shall propagate the current search query across pagination links on the tags page.`

#### Event-Driven Requirements (Tag Search Workflow)

- **TAG-01-101:** `When a user submits the tags page search form with a non-empty query, the system shall display only tags whose names contain the query as a case-insensitive substring.`
- **TAG-01-102:** `When a user submits the tags page search form with an empty query, the system shall display the unfiltered tag list.`
- **TAG-01-103:** `When a user navigates to a subsequent page of filtered tag results, the system shall preserve the original search query in the request URL.`
- **TAG-01-104:** `When a client requests the tags JSON endpoint with a non-empty query parameter, the system shall return only tag names containing the query as a case-insensitive substring.`
- **TAG-01-105:** `When a client requests the tags JSON endpoint with an empty query parameter, the system shall return the complete tag-name list.`

#### Unwanted Behaviour Requirements (Tag Search Edge Cases)

- **TAG-01-301:** `If a tag search query contains characters that have no tag-name match, then the system shall display an empty result set with no error.`
- **TAG-01-302:** `If a tag search query contains SQL LIKE wildcard characters, then the system shall defer to the database's LIKE semantics for those characters, matching the existing branches search behavior.`

---

## Business Rules

- **BR-01-001:** Tag search applies only to the repository currently being viewed; it does not search across repositories.
- **BR-01-002:** Tag search does not affect RSS or Atom feed endpoints, which always return the unfiltered tag list.
- **BR-01-003:** Tag search respects the requesting user's repository read permission (enforced by existing route middleware; no new permission check is introduced).

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Empty query string on HTML page | Full unfiltered tag list, paginated as before |
| Empty query string on JSON endpoint | Complete tag-name list, unchanged from current behavior |
| Query with only whitespace | Treated as empty (whitespace trimmed); full list returned |
| Uppercase query (e.g. `V1.2`) | Same matches as lowercase (`v1.2`); case-insensitive |
| Query with no matches | Empty result set; no error |
| Query containing `%` or `_` | DB LIKE wildcards apply (e.g. `v%` matches `v1.0.0`, `v2.0.0`); documented limitation, matches Branches |
| Pagination link on filtered results | Carries `?q=<query>&page=<n>` via `pager.SetDefaultParams(ctx)` |
| RSS/Atom feed request with `?q=` | `q` ignored; full unfiltered feed returned |

## Success Criteria

- All existing tests pass without modification (additive change only).
- New unit test `TestFindReleasesByKeyword` verifies case-insensitive substring matching at the model layer.
- New integration tests verify both HTML and JSON endpoints filter correctly for `?q=v1` and return unfiltered results for empty `q`.
- Manual verification on a repo with multiple tags confirms pagination preserves the query and that uppercase/lowercase queries return identical results.
