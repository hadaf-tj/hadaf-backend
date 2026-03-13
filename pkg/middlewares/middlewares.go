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
		var tokenString string

		// Try Authorization header first, then fall back to httpOnly cookie
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, myerrors.NewUnauthorizedErr("invalid auth header"))
				return
			}
			tokenString = headerParts[1]
		} else if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
			tokenString = cookie
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, myerrors.NewUnauthorizedErr("missing auth credentials"))
			return
		}
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

		// Verification of moderation for employees
		if claims.Role == models.RoleEmployee && !claims.IsApproved {
			c.AbortWithStatusJSON(http.StatusForbidden, myerrors.NewForbiddenErr("Ваш аккаунт ожидает подтверждения администратором"))
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("isApproved", claims.IsApproved)
		c.Next()
	}
}

// OptionalAccessToken - мягкая авторизация (не требует токена, но извлекает userID если есть)
func (m *Middleware) OptionalAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Try Authorization header first, then fall back to httpOnly cookie
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
				tokenString = headerParts[1]
			}
		}
		if tokenString == "" {
			if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
				tokenString = cookie
			}
		}
		if tokenString == "" {
			c.Set("userID", 0)
			c.Next()
			return
		}
		claims := &models.CustomClaims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.Set("userID", 0)
			c.Next()
			return
		}

		// Verification of moderation for employees (even in optional)
		if claims.Role == models.RoleEmployee && !claims.IsApproved {
			// for optional access token we just act as if they are not logged in
			c.Set("userID", 0)
			c.Next()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("isApproved", claims.IsApproved)
		c.Next()
	}
}
