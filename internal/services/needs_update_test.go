package services_test

import (
	"context"
	"errors"
	"testing"

	"shb/internal/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_UpdateNeed_getOldFails(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	d.Repo.On("GetNeedByID", ctx, 1).Return(nil, errors.New("db"))
	err := svc.UpdateNeed(ctx, &models.Need{ID: 1, Name: "x"})
	require.Error(t, err)
}

func TestService_UpdateNeed_updateFails(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	old := &models.Need{ID: 1, InstitutionID: 1, Name: "old", Unit: "u", RequiredQty: 1, ReceivedQty: 0, Urgency: "low"}
	d.Repo.On("GetNeedByID", ctx, 1).Return(old, nil)
	d.Repo.On("UpdateNeed", ctx, mock.AnythingOfType("*models.Need")).Return(errors.New("db"))
	err := svc.UpdateNeed(ctx, &models.Need{ID: 1, Name: "new", Unit: "u", RequiredQty: 1, ReceivedQty: 0, Urgency: "low"})
	require.Error(t, err)
}
