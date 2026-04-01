package repositories

import (
	"context"
	"shb/internal/models"
)

// GetAllTeamMembers retrieves all active team members ordered by sort_order
func (r *Repository) GetAllTeamMembers(ctx context.Context) ([]*models.TeamMember, error) {
	query := `
		SELECT id, full_name, role, photo_url, quote, telegram, linkedin, sort_order, is_active, created_at, updated_at
		FROM team_members
		WHERE is_active = true
		ORDER BY sort_order ASC, id ASC
	`
	rows, err := r.postgres.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.TeamMember
	for rows.Next() {
		var m models.TeamMember
		err := rows.Scan(
			&m.ID, &m.FullName, &m.Role, &m.PhotoURL, &m.Quote, &m.Telegram, &m.LinkedIn, &m.SortOrder, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		members = append(members, &m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}
