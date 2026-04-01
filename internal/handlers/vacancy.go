package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"

	myerrors "shb/pkg/myerrors"
)

// getAllVacancies возвращает список всех вакансий
func (h *Handler) getAllVacancies(c *gin.Context) {
	ctx := c.Request.Context()

	vacancies, err := h.service.GetAllVacancies(ctx)
	if err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, vacancies)
}

// getVacancyByID возвращает подробную информацию о вакансии
func (h *Handler) getVacancyByID(c *gin.Context) {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.handleError(c, myerrors.NewBadRequestErr("invalid id"))
		return
	}

	vacancy, err := h.service.GetVacancyByID(ctx, id)
	if err != nil {
		h.handleError(c, err)
		return
	}
	h.success(c, vacancy)
}
