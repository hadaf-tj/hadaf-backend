// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package repositories

import (
	"context"
	"fmt"
)

// GetPublicStats returns aggregated statistics.
func (r *Repository) GetPublicStats(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT
			(SELECT COUNT(*) FROM needs WHERE received_qty >= required_qty AND is_deleted = false) AS closed_needs,
			(SELECT COALESCE(SUM(received_qty)::int, 0) FROM needs WHERE is_deleted = false) AS people_helped,
			(SELECT COUNT(*) FROM institutions WHERE is_deleted = false) AS institutions_count
	`

	var closedNeeds, peopleHelped, institutionsCount int
	err := r.postgres.QueryRow(ctx, query).Scan(&closedNeeds, &peopleHelped, &institutionsCount)
	if err != nil {
		return nil, fmt.Errorf("get public stats: %w", err)
	}

	return map[string]int{
		"closed_needs":       closedNeeds,
		"people_helped":      peopleHelped,
		"institutions_count": institutionsCount,
	}, nil
}
