// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package handlers

import (
	"shb/internal/models"
	"shb/pkg/myerrors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// CreateEventInput holds the request body for creating a new event.
type CreateEventInput struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	EventDate     string `json:"event_date" binding:"required"` // ISO format
	InstitutionID int    `json:"institution_id" binding:"required"`
}

// getAllEvents returns a paginated list of all events.
func (h *Handler) getAllEvents(c *gin.Context) {
	ctx := c.Request.Context()

	userID := 0
	if id, exists := c.Get("userID"); exists {
		if uid, ok := id.(int); ok {
			userID = uid
		}
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "getAllEvents").Int("user_id", userID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	limit, offset, err := parseLimitOffset(c)
	if err != nil {
		h.handleError(c, err)
		return
	}

	page, err := h.service.GetAllEvents(ctx, models.EventListQuery{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int64("total", page.Total).Msg("events fetched")
	h.success(c, page)
}

func (h *Handler) getEventByID(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid id"))
		return
	}

	userID := 0
	if id, exists := c.Get("userID"); exists {
		if uid, ok := id.(int); ok {
			userID = uid
		}
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "getEventByID").Int("event_id", eventID).Int("user_id", userID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	ev, err := h.service.GetEventDetail(ctx, models.EventDetailQuery{
		EventID:      eventID,
		ViewerUserID: userID,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Msg("event fetched")
	h.success(c, ev)
}

// createEvent creates a new volunteer event.
func (h *Handler) createEvent(c *gin.Context) {
	ctx := c.Request.Context()

	creatorIDRaw, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("unauthorized"))
		return
	}
	creatorID, ok := creatorIDRaw.(int)
	if !ok || creatorID == 0 {
		h.handleError(c, myerrors.NewUnauthorizedErr("unauthorized"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "createEvent").Int("user_id", creatorID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	var input CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}

	// Parse the event date from ISO 8601 format.
	eventDate, err := time.Parse(time.RFC3339, input.EventDate)
	if err != nil {
		eventDate, err = time.Parse("2006-01-02T15:04:05", input.EventDate)
		if err != nil {
			h.handleError(c, myerrors.NewBadRequestErr("invalid date format, use ISO 8601"))
			return
		}
	}

	event := &models.Event{
		Title:         input.Title,
		Description:   input.Description,
		EventDate:     eventDate,
		InstitutionID: input.InstitutionID,
		CreatorID:     creatorID,
		Status:        "pending",
	}

	id, err := h.service.CreateEvent(ctx, event)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int("event_id", id).Msg("event created")
	h.success(c, gin.H{"id": id})
}

// joinEvent registers the authenticated user for an event.
func (h *Handler) joinEvent(c *gin.Context) {
	ctx := c.Request.Context()

	userIDRaw, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("unauthorized"))
		return
	}
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		h.handleError(c, myerrors.NewUnauthorizedErr("unauthorized"))
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid event id"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "joinEvent").Int("user_id", userID).Int("event_id", eventID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	if err := h.service.JoinEvent(ctx, eventID, userID); err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Msg("user joined event")
	h.success(c, "joined")
}

// leaveEvent cancels the authenticated user's registration for an event.
func (h *Handler) leaveEvent(c *gin.Context) {
	ctx := c.Request.Context()

	userIDRaw, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("unauthorized"))
		return
	}
	userID, ok := userIDRaw.(int)
	if !ok || userID == 0 {
		h.handleError(c, myerrors.NewUnauthorizedErr("unauthorized"))
		return
	}

	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid event id"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "leaveEvent").Int("user_id", userID).Int("event_id", eventID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	if err := h.service.LeaveEvent(ctx, eventID, userID); err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Msg("user left event")
	h.success(c, "left")
}

// getInstitutionEvents returns all events associated with a given institution.
func (h *Handler) getInstitutionEvents(c *gin.Context) {
	ctx := c.Request.Context()

	instIDStr := c.Param("id")
	institutionID, err := strconv.Atoi(instIDStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid institution id"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "getInstitutionEvents").Int("institution_id", institutionID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	events, err := h.service.GetInstitutionEvents(ctx, institutionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int("count", len(events)).Msg("institution events fetched")
	h.success(c, events)
}

// approveEvent transitions a pending event to the approved state.
func (h *Handler) approveEvent(c *gin.Context) {
	ctx := c.Request.Context()

	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid event id"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "approveEvent").Int("event_id", eventID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	if err := h.service.ApproveEvent(ctx, eventID); err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Msg("event approved")
	h.success(c, "approved")
}

// rejectEvent transitions a pending event to the rejected state.
func (h *Handler) rejectEvent(c *gin.Context) {
	ctx := c.Request.Context()

	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid event id"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "rejectEvent").Int("event_id", eventID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	if err := h.service.RejectEvent(ctx, eventID); err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Msg("event rejected")
	h.success(c, "rejected")
}
