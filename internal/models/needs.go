// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package models

import "time"

type Need struct {
	ID            int        `json:"id"`
	InstitutionID int        `json:"institution_id"`
	CategoryID    *int       `json:"category_id"`
	Name          string     `json:"name"`
	Description   *string    `json:"description"`
	Unit          string     `json:"unit"`
	RequiredQty   float64    `json:"required_qty"`
	ReceivedQty   float64    `json:"received_qty"`
	BookedQty     float64    `json:"booked_qty"`
	Urgency       string     `json:"urgency"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	IsDeleted     bool       `json:"is_deleted"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type UpdateNeedInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *string `json:"priority"`
	// Pointer fields allow partial updates: nil means "do not update this field".
}
