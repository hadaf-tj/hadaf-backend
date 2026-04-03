package services_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"shb/internal/models"
	"shb/pkg/myerrors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateBooking(t *testing.T) {
	ctx := context.Background()
	phone := "+123"
	name := "Volunteer"
	email := "inst@example.com"

	need := &models.Need{
		ID:            10,
		InstitutionID: 5,
		Name:          "Food",
		Unit:          "kg",
	}
	activeUser := &models.User{
		ID:       1,
		IsActive: true,
		Phone:    &phone,
		FullName: &name,
		Role:     models.RoleVolunteer,
	}

	tests := []struct {
		name      string
		setup     func(d testDeps)
		userID    int
		needID    int
		qty       float64
		note      string
		wantID    int
		wantErr   bool
		assertErr func(t *testing.T, err error)
	}{
		{
			name: "need not found",
			setup: func(d testDeps) {
				d.Repo.On("GetNeedByID", ctx, 10).Return(nil, fmt.Errorf("get need: %w", myerrors.ErrNotFound))
			},
			userID:  1,
			needID:  10,
			qty:     1,
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var br myerrors.BadRequestErr
				require.ErrorAs(t, err, &br)
			},
		},
		{
			name: "get need database error",
			setup: func(d testDeps) {
				d.Repo.On("GetNeedByID", ctx, 10).Return(nil, errors.New("db down"))
			},
			userID:  1,
			needID:  10,
			qty:     1,
			wantErr: true,
		},
		{
			name: "user not found",
			setup: func(d testDeps) {
				d.Repo.On("GetNeedByID", ctx, 10).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 1).Return(nil, fmt.Errorf("get user: %w", myerrors.ErrNotFound))
			},
			userID:  1,
			needID:  10,
			qty:     1,
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var br myerrors.BadRequestErr
				require.ErrorAs(t, err, &br)
			},
		},
		{
			name: "user inactive",
			setup: func(d testDeps) {
				u := *activeUser
				u.IsActive = false
				d.Repo.On("GetNeedByID", ctx, 10).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 1).Return(&u, nil)
			},
			userID:  1,
			needID:  10,
			qty:     1,
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var br myerrors.BadRequestErr
				require.ErrorAs(t, err, &br)
			},
		},
		{
			name: "quantity zero",
			setup: func(d testDeps) {
				d.Repo.On("GetNeedByID", ctx, 10).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 1).Return(activeUser, nil)
			},
			userID:  1,
			needID:  10,
			qty:     0,
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var br myerrors.BadRequestErr
				require.ErrorAs(t, err, &br)
			},
		},
		{
			name: "duplicate active booking",
			setup: func(d testDeps) {
				d.Repo.On("GetNeedByID", ctx, 10).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 1).Return(activeUser, nil)
				d.Repo.On("GetActiveBookingByUserAndNeed", ctx, 1, 10).Return(&models.Booking{ID: 99}, nil)
			},
			userID:  1,
			needID:  10,
			qty:     1,
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var ce myerrors.ConflictErr
				require.ErrorAs(t, err, &ce)
			},
		},
		{
			name: "get active booking error",
			setup: func(d testDeps) {
				d.Repo.On("GetNeedByID", ctx, 10).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 1).Return(activeUser, nil)
				d.Repo.On("GetActiveBookingByUserAndNeed", ctx, 1, 10).Return(nil, errors.New("scan"))
			},
			userID:  1,
			needID:  10,
			qty:     1,
			wantErr: true,
		},
		{
			name: "create booking error",
			setup: func(d testDeps) {
				d.Repo.On("GetNeedByID", ctx, 10).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 1).Return(activeUser, nil)
				d.Repo.On("GetActiveBookingByUserAndNeed", ctx, 1, 10).Return(nil, nil)
				d.Repo.On("CreateBooking", ctx, bookingMatcher(1, 10)).Return(0, errors.New("insert fail"))
			},
			userID:  1,
			needID:  10,
			qty:     2,
			wantErr: true,
		},
		{
			name: "success institution fetch fails email skipped",
			setup: func(d testDeps) {
				d.Repo.On("GetNeedByID", ctx, 10).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 1).Return(activeUser, nil)
				d.Repo.On("GetActiveBookingByUserAndNeed", ctx, 1, 10).Return(nil, nil)
				d.Repo.On("CreateBooking", ctx, bookingMatcher(1, 10)).Return(42, nil)
				d.Repo.On("GetInstitutionByID", ctx, 5).Return(nil, errors.New("no institution"))
			},
			userID:  1,
			needID:  10,
			qty:     2,
			wantID:  42,
			wantErr: false,
		},
		{
			name: "success with email notification",
			setup: func(d testDeps) {
				inst := &models.Institution{ID: 5, Name: "Home", Email: &email}
				d.Repo.On("GetNeedByID", ctx, 10).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 1).Return(activeUser, nil)
				d.Repo.On("GetActiveBookingByUserAndNeed", ctx, 1, 10).Return(nil, nil)
				d.Repo.On("CreateBooking", ctx, bookingMatcher(1, 10)).Return(7, nil)
				d.Repo.On("GetInstitutionByID", ctx, 5).Return(inst, nil)
				d.Email.On("SendEmail", ctx, email, mock.MatchedBy(func(s string) bool { return strings.Contains(s, "Новый") }),
					mock.MatchedBy(func(s string) bool { return strings.Contains(s, "Food") })).Return(nil)
			},
			userID:  1,
			needID:  10,
			qty:     2,
			wantID:  7,
			wantErr: false,
		},
		{
			name: "success email send failure still returns id",
			setup: func(d testDeps) {
				inst := &models.Institution{ID: 5, Name: "Home", Email: &email}
				d.Repo.On("GetNeedByID", ctx, 10).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 1).Return(activeUser, nil)
				d.Repo.On("GetActiveBookingByUserAndNeed", ctx, 1, 10).Return(nil, nil)
				d.Repo.On("CreateBooking", ctx, bookingMatcher(1, 10)).Return(8, nil)
				d.Repo.On("GetInstitutionByID", ctx, 5).Return(inst, nil)
				d.Email.On("SendEmail", ctx, email, mock.MatchedBy(func(s string) bool { return strings.Contains(s, "Новый") }),
					mock.MatchedBy(func(s string) bool { return strings.Contains(s, "Food") })).Return(errors.New("smtp down"))
			},
			userID:  1,
			needID:  10,
			qty:     2,
			wantID:  8,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, d := newTestService(t)
			tt.setup(d)

			id, err := svc.CreateBooking(ctx, tt.userID, tt.needID, tt.qty, tt.note)
			if tt.wantErr {
				require.Error(t, err)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantID, id)
		})
	}
}

func bookingMatcher(userID, needID int) interface{} {
	return mock.MatchedBy(func(b *models.Booking) bool {
		return b.UserID == userID && b.NeedID == needID && b.Status == models.BookingStatusPending
	})
}

// --- ApproveBooking ---

func TestService_ApproveBooking(t *testing.T) {
	ctx := context.Background()
	instID := 5
	booking := &models.Booking{ID: 1, NeedID: 20, UserID: 3}
	need := &models.Need{ID: 20, InstitutionID: 5,
		Name: "x", Unit: "u"}

	tests := []struct {
		name      string
		setup     func(d testDeps)
		userID    int
		bookingID int
		wantErr   bool
		assertErr func(t *testing.T, err error)
	}{
		{
			name: "get booking error",
			setup: func(d testDeps) {
				d.Repo.On("GetBookingByID", ctx, 1).Return(nil, errors.New("db"))
			},
			bookingID: 1,
			userID:    10,
			wantErr:   true,
		},
		{
			name: "forbidden volunteer",
			setup: func(d testDeps) {
				d.Repo.On("GetBookingByID", ctx, 1).Return(booking, nil)
				d.Repo.On("GetNeedByID", ctx, 20).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 10).Return(&models.User{ID: 10, Role: models.RoleVolunteer}, nil)
			},
			bookingID: 1,
			userID:    10,
			wantErr:   true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var fe myerrors.ForbiddenErr
				require.ErrorAs(t, err, &fe)
			},
		},
		{
			name: "employee wrong institution",
			setup: func(d testDeps) {
				d.Repo.On("GetBookingByID", ctx, 1).Return(booking, nil)
				d.Repo.On("GetNeedByID", ctx, 20).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 10).Return(&models.User{ID: 10, Role: models.RoleEmployee, InstitutionID: ptrInt(99)}, nil)
			},
			bookingID: 1,
			userID:    10,
			wantErr:   true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var fe myerrors.ForbiddenErr
				require.ErrorAs(t, err, &fe)
			},
		},
		{
			name: "success institution employee",
			setup: func(d testDeps) {
				d.Repo.On("GetBookingByID", ctx, 1).Return(booking, nil)
				d.Repo.On("GetNeedByID", ctx, 20).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 10).Return(&models.User{ID: 10, Role: models.RoleEmployee, InstitutionID: &instID}, nil)
				d.Repo.On("UpdateBookingStatus", ctx, 1, models.BookingStatusApproved).Return(nil)
			},
			bookingID: 1,
			userID:    10,
			wantErr:   false,
		},
		{
			name: "success super admin",
			setup: func(d testDeps) {
				d.Repo.On("GetBookingByID", ctx, 1).Return(booking, nil)
				d.Repo.On("GetNeedByID", ctx, 20).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 10).Return(&models.User{ID: 10, Role: models.RoleSuperAdmin}, nil)
				d.Repo.On("UpdateBookingStatus", ctx, 1, models.BookingStatusApproved).Return(nil)
			},
			bookingID: 1,
			userID:    10,
			wantErr:   false,
		},
		{
			name: "update status error",
			setup: func(d testDeps) {
				d.Repo.On("GetBookingByID", ctx, 1).Return(booking, nil)
				d.Repo.On("GetNeedByID", ctx, 20).Return(need, nil)
				d.Repo.On("GetUserByID", ctx, 10).Return(&models.User{ID: 10, Role: models.RoleSuperAdmin}, nil)
				d.Repo.On("UpdateBookingStatus", ctx, 1, models.BookingStatusApproved).Return(errors.New("db"))
			},
			bookingID: 1,
			userID:    10,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, d := newTestService(t)
			tt.setup(d)

			err := svc.ApproveBooking(ctx, tt.bookingID, tt.userID)
			if tt.wantErr {
				require.Error(t, err)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
				return
			}
			require.NoError(t, err)
		})
	}
}

func ptrInt(v int) *int { return &v }
