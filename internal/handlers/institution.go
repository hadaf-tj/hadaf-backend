package handlers

import (
	"shb/internal/models"
	"shb/pkg/myerrors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func (h *Handler) getAllInstitutions(c *gin.Context) {
	ctx := c.Request.Context()

	search := c.Query("search")
	iType := c.Query("type")
	sortBy := c.Query("sort")
	latStr := c.Query("lat")
	lngStr := c.Query("lng")

	var lat, lng float64
	if latStr != "" && lngStr != "" {
		lat, _ = strconv.ParseFloat(latStr, 64)
		lng, _ = strconv.ParseFloat(lngStr, 64)
	}

	limit, offset, err := parseLimitOffset(c)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "getAllInstitutions").Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	page, err := h.service.GetAllInstitutions(ctx, models.InstitutionListQuery{
		Search:  search,
		Type:    iType,
		UserLat: lat,
		UserLng: lng,
		SortBy:  sortBy,
		Limit:   limit,
		Offset:  offset,
	})
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int64("total", page.Total).Msg("institutions fetched")
	h.success(c, page)
}

func (h *Handler) createInstitution(c *gin.Context) {
	ctx := c.Request.Context()

	log := zerolog.Ctx(ctx).With().Str("handler", "createInstitution").Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	var input models.Institution
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}

	id, err := h.service.CreateInstitution(ctx, &input)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Int("institution_id", id).Msg("institution created")
	h.success(c, gin.H{"id": id})
}

func (h *Handler) getInstitutionByID(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid id"))
		return
	}

	log := zerolog.Ctx(ctx).With().Str("handler", "getInstitutionByID").Int("institution_id", id).Logger()
	ctx = log.WithContext(ctx)
	c.Request = c.Request.WithContext(ctx)

	institution, err := h.service.GetInstitutionByID(ctx, id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	log.Debug().Msg("institution fetched")
	h.success(c, institution)
}
