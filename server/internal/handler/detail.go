package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

type (
	DetailResponse struct {
		ID           uuid.UUID                      `json:"stamp_id"`
		Name         string                         `json:"stamp_name"`
		FileID       uuid.UUID                      `json:"file_id"`
		CreatorID    uuid.UUID                      `json:"creator_id"`
		IsUnicode    bool                           `json:"is_unicode"`
		CreatedAt    time.Time                      `json:"created_at"`
		UpdatedAt    time.Time                      `json:"updated_at"`
		CountMonthly int                            `json:"count_monthly"`
		CountTotal   int64                          `json:"count_total"`
		Descriptions []*repository.StampDescription `json:"descriptions"`
		Tags         []*repository.TagSummary       `json:"tags"`
	}
)

func (h *Handler) getDetails(c echo.Context) error {
	stampIDStr := c.Param("stampId")
	stampID, err := uuid.Parse(stampIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid stamp ID format")
	}

	stamps, err := h.repo.GetStampByStampID(c.Request().Context(), stampID)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Stamp not found")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get stamp details").SetInternal(err)
	}

	descriptions, err := h.repo.GetDescriptionsByStampID(c.Request().Context(), stampID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	tags, err := h.repo.GetTagsByStampID(c.Request().Context(), stampID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	res := DetailResponse{
		ID:           stamps.ID,
		Name:         stamps.Name,
		FileID:       stamps.FileID,
		CreatorID:    stamps.CreatorID,
		IsUnicode:    stamps.IsUnicode,
		CreatedAt:    stamps.CreatedAt,
		UpdatedAt:    stamps.UpdatedAt,
		CountMonthly: stamps.CountMonthly,
		CountTotal:   stamps.CountTotal,
		Descriptions: descriptions,
		Tags:         tags,
	}

	return c.JSON(http.StatusOK, res)
}
