// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package repositories

import (
	"context"
	"fmt"
	"shb/internal/models"
)

// GetVacancyByID retrieves a single active vacancy by ID
func (r *Repository) GetVacancyByID(ctx context.Context, id int) (*models.Vacancy, error) {
	query := `
		SELECT id, title, description, type, experience, workload, is_active, created_at, updated_at
		FROM vacancies
		WHERE id = $1 AND is_deleted = false AND is_active = true
	`
	var v models.Vacancy
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&v.ID, &v.Title, &v.Description, &v.Type, &v.Experience, &v.Workload, &v.IsActive, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get vacancy: %w", err)
	}
	return &v, nil
}
