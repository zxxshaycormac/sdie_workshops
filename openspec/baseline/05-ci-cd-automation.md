# 05 — CI/CD & Automation

Baseline specification of Gitea's CI/CD, webhook, automation, and background processing subsystems. All requirements describe the current (v1.22.x) system behavior.

---

## 1. Gitea Actions

**User Story:** As a developer, I want to define automated CI/CD workflows in YAML so that my code is built, tested, and deployed automatically on every push or pull request.

### Ubiquitous Requirements (Workflow Properties)

- **CI-01-001:** `The system shall parse workflow definitions from YAML files stored in the .gitea/workflows/ or .github/workflows/ directory of a repository.`
- **CI-01-002:** `The system shall support the following workflow trigger events: push, pull_request, schedule, workflow_dispatch, and repository_dispatch.`
- **CI-01-003:** `The system shall track workflow run status through the lifecycle: pending, waiting, running, success, failure, and cancelled.`
- **CI-01-004:** `The system shall support parallel job execution within a single workflow run.`
- **CI-01-005:** `The system shall support matrix build strategies that fan out jobs across defined parameter combinations.`
- **CI-01-006:** `The system shall store workflow run logs with persistent storage and indexed access for retrieval.`
- **CI-01-007:** `The system shall generate a unique task token for each job execution distributed to a runner.`
- **CI-01-008:** `The system shall require the actions unit (unit.TypeActions) to be enabled on the repository before processing workflows.`

### Event-Driven Requirements (Workflow Execution)

- **CI-01-101:** `When a trigger event matches a workflow's on clause, the system shall create a new workflow run.`
- **CI-01-102:** `When a workflow run is created, the system shall evaluate job dependencies and schedule jobs for execution.`
- **CI-01-103:** `When a job is assigned to a runner, the system shall execute steps sequentially within that job.`
- **CI-01-104:** `When a workflow run completes, the system shall update the run status to success or failure based on job outcomes.`
- **CI-01-105:** `When a user cancels a workflow run via the API, the system shall cancel all running and pending jobs in that run.`
- **CI-01-106:** `When a user reruns a workflow via the API, the system shall create a new run with the same workflow definition and trigger event.`
- **CI-01-107:** `When a workflow step references an action (uses:), the system shall resolve the action from the configured DEFAULT_ACTIONS_URL (GitHub or self-hosted).`

### State-Driven Requirements (Runner Lifecycle)

- **CI-01-701:** `While a runner is in active state, the system shall assign matching jobs to that runner based on label selection.`
- **CI-01-702:** `While a runner is in offline state, the system shall not assign new jobs to that runner.`
- **CI-01-703:** `While a task exceeds the configured ZOMBIE_TASK_TIMEOUT (default 10 minutes), the system shall mark the task as timed out.`
- **CI-01-704:** `While a task exceeds the configured ENDLESS_TASK_TIMEOUT (default 3 hours), the system shall stop the task.`

### Optional Feature Requirements (Actions Configuration)

- **CI-01-201:** `Where the actions feature is enabled via configuration, the system shall expose workflow, run, job, and artifact endpoints.`
- **CI-01-202:** `Where artifact storage is configured (local or MinIO/S3), the system shall store and retrieve workflow artifacts using the configured backend.`
- **CI-01-203:** `Where scheduled workflows are defined with cron syntax, the system shall trigger workflow runs according to the specified schedule.`
- **CI-01-204:** `Where a status badge endpoint is requested for a workflow, the system shall return an SVG image reflecting the latest run status.`
- **CI-01-205:** `Where the actions feature is disabled via configuration, the system shall hide all actions UI and API endpoints.`

### Unwanted Behaviour Requirements (Actions Errors)

- **CI-01-301:** `If a fork pull request workflow attempts to access repository secrets, then the system shall restrict access to protected secrets.`
- **CI-01-302:** `If no runner is available with matching labels for a job, then the system shall keep the job in pending state until a matching runner registers.`
- **CI-01-303:** `If an artifact exceeds the configured retention period (default 90 days), then the system shall delete the artifact.`
- **CI-01-304:** `If a job remains abandoned beyond the configured ABANDONED_JOB_TIMEOUT (default 24 hours), then the system shall cancel the job.`
- **CI-01-305:** `If the actions unit is disabled for a repository, then the system shall reject workflow trigger events for that repository.`

### Complex Requirements (Runner Registration)

- **CI-01-901:** `When a runner registration request is received with a valid registration token, the system shall create a runner record with a unique UUID, assign the specified labels, and associate the runner with the appropriate scope (global, organization, or repository).`

---

## 2. Webhooks

**User Story:** As a repository administrator, I want to configure webhook callbacks so that external systems are notified when repository events occur.

### Ubiquitous Requirements (Webhook Properties)

- **CI-02-001:** `The system shall support the following webhook payload formats: Gitea (native), Gogs, Slack, Discord, Dingtalk, Telegram, Microsoft Teams, Feishu (Lark), Matrix, WeChat Work, and Packagist.`
- **CI-02-002:** `The system shall track delivery status for each webhook invocation with success or failure indication.`
- **CI-02-003:** `The system shall store both request and response data for each webhook delivery.`
- **CI-02-004:** `The system shall support webhook configuration at three scopes: system (admin), organization, and repository.`
- **CI-02-005:** `The system shall support proxy configuration for webhook delivery in restricted network environments.`

### Event-Driven Requirements (Webhook Delivery)

- **CI-02-101:** `When a subscribed repository event occurs, the system shall enqueue an asynchronous webhook delivery to each configured hook endpoint.`
- **CI-02-102:** `When a webhook delivery is triggered, the system shall compute an HMAC-SHA256 signature of the payload and include it in the request headers.`
- **CI-02-103:** `When a webhook delivery fails, the system shall record the failure and allow manual redelivery.`
- **CI-02-104:** `When a user triggers a test delivery via the UI or API, the system shall send a test payload to the webhook endpoint.`
- **CI-02-105:** `When a user requests redelivery of a past webhook, the system shall resend the original payload to the endpoint.`
- **CI-02-106:** `When a webhook delivery times out, the system shall record the timeout as a delivery failure.`
- **CI-02-107:** `When the following repository events occur, the system shall trigger webhook delivery: push, create, delete, release, pull request, pull request review, pull request review comment, issues, issue comment, issue label, issue milestone, issue assign, repository, and package.`

### Optional Feature Requirements (Webhook Security)

- **CI-02-201:** `Where a webhook host allowlist is configured, the system shall reject webhook deliveries to hosts not on the allowlist.`
- **CI-02-202:** `Where TLS verification is skipped via configuration, the system shall deliver webhooks without validating the server certificate.`
- **CI-02-203:** `Where a proxy URL is configured, the system shall route webhook deliveries through the specified proxy.`

### Unwanted Behaviour Requirements (Webhook Errors)

- **CI-02-301:** `If a webhook delivery target is not in the allowed host list, then the system shall reject the delivery with a security error.`
- **CI-02-302:** `If the webhook queue exceeds the configured QUEUE_LENGTH (default 1000), then the system shall apply backpressure to prevent unbounded memory growth.`
- **CI-02-303:** `If a webhook endpoint returns an HTTP error status, then the system shall record the delivery as failed with the response details.`

---

## 3. Agit Flow

**User Story:** As a contributor, I want to create a pull request by pushing to a special ref so that I can propose changes without using the web UI or API.

### Ubiquitous Requirements (Agit Properties)

- **CI-03-001:** `The system shall recognize refs/for/<base-branch>/<topic-branch> as a special push target that triggers pull request creation.`
- **CI-03-002:** `The system shall prefix topic branches with the pushing user's identifier to avoid naming conflicts.`

### Event-Driven Requirements (Agit Workflow)

- **CI-03-101:** `When a user pushes to refs/for/<base-branch>/<topic-branch>, the system shall create a pull request from the topic branch targeting the specified base branch.`
- **CI-03-102:** `When an agit flow push is received, the system shall create the topic branch and pull request in a single atomic operation.`
- **CI-03-103:** `When force push is configured for the repository, the system shall allow force pushing to agit flow refs.`

### Unwanted Behaviour Requirements (Agit Errors)

- **CI-03-301:** `If the user pushing via agit flow does not have write permission on the repository, then the system shall reject the push.`
- **CI-03-302:** `If the specified base branch does not exist, then the system shall reject the agit flow push with an error.`
- **CI-03-303:** `If the topic branch name does not meet format requirements, then the system shall reject the agit flow push.`

---

## 4. Auto-merge

**User Story:** As a pull request author, I want to schedule my PR for automatic merging so that it merges as soon as all required checks and approvals are satisfied.

### Ubiquitous Requirements (Auto-merge Properties)

- **CI-04-001:** `The system shall support the following merge strategies for auto-merge: merge commit, squash merge, and rebase merge.`
- **CI-04-002:** `The system shall track the head SHA at the time auto-merge is scheduled to detect stale merge attempts.`
- **CI-04-003:** `The system shall process auto-merge requests through a queue to prevent race conditions on concurrent eligible PRs.`

### Event-Driven Requirements (Auto-merge Workflow)

- **CI-04-101:** `When a user schedules auto-merge on a pull request, the system shall record the request with the selected merge strategy and head SHA.`
- **CI-04-102:** `When all required status checks pass, required approvals are satisfied, and no merge conflicts exist, the system shall merge the pull request automatically.`
- **CI-04-103:** `When auto-merge is scheduled, the system shall post a comment on the pull request indicating the auto-merge is pending.`
- **CI-04-104:** `When a user cancels auto-merge, the system shall remove the scheduled merge and post a cancellation comment.`
- **CI-04-105:** `When the pull request head SHA changes after auto-merge is scheduled, the system shall retain the auto-merge request and re-evaluate conditions against the new SHA.`

### Unwanted Behaviour Requirements (Auto-merge Errors)

- **CI-04-301:** `If the pull request has merge conflicts at the time conditions are evaluated, then the system shall not merge and shall notify the author.`
- **CI-04-302:** `If a non-author or non-admin user attempts to schedule auto-merge, then the system shall deny the operation.`
- **CI-04-303:** `If the head SHA at merge time differs from the scheduled SHA and no new push has occurred, then the system shall abort the merge.`

---

## 5. Commit Status

**User Story:** As an external CI system, I want to report build and test results on commits so that developers can see CI status directly in the Gitea UI.

### Ubiquitous Requirements (Status Properties)

- **CI-05-001:** `The system shall support the following status states: pending, success, error, failure, and warning.`
- **CI-05-002:** `The system shall allow multiple status contexts per commit (e.g., ci/lint, ci/test, ci/build).`
- **CI-05-003:** `The system shall compute a combined status across all contexts for a given commit.`
- **CI-05-004:** `The system shall store a target URL and description with each status for linking to external CI details.`

### Event-Driven Requirements (Status Workflow)

- **CI-05-101:** `When an API request creates a status for a commit, the system shall record the status with context, state, description, and target URL.`
- **CI-05-102:** `When the combined status of all contexts for a commit is requested, the system shall return the aggregate state (failure if any context has failed).`
- **CI-05-103:** `When a status is created for a commit on a protected branch, the system shall evaluate the status against required status check rules.`

### Optional Feature Requirements (Status Integration)

- **CI-05-201:** `Where required status checks are configured on a protected branch, the system shall block pull request merging until all required contexts report success.`
- **CI-05-202:** `Where a commit status is linked to a pull request, the system shall display the status on the pull request UI.`

### Unwanted Behaviour Requirements (Status Errors)

- **CI-05-301:** `If a status creation request targets a non-existent commit SHA, then the system shall reject the request.`
- **CI-05-302:** `If an API request lacks permission to create a status on a repository, then the system shall reject the request with 403 Forbidden.`

---

## 6. Git Hooks

**User Story:** As a repository administrator, I want to configure server-side Git hooks so that custom policies are enforced on every push.

### Ubiquitous Requirements (Git Hook Properties)

- **CI-06-001:** `The system shall support the following server-side Git hooks: pre-receive, update, and post-receive.`
- **CI-06-002:** `The system shall store custom hook scripts in the repository .git/hooks directory.`
- **CI-06-003:** `The system shall allow hook script editing through both the web UI and the API.`

### Event-Driven Requirements (Git Hook Workflow)

- **CI-06-101:** `When a Git push is received, the system shall execute the pre-receive hook before accepting any refs.`
- **CI-06-102:** `When the pre-receive hook exits with a non-zero status, the system shall reject the entire push operation.`
- **CI-06-103:** `When a ref is updated, the system shall execute the update hook for each ref being modified.`
- **CI-06-104:** `When a push is accepted, the system shall execute the post-receive hook after all refs have been updated.`
- **CI-06-105:** `When an admin updates a Git hook via the API, the system shall persist the updated script immediately.`

### Optional Feature Requirements (Git Hook Security)

- **CI-06-201:** `Where Git hooks are enabled, the system shall allow repository admins to create and edit hook scripts.`
- **CI-06-202:** `Where the DISABLE_GIT_HOOKS setting is true (default), the system shall prevent creation and editing of Git hooks.`

### Unwanted Behaviour Requirements (Git Hook Errors)

- **CI-06-301:** `If a non-admin user attempts to modify a Git hook, then the system shall deny the operation.`
- **CI-06-302:** `If a hook script contains a syntax error, then the system shall accept the script but log an error when execution fails.`

---

## 7. Cron / Scheduled Tasks

**User Story:** As a system administrator, I want background tasks to run on schedule so that system maintenance is performed automatically.

### Ubiquitous Requirements (Cron Properties)

- **CI-07-001:** `The system shall provide a configurable cron framework where each task has independent ENABLED, SCHEDULE, and RUN_AT_START settings.`
- **CI-07-002:** `The system shall support standard cron syntax for scheduling each background task.`
- **CI-07-003:** `The system shall expose a cron task monitoring interface in the admin panel.`

### Event-Driven Requirements (Cron Task Execution)

- **CI-07-101:** `When the cron schedule for repository mirror sync fires (default every 10 minutes), the system shall sync all configured mirror repositories.`
- **CI-07-102:** `When the cron schedule for repository health checks fires (default daily), the system shall run health checks on all repositories.`
- **CI-07-103:** `When the cron schedule for repository stats update fires (default daily), the system shall recalculate repository statistics.`
- **CI-07-104:** `When the cron schedule for archive cleanup fires (default daily), the system shall remove expired repository archives.`
- **CI-07-105:** `When the cron schedule for zombie task cleanup fires (default every 5 minutes), the system shall stop tasks exceeding ZOMBIE_TASK_TIMEOUT.`
- **CI-07-106:** `When the cron schedule for endless task cleanup fires (default every 30 minutes), the system shall stop tasks exceeding ENDLESS_TASK_TIMEOUT.`
- **CI-07-107:** `When the cron schedule for abandoned job cleanup fires (default every 6 hours), the system shall cancel jobs exceeding ABANDONED_JOB_TIMEOUT.`
- **CI-07-108:** `When the cron schedule for scheduled workflows fires, the system shall trigger all workflow runs due at that time.`
- **CI-07-109:** `When an external user sync source (LDAP) is configured and its cron schedule fires (default daily), the system shall synchronize user data from the external source.`

### Optional Feature Requirements (Cron Configuration)

- **CI-07-201:** `Where RUN_AT_START is enabled for a cron task, the system shall execute that task once during system startup.`
- **CI-07-202:** `Where a custom schedule is provided for a cron task, the system shall override the default schedule with the custom one.`

### Unwanted Behaviour Requirements (Cron Errors)

- **CI-07-301:** `If a cron task execution fails, then the system shall log the error and continue scheduling subsequent executions.`
- **CI-07-302:** `If a cron task is disabled via configuration, then the system shall skip all scheduled executions for that task.`

---

## 8. Queue System

**User Story:** As a system operator, I want a reliable queue infrastructure so that background tasks (webhooks, indexing, mail, actions) are processed concurrently without data loss.

### Ubiquitous Requirements (Queue Properties)

- **CI-08-001:** `The system shall provide two queue types: simple FIFO queue and unique queue with key-based deduplication.`
- **CI-08-002:** `The system shall support the following queue backends: channel (in-memory), LevelDB (persistent), Redis (distributed), and dummy (no-op for testing).`
- **CI-08-003:** `The system shall use a worker pool model with dynamic sizing based on CPU count.`
- **CI-08-004:** `The system shall support batch processing of queue items with configurable batch length.`
- **CI-08-005:** `The system shall expose queue monitoring in the admin panel.`

### Event-Driven Requirements (Queue Processing)

- **CI-08-101:** `When a queue item is pushed to a unique queue, the system shall deduplicate by key and drop or replace existing items with the same key.`
- **CI-08-102:** `When a worker becomes available, the system shall dispatch the next queued item or batch to that worker.`
- **CI-08-103:** `When a queue item processing fails, the system shall retry the item according to the configured retry policy.`
- **CI-08-104:** `When the system receives a graceful shutdown signal, the system shall finish processing in-flight items before terminating.`
- **CI-08-105:** `When the admin panel is accessed, the system shall display queue statistics including queue length, worker count, and processing rate.`

### State-Driven Requirements (Queue Lifecycle)

- **CI-08-701:** `While the LevelDB backend is active, the system shall persist queue items to disk so that items survive process restarts.`
- **CI-08-702:** `While the Redis backend is active, the system shall distribute queue items across multiple Gitea instances for horizontal scaling.`
- **CI-08-703:** `While the queue length exceeds the configured maximum (default 100,000), the system shall apply backpressure to enqueue operations.`

### Optional Feature Requirements (Queue Configuration)

- **CI-08-201:** `Where the queue type is set to LevelDB, the system shall use DATADIR as the storage path for persistence.`
- **CI-08-202:** `Where the queue type is set to Redis, the system shall use CONN_STR as the connection string for the Redis instance.`
- **CI-08-203:** `Where MAX_WORKERS is configured, the system shall cap the worker pool at the specified count (default CPU count / 2, maximum 10).`

### Unwanted Behaviour Requirements (Queue Errors)

- **CI-08-301:** `If a queue backend becomes unavailable, then the system shall log the error and queue new items in memory until the backend recovers.`
- **CI-08-302:** `If a worker panics during item processing, then the system shall recover the worker and requeue the item.`
- **CI-08-303:** `If a LevelDB queue data directory is corrupted, then the system shall log the error and start with a fresh queue.`

---

## Business Rules

- **BR-05-001:** Gitea Actions requires at least one registered runner to execute workflows
- **BR-05-002:** Workflow YAML files must reside in .gitea/workflows/ or .github/workflows/ at the repository root
- **BR-05-003:** Fork pull requests have restricted access to repository secrets
- **BR-05-004:** Artifact retention defaults to 90 days and is configurable per instance
- **BR-05-005:** Webhook signatures use HMAC-SHA256 for payload integrity verification
- **BR-05-006:** Webhook host allowlist is enforced for all outgoing deliveries when configured
- **BR-05-007:** Agit flow requires write permission on the target repository
- **BR-05-008:** Auto-merge can only be scheduled by the PR author or a repository admin
- **BR-05-009:** Auto-merge tracks the head SHA to prevent merging stale code
- **BR-05-010:** Commit status contexts are independent; combined status is failure if any context fails
- **BR-05-011:** Git hooks are disabled by default (DISABLE_GIT_HOOKS=true) for security
- **BR-05-012:** Cron tasks have independent schedules and can be enabled or disabled individually
- **BR-05-013:** Queue worker pool size defaults to CPU count / 2 with a maximum of 10 workers
- **BR-05-014:** LevelDB is the default queue backend for production; channel is for single-instance development

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Workflow trigger with no registered runner | Job remains pending until a matching runner registers |
| Fork PR accesses protected secrets | Restricted access; secrets masked or unavailable |
| Artifact exceeds retention period | Automatically deleted by cron cleanup task |
| Zombie task stuck beyond timeout | Stopped by cron zombie task cleanup (default 10 minutes) |
| Webhook delivery to disallowed host | Rejected with security error |
| Webhook queue exceeds max length | Backpressure applied to prevent unbounded growth |
| Agit push to non-existent base branch | Push rejected with error message |
| Auto-merge on PR with merge conflicts | Merge blocked; author notified |
| Auto-merge scheduled by non-author/non-admin | Operation denied |
| Commit status for non-existent SHA | Request rejected |
| Git hook script with syntax error | Script saved but error logged on execution failure |
| Git hooks modified by non-admin | Operation denied |
| Cron task execution failure | Error logged; next scheduled execution proceeds normally |
| Queue backend unavailable | Items queued in memory; error logged |
| Queue worker panic during processing | Worker recovered; item requeued |
| LevelDB queue data corruption | Error logged; fresh queue started |
| Abandoned job exceeds timeout | Cancelled by cron abandoned job cleanup (default 24 hours) |

## Success Criteria

- Workflow runs begin execution within 5 seconds of trigger when a runner is available
- Webhook deliveries complete within the configured timeout (default 5 seconds)
- Webhook retry does not block the delivery queue for other hooks
- Artifact upload and download throughput meets storage backend performance limits
- Queue item processing latency is under 100ms per item for lightweight tasks (webhooks, notifications)
- Auto-merge executes within 30 seconds of all conditions being satisfied
- Cron task schedules fire within 1 second of the configured cron time
- Commit status API responds within 200ms for creation and lookup
- Combined status computation completes in under 50ms for commits with up to 100 contexts
- Runner registration completes in under 1 second with valid token
- Queue graceful shutdown completes all in-flight items without data loss
