// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package actions

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetStatusInfoList(t *testing.T) {
	got := GetStatusInfoList(context.Background())

	want := []Status{
		StatusSuccess,
		StatusFailure,
		StatusCancelled,
		StatusSkipped,
		StatusWaiting,
		StatusRunning,
		StatusBlocked,
	}
	require.Len(t, got, len(want))

	for i, status := range want {
		assert.Equal(t, int(status), got[i].Status)
		assert.Equal(t, status.String(), got[i].DisplayedStatus)
	}
}
