package handlers

import (
	"shb/internal/models"
	"shb/pkg/myerrors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateEventInput - структура для входных данных создания события
type CreateEventInput struct {
	Title         string `json:"title" binding:"required"`
	Description   string `json:"description"`
	EventDate     string `json:"event_date" binding:"required"` // ISO format
	InstitutionID int    `json:"institution_id" binding:"required"`
}

// getAllEvents возвращает список всех событий
func (h *Handler) getAllEvents(c *gin.Context) {
	ctx := c.Request.Context()

	// Получаем ID текущего пользователя (если авторизован)
	userID := 0
	if id, exists := c.Get("userID"); exists {
		if uid, ok := id.(int); ok {
			userID = uid
		}
	}

	events, err := h.service.GetAllEvents(ctx, userID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, events)
}

// createEvent создаёт новое событие
func (h *Handler) createEvent(c *gin.Context) {
	ctx := c.Request.Context()

	// Получаем ID создателя из context (установлено middleware)
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

	var input CreateEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}

	// Парсим дату
	eventDate, err := time.Parse(time.RFC3339, input.EventDate)
	if err != nil {
		// Пробуем альтернативный формат
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

	h.success(c, gin.H{"id": id})
}

// joinEvent записывает пользователя на событие
func (h *Handler) joinEvent(c *gin.Context) {
	ctx := c.Request.Context()

	// Получаем ID пользователя
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

	// Получаем ID события
	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid event id"))
		return
	}

	if err := h.service.JoinEvent(ctx, eventID, userID); err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, "joined")
}

// leaveEvent отменяет запись пользователя на событие
func (h *Handler) leaveEvent(c *gin.Context) {
	ctx := c.Request.Context()

	// Получаем ID пользователя
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

	// Получаем ID события
	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid event id"))
		return
	}

	if err := h.service.LeaveEvent(ctx, eventID, userID); err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, "left")
}

// getInstitutionEvents возвращает список событий для учреждения (модерация)
func (h *Handler) getInstitutionEvents(c *gin.Context) {
	ctx := c.Request.Context()
	instIDStr := c.Param("id")
	institutionID, err := strconv.Atoi(instIDStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid institution id"))
		return
	}
	events, err := h.service.GetInstitutionEvents(ctx, institutionID)
	if err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, events)
}

// approveEvent одобряет событие
func (h *Handler) approveEvent(c *gin.Context) {
	ctx := c.Request.Context()
	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid event id"))
		return
	}
	if err := h.service.ApproveEvent(ctx, eventID); err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, "approved")
}

// rejectEvent отклоняет событие
func (h *Handler) rejectEvent(c *gin.Context) {
	ctx := c.Request.Context()
	eventIDStr := c.Param("id")
	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid event id"))
		return
	}
	if err := h.service.RejectEvent(ctx, eventID); err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, "rejected")
}

