## 1. Data Layer

- [ ] 1.1 Add `Keyword string` field to `FindReleasesOptions` in `models/repo/release.go`
- [ ] 1.2 In `FindReleasesOptions.ToConds()` (in `models/repo/release.go`), add `if opts.Keyword != "" { cond = cond.And(builder.Like{"tag_name", opts.Keyword}) }` — matches Branches (`models/git/branch_list.go:102-104`); xorm auto-wraps `%kw%`
- [ ] 1.3 Confirm no other caller of `FindReleases` sets `Keyword` (default-empty guard holds)

## 2. Handler

- [ ] 2.1 In `routers/web/repo/release.go:TagsList`, read `kw := ctx.FormString("q")`
- [ ] 2.2 Set `opts.Keyword = kw` before calling `db.Find[repo_model.Release]`
- [ ] 2.3 Set `ctx.Data["Keyword"] = kw`
- [ ] 2.4 Replace `numTags := ctx.Data["NumTags"].(int64)` + `context.NewPagination(int(numTags), ...)` (currently `release.go:244-245`) with a filtered count: `filteredCount, err := db.Count[repo_model.Release](ctx, opts)` and pass `int(filteredCount)` to `NewPagination`. Branches does the same (`routers/web/repo/branch.go:84` uses `branchesCount` from the filtered query). The total `NumTags` shown in the header tab (`templates/repo/release_tag_header.tmpl:10`) is unchanged — it's still the repo total from middleware.
- [ ] 2.5 Verify `pager.SetDefaultParams(ctx)` is already invoked (it is, `release.go:248`) — confirm `q` survives into pagination links

## 3. Template

- [ ] 3.1 In `templates/repo/tag/list.tmpl`, restructure the listing region to mirror `templates/repo/branch/list.tmpl:70-82`: move the existing `<h4 class="ui top attached header">` block (currently lines 8-12) OUT of the `{{if .Releases}}` gate; insert a `<div class="ui attached segment"><form class="ignore-dirty" method="get">{{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}</form></div>` between the `<h4>` and the table wrapper; add an `{{else}}` branch inside `{{if .Releases}}` rendering `<div class="ui attached segment">{{ctx.Locale.Tr "search.no_results"}}</div>` so a no-match search still shows the search box and a message
- [ ] 3.2 Run `make lint-templates` to validate

## 4. Locale

- [ ] 4.1 Add `tag_kind = Search tags...` under the `[search]` section in `options/locale/locale_en-US.ini` (next to existing `branch_kind` at line 178 and `commit_kind` at line 179). Template references it as `search.tag_kind` because the INI section name becomes the key prefix.
- [ ] 4.2 Do not machine-translate other locale files (community workflow)

## 5. Tests

- [ ] 5.1 Unit: extend `models/repo/release_test.go` with cases for `FindReleasesOptions.Keyword` — (a) empty returns all tags, (b) substring matches expected subset, (c) non-matching keyword returns empty without error
- [ ] 5.2 Integration: extend or create `tests/integration/repo_tag_test.go` — GET `/tags?q=<kw>` returns 200, HTML contains matching tag rows, omits non-matching rows, pagination link contains `q=<kw>`
- [ ] 5.3 Integration: GET `/tags` (no `q`) returns same rows as before this change

## 6. Verification

- [ ] 6.1 `make lint` passes
- [ ] 6.2 `make test-backend` passes (including new unit test)
- [ ] 6.3 `make test-sqlite` passes (including new integration test)
- [ ] 6.4 Manual: start dev server, browse `/tags` on a repo with >10 tags, type a substring, confirm filter + pagination + empty-state behaviors
