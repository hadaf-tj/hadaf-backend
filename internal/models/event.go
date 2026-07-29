// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package models

import "time"

type Event struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	EventDate     time.Time  `json:"event_date"`
	InstitutionID int        `json:"institution_id"`
	CreatorID     int        `json:"creator_id"`
	Status        string     `json:"status"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	IsDeleted     bool       `json:"is_deleted"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

// EventResponse is the extended event model returned by the API (includes computed fields).
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
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

type EventParticipant struct {
	EventID int `json:"event_id"`
	UserID  int `json:"user_id"`
}
