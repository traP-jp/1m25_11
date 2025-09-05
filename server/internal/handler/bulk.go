package handler

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

func (h *Handler) BulkCreateTags(c echo.Context) error {
	var req []repository.TagInfo
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	createdTags, err := h.repo.BulkCreateTags(c.Request().Context(), req)
	if err != nil {
		log.Printf("error in BulkCreateTags repository call: %v", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create tags")
	}

	return c.JSON(http.StatusCreated, createdTags)
}

func (h *Handler) BulkAddStampMeta(c echo.Context) error {
	var req []repository.StampMetaAddition
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	err := h.repo.BulkAddStampMeta(c.Request().Context(), req)
	if err != nil {
		log.Printf("error in BulkAddStampMeta repository call: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to link tags and add descriptions")
	}

	return c.NoContent(http.StatusNoContent)
}
