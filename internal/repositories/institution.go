package repositories

import (
	"context"
	"fmt"
	"shb/internal/models"
	"shb/internal/repositories/filters"
)

func (r *Repository) GetAllInstitutions(ctx context.Context, filter filters.InstitutionFilter) ([]*models.Institution, error) {
	// Добавляем подзапрос для подсчета (COUNT)
	query := `
        SELECT 
            i.id, i.name, i.type, i.city, i.region, i.address, 
            i.phone, i.email, i.description, i.activity_hours, 
            i.latitude, i.longitude, i.created_at, i.updated_at,
            (SELECT COUNT(*) FROM needs n WHERE n.institution_id = i.id) as needs_count 
        FROM institutions i
    `
	filterQuery, args := filters.BuildGetAllInstitutionFilter(filter)
	query += filterQuery

	rows, err := r.postgres.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query institutions: %w", err)
	}
	defer rows.Close()

	var institutions []*models.Institution
	for rows.Next() {
		var i dbInstitution

		if err := rows.Scan(
			&i.ID, &i.Name, &i.Type, &i.City, &i.Region, &i.Address,
			&i.Phone, &i.Email, &i.Description, &i.ActivityHours,
			&i.Latitude, &i.Longitude, &i.CreatedAt, &i.UpdatedAt,
			&i.NeedsCount,
		); err != nil {
			return nil, fmt.Errorf("scan institution: %w", err)
		}

		institutions = append(institutions, &models.Institution{
			ID:            i.ID,
			Name:          i.Name,
			Type:          i.Type,
			City:          i.City,
			Region:        i.Region,
			Address:       i.Address,
			Phone:         i.Phone,
			Email:         i.Email,
			Description:   i.Description,
			ActivityHours: i.ActivityHours,
			Latitude:      i.Latitude,
			Longitude:     i.Longitude,
			CreatedAt:     i.CreatedAt,
			UpdatedAt:     i.UpdatedAt,
			NeedsCount:    i.NeedsCount,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	return institutions, nil
}

// CreateInstitution вставляет новое учреждение в базу
func (r *Repository) CreateInstitution(ctx context.Context, i *models.Institution) (int, error) {
	query := `
		INSERT INTO institutions (name, type, city, region, address, phone, email, description, activity_hours, latitude, longitude)
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

func (r *Repository) GetInstitutionByID(ctx context.Context, id int) (*models.Institution, error) {
	query := `
        SELECT id, name, type, city, region, address, phone, email, description, activity_hours, latitude, longitude, created_at, updated_at
        FROM institutions
        WHERE id = $1
    `
	var i models.Institution
	// ВАЖНО: Убедись, что модель Institution совпадает с полями в базе
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&i.ID, &i.Name, &i.Type, &i.City, &i.Region, &i.Address,
		&i.Phone, &i.Email, &i.Description, &i.ActivityHours,
		&i.Latitude, &i.Longitude, &i.CreatedAt, &i.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get institution by id: %w", err)
	}
	return &i, nil
}
