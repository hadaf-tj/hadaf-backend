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
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var input models.Need
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}
	input.ID = id

	// H5: Ownership check — employee can only edit needs of their own institution
	role, _ := c.Get("role")
	if role.(string) != models.RoleSuperAdmin {
		userID, _ := c.Get("userID")
		need, err := h.service.GetNeedByID(ctx, id)
		if err != nil {
			h.handleError(c, err)
			return
		}
		user, err := h.service.GetUserByID(ctx, userID.(int))
		if err != nil {
			h.handleError(c, err)
			return
		}
		if user.InstitutionID == nil || *user.InstitutionID != need.InstitutionID {
			h.handleError(c, myerrors.NewForbiddenErr("you can only edit needs of your own institution"))
			return
		}
	}

	if err := h.service.UpdateNeed(ctx, &input); err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, "updated")
}

func (h *Handler) deleteNeed(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	// H5: Ownership check — employee can only delete needs of their own institution
	role, _ := c.Get("role")
	if role.(string) != models.RoleSuperAdmin {
		userID, _ := c.Get("userID")
		need, err := h.service.GetNeedByID(ctx, id)
		if err != nil {
			h.handleError(c, err)
			return
		}
		user, err := h.service.GetUserByID(ctx, userID.(int))
		if err != nil {
			h.handleError(c, err)
			return
		}
		if user.InstitutionID == nil || *user.InstitutionID != need.InstitutionID {
			h.handleError(c, myerrors.NewForbiddenErr("you can only delete needs of your own institution"))
			return
		}
	}

	if err := h.service.DeleteNeed(ctx, id); err != nil {
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
	id, err := strconv.Atoi(idStr)

	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid institution id"))
		return
	}

	needs, err := h.service.GetNeedsByInstitution(ctx, filter, id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, needs)
}
