package services

import (
	"context"
	"shb/internal/models"
	"shb/internal/repositories/filters"
)

func (s *Service) GetAllInstitutions(ctx context.Context, filter filters.InstitutionFilter) ([]*models.Institution, error) {
	// Здесь можно добавить логирование или кэширование
	return s.repo.GetAllInstitutions(ctx, filter)
}

func (s *Service) CreateInstitution(ctx context.Context, i *models.Institution) (int, error) {
	return s.repo.CreateInstitution(ctx, i)
}

func (s *Service) GetInstitutionByID(ctx context.Context, id int) (*models.Institution, error) {
	return s.repo.GetInstitutionByID(ctx, id)
}
