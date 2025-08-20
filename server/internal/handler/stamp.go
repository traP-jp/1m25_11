package handler

import (
	"net/http"
	"time"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)
type (
	stamp struct{
		ID         uuid.UUID `json:"id"`
		Name        string `json:"name"`
		FileID      uuid.UUID `json:"file_id"`
		CreatorID   uuid.UUID  `json:"creator_id"`
		IsUnicode   bool   `json:"is_unicode"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		CountMonthly int    `json:"count_monthly"`
		CountTotal   int64  `json:"count_total"`	
	}
)

func (h *Handler) getStamps(c echo.Context) error {
	stamps, err := h.repo.GetStamps(c.Request().Context())
	if err != nil{
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.JSON(http.StatusOK, stamps)
}


func (h *Handler) getCertainStamps(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}