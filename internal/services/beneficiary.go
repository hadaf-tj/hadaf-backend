// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors
package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"shb/internal/models"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var (
	ErrUnsupportedDiagnosis   = errors.New("ERR_UNSUPPORTED_DIAGNOSIS")
	ErrAgeLimitExceeded       = errors.New("ERR_AGE_LIMIT_EXCEEDED")
	ErrApplicationQuotaExceed = errors.New("ERR_APPLICATION_QUOTA_EXHAUSTED")
	ErrInvalidBirthDate       = errors.New("ERR_INVALID_BIRTH_DATE")
)

// CreateApplication creates a medical assistance application.
//
// Validation flow:
// 1. Check diagnosis
// 2. Check child age
// 3. Check monthly quota
// 4. Save application
// 5. Increment Redis counter
func (s *Service) CreateApplication(
	ctx context.Context,
	req *models.Beneficiary,
) error {
	log := zerolog.Ctx(ctx).
		With().
		Str("service", "CreateApplication").
		Logger()

	if req == nil {
		return errors.New("beneficiary request is nil")
	}

	// 1. Diagnosis validation
	if req.Diagnosis != models.DiagnosisCerebralPalsy {
		log.Warn().
			Str("diagnosis", req.Diagnosis).
			Msg("unsupported diagnosis")
		return ErrUnsupportedDiagnosis
	}

	// 2. Age validation
	if req.BirthDate.IsZero() || req.BirthDate.After(time.Now()) {
		log.Warn().
			Time("birth_date", req.BirthDate).
			Msg("invalid birth date")
		return ErrInvalidBirthDate
	}
	if calculateAge(req.BirthDate) >= 12 {
		log.Warn().
			Time("birth_date", req.BirthDate).
			Msg("age limit exceeded")
		return ErrAgeLimitExceeded
	}

	// 3. Monthly quota
	key := fmt.Sprintf(
		"applications:pending:%s",
		time.Now().Format("2006-01"),
	)
	count, err := s.getApplicationCount(ctx, key)
	if err != nil {
		return fmt.Errorf("get application count: %w", err)
	}
	if count >= s.cfg.MaxPendingApplications {
		log.Warn().
			Int("current", count).
			Int("limit", s.cfg.MaxPendingApplications).
			Msg("monthly application quota exceeded")
		return ErrApplicationQuotaExceed
	}

	// 4. Save to database
	req.Status = models.BeneficiaryStatusPending
	if err := s.repo.CreateBeneficiary(ctx, req); err != nil {
		return fmt.Errorf("create beneficiary: %w", err)
	}

	// 5. Increment Redis quota
	if err := s.cache.Increment(ctx, key); err != nil {
		log.Error().
			Err(err).
			Str("key", key).
			Msg("failed to increment application quota")
	}

	log.Info().
		Int("beneficiary_id", req.ID).
		Msg("beneficiary application created")
	return nil
}

// calculateAge returns full years.
func calculateAge(birthDate time.Time) int {
	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.Month() < birthDate.Month() ||
		(now.Month() == birthDate.Month() &&
			now.Day() < birthDate.Day()) {
		age--
	}
	return age
}

// getApplicationCount returns monthly application count.
//
// Redis is the primary source.
// PostgreSQL is fallback if Redis has no key.
func (s *Service) getApplicationCount(
	ctx context.Context,
	key string,
) (int, error) {
	log := zerolog.Ctx(ctx).
		With().
		Str("service", "getApplicationCount").
		Logger()

	value, err := s.cache.Get(ctx, key)
	if err == nil && value != "" {
		count, convErr := strconv.Atoi(value)
		if convErr != nil {
			log.Error().
				Err(convErr).
				Str("key", key).
				Str("value", value).
				Msg("corrupt redis application counter")
			return 0, fmt.Errorf(
				"invalid redis application counter: %w",
				convErr,
			)
		}
		return count, nil
	}

	// Redis key does not exist
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("read application counter: %w", err)
	}

	// PostgreSQL fallback
	count, err := s.repo.CountPendingThisMonth(ctx)
	if err != nil {
		return 0, fmt.Errorf("count pending applications: %w", err)
	}

	// Restore Redis counter with TTL
	if err := s.cache.Set(ctx, key, count, 32*24*time.Hour); err != nil {
		log.Error().
			Err(err).
			Str("key", key).
			Msg("failed to seed application counter in cache")
	}

	return count, nil
}
