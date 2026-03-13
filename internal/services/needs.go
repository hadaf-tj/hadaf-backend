package services

import (
	"context"
	"fmt"
	"shb/internal/models"
	"shb/internal/repositories/filters"
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

// CreateNeed creates a need — institution validation is done before this layer
func (s *Service) CreateNeed(ctx context.Context, need *models.Need) (int, error) {
	id, err := s.repo.CreateNeed(ctx, need)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Service) UpdateNeed(ctx context.Context, n *models.Need) error {
	// 1. Получаем текущее состояние (можно добавить проверку прав s.checkPermission)
	oldNeed, err := s.repo.GetNeedByID(ctx, n.ID)
	if err != nil {
		return err
	}

	// 2. Обновляем запись
	// Мапим поля, которые пришли, в старую запись (или используем n напрямую, если фронт шлет всё)
	oldNeed.Name = n.Name
	oldNeed.Description = n.Description
	oldNeed.Unit = n.Unit
	oldNeed.RequiredQty = n.RequiredQty
	oldNeed.ReceivedQty = n.ReceivedQty
	oldNeed.Urgency = n.Urgency

	if err := s.repo.UpdateNeed(ctx, oldNeed); err != nil {
		return err
	}

	// 3. Пишем историю
	comment := fmt.Sprintf("Updating data: %s. Progress: %.0f/%.0f", n.Name, n.ReceivedQty, n.RequiredQty)
	_ = s.repo.CreateNeedHistory(ctx, &models.NeedsHistory{
		NeedID:  n.ID,
		Comment: &comment,
	})

	return nil
}

func (s *Service) DeleteNeed(ctx context.Context, id int) error {
	// Можно добавить проверку прав здесь

	if err := s.repo.DeleteNeed(ctx, id); err != nil {
		return err
	}

	// Пишем историю удаления
	comment := "Need is deleted (to archive)"
	_ = s.repo.CreateNeedHistory(ctx, &models.NeedsHistory{
		NeedID:  id,
		Comment: &comment,
	})

	return nil
}

func (s *Service) GetNeedsByInstitution(ctx context.Context, filter filters.NeedsFilter, institutionID int) ([]*models.Need, error) {
	return s.repo.GetNeedsByInstitution(ctx, filter, institutionID)
}

func (s *Service) GetNeedByID(ctx context.Context, id int) (*models.Need, error) {
	return s.repo.GetNeedByID(ctx, id)
}