package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)
type (
	stamp struct{
		ID          string `json:"id"`
		Name        string `json:"name"`
		FileID      string `json:"file_id"`
		CreatorID   string `json:"creator_id"`
		IsUnicode   bool   `json:"is_unicode"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
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