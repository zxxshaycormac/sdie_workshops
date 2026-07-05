# 08 — Search & Discovery

Baseline specification of Gitea's search, indexing, and discovery subsystems. All requirements describe the current (v1.22.x) system behavior.

---

## 1. Repository Search

**User Story:** As a user, I want to search for repositories by keyword and filters so that I can find relevant projects quickly.

### Ubiquitous Requirements (Search Defaults)

- **SRCH-01-001:** `The system shall match repository keywords using database LIKE (substring) queries; there is no fuzzy/exact mode toggle for repository search.`
- **SRCH-01-002:** `The system shall apply repository search filters including topic and primary language. Star count, fork count, and repository size are available as sort orders only (not filters); no license filter exists.`
- **SRCH-01-003:** `The system shall apply ownership filters including user/org UID, archived status, private status, and template status.`
- **SRCH-01-004:** `The system shall apply repository mode filters: source, fork, mirror, and collaborative.`
- **SRCH-01-005:** `The system shall support the following sort options: alphabetical, reverse alphabetical, created date, oldest, updated, recent update, least update, size, reverse size, ID, reverse ID, stars, and forks.`
- **SRCH-01-006:** `The system shall paginate repository search results using the configured EXPLORE_PAGING_NUM page size.`
- **SRCH-01-007:** `The system shall expose repository search via the API endpoint GET /repos/search with support for all filters and sort options.`

### Event-Driven Requirements (Search Execution)

- **SRCH-01-101:** `When a user submits a keyword query on the explore repos page, the system shall return matching repositories filtered by the user's visibility scope.`
- **SRCH-01-102:** `When a user selects a sort option, the system shall reorder the results according to the selected sort criteria.`
- **SRCH-01-103:** `When the setting ONLY_SHOW_RELEVANT_REPOS is enabled, the system shall filter out repositories with low relevance scores from explore results.`
- **SRCH-01-104:** `When the setting SEARCH_REPO_DESCRIPTION is enabled, the system shall include repository descriptions in the search index for matching.`

### Optional Feature Requirements (Repo Indexer)

- **SRCH-01-201:** `Where REPO_INDEXER_ENABLED is true, the system shall index repository names and descriptions for full-text search.`
- **SRCH-01-202:** `Where a default explore sort is configured via EXPLORE_PAGING_DEFAULT_SORT, the system shall use that sort order as the initial ordering on the explore repos page.`

### Unwanted Behaviour Requirements (Search Errors)

- **SRCH-01-301:** `If a user submits a search query without permission to view any matching repositories, then the system shall return an empty result set.`
- **SRCH-01-302:** `If the repository indexer is disabled and a full-text search is attempted, then the system shall fall back to database-based name matching.`

---

## 2. Code Search

**User Story:** As a developer, I want to search through source code across repositories so that I can find specific implementations, patterns, or usages.

### Ubiquitous Requirements (Code Search Properties)

- **SRCH-02-001:** `The system shall support full-text search through source code in repositories accessible to the authenticated user.`
- **SRCH-02-002:** `The system shall support two indexer backends for code search: Bleve (embedded, default) and Elasticsearch (external).`
- **SRCH-02-003:** `The system shall provide a language filter allowing users to restrict code search results to specific programming languages.`
- **SRCH-02-004:** `The system shall display syntax highlighting in code search results based on file language.`
- **SRCH-02-005:** `The system shall display context lines surrounding each search match.`
- **SRCH-02-006:** `The system shall expose code search via GET /explore/code (global) and GET /{owner}/{repo}/search (repo-scoped).`
- **SRCH-02-007:** `The system shall support two search modes for code search: fuzzy (Levenshtein distance, default) and exact (literal match), selected via the fuzzy query parameter.`

### Event-Driven Requirements (Indexing Workflow)

- **SRCH-02-101:** `When a repository is updated with a new commit, the system shall enqueue the affected files for re-indexing in the code indexer.`
- **SRCH-02-102:** `When a user submits a code search query, the system shall query the configured indexer backend and return matching file content with line numbers.`
- **SRCH-02-103:** `When the indexer encounters a file exceeding MAX_FILE_SIZE (default 1MB), the system shall skip indexing that file.`

### Optional Feature Requirements (Indexer Configuration)

- **SRCH-02-201:** `Where REPO_INDEXER_INCLUDE patterns are configured, the system shall only index files matching those glob patterns.`
- **SRCH-02-202:** `Where REPO_INDEXER_EXCLUDE patterns are configured, the system shall skip files matching those glob patterns during indexing.`
- **SRCH-02-203:** `Where REPO_INDEXER_EXCLUDE_VENDORED is enabled, the system shall skip files in vendor directories during indexing.`
- **SRCH-02-204:** `Where REPO_INDEXER_TYPE is set to elasticsearch, the system shall connect to the external Elasticsearch cluster specified by REPO_INDEXER_CONN_STR.`
- **SRCH-02-206:** `The system shall index only repositories whose type matches the configured REPO_INDEXER_REPO_TYPES list (default: sources, forks, mirrors, templates).`

### Unwanted Behaviour Requirements (Code Search Errors)

- **SRCH-02-301:** `If the code indexer is not enabled, then the system shall not display the code search tab and shall return an error for code search API requests.`
- **SRCH-02-302:** `If the external indexer backend is unreachable, then the system shall return an error for code search queries and log the connection failure.`
- **SRCH-02-303:** `If the indexer startup exceeds STARTUP_TIMEOUT, then the system shall log a timeout error and proceed without code search availability.`

---

## 3. Issue/PR Search

**User Story:** As a contributor, I want to search issues and pull requests across repositories with advanced filters so that I can find relevant work items quickly.

### Ubiquitous Requirements (Issue Search Properties)

- **SRCH-03-001:** `The system shall support searching issues and pull requests across repositories accessible to the authenticated user.`
- **SRCH-03-002:** `The system shall support state filters: open, closed, and all.`
- **SRCH-03-003:** `The system shall support filtering by labels (comma-separated label names), milestones (by name), assignee, author, and mention.`
- **SRCH-03-004:** `The system shall support filtering by type to restrict results to issues only or pull requests only.`
- **SRCH-03-005:** `The system shall support time range filters using since and before date parameters.`
- **SRCH-03-006:** `The system shall support the following sort options: oldest, recent update, least update, most comment, least comment, priority, near due date, far due date, and priority repo.`
- **SRCH-03-007:** `The system shall expose issue/PR search via the API endpoint GET /repos/issues/search with all filter and sort parameters.`

### Event-Driven Requirements (Issue Search Workflow)

- **SRCH-03-101:** `When a user navigates to GET /issues, the system shall display all issues across repositories where the user has access, filtered by the selected criteria.`
- **SRCH-03-102:** `When a user navigates to GET /pulls, the system shall display all pull requests across repositories where the user has access, filtered by the selected criteria.`
- **SRCH-03-103:** `When a user applies the "review requested" filter, the system shall show issues and PRs where a code review has been requested from the current user.`
- **SRCH-03-104:** `When a user applies the "reviewed by" filter, the system shall show issues and PRs that the current user has reviewed.`
- **SRCH-03-105:** `When a user applies an owner/team filter, the system shall restrict results to repositories owned by the specified owner or accessible to the specified team.`
- **SRCH-03-106:** `When a user submits an issue query with an empty keyword, the system shall bypass the configured indexer and use the database backend directly to ensure consistent listing.`

### Optional Feature Requirements (Issue Indexer)

- **SRCH-03-201:** `Where ISSUE_INDEXER_TYPE is set to elasticsearch, the system shall use Elasticsearch as the issue search backend.`
- **SRCH-03-202:** `Where ISSUE_INDEXER_TYPE is set to meilisearch, the system shall use Meilisearch as the issue search backend.`
- **SRCH-03-203:** `Where ISSUE_INDEXER_TYPE is set to db, the system shall use database-based issue search without an external indexer.`

### Unwanted Behaviour Requirements (Issue Search Errors)

- **SRCH-03-301:** `If a user searches issues in a repository where they lack read access, then the system shall exclude that repository's issues from results.`
- **SRCH-03-302:** `If an invalid label name is specified in the filter, then the system shall return an empty result set for that label criterion.`
- **SRCH-03-303:** `If the issue indexer backend is unreachable, then the system shall fall back to database-based issue search.`

---

## 4. User/Org Search

**User Story:** As a user, I want to search for other users and organizations so that I can find collaborators and groups to work with.

### Ubiquitous Requirements (User Search Properties)

- **SRCH-04-001:** `The system shall support keyword-based search across usernames, full names, and email addresses of users and organizations.`
- **SRCH-04-002:** `The system shall support filtering by entity type: individual users or organizations.`
- **SRCH-04-003:** `The system shall support filtering by visibility level: public, limited, and private.`
- **SRCH-04-004:** `The system shall support filtering by status: active, admin, and 2FA enabled.`
- **SRCH-04-005:** `The system shall support the following sort options: newest, oldest, alphabetical, reverse alphabetical, last login forward, last login reverse, and recent update.`
- **SRCH-04-006:** `The system shall expose user search via GET /explore/users and organization search via GET /explore/organizations.`

### Event-Driven Requirements (User Search Workflow)

- **SRCH-04-101:** `When a non-admin user searches for users, the system shall return only users with public or limited visibility.`
- **SRCH-04-102:** `When an admin user searches for users, the system shall return users across all visibility levels.`
- **SRCH-04-103:** `When a user applies a keyword filter, the system shall match against username, full name, and (where permitted) email address.`

### Optional Feature Requirements (User Search Configuration)

- **SRCH-04-201:** `Where [service].DISABLE_USERS_PAGE is enabled, the system shall disable the user explore page and redirect GET /explore/users to the repositories explore page.`

### Unwanted Behaviour Requirements (User Search Errors)

- **SRCH-04-301:** `If a non-admin user attempts to search for private users, then the system shall exclude private users from the results.`
- **SRCH-04-302:** `If email search is disabled or the requesting user lacks permission, then the system shall exclude email addresses from search matching.`

---

## 5. Global Explore

**User Story:** As a visitor, I want to browse curated views of repositories, users, organizations, and code across the instance so that I can discover interesting content.

### Ubiquitous Requirements (Explore Properties)

- **SRCH-05-001:** `The system shall provide explore tabs for repositories, users, organizations, and code.`
- **SRCH-05-002:** `The system shall paginate explore results using the configured EXPLORE_PAGING_NUM page size.`
- **SRCH-05-003:** `The system shall expose topic search via GET /explore/topics/search.`

### Event-Driven Requirements (Explore Navigation)

- **SRCH-05-101:** `When a visitor navigates to GET /explore, the system shall redirect to GET /explore/repos.`
- **SRCH-05-102:** `When a visitor navigates to GET /explore/repos, the system shall display a paginated list of public repositories sorted by the configured default sort.`
- **SRCH-05-103:** `When a visitor navigates to GET /explore/users, the system shall display a paginated list of public users.`
- **SRCH-05-104:** `When a visitor navigates to GET /explore/organizations, the system shall display a paginated list of public organizations.`
- **SRCH-05-105:** `When a visitor navigates to GET /explore/code, the system shall display the global code search interface if the code indexer is enabled.`

### Optional Feature Requirements (Explore Configuration)

- **SRCH-05-201:** `Where ONLY_SHOW_RELEVANT_REPOS is enabled, the system shall filter the explore repos view to show only repositories deemed relevant based on activity and popularity.`
- **SRCH-05-202:** `Where the code indexer is disabled, the system shall hide the code explore tab from the navigation.`

### Unwanted Behaviour Requirements (Explore Errors)

- **SRCH-05-301:** `If [service].DISABLE_USERS_PAGE is enabled and a visitor navigates to GET /explore/users, then the system shall redirect to GET /explore/repos.`
- **SRCH-05-302:** `If an unauthenticated visitor accesses explore pages, then the system shall show only public resources.`

---

## 6. Indexer System

**User Story:** As an administrator, I want pluggable indexer backends with reliable queue processing so that search indexes stay up-to-date without degrading system performance.

### Ubiquitous Requirements (Indexer Properties)

- **SRCH-06-001:** `The system shall provide separate indexer components for code (REPO_INDEXER), issues (ISSUE_INDEXER), and repository statistics (STATS_INDEXER); the STATS_INDEXER is DB-only with no backend choice and no configuration section.`
- **SRCH-06-002:** `The system shall support Bleve and Elasticsearch as backends for the code indexer, and Bleve, Elasticsearch, Meilisearch, and db as backends for the issue indexer.`
- **SRCH-06-003:** `The system shall use a worker pool with async queue processing for indexer updates.`
- **SRCH-06-004:** `The system shall deduplicate indexer queue entries to avoid redundant indexing operations.`

### Event-Driven Requirements (Indexer Lifecycle)

- **SRCH-06-101:** `When the system starts, the system shall initialize all enabled indexers and wait up to STARTUP_TIMEOUT for each to become ready.`
- **SRCH-06-102:** `When a repository update event occurs, the system shall enqueue the relevant files or issues for re-indexing.`
- **SRCH-06-103:** `When an indexer queue entry fails to process, the system shall log the error and retry the entry.`
- **SRCH-06-104:** `When the system shuts down, the system shall drain the indexer queues gracefully before terminating.`
- **SRCH-06-105:** `When an indexer is initialized for the first time (no existing index), the system shall automatically populate the index from scratch via populateRepoIndexer (code/stats) or PopulateIssueIndexer (issues); no admin reindex UI exists.`
- **SRCH-06-106:** `When repository or issue update events occur after initial population, the system shall incrementally re-index only the affected content via the indexer queue.`

### State-Driven Requirements (Indexer Health)

- **SRCH-06-701:** `While an external indexer backend is unreachable, the system shall queue indexing operations for retry and return degraded search results.`
- **SRCH-06-702:** `While initial index population is in progress, the system shall track and report the population progress.`

### Optional Feature Requirements (Indexer Configuration)

- **SRCH-06-201:** `Where REPO_INDEXER_ENABLED is false, the system shall skip code indexer initialization and disable code search functionality.`
- **SRCH-06-202:** `Where ISSUE_INDEXER_TYPE is set to bleve, the system shall store the issue index on the local filesystem at ISSUE_INDEXER_PATH.`

### Unwanted Behaviour Requirements (Indexer Errors)

- **SRCH-06-301:** `If the indexer queue exceeds its capacity, then the system shall log a warning and continue accepting operations with delayed indexing.`
- **SRCH-06-302:** `If initial index population is interrupted, then the system shall re-run it automatically on the next indexer initialization.`

---

## 7. Sitemap

**User Story:** As a search engine, I want auto-generated XML sitemaps so that I can discover and index public content on the Gitea instance.

### Ubiquitous Requirements (Sitemap Properties)

- **SRCH-07-001:** `The system shall generate XML sitemaps conforming to the sitemaps.org protocol.`
- **SRCH-07-002:** `The system shall paginate sitemaps using the configured SITEMAP_PAGING_NUM page size.`
- **SRCH-07-003:** `The system shall include last modified timestamps in sitemap entries.`
- **SRCH-07-004:** `The system shall expose sitemaps for repositories and users.`

### Event-Driven Requirements (Sitemap Generation)

- **SRCH-07-101:** `When a request is made to GET /explore/repos/sitemap-{idx}.xml, the system shall generate a sitemap page containing up to SITEMAP_PAGING_NUM public repositories.`
- **SRCH-07-102:** `When a request is made to GET /explore/users/sitemap-{idx}.xml, the system shall generate a sitemap page containing up to SITEMAP_PAGING_NUM public users.`

### Unwanted Behaviour Requirements (Sitemap Constraints)

- **SRCH-07-301:** `If a sitemap page would contain more than 50,000 URLs, then the system shall split the sitemap across multiple pages.`
- **SRCH-07-302:** `If a sitemap page would exceed 50MB in size, then the system shall split the sitemap across multiple pages.`
- **SRCH-07-303:** `If a resource is not publicly accessible, then the system shall exclude that resource from all sitemaps.`

---

## 8. Search Candidates Autocomplete (SRCH-08)

**User Story:** As a user, I want a fast typeahead search for other users when composing @mentions or selecting collaborators so that I can find the right person without leaving the form.

### Ubiquitous Requirements (Autocomplete Properties)

- **SRCH-08-001:** `The system shall expose a GET /user/search_candidates endpoint that performs a single global SearchUsers query with no repository context and no mode distinction (no separate mention/assignee modes).`
- **SRCH-08-002:** `The system shall search only active individual users (UserTypeIndividual, IsActive=true); organizations and bots are excluded.`
- **SRCH-08-003:** `The system shall return matching user ID, login name, full name, and avatar URL in the response payload via convert.ToUsers.`
- **SRCH-08-004:** `The system shall apply BuildCanSeeUserCondition visibility scoping: admins see all users; non-restricted authenticated users see public and limited-visibility users plus orgs they belong to; restricted users see only orgs they belong to and themselves; anonymous users see only public users.`

### Event-Driven Requirements (Autocomplete Workflow)

- **SRCH-08-101:** `When a user types into a typeahead field, the system shall query matching users and return up to setting.UI.MembersPagingNum results (default 20); there is no minimum character threshold.`
- **SRCH-08-102:** `When an admin requests candidates, the system shall return all matching active individual users regardless of visibility level.`
- **SRCH-08-103:** `When a non-admin requests candidates, the system shall exclude private-visibility users and private orgs they are not a member of from the results.`

### Optional Feature Requirements (Autocomplete Configuration)

- **SRCH-08-201:** `Where a user's visibility is set to private, the system shall exclude that user from candidate results for all viewers except the user themselves and admins.`

### Unwanted Behaviour Requirements (Autocomplete Errors)

- **SRCH-08-301:** `If a query matches more than setting.UI.MembersPagingNum users, then the system shall return only the first page of results with no incomplete flag in the response payload.`
- **SRCH-08-302:** `If the query string is empty, then the system shall still query the database and return matching users (no short-query suppression).`

---

## Business Rules

- **BR-08-001:** Repository search applies visibility scoping based on the authenticated user's permissions
- **BR-08-002:** Code search requires an enabled and initialized code indexer
- **BR-08-003:** Code indexer applies include/exclude patterns before the MAX_FILE_SIZE check
- **BR-08-004:** Code indexer supports only Bleve and Elasticsearch; issue indexer supports Bleve, Elasticsearch, Meilisearch, and db
- **BR-08-005:** Non-admin users can only discover public and limited-visibility users in search
- **BR-08-006:** Explore pages respect the [service].DISABLE_USERS_PAGE setting (redirects to /explore/repos, not 404) and ONLY_SHOW_RELEVANT_REPOS
- **BR-08-007:** Indexer queues use deduplication to prevent redundant indexing of the same content
- **BR-08-008:** Sitemaps are generated on-demand (not pre-built) and respect access permissions
- **BR-08-009:** Each sitemap file contains at most 50,000 URLs and must not exceed 50MB
- **BR-08-010:** Fuzzy code search uses Levenshtein distance as the default matching algorithm
- **BR-08-011:** External indexer backends require network connectivity and valid connection strings
- **BR-08-012:** Indexer startup timeout is configurable and applies per-indexer on system boot

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Code search with indexer disabled | Hide code search tab, return error for API requests |
| External indexer unreachable during search | Return error to user, log connection failure |
| External indexer unreachable during indexing | Queue operations for retry, log error |
| File exceeds MAX_FILE_SIZE during indexing | Skip file, do not index |
| Indexer startup exceeds STARTUP_TIMEOUT | Log timeout, proceed without that indexer |
| Non-admin searches for private users | Exclude private users from results silently |
| Initial index population interrupted mid-way | Re-runs automatically on next indexer initialization |
| Sitemap page exceeds 50,000 URLs | Split across multiple sitemap pages |
| Empty keyword submitted for search | Return all results with applied filters and default sort |
| Invalid label name in issue filter | Return empty result set for that label criterion |
| Code search with include/exclude conflict | Apply exclude after include; exclude takes precedence |
| Expired or invalid OAuth2 token used for search API | Return 401 Unauthorized |

## Success Criteria

- Repository search results return within 2 seconds for standard keyword queries
- Code search results include syntax-highlighted context lines with accurate line numbers
- Issue/PR cross-repository search returns results within 3 seconds for typical filter combinations
- User/org search results return within 1 second for keyword queries
- Explore pages load within 2 seconds with default pagination
- Indexer queue processes updates without blocking the main application thread
- Full index population completes without data loss and reports progress
- Sitemaps conform to the sitemaps.org XML protocol specification
- Fuzzy code search matches approximate typos within Levenshtein distance threshold
- Search result pagination provides consistent ordering across page boundaries
