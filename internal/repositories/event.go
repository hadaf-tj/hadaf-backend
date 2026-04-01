package repositories

import (
	"context"
	"errors"
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"

	"github.com/jackc/pgx/v5"
)

// CreateEvent создаёт новое событие
func (r *Repository) CreateEvent(ctx context.Context, e *models.Event) (int, error) {
	query := `
		INSERT INTO events (title, description, event_date, institution_id, creator_id, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	var id int
	err := r.postgres.QueryRow(ctx, query,
		e.Title, e.Description, e.EventDate, e.InstitutionID, e.CreatorID, e.Status,
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

// GetEventDetail возвращает одно событие в формате EventResponse (как элемент списка).
func (r *Repository) GetEventDetail(ctx context.Context, q models.EventDetailQuery) (*models.EventResponse, error) {
	const query = `
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
		WHERE e.is_deleted = false AND e.id = $2
	`

	var e models.EventResponse
	err := r.postgres.QueryRow(ctx, query, q.ViewerUserID, q.EventID).Scan(
		&e.ID, &e.Title, &e.Description, &e.EventDate, &e.InstitutionID, &e.InstitutionName,
		&e.CreatorID, &e.CreatorName, &e.ParticipantsCount, &e.IsJoined, &e.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerrors.ErrNotFound
		}
		return nil, fmt.Errorf("get event detail: %w", err)
	}
	return &e, nil
}

func (r *Repository) countEvents(ctx context.Context) (int64, error) {
	query := `
		SELECT COUNT(*)
		FROM events e
		JOIN institutions i ON e.institution_id = i.id
		JOIN users u ON e.creator_id = u.id
		WHERE e.is_deleted = false
	`
	var total int64
	err := r.postgres.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("count events: %w", err)
	}
	return total, nil
}

// GetAllEvents возвращает страницу событий с дополнительными данными.
func (r *Repository) GetAllEvents(ctx context.Context, q models.EventListQuery) (*models.EventPage, error) {
// GetAllEvents получает все события с дополнительными данными (публичная лента)
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
			e.status,
			e.created_at
			e.created_at,
			COUNT(*) OVER() AS total_count
		FROM events e
		JOIN institutions i ON e.institution_id = i.id
		JOIN users u ON e.creator_id = u.id
		WHERE e.is_deleted = false AND e.status = 'approved'
		ORDER BY e.event_date ASC
		WHERE e.is_deleted = false
		ORDER BY e.event_date ASC, e.id ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.postgres.Query(ctx, query, q.UserID, q.Limit, q.Offset)
	if err != nil {
		return nil, fmt.Errorf("get events: %w", err)
	}
	defer rows.Close()

	events := make([]*models.EventResponse, 0)
	var total int64

	for rows.Next() {
		var e models.EventResponse
		if err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.EventDate, &e.InstitutionID, &e.InstitutionName,
			&e.CreatorID, &e.CreatorName, &e.ParticipantsCount, &e.IsJoined, &e.CreatedAt,
			&total,
			&e.CreatorID, &e.CreatorName, &e.ParticipantsCount, &e.IsJoined, &e.Status, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan event: %w", err)
		}
		events = append(events, &e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}

	if len(events) == 0 {
		var errCount error
		total, errCount = r.countEvents(ctx)
		if errCount != nil {
			return nil, errCount
		}
	}

	return &models.EventPage{
		Items:  events,
		Total:  total,
		Limit:  q.Limit,
		Offset: q.Offset,
	}, nil
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

// GetInstitutionEvents получает все события, предложенные конкретному учреждению (для модерации)
func (r *Repository) GetInstitutionEvents(ctx context.Context, institutionID int) ([]*models.EventResponse, error) {
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
			false as is_joined,
			e.status,
			e.created_at
		FROM events e
		JOIN institutions i ON e.institution_id = i.id
		JOIN users u ON e.creator_id = u.id
		WHERE e.is_deleted = false AND e.institution_id = $1
		ORDER BY e.created_at DESC
	`

	rows, err := r.postgres.Query(ctx, query, institutionID)
	if err != nil {
		return nil, fmt.Errorf("get institution events: %w", err)
	}
	defer rows.Close()

	events := make([]*models.EventResponse, 0)
	for rows.Next() {
		var e models.EventResponse
		if err := rows.Scan(
			&e.ID, &e.Title, &e.Description, &e.EventDate, &e.InstitutionID, &e.InstitutionName,
			&e.CreatorID, &e.CreatorName, &e.ParticipantsCount, &e.IsJoined, &e.Status, &e.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan inst event: %w", err)
		}
		events = append(events, &e)
	}

	return events, nil
}

// UpdateEventStatus обновляет статус события
func (r *Repository) UpdateEventStatus(ctx context.Context, eventID int, status string) error {
	query := `UPDATE events SET status = $1, updated_at = NOW() WHERE id = $2 AND is_deleted = false`
	result, err := r.postgres.Exec(ctx, query, status, eventID)
	if err != nil {
		return fmt.Errorf("update event status: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("event not found")
	}
	return nil
}
