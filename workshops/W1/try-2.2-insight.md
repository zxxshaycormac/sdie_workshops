# Try 2.2 — Insight: The Regression Guard That Wasn't

Re-scoring the local tag-search implementation after `45036094a3` (test: make PaginationWithKeyword a true regression guard) landed. This commit is the only delta from [`try-2.1-insight.md`](./try-2.1-insight.md) — one test file, +14/−9 — but it exposes something the score card cannot see.

## TL;DR

| # | Dimension | Local (2.2) | PR #32045 | Δ vs. 2.1 (local) |
|---|---|:---:|:---:|---|
| 1 | Correctness | 5 | 3 | — |
| 2 | Architecture | 5 | 5 | — |
| 3 | Naming | 4 | 5 | — |
| 4 | Error handling | 5 | 5 | — |
| 5 | Context propagation | 5 | 5 | — |
| 6 | Logging | N/A | N/A | — |
| 7 | i18n | 5 | 5 | — |
| 8 | Testing | 5 | 1 | — (but see §1) |
| 9 | Security | 5 | 5 | — |
| 10 | Performance | 4 | 4 | — |
| 11 | UX polish | 3 | 5 | — |
| 12 | Maintainability | 4 | 5 | −1 (see §3) |
| 13 | Backport / scope | 4 | 4 | −1 local, +1 PR (see §3) |
| | **Average** | **4.17** | **4.33** | — |

PR #32045 edges ahead by ~0.16. The scoring moved against local versus 2.1 (which had local at 4.67 / PR 4.25) — not because the code regressed, but because a second read of the same diff landed on different judgments in three rows. More on that in §3.

## What the new commit actually did

The 2.1 version of `PaginationWithKeyword` asserted pagination *must render*:

```go
// 2.1 (3d2bb650d1):
// "q=v1" with limit=1 — multiple v1* tags match, so filtered count > 1.
req := NewRequest(t, "GET", "/user2/repo1/tags?q=v1&limit=1")
...
assert.Greater(t, pageLinks.Length(), 0,
    "pagination links should appear when filtered count exceeds page size")
```

The 2.2 commit flips both the query and the assertion:

```go
// 2.2 (45036094a3):
// "q=v1.0" with limit=1 — exactly one tag matches, filtered count = 1.
req := NewRequest(t, "GET", "/user2/repo1/tags?q=v1.0&limit=1")
...
assert.Equal(t, 0, pageLinks.Length(),
    "pagination must not render when filtered count fits in one page")
```

The first version is not a regression guard. It would pass on the **fixed** code (filtered count = 2 → pagination renders) AND on the **buggy** code (unfiltered count = 3 → pagination renders). Both branches produce `pageLinks > 0`. The test is green either way; reintroducing the `NumTags` bug would not flip it red.

The 2.2 version is a true guard. With the fix, filtered count = 1 → `TotalPages = 1` → paginate template suppresses the widget → `pageLinks == 0` → assertion holds. With the bug, unfiltered count = 3 → `TotalPages = 3` → widget renders → `pageLinks > 0` → assertion fails. The test now actually distinguishes the two implementations.

## Insights

### 1. A green test that passes on the bug is worth zero regression protection

The score card's Testing row gave local a 5 in both 2.1 and 2.2. That score did not move, because the row counts *presence and coverage shape* — subtests exist, contract is encoded, edge paths asserted. It cannot see whether any individual assertion actually distinguishes fixed from buggy behavior.

The 2.1 → 2.2 delta is the clearest example: the test file is longer in 2.1 (the `Greater` form is one line shorter than the `Equal` form, but the comment block grew). Both versions look like regression tests. One of them is decorative.

**Lesson:** before claiming a test guards behavior X, mentally reintroduce the bug and check whether the assertion flips red. If it doesn't, the test is documentation, not a guard. The 2.2 commit's comment block makes this explicit — it walks through both the fixed and buggy code paths and shows why only the fixed path satisfies `pageLinks == 0`. That comment is the artifact of doing the mental reintroduction. Every regression test should have one.

### 2. The test-quality axis is binary; the score card is continuous

There is no meaningful sense in which the 2.1 test is "70% as protective" as the 2.2 test. It is either a guard or it isn't. A bug either flips the assertion red or it doesn't. The score card's 1-3-5 anchors don't have a slot for "test exists, looks right, encodes the contract, but is operationally inert."

This is a real limitation. The Testing row is the single highest-leverage dimension in the rubric — a 4-point swing between local and PR — and the difference between a 5 and a 1 on that row is doing a lot of work. If "5" can be awarded to a test that wouldn't catch the very bug it was written to prevent, the row is overcounting safety.

**Lesson:** treat the Testing row as a necessary-not-sufficient signal. A 5 means "the test surface exists and looks right." Whether the tests actually catch regressions is a separate question that has to be answered by the mental-reintroduction check from §1, not by reading the test file.

### 3. Two careful reads of the same diff produced different scores — that's the point

Three rows flipped between 2.1 and 2.2, even though the production code is identical (the only new commit is a test edit):

- **Security (5/5 in both versions now; I had 4/4 in my first pass this session).** The `builder.Like` quirk — `%` matches everything — exists in both implementations. 2.1 treated it as an accepted quirk (parameterised query = no injection); my first pass treated it as a defect worth −1. After sitting with it: the PR even documents the behavior in `tag_tooltip`, so treating it as a bug overstates the issue. 5/5 is the honest score.
- **Maintainability (local 5 → 4).** The `if listOptions.Page <= 0 { listOptions.Page = 1 }` guard is defensible but speculative — nothing in the original code required it, nothing in the PR adds it, and it expands the diff beyond the feature. 2.1 credited it as "fell out of writing the pagination test"; 2.2 reads it as scope creep. Both readings are defensible.
- **Backport / scope (local 5 → 4, PR 3 → 4).** 2.1 dinged the PR for restructuring the template (moving `{{if .Releases}}`, adding `TagCount` to the header) beyond what the issue required. On re-read, those template changes are exactly what creates the UX win (empty state, count display) — they're in scope for *the feature the PR ships*, even if they're out of scope for *the minimal feature*. Both implementations land at 4: local for the `Page <= 0` drive-by, PR for the template restructuring.

**Lesson:** the score card is a structured way to find divergences, not a measurement instrument. A 0.5-point swing between two reviewers on the same diff is signal about the diff (it has judgment-call edges), not noise about the reviewers. When the gap between two implementations is smaller than the inter-reviewer variance, the rubric is telling you they're effectively tied.

### 4. The naming and UX gaps from 2.1 are still open

The 2.2 commit doubled down on Testing; it did not touch the two dimensions where PR #32045 is genuinely ahead:

- **Naming (row 3).** `Keyword string` vs. `NamePattern optional.Option[string]`. The field semantically *is* a LIKE pattern — the `strings.ToLower` + `lower_tag_name` machinery only makes sense for pattern matching. `Keyword` reads as free text and hides the wildcard semantics. `FormString` vs. `FormTrim` is the same neighborhood-convention point 2.0 made.
- **UX polish (row 11).** `?q=zzz-not-a-real-tag` still renders an empty table segment with no "no results" message. Four lines of template would close it.

These are the changes a 2.3 commit would make. None of them touch the query path, the test surface, or the regression guard.

## Net

The 2.2 commit is the most important single commit in the series — not because it changed code, but because it examined a test that looked like a guard and fixed the assertion that wasn't. Without it, the local implementation's headline advantage (Testing = 5 vs. PR's 1) would have been partly fictional.

With the fix, local's Testing row is honestly earned. But the score card still puts PR #32045 slightly ahead (4.33 vs. 4.17), because Naming and UX polish — the two dimensions where the PR is genuinely better — outweigh the testing gap in the rubric's point allocation. Whether that weighting is right is a separate question. The score card is a tool for finding where two implementations diverge; it doesn't tell you which divergence matters most. That's still a human call.

The arc across 2.0 → 2.1 → 2.2: in 2.0, local lost on dimensions that hurt users (Correctness, UX). In 2.1, those got fixed. In 2.2, the test that claimed to guard the 2.0 fix got audited and tightened. The next iteration, if there is one, is back to dimensions that hurt reviewers and users (Naming, UX) — the cosmetic gap that's been visible since 2.0 and still isn't closed.
