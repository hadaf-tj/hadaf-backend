package services

import (
	"context"
	"shb/internal/models"
)

func (s *Service) GetAllInstitutions(ctx context.Context, q models.InstitutionListQuery) (*models.InstitutionPage, error) {
	return s.repo.GetAllInstitutions(ctx, q)
}

func (s *Service) CreateInstitution(ctx context.Context, i *models.Institution) (int, error) {
	return s.repo.CreateInstitution(ctx, i)
}

func (s *Service) GetInstitutionByID(ctx context.Context, id int) (*models.Institution, error) {
	return s.repo.GetInstitutionByID(ctx, id)
}
