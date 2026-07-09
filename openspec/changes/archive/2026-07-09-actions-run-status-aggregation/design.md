## Context

A Gitea Actions workflow **run** owns many **jobs**, each carrying its own `Status`. The run exposes a single aggregated `Status` that drives the badge in the runs list, the run detail page, and the SVG workflow badge. Today that aggregate is computed by one unexported function — `aggregateJobStatus` in `models/actions/run_job.go` — which is invoked from `UpdateRunJob` on every job-status change and **persisted** into the `action_run.status` column. Every downstream consumer (templates, Vue view, badge endpoint, repo counters) reads that persisted value verbatim; nothing recomputes it on the fly.

The current aggregator tracks only three booleans (`allDone`, `allWaiting`, `hasFailure`) and can therefore emit only `success`, `failure`, `waiting`, or `running`. It collapses the other lifecycle states that `CICD-01-003` already enumerates: `blocked` becomes `running`, `cancelled` becomes `failure`, and `skipped` becomes `success`. The result is a badge that misrepresents reality in common scenarios, and a spec (`CICD-01-104`) that hard-codes the same binary terminal model.

## Goals / Non-Goals

**Goals:**
- Make the run badge faithfully reflect `blocked`, `cancelled`, and `skipped` as distinct states, alongside `success`, `failure`, `running`, and `waiting`.
- Define a single, documented precedence for combining job statuses into a run status, matching GitHub Actions semantics.
- Correct `CICD-01-104` and add the missing aggregation requirement so the spec and implementation agree.
- Guard the behavior with a table-driven unit test that covers the full status truth table.

**Non-Goals:**
- Backfilling/re-aggregating **historical** runs already persisted with the old collapsed statuses (no data migration).
- Changing the per-job **commit-status** mapping (`services/actions/commit_status.go`), which is independent and already correct.
- Frontend changes — `status.tmpl` and `ActionRunStatus.vue` already render all eight statuses.
- Changing the `action_run.status` column type or adding a DB migration (the existing `int` enum already holds all values).

## Decisions

**D1 — Terminal precedence: `failure → cancelled → success → skipped`.**
A `failure` dominates a `cancelled` because the canonical case is "a job failed, so its `needs`-dependent siblings were torn down (cancelled)" — the actionable root cause is the failure, mirroring GitHub Actions. A `cancelled` dominates `success` (confirmed as the desired UX: a run with some succeeded and some cancelled jobs reads as cancelled, i.e. the user tore the run down). A run where nothing actually executed (only `skipped` jobs) is reported as `skipped`, not `success`.

**D2 — In-progress precedence: `running → blocked → waiting`.**
An executing job dominates everything else (the run is observably active). A `blocked` job (awaiting environment approval or unmet `needs`) dominates a merely `waiting` (queued) job, because blocked requires user action and should surface distinctly rather than as a generic spinner — this is what `CICD-01-306` already implies ("hold the run in blocked state").

**D3 — `unknown` falls back to `running`.**
`StatusUnknown` (the zero value) is not a real lifecycle state. It is non-terminal, so the run must stay in-progress; `running` is retained as the fallback to preserve prior behavior. Empty job sets are degenerate and fall through to `skipped` ("nothing ran").

**D4 — Persisted, not recomputed.**
We keep the existing architecture: recompute on job change, persist once, read everywhere. Recomputing on read would touch many consumers and risk divergence; the persisted model is correct once the aggregator is correct.

**D5 — Expose the new statuses in the runs-list filter.**
`GetStatusInfoList` is extended so users can filter by `blocked`/`cancelled`/`skipped`, which the aggregator can now legitimately produce. The `status=` query param already filters directly against the `action_run.status` column, so no query changes are needed.

## Risks / Trade-offs

- **Historical runs keep old statuses.** Runs persisted before this change still show the collapsed status (e.g. an old cancelled run appears as `failure`). → Mitigation: documented as a non-goal; a one-off re-aggregation script can be run later if needed, but is not required for correctness going forward.
- **Counter semantics shift slightly.** `updateRepoRunsNumbers` counts "closed" runs as status ∈ {success, failure, cancelled, skipped}. The set is unchanged, but the distribution becomes more accurate (cancelled runs now counted as cancelled, not failure). → No action; this is strictly an improvement.
- **No spec partitioning into per-precedence-arm requirements.** We express each precedence as a single requirement (mirroring the repo's existing `CICD-07-101` aggregation style) rather than one requirement per arm. → Trade-off: slightly less granular requirement-to-test mapping, but the table-driven unit test (`run_job_test.go`) already covers every arm explicitly.

## Migration Plan

No schema migration. Deploy is code-only. Existing in-flight runs will be re-aggregated correctly on their next job-status change; completed historical runs are unaffected in behavior (they remain terminal) but retain their previously-computed (possibly collapsed) status label.

## Open Questions

- Should a one-off administrative backfill re-aggregate historical runs? Decision deferred (non-goal for this change); revisit if users report confusing historical badges.
