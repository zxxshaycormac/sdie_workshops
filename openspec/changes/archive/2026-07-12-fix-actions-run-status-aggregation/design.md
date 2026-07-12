## Context

Control path traced for the visible Actions pages:

- Route: `routers/web/web.go` mounts `/{username}/{reponame}/actions` to
  `routers/web/repo/actions.List` and `/{username}/{reponame}/actions/runs/{run}`
  to `View`/`ViewPost`.
- List handler: `routers/web/repo/actions/actions.go` reads the optional
  `status` query parameter into `actions_model.FindRunOptions.Status` and passes
  `actions_model.GetStatusInfoList(ctx)` to the status dropdown.
- Persistence/query: `models/actions/run_list.go` maps selected statuses into
  SQL conditions through `FindRunOptions.ToConds`.
- Aggregation: `models/actions/run.go` creates initial `ActionRunJob` rows in
  `InsertRun`, and `models/actions/run_job.go` calls `aggregateJobStatus` inside
  `UpdateRunJob` after job status changes. Both paths must persist
  `ActionRun.Status` from the same aggregation rule.
- Rendering: `templates/repo/actions/status.tmpl` and
  `web_src/js/components/ActionRunStatus.vue` already accept success, skipped,
  waiting, blocked, running, failure, cancelled, and unknown. No generated
  assets are edited.

Current state:

- `aggregateJobStatus` distinguishes only all-done success/failure,
  all-waiting waiting, and fallback running.
- `InsertRun` accepts the caller-provided initial run status even when all
  created jobs are blocked.
- Cancelled jobs are treated as failure in the run aggregate.
- All-skipped and success-plus-skipped jobs are not distinguished.
- Blocked is not returned from aggregation and is missing from
  `GetStatusInfoList`, so blocked runs cannot be selected in the list filter.

## Goals / Non-Goals

**Goals:**

- Make run aggregation cover the MECE status branches: empty/unknown,
  all skipped, success with skipped, failure, cancelled, waiting, running, and
  blocked.
- Keep pending statuses pending until all jobs are terminal, with running and
  waiting taking precedence over blocked.
- Expose cancelled, skipped, and blocked in the Actions status filter.
- Add focused model tests before production changes and verify through the
  package-level Go test harness plus strict OpenSpec validation.

**Non-Goals:**

- Do not change how jobs, tasks, or steps transition into their statuses.
- Do not change route shape, request/response field names, templates, Vue
  components, generated assets, or database schema.
- Do not add broad end-to-end coverage for this isolated model behavior unless
  focused tests reveal a wider routed integration risk.

## Decisions

1. Use model-level aggregation as the source of truth.

   `ActionRun.Status` is initialized by `InsertRun` and later persisted by
   `UpdateRunJob`, and both list and detail pages already consume that stored
   status. Fixing aggregation at this layer avoids divergent UI-only status
   calculations and keeps badges/details/list behavior aligned.

   Alternative rejected: derive corrected status in `routers/web/repo/actions`
   during rendering. That would leave persisted run status, badge output, and
   filtering inconsistent.

2. Aggregate with explicit status precedence.

   The intended precedence is:

   - no jobs: unknown
   - all skipped: skipped
   - all success or skipped with at least one success: success
   - any running: running
   - any waiting: waiting
   - any blocked: blocked
   - any cancelled: cancelled
   - any failure: failure
   - otherwise: unknown

   This preserves cancelled and blocked as first-class run states while still
   treating running/waiting/blocked jobs as pending work.

   Alternative rejected: only add skipped/cancelled/blocked to the old
   `allDone/allWaiting` booleans. That keeps the hidden precedence implicit and
   makes mixed statuses easier to regress.

3. Update filter status data in the model helper.

   `GetStatusInfoList` is the data source for the status dropdown. Adding the
   missing existing enum values there makes the list page able to filter those
   persisted run states without template changes.

   Alternative rejected: hard-code extra dropdown entries in the template. That
   duplicates enum knowledge outside `models/actions`.

## Risks / Trade-offs

- Existing persisted runs with incorrectly aggregated statuses will not be
  backfilled. Mitigation: this change fixes future updates without schema or
  migration risk; old rows can be corrected by rerun/status updates if needed.
- Mixed terminal precedence can be surprising when failure and cancelled appear
  together. Mitigation: tests document cancelled as distinct and terminal when
  no pending job remains.
- UI templates are not directly changed. Mitigation: traced template and Vue
  components already render cancelled, skipped, and blocked status names.

## Migration Plan

Deploy the code change normally. Rollback is reverting the model helper and
aggregation changes; no schema or generated assets are involved.

## Open Questions

None.
