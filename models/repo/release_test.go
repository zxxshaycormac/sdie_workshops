// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"sort"
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

	opts := FindReleasesOptions{
		IncludeDrafts: true,
		IncludeTags:   true,
		HasSha1:       optional.Some(true),
		RepoID:        1,
	}

	// Substring match: "v1" matches v1.0 and v1.1 in repo 1 fixtures.
	// Does NOT match delete-tag or draft-release.
	opts.Keyword = "v1"
	releases, err := db.Find[Release](db.DefaultContext, opts)
	assert.NoError(t, err)

	gotNames := make([]string, 0, len(releases))
	for _, r := range releases {
		gotNames = append(gotNames, r.TagName)
	}
	sort.Strings(gotNames)
	assert.Equal(t, []string{"v1.0", "v1.1"}, gotNames)

	// Non-matching keyword returns no rows.
	opts.Keyword = "nonexistent-tag-name-xyz"
	releases, err = db.Find[Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Empty(t, releases)

	// Empty keyword disables the filter — count > 0 (sanity check that the
	// empty-keyword path is not accidentally excluding everything).
	opts.Keyword = ""
	count, err := db.Count[Release](db.DefaultContext, opts)
	assert.NoError(t, err)
	assert.Greater(t, count, int64(0))
}
