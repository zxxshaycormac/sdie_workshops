// Copyright 2021 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"code.gitea.io/gitea/models"
	auth_model "code.gitea.io/gitea/models/auth"
	"code.gitea.io/gitea/models/db"
	git_model "code.gitea.io/gitea/models/git"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/models/unittest"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/git"
	api "code.gitea.io/gitea/modules/structs"
	"code.gitea.io/gitea/services/release"
	"code.gitea.io/gitea/tests"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNewTagProtected(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})

	t.Run("Code", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()

		err := release.CreateNewTag(git.DefaultContext, owner, repo, "master", "t-first", "first tag")
		assert.NoError(t, err)

		err = release.CreateNewTag(git.DefaultContext, owner, repo, "master", "v-2", "second tag")
		assert.Error(t, err)
		assert.True(t, models.IsErrProtectedTagName(err))

		err = release.CreateNewTag(git.DefaultContext, owner, repo, "master", "v-1.1", "third tag")
		assert.NoError(t, err)
	})

	t.Run("Git", func(t *testing.T) {
		onGiteaRun(t, func(t *testing.T, u *url.URL) {
			httpContext := NewAPITestContext(t, owner.Name, repo.Name)

			dstPath := t.TempDir()

			u.Path = httpContext.GitPath()
			u.User = url.UserPassword(owner.Name, userPassword)

			doGitClone(dstPath, u)(t)

			_, _, err := git.NewCommand(git.DefaultContext, "tag", "v-2").RunStdString(&git.RunOpts{Dir: dstPath})
			assert.NoError(t, err)

			_, _, err = git.NewCommand(git.DefaultContext, "push", "--tags").RunStdString(&git.RunOpts{Dir: dstPath})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "Tag v-2 is protected")
		})
	})

	t.Run("GitTagForce", func(t *testing.T) {
		onGiteaRun(t, func(t *testing.T, u *url.URL) {
			httpContext := NewAPITestContext(t, owner.Name, repo.Name)

			dstPath := t.TempDir()

			u.Path = httpContext.GitPath()
			u.User = url.UserPassword(owner.Name, userPassword)

			doGitClone(dstPath, u)(t)

			_, _, err := git.NewCommand(git.DefaultContext, "tag", "v-1.1", "-m", "force update", "--force").RunStdString(&git.RunOpts{Dir: dstPath})
			require.NoError(t, err)

			_, _, err = git.NewCommand(git.DefaultContext, "push", "--tags").RunStdString(&git.RunOpts{Dir: dstPath})
			require.NoError(t, err)

			_, _, err = git.NewCommand(git.DefaultContext, "tag", "v-1.1", "-m", "force update v2", "--force").RunStdString(&git.RunOpts{Dir: dstPath})
			require.NoError(t, err)

			_, _, err = git.NewCommand(git.DefaultContext, "push", "--tags").RunStdString(&git.RunOpts{Dir: dstPath})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "the tag already exists in the remote")

			_, _, err = git.NewCommand(git.DefaultContext, "push", "--tags", "--force").RunStdString(&git.RunOpts{Dir: dstPath})
			require.NoError(t, err)
			req := NewRequestf(t, "GET", "/%s/releases/tag/v-1.1", repo.FullName())
			resp := MakeRequest(t, req, http.StatusOK)
			htmlDoc := NewHTMLParser(t, resp.Body)
			tagsTab := htmlDoc.Find(".release-list-title")
			assert.Contains(t, tagsTab.Text(), "force update v2")
		})
	})

	// Cleanup
	releases, err := db.Find[repo_model.Release](db.DefaultContext, repo_model.FindReleasesOptions{
		IncludeTags: true,
		TagNames:    []string{"v-1", "v-1.1"},
		RepoID:      repo.ID,
	})
	assert.NoError(t, err)

	for _, release := range releases {
		_, err = db.DeleteByID[repo_model.Release](db.DefaultContext, release.ID)
		assert.NoError(t, err)
	}

	protectedTags, err := git_model.GetProtectedTags(db.DefaultContext, repo.ID)
	assert.NoError(t, err)

	for _, protectedTag := range protectedTags {
		err = git_model.DeleteProtectedTag(db.DefaultContext, protectedTag)
		assert.NoError(t, err)
	}
}

func TestRepushTag(t *testing.T) {
	onGiteaRun(t, func(t *testing.T, u *url.URL) {
		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
		owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})
		session := loginUser(t, owner.LowerName)
		token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteRepository)

		httpContext := NewAPITestContext(t, owner.Name, repo.Name)

		dstPath := t.TempDir()

		u.Path = httpContext.GitPath()
		u.User = url.UserPassword(owner.Name, userPassword)

		doGitClone(dstPath, u)(t)

		// create and push a tag
		_, _, err := git.NewCommand(git.DefaultContext, "tag", "v2.0").RunStdString(&git.RunOpts{Dir: dstPath})
		assert.NoError(t, err)
		_, _, err = git.NewCommand(git.DefaultContext, "push", "origin", "--tags", "v2.0").RunStdString(&git.RunOpts{Dir: dstPath})
		assert.NoError(t, err)
		// create a release for the tag
		createdRelease := createNewReleaseUsingAPI(t, session, token, owner, repo, "v2.0", "", "Release of v2.0", "desc")
		assert.False(t, createdRelease.IsDraft)
		// delete the tag
		_, _, err = git.NewCommand(git.DefaultContext, "push", "origin", "--delete", "v2.0").RunStdString(&git.RunOpts{Dir: dstPath})
		assert.NoError(t, err)
		// query the release by API and it should be a draft
		req := NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/releases/tags/%s", owner.Name, repo.Name, "v2.0"))
		resp := MakeRequest(t, req, http.StatusOK)
		var respRelease *api.Release
		DecodeJSON(t, resp, &respRelease)
		assert.True(t, respRelease.IsDraft)
		// re-push the tag
		_, _, err = git.NewCommand(git.DefaultContext, "push", "origin", "--tags", "v2.0").RunStdString(&git.RunOpts{Dir: dstPath})
		assert.NoError(t, err)
		// query the release by API and it should not be a draft
		req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/releases/tags/%s", owner.Name, repo.Name, "v2.0"))
		resp = MakeRequest(t, req, http.StatusOK)
		DecodeJSON(t, resp, &respRelease)
		assert.False(t, respRelease.IsDraft)
	})
}

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

	t.Run("CaseInsensitive", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		// lower_tag_name column stores lowercase; uppercase query must still match.
		req := NewRequest(t, "GET", "/user2/repo1/tags?q=V1.1")
		resp := session.MakeRequest(t, req, http.StatusOK)
		htmlDoc := NewHTMLParser(t, resp.Body)
		tagLinks := htmlDoc.doc.Find(".tag-list-row-link")
		assert.GreaterOrEqual(t, tagLinks.Length(), 1, "uppercase query should still match via lower_tag_name")
		tagLinks.Each(func(i int, s *goquery.Selection) {
			assert.Contains(t, strings.ToLower(s.Text()), "v1.1", "case-insensitive match failed")
		})
	})

	t.Run("PaginationWithKeyword", func(t *testing.T) {
		defer tests.PrintCurrentTest(t)()
		// Regression guard: when the keyword narrows the result set to fit in a
		// single page, the pagination widget MUST NOT render. Before the fix,
		// TagsList used ctx.Data["NumTags"] (unfiltered total = 3 for repo1),
		// so a 1-result search with limit=1 would still render pagination for
		// 3 pages — clicking pages 2 or 3 returned empty results. After the fix,
		// TagsList uses the filtered count (1), so pagination is suppressed.
		// repo1 fixtures: exactly one tag (v1.0) matches q=v1.0 with HasSha1=true.
		req := NewRequest(t, "GET", "/user2/repo1/tags?q=v1.0&limit=1")
		resp := session.MakeRequest(t, req, http.StatusOK)
		htmlDoc := NewHTMLParser(t, resp.Body)

		// One matching tag fits within page size 1.
		visibleCount := htmlDoc.doc.Find(".tag-list-row-link").Length()
		assert.Equal(t, 1, visibleCount, "exactly one tag should match v1.0")

		// templates/base/paginate.tmpl renders pagination only when TotalPages > 1.
		// With filtered count=1 and limit=1, TotalPages=1, so no pagination renders.
		// The buggy code (unfiltered count=3) would have rendered pagination here.
		pageLinks := htmlDoc.doc.Find(".pagination a")
		assert.Equal(t, 0, pageLinks.Length(), "pagination must not render when filtered count fits in one page")
	})
}
