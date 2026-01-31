package models

import "time"

type Event struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	EventDate     time.Time  `json:"event_date"`
	InstitutionID int        `json:"institution_id"`
	CreatorID     int        `json:"creator_id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	IsDeleted     bool       `json:"is_deleted"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

// EventResponse - расширенная модель для ответа API (с дополнительными полями)
type EventResponse struct {
	ID                int       `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	EventDate         time.Time `json:"event_date"`
	InstitutionID     int       `json:"institution_id"`
	InstitutionName   string    `json:"institution_name"`
	CreatorID         int       `json:"creator_id"`
	CreatorName       string    `json:"creator_name"`
	CreatorAvatar     *string   `json:"creator_avatar"`
	ParticipantsCount int       `json:"participants_count"`
	IsJoined          bool      `json:"is_joined"`
	CreatedAt         time.Time `json:"created_at"`
}

type EventParticipant struct {
	EventID int `json:"event_id"`
	UserID  int `json:"user_id"`
}
