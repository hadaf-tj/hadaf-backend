// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services_test

import (
	"context"
	"errors"
	"testing"

	"shb/internal/models"

	"github.com/stretchr/testify/require"
)

func TestService_DeleteNeed_repoError(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	d.Repo.On("DeleteNeed", ctx, 1).Return(errors.New("db"))
	err := svc.DeleteNeed(ctx, 1)
	require.Error(t, err)
}

func TestService_CreateNeed_error(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	n := &models.Need{InstitutionID: 1, Name: "n", Unit: "u", RequiredQty: 1, Urgency: "low"}
	d.Repo.On("CreateNeed", ctx, n).Return(0, errors.New("db"))
	_, err := svc.CreateNeed(ctx, n)
	require.Error(t, err)
}
