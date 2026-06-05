// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package tokens

import (
	"context"
	"shb/internal/models"
)

type ITokenIssuer interface {
	// IssueTokens issues a new access/refresh token pair for the given user.
	IssueTokens(ctx context.Context, id int, role string, isApproved bool) (string, string, error)
	// VerifyToken validates a token and returns its claims.
	VerifyToken(ctx context.Context, token string) (*models.CustomClaims, error)
}
