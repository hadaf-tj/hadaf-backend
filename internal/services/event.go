// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shb/internal/models"
	"shb/pkg/myerrors"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
)

// CreateEvent creates a new volunteer event. The event date must be in the
// future and the referenced institution must exist.
func (s *Service) CreateEvent(ctx context.Context, e *models.Event) (int, error) {
	log := zerolog.Ctx(ctx).With().Str("service", "CreateEvent").Int("creator_id", e.CreatorID).Logger()

	if e.EventDate.Before(time.Now()) {
		return 0, myerrors.NewBadRequestErr("event date must be in the future")
	}

	if _, err := s.repo.GetInstitutionByID(ctx, e.InstitutionID); err != nil {
		return 0, myerrors.NewBadRequestErr("institution not found")
	}

	id, err := s.repo.CreateEvent(ctx, e)
	if err != nil {
		return 0, fmt.Errorf("create event: %w", err)
	}

	log.Info().Int("event_id", id).Int("institution_id", e.InstitutionID).Msg("event created")
	return id, nil
}

// GetAllEvents retrieves a paginated list of events.
func (s *Service) GetAllEvents(ctx context.Context, q models.EventListQuery) (*models.EventPage, error) {
	return s.repo.GetAllEvents(ctx, q)
}

// GetEventByID retrieves a single event by its primary key.
func (s *Service) GetEventByID(ctx context.Context, id int) (*models.Event, error) {
	return s.repo.GetEventByID(ctx, id)
}

// GetEventDetail returns the full event card including join status for the
// requesting user.
func (s *Service) GetEventDetail(ctx context.Context, q models.EventDetailQuery) (*models.EventResponse, error) {
	ev, err := s.repo.GetEventDetail(ctx, q)
	if err != nil {
		if errors.Is(err, myerrors.ErrNotFound) || errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get event detail: %w", myerrors.ErrNotFound)
		}
		return nil, err
	}
	return ev, nil
}

// JoinEvent registers a user for an upcoming event. Past events cannot be
// joined.
func (s *Service) JoinEvent(ctx context.Context, eventID, userID int) error {
	log := zerolog.Ctx(ctx).With().Str("service", "JoinEvent").Int("event_id", eventID).Int("user_id", userID).Logger()

	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return myerrors.NewBadRequestErr("event not found")
	}

	if event.EventDate.Before(time.Now()) {
		return myerrors.NewBadRequestErr("cannot join past event")
	}

	if err := s.repo.JoinEvent(ctx, eventID, userID); err != nil {
		return err
	}

	log.Info().Msg("user joined event")
	return nil
}

// LeaveEvent cancels a user's registration for an event.
func (s *Service) LeaveEvent(ctx context.Context, eventID, userID int) error {
	log := zerolog.Ctx(ctx).With().Str("service", "LeaveEvent").Int("event_id", eventID).Int("user_id", userID).Logger()

	if _, err := s.repo.GetEventByID(ctx, eventID); err != nil {
		return myerrors.NewBadRequestErr("event not found")
	}

	if err := s.repo.LeaveEvent(ctx, eventID, userID); err != nil {
		return err
	}

	log.Info().Msg("user left event")
	return nil
}

// GetInstitutionEvents returns all events associated with a given institution
// (used by institution staff for moderation).
func (s *Service) GetInstitutionEvents(ctx context.Context, institutionID int) ([]*models.EventResponse, error) {
	return s.repo.GetInstitutionEvents(ctx, institutionID)
}

// ApproveEvent transitions a pending event to the approved state.
func (s *Service) ApproveEvent(ctx context.Context, eventID int) error {
	log := zerolog.Ctx(ctx).With().Str("service", "ApproveEvent").Int("event_id", eventID).Logger()
	if err := s.repo.UpdateEventStatus(ctx, eventID, "approved"); err != nil {
		return err
	}
	log.Info().Msg("event approved")
	return nil
}

// RejectEvent transitions a pending event to the rejected state.
func (s *Service) RejectEvent(ctx context.Context, eventID int) error {
	log := zerolog.Ctx(ctx).With().Str("service", "RejectEvent").Int("event_id", eventID).Logger()
	if err := s.repo.UpdateEventStatus(ctx, eventID, "rejected"); err != nil {
		return err
	}
	log.Info().Msg("event rejected")
	return nil
}
