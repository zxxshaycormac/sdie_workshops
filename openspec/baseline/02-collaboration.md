# 02 — Collaboration

Baseline specification of Gitea's collaboration subsystems: issues, pull requests, wiki, projects, milestones, labels, reactions, notifications, time tracking, and comments. All requirements describe the current (v1.22.x) system behavior.

---

## 1. Issues (COLL-01)

**User Story:** As a contributor, I want to create and manage issues so that I can report bugs, request features, and coordinate work within a repository.

#### Ubiquitous Requirements (Issue Properties)

- **COLL-01-001:** `The system shall require the Issues unit to be enabled on a repository before issues can be created or viewed.`
- **COLL-01-002:** `The system shall assign sequential issue numbers per repository that are never reused after deletion.`
- **COLL-01-003:** `The system shall store the following properties per issue: title, body (markdown), author, assignees, labels, milestone, state (open/closed), creation timestamp, and update timestamp.`
- **COLL-01-004:** `The system shall soft-delete issues rather than permanently removing them from the database.`
- **COLL-01-005:** `The system shall render issue bodies as markdown with support for cross-references, mentions, and emoji.`
- **COLL-01-006:** `The system shall support multiple assignees per issue.`
- **COLL-01-007:** `The system shall auto-link references of the form #issue, !PR, @user, and commit SHAs within issue text.`

#### Event-Driven Requirements (Issue Lifecycle)

- **COLL-01-101:** `When a user with write permission submits a valid issue form, the system shall create a new issue in the open state.`
- **COLL-01-102:** `When a user closes an issue, the system shall set the issue state to closed and record the closing event with author attribution.`
- **COLL-01-103:** `When a user reopens a closed issue, the system shall set the issue state to open and record the reopening event.`
- **COLL-01-104:** `When an issue is assigned a milestone, the system shall update the milestone's progress calculation to reflect the change.`
- **COLL-01-105:** `When a user with write permission pins an issue, the system shall display the issue at the top of the issue list regardless of sort order.`
- **COLL-01-106:** `When a user with write permission locks an issue, the system shall prevent non-collaborators from adding comments and shall record an optional lock reason.`
- **COLL-01-107:** `When a user unlocks a previously locked issue, the system shall restore commenting ability for all permitted users.`
- **COLL-01-108:** `When a user edits an issue body, the system shall preserve the previous content in the content history.`
- **COLL-01-109:** `When a user creates an issue using an issue template, the system shall pre-populate the issue body with the template content.`

#### Optional Feature Requirements (Issue Dependencies)

- **COLL-01-201:** `Where issue dependencies are enabled, the system shall allow users to declare blocking relationships between issues.`
- **COLL-01-202:** `Where an issue has blocking dependencies that are not yet closed, the system shall display a warning indicator on the dependent issue.`
- **COLL-01-203:** `Where issue templates are defined in .gitea/ISSUE_TEMPLATE/ or .github/ISSUE_TEMPLATE/, the system shall present the template chooser in the new issue form.`

#### Unwanted Behaviour Requirements (Issue Errors)

- **COLL-01-301:** `If a user without write permission attempts to lock an issue, then the system shall deny the operation.`
- **COLL-01-302:** `If a user attempts to create an issue on a repository with the Issues unit disabled, then the system shall return an error.`
- **COLL-01-303:** `If a non-collaborator attempts to comment on a locked issue, then the system shall reject the comment with an explanation.`
- **COLL-01-304:** `If a circular dependency between issues is detected, then the system shall reject the dependency creation.`

---

## 2. Pull Requests (COLL-02)

**User Story:** As a developer, I want to propose, review, and merge code changes so that my contributions are integrated into the project through a controlled workflow.

#### Ubiquitous Requirements (PR Properties)

- **COLL-02-001:** `The system shall require the PullRequests unit to be enabled on a repository before pull requests can be created or viewed.`
- **COLL-02-002:** `The system shall store the following properties per PR: title, body (markdown), author, head branch, base branch, state (open/closed/merged), merge base SHA, and assignees.`
- **COLL-02-003:** `The system shall assign sequential PR numbers per repository (sharing the same counter as issues).`
- **COLL-02-004:** `The system shall support four merge strategies: merge commit, squash merge, rebase merge, and manually merged.`
- **COLL-02-005:** `The system shall detect and display merge conflicts between the head and base branches.`

#### Event-Driven Requirements (PR Lifecycle)

- **COLL-02-101:** `When a user creates a pull request, the system shall compute the diff between the head and base branches and store the merge base.`
- **COLL-02-102:** `When a user merges a PR using merge commit, the system shall create a merge commit on the base branch preserving all commits from the head branch.`
- **COLL-02-103:** `When a user merges a PR using squash merge, the system shall combine all commits into a single commit on the base branch.`
- **COLL-02-104:** `When a user merges a PR using rebase merge, the system shall rebase all commits onto the base branch without creating a merge commit.`
- **COLL-02-105:** `When all required checks pass and approvals are satisfied, the system shall allow the PR to be merged.`
- **COLL-02-106:** `When a user enables auto-merge on a PR, the system shall automatically merge the PR once all required checks and approvals are satisfied.`
- **COLL-02-107:** `When a user requests a reviewer, the system shall send a notification to the requested reviewer and display the pending review status.`
- **COLL-02-108:** `When a reviewer submits an approve review, the system shall record the approval and update the PR's review status.`
- **COLL-02-109:** `When a reviewer submits a request-changes review, the system shall record the rejection and block merging until a new approval is received.`
- **COLL-02-110:** `When a user updates the base branch into the head branch, the system shall merge or rebase the base into the head to resolve divergence.`
- **COLL-02-111:** `When a PR is merged, the system shall offer to delete the head branch if it is no longer needed.`

#### State-Driven Requirements (PR Merge Blocking)

- **COLL-02-701:** `While a PR has merge conflicts, the system shall prevent merging and display the conflicting files.`
- **COLL-02-702:** `While a PR title contains WIP or draft markers, the system shall prevent merging.`
- **COLL-02-703:** `While required status checks are failing or pending, the system shall prevent merging.`

#### Optional Feature Requirements (PR Templates & Advanced)

- **COLL-02-201:** `Where a PR template exists in .gitea/ or .github/, the system shall pre-populate the PR body with the template content.`
- **COLL-02-202:** `Where branch protection rules require signed commits, the system shall verify commit signatures before allowing merge.`
- **COLL-02-203:** `Where auto-merge is configured with at least one check or approval requirement, the system shall queue the PR for automatic merge.`

#### Unwanted Behaviour Requirements (PR Errors)

- **COLL-02-301:** `If a user attempts to merge a PR with conflicts, then the system shall reject the merge and display the conflicting files.`
- **COLL-02-302:** `If a user attempts to merge a draft/WIP PR, then the system shall reject the merge.`
- **COLL-02-303:** `If the base branch has advanced and the head branch cannot be cleanly merged, then the system shall warn the user and suggest updating the base.`
- **COLL-02-304:** `If a fork PR attempts to access Actions secrets from the target repository, then the system shall restrict access to those secrets.`

---

## 3. Wiki (COLL-03)

**User Story:** As a collaborator, I want to create and maintain wiki pages so that project documentation lives alongside the codebase.

#### Ubiquitous Requirements (Wiki Properties)

- **COLL-03-001:** `The system shall store wiki content as a separate Git repository alongside the main repository.`
- **COLL-03-002:** `The system shall render wiki pages as markdown.`
- **COLL-03-003:** `The system shall designate a default wiki page named "Home".`

#### Event-Driven Requirements (Wiki Lifecycle)

- **COLL-03-101:** `When a user creates a new wiki page, the system shall commit the page content to the wiki Git repository.`
- **COLL-03-102:** `When a user edits a wiki page, the system shall record the change as a new commit in the wiki repository preserving full edit history.`
- **COLL-03-103:** `When a user deletes a wiki page, the system shall remove the page from the wiki repository.`
- **COLL-03-104:** `When a user creates a page named "_Sidebar", the system shall display it as navigation in the wiki view.`
- **COLL-03-105:** `When a user creates a page named "_Footer", the system shall display it as footer content in the wiki view.`

#### Optional Feature Requirements (Wiki Configuration)

- **COLL-03-201:** `Where an external wiki URL is configured for a repository, the system shall redirect the wiki tab to the external URL.`
- **COLL-03-202:** `Where the Wiki unit is disabled, the system shall hide the wiki tab from the repository navigation.`

#### Unwanted Behaviour Requirements (Wiki Errors)

- **COLL-03-301:** `If a user without write permission attempts to create or edit a wiki page, then the system shall deny the operation.`
- **COLL-03-302:** `If a wiki page name conflicts with the reserved names "_Sidebar" or "_Footer" used for navigation, then the system shall treat it as a navigation element rather than a regular page.`

---

## 4. Projects / Kanban (COLL-04)

**User Story:** As a project organizer, I want to manage work with boards and columns so that I can visually track progress of issues and pull requests.

#### Ubiquitous Requirements (Project Properties)

- **COLL-04-001:** `The system shall require the Projects unit to be enabled on a repository before projects can be created.`
- **COLL-04-002:** `The system shall support three project scopes: repository-level, organization-level, and user-level.`
- **COLL-04-003:** `The system shall represent issues and pull requests as cards within project columns.`
- **COLL-04-004:** `The system shall store the following properties per project: title, description, scope (repo/org/user), and state (open/closed).`

#### Event-Driven Requirements (Project Lifecycle)

- **COLL-04-101:** `When a user creates a new project, the system shall initialize it with default columns: Backlog, In Progress, and Done.`
- **COLL-04-102:** `When a user adds a column to a project, the system shall append it to the column list.`
- **COLL-04-103:** `When a user moves a card between columns, the system shall update the card's column assignment and sort position.`
- **COLL-04-104:** `When a user reorders columns, the system shall persist the new column ordering.`
- **COLL-04-105:** `When a user closes a project, the system shall set the project state to closed and preserve all column and card data.`
- **COLL-04-106:** `When a user reopens a closed project, the system shall restore the project to the open state with its existing columns and cards.`
- **COLL-04-107:** `When an issue or PR referenced by a card is closed, the system shall update the card's display state accordingly.`

#### Unwanted Behaviour Requirements (Project Errors)

- **COLL-04-301:** `If a user without project management permission attempts to modify a project board, then the system shall deny the operation.`
- **COLL-04-302:** `If a user attempts to delete a column that still contains cards, then the system shall require the user to move or remove the cards first.`
- **COLL-04-303:** `If a user attempts to add a card for an issue or PR from a different repository to a repository-scoped project, then the system shall reject the addition.`

---

## 5. Milestones (COLL-05)

**User Story:** As a project manager, I want to define milestones with deadlines so that I can track progress toward goals across issues and pull requests.

#### Ubiquitous Requirements (Milestone Properties)

- **COLL-05-001:** `The system shall store the following properties per milestone: name, description, due date, state (open/closed), and repository reference.`
- **COLL-05-002:** `The system shall calculate milestone progress as the ratio of closed issues to total issues assigned to the milestone.`

#### Event-Driven Requirements (Milestone Lifecycle)

- **COLL-05-101:** `When a user creates a milestone with a due date, the system shall track and display the due date on the milestone list.`
- **COLL-05-102:** `When all issues assigned to a milestone are closed, the system shall display the milestone as fully complete (100% progress).`
- **COLL-05-103:** `When a user closes a milestone, the system shall set its state to closed without affecting the state of its assigned issues.`
- **COLL-05-104:** `When a user reopens a closed milestone, the system shall set its state to open.`
- **COLL-05-105:** `When a user deletes a milestone, the system shall unassign all issues and PRs from the milestone before removal.`

#### State-Driven Requirements (Overdue Tracking)

- **COLL-05-701:** `While a milestone's due date has passed and the milestone is still open, the system shall display an overdue indicator on the milestone.`

#### Unwanted Behaviour Requirements (Milestone Errors)

- **COLL-05-301:** `If a user without write permission attempts to create or modify a milestone, then the system shall deny the operation.`
- **COLL-05-302:** `If a milestone name duplicates an existing milestone name in the same repository, then the system shall reject the creation.`

---

## 6. Labels (COLL-06)

**User Story:** As a repository maintainer, I want to create and apply labels so that issues and pull requests are categorized and filterable.

#### Ubiquitous Requirements (Label Properties)

- **COLL-06-001:** `The system shall store the following properties per label: name, color (hex), description, and scope (repository or organization).`
- **COLL-06-002:** `The system shall support two label scopes: repository-specific labels and organization-level labels shared across all org repos.`
- **COLL-06-003:** `The system shall allow multiple labels per issue or pull request.`

#### Event-Driven Requirements (Label Lifecycle)

- **COLL-06-101:** `When a user creates a label, the system shall store the label with its name, color, and optional description.`
- **COLL-06-102:** `When a user applies a label to an issue or PR, the system shall record a label-added event in the issue/PR timeline.`
- **COLL-06-103:** `When a user removes a label from an issue or PR, the system shall record a label-removed event.`
- **COLL-06-104:** `When a user deletes a label, the system shall remove it from all issues and PRs that had the label applied.`
- **COLL-06-105:** `When a user initializes labels from a template, the system shall create all labels defined in the selected template set.`

#### Optional Feature Requirements (Exclusive Labels)

- **COLL-06-201:** `Where a label name contains a forward slash (e.g., "priority/high"), the system shall treat the prefix before the slash as a scope and enforce exclusivity such that only one label per scope may be applied to a single issue or PR.`
- **COLL-06-202:** `Where a user applies a scoped label that conflicts with an existing label in the same scope, the system shall replace the existing scoped label with the new one.`

#### Unwanted Behaviour Requirements (Label Errors)

- **COLL-06-301:** `If a user without write permission attempts to create, edit, or delete a label, then the system shall deny the operation.`
- **COLL-06-302:** `If a label name duplicates an existing label in the same scope (repo or org), then the system shall reject the creation.`
- **COLL-06-303:** `If a user attempts to apply an organization-level label to an issue in a repository outside that organization, then the system shall reject the application.`

---

## 7. Reactions (COLL-07)

**User Story:** As a user, I want to react to issues and comments with emoji so that I can express feedback without writing a full comment.

#### Ubiquitous Requirements (Reaction Properties)

- **COLL-07-001:** `The system shall support a configurable set of emoji reaction types: +1, -1, laugh, hooray, confused, heart, rocket, and eyes.`
- **COLL-07-002:** `The system shall allow each user to add at most one reaction of a given emoji type per reactable item.`
- **COLL-07-003:** `The system shall record user attribution for each reaction.`
- **COLL-07-004:** `The system shall support reactions on issues, issue comments, PR comments, and review comments.`

#### Event-Driven Requirements (Reaction Lifecycle)

- **COLL-07-101:** `When a user adds a reaction, the system shall store the reaction with the user ID, emoji type, and target item reference.`
- **COLL-07-102:** `When a user removes their own reaction, the system shall delete the reaction record and update the displayed count.`

#### Unwanted Behaviour Requirements (Reaction Errors)

- **COLL-07-301:** `If a user attempts to add a duplicate reaction (same user, same emoji, same item), then the system shall reject the addition.`
- **COLL-07-302:** `If a user attempts to add a reaction to a locked issue, then the system shall reject the reaction.`
- **COLL-07-303:** `If a user attempts to use a reaction type not in the supported set, then the system shall reject the reaction.`

---

## 8. Notifications (COLL-08)

**User Story:** As a user, I want to receive and manage notifications about activity across repositories so that I can stay informed about relevant changes.

#### Ubiquitous Requirements (Notification Properties)

- **COLL-08-001:** `The system shall generate notifications for the following sources: issues, pull requests, commits, repository events, mentions, and assignments.`
- **COLL-08-002:** `The system shall track notification states: unread, read, and pinned.`
- **COLL-08-003:** `The system shall support three repository subscription modes: all events, participating only, and ignore.`
- **COLL-08-004:** `The system shall deliver notifications via web interface and optionally via email.`

#### Event-Driven Requirements (Notification Lifecycle)

- **COLL-08-101:** `When an event occurs that matches a user's subscription, the system shall create a notification for that user.`
- **COLL-08-102:** `When a user is @mentioned in an issue or PR, the system shall create a notification regardless of subscription mode.`
- **COLL-08-103:** `When a user is assigned to an issue or PR, the system shall create a notification.`
- **COLL-08-104:** `When a user reads a notification, the system shall mark it as read.`
- **COLL-08-105:** `When a user marks all notifications as read, the system shall set all unread notifications to read state.`
- **COLL-08-106:** `When a user watches a repository, the system shall subscribe the user to notifications for that repository.`
- **COLL-08-107:** `When a user unwatches a repository, the system shall stop generating notifications for that repository.`
- **COLL-08-108:** `When a user comments on an issue or PR, the system shall auto-subscribe the user to that issue or PR.`
- **COLL-08-109:** `When email notifications are enabled and a notification is created, the system shall send an email to the user.`

#### Optional Feature Requirements (Notification Configuration)

- **COLL-08-201:** `Where ENABLE_NOTIFY_MAIL is enabled, the system shall send email notifications for subscribed events.`
- **COLL-08-202:** `Where AUTO_WATCH_ON_CHANGES is enabled, the system shall automatically watch repositories when a user pushes to them.`
- **COLL-08-203:** `Where AUTO_WATCH_REPOS is enabled, the system shall automatically watch repositories when a user creates or forks them.`

#### Unwanted Behaviour Requirements (Notification Errors)

- **COLL-08-301:** `If a user with notifications disabled attempts to view notifications, then the system shall display an empty notification list.`
- **COLL-08-302:** `If a user has ignored a repository, then the system shall not create notifications for that user from that repository even on mention.`

---

## 9. Timetracking (COLL-09)

**User Story:** As a contributor, I want to track time spent on issues so that effort is visible and reportable.

#### Ubiquitous Requirements (Timetracking Properties)

- **COLL-09-001:** `The system shall require the Issues unit to be enabled for time tracking features to be available.`
- **COLL-09-002:** `The system shall allow at most one active stopwatch per user at any time.`
- **COLL-09-003:** `The system shall store tracked time entries with the following properties: user, duration, creation timestamp, and optional description.`

#### Event-Driven Requirements (Timetracking Lifecycle)

- **COLL-09-101:** `When a user starts a stopwatch on an issue, the system shall record the start time and associate it with the user and issue.`
- **COLL-09-102:** `When a user stops a stopwatch on an issue, the system shall compute the elapsed time and create a tracked time entry.`
- **COLL-09-103:** `When a user manually adds time to an issue, the system shall create a tracked time entry with the specified duration and description.`
- **COLL-09-104:** `When a user deletes a tracked time entry, the system shall remove it and update the issue's total time.`
- **COLL-09-105:** `When a user sets a time estimate on an issue, the system shall store the estimate and display it alongside the tracked time.`

#### Unwanted Behaviour Requirements (Timetracking Errors)

- **COLL-09-301:** `If a user with an active stopwatch on one issue attempts to start a stopwatch on another issue, then the system shall reject the new stopwatch and prompt the user to stop the existing one first.`
- **COLL-09-302:** `If a user attempts to add a negative time duration, then the system shall reject the entry.`
- **COLL-09-303:** `If a user without write permission attempts to delete another user's time entry, then the system shall deny the operation.`

---

## 10. Comments (COLL-10)

**User Story:** As a collaborator, I want to comment on issues and pull requests so that I can participate in discussions, provide feedback, and record contextual events.

#### Ubiquitous Requirements (Comment Properties)

- **COLL-10-001:** `The system shall support the following comment types: plain text/markdown discussion, code review comment, review summary (approve/request-changes/comment), and system event comments (label change, milestone change, assignee change, title change, lock/unlock, dependency change, time tracking event, reopen/close, project change).`
- **COLL-10-002:** `The system shall render comment bodies as markdown with support for cross-references, mentions, emoji, and file attachments.`
- **COLL-10-003:** `The system shall record the author and creation timestamp for every comment.`
- **COLL-10-004:** `The system shall support file attachments in comments via upload.`

#### Event-Driven Requirements (Comment Lifecycle)

- **COLL-10-101:** `When a user submits a comment on an issue or PR, the system shall create a comment record and notify subscribed users.`
- **COLL-10-102:** `When a user edits a comment, the system shall preserve the previous content in the content history.`
- **COLL-10-103:** `When a user deletes a comment, the system shall remove the comment and any associated reactions.`
- **COLL-10-104:** `When a reviewer submits a code review comment on a specific diff line, the system shall store the comment with the file path, line number, and diff hash reference.`
- **COLL-10-105:** `When a reviewer submits a review summary with verdict (approve/request-changes/comment), the system shall record the verdict and any summary comment.`
- **COLL-10-106:** `When a system event occurs (label added, assignee changed, milestone set, title changed, etc.), the system shall create an event comment in the issue or PR timeline.`

#### Optional Feature Requirements (Comment Features)

- **COLL-10-201:** `Where a comment references a specific commit SHA, the system shall auto-link the reference to the commit detail view.`
- **COLL-10-202:** `Where a comment references an issue or PR number (e.g., #42, !7), the system shall auto-link the reference to the corresponding entity.`

#### Unwanted Behaviour Requirements (Comment Errors)

- **COLL-10-301:** `If a user attempts to comment on a locked issue, then the system shall reject the comment unless the user has write permission.`
- **COLL-10-302:** `If a user attempts to edit another user's comment, then the system shall deny the operation unless the editor is an admin or the repository owner.`
- **COLL-10-303:** `If a user attempts to delete another user's comment, then the system shall deny the operation unless the deleter is an admin or the repository owner.`
- **COLL-10-304:** `If an uploaded attachment exceeds the configured maximum file size, then the system shall reject the upload.`

---

## Business Rules

- **BR-02-001:** Issue and PR numbers share a single sequential counter per repository and are never reused after deletion.
- **BR-02-002:** Issues are soft-deleted; deleted issues remain in the database but are not visible in the UI.
- **BR-02-003:** PR merge strategy selection is available to users with write permission; repository settings may restrict which strategies are permitted.
- **BR-02-004:** Wiki content is stored as a separate Git repository and is fully cloneable with standard Git tools.
- **BR-02-005:** Project cards are references to issues and PRs, not copies; changes to the underlying issue or PR are reflected on the card.
- **BR-02-006:** Milestone progress is auto-calculated as closed_issues / total_issues assigned to the milestone.
- **BR-02-007:** Exclusive/scoped labels (format "scope/name") enforce at most one label per scope on a given issue or PR.
- **BR-02-008:** Each user may have at most one reaction of a given emoji type per reactable item.
- **BR-02-009:** A user may have at most one active stopwatch at a time across all issues and repositories.
- **BR-02-010:** Notification delivery respects the user's subscription mode and the repository's watch configuration.
- **BR-02-011:** WIP/draft PRs cannot be merged regardless of user permissions.
- **BR-02-012:** Comment editing is restricted to the original author, repository owner, or admin.
- **BR-02-013:** Auto-merge requires at least one status check or approval rule to be configured.
- **BR-02-014:** Fork PRs cannot access Actions secrets from the target repository.

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Issue number reuse after deletion | Reject; numbers are never reused |
| Merge PR with unresolved conflicts | Reject merge, display conflicting files |
| Comment on locked issue by non-collaborator | Reject comment, display lock reason |
| Start second stopwatch while one is active | Reject, prompt user to stop existing stopwatch first |
| Apply exclusive label when one from same scope exists | Replace existing scoped label with new one |
| Create circular issue dependency | Reject dependency creation |
| Wiki page name collision with _Sidebar/_Footer | Treat as navigation element, not regular page |
| Delete project column with cards remaining | Require moving or removing cards first |
| PR base branch diverged from head | Warn user, suggest updating base into head |
| Auto-merge when checks never complete | PR remains queued until checks resolve or auto-merge is cancelled |
| Add duplicate reaction (same user, emoji, item) | Reject the addition silently |
| Mention user who has ignored the repository | Do not create notification for that user |
| Upload attachment exceeding max file size | Reject upload with size limit error |
| Edit another user's comment (non-admin) | Deny operation with permission error |
| External wiki configured on repository | Redirect wiki tab to the external URL |

## Success Criteria

- Issue creation completes within 2 seconds for standard markdown bodies
- PR diff computation completes within 5 seconds for changes up to 1000 files
- Wiki page rendering completes within 1 second per page
- Project board with 200 cards loads within 3 seconds
- Milestone progress calculation is real-time and reflects issue state changes immediately
- Label filtering on issue lists returns results within 1 second for repositories with up to 10,000 issues
- Notification delivery occurs within 10 seconds of the triggering event
- Email notification dispatch occurs within 30 seconds when ENABLE_NOTIFY_MAIL is enabled
- Stopwatch time tracking is accurate to within 1 second
- Comment rendering with markdown, mentions, and cross-references completes within 500ms
- Reaction add/remove operations complete within 200ms
