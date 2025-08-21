package handler

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func (h *Handler) getStamps(c echo.Context) error {
	stamps, err := h.repo.GetStampDetails(c.Request().Context())
	if err != nil{
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.JSON(http.StatusOK, stamps)
}


func (h *Handler) getCertainStamps(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}