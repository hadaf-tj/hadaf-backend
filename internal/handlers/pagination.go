package handlers

import (
	"strconv"

	"shb/pkg/constants"
	"shb/pkg/myerrors"

	"github.com/gin-gonic/gin"
)

func parseLimitOffset(c *gin.Context) (limit, offset int, err error) {
	offsetStr := c.Query("offset")
	if offsetStr == "" {
		offset = 0
	} else {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return 0, 0, myerrors.NewBadRequestErr("invalid offset")
		}
	}

	limitStr := c.Query("limit")
	if limitStr == "" {
		limit = constants.DefaultPageLimit
		return limit, offset, nil
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return 0, 0, myerrors.NewBadRequestErr("invalid limit")
	}
	if limit > constants.MaxPageLimit {
		return 0, 0, myerrors.NewBadRequestErr("limit exceeds maximum")
	}

	return limit, offset, nil
}
