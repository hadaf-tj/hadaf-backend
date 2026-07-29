// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services

import (
	"context"
	"shb/internal/models"
)

// GetAllVacancies retrieves all active vacancies.
func (s *Service) GetAllVacancies(ctx context.Context) ([]*models.Vacancy, error) {
	return s.repo.GetAllVacancies(ctx)
}

// GetVacancyByID retrieves a vacancy by ID.
func (s *Service) GetVacancyByID(ctx context.Context, id int) (*models.Vacancy, error) {
	return s.repo.GetVacancyByID(ctx, id)
}
