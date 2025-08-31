package handler

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type getRankingParams struct {
	Since *time.Time `query:"since"`
	Until *time.Time `query:"until"`
}

type errorResponse struct {
	Message string `json:"message"`
}

type rankingResponse struct {
	StampID       uuid.UUID `json:"stamp_id"`
	BodyCount     int       `json:"body_count"`
	ReactionCount int       `json:"reaction_count"`
}

func (h *Handler) getRanking(c echo.Context) error {
	var params getRankingParams
	if err := c.Bind(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	since := time.Now().AddDate(0, 0, -30)
	if params.Since != nil {
		since = *params.Since
	}

	until := time.Now().AddDate(0, 0, -1)
	if params.Until != nil {
		until = *params.Until
	}

	if since.After(until) {
		return echo.NewHTTPError(http.StatusBadRequest, errorResponse{
			Message: "Invalid date range: `since` must be before or on `until`.",
		})
	}

	rankingResults, err := h.repo.GetRanking(c.Request().Context(), since, until)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	response := make([]rankingResponse, 0, len(rankingResults))
	for _, result := range rankingResults {
		response = append(response, rankingResponse{
			StampID:       result.StampID,
			BodyCount:     result.MessageCount,
			ReactionCount: result.ReactionCount,
		})
	}

	return c.JSON(http.StatusOK, response)
}

