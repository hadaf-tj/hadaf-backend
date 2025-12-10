package repositories

import (
	"context"
	"fmt"
	"shb/internal/models"
)

func (r *Repository) CreateNeed(ctx context.Context, n *models.Need) (int, error) {
	query := `
		INSERT INTO needs (institution_id, category_id, name, description, unit, required_qty, received_qty, urgency)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	var id int
	err := r.postgres.QueryRow(ctx, query,
		n.InstitutionID, n.CategoryID, n.Name, n.Description, n.Unit, n.RequiredQty, n.ReceivedQty, n.Urgency,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("create need: %w", err)
	}
	return id, nil
}

func (r *Repository) GetNeedByID(ctx context.Context, id int) (*models.Need, error) {
	query := `SELECT id, institution_id, name, unit, required_qty, received_qty, urgency FROM needs WHERE id = $1`
	var n models.Need
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&n.ID, &n.InstitutionID, &n.Name, &n.Unit, &n.RequiredQty, &n.ReceivedQty, &n.Urgency,
	)
	if err != nil {
		return nil, fmt.Errorf("get need: %w", err)
	}
	return &n, nil
}

func (r *Repository) UpdateNeed(ctx context.Context, n *models.Need) error {
	query := `
		UPDATE needs 
		SET name=$1, description=$2, unit=$3, required_qty=$4, received_qty=$5, urgency=$6, updated_at=NOW()
		WHERE id=$7
	`
	_, err := r.postgres.Exec(ctx, query,
		n.Name, n.Description, n.Unit, n.RequiredQty, n.ReceivedQty, n.Urgency, n.ID,
	)
	if err != nil {
		return fmt.Errorf("update need: %w", err)
	}
	return nil
}

func (r *Repository) DeleteNeed(ctx context.Context, id int) error {
	query := `DELETE FROM needs WHERE id = $1`
	_, err := r.postgres.Exec(ctx, query, id)
	return err
}

func (r *Repository) GetNeedsByInstitution(ctx context.Context, institutionID int) ([]*models.Need, error) {
	query := `
		SELECT id, institution_id, name, description, unit, required_qty, received_qty, urgency, created_at
		FROM needs
		WHERE institution_id = $1
		ORDER BY urgency = 'high' DESC, created_at DESC
	`
	rows, err := r.postgres.Query(ctx, query, institutionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var needs []*models.Need
	for rows.Next() {
		var n models.Need
		if err := rows.Scan(
			&n.ID, &n.InstitutionID, &n.Name, &n.Description, &n.Unit, &n.RequiredQty, &n.ReceivedQty, &n.Urgency, &n.CreatedAt,
		); err != nil {
			return nil, err
		}
		needs = append(needs, &n)
	}
	return needs, nil
}