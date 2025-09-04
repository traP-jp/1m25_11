package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type rankingResponse struct {
	StampID      uuid.UUID `json:"stamp_id"`
	MonthlyCount int       `json:"month_count"`
	TotalCount   int       `json:"total_count"`
}

func (h *Handler) getRanking(c echo.Context) error {

	rankingResults, err := h.repo.GetRanking(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.JSON(http.StatusOK, rankingResults)
}
