package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	myerrors "shb/pkg/myerrors"
)

// getAllTeamMembers возвращает список всех участников команды
func (h *Handler) getAllTeamMembers(c *gin.Context) {
	ctx := c.Request.Context()

	members, err := h.service.GetAllTeamMembers(ctx)
	if err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, members)
}

// getTeamMemberByID возвращает подробную информацию о члене команды
func (h *Handler) getTeamMemberByID(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid id"))
		return
	}

	member, err := h.service.GetTeamMemberByID(ctx, id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, member)
}
