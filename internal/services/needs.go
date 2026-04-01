package services

import (
	"context"
	"fmt"
	"shb/internal/models"
	"shb/internal/repositories/filters"
	"shb/pkg/myerrors"

	"github.com/rs/zerolog"
)

// checkPermission implementation...
func (s *Service) checkPermission(ctx context.Context, institutionID int) error {
	role, ok := ctx.Value("role").(string)
	if !ok {
		return myerrors.NewUnauthorizedErr("role not found in context")
	}
	if role == models.RoleSuperAdmin {
		return nil
	}

	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return myerrors.NewUnauthorizedErr("user id not found")
	}

	userDB, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if role != models.RoleEmployee {
		return myerrors.NewForbiddenErr("only employees can edit needs")
	}

	if userDB != nil && userDB.InstitutionID != nil && *userDB.InstitutionID != institutionID {
		return myerrors.NewForbiddenErr("you can only manage your own institution")
	}

	return nil
}

// CreateNeed creates a need — institution validation is done before this layer
func (s *Service) CreateNeed(ctx context.Context, need *models.Need) (int, error) {
	log := zerolog.Ctx(ctx).With().Str("service", "CreateNeed").Int("institution_id", need.InstitutionID).Logger()

	id, err := s.repo.CreateNeed(ctx, need)
	if err != nil {
		return 0, err
	}

	log.Info().Int("need_id", id).Str("name", need.Name).Msg("need created")
	return id, nil
}

func (s *Service) UpdateNeed(ctx context.Context, n *models.Need) error {
	log := zerolog.Ctx(ctx).With().Str("service", "UpdateNeed").Int("need_id", n.ID).Logger()

	oldNeed, err := s.repo.GetNeedByID(ctx, n.ID)
	if err != nil {
		return err
	}

	oldNeed.Name = n.Name
	oldNeed.Description = n.Description
	oldNeed.Unit = n.Unit
	oldNeed.RequiredQty = n.RequiredQty
	oldNeed.ReceivedQty = n.ReceivedQty
	oldNeed.Urgency = n.Urgency

	if err := s.repo.UpdateNeed(ctx, oldNeed); err != nil {
		return err
	}

	comment := fmt.Sprintf("Updating data: %s. Progress: %.0f/%.0f", n.Name, n.ReceivedQty, n.RequiredQty)
	_ = s.repo.CreateNeedHistory(ctx, &models.NeedsHistory{
		NeedID:  n.ID,
		Comment: &comment,
	})

	log.Info().Str("name", n.Name).Msg("need updated")
	return nil
}

func (s *Service) DeleteNeed(ctx context.Context, id int) error {
	log := zerolog.Ctx(ctx).With().Str("service", "DeleteNeed").Int("need_id", id).Logger()

	if err := s.repo.DeleteNeed(ctx, id); err != nil {
		return err
	}

	comment := "Need is deleted (to archive)"
	_ = s.repo.CreateNeedHistory(ctx, &models.NeedsHistory{
		NeedID:  id,
		Comment: &comment,
	})

	log.Info().Msg("need deleted")
	return nil
}

func (s *Service) GetNeedsByInstitution(ctx context.Context, filter filters.NeedsFilter, institutionID int) ([]*models.Need, error) {
	return s.repo.GetNeedsByInstitution(ctx, filter, institutionID)
}

func (s *Service) GetNeedByID(ctx context.Context, id int) (*models.Need, error) {
	return s.repo.GetNeedByID(ctx, id)
}
