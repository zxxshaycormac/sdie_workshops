# W2 Try-1.0 Insight — Actions Run Status Aggregation

Comparative review of the local fix (`10e5457500` — *fix(actions): correct run status aggregation for blocked/cancelled/skipped*) against the upstream canonical fix (go-gitea/gitea PR #32859 — *Fix incomplete Actions status aggregations*, merged 2024-12-16), scored with `workshops/evaluation/score-card.md`.

Both address issue #32857: `aggregateJobStatus` only emitted `success`/`failure`/`waiting`/`running`, collapsing `blocked→running`, `cancelled→failure`, and `skipped→success`.

## Score summary

8 of 13 dimensions apply (4 Error handling, 5 Context, 6 Logging, 7 i18n, 9 Security are N/A — pure function, no I/O, no new strings, no trust boundary).

| # | Dimension | Local `10e5457` | PR #32859 | Evidence / divergence |
|---|-----------|:---:|:---:|---|
| 1 | Correctness | **4** | **4** | Local handles all 7 statuses; gap — failure doesn't win while run in-progress (`[Failure, Blocked]→Blocked`), untested. PR does fail-fast (`[Failure, Running]→Failure`, tested); gap — `all-skipped→Success` and filter not updated. |
| 2 | Architecture | **5** | **5** | Both stay in `models/actions`; no upward imports. |
| 3 | Naming | **5** | **5** | `aggregateJobStatus`/`AggregateJobStatus`, `hasX` booleans, snake_case files. |
| 4 | Error handling | N/A | N/A | Pure function, no error returns. |
| 5 | Context propagation | N/A | N/A | Changed function is pure; `GetStatusInfoList` ctx preserved. |
| 6 | Logging | N/A | N/A | No log lines added. |
| 7 | i18n | N/A | N/A | No new strings; both reuse existing `actions.status.*` keys. |
| 8 | Testing | **4** | **4** | Complementary. Local: 22 cases, strong on in-progress mixes + empty. PR: 22 cases, strong on failure interactions / fail-fast. Neither tests `GetStatusInfoList`. |
| 9 | Security | N/A | N/A | No trust boundary touched. |
| 10 | Performance | **5** | **5** | Single O(n) pass; no DB calls added. |
| 11 | UX polish | **3** | **4** | Local: updates the filter (PR doesn't) but leaves `cancelled` reusing the failure icon. PR: adds dedicated `octicon-stop` for cancelled, tightens Vue types, fixes title layout — but filter stays at 4 statuses. |
| 12 | Maintainability | **4** | **4** | Local: two-phase switch readable; `case StatusSkipped:` is a silent no-op. PR: flat accumulate-then-switch; exports `AggregateJobStatus` with no external caller (speculative API widening). |
| 13 | Backport / scope | **5** | **4** | Local: minimal diff, function unexported. PR: exports function + drive-by template `Iif` refactor + Vue type tightening + adjacent layout fix. |
| | **Average (8 dims)** | **4.38** | **4.38** | Tie in aggregate, complementary in profile. |

## Key insights

**1. Same total, opposite strengths.** The 4.38 tie masks that each is half of a complete fix. The ideal change is the union: local's filter update + scope discipline + all-skipped semantics, plus PR's fail-fast + dedicated cancelled icon.

**2. The most consequential divergence is fail-fast.** PR: `[Failure, Running]→Failure` (GitHub-compatible early signal). Local: failure only surfaces once *every* job is terminal, so `[Failure, Blocked]→Blocked` — actively misleading (tells the user to unblock a run that's already doomed). The local design doc (D1–D2) frames this as a deliberate two-phase choice, but it's the choice most likely to confuse operators in realistic multi-job workflows. This is the one place the local design should yield to the PR's behavior.

**3. Each caught something the other missed.**
- Local wins: `GetStatusInfoList` exposes blocked/cancelled/skipped in the filter (PR's aggregator produces them but the UI never offers them); all-skipped→skipped is more faithful than PR's all-skipped→success.
- PR wins: cancelled gets a distinct icon; fail-fast; tighter Vue prop types.

**4. Scope discipline favors the local.** For the `1.22.x` backport-focused branch, the local's unexported function + 3-file code diff is trivially cherry-pickable. The PR's 6 files (SVG import, template mini-refactor, layout fix in an adjacent component) are fine for trunk but heavier for backport.

**5. Process artifacts didn't move the score — by design.** The local's openspec design.md (D1–D5) is why the choices are auditable, but per the score card's "deliberately not a dimension," it earns no rubric credit. Its value is explanatory, not scoring.

## Recommendation

Take the local as the base (scope + filter + all-skipped) and graft on the PR's fail-fast precedence and the `octicon-stop` icon for cancelled.
