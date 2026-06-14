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
