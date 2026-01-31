package repositories

import (
	"context"
	"fmt"
	"shb/internal/models"
)

// CreateEvent создаёт новое событие
func (r *Repository) CreateEvent(ctx context.Context, e *models.Event) (int, error) {
	query := `
		INSERT INTO events (title, description, event_date, institution_id, creator_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int
	err := r.postgres.QueryRow(ctx, query,
		e.Title, e.Description, e.EventDate, e.InstitutionID, e.CreatorID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("create event: %w", err)
	}
	return id, nil
}

// GetEventByID получает событие по ID
func (r *Repository) GetEventByID(ctx context.Context, id int) (*models.Event, error) {
	query := `
		SELECT id, title, description, event_date, institution_id, creator_id, created_at
		FROM events 
		WHERE id = $1 AND is_deleted = false
	`
	var e models.Event
	err := r.postgres.QueryRow(ctx, query, id).Scan(
		&e.ID, &e.Title, &e.Description, &e.EventDate, &e.InstitutionID, &e.CreatorID, &e.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get event: %w", err)
	}
	return &e, nil
}

// GetAllEvents получает все события с дополнительными данными
func (r *Repository) GetAllEvents(ctx context.Context, userID int) ([]*models.EventResponse, error) {
	query := `
		SELECT 
			e.id,
			e.title,
			e.description,
			e.event_date,
			e.institution_id,
			i.name as institution_name,
			e.creator_id,
			COALESCE(u.full_name, u.phone, 'Организатор') as creator_name,
			(SELECT COUNT(*) FROM event_participants ep WHERE ep.event_id = e.id) as participants_count,
			EXISTS(SELECT 1 FROM event_participants ep WHERE ep.event_id = e.id AND ep.user_id = $1) as is_joined,
			e.created_at
		FROM events e
		JOIN institutions i ON e.institution_id = i.id
		JOIN users u ON e.creator_id = u.id
		WHERE e.is_deleted = false
		ORDER BY e.event_date ASC
	`

	rows, err := r.postgres.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	defer rows.Close()

	events := make([]*models.EventResponse, 0)

	for rows.Next() {
		var e models.EventResponse
		if err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.EventDate, &e.InstitutionID, &e.InstitutionName,
			&e.CreatorID, &e.CreatorName, &e.ParticipantsCount, &e.IsJoined, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, &e)
	}

	return events, nil
}

// JoinEvent записывает пользователя на событие
func (r *Repository) JoinEvent(ctx context.Context, eventID, userID int) error {
	query := `
		INSERT INTO event_participants (event_id, user_id)
		VALUES ($1, $2)
		ON CONFLICT (event_id, user_id) DO NOTHING
	`
	_, err := r.postgres.Exec(ctx, query, eventID, userID)
	if err != nil {
		return fmt.Errorf("join event: %w", err)
	}
	return nil
}

// LeaveEvent отменяет запись пользователя на событие
func (r *Repository) LeaveEvent(ctx context.Context, eventID, userID int) error {
	query := `DELETE FROM event_participants WHERE event_id = $1 AND user_id = $2`
	result, err := r.postgres.Exec(ctx, query, eventID, userID)
	if err != nil {
		return fmt.Errorf("leave event: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("participant not found")
	}
	return nil
}
