# 08 — Search & Discovery

## 1. Repository Search

### What
Find repositories across the Gitea instance using filters and sorting. Global and repo-scoped.

### Search Modes
- Fuzzy: Damerau-Levenshtein distance (default)
- Exact: exact keyword match

### Filters
- Topic, language, stars, forks, size, license
- User/org (uid), archived, private, template
- Mode: source, fork, mirror, collaborative

### Sort Options
- Alphabetical (alpha/reversealpha)
- Created date (created/oldest)
- Updated date (updated/recentupdate/leastupdate)
- Size (size/reversesize)
- ID (id/id_reverse)
- Stars, forks

### UI Routes
- GET /explore/repos — Explore repos
- GET /explore/repos/sitemap-{idx}.xml — Sitemap

### API Endpoints
- GET /repos/search — Search repos (supports all filters)

### Config
- [indexer] REPO_INDEXER_ENABLED — enable repo indexing
- [ui] SEARCH_REPO_DESCRIPTION — include descriptions
- [ui] EXPLORE_DEFAULT_SORT — default sort order
- [ui] ONLY_SHOW_RELEVANT_REPOS — relevant repos only

---

## 2. Code Search

### What
Full-text search through source code across accessible repositories.

### Indexer Backends
- Bleve (embedded, default)
- Elasticsearch (external)
- Meilisearch (external)

### Search Features
- Fuzzy and exact matching
- Language filter
- Syntax highlighting in results
- Context lines display
- File inclusion/exclusion patterns
- Vendored file exclusion

### UI Routes
- GET /explore/code — Global code search
- GET /{owner}/{repo}/search — Repo-scoped code search

### Config
- [indexer] REPO_INDEXER_ENABLED — enable code indexing
- [indexer] REPO_INDEXER_TYPE — backend (bleve/elasticsearch/meilisearch)
- [indexer] REPO_INDEXER_PATH — Bleve index path
- [indexer] REPO_INDEXER_CONN_STR — external connection string
- [indexer] REPO_INDEXER_NAME — index name
- [indexer] REPO_INDEXER_INCLUDE — include patterns
- [indexer] REPO_INDEXER_EXCLUDE — exclude patterns
- [indexer] REPO_INDEXER_EXCLUDE_VENDORED — exclude vendored
- [indexer] MAX_FILE_SIZE — max file size to index (default 1MB)
- [indexer] STARTUP_TIMEOUT — indexer startup timeout

### Constraints
- Requires enabled indexer for advanced features
- Large repos may have indexing delays
- External indexers need network connectivity

---

## 3. Issue/PR Search

### What
Search issues and pull requests across repos with advanced filters.

### Filters
- State: open, closed, all
- Labels (by name, comma-separated)
- Milestones (by name)
- Assignee: assigned to current user
- Created by current user
- Mentioned current user
- Review requested from current user
- Reviewed by current user
- Type: issues or pulls
- Owner/team filter
- Time range: since/before dates

### Sort Options
- Oldest, recent update, least update
- Most comment, least comment
- Priority
- Near/far due date
- Priority repo

### UI Routes
- GET /issues — User's issues across repos
- GET /pulls — User's PRs across repos
- GET /{owner}/{repo}/issues — Repo issues
- GET /{owner}/{repo}/pulls — Repo PRs

### API Endpoints
- GET /repos/issues/search — Global issue/PR search

### Config
- [indexer] ISSUE_INDEXER_TYPE — backend
- [indexer] ISSUE_INDEXER_PATH — Bleve path
- [indexer] ISSUE_INDEXER_CONN_STR — external connection
- [indexer] ISSUE_INDEXER_NAME — index name

---

## 4. User/Org Search

### What
Find users and organizations within the instance.

### Filters
- Keyword (username, full name, email)
- Type: individual or organization
- Visibility: public, limited, private
- Status: active, admin, 2FA enabled

### Sort Options
- Newest, oldest
- Alphabetical (forward/reverse)
- Last login (forward/reverse)
- Recent update

### UI Routes
- GET /explore/users — User explore
- GET /explore/organizations — Org explore
- GET /explore/users/sitemap-{idx}.xml — Sitemap

### Config
- [service] EXPLORE_DISABLE_USERS_PAGE — disable user explore
- [ui] EXPLORE_PAGING_NUM — page size

### Constraints
- Non-admin: only see public and limited users
- Email search optional/disablable

---

## 5. Global Explore

### What
Curated views of repos, users, orgs, and code across the instance.

### UI Routes
- GET /explore — Redirects to /explore/repos
- GET /explore/repos — Repos
- GET /explore/users — Users
- GET /explore/organizations — Orgs
- GET /explore/code — Code
- GET /explore/topics/search — Topic search

### Features
- Relevant repos filtering
- Sitemap generation (SEO)
- Multiple sort options
- Configurable paging

### Config
- [ui] EXPLORE_PAGING_NUM — page size
- [ui] SITEMAP_PAGING_NUM — sitemap page size
- [ui] ONLY_SHOW_RELEVANT_REPOS — filter relevant repos

---

## 6. Indexer System

### What
Pluggable search backend for code, issues, repos, and stats.

### Components
- **Code Indexer**: full-text code search (Bleve/Elasticsearch/Meilisearch)
- **Issue Indexer**: issue/PR search (Bleve/Elasticsearch/Meilisearch)
- **Stats Indexer**: repo statistics tracking

### Indexer Queue
- Worker pool for async updates
- Unique queue for deduplication
- Error handling and retry
- Graceful shutdown

### Rebuild
- Full index rebuild capability
- Incremental rebuild
- Progress tracking

### Config (all in [indexer] section)
- ISSUE_INDEXER_TYPE / ISSUE_INDEXER_PATH / ISSUE_INDEXER_CONN_STR / ISSUE_INDEXER_NAME
- REPO_INDEXER_ENABLED / REPO_INDEXER_TYPE / REPO_INDEXER_PATH / REPO_INDEXER_CONN_STR
- REPO_INDEXER_INCLUDE / REPO_INDEXER_EXCLUDE / REPO_INDEXER_EXCLUDE_VENDORED
- MAX_FILE_SIZE / STARTUP_TIMEOUT

---

## 7. Sitemap

### What
Auto-generated XML sitemaps for SEO covering repos, users, orgs.

### Routes
- /explore/repos/sitemap-{idx}.xml
- /explore/users/sitemap-{idx}.xml
- /explore/organizations/sitemap-{idx}.xml

### Features
- On-demand generation
- Paginated for large instances
- Standard XML format
- Last modified timestamps

### Config
- [ui] SITEMAP_PAGING_NUM — items per sitemap file

### Constraints
- Max 50,000 URLs per file
- Max 50MB per file
- Respects access permissions
