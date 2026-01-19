package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"shb/internal/models"
	"shb/pkg/myerrors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	jwtSecret string
}

// NewMiddleware теперь принимает секрет как аргумент
func NewMiddleware(jwtSecret string) *Middleware {
	return &Middleware{
		jwtSecret: jwtSecret,
	}
}

// AuthMiddleware (код остается прежним, но использует m.jwtSecret)
func (m *Middleware) AuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, myerrors.NewUnauthorizedErr("empty auth header"))
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, myerrors.NewUnauthorizedErr("invalid auth header"))
			return
		}

		tokenString := headerParts[1]
		claims := &models.CustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			// Используем локальное поле
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, myerrors.NewUnauthorizedErr("invalid token"))
			return
		}

		// RBAC
		if len(roles) > 0 {
			roleAllowed := false
			for _, role := range roles {
				if role == claims.Role {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				c.AbortWithStatusJSON(http.StatusForbidden, myerrors.NewForbiddenErr("access denied"))
				return
			}
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
