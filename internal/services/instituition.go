package services

import (
	"context"
	"shb/internal/models"
)

func (s *Service) GetAllInstitutions(ctx context.Context, city string) ([]*models.Institution, error) {
	// Здесь можно добавить логирование или кэширование
	return s.repo.GetAllInstitutions(ctx, city)
}

func (s *Service) CreateInstitution(ctx context.Context, i *models.Institution) (int, error) {
	return s.repo.CreateInstitution(ctx, i)
}