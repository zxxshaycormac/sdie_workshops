## Why

Actions workflow runs aggregate their displayed status from all jobs. The current
aggregation and list filtering only account for success, failure, waiting, and
running, so skipped, cancelled, and blocked runs are misrepresented or hidden on
the Actions pages.

## What Changes

- Aggregate workflow run status across all job statuses, including cancelled,
  skipped, and blocked.
- Preserve terminal status distinctions so cancelled runs display as cancelled
  and all-skipped runs display as skipped.
- Expose blocked, skipped, and cancelled in the Actions run status filter.
- Add focused model tests for the status aggregation and filter status list.

Non-goals:

- No change to how individual job, task, or step statuses are produced.
- No database schema or migration change.
- No public API field shape, route, or Swagger contract change.
- No frontend asset output edits under `public/assets/`.

## Capabilities

### New Capabilities

- `actions-run-status-aggregation`: Browser Actions pages display workflow run
  status from the full set of job statuses and expose those statuses in filters.

### Modified Capabilities

- None.

## Impact

- Affected layers: `models/actions` aggregation and status filter data consumed
  by `routers/web/repo/actions`, `templates/repo/actions`, and
  `web_src/js/components` status rendering.
- Compatibility surface: persisted `action_run.status` values use existing
  enum values; rendered pages and existing JSON response fields keep their
  current names and types.
- Observed facts: `models/actions/run_job.go` owns run status aggregation after
  job status updates, and `models/actions/run_list.go` owns status filter
  entries for the Actions list page.
- Assumption to verify during implementation: the existing template and Vue
  status components already render cancelled, skipped, and blocked icons once
  the model surfaces those statuses.
