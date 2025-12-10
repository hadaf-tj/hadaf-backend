package services

import (
	"context"
	"shb/internal/models"
	"shb/pkg/myerrors"
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

// CHANGED: Returns (int, error) now
func (s *Service) CreateNeed(ctx context.Context, need *models.Need) (int, error) {
	id, err := s.repo.CreateNeed(ctx, need)
	if err != nil {
		return 0, err
	}
	return id, nil 
}

func (s *Service) UpdateNeed(ctx context.Context, n *models.Need) error {
	existing, err := s.repo.GetNeedByID(ctx, n.ID)
	if err != nil {
		return err
	}
    // Uncomment when permissions are needed
	// if err := s.checkPermission(ctx, existing.InstitutionID); err != nil { return err }

	existing.Name = n.Name
	existing.Description = n.Description
	existing.Unit = n.Unit
	existing.RequiredQty = n.RequiredQty
	existing.ReceivedQty = n.ReceivedQty
	existing.Urgency = n.Urgency

	return s.repo.UpdateNeed(ctx, existing)
}

func (s *Service) DeleteNeed(ctx context.Context, id int) error {
	return s.repo.DeleteNeed(ctx, id)
}

func (s *Service) GetNeedsByInstitution(ctx context.Context, institutionID int) ([]*models.Need, error) {
	return s.repo.GetNeedsByInstitution(ctx, institutionID)
}