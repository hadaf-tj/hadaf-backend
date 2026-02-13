package services

import (
	"context"
	"shb/internal/models"
	"shb/internal/repositories/filters"
)

func (s *Service) GetAllInstitutions(ctx context.Context, filter filters.InstitutionFilter) ([]*models.Institution, error) {
	// Здесь можно добавить логирование или кэширование
	var lat, lng float64
	if filter.Lat != nil {
		lat = *filter.Lat
	}
	if filter.Lng != nil {
		lng = *filter.Lng
	}
	return s.repo.GetAllInstitutions(ctx, filter.Name, filter.Type, lat, lng, filter.OrderBy)
}

func (s *Service) CreateInstitution(ctx context.Context, i *models.Institution) (int, error) {
	return s.repo.CreateInstitution(ctx, i)
}

func (s *Service) GetInstitutionByID(ctx context.Context, id int) (*models.Institution, error) {
	return s.repo.GetInstitutionByID(ctx, id)
}
