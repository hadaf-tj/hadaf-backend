// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package email

import "context"

// IEmailAdapter defines the contract for sending email messages.
type IEmailAdapter interface {
	// SendEmail delivers an email to the specified address.
	SendEmail(ctx context.Context, to, subject, body string) error
}
