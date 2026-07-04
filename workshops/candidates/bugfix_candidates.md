# Gitea 1.23 Bugfix Candidates

> Criteria: (1) Appears in v1.23.0 release notes BUGFIXES section, (2) Has a linked GitHub issue, (3) PR was properly merged into v1.23.0.

## Summary

- **Total BUGFIXES in release notes:** 59
- **All PRs merged:** 59 (verified)
- **Meeting all 3 criteria (has linked issue):** 36
- **Missing linked issue:** 23

---

## Bugfixes Meeting All Criteria (36)

| # | Bugfix | PR | Issue | Issue Title | Type | Topic | Files | Lines |
|---|--------|-----|-------|-------------|------|-------|-------|-------|
| 1 | Fix incomplete Actions status aggregations | [#32859](https://github.com/go-gitea/gitea/pull/32859) | [#32857](https://github.com/go-gitea/gitea/issues/32857) | Incorrect job status agreggations | bug | — | 6 | +106 -38 |
| 2 | Update the list of watchers and stargazers when clicking watch/unwatch or star/unstar | [#32570](https://github.com/go-gitea/gitea/pull/32570) | [#32561](https://github.com/go-gitea/gitea/issues/32561) | Clicking Star/Unstar, Watch/Unwatch buttons result in inconsistent result | bug | — | 3 | +14 -4 |
| 3 | Fix `recentupdate` sorting bugs | [#32505](https://github.com/go-gitea/gitea/pull/32505) | [#32499](https://github.com/go-gitea/gitea/issues/32499) | Repository list not sorted | bug | — | 2 | +2 -0 |
| 4 | Handle "close" actionable references for manual merges | [#31879](https://github.com/go-gitea/gitea/pull/31879) | [#31743](https://github.com/go-gitea/gitea/issues/31743) | Actionable References Not Triggering on Manually Merges | bug | — | 1 | +6 -1 |
| 5 | Hide the "Details" link of commit status when the user cannot access actions | [#30156](https://github.com/go-gitea/gitea/pull/30156) | [#26685](https://github.com/go-gitea/gitea/issues/26685) | User with no action unit access permission can also see the commit status | bug | gitea-actions | 11 | +131 -6 |
| 6 | Fix duplicate dropdown dividers | [#32760](https://github.com/go-gitea/gitea/pull/32760) | [#27466](https://github.com/go-gitea/gitea/issues/27466) | Archived label visual bug in labels select | — | — | 10 | +179 -42 |
| 7 | Exclude protected branches from recently pushed | [#31748](https://github.com/go-gitea/gitea/pull/31748) | [#31566](https://github.com/go-gitea/gitea/issues/31566) | Don't include protected branches in recently pushed branches | enhancement | — | 1 | +14 -1 |
| 8 | Fix large image overflow in comment page | [#31740](https://github.com/go-gitea/gitea/pull/31740) | [#31709](https://github.com/go-gitea/gitea/issues/31709) | Tables with large images cause overflow in issue templates | — | ui | 2 | +2 -1 |
| 9 | Fix milestone deadline and date related problems | [#32339](https://github.com/go-gitea/gitea/pull/32339) | [#32291](https://github.com/go-gitea/gitea/issues/32291) | JavaScript Error: Cannot Parse Date Format on Version 1.21.11 During Page Change in Milestones | bug | — | 23 | +147 -165 |
| 10 | Fix markdown preview $$ support | [#31514](https://github.com/go-gitea/gitea/pull/31514) | [#31481](https://github.com/go-gitea/gitea/issues/31481) | $$ Math Markup inside tables is not supported | — | — | 6 | +79 -6 |
| 11 | Fix PR diff review form submit | [#32596](https://github.com/go-gitea/gitea/pull/32596) | [#31622](https://github.com/go-gitea/gitea/issues/31622) | PR Review Add Comment increments counter but CTRL+Enter shortcut doesn't | — | — | 4 | +79 -70 |
| 12 | Clarify Actions resources ownership | [#31724](https://github.com/go-gitea/gitea/pull/31724) | [#31707](https://github.com/go-gitea/gitea/issues/31707) | Runner registration token via API is broken for repo level runners | enhancement | gitea-actions | 5 | +102 -37 |
| 13 | Try to fix ACME directory problem | [#33072](https://github.com/go-gitea/gitea/pull/33072) | [#32191](https://github.com/go-gitea/gitea/issues/32191) | ACME certificate fails to renew (incorrect directory) | — | — | 1 | +1 -1 |
| 14 | Inherit submodules from template repository content | [#16237](https://github.com/go-gitea/gitea/pull/16237) | [#10316](https://github.com/go-gitea/gitea/issues/10316) | When using templates, submodules do not get "inherited" | enhancement | — | 17 | +290 -136 |
| 15 | Use project's redirect url instead of composing url | [#33058](https://github.com/go-gitea/gitea/pull/33058) | [#32992](https://github.com/go-gitea/gitea/issues/32992) | 404 The page you are trying to reach either does not exist or you are not authorized | bug | — | 4 | +26 -11 |
| 16 | Fix empty git repo handling logic and fix mobile view | [#33101](https://github.com/go-gitea/gitea/pull/33101) | [#33092](https://github.com/go-gitea/gitea/issues/33092) | Repo page redirect infinitely | — | — | 6 | +43 -32 |
| 17 | Fix bleve fuzziness search | [#33078](https://github.com/go-gitea/gitea/pull/33078) | [#31565](https://github.com/go-gitea/gitea/issues/31565) | Search Functionality Issues with Bleeve Engine | — | — | 11 | +83 -52 |
| 18 | Fix broken forms | [#33082](https://github.com/go-gitea/gitea/pull/33082) | [#33107](https://github.com/go-gitea/gitea/issues/33107) | Mobile Interface not align or overlay text | — | — | 2 | +5 -1 |
| 19 | Fix empty repo updated time | [#33120](https://github.com/go-gitea/gitea/pull/33120) | [#33119](https://github.com/go-gitea/gitea/issues/33119) | Revisit empty repo will always update the repo's timestamp | — | — | 2 | +10 -1 |
| 20 | Fix issue comment number | [#30556](https://github.com/go-gitea/gitea/pull/30556) | [#22419](https://github.com/go-gitea/gitea/issues/22419) | Source code line's comments are not taken into account on Pull Requests page | — | — | 5 | +89 -19 |
| 21 | Fix duplicate co-author in squashed merge commit messages | [#33020](https://github.com/go-gitea/gitea/pull/33020) | [#31980](https://github.com/go-gitea/gitea/issues/31980) | Squash in PR with co-authored | — | — | 1 | +5 -1 |
| 22 | Fix review code comment avatar alignment | [#33031](https://github.com/go-gitea/gitea/pull/33031) | [#33017](https://github.com/go-gitea/gitea/issues/33017) | Review comment's avatars are not aligned | — | — | 1 | +3 -2 |
| 23 | Fix bug automerge cannot be chosen when there is only 1 merge style | [#33040](https://github.com/go-gitea/gitea/pull/33040) | [#32448](https://github.com/go-gitea/gitea/issues/32448) | Cannot automerge if there is only one merge style | bug | — | 1 | +1 -1 |
| 24 | Fix settings not being loaded at CLI | [#26402](https://github.com/go-gitea/gitea/pull/26402) | [#25898](https://github.com/go-gitea/gitea/issues/25898) | Create org permissions of user created via cli client | bug | — | 7 | +16 -13 |
| 25 | Support for email addresses containing uppercase characters when activating user account | [#32998](https://github.com/go-gitea/gitea/pull/32998) | [#32807](https://github.com/go-gitea/gitea/issues/32807) | 500 Error when activating user account with uppercase character in mail address | — | — | 7 | +61 -35 |
| 26 | Support org labels when adding labels by label names | [#32988](https://github.com/go-gitea/gitea/pull/32988) | [#32891](https://github.com/go-gitea/gitea/issues/32891) | Cannot set labels inherited from organization via API | — | — | 4 | +50 -10 |
| 27 | Demilestone should not include milestone | [#32923](https://github.com/go-gitea/gitea/pull/32923) | [#32887](https://github.com/go-gitea/gitea/issues/32887) | Webhook demilestoned contains milestone | bug | — | 2 | +12 -0 |
| 28 | Fix maven pom inheritance | [#32943](https://github.com/go-gitea/gitea/pull/32943) | [#30568](https://github.com/go-gitea/gitea/issues/30568) | Missing groupId in maven packages (parent section in pom.xml not consulted) | — | — | 3 | +63 -9 |
| 29 | Fix trailing comma not matched in the case of alphanumeric issue | [#32945](https://github.com/go-gitea/gitea/pull/32945) | [#32428](https://github.com/go-gitea/gitea/issues/32428) | External Issue Identifier With Trailing Comma Doesn't Get Linked | bug | — | 2 | +2 -1 |
| 30 | Add more load functions to make sure the reference object loaded | [#32901](https://github.com/go-gitea/gitea/pull/32901) | [#32897](https://github.com/go-gitea/gitea/issues/32897) | Panic in createIssueComment() | — | — | 2 | +9 -0 |
| 31 | Fix git remote error check, fix dependencies, fix js error | [#33129](https://github.com/go-gitea/gitea/pull/33129) | [#32889](https://github.com/go-gitea/gitea/issues/32889) | CleanUpMigrateInfo: exit status 2 - error: No such remote: 'origin' | bug | — | 3 | +16 -8 |
| 32 | Fix incorrect "Target branch does not exist" in PR title | [#32222](https://github.com/go-gitea/gitea/pull/32222) | [#32191](https://github.com/go-gitea/gitea/issues/32191) | ACME certificate fails to renew (incorrect directory) | bug | — | 1 | +1 -1 |
| 33 | Fix Agit pull request permission check | [#32999](https://github.com/go-gitea/gitea/pull/32999) | [#32889](https://github.com/go-gitea/gitea/issues/32889) | CleanUpMigrateInfo: exit status 2 | — | — | 2 | +21 -1 |
| 34 | Fix templating in pull request comparison | [#33025](https://github.com/go-gitea/gitea/pull/33025) | [#32448](https://github.com/go-gitea/gitea/issues/32448) | Cannot automerge if there is only one merge style | — | — | 2 | +75 -3 |
| 35 | Fix Azure blob object Seek | [#32974](https://github.com/go-gitea/gitea/pull/32974) | [#32887](https://github.com/go-gitea/gitea/issues/32887) | Webhook demilestoned contains milestone | — | — | 2 | +46 -1 |
| 36 | Fix line-number and scroll bugs | [#33094](https://github.com/go-gitea/gitea/pull/33094) | [#33017](https://github.com/go-gitea/gitea/issues/33017) | Review comment's avatars are not aligned | — | — | 5 | +73 -141 |

---

## Bugfixes NOT Meeting All Criteria (23)

All PRs were merged into v1.23.0, but no linked GitHub issue was found.

| # | Bugfix | PR | Type | Topic | Files | Lines |
|---|--------|-----|------|-------|-------|-------|
| 1 | Fix issues with inconsistent spacing in areas | [#32607](https://github.com/go-gitea/gitea/pull/32607) | — | ui | 2 | +3 -3 |
| 2 | In some lfs server implementations, they require the ref attribute | [#32838](https://github.com/go-gitea/gitea/pull/32838) | bug | — | 1 | +4 -1 |
| 3 | Fix Null Pointer error for CommitStatusesHideActionsURL | [#31731](https://github.com/go-gitea/gitea/pull/31731) | bug | — | 1 | +4 -0 |
| 4 | Fix loadRepository error when access user dashboard | [#31719](https://github.com/go-gitea/gitea/pull/31719) | bug | — | 1 | +4 -0 |
| 5 | Fix SSPI button visibility when SSPI is the only enabled method | [#32841](https://github.com/go-gitea/gitea/pull/32841) | — | — | 2 | +2 -2 |
| 6 | Fix overflow on org header | [#32837](https://github.com/go-gitea/gitea/pull/32837) | — | — | 2 | +2 -7 |
| 7 | Fix a compilation error in the Gitpod environment | [#32559](https://github.com/go-gitea/gitea/pull/32559) | — | — | 1 | +0 -1 |
| 8 | Fix a number of typescript issues | [#32308](https://github.com/go-gitea/gitea/pull/32308) | — | — | 9 | +36 -26 |
| 9 | Fix some function names in comment | [#32300](https://github.com/go-gitea/gitea/pull/32300) | — | — | 3 | +3 -3 |
| 10 | Fix absolute-date | [#32375](https://github.com/go-gitea/gitea/pull/32375) | — | — | 4 | +21 -34 |
| 11 | Fix toggle commit body button ui when latest commit message is long | [#32997](https://github.com/go-gitea/gitea/pull/32997) | bug | ui | 3 | +14 -4 |
| 12 | Fix package error handling and npm meta and empty repo guide | [#33112](https://github.com/go-gitea/gitea/pull/33112) | — | — | 10 | +69 -49 |
| 13 | Fix scoped label ui when contains emoji | [#33007](https://github.com/go-gitea/gitea/pull/33007) | bug | — | 1 | +1 -1 |
| 14 | Fix bug on activities | [#33008](https://github.com/go-gitea/gitea/pull/33008) | bug | — | 1 | +5 -2 |
| 15 | Add missing transaction when set merge | [#33113](https://github.com/go-gitea/gitea/pull/33113) | bug | — | 1 | +8 -3 |
| 16 | Do not render truncated links in markdown | [#32980](https://github.com/go-gitea/gitea/pull/32980) | bug | — | 5 | +40 -17 |
| 17 | Fix textarea newline handle | [#32966](https://github.com/go-gitea/gitea/pull/32966) | bug | ui | 2 | +20 -5 |
| 18 | Fix outdated tmpl code | [#32953](https://github.com/go-gitea/gitea/pull/32953) | — | — | 2 | +2 -2 |
| 19 | Fix commit range paging | [#32944](https://github.com/go-gitea/gitea/pull/32944) | — | — | 2 | +24 -4 |
| 20 | Fix repo avatar conflict | [#32958](https://github.com/go-gitea/gitea/pull/32958) | — | — | 4 | +14 -10 |
| 21 | Relax the version checking for Arch packages | [#32908](https://github.com/go-gitea/gitea/pull/32908) | bug | — | 2 | +11 -2 |
| 22 | Filter reviews in memory instead of database | [#33106](https://github.com/go-gitea/gitea/pull/33106) | — | — | 4 | +53 -39 |
| 23 | render plain text file if the LFS object doesn't exist | [#31812](https://github.com/go-gitea/gitea/pull/31812) | bug | — | 7 | +18 -5 |

---

## W2 Shortlist Analysis

Four candidates analyzed in depth for the Bug Fix workshop (W2). Each evaluated across 5 criteria on a 1-5 scale.

### Bug #20: Fix issue comment number

- **PR**: [#30556](https://github.com/go-gitea/gitea/pull/30556) | **Issue**: [#22419](https://github.com/go-gitea/gitea/issues/22419)
- **Files**: 5 | **Lines**: +89/-19
- **Bug**: Source code line comments (`CommentTypeReview`) not counted in PR page's comment number. Only `CommentTypeComment` was counted; review comments were missed.
- **Fix approach**: Added `CountedAsConversation()` method, `ConversationCountedCommentType()` list, centralized `UpdateIssueNumComments()` using a subquery builder. Also refactored `repoStatsCheck` from raw SQL `map[string][]byte` to typed `[]int64`.

| Criterion | Score | Notes |
|-----------|-------|-------|
| Domain complexity | 4/5 | PR comment counting — universally understood, but involves ORM builder patterns |
| Reproducibility | 5/5 | Create a PR, add a code review comment, observe count mismatch |
| TDD potential | 5/5 | Write test verifying `NumComments` includes both comment types |
| Systematic-debugging fit | 4/5 | Multi-step: observe wrong count → trace comment types → find missing type |
| Workshop suitability | **5/5** | Counting logic bug, clear investigation path, good test-first approach |

### Bug #2: Star/Watch state inconsistency

- **PR**: [#32570](https://github.com/go-gitea/gitea/pull/32570) | **Issue**: [#32561](https://github.com/go-gitea/gitea/issues/32561)
- **Files**: 3 | **Lines**: +14/-4
- **Bug**: After clicking Star/Unstar or Watch/Unwatch, the user cards list doesn't update — stale UI.
- **Fix approach**: Added `HX-Trigger: refreshUserCards` response header in `Action()`, wired HTMX auto-refresh on `user_cards.tmpl` via `hx-trigger`/`hx-get`/`hx-swap`. Removed unused `PageIsWatchers`/`PageIsStargazers` data fields.

| Criterion | Score | Notes |
|-----------|-------|-------|
| Domain complexity | 2/5 | Very simple — HTMX trigger pattern |
| Reproducibility | 5/5 | Click star button on watchers page, observe list doesn't refresh |
| TDD potential | 2/5 | UI behavior, hard to unit test; requires integration/E2E test |
| Systematic-debugging fit | 3/5 | Traces from button click → API → response headers → HTMX, but fix is only 3 lines |
| Workshop suitability | **3/5** | Relatable but too simple for 2.5h; low TDD potential |

### Bug #10: Fix markdown preview $$ support

- **PR**: [#31514](https://github.com/go-gitea/gitea/pull/31514) | **Issue**: [#31481](https://github.com/go-gitea/gitea/issues/31481)
- **Files**: 6 (1 new) | **Lines**: +79/-6
- **Bug**: `$$A + B$$ test` — text after inline `$$` block gets swallowed. Block parser greedily matched `$$...$$` even when followed by text, and no inline `$$` parser existed.
- **Fix approach**: Added guard in `block_parser.go` to reject `$$...$$text` (not a true block). Created new `InlineBlock` AST node type, new `defaultDualDollarParser` with higher priority (502) than single-dollar (503), renderer adds `display` class for inline-block math.

| Criterion | Score | Notes |
|-----------|-------|-------|
| Domain complexity | 3/5 | Markdown parser internals — goldmark AST/visitor pattern, moderate learning curve |
| Reproducibility | 5/5 | Type `$$a$$ test` in markdown preview, observe text disappears |
| TDD potential | 5/5 | Perfect: test cases are exact string pairs (`"$$a$$ test"` → expected HTML) |
| Systematic-debugging fit | 4/5 | Trace from markdown input → block parser → inline parser → renderer |
| Workshop suitability | **4/5** | Excellent TDD story, but parser domain may need scaffolding for newcomers |

### Bug #17: Fix bleve fuzziness search

- **PR**: [#33078](https://github.com/go-gitea/gitea/pull/33078) | **Issue**: [#31565](https://github.com/go-gitea/gitea/issues/31565)
- **Files**: 11 (1 new) | **Lines**: +83/-52
- **Bug**: Bleve search engine's fuzzy matching causes performance problems — fuzziness too aggressive by default.
- **Fix approach**: Added configurable `TYPE_BLEVE_MAX_FUZZINESS` setting (default 0 = disabled). Extracted `PrepareCodeSearch()` helper to deduplicate 3 identical search handlers. Also fixed indexer queue to not re-queue failed items.

| Criterion | Score | Notes |
|-----------|-------|-------|
| Domain complexity | 3/5 | Full-text search / Bleve — moderately specialized |
| Reproducibility | 3/5 | Needs indexer setup, performance issue not always visible |
| TDD potential | 3/5 | Existing fuzziness tests can be extended, but core fix is a config change |
| Systematic-debugging fit | 3/5 | More of a performance/config fix than classic logic bug; also bundled with refactoring |
| Workshop suitability | **3/5** | Too many files, fix is small but surrounded by refactoring; harder to isolate as "bug" exercise |

### Bug #1: Fix incomplete Actions status aggregations

- **PR**: [#32859](https://github.com/go-gitea/gitea/pull/32859) | **Issue**: [#32857](https://github.com/go-gitea/gitea/issues/32857)
- **Files**: 6 (1 new test) | **Lines**: +106/-38
- **Bug**: `AggregateJobStatus()` only handled 4 states (Failure, Success, Waiting, Running) but ignored Cancelled, Blocked, and Skipped. For example, all jobs skipped → showed "Success" instead of "Skipped"; cancelled jobs → showed "Failure" instead of "Cancelled".
- **Root cause**: Boolean-based logic (`allDone`, `allWaiting`, `hasFailure`) couldn't represent all status combinations. Cancelled was lumped with Failure, Skipped/Blocked were invisible.
- **Fix approach**:
  1. **Core** (`models/actions/run_job.go`): Replaced with flag-based `switch` — `allSuccessOrSkipped`, `hasFailure`, `hasCancelled`, `hasSkipped`, `hasWaiting`, `hasRunning`, `hasBlocked`. Priority: Success > Failure > Running > Waiting > Blocked > Cancelled > Skipped.
  2. **Test** (NEW: `run_job_status_test.go`, 64 lines): Comprehensive table-driven tests covering all status combinations (21 cases).
  3. **Frontend** (`ActionRunStatus.vue`, `status.tmpl`): Separated `cancelled` from `failure` with its own icon (octicon-stop vs octicon-x-circle-fill). Minor CSS cleanup.

| Criterion | Score | Notes |
|-----------|-------|-------|
| Domain complexity | 3/5 | CI/CD job status state machine — concepts familiar, but aggregation priority semantics need understanding |
| Reproducibility | 4/5 | Issue provides demo repo. Create workflow with all-skipped jobs → observe wrong aggregate status |
| TDD potential | 5/5 | Perfect — PR itself was written test-first: new 64-line test file with 21 table-driven cases |
| Systematic-debugging fit | 5/5 | Classic state machine bug: observe wrong status → trace aggregation function → enumerate all states → find gaps → design priority |
| Workshop suitability | **5/5** | State machine bug is a universal pattern; Go backend core is self-contained; frontend changes are separable and cosmetic |

---

## W2 Recommendation

| Priority | Bug | Role | Rationale |
|----------|-----|------|-----------|
| **Primary** | **#1 Fix Actions status aggregation** | Main exercise | State machine bug — universal pattern, perfect TDD (21 test cases), clean Go-only core fix, excellent systematic-debugging (enumerate states → find gaps → design priority). Frontend changes separable. |
| **Primary** | **#20 Fix issue comment number** | Alternative exercise | Counting logic bug, clear reproduction, excellent TDD. More "backend data" flavor vs #1's "state machine" flavor. Choose based on audience background. |
| **Alternative** | **#10 Fix markdown $$ support** | Advanced option | Strong TDD story (string→HTML test cases), but parser internals need more scaffolding for newcomers. Good if audience has parser experience. |
| **Supplement** | **#2 Star/Watch inconsistency** | Warm-up / demo | Quick reproduction for demonstrating systematic-debugging flow at session start. Too simple for full 2.5h. |
| **Not recommended** | **#17 Bleve fuzziness** | Skip | More config/refactor than bug; 11 files too large; hard to isolate the "bug" from surrounding changes. |
