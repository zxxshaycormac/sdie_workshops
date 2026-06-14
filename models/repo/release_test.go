// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"testing"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/models/unittest"
	"code.gitea.io/gitea/modules/optional"

	"github.com/stretchr/testify/assert"
)

func TestMigrate_InsertReleases(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	a := &Attachment{
		UUID: "a0eebc91-9c0c-4ef7-bb6e-6bb9bd380a12",
	}
	r := &Release{
		Attachments: []*Attachment{a},
	}

	err := InsertReleases(db.DefaultContext, r)
	assert.NoError(t, err)
}

func TestFindReleasesByKeyword(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// repo 57 has tags: v1.0, v1.1, v2.0, non-existing-target-branch, empty-target-branch (all with sha1)
	baseOpts := FindReleasesOptions{
		IncludeDrafts: true,
		IncludeTags:   true,
		HasSha1:       optional.Some(true),
		RepoID:        57,
	}

	t.Run("no keyword returns all", func(t *testing.T) {
		releases, err := db.Find[Release](db.DefaultContext, baseOpts)
		assert.NoError(t, err)
		assert.Len(t, releases, 5)
	})

	t.Run("substring v1 matches v1.0 and v1.1", func(t *testing.T) {
		opts := baseOpts
		opts.Keyword = "v1"
		releases, err := db.Find[Release](db.DefaultContext, opts)
		assert.NoError(t, err)
		tagNames := releaseTagNames(releases)
		assert.ElementsMatch(t, []string{"v1.0", "v1.1"}, tagNames)
	})

	t.Run("substring is case-insensitive", func(t *testing.T) {
		opts := baseOpts
		opts.Keyword = "V1"
		releases, err := db.Find[Release](db.DefaultContext, opts)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"v1.0", "v1.1"}, releaseTagNames(releases))
	})

	t.Run("no matches returns empty", func(t *testing.T) {
		opts := baseOpts
		opts.Keyword = "nonexistent-substring"
		releases, err := db.Find[Release](db.DefaultContext, opts)
		assert.NoError(t, err)
		assert.Empty(t, releases)
	})
}

func releaseTagNames(releases []*Release) []string {
	names := make([]string, 0, len(releases))
	for _, r := range releases {
		names = append(names, r.TagName)
	}
	return names
}
