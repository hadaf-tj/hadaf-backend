// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package repositories

import (
	"context"
	"fmt"

	"shb/internal/models"
)

const donationHistoryFilter = `
	FROM ledger l
	JOIN medical_campaigns mc ON mc.id = l.campaign_id
	WHERE l.donor_user_id = $1
	  AND l.type = $2
`

const donationHistoryOrder = "ORDER BY l.created_at DESC, l.id DESC"

// GetDonationsByUser returns only recorded donation ledger entries belonging to
// one donor. Campaign rows are deliberately not filtered by is_deleted: a
// donor's financial history must remain visible after a campaign is archived.
func (r *Repository) GetDonationsByUser(
	ctx context.Context,
	userID, limit, offset int,
) (*models.DonationHistoryPage, error) {
	page := &models.DonationHistoryPage{
		Items:  make([]*models.DonationHistoryItem, 0),
		Limit:  limit,
		Offset: offset,
	}

	if err := r.postgres.QueryRow(ctx, "SELECT COUNT(*) "+donationHistoryFilter,
		userID, models.LedgerTypeDonation,
	).Scan(&page.Total); err != nil {
		return nil, fmt.Errorf("count donations by user: %w", err)
	}

	rows, err := r.postgres.Query(ctx, `
		SELECT l.id, l.campaign_id, mc.title, l.amount, l.currency, l.created_at
		`+donationHistoryFilter+`
		`+donationHistoryOrder+`
		LIMIT $3 OFFSET $4
	`, userID, models.LedgerTypeDonation, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("get donations by user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.DonationHistoryItem{
			Status: models.DonationStatusConfirmed,
		}
		if err := rows.Scan(
			&item.ID,
			&item.CampaignID,
			&item.CampaignTitle,
			&item.Amount,
			&item.Currency,
			&item.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan donation history item: %w", err)
		}
		page.Items = append(page.Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate donation history: %w", err)
	}

	return page, nil
}
