// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package repositories

import (
	"context"
	"errors"
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"

	"github.com/jackc/pgx/v5"
)

// CreateEvent creates a new event.
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

// GetEventByID gets an event by ID.
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

// GetEventDetail returns an event formatted as EventResponse.
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
			COALESCE(u.full_name, u.phone, 'Organizer') as creator_name,
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

// GetAllEvents returns a page of events with additional metadata.
func (r *Repository) GetAllEvents(ctx context.Context, q models.EventListQuery) (*models.EventPage, error) {
	query := `
		SELECT
			e.id,
			e.title,
			e.description,
			e.event_date,
			e.institution_id,
			i.name as institution_name,
			e.creator_id,
			COALESCE(u.full_name, u.phone, 'Organizer') as creator_name,
			(SELECT COUNT(*) FROM event_participants ep WHERE ep.event_id = e.id) as participants_count,
			EXISTS(SELECT 1 FROM event_participants ep WHERE ep.event_id = e.id AND ep.user_id = $1) as is_joined,
			e.status,
			e.created_at,
			COUNT(*) OVER() AS total_count
		FROM events e
		JOIN institutions i ON e.institution_id = i.id
		JOIN users u ON e.creator_id = u.id
		WHERE e.is_deleted = false AND e.status = 'approved'
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
			&e.CreatorID, &e.CreatorName, &e.ParticipantsCount, &e.IsJoined, &e.Status, &e.CreatedAt,
			&total,
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

// JoinEvent registers a user for an event.
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

// LeaveEvent unregisters a user from an event.
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

// GetInstitutionEvents returns all events proposed by a specific institution for moderation.
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
			COALESCE(u.full_name, u.phone, 'Organizer') as creator_name,
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

// UpdateEventStatus updates the event status.
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
