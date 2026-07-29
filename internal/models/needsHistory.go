// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package models

import "time"

type NeedsHistory struct {
	ID        int        `json:"id"`
	NeedID    int        `json:"need_id"`
	Comment   *string    `json:"comment"`
	CreatedAt time.Time  `json:"created_at"`
	IsDeleted bool       `json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at"`
}
