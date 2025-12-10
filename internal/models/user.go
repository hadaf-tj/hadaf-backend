package models

import "time"

// Роли пользователей
const (
	RoleSuperAdmin = "super_admin"
	RoleEmployee   = "employee"
	RoleDonor      = "donor"
)

type User struct {
	ID            int       `db:"id"`
	InstitutionID *int      `db:"institution_id"` // Новое поле (может быть nil)
	FullName      string    `db:"full_name"`
	Phone         string    `db:"phone"`
	Email         string    `db:"email"`
	Password      string    `db:"password"`
	Role          string    `db:"role"`
	IsActive      bool      `db:"is_active"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type OTP struct {
	ID         int
	Attempt    int
	Receiver   string
	Method     string
	OTPCode    string
	IsVerified bool
	SentAt     time.Time
	ExpiresAt  time.Time
}