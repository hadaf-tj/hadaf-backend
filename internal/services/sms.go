// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services

import (
	"context"
	"shb/pkg/external/sms/smsProvider"
)

// CheckSMSBalance checks the SMS account balance
func (s *Service) CheckSMSBalance(ctx context.Context) (*smsProvider.BalanceResult, error) {
	return s.sms.CheckBalance(ctx)
}
