// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package pgx

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"shb/pkg/constants"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPgxPool() (*pgxpool.Pool, error) {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		return nil, fmt.Errorf("POSTGRES_HOST environment variable is required")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		return nil, fmt.Errorf("POSTGRES_PORT environment variable is required")
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		return nil, fmt.Errorf("POSTGRES_USER environment variable is required")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		return nil, fmt.Errorf("POSTGRES_PASSWORD environment variable is required")
	}

	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		return nil, fmt.Errorf("POSTGRES_DB environment variable is required")
	}

	sslMode := os.Getenv("POSTGRES_SSL_MODE")
	if sslMode == "" {
		sslMode = constants.SSLModeDisable
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		url.QueryEscape(password),
		host,
		port,
		dbname,
		sslMode,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	// Verify the connection is reachable.
	if err = pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return pool, nil
}
