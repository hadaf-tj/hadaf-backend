// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package rateLimiter

import "context"

// IRateLimiter defines the rate-limiting contract.
type IRateLimiter interface {
	// Allow checks whether the caller identified by key is within the allowed
	// limit for the given time window.
	Allow(ctx context.Context, key string, limit int, ttlMinutes int) (bool, error)

	// ResetAttempts clears the stored attempt counter for the given key.
	ResetAttempts(ctx context.Context, key string) error
}
