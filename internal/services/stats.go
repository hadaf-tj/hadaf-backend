package services

import (
	"context"
	"fmt"
)

// GetPublicStats возвращает публичную статистику для лендинга
func (s *Service) GetPublicStats(ctx context.Context) (map[string]int, error) {
	stats, err := s.repo.GetPublicStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("get public stats: %w", err)
	}
	return stats, nil
}
