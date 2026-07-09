// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"testing"
)

func jobsWithStatuses(statuses ...Status) []*ActionRunJob {
	jobs := make([]*ActionRunJob, 0, len(statuses))
	for i, s := range statuses {
		jobs = append(jobs, &ActionRunJob{ID: int64(i + 1), Status: s})
	}
	return jobs
}

func TestAggregateJobStatus(t *testing.T) {
	cases := []struct {
		name string
		jobs []Status
		want Status
	}{
		// Single-job runs: each status maps to itself.
		{"single success", []Status{StatusSuccess}, StatusSuccess},
		{"single failure", []Status{StatusFailure}, StatusFailure},
		{"single cancelled", []Status{StatusCancelled}, StatusCancelled},
		{"single skipped", []Status{StatusSkipped}, StatusSkipped},
		{"single waiting", []Status{StatusWaiting}, StatusWaiting},
		{"single running", []Status{StatusRunning}, StatusRunning},
		{"single blocked", []Status{StatusBlocked}, StatusBlocked},
		{"single unknown falls back to running", []Status{StatusUnknown}, StatusRunning},

		// Terminal aggregates (all jobs done).
		{"all skipped", []Status{StatusSkipped, StatusSkipped}, StatusSkipped},
		{"success and skipped", []Status{StatusSuccess, StatusSkipped}, StatusSuccess},
		{"failure dominates success", []Status{StatusFailure, StatusSuccess}, StatusFailure},
		{"failure dominates cancelled", []Status{StatusFailure, StatusCancelled}, StatusFailure},
		{"success and cancelled is cancelled", []Status{StatusSuccess, StatusCancelled}, StatusCancelled},
		{"cancelled and skipped is cancelled", []Status{StatusCancelled, StatusSkipped}, StatusCancelled},

		// In-progress aggregates (not all jobs done).
		{"all blocked", []Status{StatusBlocked, StatusBlocked}, StatusBlocked},
		{"blocked dominates waiting", []Status{StatusBlocked, StatusWaiting}, StatusBlocked},
		{"running dominates blocked", []Status{StatusRunning, StatusBlocked}, StatusRunning},
		{"running dominates waiting", []Status{StatusRunning, StatusWaiting}, StatusRunning},
		{"waiting with done job", []Status{StatusWaiting, StatusSuccess}, StatusWaiting},
		{"blocked with done job", []Status{StatusBlocked, StatusSuccess}, StatusBlocked},

		// Degenerate case.
		{"empty jobs", []Status{}, StatusSkipped},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := aggregateJobStatus(jobsWithStatuses(tc.jobs...))
			if got != tc.want {
				t.Errorf("aggregateJobStatus(%v) = %s, want %s", tc.jobs, got, tc.want)
			}
		})
	}
}
