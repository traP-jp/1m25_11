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
		since = time.Time(params.Since)
	}

	until := time.Now().AddDate(0, 0, -1)
	if params.Until != nil {
		until = time.Time(params.Until)
	}

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	limit := 50
	if params.Limit != nil {
		limit = *params.Limit
	}

	if since.After(until) {
		return echo.NewHTTPError(http.StatusBadRequest, api.Error{
			Message: "Invalid date range: `since` must be before or on `until`.",
		})
	}
	rankingResults, err := h.Repo.GetStampCount(c.Request().Context(), since, until, &limit, &offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, api.Error{
			Message: fmt.Sprintf("Failed to get stamp ranking: %s", err.Error()),
		})
	}

	var response []api.RankingResult
	for _, result := range rankingResults {
		stampSummary, err := h.Repo.GetStampSummaryByID(c.Request().Context(), types.UUID(result.StampID))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, api.Error{
				Message: fmt.Sprintf("Failed to get stamp summary for ID %s: %s", result.StampID, err.Error()),
			})
		}

		response = append(response, api.RankingResult{
			Stamp: api.StampSummary{
				Id:     stampSummary.Id,
				Name:   stampSummary.Name,
				FileId: stampSummary.FileId,
			},
			Count: result.Count,
		})
	}

	return c.String(http.StatusOK, "pong")
}
