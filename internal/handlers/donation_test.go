// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"shb/internal/configs"
	"shb/internal/models"
	"shb/pkg/metrics"
	"shb/pkg/middlewares"
	servicemock "shb/pkg/mocks/services"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const donationTestJWTSecret = "donation-test-secret"

func newDonationTestHandler(t *testing.T, svc *servicemock.MockIService) *Handler {
	t.Helper()
	log := zerolog.Nop()
	return NewHandler(
		svc,
		nil,
		middlewares.NewMiddleware(donationTestJWTSecret),
		metrics.New(),
		&log,
		&configs.Config{App: configs.AppConfig{Env: "production"}, Tracing: configs.TracingConfig{ServiceName: "shb-test"}},
	)
}

func donationAuthHeader(t *testing.T, userID int) string {
	t.Helper()
	claims := models.CustomClaims{
		UserID:     userID,
		Role:       models.RoleVolunteer,
		IsApproved: true,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(donationTestJWTSecret))
	require.NoError(t, err)
	return "Bearer " + token
}

func TestHandler_GetMyDonationsRequiresAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := servicemock.NewMockIService(t)
	h := newDonationTestHandler(t, svc)

	w := httptest.NewRecorder()
	h.InitRoutes().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/donations/my", nil))
	require.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestHandler_GetMyDonationsUsesAuthenticatedUserOnly(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := servicemock.NewMockIService(t)
	page := &models.DonationHistoryPage{
		Items:  []*models.DonationHistoryItem{{ID: 11, CampaignID: 101, CampaignTitle: "Donor one campaign", Amount: 30, Currency: "TJS", Status: models.DonationStatusConfirmed}},
		Total:  1,
		Limit:  20,
		Offset: 0,
	}
	svc.On("GetMyDonations", mock.Anything, 1, 20, 0).Return(page, nil).Once()
	h := newDonationTestHandler(t, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/donations/my?user_id=2", nil)
	req.Header.Set("Authorization", donationAuthHeader(t, 1))
	w := httptest.NewRecorder()
	h.InitRoutes().ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Data models.DonationHistoryPage `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Len(t, body.Data.Items, 1)
	require.Equal(t, 11, body.Data.Items[0].ID)
	require.Equal(t, "Donor one campaign", body.Data.Items[0].CampaignTitle)
}

func TestHandler_GetMyDonationsReturnsPaginatedSafeHistory(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := servicemock.NewMockIService(t)
	page := &models.DonationHistoryPage{
		Items: []*models.DonationHistoryItem{
			{ID: 22, CampaignID: 202, CampaignTitle: "New campaign", Amount: 70, Currency: "TJS", Status: models.DonationStatusConfirmed, CreatedAt: time.Date(2026, 9, 18, 12, 0, 0, 0, time.UTC)},
			{ID: 21, CampaignID: 201, CampaignTitle: "Old campaign", Amount: 50, Currency: "TJS", Status: models.DonationStatusConfirmed, CreatedAt: time.Date(2026, 9, 17, 12, 0, 0, 0, time.UTC)},
		},
		Total:  4,
		Limit:  2,
		Offset: 2,
	}
	svc.On("GetMyDonations", mock.Anything, 7, 2, 2).Return(page, nil).Once()
	h := newDonationTestHandler(t, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/donations/my?limit=2&offset=2", nil)
	req.Header.Set("Authorization", donationAuthHeader(t, 7))
	w := httptest.NewRecorder()
	h.InitRoutes().ServeHTTP(w, req)
	require.Equal(t, http.StatusOK, w.Code)

	var body struct {
		Data models.DonationHistoryPage `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Equal(t, 4, body.Data.Total)
	require.Equal(t, 2, body.Data.Limit)
	require.Equal(t, 2, body.Data.Offset)
	require.Equal(t, []int{22, 21}, []int{body.Data.Items[0].ID, body.Data.Items[1].ID})
	require.Equal(t, "New campaign", body.Data.Items[0].CampaignTitle)
	require.Equal(t, models.DonationStatusConfirmed, body.Data.Items[0].Status)
	require.NotContains(t, w.Body.String(), "payment_ref")
	require.NotContains(t, w.Body.String(), "donor_user_id")
	require.NotContains(t, w.Body.String(), "donor_message")
	require.NotContains(t, w.Body.String(), "description")
	require.NotContains(t, w.Body.String(), "is_anonymous")
}

func TestHandler_GetMyDonationsRejectsInvalidPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := servicemock.NewMockIService(t)
	h := newDonationTestHandler(t, svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/donations/my?offset=-1", nil)
	req.Header.Set("Authorization", donationAuthHeader(t, 7))
	w := httptest.NewRecorder()
	h.InitRoutes().ServeHTTP(w, req)
	require.Equal(t, http.StatusBadRequest, w.Code)
}
