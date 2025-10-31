package repositories

import (
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
