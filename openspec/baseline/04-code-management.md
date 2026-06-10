# 04 — Code Management

## 1. Repository CRUD

### What
Full lifecycle management of Git repositories: creation, deletion, transfer, rename/redirect, forking, template repos, visibility, settings, archiving.

### Behaviors
- Create repo: user-scoped or org-scoped
- Generate from template: creates new repo from template repo
- Visibility: Public, Private, Limited
- Trust models: default, committer, collaborator, collaborator-committer
- Archive repo: makes it read-only (no pushes, issues, PRs)
- Transfer repo: change owner (user or org)
- Rename repo: creates redirect from old name
- Push-to-create: create repo by pushing to a nonexistent path

### UI Routes
- GET /repo/create — Create repo page
- POST /repo/create — Create repo
- GET /{owner}/{repo}/settings — Settings
- POST /{owner}/{repo}/settings — Update settings
- POST /{owner}/{repo}/settings/delete — Delete repo
- POST /{owner}/{repo}/settings/transfer — Transfer repo
- POST /{owner}/{repo}/settings/archive — Archive/unarchive

### API Endpoints
- GET /repos/search — Search repos
- POST /user/repos — Create user repo
- POST /orgs/{org}/repos — Create org repo
- POST /repos/{owner}/{template}/generate — Generate from template
- GET /repos/{owner}/{repo} — Get repo
- GET /repositories/{id} — Get repo by ID
- PATCH /repos/{owner}/{repo} — Edit repo
- DELETE /repos/{owner}/{repo} — Delete repo
- POST /repos/{owner}/{repo}/transfer — Transfer

### Config
- [repository] ROOT — root path for repo storage
- [repository] DEFAULT_BRANCH — default branch name
- [repository] DEFAULT_PRIVATE — default visibility
- [repository] FORCE_PRIVATE — force all repos private
- [repository] OBJECT_FORMAT — sha1 or sha256
- [repository] ALLOW_PUSH_CREATE_USER — push-to-create for users
- [repository] ALLOW_PUSH_CREATE_ORG — push-to-create for orgs
- [repository] ENABLE_PUSH_CREATE_USER — enable push creation
- [repository] MAX_CREATION_LIMIT — max repos per user

### Constraints
- Repo names: alphanumeric + dots + dashes, no reserved names (., .., -)
- Reserved patterns: *.git, *.wiki, *.rss, *.atom
- Archived repos cannot receive pushes, issues, or PRs
- Fork visibility syncs with base repo

---

## 2. Git Operations

### What
Core Git protocol support: clone, push, pull via Smart HTTP and SSH. Git ref management, blame, tree/blob operations, archive downloads.

### Behaviors
- Smart HTTP protocol (git-upload-pack, git-receive-pack)
- SSH protocol (built-in SSH server)
- Archive downloads: zip, tar.gz
- Raw file download
- HTTP basic auth and token auth for git operations

### Config
- [repository] DISABLE_HTTP_GIT — disable HTTP git ops
- [server] SSH_DOMAIN — SSH clone URL domain
- [server] SSH_PORT — SSH port
- [server] SSH_LISTEN_PORT — SSH listen port
- [repository] USE_COMPAT_SSH_URI — compatible SSH URI format

### Constraints
- Archived repos are read-only (no push)
- Mirror repos cannot be pushed to directly
- 2FA users must use token for HTTP operations

---

## 3. Branch Management

### What
Branch CRUD, default branch, protected branches and tags with granular rules.

### Behaviors
- Create, rename, delete branches
- Set default branch
- Protected branches:
  - Push allowlist (users/teams)
  - Merge allowlist (users/teams)
  - Required approvals (count, whitelist)
  - Dismiss stale approvals on push
  - Required status checks (contexts)
  - Require signed commits
  - File pattern protection (protected/unprotected patterns)
  - Block merge on rejected reviews
- Protected tags: create/delete allowlists

### UI Routes
- GET /{owner}/{repo}/branches — Branch list
- POST /{owner}/{repo}/branches — Create branch
- DELETE /{owner}/{repo}/branches/{branch} — Delete branch
- GET /{owner}/{repo}/settings/branches/{branch} — Protection settings
- POST /{owner}/{repo}/settings/branches/{branch} — Update protection

### API Endpoints
- GET /repos/{owner}/{repo}/branches — List branches
- POST /repos/{owner}/{repo}/branches — Create branch
- GET /repos/{owner}/{repo}/branches/{branch} — Get branch
- DELETE /repos/{owner}/{repo}/branches/{branch} — Delete branch
- GET /repos/{owner}/{repo}/branches/{branch}/protection — Get protection
- PUT /repos/{owner}/{repo}/branches/{branch}/protection — Set protection
- DELETE /repos/{owner}/{repo}/branches/{branch}/protection — Remove protection
- GET /repos/{owner}/{repo}/tags — List tags
- POST /repos/{owner}/{repo}/tags — Create tag
- DELETE /repos/{owner}/{repo}/tags/{tag} — Delete tag

### Config
- [repository] DEFAULT_BRANCH — default branch name
- [repository.signing] DEFAULT_TRUST_MODEL — trust model for signing

---

## 4. Code Review/Browsing

### What
File browser, blame view, commit history, commit diff, compare view, contributor graphs, find file, rendering (markup, CSV, PDF, code highlighting), line linking.

### Behaviors
- Tree/blob navigation with commit context
- Blame view with line-level authorship
- Commit history with diff, patch, and verification
- Compare view between branches/tags/commits
- Contributor graph, code frequency graph, recent commits graph
- Find file (fuzzy search within repo)
- Rendering: Markdown, AsciiDoc, reStructuredText, Org mode, Jupyter, CSV, PDF
- Code highlighting via Chroma (100+ languages)
- Line linking for easy reference

### UI Routes
- GET /{owner}/{repo}/src/branch/{branch}/{path} — File browser
- GET /{owner}/{repo}/blame/{branch}/{path} — Blame view
- GET /{owner}/{repo}/commit/{sha} — Commit detail
- GET /{owner}/{repo}/compare/{base}...{head} — Compare view
- GET /{owner}/{repo}/graphs/contributors — Contributors graph
- GET /{owner}/{repo}/graphs/code_frequency — Code frequency
- GET /{owner}/{repo}/graphs/commit_activity — Recent commits
- GET /{owner}/{repo}/find/{branch} — Find file

### API Endpoints
- GET /repos/{owner}/{repo}/contents/{path} — Get file content
- GET /repos/{owner}/{repo}/git/blobs/{sha} — Get blob
- GET /repos/{owner}/{repo}/git/trees/{sha} — Get tree
- GET /repos/{owner}/{repo}/git/commits/{sha} — Get commit
- GET /repos/{owner}/{repo}/git/refs — List refs
- GET /repos/{owner}/{repo}/git/refs/{ref} — Get ref
- GET /repos/{owner}/{repo}/commits — List commits
- GET /repos/{owner}/{repo}/commits/{sha} — Get commit
- GET /repos/{owner}/{repo}/compare/{base}...{head} — Compare

### Config
- [ui] MAX_DISPLAY_FILE_SIZE — max file for display
- [repository] RENDERED_MAX_FILE_SIZE — max file for rendering
- [git] MAX_GIT_DIFF_LINES — max diff lines
- [git] MAX_GIT_DIFF_LINE_CHARACTERS — max diff line chars
- [git] MAX_GIT_DIFF_FILES — max diff files

---

## 5. Web Editor

### What
Web-based file create/edit/delete with commit message, branch creation, and PR flow.

### Behaviors
- Create new file, edit existing, delete file
- Commit with message and description
- Choose: commit directly or create new branch
- Auto-create PR when editing on new branch
- Patch/cherry-pick support
- Directory creation
- Upload files via drag-and-drop

### UI Routes
- GET /{owner}/{repo}/new/{branch} — New file form
- GET /{owner}/{repo}/edit/{branch}/{path} — Edit file
- GET /{owner}/{repo}/delete/{branch}/{path} — Delete file
- POST /{owner}/{repo}/upload — Upload files
- GET /{owner}/{repo}/patch — Cherry-pick/patch

### API Endpoints
- POST /repos/{owner}/{repo}/contents/{path} — Create/update file
- DELETE /repos/{owner}/{repo}/contents/{path} — Delete file
- GET /repos/{owner}/{repo}/contents/{path} — Get file

### Config
- [repository.upload] ENABLED — enable uploads
- [repository.upload] TEMP_PATH — temp upload path
- [repository.upload] MAX_SIZE — max upload size

---

## 6. Forking

### What
Create personal copy of a repository. Fork network tracking.

### Behaviors
- Fork to user account or org
- Fork network visualization
- Fork count tracking
- Base repository relationship maintained

### UI Routes
- GET /{owner}/{repo}/fork — Fork page
- POST /{owner}/{repo}/fork — Create fork

### API Endpoints
- POST /repos/{owner}/{repo}/forks — Create fork
- GET /repos/{owner}/{repo}/forks — List forks

### Constraints
- Fork visibility syncs with base repo (private base → private fork)
- Cannot fork to same owner/repo name

---

## 7. Mirroring

### What
Repository mirroring: pull mirrors (sync from remote) and push mirrors (sync to remote).

### Behaviors
- Pull mirror: periodic sync from remote, manual trigger
- Push mirror: sync to remote on push or scheduled
- LFS mirroring support
- Authentication for private remotes
- GitHub OAuth for private mirror sources

### UI Routes
- GET /{owner}/{repo}/settings/mirror — Mirror settings
- POST /{owner}/{repo}/settings/mirror — Configure mirror
- POST /{owner}/{repo}/settings/mirror/sync — Trigger sync

### API Endpoints
- GET /repos/{owner}/{repo}/mirror — Get mirror info
- PATCH /repos/{owner}/{repo}/mirror — Update mirror
- POST /repos/{owner}/{repo}/mirror/sync — Trigger sync

### Config
- [mirror] MIN_INTERVAL — minimum sync interval
- [mirror] ALLOW_LOCAL_NETWORKS — allow local network mirrors
- [mirror] LFS_ENABLED — enable LFS mirroring
- [mirror] GITHUB_OAUTH — GitHub OAuth token

---

## 8. Migration

### What
Import repositories from external platforms: GitHub, GitLab, Gogs, OneDev, Codebase, GitBucket, plain Git.

### Behaviors
- Migrates: git data, issues, PRs, labels, milestones, releases, wiki, comments, avatars
- OAuth support for private sources
- LFS migration
- Rollback on failure
- Migration progress tracking

### UI Routes
- GET /repo/migrate — Migration page
- POST /repo/migrate — Start migration

### API Endpoints
- POST /repo/migrate — Migrate repository

### Config
- [migration] ALLOWED_DOMAINS — allowlist for migration sources
- [migration] MAX_FILE_SIZE — max migration file size

### Constraints
- Source URL validated against allowlist
- Local filesystem migration for privileged users only
- Timeout handling for large migrations

---

## 9. Topics

### What
Repository tags/topics for categorization and discovery.

### Behaviors
- Add/remove topics on repos
- Topic search across instance
- Topic name validation

### UI Routes
- GET /{owner}/{repo}/settings/topics — Manage topics
- POST /{owner}/{repo}/settings/topics — Update topics

### API Endpoints
- GET /repos/{owner}/{repo}/topics — Get topics
- PUT /repos/{owner}/{repo}/topics — Set topics
- GET /topics/search — Search topics

### Constraints
- Max 35 chars, lowercase alphanumeric with dots/dashes
- Pattern: ^[a-z0-9][-.a-z0-9]*$

---

## 10. Collaborators

### What
Per-repository collaborator management with individual permission levels.

### Behaviors
- Add/remove collaborators
- Set permission level per collaborator: Read, Write, Admin
- Team-based access for org repos

### UI Routes
- GET /{owner}/{repo}/settings/collaboration — Manage collaborators
- POST /{owner}/{repo}/settings/collaboration — Add/update collaborator

### API Endpoints
- GET /repos/{owner}/{repo}/collaborators — List collaborators
- PUT /repos/{owner}/{repo}/collaborators/{user} — Add collaborator
- DELETE /repos/{owner}/{repo}/collaborators/{user} — Remove collaborator
- GET /repos/{owner}/{repo}/collaborators/{user}/permission — Get permission

---

## 11. Starring & Watching

### What
Star (bookmark) and watch (notification subscription) for repositories.

### Behaviors
- Star: bookmark repo, star count tracked
- Watch modes: None, Normal (all events), Dont, Auto (auto-watch on push)
- Watch drives notification delivery

### UI Routes
- Star/unstar buttons on repo page
- GET /user/starred — User's starred repos
- GET /user/subscriptions — User's watched repos

### API Endpoints
- PUT /repos/{owner}/{repo}/star — Star repo
- DELETE /repos/{owner}/{repo}/star — Unstar
- PUT /repos/{owner}/{repo}/subscription — Watch
- DELETE /repos/{owner}/{repo}/subscription — Unwatch
- GET /user/starred — List starred
- GET /repos/{owner}/{repo}/subscribers — List watchers

### Config
- [service] AUTO_WATCH_ON_CHANGES — auto-watch on push
- [service] AUTO_WATCH_REPOS — auto-watch repos
