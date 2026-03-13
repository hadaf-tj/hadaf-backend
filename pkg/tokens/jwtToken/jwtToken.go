package jwtToken

import (
	"context"
	"fmt"
	"shb/internal/models"
	"shb/pkg/constants"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtTokenIssuer struct {
	secretKey  string // Используем один ключ для простоты и совместимости с middleware
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// Теперь принимаем конфиг аргументами
func NewJwtTokenIssuer(secretKey string, accessTTL, refreshTTL time.Duration) *JwtTokenIssuer {
	return &JwtTokenIssuer{
		secretKey:  secretKey,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (j *JwtTokenIssuer) IssueTokens(ctx context.Context, id int, role string) (string, string, error) {
	now := time.Now().UTC()

	accessClaims := models.CustomClaims{
		UserID: id,
		Role:   role, // <--- Добавили роль
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   constants.AccessSubject,
			ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("uid:%d", id),
		},
	}

	refreshClaims := models.CustomClaims{
		UserID: id,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   constants.RefreshSubject,
			ExpiresAt: jwt.NewNumericDate(now.Add(j.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        fmt.Sprintf("uid:%d", id),
		},
	}

	// Подписываем одним и тем же ключом, который ждет Middleware
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).
		SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).
		SignedString([]byte(j.secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (j *JwtTokenIssuer) VerifyToken(ctx context.Context, tokenStr string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}

	claims, ok := token.Claims.(*models.CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}