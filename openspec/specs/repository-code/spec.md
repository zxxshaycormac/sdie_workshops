# 03 — Repository & Code

Baseline specification of Gitea's repository lifecycle, git operations, branching, code browsing, web editor, mirroring, migration, topics, starring/watching, wiki (as a git repository), repo templates, repo activity, code frequency, repo search, and code search subsystems. All requirements describe the current (v1.22.x) system behavior.

> **Boundary rule:** Features that manage the source code artifact and its git lifecycle.
> **Codebase anchor:** `services/repository`, `modules/git`, `services/mirror`, `services/migrations`, `services/wiki`.
> **Cross-domain references:** Required approvals → Domain 04 (Collaboration). Required status checks → Domain 05 (CI/CD). Push access restrictions → Domain 02 (Access Control). Releases/packages → Domain 06 (Packages & Releases). Collaborators/deploy keys → Domain 02. Issues/PRs/projects → Domain 04. Webhooks → Domain 05.

---

## 1. Repository CRUD (REPO-01)

**User Story:** As a user, I want to create and manage repositories so that I can organize and share my code projects.

### Ubiquitous Requirements (Repository Properties)

- **REPO-01-001:** `The system shall assign each repository a unique combination of owner name and repository name.`
- **REPO-01-002:** `The system shall support two repository visibility levels: public and private (IsPrivate bool). Three-level visibility (Public/Limited/Private) is a user/organization attribute that gates discovery and membership, not a repository attribute.`
- **REPO-01-003:** `The system shall support four commit signature trust models: default, committer, collaborator, and collaborator-committer.`
- **REPO-01-004:** `The system shall support two object formats for repository storage: SHA-1 and SHA-256.`
- **REPO-01-005:** `The system shall store repositories under a configurable root path.`
- **REPO-01-006:** `The system shall assign a configurable default branch name to newly created repositories.`
- **REPO-01-007:** `The system shall reject repository names that match reserved names (".", "..", "-") or reserved patterns ("*.git", "*.wiki", "*.rss", "*.atom").`
- **REPO-01-008:** `The system shall restrict repository names to alphanumeric characters, dots, and dashes.`
- **REPO-01-009:** `The system shall support adoption of unadopted repositories (pre-existing bare Git directories on disk under ROOT without a database record) by administrators, importing them into the database.`

### Event-Driven Requirements (Repository Lifecycle)

- **REPO-01-101:** `When a user creates a repository, the system shall initialize the repository under the user's namespace or a specified organization's namespace.`
- **REPO-01-102:** `When a user creates a repository from a template, the system shall copy the template's Git data, branches, and optionally its topics, Git hooks, and webhooks.`
- **REPO-01-103:** `When a user deletes a repository, the system shall remove the repository and all associated data (issues, pull requests, releases, wiki, settings).`
- **REPO-01-104:** `When a user transfers a repository to another owner, the system shall update the repository namespace and adjust all access permissions.`
- **REPO-01-105:** `When a user renames a repository, the system shall create a persistent redirect from the old repository path to the new path.`
- **REPO-01-106:** `When a user archives a repository, the system shall set the repository to read-only mode, preventing pushes, new issues, and new pull requests.`
- **REPO-01-107:** `When a user unarchives a repository, the system shall restore full read-write functionality.`
- **REPO-01-108:** `When a user pushes to a nonexistent repository path and push-to-create is enabled, the system shall create a new repository at that path.`

### Optional Feature Requirements (Repository Configuration)

- **REPO-01-201:** `Where FORCE_PRIVATE is enabled, the system shall force all new repositories to be private regardless of user selection.`
- **REPO-01-202:** `Where MAX_CREATION_LIMIT is configured for a user, the system shall prevent the user from creating repositories beyond that limit.`
- **REPO-01-203:** `Where ALLOW_PUSH_CREATE_USER is enabled, the system shall allow individual users to create repositories by pushing to a nonexistent path.`
- **REPO-01-204:** `Where ALLOW_PUSH_CREATE_ORG is enabled, the system shall allow organization members to create repositories by pushing to a nonexistent path under the organization.`

### Unwanted Behaviour Requirements (Repository Errors)

- **REPO-01-301:** `If a user attempts to create a repository with a reserved or invalid name, then the system shall reject the creation with a validation error.`
- **REPO-01-302:** `If a user attempts to push to an archived repository, then the system shall reject the push with a read-only error.`
- **REPO-01-303:** `If a user attempts to create a repository beyond their creation limit, then the system shall reject the creation with a limit-exceeded error.`
- **REPO-01-304:** `If a user transfers a repository to an owner where a repository with the same name already exists, then the system shall reject the transfer.`

---

## 2. Git Operations (REPO-02)

**User Story:** As a developer, I want to clone, push, and pull repositories using standard Git protocols so that I can work with code locally and synchronize changes.

> **Cross-reference:** Git transport (HTTP smart protocol, SSH git commands) → Domain 07 (Platform Services) §Git Transport. SSH server details → Domain 07.

### Ubiquitous Requirements (Protocol Support)

- **REPO-02-001:** `The system shall support Git operations over HTTP and SSH protocols (see Domain 07 for SSH server and Smart HTTP protocol details).`
- **REPO-02-002:** `The system shall support HTTP basic authentication and token authentication for Git operations over HTTP.`
- **REPO-02-003:** `The system shall provide archive downloads in zip, tar.gz, and bundle formats for any branch, tag, or commit.`
- **REPO-02-004:** `The system shall provide raw file downloads for individual files in the repository.`
- **REPO-02-005:** `The system shall generate clone URLs for both HTTP and SSH protocols on the repository page.`
- **REPO-02-006:** `The system shall expose each repository's wiki as a standalone Git repository cloneable via the .wiki suffix on the repository path (over HTTP and SSH).`

### Event-Driven Requirements (Git Protocol Workflow)

- **REPO-02-101:** `When a client initiates a git clone over HTTP, the system shall serve the repository via the Smart HTTP protocol (see Domain 07 for full protocol behavior).`
- **REPO-02-102:** `When a client initiates a git push over HTTP, the system shall authenticate the user and accept the push if authorized.`
- **REPO-02-103:** `When a client initiates a Git operation over SSH, the system shall authenticate the user via their registered SSH key (see Domain 07 for SSH server behavior).`
- **REPO-02-104:** `When a user requests an archive download, the system shall generate and serve the archive for the specified ref.`
- **REPO-02-105:** `When a user with 2FA enabled attempts HTTP Git operations, the system shall require a personal access token instead of a password.`

### Optional Feature Requirements (Protocol Configuration)

- **REPO-02-201:** `Where DISABLE_HTTP_GIT is enabled, the system shall reject all Git operations over HTTP while keeping the web UI functional.`
- **REPO-02-202:** `Where USE_COMPAT_SSH_URI is enabled, the system shall generate SSH clone URLs in the compatible URI format (ssh://user@host/path).`

### Unwanted Behaviour Requirements (Git Operation Errors)

- **REPO-02-301:** `If a user pushes to an archived repository, then the system shall reject the push.`
- **REPO-02-302:** `If a user pushes to a mirror repository, then the system shall reject the push (mirror repos are read-only).`
- **REPO-02-303:** `If a user with 2FA enabled provides a password for HTTP Git operations, then the system shall reject the authentication.`
- **REPO-02-304:** `If an SSH key is not registered in the system, then the system shall reject the SSH connection (see Domain 07 for SSH authentication details).`

---

## 3. Branch Management (REPO-03)

**User Story:** As a maintainer, I want to manage branches and configure protection rules so that I can enforce quality standards on critical branches.

> **Cross-reference (Protected Branch):** Protected branch configuration lives here as a git concept, but its rules reference other domains: Required approvals → Domain 04 (Collaboration). Required status checks → Domain 05 (CI/CD). Push access restrictions → Domain 02 (Access Control).

### Ubiquitous Requirements (Branch Properties)

- **REPO-03-001:** `The system shall allow each repository to have a configurable default branch.`
- **REPO-03-002:** `The system shall support creating, renaming, and deleting branches.`
- **REPO-03-003:** `The system shall support protected tags with configurable create and delete allowlists.`
- **REPO-03-004:** `The system shall provide tag browsing and listing, including tag metadata (commit, signature, message) via the releases/tags UI and API.`

### Event-Driven Requirements (Branch Lifecycle)

- **REPO-03-101:** `When a user creates a branch, the system shall create it from the specified source ref.`
- **REPO-03-102:** `When a user deletes a branch, the system shall remove the branch reference and optionally prune unreferenced objects.`
- **REPO-03-103:** `When a user changes the default branch, the system shall update the repository's HEAD reference to the new default.`
- **REPO-03-104:** `When a user creates a protected branch rule, the system shall enforce push and merge restrictions on that branch.`
- **REPO-03-105:** `When a user deletes a protected branch rule, the system shall remove all associated restrictions immediately.`

### State-Driven Requirements (Protected Branch Enforcement)

- **REPO-03-701:** `While a branch is protected with a push allowlist, the system shall reject pushes from users and teams not on the allowlist.`
- **REPO-03-702:** `While a branch is protected with a merge allowlist, the system shall restrict merge pull requests to users and teams on the allowlist.`
- **REPO-03-703:** `While a branch requires a minimum number of approvals, the system shall prevent merging until the required number of approved reviews is met.`
- **REPO-03-704:** `While a branch requires status checks, the system shall prevent merging until all required status check contexts pass.`
- **REPO-03-705:** `While a branch requires signed commits, the system shall reject any push containing unsigned commits.`
- **REPO-03-706:** `While stale approval dismissal is enabled on a protected branch, the system shall dismiss existing approvals when new commits are pushed to the pull request.`
- **REPO-03-707:** `While block-on-rejected-reviews is enabled on a protected branch, the system shall prevent merging when a pull request has any rejected review.`
- **REPO-03-708:** `While file pattern protection is configured, the system shall enforce or relax protection rules for matching file paths.`

### Unwanted Behaviour Requirements (Branch Errors)

- **REPO-03-301:** `If a user without push allowlist access attempts to push to a protected branch, then the system shall reject the push.`
- **REPO-03-302:** `If a user attempts to delete a protected branch, then the system shall reject the deletion.`
- **REPO-03-303:** `If a user attempts to create a tag matching a protected tag pattern without being on the allowlist, then the system shall reject the tag creation.`

---

## 4. Code Review / Browsing (REPO-04)

**User Story:** As a developer, I want to browse code, view commit history, and compare branches so that I can understand and review the codebase.

### Ubiquitous Requirements (Browsing Properties)

- **REPO-04-001:** `The system shall provide a file tree and blob viewer with commit context for each path.`
- **REPO-04-002:** `The system shall provide a blame view with line-level authorship attribution for any file.`
- **REPO-04-003:** `The system shall provide commit history with diff, patch, and signature verification status.`
- **REPO-04-004:** `The system shall render Markdown, AsciiDoc, reStructuredText, Org mode, and Jupyter notebook files.`
- **REPO-04-005:** `The system shall render CSV files with tabular display and PDF files inline.`
- **REPO-04-006:** `The system shall apply syntax highlighting via Chroma for over 100 programming languages.`
- **REPO-04-007:** `The system shall support line linking via URL fragment for direct reference to specific lines in a file.`

### Event-Driven Requirements (Browsing Workflow)

- **REPO-04-101:** `When a user navigates to a file path in a branch, the system shall display the file tree or blob content at that path.`
- **REPO-04-102:** `When a user requests a blame view, the system shall display line-by-line authorship with commit references.`
- **REPO-04-103:** `When a user navigates to a commit, the system shall display the commit metadata, diff, and signature verification.`
- **REPO-04-104:** `When a user opens a compare view between two refs, the system shall display the commits, changed files, and diff between them.`
- **REPO-04-105:** `When a user accesses the contributor graph, the system shall display contributor statistics and activity charts.`
- **REPO-04-106:** `When a user uses the find-file feature, the system shall perform fuzzy search against filenames in the repository.`
- **REPO-04-107:** `When a user clicks a line number in a file view, the system shall update the URL to include a line-link fragment.`

### State-Driven Requirements (Display Limits)

- **REPO-04-701:** `While a file exceeds the configured MAX_DISPLAY_FILE_SIZE (setting.UI.MaxDisplayFileSize), the system shall display a notice instead of rendering the file content.`
- **REPO-04-702:** *(removed — RENDERED_MAX_FILE_SIZE does not exist; only setting.UI.MaxDisplayFileSize governs display limits)*
- **REPO-04-703:** `While a diff exceeds the configured MAX_GIT_DIFF_LINES, MAX_GIT_DIFF_LINE_CHARACTERS, or MAX_GIT_DIFF_FILES, the system shall truncate the diff display.`

### Unwanted Behaviour Requirements (Browsing Errors)

- **REPO-04-301:** `If a user navigates to a nonexistent path in a repository, then the system shall return a 404 error.`
- **REPO-04-302:** `If a compare view references an invalid base or head ref, then the system shall return an error with a descriptive message.`
- **REPO-04-303:** `If a file is binary and cannot be rendered as text, then the system shall display a binary file notice and offer download.`

---

## 5. Web Editor (REPO-05)

**User Story:** As a contributor, I want to create, edit, and delete files directly in the browser so that I can make quick changes without cloning the repository.

### Ubiquitous Requirements (Editor Capabilities)

- **REPO-05-001:** `The system shall support creating new files, editing existing files, and deleting files through the web interface.`
- **REPO-05-002:** `The system shall require a commit message for every file change operation.`
- **REPO-05-003:** `The system shall support file uploads via drag-and-drop in the web interface.`

### Event-Driven Requirements (Editor Workflow)

- **REPO-05-101:** `When a user creates a new file via the web editor, the system shall commit the file to the target branch or a new branch.`
- **REPO-05-102:** `When a user edits an existing file via the web editor, the system shall create a commit with the changes on the target branch or a new branch.`
- **REPO-05-103:** `When a user deletes a file via the web editor, the system shall create a commit removing the file from the target branch or a new branch.`
- **REPO-05-104:** `When a user chooses to commit to a new branch, the system shall create the branch from the current ref and commit the change.`
- **REPO-05-105:** `When a user commits to a new branch and the repository allows pull requests, the system shall prompt the user to create a pull request.`
- **REPO-05-106:** `When a user uploads files, the system shall commit all uploaded files in a single commit to the target branch or a new branch.`
- **REPO-05-107:** `When a user submits a cherry-pick or patch operation, the system shall apply the specified commit to the target branch.`

### Optional Feature Requirements (Editor Configuration)

- **REPO-05-201:** `Where repository uploads are enabled, the system shall accept file uploads up to the configured maximum size.`
- **REPO-05-202:** `Where directory creation is requested, the system shall create intermediate directories as needed in the repository path.`

### Unwanted Behaviour Requirements (Editor Errors)

- **REPO-05-301:** `If a user attempts to edit a file in an archived repository, then the system shall reject the edit.`
- **REPO-05-302:** `If a user attempts to commit to a protected branch without authorization, then the system shall reject the commit and suggest using a new branch.`
- **REPO-05-303:** `If an uploaded file exceeds the configured maximum size, then the system shall reject the upload.`
- **REPO-05-304:** `If a cherry-pick operation results in merge conflicts, then the system shall report the conflicts and not complete the operation.`

---

## 6. Forking (REPO-06)

**User Story:** As a contributor, I want to fork a repository so that I can work on changes independently and submit them back via pull request.

### Ubiquitous Requirements (Fork Properties)

- **REPO-06-001:** `The system shall maintain a fork relationship between the forked repository and its base repository.`
- **REPO-06-002:** `The system shall track the fork count for each repository.`
- **REPO-06-003:** `The system shall support forking to a user's personal account or to an organization the user belongs to.`

### Event-Driven Requirements (Fork Workflow)

- **REPO-06-101:** `When a user forks a repository, the system shall create a copy of the repository under the user's or organization's namespace.`
- **REPO-06-102:** `When a user forks a private repository, the system shall create the fork as a private repository.`
- **REPO-06-103:** `When a user navigates to the fork list of a repository, the system shall display all known forks in the fork network.`

### State-Driven Requirements (Fork Visibility)

- **REPO-06-701:** `While the base repository is private, the system shall ensure the fork remains private regardless of user preference.`

### Unwanted Behaviour Requirements (Fork Errors)

- **REPO-06-301:** `If a user attempts to fork a repository to a namespace where a repository with the same name already exists, then the system shall reject the fork.`
- **REPO-06-302:** `If a user does not have read access to a repository, then the system shall reject the fork attempt.`
- **REPO-06-303:** `If a user attempts to fork a repository they do not have visibility of, then the system shall return a 404 error.`

---

## 7. Mirroring (REPO-07)

**User Story:** As a repository administrator, I want to set up mirroring so that my repository stays synchronized with an external Git host.

### Ubiquitous Requirements (Mirror Properties)

- **REPO-07-001:** `The system shall support pull mirrors that periodically synchronize from a remote repository.`
- **REPO-07-002:** `The system shall support push mirrors that synchronize to a remote repository on push or on schedule.`
- **REPO-07-003:** `The system shall support LFS object mirroring when LFS is enabled.`

### Event-Driven Requirements (Mirror Workflow)

- **REPO-07-101:** `When a pull mirror sync is triggered (scheduled or manual), the system shall fetch all refs and objects from the remote repository.`
- **REPO-07-102:** `When a push mirror sync is triggered (on push or scheduled), the system shall push all refs and objects to the configured remote.`
- **REPO-07-103:** `When a user manually triggers a mirror sync, the system shall initiate the synchronization immediately.`
- **REPO-07-104:** `When a mirror requires authentication, the system shall use the configured credentials (password, token, or SSH key).`
- **REPO-07-105:** `When a pull mirror sync completes, the system shall update the repository's refs to match the remote.`

### Optional Feature Requirements (Mirror Configuration)

- **REPO-07-201:** `Where LFS mirroring is enabled, the system shall synchronize LFS objects along with Git data during mirror operations.`
- **REPO-07-202:** `Where ALLOW_LOCALNETWORKS is enabled in [migrations], the system shall permit mirror URLs pointing to local network addresses.`
- **REPO-07-203:** `Where GitHub OAuth is configured for mirroring, the system shall use the OAuth token to authenticate against private GitHub repositories.`

### Unwanted Behaviour Requirements (Mirror Errors)

- **REPO-07-301:** `If a mirror sync fails, then the system shall log the error and preserve the last successfully synced state.`
- **REPO-07-302:** `If a mirror URL is invalid or unreachable, then the system shall report the failure and not modify the local repository.`
- **REPO-07-303:** `If ALLOW_LOCALNETWORKS is disabled in [migrations] and a mirror URL points to a local network address, then the system shall reject the mirror configuration.`
- **REPO-07-304:** `If a user pushes to a pull-mirrored repository, the system shall reject the push.`

---

## 8. Migration (REPO-08)

**User Story:** As a user, I want to migrate my repository from another platform so that I can move my project to Gitea with its history and metadata preserved.

### Ubiquitous Requirements (Migration Sources)

- **REPO-08-001:** `The system shall support migration from the following platforms: Gitea, GitHub, GitLab, Gogs, OneDev, Codebase, GitBucket, and plain Git.`
- **REPO-08-002:** `The system shall migrate Git data (commits, branches, tags) as part of every migration.`
- **REPO-08-003:** `The system shall track migration progress and report status to the user.`
- **REPO-08-004:** `The system shall support backup dump (export) and restore (import) migration paths via the dump/restore downloader, allowing full repository export to an archive and subsequent import.`

### Event-Driven Requirements (Migration Workflow)

- **REPO-08-101:** `When a user initiates a migration, the system shall clone the source repository and import the selected data items.`
- **REPO-08-102:** `When a migration includes issues, the system shall import issues with their labels, milestones, comments, and attachments.`
- **REPO-08-103:** `When a migration includes pull requests, the system shall import pull requests with their reviews, comments, and merge status.`
- **REPO-08-104:** `When a migration includes releases, the system shall import releases with their tags, assets, and descriptions.`
- **REPO-08-105:** `When a migration includes wiki, the system shall import the wiki content as a separate Git repository.`
- **REPO-08-106:** `When a migration includes LFS objects, the system shall download and store LFS objects from the source.`
- **REPO-08-107:** `When a migration fails, the system shall roll back all partially imported data and report the error.`

### Optional Feature Requirements (Migration Configuration)

- **REPO-08-201:** `Where ALLOWED_DOMAINS is configured for migration, the system shall reject migration from source URLs not on the allowlist.`
- **REPO-08-202:** `Where OAuth authentication is configured for a migration source, the system shall use the OAuth token to access private source data.`
- **REPO-08-203:** `Where local filesystem migration is enabled, the system shall allow migration from local Git repository paths for privileged users.`

### Unwanted Behaviour Requirements (Migration Errors)

- **REPO-08-301:** `If a migration source URL is not on the configured allowlist, then the system shall reject the migration.`
- **REPO-08-302:** `If a migration exceeds the configured maximum file size, then the system shall abort the migration with an error.`
- **REPO-08-303:** `If a migration times out, then the system shall abort the migration and roll back partially imported data.`
- **REPO-08-304:** `If authentication to the migration source fails, then the system shall reject the migration with a credentials error.`

---

## 9. Topics (REPO-09)

**User Story:** As a repository owner, I want to add topics to my repository so that others can discover it by category.

### Ubiquitous Requirements (Topic Properties)

- **REPO-09-001:** `The system shall allow each repository to have zero or more topics assigned.`
- **REPO-09-002:** `The system shall validate topic names against the pattern ^[a-z0-9][-.a-z0-9]*$ with a maximum length of 35 characters.`
- **REPO-09-003:** `The system shall provide instance-wide topic search across all visible repositories.`

### Event-Driven Requirements (Topic Lifecycle)

- **REPO-09-101:** `When a user adds a topic to a repository, the system shall associate the topic with the repository.`
- **REPO-09-102:** `When a user removes a topic from a repository, the system shall disassociate the topic from the repository.`
- **REPO-09-103:** `When a user searches for a topic, the system shall return all repositories with that topic visible to the user.`

### Unwanted Behaviour Requirements (Topic Errors)

- **REPO-09-301:** `If a user submits a topic name that does not match the validation pattern, then the system shall reject the topic with a validation error.`
- **REPO-09-302:** `If a user submits a topic name exceeding 35 characters, then the system shall reject the topic.`
- **REPO-09-303:** `If a user without write access attempts to modify a repository's topics, then the system shall deny the operation.`

---

## 10. Starring & Watching (REPO-10)

**User Story:** As a user, I want to star repositories to bookmark them and watch repositories to receive notifications about their activity.

### Ubiquitous Requirements (Star Properties)

- **REPO-10-001:** `The system shall allow each user to star any visible repository.`
- **REPO-10-002:** `The system shall track and display the star count for each repository.`
- **REPO-10-003:** `The system shall support four watch modes: None, Normal (all events), Dont (unsubscribe), and Auto (auto-watch on push).`

### Event-Driven Requirements (Star Workflow)

- **REPO-10-101:** `When a user stars a repository, the system shall add the repository to the user's starred list and increment the star count.`
- **REPO-10-102:** `When a user unstars a repository, the system shall remove the repository from the user's starred list and decrement the star count.`
- **REPO-10-103:** `When a user requests their starred repository list, the system shall return all repositories starred by the user.`

### Event-Driven Requirements (Watch Workflow)

- **REPO-10-104:** `When a user watches a repository in Normal mode, the system shall deliver notifications for all repository events to the user.`
- **REPO-10-105:** `When a user sets a repository to Dont mode, the system shall stop delivering notifications for that repository even if the user would otherwise receive them.`
- **REPO-10-106:** `When a user pushes to a repository and AUTO_WATCH_ON_CHANGES is enabled, the system shall automatically set the user's watch mode to Normal for that repository.`

### Optional Feature Requirements (Watch Configuration)

- **REPO-10-201:** `Where AUTO_WATCH_NEW_REPOS is enabled, the system shall automatically set watch mode to Normal when a user creates or forks a repository.`
- **REPO-10-202:** `Where AUTO_WATCH_ON_CHANGES is enabled, the system shall automatically set watch mode to Normal when a user pushes to a repository.`
- **REPO-10-203:** `Where DISABLE_STARS is enabled in [repository], the system shall disable the star feature instance-wide, hiding star UI and rejecting star operations.`

### Unwanted Behaviour Requirements (Star & Watch Errors)

- **REPO-10-301:** `If a user attempts to star a repository they do not have visibility of, then the system shall return a 404 error.`
- **REPO-10-302:** `If a user attempts to star a repository they have already starred, then the system shall ignore the duplicate action.`
- **REPO-10-303:** `If a user attempts to unstar a repository they have not starred, then the system shall ignore the action.`

---

## 11. Anonymous Git Clone for Public Repositories (REPO-11)

**User Story:** As an unauthenticated user, I want to clone public repositories so that I can inspect open-source code without creating an account.

### Ubiquitous Requirements (Anonymous Access Properties)

- **REPO-11-001:** `The system shall permit unauthenticated git-upload-pack requests against any public repository (visibility = public, not private).`
- **REPO-11-002:** `The system shall reject all unauthenticated git-receive-pack requests regardless of repository visibility.`
- **REPO-11-003:** `The system shall treat any anonymous clone as if performed by a user with read-only permission.`

### Event-Driven Requirements (Anonymous Access Workflow)

- **REPO-11-101:** `When an unauthenticated client initiates git-upload-pack on a public repository, the system shall serve the repository without prompting for credentials.`
- **REPO-11-102:** `When an unauthenticated client initiates a Git operation on a private repository, the system shall challenge with HTTP 401 and request credentials.`
- **REPO-11-103:** `When the organization owning a public repository has been marked as not publicly visible, the system shall reject anonymous pulls of that organization's repositories.`

### Optional Feature Requirements (Anonymous Access Configuration)

- **REPO-11-201:** `Where DISABLE_HTTP_GIT is enabled, the system shall reject all anonymous Git HTTP operations regardless of repository visibility.`
- **REPO-11-202:** `Where REQUIRE_SIGNIN_VIEW is enabled, the system shall reject anonymous access to all web and Git HTTP endpoints.`

### Unwanted Behaviour Requirements (Anonymous Access Errors)

- **REPO-11-301:** `If an anonymous client requests git-receive-pack, then the system shall respond with HTTP 401 Unauthorized.`
- **REPO-11-302:** `If an anonymous client attempts to access an archived public repository for write, then the system shall reject the request.`
- **REPO-11-303:** `If an anonymous pull exceeds the configured per-IP rate limit, then the system shall apply backpressure or reject subsequent requests.`

---

## 12. Repository Transfer (REPO-12)

**User Story:** As a repository owner, I want to transfer ownership of a repository to another user or organization so that responsibility for the repository changes hands cleanly.

### Ubiquitous Requirements (Transfer Properties)

- **REPO-12-001:** `The system shall track pending repository transfers as a distinct state with initiator, recipient, and target teams.`
- **REPO-12-002:** `The system shall require explicit acceptance by the recipient before a transfer completes.`
- **REPO-12-003:** `The system shall create a persistent redirect from the old repository path to the new path upon completion.`
- **REPO-12-004:** `The system shall preserve all issues, pull requests, releases, wiki content, and settings across the transfer.`

### Event-Driven Requirements (Transfer Workflow)

- **REPO-12-101:** `When an owner initiates a transfer to a recipient, the system shall create a pending transfer record and notify the recipient.`
- **REPO-12-102:** `When the recipient accepts the transfer, the system shall reassign repository ownership, update paths, and create a redirect.`
- **REPO-12-103:** `When the recipient rejects the transfer, the system shall cancel the pending record and notify the initiator.`
- **REPO-12-104:** `When the initiator cancels a pending transfer, the system shall remove the pending record.`
- **REPO-12-105:** `When the recipient accepts a transfer to an organization and the initiator specified target teams, the system shall grant those teams access to the transferred repository.`
- **REPO-12-106:** `When a transfer completes, the system shall transfer all deploy keys, webhooks, and per-repository configuration to the new owner.`

### State-Driven Requirements (Pending Transfer)

- **REPO-12-701:** `While a transfer is pending, the system shall continue serving the repository from the original path.`
- **REPO-12-702:** `While a transfer is pending, the system shall block the initiator from deleting the repository.`

### Unwanted Behaviour Requirements (Transfer Errors)

- **REPO-12-301:** `If the recipient does not have permission to own repositories of the requested visibility, then the system shall reject the acceptance.`
- **REPO-12-302:** `If a non-owner attempts to initiate or cancel a transfer, then the system shall deny the operation.`
- **REPO-12-303:** `If the recipient namespace already contains a repository with the same name, then the system shall reject the acceptance.`
- **REPO-12-304:** `If the recipient rejects the transfer after the initiator has lost ownership (e.g. account deletion), then the system shall preserve the repository under a ghost owner.`

---

## 13. Branch Rename & Restore (REPO-13)

**User Story:** As a maintainer, I want to rename branches and restore deleted branches so that I can recover from mistakes and reorganize branch structure.

### Ubiquitous Requirements (Branch Rename Properties)

- **REPO-13-001:** `The system shall support renaming any branch in a repository.`
- **REPO-13-002:** `The system shall create a redirect from the old branch name to the new branch name after rename.`
- **REPO-13-003:** `The system shall support restoring recently-deleted branches from their reflog entries.`

### Event-Driven Requirements (Branch Lifecycle)

- **REPO-13-101:** `When a user renames a branch, the system shall update the ref, create a redirect, and update all internal references including default-branch pointer if applicable.`
- **REPO-13-102:** `When a user renames the default branch, the system shall update the repository's HEAD reference atomically.`
- **REPO-13-103:** `When a user requests restoration of a deleted branch, the system shall resurrect the ref from the most recent reflog entry.`
- **REPO-13-104:** `When a user restores a branch via the deleted-branches UI, the system shall recreate the ref and any associated branch-protection rules.`

### Optional Feature Requirements (Rename Constraints)

- **REPO-13-201:** `Where a branch is protected, the system shall require admin permission to rename it.`
- **REPO-13-202:** `Where a renamed branch is referenced by open pull requests, the system shall update the PR base or head references accordingly.`

### Unwanted Behaviour Requirements (Branch Rename Errors)

- **REPO-13-301:** `If a user attempts to rename a branch to a name that already exists, then the system shall reject the rename.`
- **REPO-13-302:** `If a non-admin user attempts to rename a protected branch, then the system shall deny the operation.`
- **REPO-13-303:** `If a user attempts to restore a branch whose reflog entry has expired, then the system shall report that no restorable state exists.`

---

## 14. Repository Activity (REPO-14)

**User Story:** As a maintainer, I want a per-repository activity dashboard so that I can see contribution trends, code frequency, and recent commits at a glance.

### Ubiquitous Requirements (Activity Properties)

- **REPO-14-001:** `The system shall provide an Activity tab on each repository exposing time-bucketed contribution views.`
- **REPO-14-002:** `The system shall render a contributors view showing per-author commit, addition, and deletion counts.`
- **REPO-14-003:** `The system shall render a code-frequency chart showing cumulative additions versus deletions over time.`
- **REPO-14-004:** `The system shall render a recent-commits view showing commit frequency per time bucket.`

### Event-Driven Requirements (Activity Workflow)

- **REPO-14-101:** `When a user navigates to /{owner}/{repo}/activity, the system shall render the activity dashboard for the default time window.`
- **REPO-14-102:** `When a user selects a time period (daily, halfweekly, weekly, monthly, quarterly, semiyearly, or yearly), the system shall re-bucket and re-render the activity charts.`
- **REPO-14-103:** `When a client requests /{owner}/{repo}/activity/code-frequency/data, the system shall return the code-frequency series as JSON.`
- **REPO-14-104:** `When a client requests /{owner}/{repo}/activity/recent-commits/data, the system shall return the recent-commits series as JSON.`
- **REPO-14-105:** `When a user navigates to /{owner}/{repo}/activity/authors, the system shall render the contributor breakdown with author filtering.`

### Optional Feature Requirements (Activity Configuration)

- **REPO-14-201:** `Where the repository is empty, the system shall display an empty-state message in place of activity charts.`
- **REPO-14-202:** `Where the user lacks read access to a private repository, the system shall return 404 for activity routes.`

### Unwanted Behaviour Requirements (Activity Errors)

- **REPO-14-301:** `If the activity data exceeds the maximum rendering window, then the system shall truncate to the most recent N buckets.`
- **REPO-14-302:** `If the activity data endpoint is requested with an invalid time bucket size, then the system shall return 400 Bad Request.`

---

## 15. Repository Templates (REPO-15)

**User Story:** As a repository owner, I want to mark a repository as a template and generate new repositories from it with placeholder substitution so that I can bootstrap standardized project layouts.

### Ubiquitous Requirements (Template Properties)

- **REPO-15-001:** `The system shall support an IsTemplate flag on each repository designating it as eligible for generate-from-template operations.`
- **REPO-15-002:** `The system shall expose template repositories as available templates when a user initiates generate-from-template creation.`
- **REPO-15-003:** `The system shall support selective generation of the following template components: Git content, topics, Git hooks, webhooks, avatar, issue labels, and protected branch rules (GenerateRepoOptions).`
- **REPO-15-004:** `The system shall require at least one component to be selected for generation (GenerateRepoOptions.IsValid).`
- **REPO-15-005:** `The system shall preserve the template repository's object format name (SHA-1 or SHA-256) and commit signature trust model in the generated repository.`
- **REPO-15-006:** `The system shall copy LFS objects from the template repository to the generated repository when Git content generation is requested.`

### Event-Driven Requirements (Template Generation Workflow)

- **REPO-15-101:** `When a user generates a repository from a template, the system shall clone the template's default branch at depth 1 into a temporary working directory and commit it to the new repository.`
- **REPO-15-102:** `When the template repository contains a .gitea/template configuration file, the system shall parse it for glob patterns identifying files subject to placeholder substitution.`
- **REPO-15-103:** `When a file path matches a configured .gitea/template glob, the system shall apply placeholder substitution to both the file content and the file path.`
- **REPO-15-104:** `When no default branch is specified for the generated repository, the system shall inherit the template repository's default branch name.`
- **REPO-15-105:** `When a generated repository is committed, the system shall create the initial commit attributed to the generating user as both author and committer.`
- **REPO-15-106:** `When generation completes, the system shall compute and store the repository size for the generated repository (UpdateRepoSize).`

### Optional Feature Requirements (Placeholder Substitution)

- **REPO-15-201:** `Where placeholder substitution is applied, the system shall expand the following variables in matching file contents and paths: REPO_NAME, TEMPLATE_NAME, REPO_DESCRIPTION, TEMPLATE_DESCRIPTION, REPO_OWNER, TEMPLATE_OWNER, REPO_LINK, TEMPLATE_LINK, REPO_HTTPS_URL, TEMPLATE_HTTPS_URL, REPO_SSH_URL, and TEMPLATE_SSH_URL.`
- **REPO-15-202:** `Where a placeholder variable is subject to case transformation, the system shall support the following transformers: SNAKE, KEBAB, CAMEL, PASCAL, LOWER, UPPER, and TITLE (e.g., ${REPO_NAME_SNAKE}).`
- **REPO-15-203:** `Where a placeholder variable is used in a file path, the system shall sanitize the expanded value to a valid OS filename, replacing reserved characters and sequences (including ".." to prevent directory traversal) with underscores.`

### Unwanted Behaviour Requirements (Template Errors)

- **REPO-15-301:** `If a user attempts to generate a repository from a non-template repository, then the system shall reject the operation.`
- **REPO-15-302:** `If a generated repository path already exists on disk, then the system shall reject the generation with ErrRepoFilesAlreadyExist.`
- **REPO-15-303:** `If a .gitea/template glob expression is malformed, then the system shall skip the offending line and continue processing remaining globs.`

---

## 16. Wiki (REPO-16)

**User Story:** As a collaborator, I want to create and maintain wiki pages so that project documentation lives alongside the codebase as a versioned git repository.

> **Boundary note:** The wiki IS a git repository (`models/repo/wiki.go`, cloneable via the `.wiki` suffix). It is documented here as part of the code artifact lifecycle. Wiki page editing is a collaboration activity but the wiki itself is a code artifact.

### Ubiquitous Requirements (Wiki Properties)

- **REPO-16-001:** `The system shall store wiki content as a separate Git repository alongside the main repository.`
- **REPO-16-002:** `The system shall render wiki pages as markdown.`
- **REPO-16-003:** `The system shall designate a default wiki page named "Home".`

### Event-Driven Requirements (Wiki Lifecycle)

- **REPO-16-101:** `When a user creates a new wiki page, the system shall commit the page content to the wiki Git repository.`
- **REPO-16-102:** `When a user edits a wiki page, the system shall record the change as a new commit in the wiki repository preserving full edit history.`
- **REPO-16-103:** `When a user deletes a wiki page, the system shall remove the page from the wiki repository.`
- **REPO-16-104:** `When a user creates a page named "_Sidebar", the system shall display it as navigation in the wiki view.`
- **REPO-16-105:** `When a user creates a page named "_Footer", the system shall display it as footer content in the wiki view.`

### Optional Feature Requirements (Wiki Configuration)

- **REPO-16-201:** `Where an external wiki URL is configured for a repository, the system shall redirect the wiki tab to the external URL.`
- **REPO-16-202:** `Where the Wiki unit is disabled, the system shall hide the wiki tab from the repository navigation.`

### Unwanted Behaviour Requirements (Wiki Errors)

- **REPO-16-301:** `If a user without write permission attempts to create or edit a wiki page, then the system shall deny the operation.`
- **REPO-16-302:** `If a wiki page name conflicts with the reserved names "_Sidebar" or "_Footer" used for navigation, then the system shall treat it as a navigation element rather than a regular page.`

---

## 17. Search

This section documents repo search, code search, and the repository-explore portions of the global explore feature as capabilities of the Repository & Code domain. Indexer backend infrastructure is documented in Domain 08 (Administration & Operations).

### 17a. Repository Search (REPO-17)

**User Story:** As a user, I want to search for repositories by keyword and filters so that I can find relevant projects quickly.

#### Ubiquitous Requirements (Search Defaults)

- **REPO-17-001:** `The system shall match repository keywords using database LIKE (substring) queries; there is no fuzzy/exact mode toggle for repository search.`
- **REPO-17-002:** `The system shall apply repository search filters including topic and primary language. Star count, fork count, and repository size are available as sort orders only (not filters); no license filter exists.`
- **REPO-17-003:** `The system shall apply ownership filters including user/org UID, archived status, private status, and template status.`
- **REPO-17-004:** `The system shall apply repository mode filters: source, fork, mirror, and collaborative.`
- **REPO-17-005:** `The system shall support the following sort options: alphabetical, reverse alphabetical, created date, oldest, updated, recent update, least update, size, reverse size, ID, reverse ID, stars, and forks.`
- **REPO-17-006:** `The system shall paginate repository search results using the configured EXPLORE_PAGING_NUM page size.`
- **REPO-17-007:** `The system shall expose repository search via the API endpoint GET /repos/search with support for all filters and sort options.`

#### Event-Driven Requirements (Search Execution)

- **REPO-17-101:** `When a user submits a keyword query on the explore repos page, the system shall return matching repositories filtered by the user's visibility scope.`
- **REPO-17-102:** `When a user selects a sort option, the system shall reorder the results according to the selected sort criteria.`
- **REPO-17-103:** `When the setting ONLY_SHOW_RELEVANT_REPOS is enabled, the system shall filter out repositories with low relevance scores from explore results.`
- **REPO-17-104:** `When the setting SEARCH_REPO_DESCRIPTION is enabled, the system shall include repository descriptions in the search index for matching.`

#### Optional Feature Requirements (Repo Indexer)

- **REPO-17-201:** `Where REPO_INDEXER_ENABLED is true, the system shall index repository names and descriptions for full-text search.`
- **REPO-17-202:** `Where a default explore sort is configured via EXPLORE_PAGING_DEFAULT_SORT, the system shall use that sort order as the initial ordering on the explore repos page.`

#### Unwanted Behaviour Requirements (Search Errors)

- **REPO-17-301:** `If a user submits a search query without permission to view any matching repositories, then the system shall return an empty result set.`
- **REPO-17-302:** `If the repository indexer is disabled and a full-text search is attempted, then the system shall fall back to database-based name matching.`

### 17b. Code Search (REPO-18)

**User Story:** As a developer, I want to search through source code across repositories so that I can find specific implementations, patterns, or usages.

#### Ubiquitous Requirements (Code Search Properties)

- **REPO-18-001:** `The system shall support full-text search through source code in repositories accessible to the authenticated user.`
- **REPO-18-002:** `The system shall support two indexer backends for code search: Bleve (embedded, default) and Elasticsearch (external). (Indexer backend configuration → Domain 08.)`
- **REPO-18-003:** `The system shall provide a language filter allowing users to restrict code search results to specific programming languages.`
- **REPO-18-004:** `The system shall display syntax highlighting in code search results based on file language.`
- **REPO-18-005:** `The system shall display context lines surrounding each search match.`
- **REPO-18-006:** `The system shall expose code search via GET /explore/code (global) and GET /{owner}/{repo}/search (repo-scoped).`
- **REPO-18-007:** `The system shall support two search modes for code search: fuzzy (Levenshtein distance, default) and exact (literal match), selected via the fuzzy query parameter.`

#### Event-Driven Requirements (Indexing Workflow)

- **REPO-18-101:** `When a repository is updated with a new commit, the system shall enqueue the affected files for re-indexing in the code indexer.`
- **REPO-18-102:** `When a user submits a code search query, the system shall query the configured indexer backend and return matching file content with line numbers.`
- **REPO-18-103:** `When the indexer encounters a file exceeding MAX_FILE_SIZE (default 1MB), the system shall skip indexing that file.`

#### Optional Feature Requirements (Indexer Configuration)

- **REPO-18-201:** `Where REPO_INDEXER_INCLUDE patterns are configured, the system shall only index files matching those glob patterns.`
- **REPO-18-202:** `Where REPO_INDEXER_EXCLUDE patterns are configured, the system shall skip files matching those glob patterns during indexing.`
- **REPO-18-203:** `Where REPO_INDEXER_EXCLUDE_VENDORED is enabled, the system shall skip files in vendor directories during indexing.`
- **REPO-18-204:** `Where REPO_INDEXER_TYPE is set to elasticsearch, the system shall connect to the external Elasticsearch cluster specified by REPO_INDEXER_CONN_STR.`
- **REPO-18-205:** `The system shall index only repositories whose type matches the configured REPO_INDEXER_REPO_TYPES list (default: sources, forks, mirrors, templates).`

#### Unwanted Behaviour Requirements (Code Search Errors)

- **REPO-18-301:** `If the code indexer is not enabled, then the system shall not display the code search tab and shall return an error for code search API requests.`
- **REPO-18-302:** `If the external indexer backend is unreachable, then the system shall return an error for code search queries and log the connection failure.`
- **REPO-18-303:** `If the indexer startup exceeds STARTUP_TIMEOUT, then the system shall log a timeout error and proceed without code search availability.`

### 17c. Explore Repos (REPO-19)

**User Story:** As a visitor, I want to browse curated views of public repositories and code across the instance so that I can discover interesting content.

> **Boundary note:** Explore users/organizations → Domains 01/02. This subsection covers only the repository and code explore surfaces.

#### Ubiquitous Requirements (Explore Properties)

- **REPO-19-001:** `The system shall provide explore tabs for repositories and code. (User and organization explore tabs → Domains 01/02.)`
- **REPO-19-002:** `The system shall paginate explore results using the configured EXPLORE_PAGING_NUM page size.`
- **REPO-19-003:** `The system shall expose topic search via GET /explore/topics/search.`

#### Event-Driven Requirements (Explore Navigation)

- **REPO-19-101:** `When a visitor navigates to GET /explore, the system shall redirect to GET /explore/repos.`
- **REPO-19-102:** `When a visitor navigates to GET /explore/repos, the system shall display a paginated list of public repositories sorted by the configured default sort.`
- **REPO-19-103:** `When a visitor navigates to GET /explore/code, the system shall display the global code search interface if the code indexer is enabled.`

#### Optional Feature Requirements (Explore Configuration)

- **REPO-19-201:** `Where ONLY_SHOW_RELEVANT_REPOS is enabled, the system shall filter the explore repos view to show only repositories deemed relevant based on activity and popularity.`
- **REPO-19-202:** `Where the code indexer is disabled, the system shall hide the code explore tab from the navigation.`

#### Unwanted Behaviour Requirements (Explore Errors)

- **REPO-19-301:** `If an unauthenticated visitor accesses explore pages, then the system shall show only public resources.`

---

## Configuration Reference

The following INI sections configure repository & code behaviors.

### [repository] Section

- **ROOT**: Absolute or relative path under which all repository Git data is stored (default `<AppDataPath>/gitea-repositories`).
- **DEFAULT_PRIVATE**: Default visibility for new repositories (`private`, `public`, `last`; default `last`).
- **FORCE_PRIVATE**: Force all new repositories to be private (default `false`).
- **MAX_CREATION_LIMIT**: Global per-user repository creation limit (default `-1` = unlimited).
- **DEFAULT_BRANCH**: Default branch name for new repositories (default `main`).
- **DISABLE_HTTP_GIT**: Disable Git operations over HTTP (default `false`).
- **USE_COMPAT_SSH_URI**: Generate SSH clone URLs in compatible `ssh://` format (default `false`).
- **DISABLE_STARS**: Disable the star feature instance-wide (default `false`).
- **DISABLE_MIGRATIONS**: Disable repository migration (default `false`).
- **ALLOW_ADOPTION_OF_UNADOPTED_REPOSITORIES**: Allow admins to adopt unadopted repositories (default `false`).
- **ALLOW_DELETE_OF_UNADOPTED_REPOSITORIES**: Allow admins to delete unadopted repositories (default `false`).

### [repository.upload] Section

- **ENABLED**: Enable file uploads in the web editor.
- **TEMP_PATH**: Temporary upload staging path.
- **ALLOWED_TYPES**: Comma-separated MIME type allowlist for uploads (default empty = allow all).
- **FILE_MAX_SIZE**: Maximum upload file size in MB (default 3).
- **MAX_FILES**: Maximum files per single upload batch (default 5).

### [repository.local] Section

- **LOCAL_COPY_PATH**: Path used for staging local working copies during web-editor operations.

### [repository.pull-request] Section

- **WORK_IN_PROGRESS_PREFIXES**: Comma-separated title prefixes treated as WIP markers (default `WIP:,[WIP]`).
- **CLOSE_KEYWORDS**: Comma-separated keywords that auto-close issues when used in PR body or merge commit.
- **REOPEN_KEYWORDS**: Comma-separated keywords that auto-reopen issues.

### [git] Section

- **PATH**: Path to the git binary (default `git`).
- **HOME_PATH**: Home directory used for git operations (default `home`, resolved under AppDataPath).
- **DISABLE_DIFF_HIGHLIGHT**: Disable syntax-highlighted diff rendering in the web UI (default `false`).
- **MAX_GIT_DIFF_LINES**: Maximum lines per diff file before truncation (default 1000).
- **MAX_GIT_DIFF_LINE_CHARACTERS**: Maximum characters per diff line (default 5000).
- **MAX_GIT_DIFF_FILES**: Maximum files shown in a diff before truncation (default 100).
- **COMMITS_RANGE_SIZE**: Number of commits to fetch per batch/page (default 50).
- **BRANCHES_RANGE_SIZE**: Number of branches to fetch per batch/page (default 20).
- **VERBOSE_PUSH**: Render verbose push output when enabled (default `true`).
- **VERBOSE_PUSH_DELAY**: Delay before verbose push output begins (default `5s`).
- **GC_ARGS**: Arguments passed to `git gc` when invoked by Gitea.
- **ENABLE_AUTO_GIT_WIRE_PROTOCOL**: Enable Git Wire Protocol (protocol v2) by default (default `true`).
- **PULL_REQUEST_PUSH_MESSAGE**: Enable push messages for pull request operations (default `true`).
- **LARGE_OBJECT_THRESHOLD**: Threshold above which objects are treated as large (default `1048576`).
- **DISABLE_PARTIAL_CLONE**: Disable partial clone support (default `false`).
- **DISABLE_CORE_PROTECT_NTFS**: Disable `core.protectNTFS` (default `false`).

### [git.timeout] Section

Nested under `[git]`, configures per-operation timeouts (all in seconds):

- **DEFAULT**: Default git operation timeout (default 360).
- **MIGRATE**: Migration timeout (default 600).
- **MIRROR**: Mirror sync timeout (default 300).
- **CLONE**: Clone timeout (default 300).
- **PULL**: Pull timeout (default 300).
- **GC**: Garbage collection timeout (default 60).

### [git.config] Section

Arbitrary `key=value` pairs written to the git config of each repository. Examples:
- **diff.algorithm**
- **core.logAllRefUpdates**
- **receive.advertisePushOptions**

### [git.reflog] Section

Deprecated since v1.21. Keys are remapped to `[git.config]` entries:

- **ENABLED**: *(deprecated, since v1.21)* → remapped to `[git.config]` `core.logAllRefUpdates`.
- **EXPIRATION**: *(deprecated, since v1.21)* → remapped to `[git.config]` `gc.reflogExpire`.

Use `[git.config]` directly to set `core.logAllRefUpdates` (default `true`) and `gc.reflogExpire` (default `90`).

### [mirror] Section

- **ENABLED**: Enable the mirror feature (default `true`). When disabled, both pull and push mirror creation are blocked.
- **DISABLE_NEW_PULL**: Disable creation of new pull mirrors (default `false`).
- **DISABLE_NEW_PUSH**: Disable creation of new push mirrors (default `false`).
- **DEFAULT_INTERVAL**: Default sync interval for new pull mirrors (default `8h`).
- **MIN_INTERVAL**: Minimum allowed sync interval for pull mirrors (default `10m`).

### [migrations] Section

- **MAX_ATTEMPTS**: Maximum retry attempts for migration operations (default 3).
- **RETRY_BACKOFF**: Backoff in seconds between retry attempts (default 3).
- **ALLOWED_DOMAINS**: Comma-separated allowlist of domains permitted as migration/mirror sources (default empty = allow all).
- **BLOCKED_DOMAINS**: Comma-separated blocklist of domains forbidden as migration/mirror sources (default empty).
- **ALLOW_LOCALNETWORKS**: Permit migration/mirror URLs pointing to local network addresses (default `false`).
- **SKIP_TLS_VERIFY**: Skip TLS certificate verification when fetching from migration sources (default `false`).

### [ui] Section

- **MAX_DISPLAY_FILE_SIZE**: Maximum file size (in bytes) for inline rendering in the file viewer (default `8388608` = 8 MB). This is the sole display-size limit; there is no separate `RENDERED_MAX_FILE_SIZE`.

### [indexer] Section (code-search keys)

> **Cross-reference:** Indexer backend infrastructure (Bleve/Elasticsearch/Meilisearch/db), queue processing, and startup behavior → Domain 08 (Administration & Operations). Only the code-search-scoped keys are listed here.

- **REPO_INDEXER_ENABLED**: Enable the code (repository) indexer (default `false`).
- **REPO_INDEXER_TYPE**: Code indexer backend type (`bleve` or `elasticsearch`; default `bleve`).
- **REPO_INDEXER_PATH**: Local filesystem path for the Bleve code index.
- **REPO_INDEXER_CONN_STR**: Elasticsearch connection string (when `REPO_INDEXER_TYPE = elasticsearch`).
- **REPO_INDEXER_INCLUDE**: Comma-separated glob patterns; only matching files are indexed.
- **REPO_INDEXER_EXCLUDE**: Comma-separated glob patterns; matching files are skipped.
- **REPO_INDEXER_EXCLUDE_VENDORED**: Skip files in vendor directories (default `true`).
- **REPO_INDEXER_REPO_TYPES**: Repository types to index (default `sources,forks,mirrors,templates`).
- **MAX_FILE_SIZE**: Maximum file size (in bytes) for code-indexer ingestion (default `1048576`).
- **STARTUP_TIMEOUT**: Per-indexer startup timeout (default `30s`).

---

## Business Rules

- **BR-03-001:** Repository names must be unique within an owner's namespace (case-insensitive)
- **BR-03-002:** Repository names are restricted to alphanumeric characters, dots, and dashes
- **BR-03-003:** Reserved repository names (".", "..", "-") and patterns ("*.git", "*.wiki", "*.rss", "*.atom") are prohibited
- **BR-03-004:** Archived repositories are fully read-only: no pushes, no new issues, no new pull requests
- **BR-03-005:** Fork visibility must be at least as restrictive as the base repository
- **BR-03-006:** Mirror repositories are read-only for pull mirrors; push mirrors are outbound only
- **BR-03-007:** Protected branch rules apply to all users except those explicitly on the allowlist (push access restrictions → Domain 02)
- **BR-03-008:** Repository renames create persistent redirects from the old path to the new path
- **BR-03-009:** Starring and watching are per-user, per-repository actions with no cross-user effects
- **BR-03-010:** Topic names must match ^[a-z0-9][-.a-z0-9]*$ and be at most 35 characters
- **BR-03-011:** Migration source domains must be on the configured allowlist when ALLOWED_DOMAINS is set
- **BR-03-012:** Wiki content is stored as a separate Git repository and is fully cloneable with standard Git tools
- **BR-03-013:** Repository search applies visibility scoping based on the authenticated user's permissions
- **BR-03-014:** Code search requires an enabled and initialized code indexer
- **BR-03-015:** Code indexer applies include/exclude patterns before the MAX_FILE_SIZE check
- **BR-03-016:** Fuzzy code search uses Levenshtein distance as the default matching algorithm
- **BR-03-017:** Template placeholder substitution is driven by `.gitea/template` globs; only files matching configured globs are transformed
- **BR-03-018:** Template filename substitution sanitizes expanded values to prevent directory traversal and reserved filename conflicts

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Push to archived repository | Reject with read-only error |
| Push to pull-mirror repository | Reject the push (mirror repos are synchronized from remote only) |
| Fork to existing owner/name | Reject fork with name conflict error |
| Repository rename collision | Reject rename, keep current name |
| Transfer to owner with same-named repo | Reject transfer with name conflict error |
| Delete default branch | Allow deletion; user must set a new default branch |
| Delete protected branch | Reject deletion while protection rule is active |
| Push to protected branch without allowlist | Reject push with protection rule error |
| Merge PR without required approvals | Reject merge until approval count is met (see Domain 04) |
| Merge PR with rejected review (block enabled) | Reject merge while a rejected review exists (see Domain 04) |
| Migrate repository exceeding max file size | Abort migration and roll back |
| Migration source unreachable | Abort migration, report error, roll back |
| Star already-starred repository | Ignore duplicate (idempotent) |
| File exceeds MAX_DISPLAY_FILE_SIZE | Display notice instead of content |
| Diff exceeds configured line/file limits | Truncate diff display |
| Cherry-pick with merge conflicts | Report conflicts, do not complete operation |
| Topic name validation failure | Reject topic with validation error |
| 2FA user HTTP Git with password | Reject authentication, require token |
| Wiki page name collision with _Sidebar/_Footer | Treat as navigation element, not regular page |
| External wiki configured on repository | Redirect wiki tab to the external URL |
| Generate from template with pre-existing repo path | Reject with ErrRepoFilesAlreadyExist |
| Malformed .gitea/template glob | Skip the offending line, continue processing |
| Code search with indexer disabled | Hide code search tab, return error for API requests |
| External indexer unreachable during search | Return error to user, log connection failure |
| File exceeds MAX_FILE_SIZE during indexing | Skip file, do not index |
| Empty keyword submitted for repo search | Return all results with applied filters and default sort |
| Code search with include/exclude conflict | Apply exclude after include; exclude takes precedence |

## Success Criteria

- Repository creation completes within 5 seconds for standard configuration
- Git clone over HTTP begins streaming within 3 seconds
- Git clone over SSH begins streaming within 3 seconds
- Archive download generation completes within 10 seconds for repositories under 1 GB
- File browser renders any path within 2 seconds
- Blame view renders within 5 seconds for files under 10,000 lines
- Protected branch enforcement adds no measurable latency to push operations
- Fork operation completes within 10 seconds
- Mirror sync startup occurs within the configured MIN_INTERVAL
- Migration progress is reported within 5 seconds of status change
- Topic search returns results within 1 second for instances with up to 10,000 repositories
- Star/watch toggle operations complete within 500ms
- Web editor commit completes within 3 seconds
- Wiki page rendering completes within 1 second per page
- Template generation (including placeholder substitution) completes within 10 seconds for repositories under 100 MB
- Repository search results return within 2 seconds for standard keyword queries
- Code search results include syntax-highlighted context lines with accurate line numbers
- Explore pages load within 2 seconds with default pagination
- Fuzzy code search matches approximate typos within Levenshtein distance threshold
- Search result pagination provides consistent ordering across page boundaries
