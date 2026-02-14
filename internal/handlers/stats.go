package handlers

import (
	"github.com/gin-gonic/gin"
)

// getStats возвращает публичную статистику для лендинга
func (h *Handler) getStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.service.GetPublicStats(ctx)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, stats)
}
