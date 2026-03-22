package handlers

import (
	"shb/internal/models"
	"shb/pkg/myerrors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllInstitutions(c *gin.Context) {
	ctx := c.Request.Context()
	
	// Читаем параметры
	search := c.Query("search") // Поиск по имени или городу
	iType := c.Query("type")    // Тип
	sortBy := c.Query("sort")   // 'needs_desc' или 'distance'
	
	// Координаты пользователя (если он разрешил геолокацию)
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
		h.logger.Error().Err(err).Msg("failed to get institutions")
		h.handleError(c, myerrors.ErrGeneral)
		return
	}

	h.success(c, page)
}

func (h *Handler) createInstitution(c *gin.Context) {
	var input models.Institution
	if err := c.ShouldBindJSON(&input); err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid input"))
		return
	}

	// Валидацию можно добавить здесь или в сервисе

	id, err := h.service.CreateInstitution(c.Request.Context(), &input)
	if err != nil {
		h.logger.Error().Err(err).Msg("failed to create institution")
		h.handleError(c, myerrors.ErrGeneral)
		return
	}

	h.success(c, gin.H{"id": id})
}

func (h *Handler) getInstitutionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid id"))
		return
	}

	institution, err := h.service.GetInstitutionByID(c.Request.Context(), id)
	if err != nil {
		// Здесь можно добавить проверку на sql.ErrNoRows и возвращать 404
		h.logger.Error().Err(err).Int("id", id).Msg("failed to get institution")
		h.handleError(c, myerrors.ErrGeneral)
		return
	}

	h.success(c, institution)
}
