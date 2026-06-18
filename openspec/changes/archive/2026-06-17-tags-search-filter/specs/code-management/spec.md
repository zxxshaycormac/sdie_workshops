## ADDED Requirements

### 12. Repository Tags Page Browsing

**User Story:** As a maintainer, I want to search the Tags page by name so that I can quickly locate a specific release tag in a repository that has many tags.

#### Ubiquitous Requirements (Tags Page Properties)

- **REPO-12-001:** `The system shall expose a tag-name keyword filter on the repository Tags page via the q query parameter.`
- **REPO-12-002:** `The system shall perform tag-name filtering as a case-insensitive substring match against the lower_tag_name column.`
- **REPO-12-003:** `The system shall render a search input on the Tags page above the tag list using the shared search/combo template.`
- **REPO-12-004:** `The system shall preserve an empty-keyword query as a no-op, returning the full unfiltered tag list.`

#### Event-Driven Requirements (Tags Page Workflow)

- **REPO-12-101:** `When the user submits the Tags page search form with a keyword, the system shall issue an HTTP GET to the Tags page URL with the q parameter set to the keyword.`
- **REPO-12-102:** `When the user paginates a filtered Tags result set, the system shall carry the active keyword through every page link via the SetDefaultParams pagination contract.`
- **REPO-12-103:** `When the user submits an empty keyword on the Tags page, the system shall return the full tag list with no LIKE clause applied.`

#### State-Driven Requirements (Filtered Pagination)

- **REPO-12-701:** `While a non-empty keyword is active on the Tags page, the system shall compute the pager's total count from the filtered result set rather than the unfiltered repository tag count.`
- **REPO-12-702:** `While the filtered result set fits within a single page, the system shall suppress the pagination widget entirely.`

#### Unwanted Behaviour Requirements (Tags Page Edge Cases)

- **REPO-12-301:** `If the user-supplied keyword contains SQL LIKE wildcard characters (percent or underscore), the system shall treat those characters as literal substrings rather than pattern operators.`
- **REPO-12-302:** `If a non-empty keyword matches no tags, the system shall render the search form and the empty-state message without a pagination widget.`
- **REPO-12-303:** `If the page query parameter is missing or non-positive on the Tags page, the system shall normalize it to page 1 before issuing the underlying query.`

---

## Business Rules

- **BR-12-001:** Tag-name filtering on the Tags page applies only to the current repository.
- **BR-12-002:** Tag-name filtering respects the caller's existing repository read permissions — the q parameter never broadens access.
- **BR-12-003:** The Releases page handler (`/{owner}/{repo}/releases`) is unaffected by the Tags page keyword parameter and continues to render all releases regardless of any q value.

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|-----------------|
| Empty or missing `q` parameter | `Keyword == ""` → LIKE clause skipped → returns full unfiltered tag list (backward compatible with pre-feature behavior). |
| `q` contains `%` or `_` (SQL wildcards) | `builder.Like` parameterizes and escapes; characters treated as literals. Searching `100%` does not match every tag. |
| `q` contains whitespace | `ctx.FormString` preserves whitespace; passed verbatim to LIKE. No trimming. |
| Mixed-case `q` (e.g., `V1.2`) | Handler passes raw keyword to model; `strings.ToLower` applied before LIKE on `lower_tag_name`. `V1.2` matches `v1.2.3`. |
| Filtered result count exceeds `limit` | Pagination applies; `SetDefaultParams(ctx)` carries `q` through every page link so the keyword is preserved across pages. |
| Zero results on a non-empty `q` | Template renders the existing empty state. Search form remains visible because it lives outside the `{{if .Releases}}` block. |
| Missing or non-positive `page` parameter | Handler normalizes `Page <= 0` to `1` before issuing `FindAndCount`, matching the Branches page convention. |
| Releases page receives a stray `q` | Releases handler does not read `q`; `Keyword` defaults to empty string; LIKE clause skipped; Releases behavior unchanged. |

## Success Criteria

- Visiting `/{owner}/{repo}/tags?q=<known-tag-name-fragment>` returns only tags whose `lower_tag_name` contains the lowercased keyword.
- The search input is visible on every Tags page render, including zero-result and zero-keyword states.
- Paginating a filtered result set preserves the keyword across pages and never produces empty interior pages.
- The Releases page (`/{owner}/{repo}/releases`) renders identically to its pre-change behavior.
- All existing tag/release integration tests continue to pass without modification.
