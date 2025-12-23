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

func (j *JwtTokenIssuer) IssueTokens(ctx context.Context, id int) (string, string, error) {
	now := time.Now().UTC()

	// ВАЖНО: Middleware ожидает поле "role" в claims, иначе не пустит к защищенным роутам.
	// Пока хардкодим "employee", так как у нас сейчас все - сотрудники. 
	// В идеале метод должен принимать role аргументом.
	role := models.RoleEmployee 

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