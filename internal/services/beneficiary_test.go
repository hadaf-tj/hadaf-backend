// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services_test

import (
	"context"
	"testing"
	"time"

	"shb/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	// Current month quota
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

	assert.NoError(t, err)

	deps.Repo.AssertExpectations(t)
	deps.Cache.AssertExpectations(t)
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

	assert.Error(t, err)

	// PostgreSQL should never be called
	deps.Repo.AssertNotCalled(
		t,
		"CreateBeneficiary",
		mock.Anything,
		req,
	)

	// Redis increment should never happen
	deps.Cache.AssertNotCalled(
		t,
		"Increment",
		mock.Anything,
		mock.Anything,
	)
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

	assert.Error(t, err)

	// Database must not be touched
	deps.Repo.AssertNotCalled(
		t,
		"CreateBeneficiary",
		mock.Anything,
		req,
	)

	// No increment
	deps.Cache.AssertNotCalled(
		t,
		"Increment",
		mock.Anything,
		mock.Anything,
	)
}
