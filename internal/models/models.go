package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CustomClaims struct {
	UserID     int    `json:"user_id"`
	Role       string `json:"role"`
	IsApproved bool   `json:"is_approved"`

	jwt.RegisteredClaims
}
