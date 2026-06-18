# Try 2.1 — Summary

Local vs. PR #32045 after the pagination fix: **4.67 vs 4.25**. Gap flipped.

**What landed:** `3d7acb5aac` switched `db.Find` → `db.FindAndCount` and built the pager from the filtered count. `3d2bb650d1` added `CaseInsensitive` and `PaginationWithKeyword` tests — the second one is the regression guard that would have caught the 2.0 bug.

**Where local now wins:**
- **Testing (5 vs 1):** PR #32045 ships zero tests. Local pins the contract.
- **Backport/scope (5 vs 3):** PR restructures the template beyond the issue. Local stays surgical.
- **Correctness (5 vs 4):** local adds a `Page <= 0 → 1` guard the PR doesn't have.

**Where PR #32045 still wins:**
- **Naming (5 vs 4):** `NamePattern optional.Option[string]` documents the LIKE-pattern semantics; `Keyword` hides it. Also: `FormTrim` is the sibling-handler convention; local uses `FormString`.
- **UX polish (5 vs 3):** PR adds an explicit `no_results_found` empty state and shows filtered count in the header. Local still renders a silent empty table when `?q=` matches nothing.

**Insight:** the gap moved from "does it work for the user" (Correctness/UX) to "is it nice to review" (Naming/UX cosmetics). That's the right direction — contract-level behavior is now equal-or-better than upstream.

**Path to 2.2 (3 commits, all cosmetic):** rename `Keyword` → `NamePattern`, switch `FormString` → `FormTrim`, add the 4-line empty-state branch. After that, local is strictly better on every non-N/A dimension.

Full version: [`try-2.1-insight.md`](./try-2.1-insight.md)
