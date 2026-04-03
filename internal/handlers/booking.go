package handlers

import (
	"fmt"
	"shb/internal/models"
	"shb/pkg/myerrors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func (h *Handler) createBooking(c *gin.Context) {
	ctx := c.Request.Context()

	role, _ := c.Get("role")
	if role.(string) == models.RoleEmployee {
		h.handleError(c, myerrors.NewForbiddenErr("institution employees cannot create bookings"))
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("user not authenticated"))
		return
	}
	userIDInt, ok := userID.(int)
	if !ok {
		h.handleError(c, myerrors.NewUnauthorizedErr("invalid user ID"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "createBooking").Int("user_id", userIDInt).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	limitKey := fmt.Sprintf("booking_create:%d", userIDInt)
	allowed, err := h.limiter.Allow(ctx, limitKey, 5, 3600)
	if err != nil {
		log.Error().Err(err).Msg("rate limiter error")
	} else if !allowed {
		h.handleError(c, myerrors.NewTooManyRequestsErr("Вы достигли лимита создания обещаний. Пожалуйста, попробуйте позже."))
		return
	}

	var input struct {
		NeedID   int     `json:"need_id" binding:"required"`
		Quantity float64 `json:"quantity" binding:"required,gt=0"`
		Note     string  `json:"note"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}

	bookingID, err := h.service.CreateBooking(ctx, userIDInt, input.NeedID, input.Quantity, input.Note)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int("booking_id", bookingID).Msg("booking created")
	h.success(c, gin.H{"id": bookingID})
}

func (h *Handler) approveBooking(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("user not authenticated"))
		return
	}
	userIDInt, ok := userID.(int)
	if !ok {
		h.handleError(c, myerrors.NewUnauthorizedErr("invalid user ID"))
		return
	}

	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid booking ID"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "approveBooking").Int("user_id", userIDInt).Int("booking_id", bookingID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	if err := h.service.ApproveBooking(ctx, bookingID, userIDInt); err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Msg("booking approved")
	h.success(c, "booking approved")
}

func (h *Handler) rejectBooking(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("user not authenticated"))
		return
	}
	userIDInt, ok := userID.(int)
	if !ok {
		h.handleError(c, myerrors.NewUnauthorizedErr("invalid user ID"))
		return
	}

	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid booking ID"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "rejectBooking").Int("user_id", userIDInt).Int("booking_id", bookingID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	if err := h.service.RejectBooking(ctx, bookingID, userIDInt); err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Msg("booking rejected")
	h.success(c, "booking rejected")
}

func (h *Handler) completeBooking(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("user not authenticated"))
		return
	}
	userIDInt, ok := userID.(int)
	if !ok {
		h.handleError(c, myerrors.NewUnauthorizedErr("invalid user ID"))
		return
	}

	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid booking ID"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "completeBooking").Int("user_id", userIDInt).Int("booking_id", bookingID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	if err := h.service.CompleteBooking(ctx, bookingID, userIDInt); err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Msg("booking completed")
	h.success(c, "booking completed")
}

func (h *Handler) getInstitutionBookings(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	institutionID, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid institution ID"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "getInstitutionBookings").Int("institution_id", institutionID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	role, _ := c.Get("role")
	if role.(string) != models.RoleSuperAdmin {
		userID, _ := c.Get("userID")
		user, err := h.service.GetUserByID(ctx, userID.(int))
		if err != nil {
			h.handleError(c, err)
			return
		}
		if user.InstitutionID == nil || *user.InstitutionID != institutionID {
			h.handleError(c, myerrors.NewForbiddenErr("you can only view bookings of your own institution"))
			return
		}
	}

	bookings, err := h.service.GetBookingsByInstitution(ctx, institutionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int("count", len(bookings)).Msg("institution bookings fetched")
	h.success(c, bookings)
}

func (h *Handler) getMyBookings(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("user not authenticated"))
		return
	}
	userIDInt, ok := userID.(int)
	if !ok {
		h.handleError(c, myerrors.NewUnauthorizedErr("invalid user ID"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "getMyBookings").Int("user_id", userIDInt).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	bookings, err := h.service.GetBookingsByUser(ctx, userIDInt)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int("count", len(bookings)).Msg("user bookings fetched")
	h.success(c, bookings)
}

func (h *Handler) cancelMyBooking(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("user not authenticated"))
		return
	}

	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid booking ID"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "cancelMyBooking").Int("user_id", userID.(int)).Int("booking_id", bookingID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	if err := h.service.CancelMyBooking(ctx, bookingID, userID.(int)); err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Msg("booking cancelled")
	h.success(c, "booking cancelled")
}

func (h *Handler) updateMyBooking(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("userID")
	if !exists {
		h.handleError(c, myerrors.NewUnauthorizedErr("user not authenticated"))
		return
	}

	idStr := c.Param("id")
	bookingID, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid booking ID"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "updateMyBooking").Int("user_id", userID.(int)).Int("booking_id", bookingID).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	var input struct {
		Quantity float64 `json:"quantity" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}

	if err := h.service.UpdateMyBooking(ctx, bookingID, userID.(int), input.Quantity); err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Float64("quantity", input.Quantity).Msg("booking updated")
	h.success(c, "booking updated")
}
