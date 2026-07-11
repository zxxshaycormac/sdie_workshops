// Copyright 2023 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package repo

import (
	"testing"

	"code.gitea.io/gitea/models/db"
	"code.gitea.io/gitea/models/unittest"

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

func TestFindReleasesByTagNameKeyword(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	releases, err := db.Find[Release](db.DefaultContext, FindReleasesOptions{
		ListOptions:    db.ListOptionsAll,
		RepoID:         1,
		IncludeDrafts:  true,
		IncludeTags:    true,
		TagNameKeyword: "V1",
	})
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"v1.0", "v1.1"}, releaseTagNames(releases))

	releases, err = db.Find[Release](db.DefaultContext, FindReleasesOptions{
		ListOptions:    db.ListOptionsAll,
		RepoID:         1,
		IncludeDrafts:  true,
		IncludeTags:    true,
		TagNameKeyword: "testing-release",
	})
	assert.NoError(t, err)
	assert.Empty(t, releases, "release titles must not participate in tag-name search")
}

func releaseTagNames(releases []*Release) []string {
	names := make([]string, 0, len(releases))
	for _, release := range releases {
		names = append(names, release.TagName)
	}
	return names
}
