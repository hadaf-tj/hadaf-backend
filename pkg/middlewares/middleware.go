package middlewares

import (
	"net/http"
	"shb/internal/configs"
	"shb/internal/models"
	"shb/pkg/constants"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	accessSecret  string
	refreshSecret string
}

func NewMiddleware() *Middleware {
	cfg, err := configs.InitConfigs()
	if err != nil {
		panic("failed init new middleware: " + err.Error())
	}

	return &Middleware{
		accessSecret:  cfg.Security.AccessTokenSecret,
		refreshSecret: cfg.Security.RefreshTokenSecret,
	}
}

func (m *Middleware) AccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := m.extractTokenFromHeader(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims := &models.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.accessSecret), nil
		})

		if err != nil || !token.Valid || claims.Subject != constants.AccessSubject {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("claims", claims)
		c.Next()
	}
}

func (m *Middleware) RefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := m.extractTokenFromHeader(c)
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims := &models.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.refreshSecret), nil
		})

		if err != nil || !token.Valid || claims.Subject != constants.RefreshSubject {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("claims", claims)
		c.Next()
	}
}

func (m *Middleware) extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}
	return ""
}

func (m *Middleware) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,PATCH,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding"+
			", X-CSRF-Token, Authorization, lang, Accept, accept")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}

		c.Next()
	}
}

// OptionalAccessToken - мягкая авторизация (не требует токена, но извлекает userID если есть)
func (m *Middleware) OptionalAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := m.extractTokenFromHeader(c)
		if tokenStr == "" {
			c.Set("user_id", 0)
			c.Next()
			return
		}

		claims := &models.CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(m.accessSecret), nil
		})

		if err != nil || !token.Valid || claims.Subject != constants.AccessSubject {
			c.Set("user_id", 0)
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("claims", claims)
		c.Next()
	}
}

