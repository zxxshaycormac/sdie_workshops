## MODIFIED Requirements

### Event-Driven Requirements (Workflow Execution)

- **CICD-01-101:** `When a trigger event matches a workflow's on clause, the system shall create a new workflow run.`
- **CICD-01-102:** `When a workflow run is created, the system shall evaluate job dependencies and schedule jobs for execution.`
- **CICD-01-103:** `When a job is assigned to a runner, the system shall execute steps sequentially within that job.`
- **CICD-01-104:** `When a workflow run completes (all jobs in a terminal state), the system shall set the run status to the most severe job outcome with precedence: failure if any job failed, otherwise cancelled if any job was cancelled, otherwise success if any job succeeded, otherwise skipped.`
- **CICD-01-105:** `When a user cancels a workflow run via the API, the system shall cancel all running and pending jobs in that run.`
- **CICD-01-106:** `When a user reruns a workflow via the API, the system shall create a new run with the same workflow definition and trigger event.`
- **CICD-01-107:** `When a workflow step references an action (uses:), the system shall resolve the action from the configured DEFAULT_ACTIONS_URL (GitHub or self-hosted).`

---

## ADDED Requirements

### 1. Run Status Aggregation

**User Story:** As a developer, I want the workflow run status badge to accurately reflect the combined state of every job in the run, so that I can tell at a glance whether a run is executing, blocked awaiting approval, failed, cancelled, or skipped.

#### State-Driven Requirements (In-Progress Run Status)

- **CICD-01-705:** `While a workflow run has any non-terminal job, the system shall set the run status with precedence: running if any job is running, otherwise blocked if any job is blocked (awaiting environment approval or unmet job needs), otherwise waiting.`

---

## Edge Cases & Error Handling

| Scenario | System Behavior |
|----------|----------------|
| All jobs `skipped` (e.g. all `if:` conditions false) | Run status is `skipped` (not `success`) |
| A job is cancelled and no job failed | Run status is `cancelled` (not `failure`) |
| A failed job caused dependent jobs to be cancelled | Run status is `failure` (failure dominates cancelled) |
| All jobs `blocked` awaiting approval | Run status is `blocked` (not `running`) |
| Mix of `blocked` and `waiting` jobs | Run status is `blocked` |
| Mix of `running` and `blocked` jobs | Run status is `running` |
| Job in `unknown` (zero-value) non-terminal state | Run falls back to `running` |
| Empty job set (degenerate) | Run status is `skipped` |
