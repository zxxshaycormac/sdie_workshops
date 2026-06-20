# Tags Page Keyword Search Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add Branches-style keyword search to the repository Tags page (`/{owner}/{repo}/tags?q=...`) so users can filter tags by name.

**Architecture:** Server-side substring filter via a new `Keyword` field on `FindReleasesOptions` (emitting `builder.Like("tag_name", kw)`), wired through the existing `TagsList` handler, rendered via the existing `shared/search/combo` template partial. No new packages, no new routes, no migrations, no JS.

**Tech Stack:** Go 1.22.x, XORM (`xorm.io/builder`), Chi router, Go HTML templates, Fomantic UI, Gitea's `unittest` + `tests/integration` frameworks.

**Spec:** `docs/superpowers/specs/2026-06-20-tags-search-design.md`

**Repo conventions (from `CLAUDE.md`, `design.md`, directory `CLAUDE.md` files):**
- SQL built only via `builder.*` or XORM DSL — never interpolate user input.
- Handlers: `func Xxx(ctx *context.Context)`, errors via `ctx.ServerError` / `ctx.NotFound`.
- Templates: every user-visible string via `{{ctx.Locale.Tr "key"}}`; no inline `<script>`.
- Locale: add to `options/locale/locale_en-US.ini` only; community translates the rest.
- Tests: DB-backed tests use `unittest.PrepareTestDatabase()`; integration tests use `defer tests.PrepareTestEnv(t)()`.
- Commits: no `Co-Authored-By` trailer in this repo.

---

## File Structure

| File | Change | Responsibility |
|---|---|---|
| `models/repo/release.go` | Modify lines 229-266 (add field + clause) | Data-layer keyword filter |
| `models/repo/release_test.go` | Append new test | Unit test for the `Keyword` clause |
| `routers/web/repo/release.go` | Modify `TagsList` (lines 205-251) | Handler reads `q`, sets `Keyword`, computes filtered count |
| `templates/repo/tag/list.tmpl` | Modify lines 7-63 (insert search form + empty state) | UI for search box + no-match message |
| `options/locale/locale_en-US.ini` | Add 2 keys (lines ~170, ~991+) | Locale strings |
| `tests/integration/repo_tag_test.go` | Append new test | Integration test for `/tags?q=...` |

Each task below produces a self-contained, committable change.

---

## Task 1: Data-layer `Keyword` field — failing test first (Red)

**Files:**
- Modify: `models/repo/release_test.go` (append after `TestMigrate_InsertReleases`, line 27)

- [ ] **Step 1: Read the existing test file to confirm structure**

Run: `head -30 /Users/genewu/github/gitea/models/repo/release_test.go`
Expected: shows the existing `TestMigrate_InsertReleases` function and the existing imports (`models/db`, `models/unittest`, `github.com/stretchr/testify/assert`). We will reuse these imports.

- [ ] **Step 2: Confirm fixture data for the assertion**

Run: `grep -A1 "repo_id: 1" /Users/genewu/github/gitea/models/fixtures/release.yml | grep tag_name`
Expected output contains `v1.0`, `v1.1`, `delete-tag`, `draft-release` — these are the repo_id=1 tags our test will filter against.

- [ ] **Step 3: Append the failing test**

Edit `models/repo/release_test.go`. The current imports are:
```go
import (
	"testing"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/models/unittest"

	"github.com/stretchr/testify/assert"
)
```

Replace the import block with:
```go
import (
	"sort"
	"testing"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/models/unittest"
	"code.gitea.io/gitea/modules/optional"

	"github.com/stretchr/testify/assert"
)
```

Then append this test after the closing brace of `TestMigrate_InsertReleases` (currently line 27):

```go
func TestFindReleasesOptions_KeywordFilter(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	opts := repo_model.FindReleasesOptions{
		IncludeDrafts: true,
		IncludeTags:   true,
		HasSha1:       optional.Some(true),
		RepoID:        1,
	}

	// Substring match: "v1" matches v1.0 and v1.1 in repo 1 fixtures.
	// Does NOT match delete-tag or draft-release.
	opts.Keyword = "v1"
	releases, err := db.Find[repo_model.Release](db.DefaultContext, opts)
	assert.NoError(t, err)

	gotNames := make([]string, 0, len(releases))
	for _, r := range releases {
		gotNames = append(gotNames, r.TagName)
	}
	sort.Strings(gotNames)
	assert.Equal(t, []string{"v1.0", "v1.1"}, gotNames, "keyword 'v1' should match v1.0 and v1.1 only")

	// Non-matching keyword returns no rows.
	opts.Keyword = "nonexistent-tag-name-xyz"
	releases, err = db.Find[repo_model.Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Empty(t, releases)

	// Empty keyword disables the filter — count > 0 (sanity check that the
	// empty-keyword path is not accidentally excluding everything).
	opts.Keyword = ""
	count, err := db.Count[repo_model.Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0))
}
```

You also need to add `repo_model` to the imports. The final import block:
```go
import (
	"sort"
	"testing"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/models/unittest"
	"code.gitea.io/gitea/modules/optional"

	"github.com/stretchr/testify/assert"
)
```

**Note on the import alias:** the package declares `package repo` and other files in the same package already import it as `repo_model`. Check an existing sibling file (e.g. `models/repo/repo_test.go`) — if it imports `code.gitea.io/gitea/models/repo` as `repo_model`, use the same alias. If the package is being tested from inside itself (which it is — same `package repo`), you cannot import the package itself; **instead, drop the `repo_model.` prefix and use `FindReleasesOptions`, `Release` directly** since the test lives in `package repo`.

**Corrected version for in-package test:**
```go
import (
	"sort"
	"testing"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/models/unittest"
	"code.gitea.io/gitea/modules/optional"

	"github.com/stretchr/testify/assert"
)

func TestFindReleasesOptions_KeywordFilter(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	opts := FindReleasesOptions{
		IncludeDrafts: true,
		IncludeTags:   true,
		HasSha1:       optional.Some(true),
		RepoID:        1,
	}

	opts.Keyword = "v1"
	releases, err := db.Find[Release](db.DefaultContext, opts)
	assert.NoError(t, err)

	gotNames := make([]string, 0, len(releases))
	for _, r := range releases {
		gotNames = append(gotNames, r.TagName)
	}
	sort.Strings(gotNames)
	assert.Equal(t, []string{"v1.0", "v1.1"}, gotNames)

	opts.Keyword = "nonexistent-tag-name-xyz"
	releases, err = db.Find[Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Empty(t, releases)

	opts.Keyword = ""
	count, err := db.Count[Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0))
}
```

- [ ] **Step 4: Run the test to verify it FAILS (red)**

Run: `cd /Users/genewu/github/gitea && go test -run TestFindReleasesOptions_KeywordFilter ./models/repo/`
Expected: compile error `unknown field 'Keyword' in struct literal of type FindReleasesOptions` — that's the red state we want. The field doesn't exist yet.

If you instead see a runtime failure (field exists but logic wrong), stop and re-check Step 3.

- [ ] **Step 5: Do NOT commit yet** — implementation comes in Task 2.

---

## Task 2: Implement the `Keyword` clause (Green) + commit

**Files:**
- Modify: `models/repo/release.go` (lines 229-266)

- [ ] **Step 1: Read the current FindReleasesOptions and ToConds**

Run: `sed -n '228,270p' /Users/genewu/github/gitea/models/repo/release.go`
Expected: shows the struct definition with fields `RepoID, IncludeDrafts, IncludeTags, IsPreRelease, IsDraft, TagNames, HasSha1`, and `ToConds()` with the existing `builder.In("tag_name", opts.TagNames)` clause.

- [ ] **Step 2: Add the `Keyword` field to the struct**

Edit `models/repo/release.go`. Replace this block (around lines 229-238):

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

with:

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
	Keyword       string // substring match on tag_name; empty disables the filter
	HasSha1       optional.Option[bool] // useful to find draft releases which are created with existing tags
}
```

- [ ] **Step 3: Add the `Like` clause to ToConds**

In the same file, find this block in `ToConds()` (around line 249):

```go
	if len(opts.TagNames) > 0 {
		cond = cond.And(builder.In("tag_name", opts.TagNames))
	}
```

Immediately after it (still inside `ToConds()`, before the `IsPreRelease` check), insert:

```go
	if opts.Keyword != "" {
		cond = cond.And(builder.Like("tag_name", opts.Keyword))
	}
```

- [ ] **Step 4: Run the unit test to verify it PASSES (green)**

Run: `cd /Users/genewu/github/gitea && go test -run TestFindReleasesOptions_KeywordFilter ./models/repo/`
Expected: `ok code.gitea.io/gitea/models/repo X.XXXs` — test passes.

If it fails with `v1` matching `delete-tag` or `draft-release`, your `Like` clause is not being applied — re-check Step 3.

- [ ] **Step 5: Run the full models/repo test suite to verify no regression**

Run: `cd /Users/genewu/github/gitea && go test ./models/repo/`
Expected: all tests pass. Existing callers of `FindReleasesOptions` (e.g. `GetTagNamesByRepoID`, `Releases` handler tests) should be unaffected because `Keyword` defaults to `""`.

- [ ] **Step 6: Commit**

```bash
cd /Users/genewu/github/gitea
git add models/repo/release.go models/repo/release_test.go
git commit -m "$(cat <<'EOF'
feat(repo): add Keyword filter to FindReleasesOptions

Adds a Keyword string field that emits builder.Like("tag_name", kw)
in ToConds(), mirroring FindBranchOptions.Keyword. Existing callers
are unaffected because the field defaults to empty and the new
clause is skipped when empty.
EOF
)"
```

---

## Task 3: Integration test for `/tags?q=...` — failing first (Red)

**Files:**
- Modify: `tests/integration/repo_tag_test.go` (append after `TestRepushTag`, currently line 166)

- [ ] **Step 1: Confirm the test infrastructure imports**

Run: `head -30 /Users/genewu/github/gitea/tests/integration/repo_tag_test.go`
Expected: imports include `tests`, `unittest`, `repo_model`, `user_model`, `release`, `git`, `net/http`, `net/url`. We will reuse these.

- [ ] **Step 2: Append the failing integration test**

Edit `tests/integration/repo_tag_test.go`. Append this function after `TestRepushTag`:

```go
func TestTagsListSearch(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	// Create three tags with distinct name prefixes so we can verify the
	// filter narrows the list. Clean them up at the end so the test is hermetic.
	tagNames := []string{"alpha-search-1", "alpha-search-2", "beta-search-1"}
	for _, name := range tagNames {
		err := release.CreateNewTag(git.DefaultContext, owner, repo, "master", name, "test tag "+name)
		assert.NoError(t, err)
	}
	defer func() {
		// Clean up the release rows created by CreateNewTag. Best-effort;
		// errors are ignored because the test is ending.
		releases, err := db.Find[repo_model.Release](db.DefaultContext, repo_model.FindReleasesOptions{
			RepoID:   repo.ID,
			TagNames: tagNames,
		})
		if err != nil {
			return
		}
		for _, r := range releases {
			_, _ = db.DeleteByID[repo_model.Release](db.DefaultContext, r.ID)
		}
	}()

	t.Run("substring match", func(t *testing.T) {
		req := NewRequestf(t, "GET", "/%s/tags?q=alpha-search", repo.FullName())
		resp := MakeRequest(t, req, http.StatusOK)
		doc := NewHTMLParser(t, resp.Body)

		// Both alpha-search-* tags appear in the tags table.
		text := doc.Find(`#tags-table`).Text()
		assert.Contains(t, text, "alpha-search-1")
		assert.Contains(t, text, "alpha-search-2")

		// The beta tag does NOT appear.
		assert.NotContains(t, text, "beta-search-1")

		// The search input is pre-filled with the query.
		inputVal, ok := doc.Find(`input[name="q"]`).Attr("value")
		assert.True(t, ok)
		assert.Equal(t, "alpha-search", inputVal)
	})

	t.Run("no match shows empty state", func(t *testing.T) {
		req := NewRequestf(t, "GET", "/%s/tags?q=no-such-tag-xyz-123", repo.FullName())
		resp := MakeRequest(t, req, http.StatusOK)
		doc := NewHTMLParser(t, resp.Body)

		text := doc.Find(`.page-content`).Text()
		assert.Contains(t, text, "No tags match")
	})

	t.Run("absent q preserves full list", func(t *testing.T) {
		req := NewRequestf(t, "GET", "/%s/tags", repo.FullName())
		resp := MakeRequest(t, req, http.StatusOK)
		doc := NewHTMLParser(t, resp.Body)

		text := doc.Find(`#tags-table`).Text()
		assert.Contains(t, text, "alpha-search-1")
		assert.Contains(t, text, "beta-search-1")
	})
}
```

**Imports required** — add to the existing import block if not already present:
- `"code.gitea.io/gitea/models/db"` — for `db.DeleteByID`, `db.Find`, `db.DefaultContext`

- [ ] **Step 3: Verify the test compiles**

Run: `cd /Users/genewu/github/gitea && go build ./tests/integration/`
Expected: compiles cleanly. If `strings` is referenced in the chosen cleanup variant, ensure `"strings"` is in the imports.

- [ ] **Step 4: Run the integration test to verify it FAILS (red)**

Run: `cd /Users/genewu/github/gitea && go test -run TestTagsListSearch ./tests/integration/`
Expected: at minimum, the "no match shows empty state" sub-test fails because the locale key doesn't exist yet (the empty-state branch isn't rendered). The "substring match" sub-test may also fail because the handler doesn't yet read `q`.

The test must compile and run — the *assertions* are what fail. If it fails to compile, fix the imports/helpers before proceeding.

- [ ] **Step 5: Do NOT commit yet** — implementation comes in Task 4.

---

## Task 4: Handler + template + locale (Green) + commit

**Files:**
- Modify: `options/locale/locale_en-US.ini` (two locations)
- Modify: `templates/repo/tag/list.tmpl` (insert search form + empty state)
- Modify: `routers/web/repo/release.go` (lines 205-251)

- [ ] **Step 1: Add locale keys**

Edit `options/locale/locale_en-US.ini`.

Find the `[search]` section (around line 162). After the line `branch_kind = Search branches...` (around line 175), add:

```ini
tag_kind = Search tags...
```

Find the `[repo]` section (starts at line 991) and locate the existing `repo.release.*` keys. Run: `grep -n "release.tags\b\|release.delete_tag\b" /Users/genewu/github/gitea/options/locale/locale_en-US.ini` to find them, then add immediately after one of those existing keys:

```ini
release.tags.no_match = No tags match "%s"
```

The exact line will vary; place it adjacent to other `release.tags.*` keys for translator readability.

- [ ] **Step 2: Update the template — add the search form**

Edit `templates/repo/tag/list.tmpl`. Currently lines 7-12 are:

```html
		{{if .Releases}}
		<h4 class="ui top attached header">
			<div class="five wide column tw-flex tw-items-center">
				{{svg "octicon-tag" 16 "tw-mr-1"}}{{ctx.Locale.Tr "repo.release.tags"}}
			</div>
		</h4>
```

Insert the search form BEFORE the `{{if .Releases}}` block. The new structure is:

```html
		<div class="ui attached segment">
			<form class="ignore-dirty" method="get">
				{{template "shared/search/combo" dict "Value" .Keyword "Placeholder" (ctx.Locale.Tr "search.tag_kind")}}
			</form>
		</div>

		{{if .Releases}}
		<h4 class="ui top attached header">
			<div class="five wide column tw-flex tw-items-center">
				{{svg "octicon-tag" 16 "tw-mr-1"}}{{ctx.Locale.Tr "repo.release.tags"}}
			</div>
		</h4>
```

- [ ] **Step 3: Update the template — add the no-match empty state**

In the same file, find the closing `{{end}}` of the `{{if .Releases}}` block (around line 61, right before `{{template "base/paginate" .}}`). Replace:

```html
			</div>
		</div>
		{{end}}

		{{template "base/paginate" .}}
```

with:

```html
			</div>
		</div>
		{{else if .Keyword}}
		<div class="ui attached segment center">
			{{ctx.Locale.Tr "repo.release.tags.no_match" .Keyword}}
		</div>
		{{end}}

		{{template "base/paginate" .}}
```

The `{{else if .Keyword}}` arm fires only when the keyword is non-empty AND no releases matched. When `q` is absent and the repo simply has zero tags, the original empty behavior (no message) is preserved.

- [ ] **Step 4: Update the `TagsList` handler**

Edit `routers/web/repo/release.go`. Replace the entire `TagsList` function body (currently lines 205-251):

```go
// TagsList render tags list page
func TagsList(ctx *context.Context) {
	ctx.Data["PageIsTagList"] = true
	ctx.Data["Title"] = ctx.Tr("repo.release.tags")
	ctx.Data["IsViewBranch"] = false
	ctx.Data["IsViewTag"] = true
	// Disable the showCreateNewBranch form in the dropdown on this page.
	ctx.Data["CanCreateBranch"] = false
	ctx.Data["HideBranchesInDropdown"] = true
	ctx.Data["CanCreateRelease"] = ctx.Repo.CanWrite(unit.TypeReleases) && !ctx.Repo.Repository.IsArchived

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

	opts := repo_model.FindReleasesOptions{
		ListOptions: listOptions,
		// for the tags list page, show all releases with real tags (having real commit-id),
		// the drafts should also be included because a real tag might be used as a draft.
		IncludeDrafts: true,
		IncludeTags:   true,
		HasSha1:       optional.Some(true),
		RepoID:        ctx.Repo.Repository.ID,
	}

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

	ctx.Data["PageIsViewCode"] = !ctx.Repo.Repository.UnitEnabled(ctx, unit.TypeReleases)
	ctx.HTML(http.StatusOK, tplTagsList)
}
```

with:

```go
// TagsList render tags list page
func TagsList(ctx *context.Context) {
	ctx.Data["PageIsTagList"] = true
	ctx.Data["Title"] = ctx.Tr("repo.release.tags")
	ctx.Data["IsViewBranch"] = false
	ctx.Data["IsViewTag"] = true
	// Disable the showCreateNewBranch form in the dropdown on this page.
	ctx.Data["CanCreateBranch"] = false
	ctx.Data["HideBranchesInDropdown"] = true
	ctx.Data["CanCreateRelease"] = ctx.Repo.CanWrite(unit.TypeReleases) && !ctx.Repo.Repository.IsArchived

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

	kw := ctx.FormString("q")

	opts := repo_model.FindReleasesOptions{
		ListOptions: listOptions,
		// for the tags list page, show all releases with real tags (having real commit-id),
		// the drafts should also be included because a real tag might be used as a draft.
		IncludeDrafts: true,
		IncludeTags:   true,
		HasSha1:       optional.Some(true),
		RepoID:        ctx.Repo.Repository.ID,
		Keyword:       kw,
	}

	releases, err := db.Find[repo_model.Release](ctx, opts)
	if err != nil {
		ctx.ServerError("GetReleasesByRepoID", err)
		return
	}

	ctx.Data["Releases"] = releases
	ctx.Data["Keyword"] = kw

	// Under search, the middleware-set NumTags (which counts ALL tags for the
	// repo header badge) is wrong for pagination — compute the filtered total.
	// When the keyword is empty, reuse NumTags to avoid the extra COUNT(*).
	var total int64
	if kw != "" {
		total, err = db.Count[repo_model.Release](ctx, opts)
		if err != nil {
			ctx.ServerError("CountReleases", err)
			return
		}
	} else {
		total = ctx.Data["NumTags"].(int64)
	}
	pager := context.NewPagination(int(total), opts.PageSize, opts.Page, 5)
	pager.SetDefaultParams(ctx)
	ctx.Data["Page"] = pager

	ctx.Data["PageIsViewCode"] = !ctx.Repo.Repository.UnitEnabled(ctx, unit.TypeReleases)
	ctx.HTML(http.StatusOK, tplTagsList)
}
```

- [ ] **Step 5: Verify it compiles**

Run: `cd /Users/genewu/github/gitea && go build ./routers/web/repo/`
Expected: compiles cleanly.

- [ ] **Step 6: Run the integration test to verify it PASSES (green)**

Run: `cd /Users/genewu/github/gitea && go test -run TestTagsListSearch ./tests/integration/`
Expected: `ok code.gitea.io/gitea/tests/integration X.XXXs` — all three sub-tests pass.

If "no match shows empty state" still fails, verify the locale key path matches exactly: `repo.release.tags.no_match`. The `[repo]` section in the ini maps to the `repo.` prefix; the key is then `release.tags.no_match` within that section.

- [ ] **Step 7: Run the existing release_test.go to verify no regression**

Run: `cd /Users/genewu/github/gitea && go test -run TestCreateNewTagProtected ./tests/integration/`
Expected: still passes. The Releases handler is untouched; only `TagsList` changed.

- [ ] **Step 8: Commit**

```bash
cd /Users/genewu/github/gitea
git add options/locale/locale_en-US.ini templates/repo/tag/list.tmpl routers/web/repo/release.go tests/integration/repo_tag_test.go
git commit -m "$(cat <<'EOF'
feat(repo): add keyword search to Tags page

Adds a Branches-style search box to /{owner}/{repo}/tags. The
TagsList handler reads ?q=<kw>, passes it through FindReleasesOptions
to a builder.Like filter on tag_name, computes a filtered count for
pagination, and re-populates the input via ctx.Data["Keyword"]. A
no-match query renders a guided empty state.

Tags-only scope; the Releases page is unchanged but the data-layer
Keyword field is reusable for future parity.
EOF
)"
```

---

## Task 5: Lint, format, and final verification

**Files:** none modified in this task.

- [ ] **Step 1: Run the Go formatter and module tidy**

Run: `cd /Users/genewu/github/gitea && make fmt && make tidy`
Expected: no changes (or only whitespace changes which are fine).

- [ ] **Step 2: Run the linters**

Run: `cd /Users/genewu/github/gitea && make lint-go && make lint-templates`
Expected: clean. If `lint-templates` flags the new `{{else if .Keyword}}` arm, double-check the syntax matches an existing `{{else if ...}}` use elsewhere in the templates tree (e.g. `grep -rn "{{else if" templates/repo/`).

- [ ] **Step 3: Run the full backend test suite**

Run: `cd /Users/genewu/github/gitea && make test-backend`
Expected: all tests pass. If anything fails that is unrelated to tags (e.g. flaky integration test), note it and proceed — only failures touching the files in this plan block completion.

- [ ] **Step 4: Build the server**

Run: `cd /Users/genewu/github/gitea && TAGS="bindata sqlite sqlite_unlock_notify" make build`
Expected: builds cleanly, producing `./gitea`.

- [ ] **Step 5: Manual browser verification**

Start the server in dev mode:
```bash
cd /Users/genewu/github/gitea
GITEA_RUN_MODE=dev ./gitea web
```

In a browser, navigate to a test repo's tags page (e.g. `http://localhost:3000/<owner>/<repo>/tags`). Verify each of these manually:

1. **Empty query (today's behavior preserved):** visit `/tags` with no `?q=`. All tags are listed; pagination shows the unfiltered total.
2. **Matching query:** type a substring that matches some tag names (e.g. `v1`) into the search box and submit. Only matching tags are shown. The search input retains the value. Pagination reflects the filtered count.
3. **Non-matching query:** type a string that matches nothing (e.g. `zzz`). The page renders the "No tags match 'zzz'" message; no pagination links.
4. **Pagination carries the query:** with a matching query that spans multiple pages, click page 2. The URL should retain `?q=<kw>&page=2`. The filtered set, not the full set, continues.
5. **Clearing the search:** empty the input and submit. `?q=` is sent; the handler treats empty as "no filter"; full list returns.

If any step fails, return to Task 4 and fix the relevant piece before claiming Done.

- [ ] **Step 6: Commit lint/format changes if any**

If Steps 1-2 produced any changes:
```bash
cd /Users/genewu/github/gitea
git add -p
git commit -m "chore: apply fmt/tidy to tags search implementation"
```

If nothing changed, skip this step.

---

## Definition of Done (per spec §8.4)

- [ ] Task 1 + 2 complete — `FindReleasesOptions.Keyword` covered by unit test
- [ ] Task 3 + 4 complete — `TagsList` search flow covered by integration test
- [ ] `make test-backend` passes
- [ ] `make lint-go` and `make lint-templates` pass
- [ ] Manual browser verification (Task 5 Step 5) — all 5 scenarios behave correctly
- [ ] Locale keys added to `locale_en-US.ini` only
- [ ] Acceptance scenarios from spec §9 pass (covered by Task 3 sub-tests)

---

## Open Follow-ups (not blocking Done)

- **Migrate design into the OpenSpec change directory.** The scaffold at `openspec/changes/tags-search-filter/` is currently empty (only `.openspec.yaml`). After Done, populate `proposal.md`, `design.md`, `specs/` deltas, and `tasks.md` from this plan so the change directory is no longer a stub. Can be a separate commit.
- **E2E (Playwright) test.** Only if a Branches-search E2E exists to mirror — check `tests/e2e/` during execution. The integration test in Task 3 covers the contract; E2E is optional polish.
- **Cross-database case-sensitivity note.** SQL `LIKE` is case-insensitive on MySQL/SQLite and case-sensitive on PostgreSQL by default. The unit test in Task 1 deliberately avoids asserting mixed-case queries. If the project later requires consistent case-insensitive matching across backends, that is a separate change (ILIKE on Postgres, or `LOWER()` both sides).
