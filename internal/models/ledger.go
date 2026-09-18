// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package models

import "time"

const (
	LedgerTypeDonation          = "donation"
	LedgerTypePaymentToProvider = "payment_to_provider"
	LedgerTypeOverflowToGeneral = "overflow_to_general"
	LedgerTypeRefund            = "refund"

	// DonationStatusConfirmed is the only donation status currently represented
	// by ledger rows. Pending payment attempts are not persisted in this schema.
	DonationStatusConfirmed = "confirmed"
)

type LedgerEntry struct {
	ID           int       `json:"id"`
	CampaignID   *int      `json:"campaign_id"`
	DonorUserID  *int      `json:"donor_user_id"`
	Type         string    `json:"type"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	PaymentRef   *string   `json:"payment_ref"`
	Description  *string   `json:"description"`
	IsAnonymous  bool      `json:"is_anonymous"`
	DonorMessage *string   `json:"donor_message"`
	CreatedAt    time.Time `json:"created_at"`
}

// DonationHistoryItem is the safe donor-facing view of a recorded donation.
// It intentionally excludes donor identifiers, payment references, messages,
// and internal accounting descriptions.
type DonationHistoryItem struct {
	ID            int       `json:"id"`
	CampaignID    int       `json:"campaign_id"`
	CampaignTitle string    `json:"campaign_title"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// DonationHistoryPage is a paginated donor-facing donation history.
type DonationHistoryPage struct {
	Items  []*DonationHistoryItem `json:"items"`
	Total  int                    `json:"total"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
}
