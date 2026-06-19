# Try-3 Insight: Local Implementation vs Upstream PR #32045

Comparison of the W1 Tag-search exercise (local `training/w1-try-3` branch, 7 commits since `689ad86`) against the actual upstream fix at https://github.com/go-gitea/gitea/pull/32045 (merged 2024-09-17).

Both implement the same feature (#31998): keyword search on the repository Tags page.

> **Revision note.** This is the second pass. The first pass scored local Correctness at 5 — that was wrong. Local filters on `tag_name` (case-sensitive on PostgreSQL), but every test runs on SQLite where LIKE is case-insensitive by default. All 9 subtests pass while the feature is broken on PostgreSQL for case-mismatched queries. Correctness is downgraded to 4. PR average was also miscounted (3.9 → actually 4.25).

## Scope compared

| | Local try-3 | PR #32045 |
|---|---|---|
| Commits | 7 (5 functional + 2 docs) | 1 |
| Files changed | 11 (incl. openspec docs) | 4 |
| LOC added | ~1020 (incl. 669-line plan.md) | 33 |
| LOC deleted | 3 | 7 |

## Score Card

| # | Dimension | Local | PR #32045 | Divergence |
|---|---|:---:|:---:|---|
| 1 | Correctness | **4** | **3** | Local: PostgreSQL case-mismatch bug hidden by SQLite tests. PR: happy path works in prod, edge cases unverified by tests |
| 2 | Architecture | 5 | 5 | Identical layering |
| 3 | Naming | **4** | **5** | PR's `NamePattern optional.Option[string]` matches siblings; local's `Keyword string` is the only plain-string field on this struct |
| 4 | Error handling | **5** | **4** | Local: `ServerError("CountReleases", err)` — accurate. PR: `ServerError("GetReleasesByRepoID", err)` — misleading (wraps the Count call) |
| 5 | Context propagation | 5 | 5 | Both thread ctx correctly |
| 6 | Logging | N/A | N/A | Neither adds logs |
| 7 | i18n | 4 | **5** | PR adds `tag_tooltip` documenting `%` wildcard; local adds only `tag_kind` |
| 8 | Testing | **5** | **1** | Local: 64-line unit test (4 subtests) + 74-line integration test (5 subtests). PR: zero tests |
| 9 | Security | 5 | 5 | Both use parameterized `builder.Like` |
| 10 | Performance | 4 | 4 | Both add one extra `db.Count` call. `tag_name` is indexed and `lower_tag_name` is not, but substring LIKE (`%foo%`) can't use a B-tree index anyway, so neither benefits |
| 11 | UX polish | **3** | **5** | PR distinguishes "repo has no tags" vs "search matched nothing" (`{{if .NumTags}}`), shows filtered count in header, exposes `%` wildcard via tooltip. Local conflates the two empty states and shows no count |
| 12 | Maintainability | 5 | 4 | Local's comment on `filteredCount` explains WHY; PR's `NumTags` branch could use one |
| 13 | Backport / scope | 5 | 5 | Both minimal, cherry-pickable |
|   | **Average** | **4.50** | **4.25** | |

## Insights

### 1. Tests give false confidence — the biggest meta-issue

Local has 138 lines of tests across 9 subtests. They all pass. The feature is still broken on PostgreSQL. The reason: every test runs against SQLite, where `LIKE` is case-insensitive by default. The contract the tests assert ("substring `v1` returns `v1.0` and `v1.1`") holds on SQLite and MySQL (case-insensitive collations) but fails on PostgreSQL (case-sensitive LIKE).

**Lesson.** Test coverage is necessary but not sufficient. A test suite that runs only on one DB backend cannot validate DB-backend-specific behavior. Gitea ships on MySQL, PostgreSQL, and SQLite — the local test matrix covers one of three. This is invisible from the test count alone.

The PR's zero tests is a different failure mode (no verification at all), but at least it doesn't pretend to verify what it doesn't.

### 2. Local's correctness bug is real and specific

`builder.Like{"tag_name", opts.Keyword}` on PostgreSQL with case-mismatched input returns zero rows. A user typing `V1` against tags stored as `v1.0`, `v1.1` gets nothing. The fix is mechanical: filter on `lower_tag_name` with `strings.ToLower(opts.Keyword)`, matching the upstream PR and the existing `GetReleaseByTag` / `release.go:154,188` pattern which already lower-cases tag lookups.

### 3. PR wins UX, local wins verification (after correction)

PR's `{{if .NumTags}}...<p>no_results_found</p>...{{end}}` correctly tells the user "you searched for something that doesn't exist" *without* lying when the repo genuinely has no tags. Local always shows the no-results string, which is wrong for empty repos. PR's tooltip (`%` wildcard) also exposes LIKE semantics power users would want.

### 4. Naming consistency matters more than it looks

`FindReleasesOptions` is full of `optional.Option[T]` fields (`IsDraft`, `IsPreRelease`, `HasSha1`, `TagNames`). PR's `NamePattern optional.Option[string]` fits the family. Local's `Keyword string` is the only plain-string sibling and forces a `!= ""` guard in `ToConds` instead of the `.Has()` check used by every other field in the same function.

### 5. The 669-line `plan.md` is invisible to the score card — by design

Item 160 of `score-card.md` explicitly excludes process artifacts because human contributors don't produce them. The artifacts helped *produce* better code (visible in the comment discipline and test coverage), but they aren't themselves a scored dimension. They also didn't prevent the case-sensitivity bug — the plan can describe the right intent and the implementation can still diverge. Plan fidelity is its own quality gate, separate from code quality.

### 6. Two failure modes of "no tests" vs "tests that pass on the wrong DB"

| | Local | PR |
|---|---|---|
| Tests present | Yes (138 lines) | No |
| Tests catch the case bug | No (SQLite hides it) | N/A |
| Production behavior | Broken on PostgreSQL | Correct everywhere |
| Risk profile | False confidence — looks safe, isn't | Honest unknown — looks unverified, is |

A reviewer applying only the score card would correctly identify the PR as under-tested (dim 8 = 1) but might miss that local is also correctness-broken (dim 1) unless they read the implementation closely. The score card rewards test *presence*; it cannot reward test *fidelity*.

## Action items to bring local to upstream parity

- Switch `Keyword` → `NamePattern optional.Option[string]`; filter on `lower_tag_name` with lowercased input. Closes the PostgreSQL bug and matches the surrounding idiom.
- Add a `search.tag_tooltip` locale key.
- Replace the unconditional empty state with the `{{if .NumTags}}` guard so empty repos don't show "no results".
- Keep the accurate `ServerError("CountReleases", err)` (already correct in local; the upstream PR's `GetReleasesByRepoID` label is misleading).
- Keep the test coverage — but add a comment in the integration test acknowledging the SQLite-only fidelity gap, and consider asserting the lowercased-input path explicitly so a future Postgres-CI run would surface the bug.

## What the first pass got wrong

| Item | First pass | Corrected | Why |
|---|---|---|---|
| Local Correctness | 5 | 4 | Case-sensitivity bug on PostgreSQL |
| PR average | 3.9 | 4.25 | Arithmetic error (51/12) |
| Margin (local − PR) | 0.6 | 0.25 | Local still wins, but narrower |

The relative ordering is unchanged, but the gap is roughly half what the first pass implied. The PR is closer to local than it first appeared once you account for the case-sensitivity bug dragging local down.
