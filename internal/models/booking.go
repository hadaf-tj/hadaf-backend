// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package models

import "time"

// Booking represents a volunteer's commitment to fulfil a specific institution need.
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

	// Enriched fields for frontend
	NeedName        string `json:"need_name"`
	InstitutionName string `json:"institution_name"`
	InstitutionID   int    `json:"institution_id"`
}

// Booking status constants.
const (
	BookingStatusPending   = "pending"
	BookingStatusApproved  = "approved"
	BookingStatusRejected  = "rejected"
	BookingStatusCompleted = "completed"
)
