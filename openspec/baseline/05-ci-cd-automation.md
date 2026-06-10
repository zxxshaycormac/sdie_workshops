# 05 — CI/CD & Automation

## 1. Gitea Actions

### What
GitHub Actions-compatible workflow automation system. YAML workflows in `.github/workflows/` with runners, jobs, steps, artifacts, secrets, and variables.

### Behaviors
- **Workflows**: YAML definition, trigger on push/PR/schedule/manual
- **Runs**: workflow execution tracking with status (pending/running/success/failure/cancelled)
- **Jobs**: parallel work units, matrix builds supported
- **Tasks**: individual job executions distributed to runners with unique tokens
- **Steps**: sequential commands within a job (shell or action references)
- **Runners**:
  - Registration via tokens with UUID
  - Scopes: global (admin), organization, repository
  - Labels for job matching
  - Status: active, idle, offline
  - Native, Docker, or VM-based
- **Artifacts**:
  - Upload/download during workflow
  - Configurable retention (default 90 days)
  - Storage backend: local, MinIO/S3
- **Variables**: repo/org-level environment variables (non-secret)
- **Secrets**: encrypted repo/org-level secrets (separate from regular repo secrets)
- **Scheduled Workflows**: cron syntax (GitHub Actions compatible)
- **Status Badges**: SVG endpoint for CI status
- **Logs**: persistent storage with indexed access

### UI Routes
- GET /{owner}/{repo}/actions — Workflow/run list
- GET /{owner}/{repo}/actions/runs/{index} — Run detail
- GET /{owner}/{repo}/actions/runs/{index}/jobs/{job} — Job detail
- GET /{owner}/{repo}/actions/runs/{index}/artifacts — Artifact list
- GET /{owner}/{repo}/actions/runners — Runner management
- GET /{owner}/{repo}/settings/actions/secrets — Secrets
- GET /{owner}/{repo}/settings/actions/variables — Variables

### API Endpoints
- GET /repos/{owner}/{repo}/actions/runs — List runs
- GET /repos/{owner}/{repo}/actions/runs/{id} — Get run
- GET /repos/{owner}/{repo}/actions/runs/{id}/jobs — List jobs
- POST /repos/{owner}/{repo}/actions/runs/{id}/rerun — Rerun
- POST /repos/{owner}/{repo}/actions/runs/{id}/cancel — Cancel
- GET /repos/{owner}/{repo}/actions/artifacts — List artifacts
- GET /repos/{owner}/{repo}/actions/artifacts/{id} — Get artifact
- DELETE /repos/{owner}/{repo}/actions/artifacts/{id} — Delete artifact
- POST /repos/{owner}/{repo}/actions/registrations/runner — Register runner
- DELETE /repos/{owner}/{repo}/actions/runners/{id} — Remove runner
- GET /repos/{owner}/{repo}/actions/secrets — List secrets
- PUT /repos/{owner}/{repo}/actions/secrets/{name} — Create/update secret
- DELETE /repos/{owner}/{repo}/actions/secrets/{name} — Delete secret
- GET /repos/{owner}/{repo}/actions/variables — List variables
- POST /repos/{owner}/{repo}/actions/variables — Create variable
- PATCH /repos/{owner}/{repo}/actions/variables/{name} — Update variable
- DELETE /repos/{owner}/{repo}/actions/variables/{name} — Delete variable

### Config
- [actions] ENABLED — enable/disable
- [actions] DEFAULT_ACTIONS_URL — GitHub or self-hosted action source
- [actions] ARTIFACT_RETENTION_DAYS — retention (default 90)
- [actions] LOG_STORAGE — log storage backend
- [actions] ARTIFACT_STORAGE — artifact storage backend
- [actions] ZOMBIE_TASK_TIMEOUT — stuck task timeout (default 10m)
- [actions] ENDLESS_TASK_TIMEOUT — endless task timeout (default 3h)
- [actions] ABANDONED_JOB_TIMEOUT — abandoned job timeout (default 24h)

### Constraints
- Requires unit.TypeActions enabled on repo
- Fork PRs have restricted permissions for secrets
- Requires at least one registered runner

---

## 2. Webhooks

### What
HTTP callbacks triggered by repository events. Multiple integration formats supported.

### Webhook Types
- Gitea (native format)
- Gogs (compatible)
- Slack
- Discord
- Dingtalk
- Telegram
- Microsoft Teams
- Feishu (Lark)
- Matrix
- WeChat Work
- Packagist
- General (custom payload)

### Webhook Events
- Push (branch/tag push)
- Create (branch/tag creation)
- Delete (branch/tag deletion)
- Release (publish/update/delete)
- Pull request (open/close/merge/sync)
- Pull request review (submitted)
- Pull request review comment
- Issues (open/close/edit/reopen)
- Issue comment (create/delete/edit)
- Issue label (attach/detach)
- Issue milestone (attach/detach)
- Issue assign (assign/unassign)
- Repository (create/fork/delete)
- Package (create/delete)

### Behaviors
- Asynchronous delivery with retry queue
- Configurable delivery timeout
- Request/response logging
- Status tracking (success/failure)
- HMAC-SHA256 signature verification
- Proxy support for restricted networks
- Host allowlist for security
- Test delivery endpoint

### UI Routes
- GET /{owner}/{repo}/settings/hooks — List webhooks
- GET /{owner}/{repo}/settings/hooks/{type}/new — Create webhook
- GET /{owner}/{repo}/settings/hooks/{id} — Edit webhook
- POST /{owner}/{repo}/settings/hooks/{id}/test — Test delivery
- GET /{owner}/{repo}/settings/hooks/{id}/deliveries — Delivery log
- POST /{owner}/{repo}/settings/hooks/{id}/redeliver — Redeliver
- System webhooks: GET /admin/hooks
- Org webhooks: GET /:org/settings/hooks

### API Endpoints
- GET /repos/{owner}/{repo}/hooks — List webhooks
- POST /repos/{owner}/{repo}/hooks — Create webhook
- GET /repos/{owner}/{repo}/hooks/{id} — Get webhook
- PATCH /repos/{owner}/{repo}/hooks/{id} — Update webhook
- DELETE /repos/{owner}/{repo}/hooks/{id} — Delete webhook
- POST /repos/{owner}/{repo}/hooks/{id}/tests — Test webhook

### Config
- [webhook] QUEUE_LENGTH — pending deliveries (default 1000)
- [webhook] DELIVER_TIMEOUT — delivery timeout (default 5s)
- [webhook] SKIP_TLS_VERIFY — skip TLS verification
- [webhook] ALLOWED_HOST_LIST — host allowlist
- [webhook] PAGING_NUM — webhooks per page (default 10)
- [webhook] PROXY_URL — proxy for deliveries

---

## 3. Agit Flow

### What
Server-side Git flow for creating PRs via push. Uses `refs/for/<base>/<topic>` ref format.

### Behaviors
- Push to `refs/for/<base-branch>/<topic-branch>` creates a PR automatically
- Topic branches get user prefix (user/topic-branch)
- Auto-creates pull request for topic branch
- Supports force push configuration
- Creates branch and PR in single push operation

### Constraints
- Requires repository write permission
- Base branch must exist
- Topic branch name format requirements
- No dedicated UI or API (server-side feature)

---

## 4. Auto-merge

### What
Automatic merging of PRs when required checks and approvals are satisfied.

### Behaviors
- Schedule auto-merge on a PR
- Conditions: passing status checks, required approvals, no conflicts
- Merge strategies: merge, squash, rebase
- SHA tracking to prevent stale merges
- Queue-based processing to prevent race conditions
- Auto-comment on scheduling/cancellation

### API Endpoints
- POST /repos/{owner}/{repo}/pulls/{index}/merge (with MergeWhenChecksSucceed=true) — Schedule auto-merge
- DELETE /repos/{owner}/{repo}/pulls/{index}/merge — Cancel auto-merge

### Constraints
- PR must not have merge conflicts
- Required checks must be configured
- Only the PR author or repo admin can schedule

---

## 5. Commit Status

### What
CI integration status tracking for commits. External CI systems report build/test results.

### Behaviors
- Status states: pending, success, error, failure, warning
- Multiple contexts per commit (e.g., ci/travis, ci/lint)
- Combined status view across all contexts
- Target URL for status details
- Descriptive messages
- Required status checks on protected branches

### API Endpoints
- POST /repos/{owner}/{repo}/statuses/{sha} — Create status
- GET /repos/{owner}/{repo}/statuses/{sha} — List statuses
- GET /repos/{owner}/{repo}/commits/{sha}/status — Combined status
- GET /repos/{owner}/{repo}/commits/{sha}/statuses — List per-context

---

## 6. Git Hooks

### What
Server-side Git hooks management: pre-receive, update, post-receive.

### Behaviors
- Standard Git hooks: pre-receive, update, post-receive
- Read/edit via API and UI
- Custom scripts for repository policy enforcement
- Stored in repository `.git/hooks` directory

### UI Routes
- GET /{owner}/{repo}/settings/git-hooks — List hooks
- GET /{owner}/{repo}/settings/git-hooks/{name} — Edit hook

### API Endpoints
- GET /repos/{owner}/{repo}/hooks/git — List Git hooks
- GET /repos/{owner}/{repo}/hooks/git/{id} — Get hook
- PATCH /repos/{owner}/{repo}/hooks/git/{id} — Update hook

### Config
- [security] DISABLE_GIT_HOOKS — disable git hooks (default true for security)

---

## 7. Cron/Scheduled Tasks

### What
Background task scheduling for system maintenance operations.

### Tasks
- Repository mirror sync (every 10 minutes)
- Repository health checks (daily)
- Repository stats update (daily)
- Repository archive cleanup (daily)
- External user sync (daily, LDAP)
- Actions cleanup:
  - Stop zombie tasks (every 5 minutes)
  - Stop endless tasks (every 30 minutes)
  - Cancel abandoned jobs (every 6 hours)
  - Run scheduled workflows (variable schedule)

### UI Routes
- GET /admin/monitor/cron — Cron task list (admin)

### Config
- [cron] section: each task has ENABLED, SCHEDULE, RUN_AT_START
- Custom schedule per task using cron syntax

---

## 8. Queue System

### What
Concurrent task processing infrastructure. Used by webhooks, indexing, mailer, actions, etc.

### Queue Types
- Simple Queue: standard FIFO
- Unique Queue: prevents duplicate items (key-based dedup)

### Backends
- Channel: in-memory, single instance
- LevelDB: persistent, single instance (default for production)
- Redis: distributed, multi-instance
- Dummy: immediate no-op (testing)

### Features
- Worker pool with dynamic sizing (based on CPU count)
- Batch processing
- Error handling and retry
- Graceful shutdown
- Queue monitoring via admin panel

### Config
- [queue] TYPE — backend type
- [queue] DATADIR — LevelDB path
- [queue] LENGTH — max queue size (default 100000)
- [queue] BATCH_LENGTH — items per batch (default 20)
- [queue] MAX_WORKERS — worker count (default CPU/2, max 10)
- [queue] CONN_STR — Redis connection string
