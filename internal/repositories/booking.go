package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shb/internal/models"
	"time"
)

func (r *Repository) CreateBooking(ctx context.Context, booking *models.Booking) (int, error) {
	query := `
		INSERT INTO bookings (user_id, need_id, quantity, note, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	var id int
	var createdAt time.Time
	err := r.postgres.QueryRow(ctx, query,
		booking.UserID, booking.NeedID, booking.Quantity, booking.Note, booking.Status,
	).Scan(&id, &createdAt)

	if err != nil {
		return 0, fmt.Errorf("create booking: %w", err)
	}
	booking.ID = id
	booking.CreatedAt = createdAt
	return id, nil
}

func (r *Repository) GetBookingByID(ctx context.Context, id int) (*models.Booking, error) {
	query := `
		SELECT id, user_id, need_id, quantity, note, status, created_at, updated_at, is_deleted, deleted_at
		FROM bookings
		WHERE id = $1 AND is_deleted = false
	`
	var b dbBooking
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.UserID, &b.NeedID, &b.Quantity, &b.Note, &b.Status,
		&b.CreatedAt, &b.UpdatedAt, &b.IsDeleted, &b.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("booking not found: %w", err)
		}
		return nil, fmt.Errorf("get booking by id: %w", err)
	}
	return b.ToDomain(), nil
}

func (r *Repository) GetBookingsByNeed(ctx context.Context, needID int) ([]*models.Booking, error) {
	query := `
		SELECT id, user_id, need_id, quantity, note, status, created_at, updated_at, is_deleted, deleted_at
		FROM bookings
		WHERE need_id = $1 AND is_deleted = false
		ORDER BY created_at DESC
	`
	rows, err := r.postgres.Query(ctx, query, needID)
	if err != nil {
		return nil, fmt.Errorf("get bookings by need: %w", err)
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var b dbBooking
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.NeedID, &b.Quantity, &b.Note, &b.Status,
			&b.CreatedAt, &b.UpdatedAt, &b.IsDeleted, &b.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}
		bookings = append(bookings, b.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return bookings, nil
}

func (r *Repository) GetBookingsByUser(ctx context.Context, userID int) ([]*models.Booking, error) {
	query := `
		SELECT id, user_id, need_id, quantity, note, status, created_at, updated_at, is_deleted, deleted_at
		FROM bookings
		WHERE user_id = $1 AND is_deleted = false
		ORDER BY created_at DESC
	`
	rows, err := r.postgres.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get bookings by user: %w", err)
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var b dbBooking
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.NeedID, &b.Quantity, &b.Note, &b.Status,
			&b.CreatedAt, &b.UpdatedAt, &b.IsDeleted, &b.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}
		bookings = append(bookings, b.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return bookings, nil
}

func (r *Repository) UpdateBookingStatus(ctx context.Context, bookingID int, status string) error {
	query := `
		UPDATE bookings 
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND is_deleted = false
	`
	result, err := r.postgres.Exec(ctx, query, status, bookingID)
	if err != nil {
		return fmt.Errorf("update booking status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("booking not found or already deleted")
	}
	return nil
}

func (r *Repository) GetBookingsByInstitution(ctx context.Context, institutionID int) ([]*models.Booking, error) {
	query := `
		SELECT b.id, b.user_id, b.need_id, b.quantity, b.note, b.status, b.created_at, b.updated_at, b.is_deleted, b.deleted_at
		FROM bookings b
		INNER JOIN needs n ON b.need_id = n.id
		WHERE n.institution_id = $1 AND b.is_deleted = false AND n.is_deleted = false
		ORDER BY b.created_at DESC
	`
	rows, err := r.postgres.Query(ctx, query, institutionID)
	if err != nil {
		return nil, fmt.Errorf("get bookings by institution: %w", err)
	}
	defer rows.Close()

	var bookings []*models.Booking
	for rows.Next() {
		var b dbBooking
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.NeedID, &b.Quantity, &b.Note, &b.Status,
			&b.CreatedAt, &b.UpdatedAt, &b.IsDeleted, &b.DeletedAt,
		); err != nil {
			return nil, fmt.Errorf("scan booking: %w", err)
		}
		bookings = append(bookings, b.ToDomain())
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return bookings, nil
}
