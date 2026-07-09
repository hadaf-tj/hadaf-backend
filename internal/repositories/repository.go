// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package repositories

import (
	"context"
	"fmt"

	"shb/internal/services"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository struct {
	services.IRepository
	postgres *pgxpool.Pool
	logger   *zerolog.Logger
}

func NewRepository(postgresConn *pgxpool.Pool, log *zerolog.Logger) *Repository {
	return &Repository{postgres: postgresConn, logger: log}
}

// Ping verifies the database connection is alive. Used by readiness probes.
func (r *Repository) Ping(ctx context.Context) error {
	return r.postgres.Ping(ctx)
}

// PoolStats exposes the underlying pgx pool statistics for metrics collection.
func (r *Repository) PoolStats() *pgxpool.Stat {
	return r.postgres.Stat()
}

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
