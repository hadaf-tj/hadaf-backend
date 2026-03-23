package repositories

import (
	"context"
	"fmt"
	"shb/internal/models"
)

// GetTeamMemberByID retrieves a single team member by ID
func (r *Repository) GetTeamMemberByID(ctx context.Context, id int) (*models.TeamMember, error) {
	query := `
		SELECT id, full_name, role, photo_url, quote, telegram, linkedin, sort_order, is_active, created_at, updated_at
		FROM team_members
		WHERE id = $1 AND is_active = true
	`
	var m models.TeamMember
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&m.ID, &m.FullName, &m.Role, &m.PhotoURL, &m.Quote, &m.Telegram, &m.LinkedIn, &m.SortOrder, &m.IsActive, &m.CreatedAt, &m.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get team member: %w", err)
	}
	return &m, nil
}
