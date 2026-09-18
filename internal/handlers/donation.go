// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// getMyDonations returns the authenticated donor's recorded donation history.
// The donor identity is derived solely from auth middleware context.
func (h *Handler) getMyDonations(c *gin.Context) {
	ctx := c.Request.Context()

	userID, shouldReturn := h.mustGetUserID(c)
	if shouldReturn {
		return
	}

	limit, offset, err := parseLimitOffset(c)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log := zerolog.Ctx(ctx).With().
		Str("handler", "getMyDonations").
		Int("user_id", userID).
		Int("limit", limit).
		Int("offset", offset).
		Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	page, err := h.service.GetMyDonations(ctx, userID, limit, offset)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int("count", len(page.Items)).Msg("donation history returned")
	h.success(c, page)
}
