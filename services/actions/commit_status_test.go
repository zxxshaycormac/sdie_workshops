// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"testing"

	actions_model "code.gitea.io/gitea/models/actions"
	api "code.gitea.io/gitea/modules/structs"

	"github.com/stretchr/testify/assert"
)

func TestToCommitStatus(t *testing.T) {
	tests := []struct {
		name   string
		status actions_model.Status
		want   api.CommitStatusState
	}{
		{
			name:   "success",
			status: actions_model.StatusSuccess,
			want:   api.CommitStatusSuccess,
		},
		{
			name:   "failure",
			status: actions_model.StatusFailure,
			want:   api.CommitStatusFailure,
		},
		{
			name:   "cancelled",
			status: actions_model.StatusCancelled,
			want:   api.CommitStatusWarning,
		},
		{
			name:   "skipped",
			status: actions_model.StatusSkipped,
			want:   api.CommitStatusWarning,
		},
		{
			name:   "waiting",
			status: actions_model.StatusWaiting,
			want:   api.CommitStatusPending,
		},
		{
			name:   "blocked",
			status: actions_model.StatusBlocked,
			want:   api.CommitStatusPending,
		},
		{
			name:   "running",
			status: actions_model.StatusRunning,
			want:   api.CommitStatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, toCommitStatus(tt.status))
		})
	}
}
