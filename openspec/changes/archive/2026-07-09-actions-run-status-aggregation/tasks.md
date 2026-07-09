## 1. Run-Status Aggregation

- [x] 1.1 Rewrite `aggregateJobStatus` in `models/actions/run_job.go` with explicit per-status flags and precedence: terminal `failure → cancelled → success → skipped`; in-progress `running → blocked → waiting`; `unknown` falls back to `running`.
- [x] 1.2 Add table-driven unit test `TestAggregateJobStatus` in `models/actions/run_job_test.go` covering every status combination (blocked, cancelled, skipped, mixed, empty).
- [x] 1.3 Verify tests fail against the old logic (RED) then pass after the rewrite (GREEN); run `go test -tags 'sqlite sqlite_unlock_notify' ./models/actions/`.

## 2. Runs-List Status Filter

- [x] 2.1 Extend `GetStatusInfoList` in `models/actions/run_list.go` to expose `blocked`, `cancelled`, and `skipped` alongside the existing four statuses.

## 3. Verification

- [x] 3.1 `gofmt` and `go vet` clean on changed files.
- [x] 3.2 Build consumers (`models/actions`, `routers/web/repo/actions`, `services/actions`) with `-tags 'sqlite sqlite_unlock_notify'`.
- [x] 3.3 Confirm downstream consumers need no changes: `runs_list.tmpl`, `status.tmpl`, `view.go`, `ActionRunStatus.vue`, `badge.go`, `commit_status.go`, `updateRepoRunsNumbers`.

## 4. Specification

- [x] 4.1 Modify `CICD-01-104` to describe precedence-based terminal run-status outcomes.
- [x] 4.2 Add `CICD-01-705` (state-driven) for in-progress run-status precedence.
