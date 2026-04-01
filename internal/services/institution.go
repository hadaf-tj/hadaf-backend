package services

import (
	"context"
	"shb/internal/models"

	"github.com/rs/zerolog"
)

func (s *Service) GetAllInstitutions(ctx context.Context, q models.InstitutionListQuery) (*models.InstitutionPage, error) {
	return s.repo.GetAllInstitutions(ctx, q)
}

func (s *Service) CreateInstitution(ctx context.Context, i *models.Institution) (int, error) {
	log := zerolog.Ctx(ctx).With().Str("service", "CreateInstitution").Logger()

	id, err := s.repo.CreateInstitution(ctx, i)
	if err != nil {
		return 0, err
	}

	log.Info().Int("institution_id", id).Str("name", i.Name).Msg("institution created")
	return id, nil
}

func (s *Service) GetInstitutionByID(ctx context.Context, id int) (*models.Institution, error) {
	return s.repo.GetInstitutionByID(ctx, id)
}
