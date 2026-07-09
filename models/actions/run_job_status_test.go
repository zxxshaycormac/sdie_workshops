// Copyright 2025 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func jobsWithStatuses(statuses ...Status) []*ActionRunJob {
	jobs := make([]*ActionRunJob, len(statuses))
	for i, s := range statuses {
		jobs[i] = &ActionRunJob{Status: s}
	}
	return jobs
}

func TestAggregateJobStatus(t *testing.T) {
	tests := []struct {
		name     string
		jobs     []Status
		expected Status
	}{
		// --- Single-status: each of the 8 enum values in isolation ---
		{name: "single unknown", jobs: []Status{StatusUnknown}, expected: StatusRunning},
		{name: "single success", jobs: []Status{StatusSuccess}, expected: StatusSuccess},
		{name: "single failure", jobs: []Status{StatusFailure}, expected: StatusFailure},
		{name: "single cancelled", jobs: []Status{StatusCancelled}, expected: StatusCancelled},
		{name: "single skipped", jobs: []Status{StatusSkipped}, expected: StatusSkipped},
		{name: "single waiting", jobs: []Status{StatusWaiting}, expected: StatusWaiting},
		{name: "single running", jobs: []Status{StatusRunning}, expected: StatusRunning},
		{name: "single blocked", jobs: []Status{StatusBlocked}, expected: StatusWaiting},

		// --- All-done mixes: terminal statuses only ---
		{name: "all success", jobs: []Status{StatusSuccess, StatusSuccess}, expected: StatusSuccess},
		{name: "success + failure", jobs: []Status{StatusSuccess, StatusFailure}, expected: StatusFailure},
		{name: "success + cancelled", jobs: []Status{StatusSuccess, StatusCancelled}, expected: StatusCancelled},
		{name: "success + skipped", jobs: []Status{StatusSuccess, StatusSkipped}, expected: StatusSuccess},
		{name: "failure + skipped", jobs: []Status{StatusFailure, StatusSkipped}, expected: StatusFailure},
		{name: "cancelled + skipped", jobs: []Status{StatusCancelled, StatusSkipped}, expected: StatusCancelled},
		{name: "all skipped", jobs: []Status{StatusSkipped, StatusSkipped}, expected: StatusSkipped},
		{name: "failure + cancelled (failure wins)", jobs: []Status{StatusFailure, StatusCancelled}, expected: StatusFailure},
		{name: "success + failure + cancelled (failure wins)", jobs: []Status{StatusSuccess, StatusFailure, StatusCancelled}, expected: StatusFailure},

		// --- Non-terminal mixes: at least one job still active ---
		{name: "success + waiting", jobs: []Status{StatusSuccess, StatusWaiting}, expected: StatusWaiting},
		{name: "success + blocked", jobs: []Status{StatusSuccess, StatusBlocked}, expected: StatusWaiting},
		{name: "success + running", jobs: []Status{StatusSuccess, StatusRunning}, expected: StatusRunning},
		{name: "failure + running", jobs: []Status{StatusFailure, StatusRunning}, expected: StatusRunning},
		{name: "cancelled + running", jobs: []Status{StatusCancelled, StatusRunning}, expected: StatusRunning},
		{name: "waiting + blocked", jobs: []Status{StatusWaiting, StatusBlocked}, expected: StatusWaiting},
		{name: "running + blocked", jobs: []Status{StatusRunning, StatusBlocked}, expected: StatusRunning},
		{name: "failure + waiting + blocked", jobs: []Status{StatusFailure, StatusWaiting, StatusBlocked}, expected: StatusWaiting},

		// --- Empty ---
		{name: "empty", jobs: []Status{}, expected: StatusSuccess},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobs := jobsWithStatuses(tt.jobs...)
			got := aggregateJobStatus(jobs)
			assert.Equal(t, tt.expected, got, "aggregateJobStatus(%v)", tt.jobs)
		})
	}
}
