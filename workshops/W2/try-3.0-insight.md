# W2 Try 3.0 — Insight

Comparison of the local change on branch `training/w2-try-3` (HEAD `2aceedc05b`) against upstream PR [go-gitea/gitea#32859](https://github.com/go-gitea/gitea/pull/32859) ("Fix incomplete Actions status aggregations"), scored against `workshops/evaluation/score-card.md`.

## Score summary

| # | Dimension | Local | Upstream | Decider |
|---|---|:---:|:---:|---|
| 1 | Correctness | 4 | 5 | Upstream's fail-fast semantics are documented and align with GitHub; local's "all skipped → Skipped" is defensible but undocumented |
| 2 | Architecture & layering | 5 | 5 | Both clean — models stay CRUD-shaped, no upward imports |
| 3 | Naming | 5 | 4 | Upstream's `allSuccessOrSkipped` is awkward; renaming to exported `AggregateJobStatus` widens API surface for no clear reason |
| 4 | Error handling | N/A | N/A | No errors returned |
| 5 | Context propagation | 5 | 5 | Both correctly keep `aggregateJobStatus` ctx-free |
| 6 | Logging | N/A | N/A | Neither adds logging |
| 7 | i18n | N/A | N/A | No new user-visible strings; pre-existing missing `actions.status.*` keys are not introduced by either |
| 8 | Testing | 3 | 5 | Local: 7 cases, misses `{success, skipped}` and most pairwise combos. Upstream: 21 cases + a comment documenting the semantics decision |
| 9 | Security | N/A | N/A | Not touched |
| 10 | Performance | 5 | 5 | Both single-pass O(n) |
| 11 | UX polish | 4 | 3 | Upstream has cross-surface color drift: cancelled is **grey** in `.tmpl`, **yellow** in `.vue` — exactly what `design.md` §9 warns about. Local is consistent (orange both) and plumbs `:class-name` through branches that were silently dropping it |
| 12 | Maintainability | 4 | 5 | Upstream's test comment (`"Should 'running' win? ... GitHub does fail fast"`) documents the non-obvious WHY. Local's priority order is undocumented |
| 13 | Backport / scope | 3 | 4 | Local's `commit_status.go` change introduces new externally-visible behavior (`CommitStatusWarning`) — separate decision from the display bug. Upstream creeps into `RepoActionView.vue` CSS and an export rename, but less significantly |
| | **Average (9 applicable)** | **4.22** | **4.56** | |

**Applicable dimensions:** 9 (Correctness, Architecture, Naming, Context, Testing, Performance, UX polish, Maintainability, Backport/scope). N/A dimensions excluded from the denominator.

## What each implementation does

### Local (branch `training/w2-try-3`) — 6 source files + 2 new test files

- `models/actions/run_job.go` — `aggregateJobStatus` rewritten with explicit per-status `switch` and one boolean per status; priority order is `Failure > Cancelled > Skipped > Success` when all done, `Running > Waiting > Blocked > Unknown` otherwise
- `models/actions/run_list.go` — `GetStatusInfoList` extended from 4 statuses to 7
- `services/actions/commit_status.go` — `toCommitStatus` remaps Cancelled/Skipped → `api.CommitStatusWarning` (new mapping upstream did not touch)
- `templates/repo/actions/status.tmpl` + `web_src/js/components/ActionRunStatus.vue` — new `octicon-stop` branch (orange) for cancelled; `:class-name` plumbed through consistently
- `web_src/js/svg.js` — registers `octicon-stop` (correctly uses `.js`, not upstream's `.ts` — this is the 1.22.x branch)
- Tests: `run_job_test.go` (7 aggregate cases + `GetStatusInfoList` assertions), `commit_status_test.go` (7 `toCommitStatus` cases)

### Upstream PR #32859 — 5 source files + 1 new test file

- `models/actions/run_job.go` — function **renamed** to exported `AggregateJobStatus`; rewritten with compact boolean ORs and a `switch { case ... }`; "failure wins" semantics (failure beats running)
- `templates/repo/actions/status.tmpl` — header collapsed with `Iif` helper; cancelled → `octicon-stop` **grey**
- `web_src/js/components/ActionRunStatus.vue` — cancelled → `octicon-stop` **yellow**; prop types tightened to literal unions; `v-else` replaces `.includes()`
- `web_src/js/components/RepoActionView.vue` — unrelated CSS alignment tweak
- `web_src/js/svg.ts` — registers `octicon-stop`
- Tests: `run_job_status_test.go` (21 cases covering the cartesian product on the three anchor statuses)

## Key insights

1. **Upstream wins on rigor, local wins on consistency.** The 0.34-point gap is driven mostly by testing (5 vs 3) and the documented semantics decision in maintainability (5 vs 4). Both are pure process discipline — adding cases and a comment would close most of the gap.

2. **Upstream ships a real visual bug.** Cancelled renders grey in the server template and yellow in the client Vue component. Same enum value, two glyphs. `design.md` §9 calls this out explicitly: *"Where rendering is duplicated across surfaces ... distinctness must hold in every copy."* Local avoids this entirely by using orange in both surfaces.

3. **Local's `toCommitStatus` change is the most interesting divergence.** Upstream left the commit-status mapping alone (Cancelled still serializes as `failure` over the API). Local remaps Cancelled/Skipped to `CommitStatusWarning`. That's *more correct* — the API now distinguishes a cancelled run from a failed one — but it is **scope creep** for a ticket titled "fix incomplete Actions status aggregations." It changes external API behavior and arguably belongs in its own PR with its own spec delta. Scored as a standalone improvement it would be a 4–5; as part of this bugfix it drags the backport score down.

4. **Neither fully resolves the glyph-collision latent bug.** Both still map `unknown` and `failure` to the same red `octicon-x-circle-fill`. `design.md` §9 allows this *only* as a deliberate, commented choice — neither file has the comment. Shared technical debt both PRs walked past.

5. **Branch adaptation is correct in local.** This is the 1.22.x branch, which still uses `web_src/js/svg.js`. Upstream's PR targeted main and modified `svg.ts`. Local correctly backports the asset registration to `svg.js` rather than blindly copying the upstream file path — easy to get wrong, done right here.

## Highest-leverage improvements to local

- Add the missing test cases — at minimum `{success, skipped}`, `{failure, running}`, `{success, cancelled}` — and a comment on the `allDone/hasSkipped → Skipped` branch explaining why that is the chosen semantics.
- Pull the `toCommitStatus` change into a follow-up PR (or add an explicit spec delta and acceptance scenario for it).
- Add a `// unknown shares failure's glyph deliberately` comment near the failure/unknown branch in both `status.tmpl` and `ActionRunStatus.vue` to satisfy `design.md` §9.
