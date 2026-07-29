// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package handlers

import (
	"github.com/gin-gonic/gin"
)

// getStats returns aggregated public statistics for the landing page.
func (h *Handler) getStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.service.GetPublicStats(ctx)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, stats)
}
