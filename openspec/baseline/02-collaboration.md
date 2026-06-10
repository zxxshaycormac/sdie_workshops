# 02 — Collaboration

## 1. Issues

### What
Bug tracking, task management, and discussion system for repositories.

### Behaviors
- Full CRUD: create, read, update, delete issues
- Assignees: multiple users per issue
- Labels: repo-level and org-level labels
- Milestones: due dates, progress tracking (open/closed ratio)
- Reactions: emoji reactions on issues
- Pinning: important issues stay at top of list
- Locking: prevent comments from non-collaborators (with reason)
- Content history: track all edits to issue body
- Templates: issue templates (bug report, feature request, custom)
- Cross-references: auto-link @user, #issue, !PR references
- Dependencies: block/unblock issues on other issues
- Watch/unwatch: per-issue notification subscription
- Issue search: filters by status, labels, milestones, assignees, author, mentioned, time range

### UI Routes
- GET /{owner}/{repo}/issues — Issue list
- GET /{owner}/{repo}/issues/new — New issue form
- GET /{owner}/{repo}/issues/{index} — Issue detail
- POST /{owner}/{repo}/issues — Create issue
- POST /{owner}/{repo}/issues/{index} — Update issue
- POST /{owner}/{repo}/issues/{index}/pin — Pin issue
- POST /{owner}/{repo}/issues/{index}/lock — Lock issue
- GET /issues — User's issues across repos
- GET /{owner}/{repo}/issues/{index}/content_history — Edit history

### API Endpoints
- GET /repos/{owner}/{repo}/issues — List issues
- POST /repos/{owner}/{repo}/issues — Create issue
- GET /repos/{owner}/{repo}/issues/{index} — Get issue
- PATCH /repos/{owner}/{repo}/issues/{index} — Update issue
- DELETE /repos/{owner}/{repo}/issues/{index} — Delete issue
- GET /repos/issues/search — Global issue search
- GET /repos/{owner}/{repo}/issues/{index}/labels — Issue labels
- POST /repos/{owner}/{repo}/issues/{index}/labels — Set labels
- GET /repos/{owner}/{repo}/issues/{index}/comments — List comments
- POST /repos/{owner}/{repo}/issues/{index}/comments — Add comment

### Constraints
- Requires Issues unit enabled on repo
- Locking requires write permission
- Deleted issues are soft-deleted
- Issue numbers are sequential per repo, never reused

---

## 2. Pull Requests

### What
Code review and merge workflow: propose changes, review, discuss, and merge.

### Behaviors
- Creation: from branch compare view or API
- Merge strategies:
  - Merge commit (default)
  - Squash merge (combine all commits)
  - Rebase merge (rebase onto base)
  - Manually merged (mark as merged externally)
- Conflict detection: automatic, shown in UI
- Review workflow:
  - Approve: reviewer approves changes
  - Request changes: reviewer requests modifications
  - Comment: reviewer leaves feedback without formal verdict
- Review comments:
  - Line-level comments on specific diff lines
  - Single comments (general, not on diff)
  - Review summary comments
- Auto-merge: schedule merge when all checks pass
- PR templates: default body template from .gitea/ or .github/
- Update base: merge base branch into head (merge or rebase)
- Diff viewer: file-by-file with whitespace options
- WIP detection: PRs with WIP in title cannot be merged
- Branch deletion: option to delete head branch after merge
- Commit listing: all commits in the PR

### UI Routes
- GET /{owner}/{repo}/pulls — PR list
- GET /{owner}/{repo}/pulls/{index} — PR detail
- GET /{owner}/{repo}/compare/{base}...{head} — Compare/PR creation
- POST /{owner}/{repo}/pulls — Create PR
- POST /{owner}/{repo}/pulls/{index}/merge — Merge PR
- GET /{owner}/{repo}/pulls/{index}/files — Diff view
- GET /{owner}/{repo}/pulls/{index}/commits — Commit list

### API Endpoints
- GET /repos/{owner}/{repo}/pulls — List PRs
- POST /repos/{owner}/{repo}/pulls — Create PR
- GET /repos/{owner}/{repo}/pulls/{index} — Get PR
- PATCH /repos/{owner}/{repo}/pulls/{index} — Update PR
- POST /repos/{owner}/{repo}/pulls/{index}/merge — Merge
- GET /repos/{owner}/{repo}/pulls/{index}/commits — List commits
- GET /repos/{owner}/{repo}/pulls/{index}/files — List changed files
- POST /repos/{owner}/{repo}/pulls/{index}/update — Update base into head
- POST /repos/{owner}/{repo}/pulls/{index}/requested_reviewers — Request review
- DELETE /repos/{owner}/{repo}/pulls/{index}/requested_reviewers — Remove review request

### Constraints
- Requires PullRequests unit enabled
- Merge blocked by: conflicts, WIP title, missing approvals, failing status checks, unsigned commits (if required)
- Auto-merge requires at least one check or approval configured
- Fork PRs have restricted permissions for Actions secrets

---

## 3. Wiki

### What
Collaborative documentation stored as a separate Git repository alongside the main repo.

### Behaviors
- Wiki is a Git repo: full history, cloneable
- Pages: markdown-rendered wiki pages
- Sidebar: custom navigation
- Footer: custom footer content
- Page CRUD: create, edit, delete pages
- Edit history: full git history of changes
- Default page: "Home"

### UI Routes
- GET /{owner}/{repo}/wiki — Wiki home
- GET /{owner}/{repo}/wiki/{page} — View page
- GET /{owner}/{repo}/wiki/_new — New page
- POST /{owner}/{repo}/wiki/_new — Create page
- GET /{owner}/{repo}/wiki/_edit/{page} — Edit page
- POST /{owner}/{repo}/wiki/{page} — Update page
- POST /{owner}/{repo}/wiki/{page}/delete — Delete page
- GET /{owner}/{repo}/wiki/_pages — Page list

### API Endpoints
- GET /repos/{owner}/{repo}/wiki — List pages
- GET /repos/{owner}/{repo}/wiki/{page} — Get page
- POST /repos/{owner}/{repo}/wiki/new — Create page
- PATCH /repos/{owner}/{repo}/wiki/{page} — Edit page
- DELETE /repos/{owner}/{repo}/wiki/{page} — Delete page

### Constraints
- Requires Wiki unit enabled (or ExternalWiki for redirect)
- Wiki repo is separate from main repo in storage

---

## 4. Projects (Kanban)

### What
Visual project management with boards, columns, and cards for organizing issues and PRs.

### Behaviors
- Project types: Repository, Organization, User
- Columns: create, edit, delete, reorder
- Cards: issues/PRs as cards, drag-and-drop between columns
- Default columns: Backlog, In Progress, Done (on creation)
- Project lifecycle: open/close/reopen
- Sorting within columns

### UI Routes
- GET /{owner}/{repo}/projects — List projects
- GET /{owner}/{repo}/projects/new — Create project
- GET /{owner}/{repo}/projects/{id} — View board
- POST /{owner}/{repo}/projects — Create project
- POST /{owner}/{repo}/projects/{id} — Update project
- GET /:org/projects — Org projects
- GET /user/projects — User projects

### API Endpoints
- GET /repos/{owner}/{repo}/projects — List projects
- POST /repos/{owner}/{repo}/projects — Create project
- GET /repos/{owner}/{repo}/projects/{id} — Get project
- PATCH /repos/{owner}/{repo}/projects/{id} — Update project
- DELETE /repos/{owner}/{repo}/projects/{id} — Delete project
- GET /repos/{owner}/{repo}/projects/{id}/columns — List columns
- POST /repos/{owner}/{repo}/projects/{id}/columns — Create column
- PATCH /repos/{owner}/{repo}/projects/{id}/columns/{column} — Update column
- DELETE /repos/{owner}/{repo}/projects/{id}/columns/{column} — Delete column
- GET /repos/{owner}/{repo}/projects/{id}/columns/{column}/cards — List cards
- POST /repos/{owner}/{repo}/projects/{id}/columns/{column}/cards — Add card

### Constraints
- Requires Projects unit enabled
- Cards are references to issues/PRs (not copies)

---

## 5. Milestones

### What
Goal tracking with deadlines and progress measurement for issues and PRs.

### Behaviors
- CRUD: create, edit, close, reopen, delete
- Due dates with overdue tracking
- Progress: auto-calculated from open/closed issue count
- Issues assigned to milestones
- Milestone-level issue list

### UI Routes
- GET /{owner}/{repo}/milestones — List milestones
- GET /{owner}/{repo}/milestones/new — Create milestone
- POST /{owner}/{repo}/milestones — Create milestone
- POST /{owner}/{repo}/milestones/{id} — Update milestone

### API Endpoints
- GET /repos/{owner}/{repo}/milestones — List milestones
- POST /repos/{owner}/{repo}/milestones — Create milestone
- GET /repos/{owner}/{repo}/milestones/{id} — Get milestone
- PATCH /repos/{owner}/{repo}/milestones/{id} — Update milestone
- DELETE /repos/{owner}/{repo}/milestones/{id} — Delete milestone

---

## 6. Labels

### What
Categorization system for issues and PRs with colors and exclusive/scoped labels.

### Behaviors
- CRUD: create, edit, delete labels
- Two levels: repo-specific and org-level (shared across org repos)
- Exclusive/scoped labels: format "scope/label", only one label per scope allowed
- Color customization (hex)
- Label templates: initialize from predefined sets
- Label application: add/remove from issues and PRs

### UI Routes
- GET /{owner}/{repo}/labels — Repo labels
- POST /{owner}/{repo}/labels — Create label
- POST /{owner}/{repo}/labels/{id}/edit — Update label
- POST /{owner}/{repo}/labels/{id}/delete — Delete label
- POST /{owner}/{repo}/labels/initialize — Init from template
- GET /:org/settings/labels — Org labels

### API Endpoints
- GET /repos/{owner}/{repo}/labels — List labels
- POST /repos/{owner}/{repo}/labels — Create label
- GET /repos/{owner}/{repo}/labels/{id} — Get label
- PATCH /repos/{owner}/{repo}/labels/{id} — Update label
- DELETE /repos/{owner}/{repo}/labels/{id} — Delete label
- GET /orgs/{org}/labels — Org labels
- POST /orgs/{org}/labels — Create org label

---

## 7. Reactions

### What
Emoji feedback on issues, comments, and PRs.

### Behaviors
- Supported emojis: configurable set (+1, -1, laugh, hooray, confused, heart, rocket, eyes)
- One reaction per user per item
- User attribution on each reaction
- Appears on: issues, issue comments, PR comments, review comments

### API Endpoints
- GET /repos/{owner}/{repo}/issues/{index}/reactions — List reactions
- POST /repos/{owner}/{repo}/issues/{index}/reactions — Add reaction
- DELETE /repos/{owner}/{repo}/issues/{index}/reactions/{id} — Remove reaction
- GET /repos/{owner}/{repo}/issues/comments/{id}/reactions — Comment reactions
- POST /repos/{owner}/{repo}/issues/comments/{id}/reactions — Add comment reaction

---

## 8. Notifications

### What
Activity tracking and alerts for events across repositories.

### Behaviors
- Notification sources: issues, PRs, commits, repos, mentions, assignments
- Watch/unwatch: per-repository subscription
- Subscription modes for repos:
  - All events
  - Participating only (mentioned, assigned, author)
  - Ignore
- Watch modes for issues/PRs: auto-subscribe on comment
- Notification states: unread, read, pinned
- Mark as read/unread, mark all as read
- Email notifications (configurable per user)
- Web notifications with real-time updates (EventSource)

### UI Routes
- GET /notifications — Notification list
- POST /notifications — Update notifications
- GET /notifications/subscriptions — Watched repos
- GET /:org/settings/notification — Org notification settings

### API Endpoints
- GET /notifications — List notifications
- GET /repos/{owner}/{repo}/notifications — Repo notifications
- PUT /notifications — Mark as read
- PUT /notifications/{id} — Mark single as read
- GET /repos/{owner}/{repo}/subscription — Check watch status
- PUT /repos/{owner}/{repo}/subscription — Watch repo
- DELETE /repos/{owner}/{repo}/subscription — Unwatch repo

### Config
- [service] ENABLE_NOTIFY_MAIL — email notifications
- [service] AUTO_WATCH_ON_CHANGES — auto-watch on push
- [service] AUTO_WATCH_REPOS — auto-watch repos

---

## 9. Timetracking

### What
Time tracking for issues with stopwatch and manual time logging.

### Behaviors
- Stopwatch: start/stop timer on issues
- Tracked time: manual time addition with description
- Time estimates: set expected time on issues
- Time reports: by issue, by user
- Add/delete time entries
- Timer tracking across sessions

### UI Routes
- POST /{owner}/{repo}/issues/{index}/times/stopwatch — Start/stop stopwatch
- GET /{owner}/{repo}/issues/{index}/times — Time log

### API Endpoints
- GET /repos/{owner}/{repo}/issues/{index}/times — List tracked times
- POST /repos/{owner}/{repo}/issues/{index}/times — Add time
- DELETE /repos/{owner}/{repo}/issues/{index}/times/{id} — Delete time
- POST /repos/{owner}/{repo}/issues/{index}/stopwatch/start — Start stopwatch
- POST /repos/{owner}/{repo}/issues/{index}/stopwatch/stop — Stop stopwatch
- GET /repos/{owner}/{repo}/issues/{index}/stopwatch — Get stopwatch

### Constraints
- Only one stopwatch per user at a time
- Requires Issues unit enabled

---

## 10. Comments

### What
Discussion system for issues and PRs with multiple comment types.

### Comment Types
- Comment: plain text/markdown discussion
- Code: line-level code review comment
- Review: general PR review feedback (approve/request-changes/comment)
- Label: label change event
- Milestone: milestone change event
- Assignee: assignee change event
- Title: title change event
- Lock/unlock: issue lock state change
- Dependency: dependency added/removed event
- Time tracking: time add/delete event
- Reopen/close: status change event
- Project: project board change event

### Behaviors
- Markdown rendering with all markup features
- Attachments: file uploads in comments
- Content history: track edits to comment body
- Cross-references: auto-link @user, #issue, !PR
- Reactions: emoji reactions on comments
- References: comment references to specific diff lines

### API Endpoints
- GET /repos/{owner}/{repo}/issues/{index}/comments — List comments
- POST /repos/{owner}/{repo}/issues/{index}/comments — Create comment
- PATCH /repos/{owner}/{repo}/issues/comments/{id} — Update comment
- DELETE /repos/{owner}/{repo}/issues/comments/{id} — Delete comment
- GET /repos/{owner}/{repo}/issues/comments/{id}/reactions — List reactions
- GET /repos/{owner}/{repo}/issues/{index}/content_history — Content history
