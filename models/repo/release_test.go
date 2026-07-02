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
