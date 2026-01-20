package handlers

import (
	"shb/pkg/myerrors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createBooking(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract user ID from JWT token (set by AuthMiddleware)
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

	h.success(c, gin.H{"id": bookingID})
}

func (h *Handler) approveBooking(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract user ID from JWT token
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

	if err := h.service.ApproveBooking(ctx, bookingID, userIDInt); err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, "booking approved")
}

func (h *Handler) rejectBooking(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract user ID from JWT token
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

	if err := h.service.RejectBooking(ctx, bookingID, userIDInt); err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, "booking rejected")
}

func (h *Handler) getInstitutionBookings(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	institutionID, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid institution ID"))
		return
	}

	bookings, err := h.service.GetBookingsByInstitution(ctx, institutionID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, bookings)
}

func (h *Handler) getMyBookings(c *gin.Context) {
	ctx := c.Request.Context()

	// Extract user ID from JWT token
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

	bookings, err := h.service.GetBookingsByUser(ctx, userIDInt)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, bookings)
}
