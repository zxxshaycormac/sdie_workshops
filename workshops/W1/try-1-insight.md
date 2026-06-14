# Tags Search Filter — Comparison with Upstream PR #32045

Upstream PR: https://github.com/go-gitea/gitea/pull/32045
Commit: `ec62061fe67956858f480944962bb19b58abef67` ("Included tag search capabilities")

## What's the same

- Filter mechanism: `builder.Like{"lower_tag_name", strings.ToLower(...)}` on the same column.
- HTTP param: `q` via `ctx.FormTrim`.
- UI: reuses `shared/search/combo` template.
- Pagination uses the filtered count, not the total tag count.

## Side-by-side

| Area | Upstream #32045 | Local branch |
|---|---|---|
| Option field | `NamePattern optional.Option[string]` | `Keyword string` |
| Counting | `db.Find` then separate `db.Count` (2 queries) | `db.FindAndCount` (1 query) |
| Locale | adds `tag_kind` + `tag_tooltip` | adds `tag_kind` only |
| Template | restructures: search box -> count header -> empty-state branch | minimal insert; keeps `{{if .Releases}}` wrapping header+table |
| `GetTagList` JSON endpoint | untouched | adds in-memory substring filter |
| Tests | none | unit + integration |
| Spec | none | openspec change tracking |

## Where local is stronger

- **Efficiency.** `FindAndCount` collapses find+count into one DB round-trip. Upstream does two.
- **Type choice.** Plain `string` is the right call here. The `optional.Option[string]` wrapper is redundant because `ToConds` already guards on `== ""`. The optional type earns its keep when the "unset" state must round-trip through serialization or composition; here it is pure ceremony.
- **Testability.** Unit tests in `models/repo/release_test.go:30` and integration tests in `tests/integration/release_test.go:243` lock in the substring + case-insensitivity contract. Upstream shipped blind.
- **Coverage.** Filtering `GetTagList` (the dropdown autocomplete endpoint) too means the search box and the typeahead behave consistently. Upstream left the dropdown unfiltered.

## Where local is weaker

1. **Missing empty state (real gap).** `templates/repo/tag/list.tmpl` keeps the old `{{if .Releases}}` wrapping both the header and the table. When a user searches for `q=zzz` and nothing matches, the entire header+table block disappears, leaving just the search box and pagination. Upstream's restructure (header always rendered with `{{.TagCount}}`, plus a `no_results_found` paragraph inside the table segment) is materially better UX.
2. **No tag count in header.** Upstream shows `{{.TagCount}} Tags` reflecting the filtered set. Local still shows just the icon. Cheap win; also part of the empty-state fix.
3. **Missing tooltip.** Upstream's `tag_tooltip = Search for matching tags. Use '%' to match any sequence of numbers.` (the wording "numbers" is wrong — should be "characters" — but the intent is right). Without it, users won't realize `LIKE` semantics mean `%` is a wildcard, and may not understand why `v1` matches `v1.0` and `v1.1` rather than only an exact `v1`.

## Suggestion

Port the template restructure from upstream — keep the local `FindAndCount`, plain `string`, tests, and `GetTagList` filter. That gets the better empty-state + count display without giving back any of the local improvements. Also consider adding a corrected `tag_tooltip` locale key.

## Score Card (standard rubric, 13 dimensions)

Scored against the reusable rubric in [`workshops/score-card.md`](../score-card.md), which maps every dimension to its source in [`design.md`](../../design.md) or [`CLAUDE.md`](../../CLAUDE.md). Scale 1-5; both implementations scored independently, then diffed. Logging is N/A for both (no log statements added). Process discipline is intentionally excluded — openspec artifacts are an AI-workflow convention, not a codebase rule, so penalizing a human-authored PR for lacking them would be unfair.

| # | Dimension | Source | Upstream | Local | Δ | Evidence |
|---|---|---|:-:|:-:|:-:|---|
| 1 | Correctness | general | 3 | 4 | +1 | Upstream behavior right but unverified; local locks contract with tests. Local loses a point for missing empty-state feedback. |
| 2 | Architecture & layering | design.md §1 | 5 | 5 | 0 | Both: router binds `q`, passes via `FindReleasesOptions`, model adds `builder.Like`. Clean downward flow. |
| 3 | Naming | design.md §2 | 5 | 5 | 0 | Both: PascalCase exports, camelCase locals, snake_case files. `NamePattern` / `Keyword` are both clear. |
| 4 | Error handling | design.md §3 | 4 | 4 | 0 | Both: `ctx.ServerError` at boundary, errors returned up. Small surface area. |
| 5 | Context propagation | design.md §4 | 5 | 5 | 0 | Both: ctx flows handler → `db.Find` / `db.FindAndCount`. |
| 6 | Logging | design.md §5 | N/A | N/A | — | Neither change adds log statements. |
| 7 | i18n | design.md §6 | 4 | 3 | −1 | Upstream adds `tag_kind` + `tag_tooltip` (wording bug: "numbers" → "characters"). Local adds `tag_kind` only. |
| 8 | Testing | design.md §7 | 1 | 5 | +4 | Upstream: zero tests. Local: unit (`models/repo/release_test.go:30`) + integration (`tests/integration/release_test.go:243`). |
| 9 | Security | design.md §8 | 4 | 4 | 0 | Both: `builder.Like` parameterizes; `ctx.FormTrim` at boundary; repo-scoped via `RepoID`. |
| 10 | Performance | general | 3 | 5 | +2 | Upstream: `Find` + `Count` = 2 round trips. Local: `FindAndCount` = 1. |
| 11 | UX polish | general | 5 | 2 | −3 | Upstream: empty state + count in header + tooltip. Local: search box only; header+table vanish together on no match. |
| 12 | Maintainability | general | 3 | 5 | +2 | Upstream: `optional.Option[string]` redundant given `ToConds`'s `== ""` guard. Local: plain `string`. |
| 13 | Backport / scope | CLAUDE.md | 5 | 4 | −1 | Upstream: 4 files, 33+/7-, narrowly scoped. Local: also filters `GetTagList` — defensible extension of #31998 but adds surface for the 1.22.x backport. |
| | **Average** | | **3.92 / 5** | **4.25 / 5** | **+0.33** | 12 applicable dimensions each (Logging N/A). |

### What changed from the ad-hoc table

Three adjustments were made to align with the source docs:

- **Logging (row 6)** — new, from `design.md §5`. N/A here, but the dimension exists so future changes that touch log lines get scored.
- **Naming (row 3)** — promoted out of Maintainability, from `design.md §2`. Lets a change have great identifiers but bad abstraction choices (or vice versa) without the scores blurring together.
- **Backport / scope (row 13)** — new, from `CLAUDE.md`'s "1.22.x branch — backport-focused, avoid unnecessary refactors." Gives upstream credit for its tighter scope and surfaces local's `GetTagList` extension as a (defensible) scope expansion.
- **Process discipline** — *removed*. Openspec folders are an AI-workflow convention; human PRs are not required to produce them. Scoring upstream down for lacking them was unfair in an AI-vs-human comparison.

### Reading the score

Removing process discipline narrowed the gap from +0.62 to **+0.33** — the lead is now concentrated in real engineering differences rather than workflow artifacts.

- **Local's structural wins** (rows 8, 10, 12) are testing, query efficiency, and type simplicity. They compound over time and cost little to keep.
- **Upstream's surface wins** (rows 7, 11, 13) are UX polish, tooltip text, and tighter scope. The big one (row 11) is one template edit away.
- **Net recommendation**: port upstream's template restructure into the local branch (closes row 11, mostly closes row 7, modest hit to row 13 from extra template surface). Leave `GetTagList` filtering in place — it's a defensible interpretation of #31998 even if upstream read the issue more narrowly. Expected local average after the port: **~4.6 / 5**.
