// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2026 Siyovush Hamidov and The Hadaf Contributors

package handlers

import "github.com/gin-gonic/gin"

func (h *Handler) getApplicationQuota(c *gin.Context) {
	quota, err := h.service.GetApplicationQuota(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, quota)
}
