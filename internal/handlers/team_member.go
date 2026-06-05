// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	myerrors "shb/pkg/myerrors"
)

// getAllTeamMembers returns all team members.
func (h *Handler) getAllTeamMembers(c *gin.Context) {
	ctx := c.Request.Context()

	members, err := h.service.GetAllTeamMembers(ctx)
	if err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, members)
}

// getTeamMemberByID returns the details of a single team member.
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
