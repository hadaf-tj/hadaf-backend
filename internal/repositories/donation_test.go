// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package repositories

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDonationHistoryQueryScopesAndOrdersLedgerEntries(t *testing.T) {
	require.Contains(t, donationHistoryFilter, "l.donor_user_id = $1")
	require.Contains(t, donationHistoryFilter, "l.type = $2")
	require.Contains(t, donationHistoryFilter, "JOIN medical_campaigns mc")
	require.NotContains(t, donationHistoryFilter, "mc.is_deleted")
	require.NotContains(t, donationHistoryFilter, "is_anonymous")
	require.Equal(t, "ORDER BY l.created_at DESC, l.id DESC", donationHistoryOrder)
}
