package services

import (
	"context"
	"shb/internal/models"
)

// GetAllVacancies получает все активные вакансии
func (s *Service) GetAllVacancies(ctx context.Context) ([]*models.Vacancy, error) {
	return s.repo.GetAllVacancies(ctx)
}

// GetVacancyByID получает вакансию по ID
func (s *Service) GetVacancyByID(ctx context.Context, id int) (*models.Vacancy, error) {
	return s.repo.GetVacancyByID(ctx, id)
}
