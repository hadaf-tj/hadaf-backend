// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package customLimiter

import (
	"context"
	"fmt"
	"shb/pkg/db/cache"
	"strconv"
	"time"
)

type RateLimiter struct {
	cache cache.ICache
}

func NewRateLimiter(cache cache.ICache) *RateLimiter {
	return &RateLimiter{cache: cache}
}

func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, ttlMinutes int) (bool, error) {
	val, err := r.cache.Get(ctx, key)
	if err != nil && err.Error() != "redis: nil" {
		return false, err
	}

	if val == "" {
		// First request — initialise the counter and set the TTL.
		err = r.cache.Set(ctx, key, 1, time.Duration(ttlMinutes)*time.Minute)
		if err != nil {
			return false, err
		}
		return true, nil
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return false, fmt.Errorf("invalid rate limiter value: %w", err)
	}

	if count >= limit {
		return false, nil
	}

	err = r.cache.Increment(ctx, key)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RateLimiter) ResetAttempts(ctx context.Context, key string) error {
	return r.cache.Delete(ctx, key)
}
