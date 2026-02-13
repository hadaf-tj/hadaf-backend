package repositories

import (
	"context"
	"fmt"
	"shb/internal/models"
)

func (r *Repository) GetAllInstitutions(ctx context.Context, search string, iType string, userLat, userLng float64, sortBy string) ([]*models.Institution, error) {
	// 1. Формируем SELECT с расчетом дистанции и количества нужд
	query := `
		SELECT 
			i.id, i.name, i.type, i.city, i.region, i.address, 
			i.phone, i.email, i.description, i.activity_hours, 
			i.latitude, i.longitude, i.created_at, i.updated_at,
			(SELECT COUNT(*) FROM needs n WHERE n.institution_id = i.id AND n.is_deleted = false) as needs_count
	`

	// Добавляем расчет дистанции, если координаты переданы
	if userLat != 0 && userLng != 0 {
		query += fmt.Sprintf(`, (6371 * acos(
				cos(radians(%f)) * cos(radians(i.latitude)) * cos(radians(i.longitude) - radians(%f)) + 
				sin(radians(%f)) * sin(radians(i.latitude))
			)) as distance`, userLat, userLng, userLat)
	} else {
		query += ", 0 as distance" // Заглушка, чтобы Scan не ломался, если координат нет
	}

	query += ` FROM institutions i WHERE i.is_deleted = false`

	// 2. Добавляем фильтры
	var args []interface{}
	idx := 1

	if search != "" {
		query += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.city ILIKE $%d)", idx, idx)
		args = append(args, "%"+search+"%")
		idx++
	}

	if iType != "" && iType != "all" {
		query += fmt.Sprintf(" AND i.type = $%d", idx)
		args = append(args, iType)
		idx++
	}

	// 3. Сортировка
	switch sortBy {
	case "distance":
		if userLat != 0 && userLng != 0 {
			query += " ORDER BY distance ASC"
		} else {
			query += " ORDER BY i.id DESC"
		}
	case "needs_desc":
		query += " ORDER BY needs_count DESC, i.id DESC"
	default:
		query += " ORDER BY i.id DESC"
	}

	// Логируем для отладки
	r.logger.Debug().Str("query", query).Msg("GetAllInstitutions SQL")

	rows, err := r.postgres.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query institutions: %w", err)
	}
	defer rows.Close()

	institutions := make([]*models.Institution, 0)
	
	for rows.Next() {
		// Используем структуру для сканирования основных полей
		// Но нам нужно создать переменные под поля, которые мы добавили в SELECT вручную
		var i dbInstitution 
		var needsCount int
		var distance float64 // Просто считываем, но пока не сохраняем в модель (если нужно - добавь Distance в модель)

		// ВАЖНО: Порядок переменных должен строго соответствовать порядку в SELECT
		err := rows.Scan(
			&i.ID, &i.Name, &i.Type, &i.City, &i.Region, &i.Address,
			&i.Phone, &i.Email, &i.Description, &i.ActivityHours,
			&i.Latitude, &i.Longitude, &i.CreatedAt, &i.UpdatedAt,
			&needsCount, &distance,
		)
		if err != nil {
			return nil, fmt.Errorf("scan institution: %w", err)
		}

		// Мапим из DB-модели в Domain-модель
		domainInst := i.ToDomain()
		
		// Вручную проставляем то, чего не было в dbInstitution
		domainInst.NeedsCount = needsCount
		
		institutions = append(institutions, domainInst)
	}

	return institutions, nil
}

// CreateInstitution вставляет новое учреждение в базу
func (r *Repository) CreateInstitution(ctx context.Context, i *models.Institution) (int, error) {
	query := `
		INSERT INTO institutions 
			(name, type, city, region, address, phone, email, description, activity_hours, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	var id int
	err := r.postgres.QueryRow(ctx, query,
		i.Name, i.Type, i.City, i.Region, i.Address,
		i.Phone, i.Email, i.Description, i.ActivityHours, i.Latitude, i.Longitude,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("create institution: %w", err)
	}

	return id, nil
}

// GetInstitutionByID получает учреждение по ID
func (r *Repository) GetInstitutionByID(ctx context.Context, id int) (*models.Institution, error) {
	query := `
        SELECT id, name, type, city, region, address, phone, email, description, activity_hours, latitude, longitude, created_at, updated_at, is_deleted, deleted_at
        FROM institutions
        WHERE id = $1
    `

	var dbI dbInstitution
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&dbI.ID, &dbI.Name, &dbI.Type, &dbI.City, &dbI.Region, &dbI.Address,
		&dbI.Phone, &dbI.Email, &dbI.Description, &dbI.ActivityHours,
		&dbI.Latitude, &dbI.Longitude, &dbI.CreatedAt, &dbI.UpdatedAt,
		&dbI.IsDeleted, &dbI.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get institution by id: %w", err)
	}

	// Конвертируем dbInstitution → models.Institution через маппер
	return dbI.ToDomain(), nil
}
