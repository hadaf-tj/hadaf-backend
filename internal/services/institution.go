package services

import (
	"context"
	"shb/internal/models"
)

func (s *Service) GetAllInstitutions(ctx context.Context, search string, iType string, userLat, userLng float64, sortBy string) ([]*models.Institution, error) {
	return s.repo.GetAllInstitutions(ctx, search, iType, userLat, userLng, sortBy)
}

func (s *Service) CreateInstitution(ctx context.Context, i *models.Institution) (int, error) {
	return s.repo.CreateInstitution(ctx, i)
}

func (s *Service) GetInstitutionByID(ctx context.Context, id int) (*models.Institution, error) {
	return s.repo.GetInstitutionByID(ctx, id)
}
