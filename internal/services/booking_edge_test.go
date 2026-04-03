package services_test

import (
	"context"
	"testing"

	"shb/internal/models"
	"shb/pkg/myerrors"

	"github.com/stretchr/testify/require"
)

func TestService_CancelMyBooking_notPending(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	d.Repo.On("GetBookingByID", ctx, 1).Return(&models.Booking{ID: 1, UserID: 2, Status: models.BookingStatusApproved}, nil)
	err := svc.CancelMyBooking(ctx, 1, 2)
	require.Error(t, err)
	var br myerrors.BadRequestErr
	require.ErrorAs(t, err, &br)
}

func TestService_UpdateMyBooking_forbidden_and_badQty(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	d.Repo.On("GetBookingByID", ctx, 1).Return(&models.Booking{ID: 1, UserID: 9, Status: models.BookingStatusPending}, nil)
	err := svc.UpdateMyBooking(ctx, 1, 2, 5)
	require.Error(t, err)
	var fe myerrors.ForbiddenErr
	require.ErrorAs(t, err, &fe)

	svc2, d2 := newTestService(t)
	d2.Repo.On("GetBookingByID", ctx, 1).Return(&models.Booking{ID: 1, UserID: 2, Status: models.BookingStatusPending}, nil)
	err = svc2.UpdateMyBooking(ctx, 1, 2, 0)
	require.Error(t, err)
	var br myerrors.BadRequestErr
	require.ErrorAs(t, err, &br)
}

func TestService_ApproveBooking_getNeedError(t *testing.T) {
	ctx := context.Background()
	svc, d := newTestService(t)
	d.Repo.On("GetBookingByID", ctx, 1).Return(&models.Booking{ID: 1, NeedID: 20}, nil)
	d.Repo.On("GetNeedByID", ctx, 20).Return(nil, myerrors.ErrNotFound)
	err := svc.ApproveBooking(ctx, 1, 10)
	require.Error(t, err)
}
