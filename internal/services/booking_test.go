package services

import (
	"context"
	"errors"
	"shb/internal/configs"
	"shb/internal/models"
	emailmock "shb/pkg/mocks/email"
	repositorymock "shb/pkg/mocks/repository"
	"shb/pkg/myerrors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateBooking(t *testing.T) {
	tests := []struct {
		name           string
		userID         int
		needID         int
		quantity       float64
		note           string
		setupMocks     func(*repositorymock.MockIRepository, *emailmock.MockIEmailAdapter)
		expectedID     int
		expectedError  error
		expectedErrMsg string
	}{
		{
			name:     "successful booking creation with email",
			userID:   1,
			needID:   1,
			quantity: 10.5,
			note:     "I can help",
			setupMocks: func(repo *repositorymock.MockIRepository, emailMock *emailmock.MockIEmailAdapter) {
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
					Name:          "Подгузники",
					Unit:          "шт",
					RequiredQty:   100,
				}
				user := &models.User{
					ID:       1,
					FullName: stringPtr("John Doe"),
					Phone:    stringPtr("+992123456789"),
					IsActive: true,
					Role:     models.RoleDonor,
				}
				institution := &models.Institution{
					ID:    1,
					Name:  "Детский дом №1",
					Email: stringPtr("director@example.com"),
				}
				booking := &models.Booking{
					UserID:   1,
					NeedID:   1,
					Quantity: 10.5,
					Note:     "I can help",
					Status:   models.BookingStatusPending,
				}

				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 1).Return(user, nil)
				repo.EXPECT().CreateBooking(mock.Anything, booking).Return(1, nil)
				repo.EXPECT().GetInstitutionByID(mock.Anything, 1).Return(institution, nil)
				emailMock.EXPECT().SendEmail(mock.Anything, "director@example.com", mock.Anything, mock.Anything).Return(nil)
			},
			expectedID: 1,
		},
		{
			name:     "successful booking creation without email",
			userID:   1,
			needID:   1,
			quantity: 5.0,
			note:     "",
			setupMocks: func(repo *repositorymock.MockIRepository, emailMock *emailmock.MockIEmailAdapter) {
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
					Name:          "Подгузники",
					Unit:          "шт",
				}
				user := &models.User{
					ID:       1,
					IsActive: true,
					Role:     models.RoleDonor,
				}
				institution := &models.Institution{
					ID:    1,
					Name:  "Детский дом №1",
					Email: nil,
				}
				booking := &models.Booking{
					UserID:   1,
					NeedID:   1,
					Quantity: 5.0,
					Note:     "",
					Status:   models.BookingStatusPending,
				}

				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 1).Return(user, nil)
				repo.EXPECT().CreateBooking(mock.Anything, booking).Return(2, nil)
				repo.EXPECT().GetInstitutionByID(mock.Anything, 1).Return(institution, nil)
			},
			expectedID: 2,
		},
		{
			name:     "need not found",
			userID:   1,
			needID:   999,
			quantity: 10.0,
			note:     "",
			setupMocks: func(repo *repositorymock.MockIRepository, emailMock *emailmock.MockIEmailAdapter) {
				repo.EXPECT().GetNeedByID(mock.Anything, 999).Return(nil, myerrors.ErrNotFound)
			},
			expectedError:  myerrors.BadRequestErr{},
			expectedErrMsg: "need not found",
		},
		{
			name:     "user not found",
			userID:   999,
			needID:   1,
			quantity: 10.0,
			note:     "",
			setupMocks: func(repo *repositorymock.MockIRepository, emailMock *emailmock.MockIEmailAdapter) {
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
					Name:          "Подгузники",
				}
				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 999).Return(nil, myerrors.ErrNotFound)
			},
			expectedError:  myerrors.BadRequestErr{},
			expectedErrMsg: "user not found",
		},
		{
			name:     "user not active",
			userID:   1,
			needID:   1,
			quantity: 10.0,
			note:     "",
			setupMocks: func(repo *repositorymock.MockIRepository, emailMock *emailmock.MockIEmailAdapter) {
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
					Name:          "Подгузники",
				}
				user := &models.User{
					ID:       1,
					IsActive: false,
				}
				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 1).Return(user, nil)
			},
			expectedError:  myerrors.BadRequestErr{},
			expectedErrMsg: "user is not active",
		},
		{
			name:     "invalid quantity",
			userID:   1,
			needID:   1,
			quantity: 0,
			note:     "",
			setupMocks: func(repo *repositorymock.MockIRepository, emailMock *emailmock.MockIEmailAdapter) {
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
					Name:          "Подгузники",
				}
				user := &models.User{
					ID:       1,
					IsActive: true,
				}
				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 1).Return(user, nil)
			},
			expectedError:  myerrors.BadRequestErr{},
			expectedErrMsg: "quantity must be greater than 0",
		},
		{
			name:     "negative quantity",
			userID:   1,
			needID:   1,
			quantity: -5.0,
			note:     "",
			setupMocks: func(repo *repositorymock.MockIRepository, emailMock *emailmock.MockIEmailAdapter) {
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
					Name:          "Подгузники",
				}
				user := &models.User{
					ID:       1,
					IsActive: true,
				}
				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 1).Return(user, nil)
			},
			expectedError:  myerrors.BadRequestErr{},
			expectedErrMsg: "quantity must be greater than 0",
		},
		{
			name:     "email send failure does not fail booking",
			userID:   1,
			needID:   1,
			quantity: 10.0,
			note:     "test",
			setupMocks: func(repo *repositorymock.MockIRepository, emailMock *emailmock.MockIEmailAdapter) {
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
					Name:          "Подгузники",
					Unit:          "шт",
				}
				user := &models.User{
					ID:       1,
					IsActive: true,
					Role:     models.RoleDonor,
				}
				institution := &models.Institution{
					ID:    1,
					Name:  "Детский дом №1",
					Email: stringPtr("director@example.com"),
				}
				booking := &models.Booking{
					UserID:   1,
					NeedID:   1,
					Quantity: 10.0,
					Note:     "test",
					Status:   models.BookingStatusPending,
				}

				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 1).Return(user, nil)
				repo.EXPECT().CreateBooking(mock.Anything, booking).Return(3, nil)
				repo.EXPECT().GetInstitutionByID(mock.Anything, 1).Return(institution, nil)
				emailMock.EXPECT().SendEmail(mock.Anything, "director@example.com", mock.Anything, mock.Anything).Return(errors.New("smtp error"))
			},
			expectedID: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repositorymock.NewMockIRepository(t)
			emailMock := emailmock.NewMockIEmailAdapter(t)
			logger := zerolog.Nop()

			service := &Service{
				cfg:    &configs.ServiceConfig{},
				logger: &logger,
				repo:   repo,
				email:  emailMock,
			}

			tt.setupMocks(repo, emailMock)

			ctx := context.Background()
			id, err := service.CreateBooking(ctx, tt.userID, tt.needID, tt.quantity, tt.note)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.IsType(t, tt.expectedError, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
				assert.Equal(t, 0, id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, id)
			}
		})
	}
}

func TestService_ApproveBooking(t *testing.T) {
	tests := []struct {
		name           string
		bookingID      int
		userID         int
		setupMocks     func(*repositorymock.MockIRepository)
		expectedError  error
		expectedErrMsg string
	}{
		{
			name:      "successful approval by super admin",
			bookingID: 1,
			userID:    10,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				booking := &models.Booking{
					ID:     1,
					NeedID: 1,
					Status: models.BookingStatusPending,
				}
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
				}
				user := &models.User{
					ID:   10,
					Role: models.RoleSuperAdmin,
				}

				repo.EXPECT().GetBookingByID(mock.Anything, 1).Return(booking, nil)
				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 10).Return(user, nil)
				repo.EXPECT().UpdateBookingStatus(mock.Anything, 1, models.BookingStatusApproved).Return(nil)
			},
		},
		{
			name:      "successful approval by employee of same institution",
			bookingID: 1,
			userID:    5,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				booking := &models.Booking{
					ID:     1,
					NeedID: 1,
					Status: models.BookingStatusPending,
				}
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
				}
				institutionID := 1
				user := &models.User{
					ID:            5,
					Role:          models.RoleEmployee,
					InstitutionID: &institutionID,
				}

				repo.EXPECT().GetBookingByID(mock.Anything, 1).Return(booking, nil)
				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 5).Return(user, nil)
				repo.EXPECT().UpdateBookingStatus(mock.Anything, 1, models.BookingStatusApproved).Return(nil)
			},
		},
		{
			name:      "booking not found",
			bookingID: 999,
			userID:    10,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				repo.EXPECT().GetBookingByID(mock.Anything, 999).Return(nil, errors.New("booking not found"))
			},
			expectedError: errors.New(""),
		},
		{
			name:      "forbidden - donor role",
			bookingID: 1,
			userID:    5,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				booking := &models.Booking{
					ID:     1,
					NeedID: 1,
				}
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
				}
				user := &models.User{
					ID:   5,
					Role: models.RoleDonor,
				}

				repo.EXPECT().GetBookingByID(mock.Anything, 1).Return(booking, nil)
				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 5).Return(user, nil)
			},
			expectedError:  myerrors.ForbiddenErr{},
			expectedErrMsg: "only employees and super admins can approve bookings",
		},
		{
			name:      "forbidden - employee of different institution",
			bookingID: 1,
			userID:    5,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				booking := &models.Booking{
					ID:     1,
					NeedID: 1,
				}
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
				}
				institutionID := 2
				user := &models.User{
					ID:            5,
					Role:          models.RoleEmployee,
					InstitutionID: &institutionID,
				}

				repo.EXPECT().GetBookingByID(mock.Anything, 1).Return(booking, nil)
				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 5).Return(user, nil)
			},
			expectedError:  myerrors.ForbiddenErr{},
			expectedErrMsg: "you can only approve bookings for your own institution",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repositorymock.NewMockIRepository(t)
			logger := zerolog.Nop()

			service := &Service{
				cfg:    &configs.ServiceConfig{},
				logger: &logger,
				repo:   repo,
			}

			tt.setupMocks(repo)

			ctx := context.Background()
			err := service.ApproveBooking(ctx, tt.bookingID, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_RejectBooking(t *testing.T) {
	tests := []struct {
		name           string
		bookingID      int
		userID         int
		setupMocks     func(*repositorymock.MockIRepository)
		expectedError  error
		expectedErrMsg string
	}{
		{
			name:      "successful rejection by super admin",
			bookingID: 1,
			userID:    10,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				booking := &models.Booking{
					ID:     1,
					NeedID: 1,
					Status: models.BookingStatusPending,
				}
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
				}
				user := &models.User{
					ID:   10,
					Role: models.RoleSuperAdmin,
				}

				repo.EXPECT().GetBookingByID(mock.Anything, 1).Return(booking, nil)
				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 10).Return(user, nil)
				repo.EXPECT().UpdateBookingStatus(mock.Anything, 1, models.BookingStatusRejected).Return(nil)
			},
		},
		{
			name:      "forbidden - donor role",
			bookingID: 1,
			userID:    5,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				booking := &models.Booking{
					ID:     1,
					NeedID: 1,
				}
				need := &models.Need{
					ID:            1,
					InstitutionID: 1,
				}
				user := &models.User{
					ID:   5,
					Role: models.RoleDonor,
				}

				repo.EXPECT().GetBookingByID(mock.Anything, 1).Return(booking, nil)
				repo.EXPECT().GetNeedByID(mock.Anything, 1).Return(need, nil)
				repo.EXPECT().GetUserByID(mock.Anything, 5).Return(user, nil)
			},
			expectedError:  myerrors.ForbiddenErr{},
			expectedErrMsg: "only employees and super admins can reject bookings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repositorymock.NewMockIRepository(t)
			logger := zerolog.Nop()

			service := &Service{
				cfg:    &configs.ServiceConfig{},
				logger: &logger,
				repo:   repo,
			}

			tt.setupMocks(repo)

			ctx := context.Background()
			err := service.RejectBooking(ctx, tt.bookingID, tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedErrMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestService_GetBookingsByInstitution(t *testing.T) {
	tests := []struct {
		name          string
		institutionID int
		setupMocks    func(*repositorymock.MockIRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name:          "successful retrieval",
			institutionID: 1,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				bookings := []*models.Booking{
					{ID: 1, UserID: 1, NeedID: 1, Status: models.BookingStatusPending},
					{ID: 2, UserID: 2, NeedID: 1, Status: models.BookingStatusApproved},
				}
				repo.EXPECT().GetBookingsByInstitution(mock.Anything, 1).Return(bookings, nil)
			},
			expectedCount: 2,
		},
		{
			name:          "empty list",
			institutionID: 1,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				repo.EXPECT().GetBookingsByInstitution(mock.Anything, 1).Return([]*models.Booking{}, nil)
			},
			expectedCount: 0,
		},
		{
			name:          "repository error",
			institutionID: 1,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				repo.EXPECT().GetBookingsByInstitution(mock.Anything, 1).Return(nil, errors.New("db error"))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repositorymock.NewMockIRepository(t)
			logger := zerolog.Nop()

			service := &Service{
				cfg:    &configs.ServiceConfig{},
				logger: &logger,
				repo:   repo,
			}

			tt.setupMocks(repo)

			ctx := context.Background()
			bookings, err := service.GetBookingsByInstitution(ctx, tt.institutionID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, bookings)
			} else {
				assert.NoError(t, err)
				assert.Len(t, bookings, tt.expectedCount)
			}
		})
	}
}

func TestService_GetBookingsByUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        int
		setupMocks    func(*repositorymock.MockIRepository)
		expectedCount int
		expectedError bool
	}{
		{
			name:   "successful retrieval",
			userID: 1,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				bookings := []*models.Booking{
					{ID: 1, UserID: 1, NeedID: 1, Status: models.BookingStatusPending},
					{ID: 2, UserID: 1, NeedID: 2, Status: models.BookingStatusApproved},
				}
				repo.EXPECT().GetBookingsByUser(mock.Anything, 1).Return(bookings, nil)
			},
			expectedCount: 2,
		},
		{
			name:   "empty list",
			userID: 1,
			setupMocks: func(repo *repositorymock.MockIRepository) {
				repo.EXPECT().GetBookingsByUser(mock.Anything, 1).Return([]*models.Booking{}, nil)
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repositorymock.NewMockIRepository(t)
			logger := zerolog.Nop()

			service := &Service{
				cfg:    &configs.ServiceConfig{},
				logger: &logger,
				repo:   repo,
			}

			tt.setupMocks(repo)

			ctx := context.Background()
			bookings, err := service.GetBookingsByUser(ctx, tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, bookings)
			} else {
				assert.NoError(t, err)
				assert.Len(t, bookings, tt.expectedCount)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
