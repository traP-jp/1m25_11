package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"
	"github.com/traP-jp/1m25_11/server/api"
)

func (h *Handler) getRanking(c echo.Context) error {
	var params api.GetStampsRankingParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid parameters: %w", err))
	}
	since := time.Now().AddDate(0, 0, -30)
	if params.Since != nil {
		since = params.Since.Time
	}

	until := time.Now().AddDate(0, 0, -1)
	if params.Until != nil {
		until = params.Until.Time
	}

	if since.After(until) {
		return echo.NewHTTPError(http.StatusBadRequest, api.Error{
			Message: "Invalid date range: `since` must be before or on `until`.",
		})
	}

	rankingResults, err := h.repo.GetStampCount(c.Request().Context(), since, until)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.Error{
			Message: fmt.Sprintf("Failed to get stamp ranking: %s", err.Error()),
		})
	}

	var response []api.RankingResult
	for _, result := range rankingResults {
		response = append(response, api.RankingResult{
			Stamp: api.Stamp{
				Id:     result.StampID,
				Name:   "",
				FileId: types.UUID{},
			},
			Count: result.ReactionCount,
		})
	}

	return c.JSON(http.StatusOK, response)
}
