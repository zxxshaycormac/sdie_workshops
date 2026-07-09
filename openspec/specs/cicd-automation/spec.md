# 05 — CI/CD & Automation

Baseline specification of Gitea's CI/CD, webhook, and server-side automation subsystems. All requirements describe the current (v1.22.x) system behavior. CI-produced artifacts and logs live in this domain; published packages and releases are documented in Domain 06 (Packages & Releases). Cron and queue infrastructure is documented in Domain 08 (Administration & Operations).

---

## 1. Gitea Actions

**User Story:** As a developer, I want to define automated CI/CD workflows in YAML so that my code is built, tested, and deployed automatically on every push or pull request.

### Ubiquitous Requirements (Workflow Properties)

- **CICD-01-001:** `The system shall parse workflow definitions from YAML files stored in the .gitea/workflows/ or .github/workflows/ directory of a repository.`
- **CICD-01-002:** `The system shall support the following workflow trigger events: push, pull_request, pull_request_target, schedule, create, delete, fork, issues, issue_comment, release, pull_request_review, pull_request_review_comment, registry_package, and gollum. The workflow_dispatch and repository_dispatch events are NOT supported (ignored by Gitea Actions).`
- **CICD-01-003:** `The system shall track workflow run status through the lifecycle: unknown, waiting (pending), running, success, failure, cancelled, skipped, and blocked.`
- **CICD-01-004:** `The system shall support parallel job execution within a single workflow run.`
- **CICD-01-005:** `The system shall support matrix build strategies that fan out jobs across defined parameter combinations.`
- **CICD-01-006:** `The system shall store workflow run logs with persistent storage and indexed access for retrieval.`
- **CICD-01-007:** `The system shall generate a unique task token for each job execution distributed to a runner.`
- **CICD-01-008:** `The system shall require the actions unit (unit.TypeActions) to be enabled on the repository before processing workflows.`

### Event-Driven Requirements (Workflow Execution)

- **CICD-01-101:** `When a trigger event matches a workflow's on clause, the system shall create a new workflow run.`
- **CICD-01-102:** `When a workflow run is created, the system shall evaluate job dependencies and schedule jobs for execution.`
- **CICD-01-103:** `When a job is assigned to a runner, the system shall execute steps sequentially within that job.`
- **CICD-01-104:** `When a workflow run completes (all jobs in a terminal state), the system shall set the run status to the most severe job outcome with precedence: failure if any job failed, otherwise cancelled if any job was cancelled, otherwise success if any job succeeded, otherwise skipped.`
- **CICD-01-105:** `When a user cancels a workflow run via the API, the system shall cancel all running and pending jobs in that run.`
- **CICD-01-106:** `When a user reruns a workflow via the API, the system shall create a new run with the same workflow definition and trigger event.`
- **CICD-01-107:** `When a workflow step references an action (uses:), the system shall resolve the action from the configured DEFAULT_ACTIONS_URL (GitHub or self-hosted).`

### State-Driven Requirements (Runner Lifecycle)

- **CICD-01-701:** `While a runner is in active state, the system shall assign matching jobs to that runner based on label selection.`
- **CICD-01-702:** `While a runner is in offline state, the system shall not assign new jobs to that runner.`
- **CICD-01-703:** `While a task exceeds the configured ZOMBIE_TASK_TIMEOUT (default 10 minutes), the system shall mark the task as timed out.`
- **CICD-01-704:** `While a task exceeds the configured ENDLESS_TASK_TIMEOUT (default 3 hours), the system shall stop the task.`

### State-Driven Requirements (Run Status)

- **CICD-01-705:** `While a workflow run has any non-terminal job, the system shall set the run status with precedence: running if any job is running, otherwise blocked if any job is blocked (awaiting environment approval or unmet job needs), otherwise waiting.`

### Optional Feature Requirements (Actions Configuration)

- **CICD-01-201:** `Where the actions feature is enabled via configuration, the system shall expose workflow, run, job, and artifact endpoints.`
- **CICD-01-202:** `Where artifact storage is configured (local or MinIO/S3), the system shall store and retrieve workflow artifacts using the configured backend.`
- **CICD-01-203:** `Where scheduled workflows are defined with cron syntax, the system shall trigger workflow runs according to the specified schedule.`
- **CICD-01-204:** `Where a status badge endpoint is requested for a workflow, the system shall return an SVG image reflecting the latest run status.`
- **CICD-01-205:** `Where the actions feature is disabled via configuration, the system shall hide all actions UI and API endpoints.`
- **CICD-01-206:** `Where a workflow uses the pull_request_target trigger, the system shall run the workflow in the context of the base branch (not the PR head) so that fork pull requests cannot execute untrusted code with write secrets.`
- **CICD-01-207:** `Where a workflow run is triggered from a fork pull request (non-pull_request_target), the system shall require manual maintainer approval before scheduling jobs unless the triggering user has write Actions permission or was previously approved.`
- **CICD-01-208:** `Where a commit message or PR title contains one of the configured SKIP_WORKFLOW_STRINGS (default [skip ci], [ci skip], [no ci], [skip actions], [actions skip]), the system shall skip the workflow run for push, pull_request, and pull_request_sync events.`
- **CICD-01-209:** `The system shall persist per-task outputs via ActionTaskOutput (keyed by TaskID + OutputKey) and a per-scope ActionTasksVersion monotonic counter (global/org/repo) so reruns reset outputs and runners can detect new tasks.`
- **CICD-01-210:** `The system shall store workflow run logs in a dedicated [actions_log] storage backend (distinct from [actions.artifacts]), resolved via getStorage(rootCfg, "actions_log", "", nil).`

### Unwanted Behaviour Requirements (Actions Errors)

- **CICD-01-301:** `If a fork pull request workflow attempts to access repository secrets, then the system shall restrict access to protected secrets.`
- **CICD-01-302:** `If no runner is available with matching labels for a job, then the system shall keep the job in pending state until a matching runner registers.`
- **CICD-01-303:** `If an artifact exceeds the configured retention period (default 90 days), then the system shall delete the artifact.`
- **CICD-01-304:** `If a job remains abandoned beyond the configured ABANDONED_JOB_TIMEOUT (default 24 hours), then the system shall cancel the job.`
- **CICD-01-305:** `If the actions unit is disabled for a repository, then the system shall reject workflow trigger events for that repository.`
- **CICD-01-306:** `If a fork pull request run has not been approved and the triggering user lacks write Actions permission and has never been approved before, then the system shall hold the run in blocked (NeedApproval) state pending manual approval.`

### Complex Requirements (Runner Registration)

- **CICD-01-901:** `When a runner registration request is received with a valid registration token, the system shall create a runner record with a unique UUID, assign the specified labels, and associate the runner with the appropriate scope (global, organization, or repository).`

---

## 2. Actions Secrets

**User Story:** As a repository administrator, I want to store encrypted secrets (deployment tokens, signing keys) at the repository, organization, or user level so that workflows can use them without exposing the values in logs or code.

### Ubiquitous Requirements (Secret Properties)

- **CICD-02-001:** `The system shall support Actions secrets at three scopes: repository, organization, and user.`
- **CICD-02-002:** `The system shall encrypt secret values at rest using the instance secret key.`
- **CICD-02-003:** `The system shall expose secret values to workflow jobs only via the secrets runtime context.`
- **CICD-02-004:** `The system shall mask secret values in workflow run logs when the value appears in command output.`

### Event-Driven Requirements (Secret Workflow)

- **CICD-02-101:** `When an authorized user submits a new secret at repository, organization, or user scope, the system shall encrypt the value and store it under the chosen name.`
- **CICD-02-102:** `When a workflow job references a secret by name, the system shall inject the resolved secret into the job environment only after resolving the precedence: repository, then organization, then user.`
- **CICD-02-103:** `When an authorized user updates an existing secret, the system shall overwrite the encrypted value without retaining the previous value.`
- **CICD-02-104:** `When an authorized user deletes a secret, the system shall remove the encrypted record immediately.`
- **CICD-02-105:** `When a workflow run completes, the system shall not persist secret values in the run log storage.`

### Optional Feature Requirements (Secret Scope Inheritance)

- **CICD-02-201:** `Where an organization-level secret has the same name as a repository-level secret, the system shall use the repository-level secret for workflow runs in that repository.`
- **CICD-02-202:** `Where a secret is configured at user scope, the system shall make the secret available to workflows running in repositories owned by that user.`

### Unwanted Behaviour Requirements (Secret Errors)

- **CICD-02-301:** `If a workflow run is triggered from a fork pull request, then the system shall restrict access to repository and organization secrets unless the maintainer has explicitly approved the run.`
- **CICD-02-302:** `If a secret name contains characters other than alphanumeric or underscore, then the system shall reject the secret creation with a validation error.`
- **CICD-02-303:** `If an unauthorized user attempts to view, create, or delete a secret, then the system shall deny the operation with 403 Forbidden.`
- **CICD-02-304:** `If a workflow references a secret name that is not defined at any scope, then the system shall substitute an empty string for the secret value and continue the run.`
- **CICD-02-305:** `If a secret name begins with the reserved prefixes GITEA_ or GITHUB_ (case-insensitive), then the system shall reject the secret creation with a validation error to prevent collision with built-in environment variables.`

---

## 3. Actions Variables

**User Story:** As a repository administrator, I want to store non-sensitive configuration values (build flags, image tags) at the repository, organization, or user level so that workflows can reference them without duplicating the values across files.

### Ubiquitous Requirements (Variable Properties)

- **CICD-03-001:** `The system shall support Actions variables at three scopes: repository, organization, and user.`
- **CICD-03-002:** `The system shall store variable values in plaintext.`
- **CICD-03-003:** `The system shall expose variable values to workflow jobs and workflow-file parsing via the vars runtime context.`
- **CICD-03-004:** `The system shall display variable values unmasked in the management UI because variables are not intended to hold secrets.`

### Event-Driven Requirements (Variable Workflow)

- **CICD-03-101:** `When an authorized user creates a variable, the system shall store the value and make it resolvable from workflow YAML via the vars context.`
- **CICD-03-102:** `When a workflow run parses YAML, the system shall substitute ${{ vars.NAME }} expressions using the resolved variable value.`
- **CICD-03-103:** `When an authorized user updates a variable, the system shall apply the change to all subsequent workflow runs without affecting in-flight runs.`
- **CICD-03-104:** `When resolving a variable name that exists at multiple scopes, the system shall apply the precedence: repository, then organization, then user.`

### Optional Feature Requirements (Variable Scope Inheritance)

- **CICD-03-201:** `Where a repository-level variable shadows an organization-level variable of the same name, the system shall use the repository-level value for runs in that repository.`
- **CICD-03-202:** `Where a user-level variable is defined, the system shall make the variable available to workflow runs in repositories owned by that user.`

### Unwanted Behaviour Requirements (Variable Errors)

- **CICD-03-301:** `If a variable name contains characters other than alphanumeric or underscore, then the system shall reject the variable creation with a validation error.`
- **CICD-03-302:** `If a workflow references an undefined variable, then the system shall substitute an empty string and continue the run.`
- **CICD-03-303:** `If an unauthorized user attempts to view, create, or delete a variable, then the system shall deny the operation with 403 Forbidden.`
- **CICD-03-304:** `If a variable name begins with the reserved prefix CI (case-insensitive), then the system shall reject the variable creation with a validation error to prevent collision with the built-in CI runtime variable.`

---

## 4. Actions Artifacts

**User Story:** As a workflow author, I want to upload artifacts (build outputs, test reports) during a workflow run and download them later so that I can share build outputs with consumers.

> **Scope note:** This section covers CI-produced artifacts generated by workflow runs. Published packages and releases are documented in Domain 06 (Packages & Releases).

### Ubiquitous Requirements (Artifact Properties)

- **CICD-04-001:** `The system shall store artifacts scoped to a specific workflow run.`
- **CICD-04-002:** `The system shall enforce a configurable retention period (default 90 days) after which artifacts are deleted.`
- **CICD-04-003:** `The system shall provide UI and API endpoints for listing, downloading, and deleting artifacts.`
- **CICD-04-004:** `The system shall store artifacts in the configured actions_artifact storage backend.`

### Event-Driven Requirements (Artifact Workflow)

- **CICD-04-101:** `When a workflow job calls the artifact upload action, the system shall persist the artifact under a unique name scoped to the run.`
- **CICD-04-102:** `When a workflow job calls the artifact download action, the system shall serve the named artifact from the same run or a previously-completed run.`
- **CICD-04-103:** `When a user navigates to a workflow run's artifacts tab, the system shall list every artifact uploaded by the run.`
- **CICD-04-104:** `When a user clicks an artifact in the UI, the system shall stream a zip archive of the artifact contents.`
- **CICD-04-105:** `When an authorized user deletes an artifact via the UI or API, the system shall remove the artifact immediately.`

### State-Driven Requirements (Artifact Retention)

- **CICD-04-701:** `While an artifact has not reached the retention limit, the system shall keep it accessible for download.`
- **CICD-04-702:** `While the configured storage backend is unreachable, the system shall fail artifact upload/download operations with a storage error and not corrupt existing artifacts.`

### Optional Feature Requirements (Artifact Configuration)

- **CICD-04-201:** `Where the retention period is overridden per-instance, the system shall apply the configured value in place of the 90-day default.`
- **CICD-04-202:** `Where SERVE_DIRECT is enabled on the storage backend, the system shall issue pre-signed direct-download URLs that bypass the Gitea server.`

### Unwanted Behaviour Requirements (Artifact Errors)

- **CICD-04-301:** `If an artifact upload exceeds the configured maximum size, then the system shall reject the upload with a size-limit error.`
- **CICD-04-302:** `If a user attempts to download an artifact from a private repository without authentication, then the system shall return 404.`
- **CICD-04-303:** `If an artifact name contains path separators or invalid characters, then the system shall reject the upload.`

---

## 5. Actions Log Storage

**User Story:** As a developer, I want to access workflow run logs so that I can debug failing CI jobs.

### Ubiquitous Requirements (Log Storage Properties)

- **CICD-LOG-001:** `The system shall store workflow run logs in a dedicated [actions_log] storage backend, distinct from the [actions.artifacts] artifact storage.`
- **CICD-LOG-002:** `The system shall resolve the log storage backend via getStorage(rootCfg, "actions_log", "", nil).`
- **CICD-LOG-003:** `The system shall NOT consult the legacy [actions] section for log storage configuration.`
- **CICD-LOG-004:** `The system shall support the same storage keys as artifact storage (STORAGE_TYPE, PATH, SERVE_DIRECT, MINIO_*).`

### Event-Driven Requirements (Log Access)

- **CICD-LOG-101:** `When a workflow step produces output, the system shall persist the log lines to the configured log storage backend.`
- **CICD-LOG-102:** `When a user requests a workflow run log, the system shall stream the stored log content from the log storage backend.`

### Unwanted Behaviour Requirements (Log Storage Errors)

- **CICD-LOG-301:** `If the log storage backend is unreachable, then the system shall fail log retrieval and not serve corrupted content.`
- **CICD-LOG-302:** `If a workflow run is cleaned up by the cron cleanup_actions task (see Domain 08), then the system shall remove the associated logs along with the run record.`

---

## 6. Webhooks

**User Story:** As a repository administrator, I want to configure webhook callbacks so that external systems are notified when repository events occur.

### Ubiquitous Requirements (Webhook Properties)

- **CICD-05-001:** `The system shall support the following webhook payload formats: Gitea (native), Gogs, Slack, Discord, Dingtalk, Telegram, Microsoft Teams, Feishu (Lark), Matrix, WeChat Work, and Packagist.`
- **CICD-05-002:** `The system shall track delivery status for each webhook invocation with success or failure indication.`
- **CICD-05-003:** `The system shall store both request and response data for each webhook delivery.`
- **CICD-05-004:** `The system shall support webhook configuration at three scopes: system (admin), organization, and repository.`
- **CICD-05-005:** `The system shall support proxy configuration for webhook delivery in restricted network environments.`

### Event-Driven Requirements (Webhook Delivery)

- **CICD-05-101:** `When a subscribed repository event occurs, the system shall enqueue an asynchronous webhook delivery to each configured hook endpoint.`
- **CICD-05-102:** `When a webhook delivery is triggered, the system shall compute an HMAC-SHA256 signature of the payload and include it in the request headers.`
- **CICD-05-103:** `When a webhook delivery fails, the system shall record the failure and allow manual redelivery.`
- **CICD-05-104:** `When a user triggers a test delivery via the UI or API, the system shall send a test payload to the webhook endpoint.`
- **CICD-05-105:** `When a user requests redelivery of a past webhook, the system shall resend the original payload to the endpoint.`
- **CICD-05-106:** `When a webhook delivery times out, the system shall record the timeout as a delivery failure.`
- **CICD-05-107:** `When the following repository events occur, the system shall trigger webhook delivery: push, create, delete, fork, issues, issue assign, issue label, issue milestone, issue comment, pull request, pull request assign, pull request label, pull request milestone, pull request comment, pull request review (approved/rejected/comment), pull request sync, pull request review request, wiki, repository, release, package, and schedule.`

### Optional Feature Requirements (Webhook Security)

- **CICD-05-201:** `Where a webhook host allowlist is configured, the system shall reject webhook deliveries to hosts not on the allowlist.`
- **CICD-05-202:** `Where TLS verification is skipped via configuration, the system shall deliver webhooks without validating the server certificate.`
- **CICD-05-203:** `Where a proxy URL is configured, the system shall route webhook deliveries through the specified proxy.`

### Unwanted Behaviour Requirements (Webhook Errors)

- **CICD-05-301:** `If a webhook delivery target is not in the allowed host list, then the system shall reject the delivery with a security error.`
- **CICD-05-302:** `If the webhook queue exceeds the configured QUEUE_LENGTH (default 1000), then the system shall apply backpressure to prevent unbounded memory growth.`
- **CICD-05-303:** `If a webhook endpoint returns an HTTP error status, then the system shall record the delivery as failed with the response details.`

---

## 7. Git Hooks

**User Story:** As a repository administrator, I want to configure server-side Git hooks so that custom policies are enforced on every push.

### Ubiquitous Requirements (Git Hook Properties)

- **CICD-06-001:** `The system shall support the following server-side Git hooks: pre-receive, update, and post-receive.`
- **CICD-06-002:** `The system shall store custom hook scripts in the repository .git/hooks directory.`
- **CICD-06-003:** `The system shall allow hook script editing through both the web UI and the API.`

### Event-Driven Requirements (Git Hook Workflow)

- **CICD-06-101:** `When a Git push is received, the system shall execute the pre-receive hook before accepting any refs.`
- **CICD-06-102:** `When the pre-receive hook exits with a non-zero status, the system shall reject the entire push operation.`
- **CICD-06-103:** `When a ref is updated, the system shall execute the update hook for each ref being modified.`
- **CICD-06-104:** `When a push is accepted, the system shall execute the post-receive hook after all refs have been updated.`
- **CICD-06-105:** `When an admin updates a Git hook via the API, the system shall persist the updated script immediately.`

### Optional Feature Requirements (Git Hook Security)

- **CICD-06-201:** `Where Git hooks are enabled, the system shall allow repository admins to create and edit hook scripts.`
- **CICD-06-202:** `Where the DISABLE_GIT_HOOKS setting is true (default), the system shall prevent creation and editing of Git hooks.`

### Unwanted Behaviour Requirements (Git Hook Errors)

- **CICD-06-301:** `If a non-admin user attempts to modify a Git hook, then the system shall deny the operation.`
- **CICD-06-302:** `If a hook script contains a syntax error, then the system shall accept the script but log an error when execution fails.`

---

## 8. Commit Status

**User Story:** As an external CI system, I want to report build and test results on commits so that developers can see CI status directly in the Gitea UI.

### Ubiquitous Requirements (Status Properties)

- **CICD-07-001:** `The system shall support the following status states: pending, success, error, failure, and warning.`
- **CICD-07-002:** `The system shall allow multiple status contexts per commit (e.g., ci/lint, ci/test, ci/build).`
- **CICD-07-003:** `The system shall compute a combined status across all contexts for a given commit.`
- **CICD-07-004:** `The system shall store a target URL and description with each status for linking to external CI details.`

### Event-Driven Requirements (Status Workflow)

- **CICD-07-101:** `When an API request creates a status for a commit, the system shall record the status with context, state, description, and target URL.`
- **CICD-07-102:** `When the combined status of all contexts for a commit is requested, the system shall return the aggregate state (failure if any context has failed).`
- **CICD-07-103:** `When a status is created for a commit on a protected branch, the system shall evaluate the status against required status check rules.`

### Optional Feature Requirements (Status Integration)

- **CICD-07-201:** `Where required status checks are configured on a protected branch, the system shall block pull request merging until all required contexts report success.`
- **CICD-07-202:** `Where a commit status is linked to a pull request, the system shall display the status on the pull request UI.`

### Unwanted Behaviour Requirements (Status Errors)

- **CICD-07-301:** `If a status creation request targets a non-existent commit SHA, then the system shall reject the request.`
- **CICD-07-302:** `If an API request lacks permission to create a status on a repository, then the system shall reject the request with 403 Forbidden.`

---

## Configuration Reference

The following INI sections configure CI/CD and automation behaviors. Cron and queue configuration → Domain 08 (Administration & Operations).

### [actions] Section

- **ENABLED**: Master switch for Gitea Actions (default true).
- **DEFAULT_ACTIONS_URL**: Default source for action resolution (`github` or `self`).
- **ARTIFACT_RETENTION_DAYS**: Days to retain workflow artifacts (default 90).
- **ZOMBIE_TASK_TIMEOUT**: Idle-then-dead task timeout (default 10m).
- **ENDLESS_TASK_TIMEOUT**: Maximum task runtime (default 3h).
- **ABANDONED_JOB_TIMEOUT**: Idle job cleanup threshold (default 24h).
- **SKIP_WORKFLOW_STRINGS**: Commit-message/PR-title substrings that skip workflow runs: `[skip ci]`, `[ci skip]`, `[no ci]`, `[skip actions]`, `[actions skip]` (applies to push and pull_request/_sync events only).

### [actions.artifacts] Section

- **STORAGE_TYPE**: Artifact storage backend (`local`, `minio`).
- **PATH**: Local filesystem path for artifact storage.
- **SERVE_DIRECT**: Use direct URLs for MinIO backend.
- **MINIO_ENDPOINT**, **MINIO_ACCESS_KEY_ID**, **MINIO_SECRET_ACCESS_KEY**, **MINIO_BUCKET**, **MINIO_LOCATION**: MinIO connection parameters.

### [actions_log] Section

Workflow run logs use a dedicated storage type, configured under `[actions_log]` (resolved via `getStorage(rootCfg, "actions_log", "", nil)`). It is distinct from `[actions.artifacts]` and supports the same storage keys (`STORAGE_TYPE`, `PATH`, `SERVE_DIRECT`, `MINIO_*`). The legacy `[actions]` section is NOT consulted for log storage.

### [webhook] Section

- **QUEUE_LENGTH**: Max queued webhook deliveries before backpressure (default 1000).
- **DELIVER_TIMEOUT**: Per-delivery HTTP timeout in seconds (default 5).
- **SKIP_TLS_VERIFY**: Skip TLS certificate verification on delivery (default false).
- **ALLOWED_HOST_LIST**: Comma-separated allowed destination host list (SSRF protection).
- **PAGING_NUM**: Paging size for webhook list UI/API (default 10).
- **PROXY_URL**: Outbound proxy URL for webhook delivery.
- **PROXY_HOSTS**: Comma-separated host patterns routed through the proxy.

---

## Business Rules

- **BR-05-001:** Gitea Actions requires at least one registered runner to execute workflows
- **BR-05-002:** Workflow YAML files must reside in .gitea/workflows/ or .github/workflows/ at the repository root
- **BR-05-003:** Fork pull requests have restricted access to repository secrets
- **BR-05-004:** Artifact retention defaults to 90 days and is configurable per instance
- **BR-05-005:** Webhook signatures use HMAC-SHA256 for payload integrity verification
- **BR-05-006:** Webhook host allowlist is enforced for all outgoing deliveries when configured
- **BR-05-007:** Commit status contexts are independent; combined status is failure if any context fails
- **BR-05-008:** Git hooks are disabled by default (DISABLE_GIT_HOOKS=true) for security

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| Workflow trigger with no registered runner | Job remains pending until a matching runner registers |
| Fork PR accesses protected secrets | Restricted access; secrets masked or unavailable |
| Artifact exceeds retention period | Automatically deleted by cron cleanup task (see Domain 08) |
| Zombie task stuck beyond timeout | Stopped by cron zombie task cleanup (default 10 minutes; see Domain 08) |
| Webhook delivery to disallowed host | Rejected with security error |
| Webhook queue exceeds max length | Backpressure applied to prevent unbounded growth |
| Commit status for non-existent SHA | Request rejected |
| Git hook script with syntax error | Script saved but error logged on execution failure |
| Git hooks modified by non-admin | Operation denied |
| Abandoned job exceeds timeout | Cancelled by cron abandoned job cleanup (default 24 hours; see Domain 08) |
| All jobs skipped (e.g. all if: conditions false) | Run status is skipped (not success) |
| A job is cancelled and no job failed | Run status is cancelled (not failure) |
| A failed job caused dependent jobs to be cancelled | Run status is failure (failure dominates cancelled) |
| All jobs blocked awaiting approval | Run status is blocked (not running) |
| Mix of blocked and waiting jobs | Run status is blocked |
| Mix of running and blocked jobs | Run status is running |
| Job in unknown (zero-value) non-terminal state | Run falls back to running |
| Empty job set (degenerate) | Run status is skipped |

## Success Criteria

- Workflow runs begin execution within 5 seconds of trigger when a runner is available
- Webhook deliveries complete within the configured timeout (default 5 seconds)
- Webhook retry does not block the delivery queue for other hooks
- Artifact upload and download throughput meets storage backend performance limits
- Commit status API responds within 200ms for creation and lookup
- Combined status computation completes in under 50ms for commits with up to 100 contexts
- Runner registration completes in under 1 second with valid token
