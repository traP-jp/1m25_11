package handler

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
	"database/sql" 

)

type (
	Stamp struct {
		ID           uuid.UUID `json:"id"`
		Name         string    `json:"name"`
		FileID       uuid.UUID `json:"file_id"`
		IsUnicode    bool      `json:"is_unicode"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		CountMonthly int       `json:"count_monthly"`
		CountTotal   int64     `json:"count_total"`
	}
	Description struct {
		StampID     uuid.UUID `json:"stamp_id"`
		Description string    `json:"description"`
		CreatorID   uuid.UUID `json:"creator_id"`
		CreatedAt   time.Time `json:"created_at"`
	}
	TagSummary struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}
	DetailResponse struct {
		ID           uuid.UUID     `json:"id"`
		Name         string        `json:"name"`
		FileID       uuid.UUID     `json:"file_id"`
		IsUnicode    bool          `json:"is_unicode"`
		CreatedAt    time.Time        `json:"created_at"`
		UpdatedAt    time.Time        `json:"updated_at"`
		CountMonthly int           `json:"count_monthly"`
		CountTotal   int64         `json:"count_total"`
		Descriptions []Description `json:"descriptions"`
		Tags         []TagSummary  `json:"tags"`
	}
)

func (h *Handler) getDetails(c echo.Context) error {
	stampIDStr := c.Param("stampID") 
	stampID, err := uuid.Parse(stampIDStr)
    if err != nil {
        if err == sql.ErrNoRows {

            return echo.NewHTTPError(http.StatusNotFound, "Stamp not found")
        }

        return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
    }

	stamps, err := h.repo.GetStampByStampID(c.Request().Context(), stampID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	returnDescriptions, err := h.repo.GetDescriptionsByStampID(c.Request().Context(), stampID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	resDescriptions := make([]Description, len(returnDescriptions))
	for i, descriptionID := range returnDescriptions {
		resDescriptions[i] = Description{
			StampID:     descriptionID.StampID,
			Description: descriptionID.Description,
			CreatorID:   descriptionID.CreatorID,
			CreatedAt:   descriptionID.CreatedAt,
		}
	}

	returnTags, err := h.repo.GetTagsByStampID(c.Request().Context(), stampID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	resTags := make([]TagSummary, len(returnTags))
	for i, tagID := range returnTags {
		resTags[i] = TagSummary{
			ID:   tagID.ID,
			Name: tagID.Name,
		}
	}

	res := DetailResponse{
		ID:           stamps.ID,
		Name:         stamps.Name,
		FileID:       stamps.FileID,
		IsUnicode:    stamps.IsUnicode,
		CreatedAt:    stamps.CreatedAt,
		UpdatedAt:    stamps.UpdatedAt,
		CountMonthly: stamps.CountMonthly,
		CountTotal:   stamps.CountTotal,
		Descriptions: resDescriptions,
		Tags:         resTags,
	}

	return c.JSON(http.StatusOK, res)
}
