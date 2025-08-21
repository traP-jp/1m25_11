package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) createDescriptions(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func (h *Handler) getDescriptions(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stamp_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	descriptions, err := h.repo.GetDescriptionsByStampID(c.Request().Context(), stampID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.JSON(http.StatusOK, descriptions)
}

func (h *Handler) updateDescriptions(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stamp_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	creatorID, err := uuid.Parse(c.Param("creator_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	err = h.repo.UpdateDescriptions(c.Request().Context(), stampID, creatorID, c.Param("description"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) deleteDescriptions(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}
