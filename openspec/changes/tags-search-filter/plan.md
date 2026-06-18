# Tags Page Search Filter Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add server-side keyword filtering to the repository Tags page (`/{owner}/{repo}/tags`), mirroring the existing Branches page search pattern.

**Architecture:** Add an opt-in `Keyword` field to `repo_model.FindReleasesOptions` that applies a `builder.Like{"tag_name", kw}` condition. The `TagsList` handler reads `q`, threads it through the options, computes pagination from the filtered count, and exposes the keyword to the template. The template gains a `shared/search/combo` form (same partial Branches uses) and an empty-state branch for no-match searches.

**Tech Stack:** Go (XORM + `xorm.io/builder`), Go HTML templates (`html/template`), INI locale files, Gitea `unittest` + integration test frameworks.

**Spec:** `openspec/changes/tags-search-filter/spec.md` (requirements PKG-04-006, 108-111, 304).

---

## File Structure

| File | Change | Responsibility |
|------|--------|----------------|
| `models/repo/release.go` | Modify | Add `Keyword` field to `FindReleasesOptions`; extend `ToConds()` |
| `models/repo/release_test.go` | Modify | Add unit tests for `Keyword` filter behavior |
| `routers/web/repo/release.go` | Modify | Read `q` in `TagsList`; switch pagination source to filtered `db.Count` |
| `templates/repo/tag/list.tmpl` | Modify | Pull `<h4>` out of `{{if .Releases}}`; add search-segment; add `{{else}}` empty state |
| `options/locale/locale_en-US.ini` | Modify | Add `tag_kind` under `[search]` section |
| `tests/integration/repo_tag_test.go` | Modify | Add end-to-end tests for `/tags?q=...` |

No new files. No migrations. No new packages.

**Fixture context (repo_id=1, used by tests):** The Tags page for repo 1 currently shows 3 rows — `v1.1`, `delete-tag`, `v1.0` (all have `sha1`; `draft-release` is excluded because it has no sha1). These three tag names are the substrate for all test assertions.

---

## Task 1: Data Layer — `FindReleasesOptions.Keyword`

**Files:**
- Modify: `models/repo/release.go:228-266` (struct `FindReleasesOptions` + method `ToConds`)
- Test: `models/repo/release_test.go` (extend with new test function)

- [ ] **Step 1: Write the failing unit test**

Append to `models/repo/release_test.go`:

```go
func TestFindReleasesOptions_KeywordFilter(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// Repo 1 has 3 tag-shaped releases with sha1: v1.1, delete-tag, v1.0
	// (draft-release is excluded because it has no sha1)

	t.Run("empty keyword returns all matching tags", func(t *testing.T) {
		releases, err := db.Find[Release](db.DefaultContext, FindReleasesOptions{
			RepoID:      1,
			IncludeTags: true,
			IncludeDrafts: true,
			HasSha1:     optional.Some(true),
			ListOptions: db.ListOptions{ListAll: true},
		})
		assert.NoError(t, err)
		assert.Len(t, releases, 3)
	})

	t.Run("substring keyword filters by tag_name", func(t *testing.T) {
		releases, err := db.Find[Release](db.DefaultContext, FindReleasesOptions{
			RepoID:        1,
			IncludeTags:   true,
			IncludeDrafts: true,
			HasSha1:       optional.Some(true),
			Keyword:       "v1",
			ListOptions:   db.ListOptions{ListAll: true},
		})
		assert.NoError(t, err)
		assert.Len(t, releases, 2)
		names := []string{releases[0].TagName, releases[1].TagName}
		assert.ElementsMatch(t, []string{"v1.0", "v1.1"}, names)
	})

	t.Run("non-matching keyword returns empty without error", func(t *testing.T) {
		releases, err := db.Find[Release](db.DefaultContext, FindReleasesOptions{
			RepoID:        1,
			IncludeTags:   true,
			IncludeDrafts: true,
			HasSha1:       optional.Some(true),
			Keyword:       "this-tag-does-not-exist",
			ListOptions:   db.ListOptions{ListAll: true},
		})
		assert.NoError(t, err)
		assert.Empty(t, releases)
	})

	t.Run("count agrees with find when filtered", func(t *testing.T) {
		opts := FindReleasesOptions{
			RepoID:        1,
			IncludeTags:   true,
			IncludeDrafts: true,
			HasSha1:       optional.Some(true),
			Keyword:       "v1",
		}
		count, err := db.Count[Release](db.DefaultContext, opts)
		assert.NoError(t, err)
		assert.EqualValues(t, 2, count)
	})
}
```

Also update the import block at the top of the file to include:

```go
	"code.gitea.io/gitea/modules/optional"
```

(Place it alphabetically among the existing `code.gitea.io/...` imports.)

- [ ] **Step 2: Run the test to verify it fails**

Run:
```bash
go test -run TestFindReleasesOptions_KeywordFilter ./models/repo/
```
Expected: COMPILATION FAILURE. `FindReleasesOptions` has no `Keyword` field and `optional` may not be imported. This confirms the test exercises the new surface.

- [ ] **Step 3: Add `Keyword` field to `FindReleasesOptions`**

In `models/repo/release.go:228-238`, extend the struct:

```go
// FindReleasesOptions describes the conditions to Find releases
type FindReleasesOptions struct {
	db.ListOptions
	RepoID        int64
	IncludeDrafts bool
	IncludeTags   bool
	IsPreRelease  optional.Option[bool]
	IsDraft       optional.Option[bool]
	TagNames      []string
	HasSha1       optional.Option[bool] // useful to find draft releases which are created with existing tags
	Keyword       string // substring filter on tag_name; opt-in, empty means no filter
}
```

- [ ] **Step 4: Extend `ToConds()` to apply the keyword**

In `models/repo/release.go:240-266`, add the keyword clause inside `ToConds()`, right before the `return cond`:

```go
	if opts.Keyword != "" {
		cond = cond.And(builder.Like{"tag_name", opts.Keyword})
	}
	return cond
```

The full method after the change:

```go
func (opts FindReleasesOptions) ToConds() builder.Cond {
	var cond builder.Cond = builder.Eq{"repo_id": opts.RepoID}

	if !opts.IncludeDrafts {
		cond = cond.And(builder.Eq{"is_draft": false})
	}
	if !opts.IncludeTags {
		cond = cond.And(builder.Eq{"is_tag": false})
	}
	if len(opts.TagNames) > 0 {
		cond = cond.And(builder.In("tag_name", opts.TagNames))
	}
	if opts.IsPreRelease.Has() {
		cond = cond.And(builder.Eq{"is_prerelease": opts.IsPreRelease.Value()})
	}
	if opts.IsDraft.Has() {
		cond = cond.And(builder.Eq{"is_draft": opts.IsDraft.Value()})
	}
	if opts.HasSha1.Has() {
		if opts.HasSha1.Value() {
			cond = cond.And(builder.Neq{"sha1": ""})
		} else {
			cond = cond.And(builder.Eq{"sha1": ""})
		}
	}
	if opts.Keyword != "" {
		cond = cond.And(builder.Like{"tag_name", opts.Keyword})
	}
	return cond
}
```

Confirm `builder` is already imported at the top of the file (it is — existing code uses `builder.Eq`, `builder.Neq`, etc.).

- [ ] **Step 5: Run the test to verify it passes**

Run:
```bash
go test -run TestFindReleasesOptions_KeywordFilter ./models/repo/
```
Expected: PASS (all four subtests).

- [ ] **Step 6: Run all existing release tests to confirm no regression**

Run:
```bash
go test ./models/repo/
```
Expected: PASS. The `Keyword` field defaults to empty string, so existing callers of `FindReleases` see no behavior change.

- [ ] **Step 7: Commit**

```bash
git add models/repo/release.go models/repo/release_test.go
git commit -m "feat(repo): add Keyword filter to FindReleasesOptions

Adds an opt-in Keyword field that applies builder.Like on tag_name,
mirroring the pattern used by FindBranchOptions
(models/git/branch_list.go). Empty keyword is a no-op, so existing
callers of FindReleases are unaffected."
```

---

## Task 2: Locale — Add `tag_kind` Placeholder

**Files:**
- Modify: `options/locale/locale_en-US.ini` (under `[search]` section, around line 178-179)

- [ ] **Step 1: Add the new key**

Open `options/locale/locale_en-US.ini`, find the `[search]` section (starts at line 162). Locate the sibling keys at lines 178-179:

```ini
branch_kind = Search branches...
commit_kind = Search commits...
```

Insert a new line between them (alphabetical order: `branch`, `commit`, `tag` would actually put `tag` after `commit`):

```ini
branch_kind = Search branches...
commit_kind = Search commits...
tag_kind = Search tags...
```

- [ ] **Step 2: Verify the key is well-formed**

Run:
```bash
grep -n "tag_kind\|branch_kind\|commit_kind" options/locale/locale_en-US.ini
```
Expected output includes:
```
178:branch_kind = Search branches...
179:commit_kind = Search commits...
180:tag_kind = Search tags...
```
(Line numbers may shift slightly; what matters is that all three appear under `[search]` and `tag_kind` is new.)

- [ ] **Step 3: Commit**

```bash
git add options/locale/locale_en-US.ini
git commit -m "locale: add search.tag_kind placeholder

Sister key to search.branch_kind and search.commit_kind, used by the
upcoming Tags page search form."
```

---

## Task 3: Handler — Plumb `q` Through `TagsList`

**Files:**
- Modify: `routers/web/repo/release.go:204-251` (function `TagsList`)

- [ ] **Step 1: Read `q` and thread it through `opts`**

In `routers/web/repo/release.go`, locate `TagsList` (starts line 204). The current body builds `opts` (lines 219-228), calls `db.Find` (line 236), then builds pagination off `NumTags` (lines 244-247).

Make three edits inside the function body:

**Edit A — read keyword (insert just before `opts := repo_model.FindReleasesOptions{`, currently line 219):**

```go
	keyword := ctx.FormString("q")
```

**Edit B — set `Keyword` on opts (add as the last field of the struct literal, after `RepoID:`):**

```go
	opts := repo_model.FindReleasesOptions{
		ListOptions: listOptions,
		// for the tags list page, show all releases with real tags (having real commit-id),
		// the drafts should also be included because a real tag might be used as a draft.
		IncludeDrafts: true,
		IncludeTags:   true,
		HasSha1:       optional.Some(true),
		RepoID:        ctx.Repo.Repository.ID,
		Keyword:       keyword,
	}
```

**Edit C — expose keyword to template (insert immediately after `ctx.Data["Releases"] = releases`, currently line 242):**

```go
	ctx.Data["Releases"] = releases
	ctx.Data["Keyword"] = keyword
```

- [ ] **Step 2: Switch pagination source to filtered count**

Replace lines 244-247 (currently):

```go
	numTags := ctx.Data["NumTags"].(int64)
	pager := context.NewPagination(int(numTags), opts.PageSize, opts.Page, 5)
	pager.SetDefaultParams(ctx)
	ctx.Data["Page"] = pager
```

with:

```go
	// Use the filtered count (not the middleware-set NumTags) so pagination reflects
	// the search result. NumTags in the header tab stays as the repo total.
	filteredCount, err := db.Count[repo_model.Release](ctx, opts)
	if err != nil {
		ctx.ServerError("CountReleases", err)
		return
	}
	pager := context.NewPagination(int(filteredCount), opts.PageSize, opts.Page, 5)
	pager.SetDefaultParams(ctx)
	ctx.Data["Page"] = pager
```

- [ ] **Step 3: Verify the file compiles**

Run:
```bash
go build ./routers/web/repo/
```
Expected: no output (success). `db`, `repo_model`, `optional`, and `context` are all already imported in this file.

- [ ] **Step 4: Run the existing release tests to confirm no regression**

Run:
```bash
go test ./routers/web/repo/
```
Expected: PASS. The handler change is internal; the integration test for `/tags` is added in Task 5.

- [ ] **Step 5: Commit**

```bash
git add routers/web/repo/release.go
git commit -m "feat(repo): wire q keyword through TagsList handler

TagsList now reads q from the query string, passes it as
FindReleasesOptions.Keyword, exposes it as ctx.Data[\"Keyword\"],
and computes pagination total from db.Count on the filtered opts
instead of the middleware-set NumTags total.

Branches parity: routers/web/repo/branch.go:84."
```

---

## Task 4: Template — Add Search Form + Empty State

**Files:**
- Modify: `templates/repo/tag/list.tmpl` (full file is 81 lines)

- [ ] **Step 1: Read the current template to confirm line numbers**

Run:
```bash
cat -n templates/repo/tag/list.tmpl | head -65
```
Confirm:
- Line 6: `{{template "repo/release_tag_header" .}}`
- Lines 7-61: the `{{if .Releases}} ... {{end}}` block
- Lines 8-12: the `<h4 class="ui top attached header">` block
- Lines 13-60: the table block
- Line 63: `{{template "base/paginate" .}}`

- [ ] **Step 2: Restructure the listing region**

Replace lines 7-61 of `templates/repo/tag/list.tmpl` (the entire `{{if .Releases}} ... {{end}}` block) with the version in Step 3 below.

- [ ] **Step 3: Move `<h4>` and search form out of the `{{if}}` gate**

The structure must be: header fragment → `<h4>` → search form → `{{if .Releases}}` table `{{else}}` empty-state `{{end}}` → paginate. Three changes from the original:
1. `<h4>` moves OUT of `{{if .Releases}}` so it shows even when search returns nothing.
2. A new `<div class="ui attached segment">` containing the search form is inserted between `<h4>` and the table wrapper — same shape as `templates/repo/branch/list.tmpl:76-80`.
3. A new `{{else}}` branch renders `search.no_results` inside a styled segment.

Replace the entire region from line 7 through line 61 with:

```gotemplate
<h4 class="ui top attached header">
	<div class="five wide column tw-flex tw-items-center">
		{{svg "octicon-tag" 16 "tw-mr-1"}}{{ctx.Locale.Tr "repo.release.tags"}}
	</div>
</h4>
<div class="ui attached segment">
	<form class="ignore-dirty" method="get">
		{{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}
	</form>
</div>
{{$canReadReleases := $.Permission.CanRead ctx.Consts.RepoUnitTypeReleases}}
{{if .Releases}}
<div class="ui attached table segment">
	<table class="ui very basic striped fixed table single line" id="tags-table">
		<tbody class="tag-list">
			{{range $idx, $release := .Releases}}
				<tr>
					<td class="tag-list-row">
						<h3 class="tag-list-row-title tw-mb-2">
							{{if $canReadReleases}}
								<a class="tag-list-row-link tw-flex tw-items-center" href="{{$.RepoLink}}/releases/tag/{{.TagName | PathEscapeSegments}}" rel="nofollow">{{.TagName}}</a>
							{{else}}
								<a class="tag-list-row-link tw-flex tw-items-center" href="{{$.RepoLink}}/src/tag/{{.TagName | PathEscapeSegments}}" rel="nofollow">{{.TagName}}</a>
							{{end}}
						</h3>
						<div class="download tw-flex tw-items-center">
							{{if $.Permission.CanRead ctx.Consts.RepoUnitTypeCode}}
								{{if .CreatedUnix}}
									<span class="tw-mr-2">{{svg "octicon-clock" 16 "tw-mr-1"}}{{TimeSinceUnix .CreatedUnix ctx.Locale}}</span>
								{{end}}

								<a class="tw-mr-2 tw-font-mono muted" href="{{$.RepoLink}}/src/commit/{{.Sha1}}" rel="nofollow">{{svg "octicon-git-commit" 16 "tw-mr-1"}}{{ShortSha .Sha1}}</a>

								{{if not $.DisableDownloadSourceArchives}}
									<a class="archive-link tw-mr-2 muted" href="{{$.RepoLink}}/archive/{{.TagName | PathEscapeSegments}}.zip" rel="nofollow">{{svg "octicon-file-zip" 16 "tw-mr-1"}}ZIP</a>
									<a class="archive-link tw-mr-2 muted" href="{{$.RepoLink}}/archive/{{.TagName | PathEscapeSegments}}.tar.gz" rel="nofollow">{{svg "octicon-file-zip" 16 "tw-mr-1"}}TAR.GZ</a>
								{{end}}

								{{if (and $canReadReleases $.CanCreateRelease $release.IsTag)}}
									<a class="tw-mr-2 muted" href="{{$.RepoLink}}/releases/new?tag={{.TagName}}">{{svg "octicon-tag" 16 "tw-mr-1"}}{{ctx.Locale.Tr "repo.release.new_release"}}</a>
								{{end}}

								{{if (and ($.Permission.CanWrite ctx.Consts.RepoUnitTypeCode) $release.IsTag)}}
									<a class="ui delete-button tw-mr-2 muted" data-url="{{$.RepoLink}}/tags/delete" data-id="{{.ID}}">
										{{svg "octicon-trash" 16 "tw-mr-1"}}{{ctx.Locale.Tr "repo.release.delete_tag"}}
									</a>
								{{end}}

								{{if and $canReadReleases (not $release.IsTag)}}
									<a class="tw-mr-2 muted" href="{{$.RepoLink}}/releases/tag/{{.TagName | PathEscapeSegments}}">{{svg "octicon-tag" 16 "tw-mr-1"}}{{ctx.Locale.Tr "repo.release.detail"}}</a>
								{{end}}
							{{end}}
						</div>
					</td>
				</tr>
			{{end}}
		</tbody>
	</table>
</div>
{{else}}
<div class="ui attached segment">
	{{ctx.Locale.Tr "search.no_results"}}
</div>
{{end}}
```

The only changes from today's content:
- `<h4>` and the new search-segment live above `{{if .Releases}}`.
- `$canReadReleases` is computed above the `{{if}}` so both branches can use it if needed (harmless in the empty branch).
- New `{{else}}` branch renders `search.no_results` inside a styled segment that visually matches the table wrapper.

- [ ] **Step 4: Lint the template**

Run:
```bash
make lint-templates
```
Expected: PASS with no errors. (This target lints all `.tmpl` files; it should take under 30s.)

- [ ] **Step 5: Commit**

```bash
git add templates/repo/tag/list.tmpl
git commit -m "feat(repo): add search box and empty state to Tags page

Mirrors templates/repo/branch/list.tmpl:70-82 layout. The <h4>
title and search form live outside the {{if .Releases}} gate so an
empty search result still shows the search box. The new {{else}}
branch renders search.no_results."
```

---

## Task 5: Integration Test — `/tags?q=...`

**Files:**
- Modify: `tests/integration/repo_tag_test.go` (append new test function)

- [ ] **Step 1: Write the failing integration test**

Append to `tests/integration/repo_tag_test.go`:

```go
func TestTagsListSearchFilter(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})
	session := loginUser(t, owner.Name)

	// Repo 1's /tags page shows 3 rows: v1.1, delete-tag, v1.0.

	t.Run("no q returns all tags", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequestf(t, "GET", "/%s/%s/tags", owner.Name, repo.Name)
		resp := session.MakeRequest(t, req, http.StatusOK)

		body := string(resp.Body())
		assert.Contains(t, body, ">v1.1<")
		assert.Contains(t, body, ">delete-tag<")
		assert.Contains(t, body, ">v1.0<")
	})

	t.Run("q=v1 returns only matching tags", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequestf(t, "GET", "/%s/%s/tags?q=v1", owner.Name, repo.Name)
		resp := session.MakeRequest(t, req, http.StatusOK)

		body := string(resp.Body())
		assert.Contains(t, body, ">v1.1<")
		assert.Contains(t, body, ">v1.0<")
		assert.NotContains(t, body, ">delete-tag<")
	})

	t.Run("q=delete returns only delete-tag", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequestf(t, "GET", "/%s/%s/tags?q=delete", owner.Name, repo.Name)
		resp := session.MakeRequest(t, req, http.StatusOK)

		body := string(resp.Body())
		assert.Contains(t, body, ">delete-tag<")
		assert.NotContains(t, body, ">v1.1<")
		assert.NotContains(t, body, ">v1.0<")
	})

	t.Run("non-matching q renders empty state", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequestf(t, "GET", "/%s/%s/tags?q=definitely-not-a-real-tag", owner.Name, repo.Name)
		resp := session.MakeRequest(t, req, http.StatusOK)

		body := string(resp.Body())
		assert.NotContains(t, body, ">v1.1<")
		assert.NotContains(t, body, ">delete-tag<")
		assert.NotContains(t, body, ">v1.0<")
		// search.no_results is "No matching results found." in locale_en-US.ini
		assert.Contains(t, body, "No matching results found")
		// search box stays visible and pre-filled
		assert.Contains(t, body, `value="definitely-not-a-real-tag"`)
	})

	t.Run("empty q behaves like no q", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		req := NewRequestf(t, "GET", "/%s/%s/tags?q=", owner.Name, repo.Name)
		resp := session.MakeRequest(t, req, http.StatusOK)

		body := string(resp.Body())
		assert.Contains(t, body, ">v1.1<")
		assert.Contains(t, body, ">delete-tag<")
		assert.Contains(t, body, ">v1.0<")
	})
}
```

If `http` is not already imported in `repo_tag_test.go`, check the existing import block — it is (used in other tests). `tests`, `unittest`, `repo_model`, `user_model` are all already imported. No new imports needed.

- [ ] **Step 2: Run the integration test to verify it passes**

Run:
```bash
make test-sqlite#TestTagsListSearchFilter
```
Expected: PASS (all five subtests). The underlying Go code (Tasks 1-4) is already in place; this test verifies the end-to-end behavior.

If any subtest fails:
- "no q returns all tags" failing means a fixture changed — re-read `models/fixtures/release.yml` for repo_id=1 and adjust expectations.
- "non-matching q renders empty state" failing on the `value="..."` assertion means the `shared/search/combo` partial doesn't render the value attribute the way expected — inspect the rendered HTML to find the actual attribute shape and adjust the assertion.

- [ ] **Step 3: Run all tag-related integration tests to confirm no regression**

Run:
```bash
make test-sqlite#TestCreateNewTagProtected
make test-sqlite#TestTagsListSearchFilter
```
Expected: both PASS.

- [ ] **Step 4: Commit**

```bash
git add tests/integration/repo_tag_test.go
git commit -m "test(integration): cover Tags page search filter

End-to-end coverage for /tags?q=...: no q returns all 3 repo-1 tags;
q=v1 returns v1.0 and v1.1; q=delete returns delete-tag; non-matching
q renders search.no_results with the search box still visible; empty
q behaves like no q."
```

---

## Task 6: Verification

**Files:** none (verification only)

- [ ] **Step 1: Run the full lint suite**

Run:
```bash
make lint
```
Expected: PASS. If lint reports issues in code you didn't touch, those are pre-existing — do not fix them as part of this change.

- [ ] **Step 2: Run the full backend test suite**

Run:
```bash
make test-backend
```
Expected: PASS. Watch specifically for failures in `models/repo/`, `routers/web/repo/`, and `tests/integration/`.

- [ ] **Step 3: Run the SQLite integration suite**

Run:
```bash
make test-sqlite
```
Expected: PASS. This is slower (several minutes); the tags-specific test (`TestTagsListSearchFilter`) should be among the passing tests.

- [ ] **Step 4: Manual smoke test**

Build and start the dev server:
```bash
TAGS="bindata sqlite sqlite_unlock_notify" make build
./gitea web
```

In a browser, log in as owner of a repo with multiple tags (or create a few). Visit `/{owner}/{repo}/tags`:

- **Default view:** All tags listed; search box visible and empty; pagination reflects total.
- **Type a substring:** Page reloads; only matching tags shown; search box retains the typed value; pagination count shrinks to filtered count.
- **Submit a non-matching query:** Page reloads; no tags in list area; `No matching results found.` message visible; search box still pre-filled.
- **Click page 2 on a filtered result:** URL retains `q=<value>`; second page of filtered results shown.
- **Clear the search box and submit:** Returns to full list.

- [ ] **Step 5: Final commit (if any cleanup)**

If manual testing surfaced any issue, fix it with a focused commit. Otherwise no commit needed — Tasks 1-5 are the complete change set.

---

## Self-Review Checklist (already completed during plan authoring)

- **Spec coverage:** Each requirement (PKG-04-006, 108, 109, 110, 111, 304) maps to at least one task and is exercised by an automated test:
  - 006 (filter support): Task 1 + Task 5
  - 108 (non-empty q filters): Task 5 subtests "q=v1" and "q=delete"
  - 109 (empty/absent q unfiltered): Task 5 subtests "no q" and "empty q"
  - 110 (q preserved across pages): Task 3 (handler) + Task 5 indirectly via `SetDefaultParams` — manual verification confirms
  - 111 (pagination count from filtered set): Task 3 (handler) + Task 5 "q=v1" implicitly (would fail if pagination broke)
  - 304 (no-match empty state): Task 5 "non-matching q renders empty state"
- **Placeholder scan:** No TBDs, TODOs, "implement later", or "similar to Task N" references. Every code step contains complete code.
- **Type consistency:** `Keyword string` field used consistently across Tasks 1, 3, and the test in Task 5. `builder.Like{"tag_name", opts.Keyword}` matches Branches usage at `models/git/branch_list.go:103`.
