# Issue: Actions Run Status Shows Wrong Result

## Summary

In Gitea Actions, a workflow **run** is made up of multiple **jobs**. The run list and run
detail pages display a single **overall status** for each run, computed by combining the
statuses of all its jobs. Users rely on this overall status to tell, at a glance, whether a
run passed, failed, was aborted, or is still in progress.

Currently this overall status is **wrong in several common situations**. The badge shown to
users does not match what actually happened in the run.

## User-Visible Symptoms

A user sets up a workflow with several jobs and triggers a run. Depending on what happens
to the jobs, the page reports an incorrect overall status:

| What actually happened in the run | What the user expects to see | What the user actually sees |
|---|---|---|
| Every job was **skipped** (e.g. an `if:` condition was false for all of them) | **Skipped** — the run never really executed | **Success** — looks like everything passed |
| One or more jobs were **cancelled** by the user mid-run | **Cancelled** — the user aborted it | **Failure** — looks like something broke |
| A job is **blocked** waiting on a dependency or concurrency slot | **Blocked** or **Waiting** | **Running** — looks like work is happening when nothing is |
| A mix of skipped and cancelled jobs | **Cancelled** (or at least not "success") | **Success** or **Failure** |

These are not cosmetic glitches. The wrong status propagates everywhere a user looks:

- The **run list** and **run detail** badges.
- The **commit status** shown next to commits and on PR checks.
- **Webhook payloads** sent to external systems that react to run outcomes.
- Any automation that gates merges on a green check.

So a user may believe a run passed when it was actually skipped, or think a run failed when
the user themselves cancelled it.

## Why This Matters

- **False success is dangerous.** A team gating merges on CI could merge code that was never
  actually tested because every job was skipped — the run proudly shows green.
- **Cancelled ≠ Failure.** A user who aborts a run knows they aborted it; seeing it reported
  as a failure pollutes failure notifications, metrics, and dashboards.
- **Blocked is invisible.** A run stuck waiting on a concurrency group or an unmet dependency
  appears to be actively running, so the user never investigates the real blocker.

## Reproduction

1. Enable Gitea Actions and register a runner.
2. Create a workflow where **all** jobs will be skipped — for example, gate every job on an
   `if:` expression that evaluates to false.
3. Trigger the workflow.
4. Open the run list. The overall status reads **Success**, even though no job ran.
5. Repeat the experiment by cancelling one job in a multi-job run: the overall status flips
   to **Failure** instead of **Cancelled**.

## Expected Behavior

The overall run status should faithfully reflect what happened across all jobs:

- A run where **every job was skipped** should be reported as **Skipped**.
- A run where **any job was cancelled** (and none failed on its own) should be reported as
  **Cancelled**, not lumped under Failure.
- A run that is **blocked** should not pretend to be running.
- The user must never see **Success** for a run that did not actually execute any work.

In short: the overall status should be a truthful summary of the jobs, distinguishable at a
glance, and consistent across the UI, commit checks, and webhooks.

## Acceptance Criteria

- [ ] An all-skipped run shows **Skipped**, not Success.
- [ ] A run with a cancelled job shows **Cancelled**, not Failure (when no job actually failed).
- [ ] A blocked run no longer masquerades as Running.
- [ ] The **Cancelled** status looks visually distinct from **Failure** in the UI (different
      icon), so a user can tell an abort apart from a real failure without reading text.
- [ ] Existing runs that genuinely succeed or fail are still reported correctly.
