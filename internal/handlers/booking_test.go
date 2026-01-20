package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shb/internal/configs"
	"shb/internal/models"
	mockSvc "shb/pkg/mocks/services"
	"shb/pkg/myerrors"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestHandler(service *mockSvc.MockIService) *Handler {
	logger := zerolog.Nop()
	cfg := &configs.Config{}
	return NewHandler(service, nil, nil, &logger, cfg)
}

func TestHandler_createBooking(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         interface{}
		body           interface{}
		setupMocks     func(*mockSvc.MockIService)
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful booking creation",
			userID: 1,
			body: map[string]interface{}{
				"need_id":  1,
				"quantity": 10.5,
				"note":     "I can help",
			},
			setupMocks: func(service *mockSvc.MockIService) {
				service.On("CreateBooking", mock.Anything, 1, 1, 10.5, "I can help").Return(1, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "user not authenticated",
			userID: nil,
			body: map[string]interface{}{
				"need_id":  1,
				"quantity": 10.5,
			},
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "invalid user ID type",
			userID: "not-an-int",
			body: map[string]interface{}{
				"need_id":  1,
				"quantity": 10.5,
			},
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "invalid input - missing need_id",
			userID: 1,
			body: map[string]interface{}{
				"quantity": 10.5,
			},
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "invalid input - invalid quantity",
			userID: 1,
			body: map[string]interface{}{
				"need_id":  1,
				"quantity": -5.0,
			},
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "service error",
			userID: 1,
			body: map[string]interface{}{
				"need_id":  1,
				"quantity": 10.5,
			},
			setupMocks: func(service *mockSvc.MockIService) {
				service.On("CreateBooking", mock.Anything, 1, 1, 10.5, "").Return(0, myerrors.NewBadRequestErr("need not found"))
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockSvc.MockIService)
			tt.setupMocks(service)

			handler := setupTestHandler(service)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			if tt.userID != nil {
				c.Set("userID", tt.userID)
			}

			handler.createBooking(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			service.AssertExpectations(t)
		})
	}
}

func TestHandler_approveBooking(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         interface{}
		bookingID      string
		setupMocks     func(*mockSvc.MockIService)
		expectedStatus int
	}{
		{
			name:      "successful approval",
			userID:    5,
			bookingID: "1",
			setupMocks: func(service *mockSvc.MockIService) {
				service.On("ApproveBooking", mock.Anything, 1, 5).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user not authenticated",
			userID:         nil,
			bookingID:      "1",
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid booking ID",
			userID:         5,
			bookingID:      "invalid",
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "service error",
			userID:    5,
			bookingID: "1",
			setupMocks: func(service *mockSvc.MockIService) {
				service.On("ApproveBooking", mock.Anything, 1, 5).Return(myerrors.NewForbiddenErr("access denied"))
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockSvc.MockIService)
			tt.setupMocks(service)

			handler := setupTestHandler(service)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/"+tt.bookingID+"/approve", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tt.bookingID}}

			if tt.userID != nil {
				c.Set("userID", tt.userID)
			}

			handler.approveBooking(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			service.AssertExpectations(t)
		})
	}
}

func TestHandler_rejectBooking(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         interface{}
		bookingID      string
		setupMocks     func(*mockSvc.MockIService)
		expectedStatus int
	}{
		{
			name:      "successful rejection",
			userID:    5,
			bookingID: "1",
			setupMocks: func(service *mockSvc.MockIService) {
				service.On("RejectBooking", mock.Anything, 1, 5).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "user not authenticated",
			userID:         nil,
			bookingID:      "1",
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid booking ID",
			userID:         5,
			bookingID:      "invalid",
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockSvc.MockIService)
			tt.setupMocks(service)

			handler := setupTestHandler(service)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings/"+tt.bookingID+"/reject", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tt.bookingID}}

			if tt.userID != nil {
				c.Set("userID", tt.userID)
			}

			handler.rejectBooking(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			service.AssertExpectations(t)
		})
	}
}

func TestHandler_getInstitutionBookings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		institutionID  string
		setupMocks     func(*mockSvc.MockIService)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:          "successful retrieval",
			institutionID: "1",
			setupMocks: func(service *mockSvc.MockIService) {
				bookings := []*models.Booking{
					{ID: 1, UserID: 1, NeedID: 1, Status: models.BookingStatusPending},
					{ID: 2, UserID: 2, NeedID: 1, Status: models.BookingStatusApproved},
				}
				service.On("GetBookingsByInstitution", mock.Anything, 1).Return(bookings, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "invalid institution ID",
			institutionID:  "invalid",
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:          "empty list",
			institutionID: "1",
			setupMocks: func(service *mockSvc.MockIService) {
				service.On("GetBookingsByInstitution", mock.Anything, 1).Return([]*models.Booking{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockSvc.MockIService)
			tt.setupMocks(service)

			handler := setupTestHandler(service)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/institutions/"+tt.institutionID+"/bookings", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req
			c.Params = gin.Params{{Key: "id", Value: tt.institutionID}}

			handler.getInstitutionBookings(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var response models.Response
				json.Unmarshal(w.Body.Bytes(), &response)
				if bookings, ok := response.Data.([]interface{}); ok {
					assert.Len(t, bookings, tt.expectedCount)
				}
			}
			service.AssertExpectations(t)
		})
	}
}

func TestHandler_getMyBookings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         interface{}
		setupMocks     func(*mockSvc.MockIService)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:   "successful retrieval",
			userID: 1,
			setupMocks: func(service *mockSvc.MockIService) {
				bookings := []*models.Booking{
					{ID: 1, UserID: 1, NeedID: 1, Status: models.BookingStatusPending},
					{ID: 2, UserID: 1, NeedID: 2, Status: models.BookingStatusApproved},
				}
				service.On("GetBookingsByUser", mock.Anything, 1).Return(bookings, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:           "user not authenticated",
			userID:         nil,
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid user ID type",
			userID:         "not-an-int",
			setupMocks:     func(service *mockSvc.MockIService) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "empty list",
			userID: 1,
			setupMocks: func(service *mockSvc.MockIService) {
				service.On("GetBookingsByUser", mock.Anything, 1).Return([]*models.Booking{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := new(mockSvc.MockIService)
			tt.setupMocks(service)

			handler := setupTestHandler(service)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/bookings/my", nil)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			if tt.userID != nil {
				c.Set("userID", tt.userID)
			}

			handler.getMyBookings(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var response models.Response
				json.Unmarshal(w.Body.Bytes(), &response)
				if bookings, ok := response.Data.([]interface{}); ok {
					assert.Len(t, bookings, tt.expectedCount)
				}
			}
			service.AssertExpectations(t)
		})
	}
}
