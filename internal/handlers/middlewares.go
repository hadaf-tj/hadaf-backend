package handlers

import (
	"context"
	"fmt"
	"strings"

	"shb/internal/models"
	"shb/pkg/constants"
	"shb/pkg/myerrors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (h *Handler) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(constants.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, constants.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// AuthMiddleware проверяет JWT токен и роли доступа
func (h *Handler) AuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			h.handleError(c, myerrors.NewUnauthorizedErr("empty auth header"))
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			h.handleError(c, myerrors.NewUnauthorizedErr("invalid auth header"))
			return
		}

		tokenString := headerParts[1]
		claims := &models.CustomClaims{}

		// Парсим токен
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Используем секретный ключ из конфига
			return []byte(h.cfg.Security.JWTSecretKey), nil
		})

		if err != nil || !token.Valid {
			h.handleError(c, myerrors.NewUnauthorizedErr("invalid token"))
			return
		}

		// Проверка ролей (RBAC)
		if len(roles) > 0 {
			roleAllowed := false
			for _, role := range roles {
				if role == claims.Role {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				h.handleError(c, myerrors.NewForbiddenErr("access denied"))
				return
			}
		}

		// Сохраняем данные пользователя в контекст Gin для дальнейшего использования
		// "userID" и "role" - это ключи, по которым мы будем доставать данные в хендлерах
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}