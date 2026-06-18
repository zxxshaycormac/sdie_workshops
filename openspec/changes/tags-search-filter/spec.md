## ADDED Requirements

### 1. Tag Listing Search

**User Story:** As a developer browsing a repository with many tags, I want to filter the Tags page by keyword so that I can quickly locate a specific tag without scanning the full list.

#### Ubiquitous Requirements (Tag Listing Filter Support)

- **PKG-04-006:** `The system shall support filtering the standalone tag list at /{owner}/{repo}/tags by a substring match against tag_name.`

#### Event-Driven Requirements (Tag Listing Filter Workflow)

- **PKG-04-108:** `When a user submits a GET request to /{owner}/{repo}/tags with a non-empty q query parameter, the system shall return only those tags whose tag_name contains the value of q as a substring.`
- **PKG-04-109:** `When a user submits a GET request to /{owner}/{repo}/tags with an empty or absent q query parameter, the system shall return the unfiltered paginated tag list.`
- **PKG-04-110:** `When the system renders the tag list pagination controls, the system shall preserve the current q value across all page links.`
- **PKG-04-111:** `When the system renders the tag list pagination controls, the system shall compute the total page count from the filtered result set, not from the repository's total tag count.`

#### Unwanted Behaviour Requirements (Tag Listing Filter Edge Cases)

- **PKG-04-304:** `If the q parameter matches no tag_name in the repository, then the system shall render the tag list page with an empty result set, the search box still visible and pre-filled with q, and a no-match message using the existing search.no_results locale key.`

---

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| `q` absent | Full paginated list, search box rendered empty (today's behavior) |
| `q` empty string (`?q=`) | Treated as absent — full paginated list |
| `q` contains SQL LIKE wildcards (`%`, `_`) | Interpreted as wildcards in the LIKE clause; consistent with Branches page behavior |
| `q` contains characters requiring URL encoding | Standard query-string decoding; no special handling |
| `q` matches no tags | Page renders with empty list area, search box still visible with `q` value, and `search.no_results` message; HTTP 200 (not 404) |
| User navigates to page 2 of a filtered result | `q` preserved in pagination link; page 2 of the filtered set is returned |
| User lands on `/tags` from `/{owner}/{repo}` ref selector | No `q` set; full paginated list as today |

---

## Success Criteria

- A repository with 50 tags, where 5 contain the substring `v1.2`, returns exactly those 5 when `?q=v1.2` is submitted on `/tags`.
- The search box on `/tags` visually matches the search box on `/branches` (same partial, same placeholder family).
- Clicking "Next" on a filtered `/tags?q=v1.2` result preserves the `q=v1.2` parameter in the URL.
- On a 50-tag repo filtered to 5 matches with page size 10, the pagination control shows 1 page (not 5).
- Empty/absent `q` produces the same SQL conditions as today's query (the new `Keyword != ""` guard adds nothing).
- Total-tag count in the page header tab remains the repo total, not the filtered count.
