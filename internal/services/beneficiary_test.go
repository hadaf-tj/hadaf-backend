// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors
package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"shb/internal/models"
	"shb/internal/services"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateApplication_Success(t *testing.T) {
	service, deps := newTestService(t)
	ctx := context.Background()

	birthDate := time.Now().AddDate(-11, -11, 0)
	req := &models.Beneficiary{
		UserID:    1,
		FullName:  "Test Child",
		BirthDate: birthDate,
		Diagnosis: models.DiagnosisCerebralPalsy,
		Status:    models.BeneficiaryStatusPending,
	}

	key := "applications:pending:" + time.Now().Format("2006-01")
	deps.Cache.
		On("Get", mock.Anything, key).
		Return("5", nil)
	deps.Repo.
		On("CreateBeneficiary", mock.Anything, req).
		Return(nil)
	deps.Cache.
		On("Increment", mock.Anything, key).
		Return(nil)

	err := service.CreateApplication(ctx, req)

	require.NoError(t, err)
	deps.Repo.AssertExpectations(t)
	deps.Cache.AssertExpectations(t)
}

func TestCreateApplication_UnsupportedDiagnosis(t *testing.T) {
	service, deps := newTestService(t)
	ctx := context.Background()

	req := &models.Beneficiary{
		UserID:    1,
		FullName:  "Wrong Diagnosis Child",
		BirthDate: time.Now().AddDate(-5, 0, 0),
		Diagnosis: "asthma",
	}

	err := service.CreateApplication(ctx, req)

	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrUnsupportedDiagnosis)
	deps.Repo.AssertNotCalled(t, "CreateBeneficiary", mock.Anything, mock.Anything)
	deps.Cache.AssertNotCalled(t, "Increment", mock.Anything, mock.Anything)
	deps.Cache.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
}

func TestCreateApplication_AgeLimitExceeded(t *testing.T) {
	service, deps := newTestService(t)
	ctx := context.Background()

	// 12 years and 1 day old
	birthDate := time.Now().AddDate(-12, 0, -1)
	req := &models.Beneficiary{
		UserID:    1,
		FullName:  "Too Old Child",
		BirthDate: birthDate,
		Diagnosis: models.DiagnosisCerebralPalsy,
	}

	err := service.CreateApplication(ctx, req)

	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrAgeLimitExceeded)
	deps.Repo.AssertNotCalled(t, "CreateBeneficiary", mock.Anything, mock.Anything)
	deps.Cache.AssertNotCalled(t, "Increment", mock.Anything, mock.Anything)
}

func TestCreateApplication_AgeLimitBoundary(t *testing.T) {
	// Exactly 12 years old today: age >= 12 must be rejected, not just > 12.
	service, deps := newTestService(t)
	ctx := context.Background()

	birthDate := time.Now().AddDate(-12, 0, 0)
	req := &models.Beneficiary{
		UserID:    1,
		FullName:  "Exactly Twelve",
		BirthDate: birthDate,
		Diagnosis: models.DiagnosisCerebralPalsy,
	}

	err := service.CreateApplication(ctx, req)

	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrAgeLimitExceeded)
	deps.Repo.AssertNotCalled(t, "CreateBeneficiary", mock.Anything, mock.Anything)
}

func TestCreateApplication_QuotaExceeded(t *testing.T) {
	service, deps := newTestService(t)
	ctx := context.Background()

	req := &models.Beneficiary{
		UserID:    1,
		FullName:  "Quota Child",
		BirthDate: time.Now().AddDate(-5, 0, 0),
		Diagnosis: models.DiagnosisCerebralPalsy,
	}

	key := "applications:pending:" + time.Now().Format("2006-01")
	// Limit reached
	deps.Cache.
		On("Get", mock.Anything, key).
		Return("10", nil)

	err := service.CreateApplication(ctx, req)

	require.Error(t, err)
	assert.ErrorIs(t, err, services.ErrApplicationQuotaExceed)
	deps.Repo.AssertNotCalled(t, "CreateBeneficiary", mock.Anything, mock.Anything)
	deps.Cache.AssertNotCalled(t, "Increment", mock.Anything, mock.Anything)
}

func TestCreateApplication_QuotaFallsBackToPostgresOnCacheMiss(t *testing.T) {
	service, deps := newTestService(t)
	ctx := context.Background()

	req := &models.Beneficiary{
		UserID:    1,
		FullName:  "Cold Cache Child",
		BirthDate: time.Now().AddDate(-5, 0, 0),
		Diagnosis: models.DiagnosisCerebralPalsy,
		Status:    models.BeneficiaryStatusPending,
	}

	key := "applications:pending:" + time.Now().Format("2006-01")
	deps.Cache.
		On("Get", mock.Anything, key).
		Return("", redis.Nil)
	deps.Repo.
		On("CountPendingThisMonth", mock.Anything).
		Return(3, nil)
	deps.Cache.
		On("Set", mock.Anything, key, 3, 32*24*time.Hour).
		Return(nil)
	deps.Repo.
		On("CreateBeneficiary", mock.Anything, req).
		Return(nil)
	deps.Cache.
		On("Increment", mock.Anything, key).
		Return(nil)

	err := service.CreateApplication(ctx, req)

	require.NoError(t, err)
	deps.Repo.AssertExpectations(t)
	deps.Cache.AssertExpectations(t)
}

func TestCreateApplication_NilRequest(t *testing.T) {
	service, deps := newTestService(t)
	ctx := context.Background()

	err := service.CreateApplication(ctx, nil)

	require.Error(t, err)
	deps.Repo.AssertNotCalled(t, "CreateBeneficiary", mock.Anything, mock.Anything)
	deps.Cache.AssertNotCalled(t, "Get", mock.Anything, mock.Anything)
}

func TestCreateApplication_RepoErrorDoesNotIncrementCache(t *testing.T) {
	service, deps := newTestService(t)
	ctx := context.Background()

	req := &models.Beneficiary{
		UserID:    1,
		FullName:  "Repo Failure Child",
		BirthDate: time.Now().AddDate(-5, 0, 0),
		Diagnosis: models.DiagnosisCerebralPalsy,
		Status:    models.BeneficiaryStatusPending,
	}

	key := "applications:pending:" + time.Now().Format("2006-01")
	deps.Cache.
		On("Get", mock.Anything, key).
		Return("5", nil)
	deps.Repo.
		On("CreateBeneficiary", mock.Anything, req).
		Return(errors.New("db unavailable"))

	err := service.CreateApplication(ctx, req)

	require.Error(t, err)
	deps.Cache.AssertNotCalled(t, "Increment", mock.Anything, mock.Anything)
	deps.Repo.AssertExpectations(t)
}
