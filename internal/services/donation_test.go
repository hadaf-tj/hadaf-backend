// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services_test

import (
	"context"
	"errors"
	"testing"

	"shb/internal/models"
	"shb/pkg/constants"

	"github.com/stretchr/testify/require"
)

func TestService_GetMyDonations(t *testing.T) {
	ctx := context.Background()
	page := &models.DonationHistoryPage{
		Items: []*models.DonationHistoryItem{{
			ID:            101,
			CampaignID:    11,
			CampaignTitle: "Urgent treatment",
			Amount:        50,
			Currency:      "TJS",
			Status:        models.DonationStatusConfirmed,
		}},
		Total:  1,
		Limit:  20,
		Offset: 0,
	}

	svc, d := newTestService(t)
	d.Repo.On("GetDonationsByUser", ctx, 7, 20, 0).Return(page, nil).Once()

	got, err := svc.GetMyDonations(ctx, 7, 20, 0)
	require.NoError(t, err)
	require.Same(t, page, got)
}

func TestService_GetMyDonationsAllowsEmptyHistory(t *testing.T) {
	ctx := context.Background()
	empty := &models.DonationHistoryPage{Items: make([]*models.DonationHistoryItem, 0), Limit: 20}

	svc, d := newTestService(t)
	d.Repo.On("GetDonationsByUser", ctx, 7, 20, 0).Return(empty, nil).Once()

	got, err := svc.GetMyDonations(ctx, 7, 20, 0)
	require.NoError(t, err)
	require.Empty(t, got.Items)
}

func TestService_GetMyDonationsValidatesInput(t *testing.T) {
	svc, _ := newTestService(t)
	ctx := context.Background()

	for _, tt := range []struct {
		name                  string
		userID, limit, offset int
	}{
		{name: "missing authenticated user", userID: 0, limit: 20, offset: 0},
		{name: "zero limit", userID: 7, limit: 0, offset: 0},
		{name: "limit exceeds maximum", userID: 7, limit: constants.MaxPageLimit + 1, offset: 0},
		{name: "negative offset", userID: 7, limit: 20, offset: -1},
	} {
		t.Run(tt.name, func(t *testing.T) {
			page, err := svc.GetMyDonations(ctx, tt.userID, tt.limit, tt.offset)
			require.Nil(t, page)
			require.Error(t, err)
		})
	}
}

func TestService_GetMyDonationsWrapsRepositoryError(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	d.Repo.On("GetDonationsByUser", ctx, 7, 20, 0).Return(nil, errors.New("database unavailable")).Once()

	page, err := svc.GetMyDonations(ctx, 7, 20, 0)
	require.Nil(t, page)
	require.ErrorContains(t, err, "get donations by user")
}
