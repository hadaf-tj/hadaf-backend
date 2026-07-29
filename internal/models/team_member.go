// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package models

import "time"

// TeamMember represents a member of the Hadaf core team.
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
