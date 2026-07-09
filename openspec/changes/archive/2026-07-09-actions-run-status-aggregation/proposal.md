## Why

The overall status badge shown for a Gitea Actions workflow **run** (in the runs list, the run detail page, and the SVG workflow badge) is wrong in common scenarios: a run whose jobs are all `blocked` (awaiting environment approval or unmet `needs`) shows a spinning "running" icon; a cancelled run shows "failure"; and a run where every job was `skipped` shows a green "success" check. The single aggregation function `aggregateJobStatus` (`models/actions/run_job.go`) could only ever emit `success`/`failure`/`waiting`/`running`, collapsing `blocked`, `cancelled`, and `skipped` into incorrect values. Its result is persisted into `action_run.status` on every job change, so the wrong value propagates to every downstream consumer. The spec also reflects this gap: requirement `CICD-01-104` hard-codes the terminal run status to "success or failure", encoding the very bug being fixed.

## What Changes

- **Run-status aggregation redefined.** `aggregateJobStatus` now surfaces `blocked`, `cancelled`, and `skipped` as first-class run states with an explicit precedence:
  - Terminal (all jobs done): **failure → cancelled → success → skipped**.
  - In-progress (any job non-terminal): **running → blocked → waiting**.
- **Spec correction.** `CICD-01-104` is modified to describe the full set of terminal outcomes with precedence, and a new state-driven requirement (`CICD-01-705`) defines the in-progress aggregation. This aligns the spec with the already-enumerated lifecycle states in `CICD-01-003` and the blocked-state expectation in `CICD-01-306`.
- **Runs-list status filter.** `GetStatusInfoList` (`models/actions/run_list.go`) now exposes `blocked`, `cancelled`, and `skipped` as filterable statuses, since the aggregator can legitimately produce them.

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `cicd-automation`: `CICD-01-104` (terminal run-status outcomes) changes from a binary "success/failure" rule to a precedence-based rule over `failure`/`cancelled`/`success`/`skipped`; a new `CICD-01-705` defines in-progress run-status precedence (`running`/`blocked`/`waiting`).

## Impact

- **Code:**
  - `models/actions/run_job.go` — `aggregateJobStatus` rewritten (behavior change, no signature change).
  - `models/actions/run_list.go` — `GetStatusInfoList` extended (no signature change).
  - `models/actions/run_job_test.go` — new table-driven test covering the corrected truth table.
- **Downstream consumers (read `action_run.status` verbatim, no code change required):**
  - Runs-list badge — `templates/repo/actions/runs_list.tmpl`, `templates/repo/actions/status.tmpl`.
  - Run detail badge — `routers/web/repo/actions/view.go:159` (JSON) → `web_src/js/components/ActionRunStatus.vue`.
  - SVG workflow badge — `routers/web/repo/actions/badge.go` + `modules/badge/badge.go` (color map already covers all 8 statuses).
  - Repo run counters — `models/actions/run.go:170` `updateRepoRunsNumbers` (closed-run set `{success,failure,cancelled,skipped}` unchanged; counting becomes more accurate).
  - Commit-status mapping — `services/actions/commit_status.go` (per-job, independent; already maps all 8 statuses correctly).
- **API / routers:** No REST API or route changes. **No Swagger annotation changes** are required (no handler signatures or response shapes change).
- **Database:** No schema migration — `action_run.status` is an existing `int` column whose enum already includes all 8 values.
- **Backward compatibility:** Historical runs already persisted with the old collapsed statuses (e.g. cancelled-as-failure) are **not** retroactively re-aggregated by this change; only future runs are affected. A one-off backfill is out of scope (see design.md "Non-Goals").
