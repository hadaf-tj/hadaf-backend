// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package models

// InstitutionPage is a paginated response for the institution list endpoint.
type InstitutionPage struct {
	Items  []*Institution `json:"items"`
	Total  int64          `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

// EventPage is a paginated response for the event list endpoint.
type EventPage struct {
	Items  []*EventResponse `json:"items"`
	Total  int64            `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}
