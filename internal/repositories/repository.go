// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package repositories

import (
	"context"
	"fmt"
	"shb/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository struct {
	postgres *pgxpool.Pool
	logger   *zerolog.Logger
}

func NewRepository(postgresConn *pgxpool.Pool, log *zerolog.Logger) *Repository {
	return &Repository{postgres: postgresConn, logger: log}
}

// CreateBeneficiary persists a new beneficiary record to the database.
func (r *Repository) CreateBeneficiary(ctx context.Context, b *models.Beneficiary) error {
	if b == nil {
		return fmt.Errorf("beneficiary is nil")
	}
	if b.UserID <= 0 {
		return fmt.Errorf("invalid user_id: %d", b.UserID)
	}
	if b.FullName == "" {
		return fmt.Errorf("full_name is required")
	}

	query := `
		INSERT INTO beneficiaries
		(user_id, full_name, birth_date, diagnosis, city, region, contact_phone, status)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, created_at
	`

	err := r.postgres.QueryRow(ctx, query,
		b.UserID,
		b.FullName,
		b.BirthDate,
		b.Diagnosis,
		b.City,
		b.Region,
		b.ContactPhone,
		b.Status,
	).Scan(&b.ID, &b.CreatedAt)

	if err != nil {
		return fmt.Errorf("insert beneficiary: %w", err)
	}

	return nil
}

// CountPendingThisMonth returns the count of pending beneficiaries created this month.
func (r *Repository) CountPendingThisMonth(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM beneficiaries
		WHERE status = 'pending'
		  AND created_at >= date_trunc('month', CURRENT_DATE)
		  AND is_deleted = false
	`

	var count int
	if err := r.postgres.QueryRow(ctx, query).Scan(&count); err != nil {
		r.logger.Error().Err(err).Msg("count pending beneficiaries failed")
		return 0, fmt.Errorf("count pending beneficiaries: %w", err)
	}

	return count, nil
}
