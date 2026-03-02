package handlers

import (
	"github.com/gin-gonic/gin"
)

// getSMSBalance возвращает баланс SMS аккаунта
func (h *Handler) getSMSBalance(c *gin.Context) {
	ctx := c.Request.Context()

	balance, err := h.service.CheckSMSBalance(ctx)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.success(c, balance)
}
