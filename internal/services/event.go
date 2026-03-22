package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"shb/internal/models"
	"shb/pkg/myerrors"
)

// CreateEvent создаёт новое событие
func (s *Service) CreateEvent(ctx context.Context, e *models.Event) (int, error) {
	// Валидация: дата события должна быть в будущем
	if e.EventDate.Before(time.Now()) {
		return 0, myerrors.NewBadRequestErr("event date must be in the future")
	}

	// Проверяем, что учреждение существует
	_, err := s.repo.GetInstitutionByID(ctx, e.InstitutionID)
	if err != nil {
		return 0, myerrors.NewBadRequestErr("institution not found")
	}

	return s.repo.CreateEvent(ctx, e)
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
		if errors.Is(err, myerrors.ErrNotFound) {
			return nil, fmt.Errorf("get event detail: %w", myerrors.ErrNotFound)
		}
		return nil, err
	}
	return ev, nil
}

// JoinEvent записывает пользователя на событие
func (s *Service) JoinEvent(ctx context.Context, eventID, userID int) error {
	// Проверяем, что событие существует и не прошло
	event, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return myerrors.NewBadRequestErr("event not found")
	}

	// Нельзя записаться на прошедшее событие
	if event.EventDate.Before(time.Now()) {
		return myerrors.NewBadRequestErr("cannot join past event")
	}

	return s.repo.JoinEvent(ctx, eventID, userID)
}

// LeaveEvent отменяет запись на событие
func (s *Service) LeaveEvent(ctx context.Context, eventID, userID int) error {
	// Проверяем, что событие существует
	_, err := s.repo.GetEventByID(ctx, eventID)
	if err != nil {
		return myerrors.NewBadRequestErr("event not found")
	}

	return s.repo.LeaveEvent(ctx, eventID, userID)
}
