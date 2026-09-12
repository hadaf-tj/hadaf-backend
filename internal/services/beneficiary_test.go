// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services_test

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_GetApplicationQuota(t *testing.T) {
	ctx := context.Background()

	for _, tt := range []struct {
		name      string
		used      int
		remaining int
		allowed   bool
	}{
		{name: "partially used", used: 4, remaining: 6, allowed: true},
		{name: "exactly exhausted", used: 10, remaining: 0, allowed: false},
		{name: "over capacity", used: 12, remaining: 0, allowed: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			svc, d := newTestService(t)
			d.Cache.On("Get", ctx, mock.MatchedBy(func(key string) bool {
				return strings.HasPrefix(key, "applications:pending:")
			})).Return(strconv.Itoa(tt.used), nil)

			quota, err := svc.GetApplicationQuota(ctx)
			require.NoError(t, err)
			require.Equal(t, 10, quota.Limit)
			require.Equal(t, tt.used, quota.Used)
			require.Equal(t, tt.remaining, quota.Remaining)
			require.Equal(t, tt.allowed, quota.Allowed)
			require.GreaterOrEqual(t, quota.Remaining, 0)
		})
	}
}

func TestService_GetApplicationQuota_PropagatesCountError(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	d.Cache.On("Get", ctx, mock.Anything).Return("", errors.New("redis unavailable"))

	quota, err := svc.GetApplicationQuota(ctx)
	require.Nil(t, quota)
	require.ErrorContains(t, err, "get application count")
}
