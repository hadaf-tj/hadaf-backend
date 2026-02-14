package models

import "time"

// Booking представляет отклик волонтера на нужду учреждения
type Booking struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	NeedID    int        `json:"need_id"`
	Quantity  float64    `json:"quantity"`
	Note      string     `json:"note"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	IsDeleted bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// Константы статусов бронирования
const (
	BookingStatusPending   = "pending"
	BookingStatusApproved  = "approved"
	BookingStatusRejected  = "rejected"
	BookingStatusCompleted = "completed"
)
