package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"shb/internal/models"
	"shb/pkg/myerrors"
)

func (r *Repository) SaveOTP(ctx context.Context, o *models.OTP) error {
	const query = `
	INSERT INTO otp (receiver, method, otp_code, sent_at, expires_at, attempt, is_verified)
	VALUES (:receiver, :method, :otp_code, :sent_at, :expires_at, :attempt, :is_verified)
	RETURNING id;
`
	args := map[string]interface{}{
		"receiver":    o.Receiver,
		"method":      o.Method,
		"otp_code":    o.OTPCode,
		"sent_at":     o.SentAt,
		"expires_at":  o.ExpiresAt,
		"attempt":     o.Attempt,
		"is_verified": o.IsVerified,
	}

	_, err := r.postgres.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("failed to save otp: %w", err)
	}

	return nil
}

func (r *Repository) GetOTP(ctx context.Context, receiver string) (*models.OTP, error) {
	const query = `
	SELECT *
	FROM otp
	WHERE receiver = :receiver
	  AND is_verified = false
	  AND expires_at > NOW()
	ORDER BY id DESC
	LIMIT 1;
`
	args := map[string]interface{}{"receiver": receiver}

	var otpDB dbOtp
	err := r.postgres.QueryRow(ctx, query, args).Scan(&otpDB)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Str("receiver", receiver).Msg("no valid OTP found in DB")
			return nil, myerrors.ErrNotFound
		}
		return nil, fmt.Errorf("otp not found or expired: %w", err)
	}

	return &models.OTP{
		ID:         otpDB.ID,
		Attempt:    otpDB.Attempt,
		Receiver:   otpDB.Receiver,
		Method:     otpDB.Method,
		OTPCode:    otpDB.OTPCode,
		IsVerified: otpDB.IsVerified,
		SentAt:     otpDB.SentAt,
		ExpiresAt:  otpDB.ExpiresAt,
	}, nil
}

func (r *Repository) MarkOTPAsVerified(ctx context.Context, otpID int) error {
	const query = `UPDATE otp SET is_verified=false, updated_at=NOW() WHERE id=:id`
	args := map[string]interface{}{
		"id": otpID,
	}

	result, err := r.postgres.Exec(ctx, query, args)
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
	query := `
	UPDATE otp
	SET attempt = attempt + 1,
	    updated_at = NOW()
	WHERE id = :id AND phone_number = :phone_number;
`
	args := map[string]interface{}{
		"id":           otpID,
		"phone_number": phone,
	}

	result, err := r.postgres.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("phone:%s failed to increase otp attempt: %w", phone, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated for otp id: %d", otpID)
	}

	return nil
}
