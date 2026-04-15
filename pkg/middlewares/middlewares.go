// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

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

// Middleware holds shared middleware state, such as the JWT signing secret.
type Middleware struct {
	jwtSecret string
}

// NewMiddleware creates a new Middleware instance with the given JWT secret.
func NewMiddleware(jwtSecret string) *Middleware {
	return &Middleware{
		jwtSecret: jwtSecret,
	}
}

// AuthMiddleware returns a Gin handler that enforces JWT authentication and,
// optionally, role-based access control. If roles are provided, the caller
// must have at least one of them; otherwise any valid token is accepted.
func (m *Middleware) AuthMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Try Authorization header first, then fall back to httpOnly cookie.
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
			return []byte(m.jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, myerrors.NewUnauthorizedErr("invalid token"))
			return
		}

		// Role-based access control.
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

		// Employees must be approved by a super-admin before they can access protected routes.
		if claims.Role == models.RoleEmployee && !claims.IsApproved {
			c.AbortWithStatusJSON(http.StatusForbidden, myerrors.NewForbiddenErr("ERR_ACCOUNT_PENDING_APPROVAL"))
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("isApproved", claims.IsApproved)
		c.Next()
	}
}

// OptionalAccessToken is a soft-authentication middleware. It attempts to
// extract and validate the access token but proceeds without error if one is
// absent or invalid, setting userID to 0.
func (m *Middleware) OptionalAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Try Authorization header first, then fall back to httpOnly cookie.
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

		// Unapproved employees are treated as unauthenticated in optional mode.
		if claims.Role == models.RoleEmployee && !claims.IsApproved {
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
