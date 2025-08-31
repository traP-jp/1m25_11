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

type countResponse struct {
	Reaction int `json:"reaction"`
	Message  int `json:"message"`
}

type stampResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	FileID uuid.UUID `json:"fileId"`
}

type rankingResultResponse struct {
	Stamp stampResponse `json:"stamp"`
	Count countResponse `json:"count"`
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

	rankingResults, err := h.repo.GetStampCount(c.Request().Context(), since, until)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	var response []rankingResultResponse
	for _, result := range rankingResults {
		response = append(response, rankingResultResponse{
			Stamp: stampResponse{
				ID:     result.StampID,
				Name:   result.Name,
				FileID: result.FileID,
			},
			Count: countResponse{
				Reaction: result.ReactionCount,
				Message:  result.MessageCount,
			},
		})
	}

	return c.JSON(http.StatusOK, response)
}
