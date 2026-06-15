# Repository Tags Page Search Filter — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add keyword search to the repository Tags page (`/{owner}/{repo}/tags`) so users can filter tags by name, matching the existing Branches page search UX.

**Architecture:** Reuse the Branches search pattern end-to-end. The `release` table already has a `lower_tag_name` column (maintained by the model), so we add a `Keyword` field to `FindReleasesOptions`, append a `builder.Like{"lower_tag_name", strings.ToLower(keyword)}` clause in `ToConds()`, read `q` from the query string in the `TagsList` handler, and render a `shared/search/combo` form in `tag/list.tmpl`. No new DB columns, no new API, no JS.

**Tech Stack:** Go (XORM + `xorm.io/builder`), Go HTML templates (`html/template`), Fomantic UI, INI locale files, Gitea integration test framework (`tests/integration/`).

**Spec:** [`docs/superpowers/specs/2026-06-15-tags-search-filter-design.md`](../specs/2026-06-15-tags-search-filter-design.md)

---

## File Structure

| File | Change | Responsibility |
|---|---|---|
| `models/repo/release.go` | Modify | Add `Keyword` field to `FindReleasesOptions`; append LIKE clause in `ToConds()` |
| `routers/web/repo/release.go` | Modify | `TagsList` handler reads `q`, sets `ctx.Data["Keyword"]`, passes keyword into opts |
| `templates/repo/tag/list.tmpl` | Modify | Render search form between header partial and tag table |
| `options/locale/locale_en-US.ini` | Modify | Add `tag_kind = Search tags...` key under `[search]` |
| `tests/integration/repo_tag_test.go` | Modify | Append `TestTagsSearch` integration test (hit + miss cases) |

No new files. No new packages. No migrations.

---

## Task 1: Write the failing integration test

**Files:**
- Modify: `tests/integration/repo_tag_test.go` (append new `TestTagsSearch` function at end of file)

- [ ] **Step 1: Append the test function**

Add this code at the end of `tests/integration/repo_tag_test.go` (after the closing `}` of `TestRepushTag` at line 166):

```go
func TestTagsSearch(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user2")

	t.Run("Hit", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		// "v1.1" exists in repo1 fixtures (models/fixtures/release.yml).
		req := NewRequest(t, "GET", "/user2/repo1/tags?q=v1.1")
		resp := session.MakeRequest(t, req, http.StatusOK)
		htmlDoc := NewHTMLParser(t, resp.Body)
		tagLinks := htmlDoc.doc.Find(".tag-list-row-link")
		assert.GreaterOrEqual(t, tagLinks.Length(), 1, "expected at least one tag matching 'v1.1'")
		// Every visible tag link's text must contain the keyword (proves filtering works).
		tagLinks.Each(func(i int, s *goquery.Selection) {
			assert.Contains(t, s.Text(), "v1.1", "non-matching tag leaked into filtered result set")
		})
	})

	t.Run("Miss", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		req := NewRequest(t, "GET", "/user2/repo1/tags?q=zzz-not-a-real-tag-xyz")
		resp := session.MakeRequest(t, req, http.StatusOK)
		htmlDoc := NewHTMLParser(t, resp.Body)
		assert.Equal(t, 0, htmlDoc.doc.Find(".tag-list-row-link").Length(), "expected zero tag rows for a non-matching keyword")
	})

	t.Run("EmptyKeywordPreservesBehavior", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		// No q param: same as before the feature shipped.
		req := NewRequest(t, "GET", "/user2/repo1/tags")
		resp := session.MakeRequest(t, req, http.StatusOK)
		htmlDoc := NewHTMLParser(t, resp.Body)
		assert.GreaterOrEqual(t, htmlDoc.doc.Find(".tag-list-row-link").Length(), 1, "expected tags to render when no keyword is supplied")
	})
}
```

Also add `"github.com/PuerkitoBio/goquery"` to the import block at the top of the file (alphabetical position: after `"code.gitea.io/gitea/tests"`). The `goquery` import is required because the test calls `s.Text()` on `*goquery.Selection` inside the `Each` callback.

- [ ] **Step 2: Run the test to verify it fails**

Run:
```bash
make test-sqlite#TestTagsSearch
```

Expected: FAIL. The `Hit` and `Miss` subtests will fail because the `q` parameter is currently ignored — all tags render regardless of keyword, so the `Hit` assertion `all visible tags should match the keyword` fails (e.g., `v1.0` appears alongside `v1.1`), and the `Miss` assertion `expected zero tag rows` fails.

The `EmptyKeywordPreservesBehavior` subtest should PASS even before implementation (it tests the unchanged baseline).

- [ ] **Step 3: Commit the failing test**

```bash
git add tests/integration/repo_tag_test.go
git commit -m "test: add failing integration test for tags page search filter"
```

---

## Task 2: Add `Keyword` field and LIKE clause to `FindReleasesOptions`

**Files:**
- Modify: `models/repo/release.go:229-238` (struct) and `models/repo/release.go:240-266` (`ToConds`)

- [ ] **Step 1: Add the `Keyword` field to the struct**

In `models/repo/release.go`, locate the `FindReleasesOptions` struct (lines 229-238):

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
}
```

Append `Keyword string` as a new field after `HasSha1`:

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
	Keyword       string
}
```

- [ ] **Step 2: Add the LIKE clause to `ToConds()`**

In the same file, locate `ToConds()` (lines 240-266). It currently ends with the `HasSha1` block and a `return cond`:

```go
	if opts.HasSha1.Has() {
		if opts.HasSha1.Value() {
			cond = cond.And(builder.Neq{"sha1": ""})
		} else {
			cond = cond.And(builder.Eq{"sha1": ""})
		}
	}
	return cond
}
```

Insert a new `Keyword` clause between the `HasSha1` block and `return cond`:

```go
	if opts.HasSha1.Has() {
		if opts.HasSha1.Value() {
			cond = cond.And(builder.Neq{"sha1": ""})
		} else {
			cond = cond.And(builder.Eq{"sha1": ""})
		}
	}
	if opts.Keyword != "" {
		cond = cond.And(builder.Like{"lower_tag_name", strings.ToLower(opts.Keyword)})
	}
	return cond
}
```

The `"strings"` package is already imported at line 14 — verify by grepping if uncertain: `grep -n '"strings"' models/repo/release.go` should print `14:	"strings"`.

- [ ] **Step 3: Verify the package compiles**

Run:
```bash
go build ./models/repo/...
```

Expected: no output (success). If it fails, the most likely cause is a typo in the struct field or a missing import.

- [ ] **Step 4: Run the integration test (still red — handler not wired)**

Run:
```bash
make test-sqlite#TestTagsSearch
```

Expected: still FAIL on `Hit` and `Miss`. The model layer is in place but the handler doesn't read `q`, so `Keyword` is always empty and `ToConds()` skips the LIKE clause. Confirming this proves the model change alone is insufficient and the test is genuinely exercising the full stack.

- [ ] **Step 5: Commit**

```bash
git add models/repo/release.go
git commit -m "feat: support keyword filter in FindReleasesOptions"
```

---

## Task 3: Wire `q` parameter through the `TagsList` handler

**Files:**
- Modify: `routers/web/repo/release.go:204-251` (`TagsList`)

- [ ] **Step 1: Read the `q` parameter**

In `routers/web/repo/release.go`, locate `TagsList` (line 204). Find the `listOptions` block (lines 215-224):

```go
	listOptions := db.ListOptions{
		Page:     ctx.FormInt("page"),
		PageSize: ctx.FormInt("limit"),
	}
	if listOptions.PageSize == 0 {
		listOptions.PageSize = setting.Repository.Release.DefaultPagingNum
	}
	if listOptions.PageSize > setting.API.MaxResponseItems {
		listOptions.PageSize = setting.API.MaxResponseItems
	}
```

Insert immediately after this block (before the `opts := repo_model.FindReleasesOptions{` literal at line 226):

```go
	keyword := ctx.FormString("q")
```

- [ ] **Step 2: Pass keyword into the opts literal**

In the same function, the `opts` literal is at lines 226-234:

```go
	opts := repo_model.FindReleasesOptions{
		ListOptions: listOptions,
		// for the tags list page, show all releases with real tags (having real commit-id),
		// the drafts should also be included because a real tag might be used as a draft.
		IncludeDrafts: true,
		IncludeTags:   true,
		HasSha1:       optional.Some(true),
		RepoID:        ctx.Repo.Repository.ID,
	}
```

Append `Keyword: keyword,` as a new field (after `RepoID:`):

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

- [ ] **Step 3: Expose keyword to the template**

In the same function, locate `ctx.Data["Releases"] = releases` (line 242). Add a new line right after it:

```go
	ctx.Data["Releases"] = releases
	ctx.Data["Keyword"] = keyword
```

(Pagination at lines 244-247 already calls `pager.SetDefaultParams(ctx)` which auto-carries `q` through page links — no change needed there.)

- [ ] **Step 4: Verify the package compiles**

Run:
```bash
go build ./routers/web/repo/...
```

Expected: no output (success).

- [ ] **Step 5: Run the integration test (still red — template missing)**

Run:
```bash
make test-sqlite#TestTagsSearch
```

Expected: still FAIL on `Hit` and `Miss`. The handler now reads `q` and passes it to the model, so the filter SHOULD work at the SQL layer — but the test asserts on HTML output, and since the test was failing before Task 3, this run confirms the test still depends on the template path. (If `Hit` now passes here, that's fine — it means the SQL filter is working. The `Miss` case may still fail if there are tag rows being rendered with a stale keyword. Either way, proceed to Task 4.)

- [ ] **Step 6: Commit**

```bash
git add routers/web/repo/release.go
git commit -m "feat: read q keyword in TagsList handler"
```

---

## Task 4: Add locale key and search form to the template

**Files:**
- Modify: `options/locale/locale_en-US.ini:180` (add new key)
- Modify: `templates/repo/tag/list.tmpl` (insert search form between lines 6 and 7)

- [ ] **Step 1: Add the `tag_kind` locale key**

In `options/locale/locale_en-US.ini`, find the `[search]` section. `branch_kind` is at line 178, `commit_kind` at line 179. Add `tag_kind` as a new line immediately after `commit_kind` (becomes line 180):

Before:
```ini
branch_kind = Search branches...
commit_kind = Search commits...
```

After:
```ini
branch_kind = Search branches...
commit_kind = Search commits...
tag_kind = Search tags...
```

- [ ] **Step 2: Insert the search form into `templates/repo/tag/list.tmpl`**

The template currently opens (lines 1-7):

```html
{{template "base/head" .}}
<div role="main" aria-label="{{.Title}}" class="page-content repository tags">
	{{template "repo/header" .}}
	<div class="ui container">
		{{template "base/alert" .}}
		{{template "repo/release_tag_header" .}}
		{{if .Releases}}
```

Insert the search form **between** the `release_tag_header` line (line 6) and the `{{if .Releases}}` line (line 7), so the form renders even when the search returns zero results:

```html
{{template "base/head" .}}
<div role="main" aria-label="{{.Title}}" class="page-content repository tags">
	{{template "repo/header" .}}
	<div class="ui container">
		{{template "base/alert" .}}
		{{template "repo/release_tag_header" .}}
		<div class="ui attached segment">
			<form class="ignore-dirty" method="get">
				{{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}
			</form>
		</div>
		{{if .Releases}}
```

The `Value` and `Placeholder` keys are the contract documented in `templates/shared/search/combo.tmpl` — same pattern as `templates/repo/branch/list.tmpl:76-80`.

- [ ] **Step 3: Run the integration test (green)**

Run:
```bash
make test-sqlite#TestTagsSearch
```

Expected: PASS for all three subtests (`Hit`, `Miss`, `EmptyKeywordPreservesBehavior`).

If `Hit` fails with "all visible tags should match the keyword", verify that `LowerTagName` values in `models/fixtures/release.yml` for repo_id=1 are actually lowercase — they should be (`v1.1`, `v1.0`, etc.).

If `Miss` fails (tag rows appear for a non-matching keyword), verify the LIKE clause in Task 2 was inserted **inside** `ToConds()` (not accidentally inside the `HasSha1.Has()` branch).

- [ ] **Step 4: Commit**

```bash
git add options/locale/locale_en-US.ini templates/repo/tag/list.tmpl
git commit -m "feat: add keyword search form to tags page"
```

---

## Task 5: Build, lint, and manual browser verification

**Files:** None modified (verification only).

- [ ] **Step 1: Run the full build**

Run:
```bash
make build
```

Expected: completes without error. This recompiles bindata (which embeds the modified template and locale file) and produces a fresh `./gitea` binary.

- [ ] **Step 2: Lint the touched surfaces**

Run:
```bash
make lint-go
make lint-templates
make lint-css
```

Expected: no errors related to the changes. Pre-existing warnings unrelated to the diff are acceptable.

If `lint-templates` complains about `shared/search/combo`, double-check the template call against `templates/repo/branch/list.tmpl:78` — the signature is `dict "Value" <val> "Placeholder" <key>`.

- [ ] **Step 3: Run the broader test suite to catch regressions**

Run:
```bash
make test-sqlite#TestViewReleases
make test-sqlite#TestViewReleasesNoLogin
make test-sqlite#TestCreateNewTagProtected
make test-sqlite#TestRepushTag
```

Expected: all PASS. These tests exercise the Releases page and tag creation paths that share `FindReleasesOptions`. If any fail, the `Keyword` field addition broke backward compatibility (most likely culprit: a test that constructs `FindReleasesOptions{}` literally and now needs the field omitted — Go's zero value handles this automatically, so failure here would indicate a real regression).

- [ ] **Step 4: Manual browser verification**

Start the dev server with the rebuilt binary against a test DB:

```bash
GITEA_RUN_MODE=dev ./gitea web
```

In a browser, navigate to a repository with multiple tags (e.g. `http://localhost:3000/user2/repo1/tags`). Confirm:

1. The search input appears between the page header and the tag table.
2. Placeholder text reads "Search tags...".
3. Typing `v1.1` and pressing Enter returns only `v1.1` rows.
4. Typing `V1.1` (uppercase) returns the same results (case-insensitive).
5. Typing `zzz-nonexistent` returns the empty state with the search form still visible.
6. The URL after submit ends with `?q=v1.1` (and pagination links, if visible, preserve the `q` parameter).
7. Visiting `/user2/repo1/tags` with no `q` parameter renders all tags as before.

- [ ] **Step 5: Commit any fix-ups (if any)**

If the manual check surfaced a bug that required code changes, fix and commit:

```bash
git add <files>
git commit -m "fix: <description>"
```

If no fix-ups are needed, no commit — Task 4 already produced the final code.

---

## Done Criteria

- [ ] All three `TestTagsSearch` subtests pass under `make test-sqlite`.
- [ ] Existing tag/release integration tests still pass (no regressions).
- [ ] `make build`, `make lint-go`, `make lint-templates`, `make lint-css` clean for the touched files.
- [ ] Manual browser check confirms: search filters tags by name, case-insensitive, empty keyword is no-op, search form persists across pagination and zero-result states.
- [ ] Five commits land on the branch in this order:
  1. `test: add failing integration test for tags page search filter`
  2. `feat: support keyword filter in FindReleasesOptions`
  3. `feat: read q keyword in TagsList handler`
  4. `feat: add keyword search form to tags page`
  5. (optional) `fix: <description>` if Step 5 surfaced anything.
