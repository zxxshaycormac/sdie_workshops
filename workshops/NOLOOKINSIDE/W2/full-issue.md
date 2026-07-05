# Issue #32857: Incorrect Actions Job Status Aggregation

> Source: [go-gitea/gitea#32857](https://github.com/go-gitea/gitea/issues/32857)
> Fix PR: [#32859](https://github.com/go-gitea/gitea/pull/32859) (**Do NOT read before completing the workshop**)
> Category: `gitea-actions` · Type: Bug · Maintenance: Corrective (behavior change)

---

## Summary

In Gitea Actions, a workflow **run** contains multiple **jobs**. The UI must aggregate
all job statuses to display a single overall run status. The current aggregation function
is a **state machine with missing transitions** — it only reasons about 4 of the 8 possible
job statuses, producing wrong run-level statuses for common scenarios.

## User-Visible Symptoms

| Scenario (all jobs in the run) | Expected run status | Actual run status |
|---|---|---|
| All `skipped` | `skipped` | `success` |
| One or more `cancelled` (rest `success`) | `cancelled` | `failure` |
| All `cancelled` | `cancelled` | `failure` |
| One `blocked` (rest `success`) | `blocked` / `waiting` | `running` |
| All `blocked` | `blocked` / `waiting` | `running` |

These are **not** cosmetic — CI consumers (PR checks, commit status badges, webhook
payloads, UI badges) all read the aggregated run status, so users are misled about whether
the run actually succeeded, was aborted, or is still pending.

## Root Cause (to be confirmed via systematic-debugging)

The aggregation logic lives in `models/actions/run_job.go:155` — `aggregateJobStatus(jobs)`.
It uses **three boolean flags** to summarize the job set:

```go
allDone    := true   // every job has reached a terminal state?
allWaiting := true   // every job is either waiting or done?
hasFailure := false  // any job failed or was cancelled?
```

### Why booleans cannot express the full state space

The `Status` enum (`models/actions/status.go:15`) defines **8** values:

| # | Status | `IsDone()` | Has a runner result? |
|---|---|---|---|
| 0 | `StatusUnknown` | no | no |
| 1 | `StatusSuccess` | yes | yes |
| 2 | `StatusFailure` | yes | yes |
| 3 | `StatusCancelled` | yes | yes |
| 4 | `StatusSkipped` | yes | yes |
| 5 | `StatusWaiting` | no | no (scheduler state) |
| 6 | `StatusRunning` | no | no (scheduler state) |
| 7 | `StatusBlocked` | no | no (scheduler state) |

The boolean reducer collapses this 8-state space into ~3 buckets and loses information:

- **`cancelled` is folded into `failure`** — `hasFailure` is set for both, so a cancelled
  run reports as `failure`.
- **`skipped` is invisible** — it counts as "done" but sets neither `hasFailure` nor any
  positive signal, so an all-skipped run reports as `success`.
- **`blocked` is indistinguishable from `running`** — neither is `IsDone()` nor `StatusWaiting`,
  so both flip `allDone=false` and `allWaiting=false`, landing on the default `running` branch.

### Frontend compounding

`web_src/js/components/ActionRunStatus.vue:37` renders `failure`, `cancelled`, and `unknown`
with the **same** red `octicon-x-circle-fill` icon, so even if the backend returned `cancelled`,
the UI could not visually distinguish it from `failure`.

## Scope

| Layer | File | Concern |
|---|---|---|
| Backend (core) | `models/actions/run_job.go` | `aggregateJobStatus` — the bug lives here |
| Backend (enum) | `models/actions/status.go` | Source of truth for the 8 statuses (read-only context) |
| Frontend | `web_src/js/components/ActionRunStatus.vue` | Status icon rendering |
| Frontend (mirror) | `templates/repo/actions/status.tmpl` | Server-side twin of the Vue component |
| Test (new) | `models/actions/run_job_status_test.go` | Table-driven coverage of all status combinations |

**Affected call path:** `UpdateRunJob` (`run_job.go:94`) recomputes `run.Status =
aggregateJobStatus(jobs)` on every job status change and persists it, so the wrong value
propagates to every downstream consumer.

## Expected Behavior (Spec — to be extracted by students)

A correct aggregator must:

1. **Enumerate all 8 statuses** explicitly — no implicit "else → running" fallthrough.
2. **Preserve distinct semantics** for `cancelled` vs. `failure` vs. `skipped`.
3. **Define a priority order** for mixed-status runs. One defensible ordering:
   `success` > `failure` > `running` > `waiting` > `blocked` > `cancelled` > `skipped`
   (a single failing job taints the whole run; a single running job means it is not done).
4. **Handle the all-`skipped` case** — the run never actually executed, so the aggregate
   should reflect that, not masquerade as `success`.
5. **Render `cancelled` distinctly** in the UI (e.g. `octicon-stop` instead of the failure icon).

The precise priority is a design decision — students must justify their choice, not copy one.

## Reproduction

1. Start a local Gitea instance with Actions enabled and a runner registered.
2. Create a workflow where **all** jobs are skipped, e.g. a job gated on an `if:` that
   evaluates false, or jobs whose `needs:` dependency was skipped.
3. Trigger the workflow and observe the run list / run detail page.
4. The aggregated run status will read **`success`** (or `running` in some edge cases) instead
   of **`skipped`**.
5. Repeat with a run where one job is cancelled mid-flight → aggregate shows `failure`.

## Acceptance Criteria

- [ ] `aggregateJobStatus` returns `skipped` when all jobs are `skipped`.
- [ ] `aggregateJobStatus` returns `cancelled` (not `failure`) when any job is `cancelled`
      and none failed.
- [ ] `aggregateJobStatus` no longer collapses `blocked` into `running`.
- [ ] A **table-driven** Go test covers every status combination (single-status and
      representative multi-status mixes) — written **before** the fix (Red phase).
- [ ] The frontend renders `cancelled` with a distinct icon from `failure`.
- [ ] All pre-existing tests in `models/actions/` still pass.

## Workshop Mapping (Reverse SDIE: E → S → I → E)

| Phase | Activity | Skill |
|---|---|---|
| **(E)val — reverse start** | Reproduce symptom → enumerate all 8 statuses → read `aggregateJobStatus` → spot the gap | `systematic-debugging` |
| **(S)pec — reverse extract** | Derive the full aggregation rule (priority table) from the status enum | — |
| **(I)mpl — fix** | Red: failing table-driven test · Green: rewrite aggregator · Refactor: simplify | `test-driven-development` |
| **(E)val — regression** | All tests green + cancelled icon distinct + original scenarios fixed | — |

## Constraints for Students

- **Do NOT** read PR #32859 or its diff. Derive the fix from the bug.
- **Do NOT** guess the solution — follow the ReAct loop: phenomenon → hypothesis →
  verification → root cause.
- **Do NOT** skip the Red phase. Write the failing test first; it is the contract that
  proves the fix and guards against regression.
- **DO** enumerate every status combination in a table before writing code. MECE is the
  self-check that you have not missed a transition.

## References

- Status enum: `models/actions/status.go:15`
- Buggy aggregator: `models/actions/run_job.go:155`
- Call site: `models/actions/run_job.go:140` (inside `UpdateRunJob`)
- Frontend icon map: `web_src/js/components/ActionRunStatus.vue:32`
- Gitea Actions status docs: https://docs.gitea.com/usage/actions/
