// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"shb/internal/configs"
	"shb/internal/models"
	"shb/pkg/metrics"
	"shb/pkg/middlewares"
	servicemock "shb/pkg/mocks/services"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_GetApplicationQuota(t *testing.T) {
	gin.SetMode(gin.TestMode)

	for _, tt := range []struct {
		name    string
		quota   *models.ApplicationQuota
		allowed bool
	}{
		{name: "allowed", quota: &models.ApplicationQuota{Limit: 10, Used: 4, Remaining: 6, Allowed: true}, allowed: true},
		{name: "exhausted", quota: &models.ApplicationQuota{Limit: 10, Used: 10, Remaining: 0, Allowed: false}, allowed: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			svc := servicemock.NewMockIService(t)
			svc.On("GetApplicationQuota", mock.Anything).Return(tt.quota, nil)
			h := &Handler{service: svc}
			r := gin.New()
			r.GET("/api/v1/beneficiaries/quota", h.getApplicationQuota)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/beneficiaries/quota", nil))
			require.Equal(t, http.StatusOK, w.Code)

			var body struct {
				Message string                  `json:"message"`
				Data    models.ApplicationQuota `json:"data"`
			}
			require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
			require.Equal(t, "Success", body.Message)
			require.Equal(t, tt.allowed, body.Data.Allowed)
		})
	}
}

func TestHandler_ApplicationQuotaRouteIsPublic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := servicemock.NewMockIService(t)
	svc.On("GetApplicationQuota", mock.Anything).Return(&models.ApplicationQuota{Limit: 10, Remaining: 10, Allowed: true}, nil)
	log := zerolog.Nop()
	h := NewHandler(
		svc,
		nil,
		middlewares.NewMiddleware("test-secret"),
		metrics.New(),
		&log,
		&configs.Config{App: configs.AppConfig{Env: "production"}, Tracing: configs.TracingConfig{ServiceName: "shb-test"}},
	)

	w := httptest.NewRecorder()
	h.InitRoutes().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/v1/beneficiaries/quota", nil))
	require.Equal(t, http.StatusOK, w.Code)
}
