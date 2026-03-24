package repositories

import (
	"context"
	"fmt"
	"shb/internal/models"
)

func (r *Repository) GetAllInstitutions(ctx context.Context, search string, iType string, userLat, userLng float64, sortBy string) ([]*models.Institution, error) {
	// Базовый SELECT с подсчетом активных нужд
	// Если переданы координаты (не 0), считаем дистанцию (в километрах) по формуле Haversine
	var args []interface{}
	idx := 1

	distanceSelect := "0 as distance"
	if userLat != 0 && userLng != 0 {
		distanceSelect = fmt.Sprintf(`
			(6371 * acos(
				cos(radians($%d)) * cos(radians(latitude)) * cos(radians(longitude) - radians($%d)) + 
				sin(radians($%d)) * sin(radians(latitude))
			)) as distance`, idx, idx+1, idx)
		args = append(args, userLat, userLng)
		idx += 2
	}

	query := fmt.Sprintf(`
		SELECT 
			i.id, i.name, i.type, i.city, i.region, i.address, 
			i.phone, i.email, i.description, i.activity_hours, 
			i.latitude, i.longitude, i.wards_count, i.created_at, i.updated_at,
			(SELECT COUNT(*) FROM needs n WHERE n.institution_id = i.id AND n.is_deleted = false) as needs_count,
			%s
		FROM institutions i
		WHERE i.is_deleted = false
	`, distanceSelect)

	// 1. Поиск (Название ИЛИ Город)
	if search != "" {
		query += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.city ILIKE $%d)", idx, idx)
		args = append(args, "%"+search+"%")
		idx++
	}

	// 2. Фильтр по типу
	if iType != "" && iType != "all" {
		query += fmt.Sprintf(" AND i.type = $%d", idx)
		args = append(args, iType)
		idx++
	}

	// 3. Сортировка
	switch sortBy {
	case "needs_desc": // Сначала те, у кого больше нужд
		query += " ORDER BY needs_count DESC, i.id DESC"
	case "distance": // Сначала ближайшие
		if userLat != 0 && userLng != 0 {
			query += " ORDER BY distance ASC"
		} else {
			query += " ORDER BY i.id DESC" // Фоллбэк
		}
	default:
		// По умолчанию просто новые
		query += " ORDER BY i.id DESC"
	}

	r.logger.Debug().Str("query", query).Msg("GetAllInstitutions SQL")

	rows, err := r.postgres.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query institutions: %w", err)
	}
	defer rows.Close()

	var institutions []*models.Institution
	for rows.Next() {
		var i dbInstitution
    	var needsCount int
  		var distance float64

		// ВАЖНО: Порядок переменных должен строго соответствовать порядку в SELECT
		if err := rows.Scan(
			&i.ID, &i.Name, &i.Type, &i.City, &i.Region, &i.Address,
			&i.Phone, &i.Email, &i.Description, &i.ActivityHours,
			&i.Latitude, &i.Longitude, &i.WardsCount, &i.CreatedAt, &i.UpdatedAt,
			&needsCount, &distance,
		); err != nil {
			return nil, fmt.Errorf("scan institution: %w", err)
		}

		// Используем маппер
		domainInst := i.ToDomain()
		domainInst.NeedsCount = needsCount
		institutions = append(institutions, domainInst)
	}

	

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
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
        SELECT id, name, type, city, region, address, phone, email, description, activity_hours,
               latitude, longitude, wards_count, created_at, updated_at, is_deleted, deleted_at
        FROM institutions
        WHERE id = $1
    `

	var dbI dbInstitution
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&dbI.ID, &dbI.Name, &dbI.Type, &dbI.City, &dbI.Region, &dbI.Address,
		&dbI.Phone, &dbI.Email, &dbI.Description, &dbI.ActivityHours,
		&dbI.Latitude, &dbI.Longitude, &dbI.WardsCount, &dbI.CreatedAt, &dbI.UpdatedAt,
		&dbI.IsDeleted, &dbI.DeletedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get institution by id: %w", err)
	}

	// Конвертируем dbInstitution → models.Institution через маппер
	return dbI.ToDomain(), nil
}
