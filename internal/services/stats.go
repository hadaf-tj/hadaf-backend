// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services

import (
	"context"
	"fmt"
)

// GetPublicStats returns public statistics for the landing page.
func (s *Service) GetPublicStats(ctx context.Context) (map[string]int, error) {
	stats, err := s.repo.GetPublicStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get public stats: %w", err)
	}
	return stats, nil
}
