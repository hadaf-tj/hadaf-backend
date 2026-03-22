package services_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"shb/internal/models"
	"shb/pkg/myerrors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_Register(t *testing.T) {
	ctx := context.Background()
	email := "newuser@example.com"
	phone := "+999"
	instID := 3

	tests := []struct {
		name          string
		role          string
		institutionID *int
		setup         func(d testDeps)
		wantErr       bool
		assertErr     func(t *testing.T, err error)
	}{
		{
			name: "duplicate email",
			setup: func(d testDeps) {
				d.Repo.On("GetUserByEmail", ctx, email).Return(&models.User{ID: 1}, nil)
			},
			role:    models.RoleVolunteer,
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var br myerrors.BadRequestErr
				require.ErrorAs(t, err, &br)
			},
		},
		{
			name: "employee without institution",
			setup: func(d testDeps) {
				d.Repo.On("GetUserByEmail", ctx, email).Return(nil, fmt.Errorf("%w", myerrors.ErrNotFound))
			},
			role:          models.RoleEmployee,
			institutionID: nil,
			wantErr:       true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var br myerrors.BadRequestErr
				require.ErrorAs(t, err, &br)
			},
		},
		{
			name: "institution not found",
			setup: func(d testDeps) {
				d.Repo.On("GetUserByEmail", ctx, email).Return(nil, fmt.Errorf("%w", myerrors.ErrNotFound))
				d.Repo.On("GetInstitutionByID", ctx, instID).Return(nil, errors.New("no rows"))
			},
			role:          models.RoleEmployee,
			institutionID: &instID,
			wantErr:       true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var br myerrors.BadRequestErr
				require.ErrorAs(t, err, &br)
			},
		},
		{
			name: "institution deleted",
			setup: func(d testDeps) {
				d.Repo.On("GetUserByEmail", ctx, email).Return(nil, fmt.Errorf("%w", myerrors.ErrNotFound))
				d.Repo.On("GetInstitutionByID", ctx, instID).Return(&models.Institution{IsDeleted: true}, nil)
			},
			role:          models.RoleEmployee,
			institutionID: &instID,
			wantErr:       true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var br myerrors.BadRequestErr
				require.ErrorAs(t, err, &br)
			},
		},
		{
			name: "success volunteer pending in cache and otp sent",
			setup: func(d testDeps) {
				d.Repo.On("GetUserByEmail", ctx, email).Return(nil, fmt.Errorf("%w", myerrors.ErrNotFound))
				d.Cache.On("Set", ctx, "pending_reg:"+email, mock.MatchedBy(func(s string) bool {
					var u models.User
					_ = json.Unmarshal([]byte(s), &u)
					return u.Role == models.RoleVolunteer
				}), 15*time.Minute).Return(nil)
				d.Repo.On("SaveOTP", ctx, mock.AnythingOfType("*models.OTP")).Return(1, nil)
				d.Email.On("SendEmail", mock.Anything, email, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Maybe()
			},
			role:    models.RoleVolunteer,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, d := newTestService(t)
			tt.setup(d)

			tok, err := svc.Register(ctx, email, phone, "password123", "Full Name", tt.role, tt.institutionID)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, tok)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
				return
			}
			require.NoError(t, err)
			require.Nil(t, tok)
		})
	}
}

func TestService_ConfirmOTP(t *testing.T) {
	ctx := context.Background()
	receiver := "user@example.com"
	otpCode := "123456"

	otpRow := &models.OTP{
		ID:       1,
		Receiver: receiver,
		OTPCode:  otpCode,
	}

	activeUser := &models.User{
		ID:         7,
		Email:      strPtr("user@example.com"),
		Role:       models.RoleVolunteer,
		IsActive:   true,
		IsApproved: true,
	}

	tests := []struct {
		name      string
		setup     func(d testDeps)
		wantErr   bool
		assertErr func(t *testing.T, err error)
	}{
		{
			name: "otp not found",
			setup: func(d testDeps) {
				d.Repo.On("GetOTP", ctx, receiver).Return(nil, errors.New("no otp"))
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var u myerrors.UnauthorizedErr
				require.ErrorAs(t, err, &u)
			},
		},
		{
			name: "wrong otp",
			setup: func(d testDeps) {
				d.Repo.On("GetOTP", ctx, receiver).Return(otpRow, nil)
				d.Repo.On("IncreaseOTPAttempt", ctx, 1, receiver).Return(nil)
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var u myerrors.UnauthorizedErr
				require.ErrorAs(t, err, &u)
			},
		},
		{
			name: "mark verified fails",
			setup: func(d testDeps) {
				d.Repo.On("GetOTP", ctx, receiver).Return(otpRow, nil)
				d.Repo.On("MarkOTPAsVerified", ctx, 1).Return(errors.New("db"))
			},
			wantErr: true,
		},
		{
			name: "user not in db and no cache",
			setup: func(d testDeps) {
				d.Repo.On("GetOTP", ctx, receiver).Return(otpRow, nil)
				d.Repo.On("MarkOTPAsVerified", ctx, 1).Return(nil)
				d.Repo.On("GetUserByEmail", ctx, receiver).Return(nil, fmt.Errorf("%w", myerrors.ErrNotFound))
				d.Cache.On("Get", ctx, "pending_reg:"+receiver).Return("", nil)
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var br myerrors.BadRequestErr
				require.ErrorAs(t, err, &br)
			},
		},
		{
			name: "employee not approved",
			setup: func(d testDeps) {
				d.Repo.On("GetOTP", ctx, receiver).Return(otpRow, nil)
				d.Repo.On("MarkOTPAsVerified", ctx, 1).Return(nil)
				u := *activeUser
				u.Role = models.RoleEmployee
				u.IsApproved = false
				d.Repo.On("GetUserByEmail", ctx, receiver).Return(&u, nil)
			},
			wantErr: true,
			assertErr: func(t *testing.T, err error) {
				t.Helper()
				var fe myerrors.ForbiddenErr
				require.ErrorAs(t, err, &fe)
			},
		},
		{
			name: "success issues tokens",
			setup: func(d testDeps) {
				d.Repo.On("GetOTP", ctx, receiver).Return(otpRow, nil)
				d.Repo.On("MarkOTPAsVerified", ctx, 1).Return(nil)
				d.Repo.On("GetUserByEmail", ctx, receiver).Return(activeUser, nil)
				d.Token.On("IssueTokens", ctx, 7, models.RoleVolunteer, true).Return("access", "refresh", nil)
				d.Repo.On("SaveRefreshToken", ctx, 7, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc, d := newTestService(t)
			tt.setup(d)

			inOTP := otpCode
			if tt.name == "wrong otp" {
				inOTP = "000000"
			}

			tok, err := svc.ConfirmOTP(ctx, receiver, inOTP)
			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, tok)
				if tt.assertErr != nil {
					tt.assertErr(t, err)
				}
				return
			}
			require.NoError(t, err)
			require.NotNil(t, tok)
			require.Equal(t, "access", tok.AccessToken)
			require.Equal(t, "refresh", tok.RefreshToken)
		})
	}
}

func strPtr(s string) *string { return &s }

func TestService_ConfirmOTP_inactiveUser_activated(t *testing.T) {
	ctx := context.Background()
	receiver := "inactive@example.com"
	otpCode := "111111"
	otpRow := &models.OTP{ID: 3, Receiver: receiver, OTPCode: otpCode}
	email := receiver
	inactive := &models.User{
		ID: 3, Email: &email, Role: models.RoleVolunteer, IsActive: false, IsApproved: true,
	}
	svc, d := newTestService(t)
	d.Repo.On("GetOTP", ctx, receiver).Return(otpRow, nil)
	d.Repo.On("MarkOTPAsVerified", ctx, 3).Return(nil)
	d.Repo.On("GetUserByEmail", ctx, receiver).Return(inactive, nil)
	d.Repo.On("ActivateUser", ctx, 3).Return(nil)
	d.Token.On("IssueTokens", ctx, 3, models.RoleVolunteer, true).Return("a", "b", nil)
	d.Repo.On("SaveRefreshToken", ctx, 3, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
	tok, err := svc.ConfirmOTP(ctx, receiver, otpCode)
	require.NoError(t, err)
	require.NotNil(t, tok)
}

func TestService_ConfirmOTP_getUser_dbError(t *testing.T) {
	ctx := context.Background()
	receiver := "db@fail.com"
	otpRow := &models.OTP{ID: 1, Receiver: receiver, OTPCode: "999999"}
	svc, d := newTestService(t)
	d.Repo.On("GetOTP", ctx, receiver).Return(otpRow, nil)
	d.Repo.On("MarkOTPAsVerified", ctx, 1).Return(nil)
	d.Repo.On("GetUserByEmail", ctx, receiver).Return(nil, errors.New("db"))
	_, err := svc.ConfirmOTP(ctx, receiver, "999999")
	require.Error(t, err)
}

func TestService_ConfirmOTP_pendingUserFromCache(t *testing.T) {
	ctx := context.Background()
	receiver := "pending@example.com"
	otpCode := "654321"
	otpRow := &models.OTP{ID: 2, Receiver: receiver, OTPCode: otpCode}

	email := receiver
	pending := models.User{
		Email: &email, Role: models.RoleVolunteer, IsActive: true, IsApproved: true,
	}
	payload, err := json.Marshal(pending)
	require.NoError(t, err)

	svc, d := newTestService(t)
	d.Repo.On("GetOTP", ctx, receiver).Return(otpRow, nil)
	d.Repo.On("MarkOTPAsVerified", ctx, 2).Return(nil)
	d.Repo.On("GetUserByEmail", ctx, receiver).Return(nil, fmt.Errorf("%w", myerrors.ErrNotFound))
	d.Cache.On("Get", ctx, "pending_reg:"+receiver).Return(string(payload), nil)
	d.Repo.On("CreateUser", ctx, mock.MatchedBy(func(u *models.User) bool {
		return u.Email != nil && *u.Email == receiver
	})).Return(nil)
	d.Cache.On("Delete", ctx, "pending_reg:"+receiver).Return(nil)
	d.Token.On("IssueTokens", ctx, 0, models.RoleVolunteer, true).Return("a", "b", nil)
	d.Repo.On("SaveRefreshToken", ctx, 0, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)

	tok, err := svc.ConfirmOTP(ctx, receiver, otpCode)
	require.NoError(t, err)
	require.NotNil(t, tok)
}
