// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package handlers

import (
	"shb/internal/models"
	"shb/internal/repositories/filters"
	"shb/pkg/myerrors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func (h *Handler) createNeed(c *gin.Context) {
	ctx := c.Request.Context()

	userID, _ := c.Get("userID")
	log := zerolog.Ctx(ctx).With().Str("handler", "createNeed").Int("user_id", userID.(int)).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	var input models.Need
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}

	// Security: always take institution_id from the JWT, not from request body.
	role, _ := c.Get("role")
	if role.(string) != models.RoleSuperAdmin {
		user, err := h.service.GetUserByID(ctx, userID.(int))
		if err != nil || user.InstitutionID == nil {
			h.handleError(c, myerrors.NewForbiddenErr("employee is not linked to any institution"))
			return
		}
		input.InstitutionID = *user.InstitutionID
	}

	id, err := h.service.CreateNeed(ctx, &input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int("need_id", id).Msg("need created")
	h.success(c, gin.H{"id": id})
}

func (h *Handler) updateNeed(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	userID, _ := c.Get("userID")

	log := zerolog.Ctx(ctx).With().Str("handler", "updateNeed").Int("user_id", userID.(int)).Int("need_id", id).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	var input models.Need
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}
	input.ID = id

	// H5: Ownership check — employee can only edit needs of their own institution
	role, _ := c.Get("role")
	if role.(string) != models.RoleSuperAdmin {
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

	log.Debug().Msg("need updated")
	h.success(c, "updated")
}

func (h *Handler) deleteNeed(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	userID, _ := c.Get("userID")

	log := zerolog.Ctx(ctx).With().Str("handler", "deleteNeed").Int("user_id", userID.(int)).Int("need_id", id).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	// H5: Ownership check — employee can only delete needs of their own institution
	role, _ := c.Get("role")
	if role.(string) != models.RoleSuperAdmin {
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

	log.Debug().Msg("need deleted")
	h.success(c, "deleted")
}

func (h *Handler) getNeedsByInstitution(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid institution id"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "getNeedsByInstitution").Int("institution_id", id).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	var filter filters.NeedsFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid need by institution input"))
		return
	}

	needs, err := h.service.GetNeedsByInstitution(ctx, filter, id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int("count", len(needs)).Msg("needs fetched")
	h.success(c, needs)
}
