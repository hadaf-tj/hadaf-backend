// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package models

// InstitutionListQuery holds pagination and filter parameters for the institution list endpoint.
type InstitutionListQuery struct {
	Search  string
	Type    string
	UserLat float64
	UserLng float64
	SortBy  string
	Limit   int
	Offset  int
}

// EventListQuery holds pagination and filter parameters for the event list endpoint.
type EventListQuery struct {
	UserID int
	Limit  int
	Offset int
}

// EventDetailQuery holds parameters for the event detail endpoint. ViewerUserID is 0 for unauthenticated requests.
type EventDetailQuery struct {
	EventID      int
	ViewerUserID int
}
