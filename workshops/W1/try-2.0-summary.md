# Try 2.0 — Summary

Local vs. PR #32045 on the rubric: **3.5 vs 4.2**.

**Where local lost:** Correctness (2 vs 5) and UX polish (2 vs 5). Same root cause.

**The bug:** `TagsList` built the pager from `ctx.Data["NumTags"]` — the *unfiltered* repo total — while the query carried `Keyword`. Page 1 rendered correctly; pages 2+ were empty when filtering. The integration test only asserted on page 1, so it passed.

**Same root cause in the template:** the inherited `{{if .Releases}}` guard hid the entire section when a filter matched nothing. User saw a search box with no feedback.

**Lessons:**
1. A failing test must assert the behavior you care about, not the happy path. Green tests ≠ correct behavior.
2. Conventions come from the sibling files, not the file you're editing. `FormString` vs `FormTrim`, `Keyword` vs `optional.Option[T]` — both wrong against neighborhood style.
3. Adding a filter adds an empty-state path. Someone has to design it, or the template's existing guards will hide the absence.

**Fix path:** swap `db.Find` → `db.FindAndCount`, use the returned count for the pager, add an empty-state branch. That's what 2.1 does.

Full version: [`try-2.0-insight.md`](./try-2.0-insight.md)
