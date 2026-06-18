# Try 2.0 — Insight: Comparing Against PR #32045

Comparing the local tag-search implementation (`825d736b77..HEAD`, 5 commits) against the upstream fix (go-gitea/gitea#32045), scored with `workshops/score-card.md`.

## TL;DR

| | Local | PR #32045 |
|---|:---:|:---:|
| Correctness | 2 | 5 |
| Architecture | 5 | 5 |
| Naming | 3 | 4 |
| Error handling | 5 | 4 |
| Context propagation | 5 | 5 |
| Logging | N/A | N/A |
| i18n | 4 | 5 |
| Testing | 4 | 2 |
| Security | 5 | 5 |
| Performance | 3 | 4 |
| UX polish | 2 | 5 |
| Maintainability | 4 | 4 |
| Backport / scope | 5 | 4 |
| **Average** | **3.5** | **4.2** |

The gap lives almost entirely in two rows: **Correctness** and **UX polish**. Both trace to the same root cause.

## The pivotal difference

`NumTags` is set in `services/context/repo.go:517` as an **unfiltered** count of all repo tags. The local handler still uses it for pagination:

```go
numTags := ctx.Data["NumTags"].(int64)
pager := context.NewPagination(int(numTags), opts.PageSize, opts.Page, 5)
```

When filtering, the *results* are correct (the query carries `Keyword`), but the *pagination widget* is built from the total count. Pages 2+ render empty when fewer tags match. PR #32045 fixes this with a dedicated filtered count:

```go
count, err := db.Count[repo_model.Release](ctx, opts)  // opts includes NamePattern
pager := context.NewPagination(int(count), opts.PageSize, opts.Page, 5)
```

The local integration test only hits page 1, so the bug slips through.

The same root cause shows up again in the template: the local version keeps the original `{{if .Releases}}` guard around the header+table, so when a filter returns nothing the user sees a floating search box with no feedback. PR #32045 restructures so the header always shows (with the filtered count) and adds an explicit `no_results_found` empty state.

## Insights

### 1. TDD caught the wrong bug

The local workflow wrote the test first (`2cd12a1ed9 test: add failing integration test`) — good discipline. But the test asserted only on the first-page result set. A pagination bug is invisible to it. A second assertion — "when filtering reduces results, the pager reflects the filtered count, not the repo total" — would have forced the `Count` query into existence.

**Lesson:** the failing test must assert the behavior you actually care about, not just the happy path. A green test suite is only as strong as the behaviors it encodes.

### 2. Conventions are established by the neighborhood, not the file

Two naming deviations in the local version are invisible if you only read `release.go`:

- `ctx.FormString("q")` — but `FormTrim` is the dominant pattern across 7 of 9 sibling handlers (`packages.go`, `commit.go`, `search.go`, `projects.go`, `issue.go` ×3).
- `Keyword string` — but every other filter field on the same struct uses `optional.Option[T]` (`IsDraft`, `HasSha1`).

PR #32045 matches both conventions.

**Lesson:** before adding an identifier, scan the sibling files and the sibling fields. The file you're editing is the weakest signal of local style.

### 3. Any new filter creates a new empty state

The PR's most user-visible improvement isn't the search box — it's the *always-visible header with count* and the explicit `no_results_found` paragraph. Both come from treating "filter returns nothing" as a first-class UI state.

The local template inherited the original `{{if .Releases}}` guard, which silently hid the entire section when the filter matched zero tags. The user is left staring at a search box and a pager with no explanation.

**Lesson:** adding a filter doesn't just add a "results" path — it adds an "empty results" path that someone must explicitly design. If you don't, the template's existing guards will quietly hide the absence.

### 4. Where the rubric bites

The score-card's row 1 (Correctness) and row 11 (UX polish) are where the gap is largest — both trace to the same root cause (not re-deriving the count for the filtered option set). Architecture, security, i18n, and naming are all close.

This matches the rubric's intent: the dimensions that catch real bugs are the ones that cost the most when skipped. A reviewer scanning only rows 2, 5, 9 would conclude the implementations are equivalent. They aren't.

### 5. Process artifacts are deliberately not scored — and that's right

The local version carries a 456-line implementation plan and an openspec change folder. The PR carries none. The score-card excludes these (see "What's deliberately not a dimension") because human contributors aren't held to that standard. Good thing — scoring on them would flip the comparison unfairly.

**What matters is the shipped diff.** The plan's quality is measured by the diff it produced, not by its own eloquence. A plan that misses the pagination count produced a diff that misses the pagination count.

## Net

The local implementation is a clean, well-tested **first pass** that misses one operational detail (filtered pagination count) and one UX state (empty results). PR #32045 is a more complete **production fix** with weaker test coverage. On the rubric, the PR wins by ~0.7, almost entirely on rows 1 and 11 — the two rows that encode "does it actually work for the user."

The encouraging part: the local workflow got rows 2, 4, 5, 7, 9, 12, 13 right. The discipline held. What it missed was not a discipline failure but a **requirements failure** — no one (including the plan) asked "what happens to pagination when filtering?" or "what does the user see when zero tags match?". Those are the questions to add to the brainstorming checklist next time.
