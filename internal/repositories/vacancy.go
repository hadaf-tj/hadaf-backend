// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package repositories

import (
	"context"
	"shb/internal/models"
)

// GetAllVacancies retrieves external facing active vacancies
func (r *Repository) GetAllVacancies(ctx context.Context) ([]*models.Vacancy, error) {
	query := `
		SELECT id, title, description, type, experience, workload, is_active, created_at, updated_at
		FROM vacancies
		WHERE is_deleted = false AND is_active = true
		ORDER BY created_at DESC
	`
	rows, err := r.postgres.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vacancies []*models.Vacancy
	for rows.Next() {
		var v models.Vacancy
		err := rows.Scan(
			&v.ID, &v.Title, &v.Description, &v.Type, &v.Experience, &v.Workload, &v.IsActive, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		vacancies = append(vacancies, &v)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return vacancies, nil
}
