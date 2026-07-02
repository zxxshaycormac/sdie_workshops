# W1 Try 5.0 — Insight: Tag Page Search (Ours vs Upstream PR #32045)

Compares the most recent change on `training/w1-try-5` (commits since `56e1dd0ae5`) with the upstream PR [go-gitea/gitea#32045](https://github.com/go-gitea/gitea/pull/32045) ("Included tag search capabilities", bsofiato), scored against [`workshops/score-card.md`](../score-card.md).

## Scope of each change

**Ours** — 2 code commits, +634 / −3 across 6 code/locale files + 2 workshop docs:

- `models/repo/release.go`: new `Keyword string` field on `FindReleasesOptions`; `cond.And(builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)})` when non-empty.
- `routers/web/repo/release.go`: bind `q` via `ctx.FormString`; swap `db.Find` → `db.FindAndCount`; pager now uses filtered `count`.
- `templates/repo/tag/list.tmpl`: insert `shared/search/combo` form above the existing `{{if .Releases}}` block — no other template surgery.
- `models/repo/release_test.go`: new `TestFindReleasesByKeyword` (5 scenarios).
- `tests/integration/release_test.go`: new `TestViewTagsListSearch` (4 scenarios).
- `options/locale/locale_en-US.ini`: `tag_kind = Search tags...`.

**PR #32045** — 1 commit, +33 / −7 across 4 files:

- `models/repo/release.go`: new `NamePattern optional.Option[string]`; same `lower_tag_name` LIKE.
- `routers/web/repo/release.go`: bind via `ctx.FormTrim("q")`; `db.Find` then separate `db.Count`; sets `Keyword` and `TagCount`.
- `templates/repo/tag/list.tmpl`: restructured — search form with tooltip; header now reads `{{.TagCount}} tags` (drops the octicon); explicit `no_results_found` branch when search yields nothing but tags exist.
- `options/locale/locale_en-US.ini`: `tag_kind` + `tag_tooltip` (the tooltip advertises `%` wildcard semantics).
- **No tests.**

## Scorecard

| # | Dimension | Ours | PR #32045 | Evidence |
|---|-----------|:----:|:---------:|----------|
| 1 | Correctness | **5** | 4 | Ours: `FindAndCount` keeps find/count atomic; unit + integration tests cover empty / uppercase / no-match / regression. Theirs: also correct, but `no_results_found` only renders `{{if .NumTags}}` — borderline when a repo genuinely has 0 tags. |
| 2 | Architecture | **5** | **5** | Both: option struct in `models/`, handler binds and dispatches, template renders. Dependencies flow downward cleanly. |
| 3 | Naming | 4 | **5** | `Keyword` slightly misleads — it's a LIKE pattern, not a fuzzy keyword. `NamePattern` is honest. |
| 4 | Error handling | **5** | **5** | Both check the error and route through `ctx.ServerError`. |
| 5 | Context propagation | **5** | **5** | `ctx` first parameter everywhere it matters; both pass it into `db.Find` / `db.Count`. |
| 6 | Logging | N/A | N/A | Neither adds log lines. |
| 7 | i18n | **5** | 4 | Both add English source first. Theirs' tooltip copy "match any sequence of **numbers**" is factually wrong — `%` matches any sequence of *characters*, and would have to be re-translated after correction. |
| 8 | Testing | **5** | 1 | Ours: 9 scenarios across unit + integration. Theirs: zero tests for a search feature. |
| 9 | Security | 4 | 4 | Both parameterize via `builder.Like`. Neither escapes `%` / `_` — same minor gap, treated as documented wildcard behavior. |
| 10 | Performance | 4 | 4 | Both issue one find + one count. No new index; relies on the existing `lower_tag_name` column. |
| 11 | UX polish | 3 | **4** | Ours: no empty state when search returns nothing — user sees only the search box. Theirs: explicit `no_results_found` message + tooltip, but loses the tag icon in the header. |
| 12 | Maintainability | **5** | 3 | Ours: plain `string` field fits the surrounding options; tests and workshop docs anchor it. Theirs: `optional.Option[string]` for a never-required search input is ceremony; no tests or docs. |
| 13 | Backport / scope | 4 | 4 | Ours: focused code diff, but two `workshops/*.md` artifacts ride along — fine in training, would need to be split for a real PR. Theirs: small drive-by header rewrite (octicon → count) expands the diff. |
| — | **Average (1–5)** | **4.67** | **4.00** | Ours: 56 / 60 across 12 applicable dims (logging N/A). Theirs: 48 / 60. |

### Score summary

- **Ours:** 56 / 60 → **4.67 / 5** across 12 applicable dimensions (Logging N/A).
- **PR #32045:** 48 / 60 → **4.00 / 5** across 12 applicable dimensions (Logging N/A).
- **Per-row deltas favoring ours:** Correctness (+1), i18n (+1), Testing (+4), Maintainability (+2).
- **Per-row deltas favoring theirs:** Naming (+1), UX polish (+1).
- **Ties:** Architecture, Error handling, Context propagation, Security, Performance, Backport/scope.

## Insights

1. **Where we clearly win: testing and i18n accuracy.** The unit + integration tests are the single biggest gap in the upstream PR — a search feature shipped with zero tests is exactly the kind of change that regresses silently. Our i18n is also tighter: their `tag_tooltip` copy is factually wrong about `%` semantics and would have to be re-translated after correction.

2. **Where they clearly win: empty-state UX.** When `q=nonexistent` returns nothing on a repo that *does* have tags, ours renders a lone search box with no message. Theirs shows `no_results_found`. This is the one feature gap worth pulling from their version — a small template edit (`{{else}}{{if .NumTags}}<p>{{ctx.Locale.Tr "no_results_found"}}</p>{{end}}{{end}}` inside the table segment).

3. **API shape is a genuine design choice, not a clear win.** `Keyword string` (ours) vs `NamePattern optional.Option[string]` (theirs). Ours is simpler and matches the `ListOptions` style; theirs is more semantically honest. Either is defensible — keep ours but consider renaming `Keyword` → `NamePattern` (plain `string`) to get the best of both.

4. **One real bug to flag in ours.** `FindAndCount` correctly runs the count with the same `opts`, but `release_tag_header.tmpl:10` still renders `{{.NumTags}} tags` from the *unfiltered* count in the tab. So searching for "v1" shows "2 tags" in the body (via pager) and "3 tags" in the header tab simultaneously. Theirs has the inverse inconsistency (header shows filtered, tab shows total). Neither is wrong, but the inconsistency should be a conscious decision, not an accident.

5. **Process artifacts in scope (row 13).** The `workshops/*-design.md` and `*-plan.md` files are committed alongside the code. The rubric explicitly excludes spec artifacts from being their own dimension, so this does not hurt the score — but for a real upstream PR they should live in a separate commit or out of tree.

## Bottom line

Ours is the stronger change, largely on the back of tests and i18n correctness. Theirs has one UX detail (the empty state) worth borrowing — and a slightly more honest field name.
