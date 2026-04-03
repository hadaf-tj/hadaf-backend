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

// CreateEvent создаёт новое событие
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

// GetAllEvents получает страницу событий.
func (s *Service) GetAllEvents(ctx context.Context, q models.EventListQuery) (*models.EventPage, error) {
	return s.repo.GetAllEvents(ctx, q)
}

// GetEventByID получает событие по ID
func (s *Service) GetEventByID(ctx context.Context, id int) (*models.Event, error) {
	return s.repo.GetEventByID(ctx, id)
}

// GetEventDetail возвращает карточку события (как в списке).
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

// JoinEvent записывает пользователя на событие
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

// LeaveEvent отменяет запись пользователя на событие
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

// GetInstitutionEvents получает события для конкретного учреждения (для модерации)
func (s *Service) GetInstitutionEvents(ctx context.Context, institutionID int) ([]*models.EventResponse, error) {
	return s.repo.GetInstitutionEvents(ctx, institutionID)
}

// ApproveEvent одобряет предложенное событие
func (s *Service) ApproveEvent(ctx context.Context, eventID int) error {
	log := zerolog.Ctx(ctx).With().Str("service", "ApproveEvent").Int("event_id", eventID).Logger()
	if err := s.repo.UpdateEventStatus(ctx, eventID, "approved"); err != nil {
		return err
	}
	log.Info().Msg("event approved")
	return nil
}

// RejectEvent отклоняет предложенное событие
func (s *Service) RejectEvent(ctx context.Context, eventID int) error {
	log := zerolog.Ctx(ctx).With().Str("service", "RejectEvent").Int("event_id", eventID).Logger()
	if err := s.repo.UpdateEventStatus(ctx, eventID, "rejected"); err != nil {
		return err
	}
	log.Info().Msg("event rejected")
	return nil
}
