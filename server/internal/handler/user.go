package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)


func (h *Handler) GetUser(c echo.Context) error {
	creatorID, err := h.getUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized).SetInternal(err)
	}
	user, err := h.repo.GetUser(c.Request().Context(), creatorID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}


	return c.JSON(http.StatusOK, user)
}
