# 04 — Code Management

Baseline specification of Gitea's repository, Git operations, branching, code review, forking, mirroring, migration, topics, collaboration, and starring/watching subsystems. All requirements describe the current (v1.22.x) system behavior.

---

## 1. Repository CRUD

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

## 2. Git Operations

**User Story:** As a developer, I want to clone, push, and pull repositories using standard Git protocols so that I can work with code locally and synchronize changes.

### Ubiquitous Requirements (Protocol Support)

- **REPO-02-001:** `The system shall support Git operations over HTTP and SSH protocols (see INT-05 for SSH server details and INT-12 for Smart HTTP protocol details).`
- **REPO-02-002:** `The system shall support HTTP basic authentication and token authentication for Git operations over HTTP.`
- **REPO-02-003:** `The system shall provide archive downloads in zip, tar.gz, and bundle formats for any branch, tag, or commit.`
- **REPO-02-004:** `The system shall provide raw file downloads for individual files in the repository.`
- **REPO-02-005:** `The system shall generate clone URLs for both HTTP and SSH protocols on the repository page.`
- **REPO-02-006:** `The system shall expose each repository's wiki as a standalone Git repository cloneable via the .wiki suffix on the repository path (over HTTP and SSH).`

### Event-Driven Requirements (Git Protocol Workflow)

- **REPO-02-101:** `When a client initiates a git clone over HTTP, the system shall serve the repository via the Smart HTTP protocol (see INT-12 for full protocol behavior).`
- **REPO-02-102:** `When a client initiates a git push over HTTP, the system shall authenticate the user and accept the push if authorized.`
- **REPO-02-103:** `When a client initiates a Git operation over SSH, the system shall authenticate the user via their registered SSH key (see INT-05 for SSH server behavior).`
- **REPO-02-104:** `When a user requests an archive download, the system shall generate and serve the archive for the specified ref.`
- **REPO-02-105:** `When a user with 2FA enabled attempts HTTP Git operations, the system shall require a personal access token instead of a password.`

### Optional Feature Requirements (Protocol Configuration)

- **REPO-02-201:** `Where DISABLE_HTTP_GIT is enabled, the system shall reject all Git operations over HTTP while keeping the web UI functional.`
- **REPO-02-202:** `Where USE_COMPAT_SSH_URI is enabled, the system shall generate SSH clone URLs in the compatible URI format (ssh://user@host/path).`

### Unwanted Behaviour Requirements (Git Operation Errors)

- **REPO-02-301:** `If a user pushes to an archived repository, then the system shall reject the push.`
- **REPO-02-302:** `If a user pushes to a mirror repository, then the system shall reject the push (mirror repos are read-only).`
- **REPO-02-303:** `If a user with 2FA enabled provides a password for HTTP Git operations, then the system shall reject the authentication.`
- **REPO-02-304:** `If an SSH key is not registered in the system, then the system shall reject the SSH connection (see INT-05 for SSH authentication details).`

---

## 3. Branch Management

**User Story:** As a maintainer, I want to manage branches and configure protection rules so that I can enforce quality standards on critical branches.

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

## 4. Code Review / Browsing

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

## 5. Web Editor

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

## 6. Forking

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

## 7. Mirroring

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

## 8. Migration

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

## 9. Topics

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

## 10. Collaborators

**User Story:** As a repository owner, I want to add collaborators with specific permission levels so that I can control who can read, write, or administer my repository.

### Ubiquitous Requirements (Collaborator Properties)

- **REPO-10-001:** `The system shall support three collaborator permission levels per repository: Read, Write, and Admin.`
- **REPO-10-002:** `The system shall list collaborators separately from organization team members.`
- **REPO-10-003:** `The system shall compute effective permissions by combining collaborator, team, and organization memberships.`

### Event-Driven Requirements (Collaborator Workflow)

- **REPO-10-101:** `When a repository owner adds a collaborator, the system shall grant the specified permission level to that user for the repository.`
- **REPO-10-102:** `When a repository owner removes a collaborator, the system shall revoke all direct repository access for that user.`
- **REPO-10-103:** `When a repository owner changes a collaborator's permission level, the system shall update the permission immediately.`
- **REPO-10-104:** `When a user is both a collaborator and a team member, the system shall grant the highest permission from all sources.`

### Unwanted Behaviour Requirements (Collaborator Errors)

- **REPO-10-301:** `If a non-owner user attempts to add a collaborator, then the system shall deny the operation.`
- **REPO-10-302:** `If a user attempts to add themselves as a collaborator, then the system shall reject the addition.`
- **REPO-10-303:** `If a collaborator is added who is blocked by the repository owner, then the system shall reject the addition.`

---

## 11. Starring & Watching

**User Story:** As a user, I want to star repositories to bookmark them and watch repositories to receive notifications about their activity.

### Ubiquitous Requirements (Star Properties)

- **REPO-11-001:** `The system shall allow each user to star any visible repository.`
- **REPO-11-002:** `The system shall track and display the star count for each repository.`
- **REPO-11-003:** `The system shall support four watch modes: None, Normal (all events), Dont (unsubscribe), and Auto (auto-watch on push).`

### Event-Driven Requirements (Star Workflow)

- **REPO-11-101:** `When a user stars a repository, the system shall add the repository to the user's starred list and increment the star count.`
- **REPO-11-102:** `When a user unstars a repository, the system shall remove the repository from the user's starred list and decrement the star count.`
- **REPO-11-103:** `When a user requests their starred repository list, the system shall return all repositories starred by the user.`

### Event-Driven Requirements (Watch Workflow)

- **REPO-11-104:** `When a user watches a repository in Normal mode, the system shall deliver notifications for all repository events to the user.`
- **REPO-11-105:** `When a user sets a repository to Dont mode, the system shall stop delivering notifications for that repository even if the user would otherwise receive them.`
- **REPO-11-106:** `When a user pushes to a repository and AUTO_WATCH_ON_CHANGES is enabled, the system shall automatically set the user's watch mode to Normal for that repository.`

### Optional Feature Requirements (Watch Configuration)

- **REPO-11-201:** `Where AUTO_WATCH_NEW_REPOS is enabled, the system shall automatically set watch mode to Normal when a user creates or forks a repository.`
- **REPO-11-202:** `Where AUTO_WATCH_ON_CHANGES is enabled, the system shall automatically set watch mode to Normal when a user pushes to a repository.`
- **REPO-11-203:** `Where DISABLE_STARS is enabled in [repository], the system shall disable the star feature instance-wide, hiding star UI and rejecting star operations.`

### Unwanted Behaviour Requirements (Star & Watch Errors)

- **REPO-11-301:** `If a user attempts to star a repository they do not have visibility of, then the system shall return a 404 error.`
- **REPO-11-302:** `If a user attempts to star a repository they have already starred, then the system shall ignore the duplicate action.`
- **REPO-11-303:** `If a user attempts to unstar a repository they have not starred, then the system shall ignore the action.`

---

## 12. Deploy Keys (REPO-12)

**User Story:** As a repository administrator, I want to register per-repository SSH public keys so that CI/CD systems can clone or push without sharing a user account's keys.

### Ubiquitous Requirements (Deploy Key Properties)

- **REPO-12-001:** `The system shall scope each deploy key to exactly one repository.`
- **REPO-12-002:** `The system shall allow each deploy key to be marked read-only or read-write.`
- **REPO-12-003:** `The system shall store the SSH public key, fingerprint, title, creation timestamp, and associated repository for each deploy key.`
- **REPO-12-004:** `The system shall allow the same SSH public key to be registered as a deploy key on multiple repositories without conflict.`

### Event-Driven Requirements (Deploy Key Workflow)

- **REPO-12-101:** `When an admin submits an SSH public key as a deploy key for a repository, the system shall store the key and compute its fingerprint.`
- **REPO-12-102:** `When an SSH connection authenticates with a deploy key, the system shall grant access scoped to the associated repository only.`
- **REPO-12-103:** `When a deploy key marked read-only attempts a git-receive-pack, the system shall reject the push.`
- **REPO-12-104:** `When an admin deletes a deploy key, the system shall remove the key and reject future SSH authentications using that key.`

### Optional Feature Requirements (Deploy Key Sources)

- **REPO-12-201:** `Where a deploy key title is provided at creation, the system shall display the title in the deploy-key management list.`
- **REPO-12-202:** `Where an admin creates a deploy key by generating a new SSH keypair, the system shall display the private key once for download and store only the public key.`

### Unwanted Behaviour Requirements (Deploy Key Errors)

- **REPO-12-301:** `If a deploy key attempts to access a repository other than its scoped repository, then the system shall reject the operation.`
- **REPO-12-302:** `If a non-admin user attempts to add or delete a deploy key, then the system shall deny the operation.`
- **REPO-12-303:** `If a submitted SSH public key is malformed, then the system shall reject the deploy-key creation with a parse error.`

---

## 13. Anonymous Git Clone for Public Repositories (REPO-13)

**User Story:** As an unauthenticated user, I want to clone public repositories so that I can inspect open-source code without creating an account.

### Ubiquitous Requirements (Anonymous Access Properties)

- **REPO-13-001:** `The system shall permit unauthenticated git-upload-pack requests against any public repository (visibility = public, not private).`
- **REPO-13-002:** `The system shall reject all unauthenticated git-receive-pack requests regardless of repository visibility.`
- **REPO-13-003:** `The system shall treat any anonymous clone as if performed by a user with read-only permission.`

### Event-Driven Requirements (Anonymous Access Workflow)

- **REPO-13-101:** `When an unauthenticated client initiates git-upload-pack on a public repository, the system shall serve the repository without prompting for credentials.`
- **REPO-13-102:** `When an unauthenticated client initiates a Git operation on a private repository, the system shall challenge with HTTP 401 and request credentials.`
- **REPO-13-103:** `When the organization owning a public repository has been marked as not publicly visible, the system shall reject anonymous pulls of that organization's repositories.`

### Optional Feature Requirements (Anonymous Access Configuration)

- **REPO-13-201:** `Where DISABLE_HTTP_GIT is enabled, the system shall reject all anonymous Git HTTP operations regardless of repository visibility.`
- **REPO-13-202:** `Where REQUIRE_SIGNIN_VIEW is enabled, the system shall reject anonymous access to all web and Git HTTP endpoints.`

### Unwanted Behaviour Requirements (Anonymous Access Errors)

- **REPO-13-301:** `If an anonymous client requests git-receive-pack, then the system shall respond with HTTP 401 Unauthorized.`
- **REPO-13-302:** `If an anonymous client attempts to access an archived public repository for write, then the system shall reject the request.`
- **REPO-13-303:** `If an anonymous pull exceeds the configured per-IP rate limit, then the system shall apply backpressure or reject subsequent requests.`

---

## 14. Repository Transfer (REPO-14)

**User Story:** As a repository owner, I want to transfer ownership of a repository to another user or organization so that responsibility for the repository changes hands cleanly.

### Ubiquitous Requirements (Transfer Properties)

- **REPO-14-001:** `The system shall track pending repository transfers as a distinct state with initiator, recipient, and target teams.`
- **REPO-14-002:** `The system shall require explicit acceptance by the recipient before a transfer completes.`
- **REPO-14-003:** `The system shall create a persistent redirect from the old repository path to the new path upon completion.`
- **REPO-14-004:** `The system shall preserve all issues, pull requests, releases, wiki content, and settings across the transfer.`

### Event-Driven Requirements (Transfer Workflow)

- **REPO-14-101:** `When an owner initiates a transfer to a recipient, the system shall create a pending transfer record and notify the recipient.`
- **REPO-14-102:** `When the recipient accepts the transfer, the system shall reassign repository ownership, update paths, and create a redirect.`
- **REPO-14-103:** `When the recipient rejects the transfer, the system shall cancel the pending record and notify the initiator.`
- **REPO-14-104:** `When the initiator cancels a pending transfer, the system shall remove the pending record.`
- **REPO-14-105:** `When the recipient accepts a transfer to an organization and the initiator specified target teams, the system shall grant those teams access to the transferred repository.`
- **REPO-14-106:** `When a transfer completes, the system shall transfer all deploy keys, webhooks, and per-repository configuration to the new owner.`

### State-Driven Requirements (Pending Transfer)

- **REPO-14-701:** `While a transfer is pending, the system shall continue serving the repository from the original path.`
- **REPO-14-702:** `While a transfer is pending, the system shall block the initiator from deleting the repository.`

### Unwanted Behaviour Requirements (Transfer Errors)

- **REPO-14-301:** `If the recipient does not have permission to own repositories of the requested visibility, then the system shall reject the acceptance.`
- **REPO-14-302:** `If a non-owner attempts to initiate or cancel a transfer, then the system shall deny the operation.`
- **REPO-14-303:** `If the recipient namespace already contains a repository with the same name, then the system shall reject the acceptance.`
- **REPO-14-304:** `If the recipient rejects the transfer after the initiator has lost ownership (e.g. account deletion), then the system shall preserve the repository under a ghost owner.`

---

## 15. Branch Rename & Restore (REPO-15)

**User Story:** As a maintainer, I want to rename branches and restore deleted branches so that I can recover from mistakes and reorganize branch structure.

### Ubiquitous Requirements (Branch Rename Properties)

- **REPO-15-001:** `The system shall support renaming any branch in a repository.`
- **REPO-15-002:** `The system shall create a redirect from the old branch name to the new branch name after rename.`
- **REPO-15-003:** `The system shall support restoring recently-deleted branches from their reflog entries.`

### Event-Driven Requirements (Branch Lifecycle)

- **REPO-15-101:** `When a user renames a branch, the system shall update the ref, create a redirect, and update all internal references including default-branch pointer if applicable.`
- **REPO-15-102:** `When a user renames the default branch, the system shall update the repository's HEAD reference atomically.`
- **REPO-15-103:** `When a user requests restoration of a deleted branch, the system shall resurrect the ref from the most recent reflog entry.`
- **REPO-15-104:** `When a user restores a branch via the deleted-branches UI, the system shall recreate the ref and any associated branch-protection rules.`

### Optional Feature Requirements (Rename Constraints)

- **REPO-15-201:** `Where a branch is protected, the system shall require admin permission to rename it.`
- **REPO-15-202:** `Where a renamed branch is referenced by open pull requests, the system shall update the PR base or head references accordingly.`

### Unwanted Behaviour Requirements (Branch Rename Errors)

- **REPO-15-301:** `If a user attempts to rename a branch to a name that already exists, then the system shall reject the rename.`
- **REPO-15-302:** `If a non-admin user attempts to rename a protected branch, then the system shall deny the operation.`
- **REPO-15-303:** `If a user attempts to restore a branch whose reflog entry has expired, then the system shall report that no restorable state exists.`

---

## 16. Repository Activity (REPO-16)

**User Story:** As a maintainer, I want a per-repository activity dashboard so that I can see contribution trends, code frequency, and recent commits at a glance.

### Ubiquitous Requirements (Activity Properties)

- **REPO-16-001:** `The system shall provide an Activity tab on each repository exposing time-bucketed contribution views.`
- **REPO-16-002:** `The system shall render a contributors view showing per-author commit, addition, and deletion counts.`
- **REPO-16-003:** `The system shall render a code-frequency chart showing cumulative additions versus deletions over time.`
- **REPO-16-004:** `The system shall render a recent-commits view showing commit frequency per time bucket.`

### Event-Driven Requirements (Activity Workflow)

- **REPO-16-101:** `When a user navigates to /{owner}/{repo}/activity, the system shall render the activity dashboard for the default time window.`
- **REPO-16-102:** `When a user selects a time period (daily, halfweekly, weekly, monthly, quarterly, semiyearly, or yearly), the system shall re-bucket and re-render the activity charts.`
- **REPO-16-103:** `When a client requests /{owner}/{repo}/activity/code-frequency/data, the system shall return the code-frequency series as JSON.`
- **REPO-16-104:** `When a client requests /{owner}/{repo}/activity/recent-commits/data, the system shall return the recent-commits series as JSON.`
- **REPO-16-105:** `When a user navigates to /{owner}/{repo}/activity/authors, the system shall render the contributor breakdown with author filtering.`

### Optional Feature Requirements (Activity Configuration)

- **REPO-16-201:** `Where the repository is empty, the system shall display an empty-state message in place of activity charts.`
- **REPO-16-202:** `Where the user lacks read access to a private repository, the system shall return 404 for activity routes.`

### Unwanted Behaviour Requirements (Activity Errors)

- **REPO-16-301:** `If the activity data exceeds the maximum rendering window, then the system shall truncate to the most recent N buckets.`
- **REPO-16-302:** `If the activity data endpoint is requested with an invalid time bucket size, then the system shall return 400 Bad Request.`

---

## Configuration Reference

The following INI sections configure code-management behaviors.

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

---

## Business Rules

- **BR-04-001:** Repository names must be unique within an owner's namespace (case-insensitive)
- **BR-04-002:** Repository names are restricted to alphanumeric characters, dots, and dashes
- **BR-04-003:** Reserved repository names (".", "..", "-") and patterns ("*.git", "*.wiki", "*.rss", "*.atom") are prohibited
- **BR-04-004:** Archived repositories are fully read-only: no pushes, no new issues, no new pull requests
- **BR-04-005:** Fork visibility must be at least as restrictive as the base repository
- **BR-04-006:** Mirror repositories are read-only for pull mirrors; push mirrors are outbound only
- **BR-04-007:** Protected branch rules apply to all users except those explicitly on the allowlist
- **BR-04-008:** Permission inheritance is strict: Owner > Admin > Write > Read > None
- **BR-04-009:** Effective permission is the highest level from all sources (collaborator, team, org membership)
- **BR-04-010:** Topic names must match ^[a-z0-9][-.a-z0-9]*$ and be at most 35 characters
- **BR-04-011:** Migration source domains must be on the configured allowlist when ALLOWED_DOMAINS is set
- **BR-04-012:** Repository renames create persistent redirects from the old path to the new path
- **BR-04-013:** Starring and watching are per-user, per-repository actions with no cross-user effects

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
| Merge PR without required approvals | Reject merge until approval count is met |
| Merge PR with rejected review (block enabled) | Reject merge while a rejected review exists |
| Migrate repository exceeding max file size | Abort migration and roll back |
| Migration source unreachable | Abort migration, report error, roll back |
| Add blocked user as collaborator | Reject the addition |
| Star already-starred repository | Ignore duplicate (idempotent) |
| File exceeds MAX_DISPLAY_FILE_SIZE | Display notice instead of content |
| Diff exceeds configured line/file limits | Truncate diff display |
| Cherry-pick with merge conflicts | Report conflicts, do not complete operation |
| Topic name validation failure | Reject topic with validation error |
| 2FA user HTTP Git with password | Reject authentication, require token |

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
- Permission computation for any repository action completes in under 50ms
- Star/watch toggle operations complete within 500ms
- Web editor commit completes within 3 seconds
