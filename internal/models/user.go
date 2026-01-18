package models

import "time"

// Роли пользователей
const (
	RoleSuperAdmin = "super_admin"
	RoleEmployee   = "employee"
	RoleDonor      = "donor"
)

type User struct {
	ID            int        `json:"id"`
	InstitutionID *int       `json:"institution_id"`
	FullName      *string    `json:"full_name"`
	Phone         *string    `json:"phone"`
	Email         *string    `json:"email"`
	Password      *string    `json:"password"`
	Role          string     `json:"role"`
	IsActive      bool       `json:"is_active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	IsDeleted     bool       `json:"is_deleted"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type OTP struct {
	ID         int        `json:"id"`
	Attempt    int        `json:"attempt"`
	Receiver   string     `json:"receiver"`
	Method     *string    `json:"method"`
	OTPCode    string     `json:"otp_code"`
	IsVerified bool       `json:"is_verified"`
	SentAt     time.Time  `json:"sent_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	IsDeleted  bool       `json:"is_deleted"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
