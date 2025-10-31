package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"shb/pkg/constants"
)

func (h *Handler) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.Request.Header.Get(constants.RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, constants.RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
