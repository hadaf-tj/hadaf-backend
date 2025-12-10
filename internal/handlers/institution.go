package handlers

import (
	"shb/internal/models"
	"shb/pkg/myerrors"

	"github.com/gin-gonic/gin"
)

// getAllInstitutions возвращает список учреждений
// @Summary Get all institutions
// @Description Возвращает список всех учреждений с возможностью фильтрации по городу
// @Tags Institutions
// @Accept json
// @Produce json
// @Param city query string false "City filter"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.ErrorResponse
// @Router /institutions [get]
func (h *Handler) getAllInstitutions(c *gin.Context) {
	ctx := c.Request.Context()
	city := c.Query("city")

	institutions, err := h.service.GetAllInstitutions(ctx, city)
	if err != nil {
		h.logger.Error().Err(err).Str("city", city).Msg("failed to get institutions")
		h.handleError(c, myerrors.ErrGeneral)
		return
	}

	h.success(c, institutions)
}

// createInstitution создает новое учреждение
// @Summary Create institution
// @Description Создает новое учреждение (для админов)
// @Tags Institutions
// @Accept json
// @Produce json
// @Param input body models.Institution true "Institution data"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /institutions [post]
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