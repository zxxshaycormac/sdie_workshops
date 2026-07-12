// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"testing"

	"code.gitea.io/gitea/models/db"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/models/unittest"
	webhook_module "code.gitea.io/gitea/modules/webhook"

	"github.com/nektos/act/pkg/jobparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertRunAggregatesInitialBlockedJobs(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	workflow := []byte(`
name: test
on: [push]
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - run: echo build
  test:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - run: echo test
`)
	jobs, err := jobparser.Parse(workflow)
	require.NoError(t, err)

	run := &ActionRun{
		Title:         "test blocked run",
		RepoID:        4,
		OwnerID:       1,
		WorkflowID:    "test.yml",
		TriggerUserID: 1,
		Ref:           "refs/heads/master",
		CommitSHA:     "c2d72f548424103f01ee1dc02889c1e2bff816b0",
		Event:         webhook_module.HookEventPush,
		NeedApproval:  true,
		Status:        StatusWaiting,
		Repo:          &repo_model.Repository{ID: 4},
	}
	require.NoError(t, InsertRun(db.DefaultContext, run, jobs))

	persistedRun, err := GetRunByID(db.DefaultContext, run.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusBlocked, run.Status)
	assert.Equal(t, StatusBlocked, persistedRun.Status)

	runJobs, err := GetRunJobsByRunID(db.DefaultContext, run.ID)
	require.NoError(t, err)
	require.Len(t, runJobs, 2)
	for _, job := range runJobs {
		assert.Equal(t, StatusBlocked, job.Status)
	}
}
