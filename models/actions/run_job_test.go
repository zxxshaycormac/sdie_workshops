// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregateJobStatus(t *testing.T) {
	tests := []struct {
		name     string
		statuses []Status
		want     Status
	}{
		{
			name: "empty jobs",
			want: StatusUnknown,
		},
		{
			name:     "all success",
			statuses: []Status{StatusSuccess, StatusSuccess},
			want:     StatusSuccess,
		},
		{
			name:     "all skipped",
			statuses: []Status{StatusSkipped, StatusSkipped},
			want:     StatusSkipped,
		},
		{
			name:     "success with skipped",
			statuses: []Status{StatusSuccess, StatusSkipped},
			want:     StatusSuccess,
		},
		{
			name:     "cancelled without pending jobs",
			statuses: []Status{StatusSuccess, StatusCancelled},
			want:     StatusCancelled,
		},
		{
			name:     "failure without cancelled or pending jobs",
			statuses: []Status{StatusFailure, StatusFailure},
			want:     StatusFailure,
		},
		{
			name:     "failure without pending jobs",
			statuses: []Status{StatusSuccess, StatusFailure},
			want:     StatusFailure,
		},
		{
			name:     "cancelled takes precedence over failure",
			statuses: []Status{StatusFailure, StatusCancelled},
			want:     StatusCancelled,
		},
		{
			name:     "running takes precedence while pending",
			statuses: []Status{StatusRunning, StatusWaiting, StatusBlocked, StatusCancelled},
			want:     StatusRunning,
		},
		{
			name:     "waiting takes precedence over blocked",
			statuses: []Status{StatusWaiting, StatusBlocked},
			want:     StatusWaiting,
		},
		{
			name:     "all blocked",
			statuses: []Status{StatusBlocked, StatusBlocked},
			want:     StatusBlocked,
		},
		{
			name:     "blocked is visible when no jobs are waiting or running",
			statuses: []Status{StatusBlocked, StatusSkipped},
			want:     StatusBlocked,
		},
		{
			name:     "unknown remains unknown when no known state applies",
			statuses: []Status{StatusUnknown},
			want:     StatusUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobs := make([]*ActionRunJob, 0, len(tt.statuses))
			for _, status := range tt.statuses {
				jobs = append(jobs, &ActionRunJob{Status: status})
			}

			assert.Equal(t, tt.want, aggregateJobStatus(jobs))
		})
	}
}
