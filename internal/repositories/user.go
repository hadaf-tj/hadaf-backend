package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"
)

func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (*models.User, error) {
	// Добавляем institution_id в выборку (SELECT *)
	const query = `
        SELECT id, institution_id, full_name, phone, email, password, role, is_active, created_at, updated_at
        FROM users
        WHERE phone = $1
        LIMIT 1;
    `

	// Используем $1 вместо :phone, так как pgx по умолчанию использует $n
	// Если ваш драйвер требует именованных параметров, оставьте как было,
	// но убедитесь, что структура dbUser соответствует полям.
	// В предыдущем коде использовались именованные, но стандартный pgx QueryRow обычно работает с $1.
	// Но у вас ранее использовалось :phone, возможно используется обертка sqlx или pgx с поддержкой имен?
	// Судя по `r.postgres.QueryRow`, это чистый pgxpool. Он использует $1, $2.
	// Давайте поправим на стандартный $1 для надежности.

	var userDB dbUser
	// Scan должен перечислять поля в том же порядке, что и в SELECT,
	// ЛИБО использовать scany/pgxscan для автоматического маппинга.
	// Судя по прошлому коду `Scan(&userDB)`, вы используете библиотеку, которая умеет сканировать в структуру (например, scany или pgxutil)?
	// Или это pgx.Row.Scan? pgx.Row.Scan требует перечисления переменных, он не умеет в структуру сам по себе без доп. библиотек.
	// Если это чистый pgx, то код `Scan(&userDB)` не сработает напрямую без перечисления полей.
	// Предположим, что у вас настроен scany или аналог. Если нет - скажите, перепишем на явное перечисление.

	// Оставим пока как было, но добавим поле в маппинг, если вы используете scany.
	// Если вы используете чистый pgx, то нужно перечислять поля.
	// Давайте напишем безопасно с перечислением:

	err := r.postgres.QueryRow(ctx, query, phone).Scan(
		&userDB.ID,
		&userDB.InstitutionID,
		&userDB.FullName,
		&userDB.Phone,
		&userDB.Email,
		&userDB.Password,
		&userDB.Role,
		&userDB.IsActive,
		&userDB.CreatedAt,
		&userDB.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // pgx возвращает pgx.ErrNoRows, но sql.ErrNoRows тоже часто используется в абстракциях
			return nil, myerrors.ErrNotFound
		}
		// Для чистого pgx проверка на no rows:
		if err.Error() == "no rows in result set" {
			return nil, myerrors.ErrNotFound
		}
		return nil, fmt.Errorf("get user by phone: %w", err)
	}

	return &models.User{
		ID:            userDB.ID,
		InstitutionID: userDB.InstitutionID,
		FullName:      userDB.FullName,
		Phone:         userDB.Phone,
		Password:      userDB.Password,
		Role:          userDB.Role,
		CreatedAt:     userDB.CreatedAt,
		UpdatedAt:     userDB.UpdatedAt,
	}, nil
}

func (r *Repository) CreateUser(ctx context.Context, user *models.User) error {
	const query = `
        INSERT INTO users (phone, role)
        VALUES ($1, $2)
        RETURNING id;
    `
	// По умолчанию роль 'donor'
	role := models.RoleDonor
	if user.Role != "" {
		role = user.Role
	}

	var id int64
	err := r.postgres.QueryRow(ctx, query, user.Phone, role).Scan(&id)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	user.ID = int(id)
	user.Role = role
	return nil
}