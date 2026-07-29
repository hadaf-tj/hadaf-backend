// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package handlers

import (
	"github.com/gin-gonic/gin"
)

// getSMSBalance returns the current SMS account balance.
func (h *Handler) getSMSBalance(c *gin.Context) {
	ctx := c.Request.Context()

	balance, err := h.service.CheckSMSBalance(ctx)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, balance)
}
