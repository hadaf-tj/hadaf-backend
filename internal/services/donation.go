// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services

import (
	"context"
	"fmt"

	"shb/internal/models"
	"shb/pkg/constants"
	"shb/pkg/myerrors"

	"github.com/rs/zerolog"
)

// GetMyDonations retrieves the authenticated donor's recorded donation
// history. The caller supplies userID only from authenticated request context.
func (s *Service) GetMyDonations(
	ctx context.Context,
	userID, limit, offset int,
) (*models.DonationHistoryPage, error) {
	if userID <= 0 {
		return nil, myerrors.NewBadRequestErr("invalid user ID")
	}
	if limit <= 0 || limit > constants.MaxPageLimit {
		return nil, myerrors.NewBadRequestErr("invalid limit")
	}
	if offset < 0 {
		return nil, myerrors.NewBadRequestErr("invalid offset")
	}

	page, err := s.repo.GetDonationsByUser(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get donations by user: %w", err)
	}

	zerolog.Ctx(ctx).Debug().
		Int("user_id", userID).
		Int("count", len(page.Items)).
		Int("limit", limit).
		Int("offset", offset).
		Msg("donation history fetched")

	return page, nil
}
