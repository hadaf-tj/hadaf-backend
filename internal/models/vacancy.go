// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package models

import "time"

// Vacancy represents a volunteer vacancy posting.
type Vacancy struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Type        string     `json:"type"`       // e.g. "Volunteering"
	Experience  string     `json:"experience"` // e.g. "1+ year"
	Workload    string     `json:"workload"`   // e.g. "1 hour per week"
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	IsDeleted   bool       `json:"-"`
}
