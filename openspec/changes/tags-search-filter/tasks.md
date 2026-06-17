## 1. Test Foundation (TDD Red Phase)

- [x] 1.1 Append `TestTagsSearch` integration test to `tests/integration/repo_tag_test.go` with three subtests: `Hit` (asserts only matching tags render for `?q=v1.1`), `Miss` (asserts zero tag rows for a non-matching keyword), and `EmptyKeywordPreservesBehavior` (asserts the no-`q` path still returns all tags). Add `github.com/PuerkitoBio/goquery` to the import block.
- [x] 1.2 Run `make test-sqlite#TestTagsSearch` to confirm the test fails in the expected TDD red-phase pattern (`Hit` and `Miss` fail because the `q` parameter is ignored; `EmptyKeywordPreservesBehavior` passes).
- [x] 1.3 Commit: `test: add failing integration test for tags page search filter`.

## 2. Data-Layer Keyword Support

- [x] 2.1 Add a `Keyword string` field to `FindReleasesOptions` in `models/repo/release.go` (after `HasSha1`).
- [x] 2.2 Append a `Keyword` clause to `FindReleasesOptions.ToConds()` between the `HasSha1` block and `return cond`: `if opts.Keyword != "" { cond = cond.And(builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)}) }`. Verify the `strings` and `xorm.io/builder` imports are already present.
- [x] 2.3 Run `go build ./models/repo/...` to confirm the package compiles.
- [x] 2.4 Run `make test-sqlite#TestTagsSearch` to confirm the test still fails (handler not wired yet; model layer alone is insufficient).
- [x] 2.5 Commit: `feat: support keyword filter in FindReleasesOptions`.

## 3. Handler Wiring

- [x] 3.1 In `routers/web/repo/release.go:TagsList`, read the keyword after the `listOptions` block: `keyword := ctx.FormString("q")`.
- [x] 3.2 Pass the keyword into the `FindReleasesOptions` literal as `Keyword: keyword,`.
- [x] 3.3 Expose the keyword to the template by adding `ctx.Data["Keyword"] = keyword` right after `ctx.Data["Releases"] = releases`. (`pager.SetDefaultParams(ctx)` already auto-carries `q` through page links.)
- [x] 3.4 Run `go build ./routers/web/repo/...` to confirm the package compiles.
- [x] 3.5 Run `make test-sqlite#TestTagsSearch`. The `Hit` and `Miss` subtests should now pass because the SQL filter is wired end-to-end.
- [x] 3.6 Commit: `feat: read q keyword in TagsList handler`.

## 4. UI Surface (Locale + Template)

- [x] 4.1 Add `tag_kind = Search tags...` to `options/locale/locale_en-US.ini` in the `[search]` section, immediately after `commit_kind`.
- [x] 4.2 Insert a `<div class="ui attached segment"><form class="ignore-dirty" method="get">{{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}</form></div>` block in `templates/repo/tag/list.tmpl` between `{{template "repo/release_tag_header" .}}` and `{{if .Releases}}` so the form renders even on empty search results.
- [x] 4.3 Run `make test-sqlite#TestTagsSearch` to confirm all three subtests still pass.
- [x] 4.4 Commit: `feat: add keyword search form to tags page`.

## 5. Pagination Count Fix (Bug Surfaced by Feature Work)

- [x] 5.1 Replace `db.Find[repo_model.Release](ctx, opts)` with `db.FindAndCount[repo_model.Release](ctx, opts)` in `TagsList`, capturing the filtered count as `total`.
- [x] 5.2 Replace `numTags := ctx.Data["NumTags"].(int64)` and `context.NewPagination(int(numTags), ...)` with `context.NewPagination(int(total), ...)`. Leave `ctx.Data["NumTags"]` itself untouched (still consumed by the repo header).
- [x] 5.3 Add page normalization before the `opts` literal: `if listOptions.Page <= 0 { listOptions.Page = 1 }`. (Required because `db.FindAndCount` does not normalize `page == 0` the way `db.Find` does — matches the Branches convention at `routers/web/repo/branch.go:50-52`.)
- [x] 5.4 Run `go build ./routers/web/repo/...` and `make test-sqlite#TestTagsSearch` to confirm everything still passes.
- [x] 5.5 Run regression tests: `make test-sqlite#TestViewReleases`, `make test-sqlite#TestCreateNewTagProtected`, `make test-sqlite#TestRepushTag`. All should pass.
- [x] 5.6 Commit: `fix: use filtered count for tags page pagination`.

## 6. Test Coverage Expansion

- [x] 6.1 Add a `CaseInsensitive` subtest to `TestTagsSearch` exercising `?q=V1.1` to verify the `lower_tag_name` + `strings.ToLower` matching path. Add `"strings"` to the test file imports.
- [x] 6.2 Add a `PaginationWithKeyword` subtest using `?q=v1.0&limit=1` to verify the pagination widget does NOT render when the filtered count fits in a single page (true regression guard: under the buggy unfiltered-count code, pagination WOULD render for 3 pages).
- [x] 6.3 Run `make test-sqlite#TestTagsSearch` to confirm all five subtests pass.
- [x] 6.4 Commit: `test: cover case-insensitive and pagination+keyword on tags search`.
- [x] 6.5 (Follow-up) Strengthen `PaginationWithKeyword` to a true regression guard by switching from `q=v1&limit=1` (false-negative — both filtered=2 and unfiltered=3 exceed the limit) to `q=v1.0&limit=1` (filtered=1 vs unfiltered=3 flips the pagination widget on/off). Commit: `test: make PaginationWithKeyword a true regression guard`.

## 7. Verification

- [x] 7.1 Run `make backend` to confirm the full backend build succeeds.
- [x] 7.2 Run `go vet ./models/repo/... ./routers/web/repo/...` to confirm no vet issues on the touched packages.
- [x] 7.3 Run `gofmt -l` on all five touched files to confirm no formatting drift.
- [x] 7.4 Run regression suite: `make test-sqlite#TestViewReleases`, `make test-sqlite#TestViewReleasesNoLogin`, `make test-sqlite#TestCreateNewTagProtected`, `make test-sqlite#TestRepushTag`, `make test-sqlite#TestTagsSearch`.
- [ ] 7.5 Manual browser verification: launch `./gitea web`, navigate to `/user2/repo1/tags`, confirm the search input renders with `Search tags...` placeholder, type `v1.1` and press Enter, confirm only matching tags appear. Type `V1.1` to confirm case-insensitive matching. Type a non-matching string to confirm the empty state. Confirm pagination links preserve the keyword.
- [ ] 7.6 (Optional) Run `make lint-go`, `make lint-templates`, `make lint-css` to confirm no linter regressions on the touched files. (Note: `make lint-go` requires network access to fetch `golangci-lint`; `make lint-templates` requires `poetry`. Both were unavailable in the local sandbox during initial verification — `go vet` and `gofmt` were used as substitutes.)

## Implementation Status

All code work is complete on branch `training/w1-try-2` (commits `2cd12a1ed9` through `45036094a3`). The two unchecked items in Section 7 are manual verification steps that the user must run.

The full brainstorming spec and step-by-step implementation plan are preserved at:
- `docs/superpowers/specs/2026-06-15-tags-search-filter-design.md`
- `docs/superpowers/plans/2026-06-15-tags-search-filter.md`
