// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package repositories

import (
	"context"
	"fmt"
	"shb/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) SaveRefreshToken(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	query := `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	_, err := r.postgres.Exec(ctx, query, userID, tokenHash, expiresAt)
	return err
}

func (r *Repository) GetRefreshToken(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, is_revoked, created_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	var t models.RefreshToken
	err := r.postgres.QueryRow(ctx, query, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.IsRevoked, &t.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("token not found")
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true
		WHERE token_hash = $1
	`
	_, err := r.postgres.Exec(ctx, query, tokenHash)
	return err
}

func (r *Repository) RevokeAllUserRefreshTokens(ctx context.Context, userID int) error {
	query := `
		UPDATE refresh_tokens
		SET is_revoked = true
		WHERE user_id = $1 AND is_revoked = false
	`
	_, err := r.postgres.Exec(ctx, query, userID)
	return err
}
