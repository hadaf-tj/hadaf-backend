package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"
	"time"

	"github.com/rs/zerolog/log"
)

func (r *Repository) SaveOTP(ctx context.Context, o *models.OTP) error {
	const query = `
    INSERT INTO otp (receiver, method, otp_code, sent_at, expires_at, attempt, is_verified)
    VALUES ($1, $2, $3, $4, $5, $6, $7);
`
	_, err := r.postgres.Exec(ctx, query,
		o.Receiver,
		o.Method,
		o.OTPCode,
		o.SentAt,
		o.ExpiresAt,
		o.Attempt,
		o.IsVerified,
	)

	if err != nil {
		return fmt.Errorf("failed to save otp: %w", err)
	}

	return nil
}

func (r *Repository) GetOTP(ctx context.Context, receiver string) (*models.OTP, error) {
	// Мы убедились по логам, что UTC время работает корректно.
	// Возвращаем проверку: expires_at > $2
	const query = `
    SELECT id, attempt, receiver, method, otp_code, is_verified, sent_at, expires_at, updated_at, is_deleted, deleted_at
    FROM otp
    WHERE receiver = $1
      AND is_verified = false
      AND expires_at > $2
    ORDER BY id DESC
    LIMIT 1;
`
	var otpDB dbOtp

	// Берем текущее время в UTC, чтобы сравнивать корректно
	currentTime := time.Now().UTC()

	// ВАЖНО: В запросе 2 параметра ($1, $2), и мы передаем ровно 2 аргумента.
	err := r.postgres.QueryRow(ctx, query, receiver, currentTime).Scan(
		&otpDB.ID, &otpDB.Attempt, &otpDB.Receiver, &otpDB.Method, &otpDB.OTPCode,
		&otpDB.IsVerified, &otpDB.SentAt, &otpDB.ExpiresAt, &otpDB.UpdatedAt,
		&otpDB.IsDeleted, &otpDB.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Если не нашли - логируем детали для отладки
			log.Warn().
				Str("receiver", receiver).
				Time("check_time_utc", currentTime).
				Msg("OTP not found (expired or wrong receiver)")
			return nil, myerrors.ErrNotFound
		}
		return nil, fmt.Errorf("otp query error: %w", err)
	}

	return otpDB.ToDomain(), nil
}

func (r *Repository) MarkOTPAsVerified(ctx context.Context, otpID int) error {
	const query = `UPDATE otp SET is_verified=true, updated_at=NOW() WHERE id=$1`

	result, err := r.postgres.Exec(ctx, query, otpID)
	if err != nil {
		return fmt.Errorf("failed to mark otp as verified: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for otp id: %d", otpID)
	}

	return nil
}

func (r *Repository) IncreaseOTPAttempt(ctx context.Context, otpID int, phone string) error {
	const query = `
    UPDATE otp
    SET attempt = attempt + 1,
        updated_at = NOW()
    WHERE id = $1 AND receiver = $2;
`
	result, err := r.postgres.Exec(ctx, query, otpID, phone)
	if err != nil {
		return fmt.Errorf("receiver:%s failed to increase otp attempt: %w", phone, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for otp id: %d", otpID)
	}

	return nil
}
