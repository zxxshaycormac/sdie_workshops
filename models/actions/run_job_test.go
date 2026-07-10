// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAggregateJobStatus(t *testing.T) {
	tests := []struct {
		name string
		jobs []*ActionRunJob
		want Status
	}{
		{
			name: "all skipped",
			jobs: []*ActionRunJob{
				{Status: StatusSkipped},
				{Status: StatusSkipped},
			},
			want: StatusSkipped,
		},
		{
			name: "cancelled and skipped",
			jobs: []*ActionRunJob{
				{Status: StatusCancelled},
				{Status: StatusSkipped},
			},
			want: StatusCancelled,
		},
		{
			name: "failure and cancelled",
			jobs: []*ActionRunJob{
				{Status: StatusFailure},
				{Status: StatusCancelled},
			},
			want: StatusFailure,
		},
		{
			name: "waiting and blocked",
			jobs: []*ActionRunJob{
				{Status: StatusWaiting},
				{Status: StatusBlocked},
			},
			want: StatusWaiting,
		},
		{
			name: "all blocked",
			jobs: []*ActionRunJob{
				{Status: StatusBlocked},
				{Status: StatusBlocked},
			},
			want: StatusBlocked,
		},
		{
			name: "running and waiting",
			jobs: []*ActionRunJob{
				{Status: StatusRunning},
				{Status: StatusWaiting},
			},
			want: StatusRunning,
		},
		{
			name: "all success",
			jobs: []*ActionRunJob{
				{Status: StatusSuccess},
				{Status: StatusSuccess},
			},
			want: StatusSuccess,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, aggregateJobStatus(tt.jobs))
		})
	}
}

func TestGetStatusInfoList(t *testing.T) {
	statusInfoList := GetStatusInfoList(t.Context())
	statuses := make([]Status, 0, len(statusInfoList))
	for _, statusInfo := range statusInfoList {
		statuses = append(statuses, Status(statusInfo.Status))
	}

	require.NotContains(t, statuses, StatusUnknown)
	assert.Contains(t, statuses, StatusSuccess)
	assert.Contains(t, statuses, StatusFailure)
	assert.Contains(t, statuses, StatusCancelled)
	assert.Contains(t, statuses, StatusSkipped)
	assert.Contains(t, statuses, StatusWaiting)
	assert.Contains(t, statuses, StatusRunning)
	assert.Contains(t, statuses, StatusBlocked)
}
