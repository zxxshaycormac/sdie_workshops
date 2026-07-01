# W1 Try-4 Insight: Tag Search vs Upstream PR #32045

Comparison of the local tag-search implementation (commits `56e1dd0..HEAD` on `training/w1-try-4`) against the upstream Gitea fix [PR #32045](https://github.com/go-gitea/gitea/pull/32045) (resolves #31998), scored with [`workshops/score-card.md`](../score-card.md).

## Scope note

Commit `56e1dd0` itself is the asdf/Go-pinning change and is unrelated to tag search. The actual W1 tag-search work lives in the 5 commits *after* it. Both implementations touch the same 4 files (`models/repo/release.go`, `routers/web/repo/release.go`, `templates/repo/tag/list.tmpl`, `options/locale/locale_en-US.ini`); the local branch additionally adds 2 test files.

## Side-by-side scorecard

| # | Dimension | Upstream #32045 | Local (try-4) | Divergence |
|---|---|:---:|:---:|---|
| 1 | Correctness | **4** | **3** | Upstream searches `lower_tag_name` (case-insensitive); local searches `tag_name` (case-sensitive — `V1` misses `v1.0`). Upstream also `FormTrim`s; local uses `FormString`. |
| 2 | Architecture & layering | 5 | 5 | Both clean — router binds, model expresses `builder.Cond`. |
| 3 | Naming | **4** | **4** | Upstream `NamePattern optional.Option[string]` matches sibling fields; local `Keyword string` is plain. Local adds a WHY-comment; upstream adds none. |
| 4 | Error handling | **3** | **4** | Upstream reuses `ctx.ServerError("GetReleasesByRepoID", err)` for the *count* call — misleading copy-paste. Local uses `"CountReleases"`. |
| 5 | Context propagation | 5 | 5 | Both thread `ctx` to every DB call. |
| 6 | Logging | N/A | N/A | Neither adds log lines. |
| 7 | i18n | 5 | 5 | Both add `search.tag_kind`. Local adds a parametrized `release.tags.no_match`; upstream adds `tag_tooltip` and reuses generic `no_results_found`. |
| 8 | Testing | **1** | **5** | Upstream ships zero tests. Local adds a unit test (`TestFindReleasesOptions_KeywordFilter`) + an integration test with 3 sub-cases (match / no-match / absent-q) and tag cleanup. |
| 9 | Security | 5 | 4 | Both parameterize via `builder.Like`. Upstream *documents* the `%` wildcard in a tooltip; local leaves the wildcard-into-`LIKE` behavior undocumented. |
| 10 | Performance | **3** | **5** | Upstream always issues a fresh `COUNT(*)` even when `NumTags` is already known. Local only counts when a keyword is set, with a comment explaining why. |
| 11 | UX polish | **4** | **4** | Upstream: tooltip + count-in-header + conditional empty state. Local: better empty-state message (`No tags match "foo"`) but no tooltip, no count-in-header. |
| 12 | Maintainability | 4 | 4 | Upstream more consistent with the `optional.Option` idiom; local has better WHY-comments. |
| 13 | Backport / scope | **4** | **4** | Upstream sneaks in a header refactor (`Tags` -> `<n> Tags`); local's template diff is smaller. |

### Score summary

- Applicable dimensions: 12 (Logging is N/A for both).
- **Upstream PR #32045**: sum = 47 / 60, **average = 3.92** on a 1-5 scale.
- **Local try-4**: sum = 52 / 60, **average = 4.33** on a 1-5 scale.
- Per-row delta (local − upstream): `+0, +0, +0, +1, +0, +4, −1, +2, +0, +0, +0, +0` = net **+4** across 12 rows, with the biggest swings on Testing (`+4`) and Performance (`+2`), offset by Correctness (`−1`) and Security (`−1`).

## The three insights that matter

**1. Case-insensitivity is the one true correctness bug.** The `lower_tag_name` column exists precisely for this (`models/repo/release.go:457,464` already use it). Searching `tag_name` directly means uppercase queries silently return nothing. The local tests pass only because the fixtures happen to use lowercase tag names — the test does not catch the gap because the test data does not exercise it. A `"V1"` search case in `TestFindReleasesOptions_KeywordFilter` would have exposed this in red.

**2. Local beats upstream on two dimensions they actually care about.**
- *Performance*: upstream's unconditional `db.Count[Release]` is a redundant round trip on every tags-page load — `NumTags` is already middleware-set. The local conditional fallback is the right call, and the comment makes the intent reviewer-proof.
- *Testing*: upstream merged with zero tests. The local integration test (`TestTagsListSearch`) is the kind of contract test `design.md` §7 asks for — substring match, empty state, absent `q`, hermetic cleanup.

**3. The `optional.Option` vs plain `string` choice is a wash — but `FormTrim` isn't.** Every other `q` search parameter in `routers/web/repo/` uses `ctx.FormTrim` (`packages.go:28`, `commit.go:188`). The local `FormString` is the odd one out and lets a pasted `" v1 "` miss matches. Cheap fix, convention-aligned.

## Suggested follow-ups (in priority order)

1. Switch `builder.Like{"tag_name", ...}` to `builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)}` and add a mixed-case test case.
2. `ctx.FormString("q")` to `ctx.FormTrim("q")`.
3. Consider renaming `Keyword` to `NamePattern` and typing it `optional.Option[string]` to match `IsPreRelease`/`IsDraft`/`HasSha1`. Or keep `Keyword` and add a one-line comment explaining why plain string is intentional.
4. Optional: steal upstream's `<count> Tags` header — it is a real UX win and a one-line template change.

## Net

Try-4 is arguably the stronger engineering artifact (tests + perf), but upstream got the case-insensitivity right that try-4 missed. Both fixes are small.
