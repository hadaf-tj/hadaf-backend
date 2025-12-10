package repositories

import (
	"time"
)

type dbUser struct {
	ID            int       `db:"id"`
	InstitutionID *int      `db:"institution_id"`
	FullName      string    `db:"full_name"`
	Phone         string    `db:"phone"`
	Email         string    `db:"email"`
	Password      string    `db:"password"`
	Role          string    `db:"role"`
	IsActive      bool      `db:"is_active"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type dbOtp struct {
	ID         int       `db:"id"`
	Attempt    int       `db:"attempt"`
	Receiver   string    `db:"receiver"`
	Method     string    `db:"method"`
	OTPCode    string    `db:"otp_code"`
	IsVerified bool      `db:"is_verified"`
	SentAt     time.Time `db:"sent_at"`
	ExpiresAt  time.Time `db:"expires_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

// Новая структура для учреждений
type dbInstitution struct {
	ID            int        `db:"id"`
	Name          string     `db:"name"`
	Type          string     `db:"type"`
	City          string     `db:"city"`
	Region        string     `db:"region"`
	Address       string     `db:"address"`
	Phone         string     `db:"phone"`
	Email         string     `db:"email"`
	Description   string     `db:"description"`
	ActivityHours string     `db:"activity_hours"`
	Latitude      *float64   `db:"latitude"`
	Longitude     *float64   `db:"longitude"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
}