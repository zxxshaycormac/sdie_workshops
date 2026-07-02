# Tags Page Search Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add keyword search/filter to the repository Tags page (`/{org}/{repo}/tags`), consistent with the existing Branches search UX.

**Architecture:** Replicate the Branches search pattern — a plain HTML GET form submits `?q=keyword`; the `TagsList` handler reads it, passes it as `Keyword` on `FindReleasesOptions`, and `ToConds()` adds a case-insensitive `LIKE` on `lower_tag_name`. Switch from `db.Find` to `db.FindAndCount` so pagination uses the filtered count.

**Tech Stack:** Go (Gitea 1.22.x), XORM/builder, Go HTML templates, Fomantic UI

## Global Constraints

- Go pinned to 1.22.x via `.tool-versions` (asdf)
- Use `xorm.io/builder` for query conditions — never interpolate user values into SQL
- All user-visible strings via locale keys (`options/locale/locale_en-US.ini`)
- Templates use `{{ctx.Locale.Tr "key"}}` — never hardcoded English
- No inline `<script>` in templates
- `TAGS="bindata sqlite sqlite_unlock_notify"` for building/testing
- This is the 1.22.x branch — keep changes minimal, no unnecessary refactors

---

## File Structure

| File | Action | Responsibility |
|------|--------|----------------|
| `models/repo/release.go` | Modify | Add `Keyword` field to `FindReleasesOptions` + `LIKE` filter in `ToConds()` |
| `models/repo/release_test.go` | Modify | Add unit test for keyword filtering |
| `routers/web/repo/release.go` | Modify | Read `q` param, pass keyword, switch to `FindAndCount`, set `ctx.Data["Keyword"]` |
| `templates/repo/tag/list.tmpl` | Modify | Insert search form using `shared/search/combo` |
| `options/locale/locale_en-US.ini` | Modify | Add `search.tag_kind` key |
| `tests/integration/release_test.go` | Modify | Add integration tests for tags page search |

---

## Task 1: Model Layer — Add Keyword Filter to FindReleasesOptions

**Files:**
- Modify: `models/repo/release.go:229-266` (struct `FindReleasesOptions` + method `ToConds`)
- Test: `models/repo/release_test.go`

**Interfaces:**
- Produces: `FindReleasesOptions.Keyword string` field — read by `TagsList` handler in Task 2
- The `ToConds()` method already exists and is called by `db.Find`/`db.FindAndCount`/`db.Count`

**Context:**
- The struct `FindReleasesOptions` is at lines 229-238. It already has fields like `RepoID`, `IncludeDrafts`, `TagNames`.
- `ToConds()` is at lines 240-266. It builds `builder.Cond` chain with `builder.Eq`, `builder.In`, `builder.Neq`.
- Imports `strings` (line 14) and `xorm.io/builder` (line 24) are already present — no new imports needed.
- The `release` model has a `LowerTagName string` column (lowercased tag name) used for matching.
- Fixtures in `models/fixtures/release.yml` include repo 1 with tags: `v1.1` (id 1, is_tag=false), `delete-tag` (id 3, is_tag=true), `draft-release` (id 4, is_draft=true, is_tag=false, no sha1), `v1.0` (id 5, is_tag=false). The tags page uses `IncludeTags: true, IncludeDrafts: true, HasSha1: true` which means only tags with a real sha1 show: `v1.1` (id 1), `delete-tag` (id 3), `v1.0` (id 5).

- [ ] **Step 1: Write the failing unit test**

Add this test to the end of `models/repo/release_test.go`:

```go
func TestFindReleasesByKeyword(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// repo 1 has tags with sha1: "v1.1" (id 1), "delete-tag" (id 3), "v1.0" (id 5)
	opts := FindReleasesOptions{
		IncludeDrafts: true,
		IncludeTags:   true,
		HasSha1:       optional.Some(true),
		RepoID:        1,
	}

	// No keyword — should return all 3 tags
	releases, err := db.Find[Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Len(t, releases, 3)

	// Keyword "v1" — case-insensitive match on "v1.1" and "v1.0"
	opts.Keyword = "v1"
	releases, err = db.Find[Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Len(t, releases, 2)
	tagNames := make([]string, 0, 2)
	for _, r := range releases {
		tagNames = append(tagNames, r.TagName)
	}
	assert.ElementsMatch(t, []string{"v1.1", "v1.0"}, tagNames)

	// Keyword "V1" — uppercase keyword, should still match (case-insensitive)
	opts.Keyword = "V1"
	releases, err = db.Find[Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Len(t, releases, 2)

	// Keyword "DELETE" — uppercase keyword matches "delete-tag"
	opts.Keyword = "DELETE"
	releases, err = db.Find[Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "delete-tag", releases[0].TagName)

	// Keyword that matches nothing
	opts.Keyword = "nonexistent-tag-xyz"
	releases, err = db.Find[Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Empty(t, releases)
}
```

You also need to update the imports in `release_test.go` to add `"code.gitea.io/gitea/modules/optional"`. The existing imports are `db`, `unittest`, and `assert`. Since the test is in `package repo` (same package as `release.go`), `FindReleasesOptions` and `Release` are referenced directly without any package alias.

- [ ] **Step 2: Run the test to verify it fails**

```bash
go test -run TestFindReleasesByKeyword ./models/repo/ -v
```
Expected: FAIL — compile error because `FindReleasesOptions` has no `Keyword` field.

- [ ] **Step 3: Add the Keyword field and filter**

In `models/repo/release.go`, add `Keyword` to the `FindReleasesOptions` struct (after line 237, before the closing brace):

```go
	Keyword       string
```

The full struct becomes:
```go
type FindReleasesOptions struct {
	db.ListOptions
	RepoID        int64
	IncludeDrafts bool
	IncludeTags   bool
	IsPreRelease  optional.Option[bool]
	IsDraft       optional.Option[bool]
	TagNames      []string
	HasSha1       optional.Option[bool] // useful to find draft releases which are created with existing tags
	Keyword       string
}
```

In `ToConds()`, add the keyword condition before the final `return cond` (after the `HasSha1` block, before line 265):

```go
	if opts.Keyword != "" {
		cond = cond.And(builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)})
	}
```

- [ ] **Step 4: Run the test to verify it passes**

```bash
go test -run TestFindReleasesByKeyword ./models/repo/ -v
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add models/repo/release.go models/repo/release_test.go
git commit -m "Add keyword filter to FindReleasesOptions for tag search"
```

---

## Task 2: Handler + Template + Locale — Wire Up the Search UI

**Files:**
- Modify: `routers/web/repo/release.go:204-251` (handler `TagsList`)
- Modify: `templates/repo/tag/list.tmpl:6-7` (insert search form)
- Modify: `options/locale/locale_en-US.ini:179` (add locale key)
- Test: `tests/integration/release_test.go` (add integration test)

**Interfaces:**
- Consumes: `FindReleasesOptions.Keyword` from Task 1
- Produces: `ctx.Data["Keyword"]` string for the template; `?q=` query param handling on `GET /tags`

**Context:**
- The `TagsList` handler is at `routers/web/repo/release.go:204-251`.
- Currently it uses `db.Find[repo_model.Release](ctx, opts)` (line 236) and reads `ctx.Data["NumTags"]` (line 244) for pagination total.
- `db.FindAndCount` (defined at `models/db/list.go:186`) returns `([]*T, int64, error)` — the int64 is the count matching the same conditions.
- The Branches handler (`routers/web/repo/branch.go:55`) reads the keyword via `kw := ctx.FormString("q")` and sets `ctx.Data["Keyword"] = kw`.
- The template `templates/repo/tag/list.tmpl` has an `<h4>` header at lines 8-12 and the table at line 14. The search form goes between them, matching `branch/list.tmpl:76-80`.
- The locale file has `branch_kind` at line 178 and `commit_kind` at line 179.
- The existing integration test `TestViewTagsList` at `tests/integration/release_test.go:220` loads `/{repo}/tags` and asserts 3 tags appear. This is the pattern to extend.

- [ ] **Step 1: Write the failing integration test**

Add these test functions to the end of `tests/integration/release_test.go`:

```go
func TestViewTagsListSearch(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	link := repo.Link() + "/tags"
	session := loginUser(t, "user1")

	// Search for "v1" — should match "v1.0" and "v1.1" (case-insensitive)
	req := NewRequest(t, "GET", link+"?q=v1")
	rsp := session.MakeRequest(t, req, http.StatusOK)
	htmlDoc := NewHTMLParser(t, rsp.Body)
	tags := htmlDoc.Find(".tag-list-row-link")
	assert.Equal(t, 2, tags.Length())
	tagNames := make([]string, 0, 2)
	tags.Each(func(i int, s *goquery.Selection) {
		tagNames = append(tagNames, s.Text())
	})
	assert.ElementsMatch(t, []string{"v1.0", "v1.1"}, tagNames)

	// Search input should reflect the current keyword
	searchInput := htmlDoc.Find(`input[name="q"]`)
	assert.Equal(t, 1, searchInput.Length())
	assert.Equal(t, "v1", searchInput.AttrOr("value", ""))

	// Uppercase keyword should also match (case-insensitive)
	req = NewRequest(t, "GET", link+"?q=V1")
	rsp = session.MakeRequest(t, req, http.StatusOK)
	htmlDoc = NewHTMLParser(t, rsp.Body)
	tags = htmlDoc.Find(".tag-list-row-link")
	assert.Equal(t, 2, tags.Length())

	// Search for nonexistent tag — empty result, page still 200
	req = NewRequest(t, "GET", link+"?q=nonexistent-xyz")
	rsp = session.MakeRequest(t, req, http.StatusOK)
	htmlDoc = NewHTMLParser(t, rsp.Body)
	tags = htmlDoc.Find(".tag-list-row-link")
	assert.Equal(t, 0, tags.Length())

	// No keyword — all tags shown (regression check)
	req = NewRequest(t, "GET", link)
	rsp = session.MakeRequest(t, req, http.StatusOK)
	htmlDoc = NewHTMLParser(t, rsp.Body)
	tags = htmlDoc.Find(".tag-list-row-link")
	assert.Equal(t, 3, tags.Length())
}
```

- [ ] **Step 2: Run the test to verify it fails**

```bash
go test -run TestViewTagsListSearch ./tests/integration/ -v
```
Expected: FAIL — the search input `input[name="q"]` doesn't exist yet (0 length), and the `?q=v1` filter isn't applied so all 3 tags appear instead of 2.

- [ ] **Step 3: Add the locale key**

In `options/locale/locale_en-US.ini`, add after line 179 (`commit_kind = Search commits...`):

```ini
tag_kind = Search tags...
```

- [ ] **Step 4: Update the TagsList handler**

In `routers/web/repo/release.go`, modify the `TagsList` function (lines 204-251).

After the `listOptions` block (after line 224), add reading the keyword:

```go
	kw := ctx.FormString("q")
```

In the `opts` struct literal (lines 226-234), add the `Keyword` field:

```go
	opts := FindReleasesOptions{
		ListOptions: listOptions,
		// for the tags list page, show all releases with real tags (having real commit-id),
		// the drafts should also be included because a real tag might be used as a draft.
		IncludeDrafts: true,
		IncludeTags:   true,
		HasSha1:       optional.Some(true),
		RepoID:        ctx.Repo.Repository.ID,
		Keyword:       kw,
	}
```

Replace the `db.Find` call and pagination block (lines 236-247). Change from:

```go
	releases, err := db.Find[repo_model.Release](ctx, opts)
	if err != nil {
		ctx.ServerError("GetReleasesByRepoID", err)
		return
	}

	ctx.Data["Releases"] = releases

	numTags := ctx.Data["NumTags"].(int64)
	pager := context.NewPagination(int(numTags), opts.PageSize, opts.Page, 5)
	pager.SetDefaultParams(ctx)
	ctx.Data["Page"] = pager
```

To:

```go
	releases, count, err := db.FindAndCount[repo_model.Release](ctx, opts)
	if err != nil {
		ctx.ServerError("GetReleasesByRepoID", err)
		return
	}

	ctx.Data["Releases"] = releases
	ctx.Data["Keyword"] = kw

	pager := context.NewPagination(int(count), opts.PageSize, opts.Page, 5)
	pager.SetDefaultParams(ctx)
	ctx.Data["Page"] = pager
```

Note: `count` from `FindAndCount` equals the total tag count when no keyword is set (no `LIKE` condition added), so this preserves existing pagination behavior when not searching.

- [ ] **Step 5: Add the search form to the template**

In `templates/repo/tag/list.tmpl`, insert the search form between the `{{template "repo/release_tag_header" .}}` line (line 6) and the `{{if .Releases}}` block (line 7). The form must be **outside** the `{{if .Releases}}` conditional so it shows even when there are zero results.

Insert after line 6:

```go
		<div class="ui attached segment">
			<form class="ignore-dirty" method="get">
				{{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}
			</form>
		</div>
```

The surrounding context will look like:

```go
		{{template "repo/release_tag_header" .}}
		<div class="ui attached segment">
			<form class="ignore-dirty" method="get">
				{{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}
			</form>
		</div>
		{{if .Releases}}
```

- [ ] **Step 6: Run the integration test to verify it passes**

```bash
go test -run TestViewTagsListSearch ./tests/integration/ -v
```
Expected: PASS

- [ ] **Step 7: Run the existing tags list test to verify no regression**

```bash
go test -run TestViewTagsList$ ./tests/integration/ -v
```
Expected: PASS — still finds 3 tags with no keyword.

- [ ] **Step 8: Commit**

```bash
git add routers/web/repo/release.go templates/repo/tag/list.tmpl options/locale/locale_en-US.ini tests/integration/release_test.go
git commit -m "Add keyword search to tags page"
```

---

## Task 3: Lint and Full Verification

**Files:** None (verification only)

- [ ] **Step 1: Run Go formatting and lint**

```bash
make fmt
```

- [ ] **Step 2: Run lint**

```bash
make lint-go
```
Expected: PASS — no new lint issues introduced.

- [ ] **Step 3: Run all related backend tests**

```bash
go test -run "TestFindReleases|TestViewTagsList" ./models/repo/ ./tests/integration/ -v
```
Expected: All tests PASS.

- [ ] **Step 4: Manual smoke test (if running locally)**

Build and run:
```bash
TAGS="bindata sqlite sqlite_unlock_notify" make build
./gitea web
```

Navigate to a repo's Tags page:
1. Verify the search box appears above the tags table
2. Type a partial tag name, press Enter — only matching tags should show
3. Verify case-insensitive: type uppercase when tag is lowercase
4. Verify pagination links preserve the `q` param in the URL
5. Clear the search — all tags should reappear

- [ ] **Step 5: Final commit if formatting changed anything**

```bash
git add -A
git diff --cached --quiet || git commit --amend --no-edit
```
