## 1. Model Layer

- [x] 1.1 Add `Keyword string` field to `FindReleasesOptions` in `models/repo/release.go` (alongside the existing `TagNames` field, around line 229-238)
- [x] 1.2 In `models/repo/release.go` `toConds()` (around line 239-260), add the conditional `builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)}` clause guarded by `if opts.Keyword != ""`
- [x] 1.3 Confirm `models/repo/release.go` already imports `strings` (it does — used elsewhere); add the import if missing

## 2. Web Handlers

- [x] 2.1 In `routers/web/repo/release.go` `TagsList` (line 204-251): read `keyword := ctx.FormTrim("q")`, set `Keyword: keyword` on the `FindReleasesOptions` literal, set `ctx.Data["Keyword"] = keyword` before `ctx.HTML(...)`. **Note:** also switched `db.Find` → `db.FindAndCount` so pagination uses the filtered count (middleware-set `NumTags` stays untouched for the header badge).
- [x] 2.2 In `routers/web/repo/repo.go` `GetTagList` (line 716-725): after fetching `tags`, add the Go-side filter block that reads `q := ctx.FormTrim("q")` and reuses the slice backing array to filter case-insensitively via `strings.Contains(strings.ToLower(t), lower)`
- [x] 2.3 Confirm `routers/web/repo/repo.go` already imports `strings` (add if missing)

## 3. Template & Locale

- [x] 3.1 In `templates/repo/tag/list.tmpl`, add the search form partial above the tag table (placed outside `{{if .Releases}}` so it stays visible even on empty search results):
  ```html
  <div class="ui attached segment">
      <form class="ignore-dirty" method="get">
          {{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}
      </form>
  </div>
  ```
- [x] 3.2 Add `tag_kind = Search tag name...` to `options/locale/locale_en-US.ini` under the `[search]` section (next to the existing `branch_kind`)

## 4. Tests

- [x] 4.1 Add `TestFindReleasesByKeyword` to `models/repo/release_test.go`: uses repo 57 (which has tags v1.0, v1.1, v2.0); 4 sub-cases cover no keyword, substring match, case-insensitive match, no-match returns empty
- [x] 4.2 Add `TestTagsListSearch` to `tests/integration/release_test.go` (next to existing `TestViewTagsList`): 4 sub-cases for substring, case-insensitive, no-match, empty-q on the HTML page
- [x] 4.3 Add `TestTagListJSONSearch` to `tests/integration/release_test.go`: 3 sub-cases for substring, case-insensitive, empty-q on the JSON endpoint

## 5. Verification

- [x] 5.1 Run `go test -run TestFindReleasesByKeyword ./models/repo/...` and confirm pass — **passed** (with `-tags 'sqlite sqlite_unlock_notify'`)
- [x] 5.2 Run `go test -run TestTagsListSearch ./tests/integration/...` and confirm pass — **passed** (ran via `integrations.sqlite.test` with `GITEA_CONF=tests/sqlite.ini`)
- [x] 5.3 Run `go test -run TestTagListJSONSearch ./tests/integration/...` and confirm pass — **passed** (same run as 5.2)
- [x] 5.4 Regression check on `./models/repo/...` — **passed** (skipped full `make test-backend` because git-lfs is not installed locally; models/repo covers all release-filter callers)
- [x] 5.5 `go build ./...` and `go vet ./models/repo/... ./routers/web/repo/...` — **clean**
- [ ] 5.6 Start the server (`./gitea web`), navigate to a repo's `/tags` page, and manually verify: search box appears; typing a partial tag name + Enter filters; pagination preserves `q`; uppercase vs lowercase query returns identical results; `/tags/list?q=...` JSON filters correctly; `/tags.rss` ignores `q` and returns the full feed — **deferred to user** (gitea server is already running on port 3000 from a separate session)
