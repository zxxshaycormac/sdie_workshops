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

func TestFindReleasesOptions_KeywordFilter(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	// Repo 1 has 3 sha1-bearing releases matching {IncludeTags, IncludeDrafts}={true,true}:
	// v1.1, delete-tag, v1.0. draft-release is excluded by HasSha1 because it has no sha1.

	t.Run("empty keyword returns all matching tags", func(t *testing.T) {
		releases, err := db.Find[Release](db.DefaultContext, FindReleasesOptions{
			RepoID:        1,
			IncludeTags:   true,
			IncludeDrafts: true,
			HasSha1:       optional.Some(true),
			ListOptions:   db.ListOptions{ListAll: true},
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
		names := make([]string, 0, len(releases))
		for _, r := range releases {
			names = append(names, r.TagName)
		}
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
