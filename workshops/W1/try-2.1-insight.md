# Try 2.1 — Insight: After the Pagination Fix

Re-scoring the local tag-search implementation against PR #32045 after `3d7acb5aac` (filtered count) and `3d2bb650d1` (edge-case tests) landed. These two commits close exactly the two gaps called out in [`try-2.0-insight.md`](./try-2.0-insight.md) — Correctness and UX polish — and flip the comparison.

## TL;DR

| # | Dimension | Local | PR #32045 |
|---|---|:---:|:---:|
| 1 | Correctness | **5** | 4 |
| 2 | Architecture | 5 | 5 |
| 3 | Naming | 4 | **5** |
| 4 | Error handling | 5 | 5 |
| 5 | Context propagation | 5 | 5 |
| 6 | Logging | N/A | N/A |
| 7 | i18n | 5 | **5** |
| 8 | Testing | **5** | 1 |
| 9 | Security | 5 | 5 |
| 10 | Performance | 4 | 3 |
| 11 | UX polish | 3 | **5** |
| 12 | Maintainability | 5 | 5 |
| 13 | Backport / scope | **5** | 3 |
| | **Average** | **4.67** | **4.25** |

The 0.7 gap from 2.0 is gone — local now leads by 0.42. But the *shape* of the gap changed: it moved out of Correctness (now equal-or-better) and into Testing + Backport/scope (where local is clearly stronger) and Naming + UX polish (where PR #32045 is still ahead).

## What the fix commits actually did

`3d7acb5aac` swaps the pagination source from the unfiltered repo total to the filtered count:

```go
// before (2.0):
releases, err := db.Find[repo_model.Release](ctx, opts)
numTags := ctx.Data["NumTags"].(int64)
pager := context.NewPagination(int(numTags), opts.PageSize, opts.Page, 5)

// after (2.1):
releases, total, err := db.FindAndCount[repo_model.Release](ctx, opts)
pager := context.NewPagination(int(total), opts.PageSize, opts.Page, 5)
```

`FindAndCount` does both the page query and the count in one call, so the pagination widget now reflects whatever the `Keyword` filter narrowed the result set to. PR #32045 solves the same problem but with two separate calls (`db.Find` then `db.Count`); functionally equivalent, slightly less idiomatic.

`3d2bb650d1` adds two integration tests that directly encode the behaviors 2.0 missed:

- `CaseInsensitive` — verifies `?q=V1.1` matches via `lower_tag_name` (proves the `strings.ToLower` machinery is wired through).
- `PaginationWithKeyword` — forces `limit=1` with `?q=v1` and asserts the pager shows multiple pages. This is the regression test for the 2.0 bug; it would have caught the original `numTags` mistake.

## Insights

### 1. The right test would have caught 2.0 — and now it does

The 2.0 insight's main lament was that the failing-test assertion was too weak: it only checked page 1. The fix isn't to write *more* tests, it's to write the *one* test that encodes the actual contract — "pager reflects filtered count." That test now exists in `PaginationWithKeyword`. PR #32045 has zero tests in its diff, so it has no regression guard at all.

**Lesson:** the score-card's Testing row is asymmetric. Going from "some tests" (3) to "tests of the contract" (5) costs one well-chosen assertion. Going from "no tests" (1) to "some tests" (3) costs the same effort but produces strictly less safety. The marginal test is worth more than the first one if it pins behavior the previous tests missed.

### 2. Two real differences remain, and they pull in opposite directions

**Where PR #32045 is better:**

- **Naming (row 3).** `NamePattern optional.Option[string]` vs. `Keyword string`. The field semantically *is* a LIKE pattern — the `lower_tag_name` + `strings.ToLower` machinery proves it. `NamePattern` documents that at the type; `Keyword` reads as "free text" and hides the wildcard semantics. Local also uses `ctx.FormString("q")` where 7 of 9 sibling handlers use `ctx.FormTrim("q")` — minor, but the same neighborhood-convention point 2.0 made.
- **UX polish (row 11).** PR #32045 adds an explicit empty-state branch (`no_results_found`) and shows the filtered count in the header. Local only adds the search form. When `?q=zzz` matches nothing, local renders an empty table with no explanation. This is the *same* UX gap 2.0 had — the fix commits addressed Correctness but not UX. Worth a 4-line template follow-up.

**Where local is better:**

- **Backport / scope (row 13).** PR #32045 restructures the template more than the issue requires (moves the header section, reorders `{{if .Releases}}` around the whole block, adds `TagCount` to the header). Local touches only what's needed. On a backport-focused 1.22.x branch, the smaller diff is the right call.
- **Correctness (row 1).** Local adds a `Page <= 0 → Page = 1` guard in `TagsList` that PR #32045 doesn't have. Not strictly a bug — `db.ListOptions{Page: 0}` returns all rows — but the guard makes the pager widget's arithmetic correct when no `page` param is supplied. A defensible scope expansion that fell out of writing the pagination test.

### 3. `FindAndCount` vs. `Find` + `Count` is not a performance question

Both approaches issue two SQL queries under the hood (XORM's `FindAndCount` runs the count separately). The Performance row (10) scores local 4 / PR 3, but the real difference is call-site cleanliness — one err return vs. two. Don't over-read this row; if Performance were the deciding axis, both implementations would score the same.

### 4. The shape of the gap moved, which is the point of iterating

In 2.0, the gap was on the dimensions that hurt users (Correctness, UX). In 2.1, the gap is on dimensions that hurt reviewers (Naming) and power-users (UX polish), but the contract-level correctness is now equal-or-better than the upstream fix. That's the right direction — the dimensions that encode "does it actually work" are no longer where local loses.

## Net

Local is now ahead on the rubric (4.67 vs. 4.25), with the win concentrated in Testing and Backport/scope. The remaining gaps are cosmetic: rename `Keyword` → `NamePattern`, switch `FormString` → `FormTrim`, add the 4-line empty-state branch. None of them touch the query path or the test surface.

If there's a 2.2, it's a 3-commit cleanup: rename, FormTrim, empty state. After that, the local implementation would be strictly better than PR #32045 on every dimension except the N/A.
