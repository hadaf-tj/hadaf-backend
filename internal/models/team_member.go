package models

import "time"

// TeamMember модель участника команды
type TeamMember struct {
	ID        int        `json:"id"`
	FullName  string     `json:"full_name"`
	Role      string     `json:"role"`
	PhotoURL  *string    `json:"photo_url"`
	Quote     *string    `json:"quote"`
	Telegram  string     `json:"telegram"`
	LinkedIn  string     `json:"linkedin"`
	SortOrder int        `json:"sort_order"`
	IsActive  bool       `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
