package jwtToken

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"shb/internal/models"
	"shb/pkg/configs"
	"shb/pkg/constants"
	"time"
)

type JwtTokenIssuer struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJwtTokenIssuer() *JwtTokenIssuer {
	accessExpire, _ := time.ParseDuration(configs.AccessExpire)
	refreshExpire, _ := time.ParseDuration(configs.RefreshExpire)
	return &JwtTokenIssuer{
		accessSecret:  configs.AccessSecret,
		refreshSecret: configs.RefreshSecret,
		accessTTL:     accessExpire * time.Hour,
		refreshTTL:    refreshExpire * time.Hour,
	}
}

func (j *JwtTokenIssuer) IssueTokens(ctx context.Context, user *models.User) (string, string, error) {
	now := time.Now().UTC()

	accessClaims := models.CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   constants.AccessSubject,
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("uid:%d", user.ID),
		},
	}

	refreshClaims := models.CustomClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   constants.RefreshSubject,
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("uid:%d", user.ID),
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(j.accessSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(j.refreshSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
