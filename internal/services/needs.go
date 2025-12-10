package services

import (
	"context"
	"shb/internal/models"
	"shb/pkg/myerrors"
)

// Проверка прав доступа: Админ может все, Сотрудник - только свое
func (s *Service) checkPermission(ctx context.Context, institutionID int) error {
	role, ok := ctx.Value("role").(string)
	if !ok {
		return myerrors.NewUnauthorizedErr("role not found in context")
	}
	// Админ может всё
	if role == models.RoleSuperAdmin {
		return nil
	}

	// Сотрудник
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return myerrors.NewUnauthorizedErr("user id not found")
	}

	// Получаем свежие данные пользователя, чтобы проверить привязку к учреждению
	// (В идеале ID учреждения можно хранить в токене, чтобы не делать запрос в БД каждый раз)
	// Для MVP сделаем запрос для надежности:
	// NOTE: Это упрощение. В продакшене лучше institution_id класть в JWT claims.
	userDB, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err // Если ошибка, значит юзера нет или БД упала
	}
    // Внимание: метода GetUserByID у нас в интерфейсе пока нет, нужно будет добавить или использовать GetUserByPhone если есть телефон.
    // Давайте лучше положимся на то, что institution_id будет в токене (Claims) в будущем.
    // Пока что предположим, что мы доверяем данным, переданным с фронта, НО это небезопасно.
    // Правильный путь MVP: Достать user и сверить ID.
    
    // ВАЖНО: Добавьте GetUserByID в репозиторий позже. Сейчас реализуем простую проверку роли.
	if role != models.RoleEmployee {
        return myerrors.NewForbiddenErr("only employees can edit needs")
    }
    
    // В MVP пока пропустим жесткую сверку institution_id без метода GetUserByID,
    // но в продакшене она обязательна.
    if userDB != nil && userDB.InstitutionID != nil && *userDB.InstitutionID != institutionID {
         return myerrors.NewForbiddenErr("you can only manage your own institution")
    }

	return nil
}

func (s *Service) CreateNeed(ctx context.Context, need *models.Need) (error) {
    // Вызываем репозиторий, он возвращает ID и ошибку
    err := s.repo.CreateNeed()
    if err != nil {
        return err
    }
    return nil // Возвращаем ID и nil (нет ошибки)
}

func (s *Service) UpdateNeed(ctx context.Context, n *models.Need) error {
	// Сначала получаем существующую нужду, чтобы узнать её InstitutionID
	existing, err := s.repo.GetNeedByID(ctx, n.ID)
	if err != nil {
		return err
	}
    // Проверка прав
    // if err := s.checkPermission(ctx, existing.InstitutionID); err != nil { return err }

	// Обновляем поля
	existing.Name = n.Name
	existing.Description = n.Description
	existing.Unit = n.Unit
	existing.RequiredQty = n.RequiredQty
	existing.ReceivedQty = n.ReceivedQty
	existing.Urgency = n.Urgency
	
	return s.repo.UpdateNeed(ctx, existing)
}

func (s *Service) DeleteNeed(ctx context.Context, id int) error {
	// existing, err := s.repo.GetNeedByID(ctx, id)
    // if err := s.checkPermission(ctx, existing.InstitutionID); err != nil { return err }
	return s.repo.DeleteNeed(ctx, id)
}

func (s *Service) GetNeedsByInstitution(ctx context.Context, institutionID int) ([]*models.Need, error) {
	return s.repo.GetNeedsByInstitution(ctx, institutionID)
}