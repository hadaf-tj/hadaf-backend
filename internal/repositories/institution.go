package repositories

import (
	"context"
	"fmt"
	"shb/internal/models"
)

func (r *Repository) countInstitutions(ctx context.Context, q models.InstitutionListQuery) (int64, error) {
	query := `SELECT COUNT(*) FROM institutions i WHERE i.is_deleted = false`
	var args []interface{}
	idx := 1

	if q.Search != "" {
		query += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.city ILIKE $%d)", idx, idx)
		args = append(args, "%"+q.Search+"%")
		idx++
	}
	if q.Type != "" && q.Type != "all" {
		query += fmt.Sprintf(" AND i.type = $%d", idx)
		args = append(args, q.Type)
	}

	var total int64
	err := r.postgres.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("count institutions: %w", err)
	}
	return total, nil
}

func (r *Repository) GetAllInstitutions(ctx context.Context, q models.InstitutionListQuery) (*models.InstitutionPage, error) {
	var args []interface{}
	idx := 1

	distanceSelect := "0::double precision as distance"
	if q.UserLat != 0 && q.UserLng != 0 {
		distanceSelect = fmt.Sprintf(`
			(6371 * acos(
				cos(radians($%d)) * cos(radians(latitude)) * cos(radians(longitude) - radians($%d)) + 
				sin(radians($%d)) * sin(radians(latitude))
			)) as distance`, idx, idx+1, idx)
		args = append(args, q.UserLat, q.UserLng)
		idx += 2
	}

	query := fmt.Sprintf(`
		SELECT 
			i.id, i.name, i.type, i.city, i.region, i.address, 
			i.phone, i.email, i.description, i.activity_hours, 
			i.latitude, i.longitude, i.created_at, i.updated_at,
			(SELECT COUNT(*) FROM needs n WHERE n.institution_id = i.id AND n.is_deleted = false) as needs_count,
			%s,
			COUNT(*) OVER() AS total_count
		FROM institutions i
		WHERE i.is_deleted = false
	`, distanceSelect)

	if q.Search != "" {
		query += fmt.Sprintf(" AND (i.name ILIKE $%d OR i.city ILIKE $%d)", idx, idx)
		args = append(args, "%"+q.Search+"%")
		idx++
	}

	if q.Type != "" && q.Type != "all" {
		query += fmt.Sprintf(" AND i.type = $%d", idx)
		args = append(args, q.Type)
		idx++
	}

	switch q.SortBy {
	case "needs_desc":
		query += " ORDER BY needs_count DESC, i.id DESC"
	case "distance":
		if q.UserLat != 0 && q.UserLng != 0 {
			query += " ORDER BY distance ASC, i.id ASC"
		} else {
			query += " ORDER BY i.id DESC"
		}
	default:
		query += " ORDER BY i.id DESC"
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", idx, idx+1)
	args = append(args, q.Limit, q.Offset)

	r.logger.Debug().Str("query", query).Msg("GetAllInstitutions SQL")

	rows, err := r.postgres.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query institutions: %w", err)
	}
	defer rows.Close()

	institutions := make([]*models.Institution, 0)
	var total int64

	for rows.Next() {
		var i dbInstitution
		var needsCount int
		var distance float64

		if err := rows.Scan(
			&i.ID, &i.Name, &i.Type, &i.City, &i.Region, &i.Address,
			&i.Phone, &i.Email, &i.Description, &i.ActivityHours,
			&i.Latitude, &i.Longitude, &i.CreatedAt, &i.UpdatedAt,
			&needsCount, &distance, &total,
		); err != nil {
			return nil, fmt.Errorf("scan institution: %w", err)
		}

		domainInst := i.ToDomain()
		domainInst.NeedsCount = needsCount
		institutions = append(institutions, domainInst)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if len(institutions) == 0 {
		var errCount error
		total, errCount = r.countInstitutions(ctx, q)
		if errCount != nil {
			return nil, errCount
		}
	}

	return &models.InstitutionPage{
		Items:  institutions,
		Total:  total,
		Limit:  q.Limit,
		Offset: q.Offset,
	}, nil
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
               latitude, longitude, created_at, updated_at, is_deleted, deleted_at
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

	return dbI.ToDomain(), nil
}
