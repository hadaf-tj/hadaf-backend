// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package sms

import (
	"context"
	"shb/pkg/external/sms/smsProvider"
)

// ISmsAdapter defines the contract for sending SMS messages.
type ISmsAdapter interface {
	// SendSms delivers an SMS to the given phone number.
	// txnID is the unique transaction identifier (typically the OTP record ID from the DB).
	SendSms(ctx context.Context, phone, message, txnID string) error
	// CheckBalance returns the current SMS account balance.
	CheckBalance(ctx context.Context) (*smsProvider.BalanceResult, error)
}
