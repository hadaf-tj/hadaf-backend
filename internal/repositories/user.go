package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"
)

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	const query = `
		SELECT *
		FROM user
		WHERE phone = :phone
		LIMIT 1;
	`

	args := map[string]interface{}{
		"phone": phone,
	}

	var userDB dbUser
	err := r.postgres.QueryRow(ctx, query, args).Scan(&userDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, myerrors.ErrNotFound
		}
		return nil, fmt.Errorf("get user by phone: %w", err)
	}

	return &models.User{
		ID:        userDB.ID,
		FullName:  userDB.FullName,
		Phone:     userDB.Phone,
		Password:  userDB.Password,
		CreatedAt: userDB.CreatedAt,
		UpdatedAt: userDB.UpdatedAt,
	}, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	const query = `
		INSERT INTO user (phone)
		VALUES (:phone)
		RETURNING id;
	`

	args := map[string]interface{}{
		"phone": user.Phone,
	}

	var id int64
	err := r.postgres.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	user.ID = int(id)
	return nil
}
