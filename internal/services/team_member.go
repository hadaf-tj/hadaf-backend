// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services

import (
	"context"
	"shb/internal/models"
)

// GetAllTeamMembers retrieves all active team members.
func (s *Service) GetAllTeamMembers(ctx context.Context) ([]*models.TeamMember, error) {
	return s.repo.GetAllTeamMembers(ctx)
}

// GetTeamMemberByID retrieves a team member by ID.
func (s *Service) GetTeamMemberByID(ctx context.Context, id int) (*models.TeamMember, error) {
	return s.repo.GetTeamMemberByID(ctx, id)
}
