package middlewares

import (
	"fmt"
	"net/http"
	"shb/internal/configs"
	"shb/internal/models"
	"shb/pkg/myerrors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	cfg *configs.Config
}

func NewMiddleware(cfg *configs.Config) *Middleware {
	return &Middleware{cfg: cfg}
}

// AuthMiddleware проверяет JWT и роли
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
			return []byte(m.cfg.Security.JWTSecretKey), nil
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