package handlers

import (
	"shb/internal/models"
	"shb/internal/repositories/filters"
	"shb/pkg/myerrors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createNeed(c *gin.Context) {
	var input models.Need
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}
	id, err := h.service.CreateNeed(c.Request.Context(), &input)
	if err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, gin.H{"id": id})
}

func (h *Handler) updateNeed(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var input models.Need
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}
	input.ID = id

	if err := h.service.UpdateNeed(c.Request.Context(), &input); err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, "updated")
}

func (h *Handler) deleteNeed(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := h.service.DeleteNeed(c.Request.Context(), id); err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, "deleted")
}

func (h *Handler) getNeedsByInstitution(c *gin.Context) {
	ctx := c.Request.Context()

	var filter filters.NeedsFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid need by institution input"))
		return
	}

	idStr := c.Param("id") // ID учреждения
	id, _ := strconv.Atoi(idStr)

	needs, err := h.service.GetNeedsByInstitution(ctx, filter, id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, needs)
}
