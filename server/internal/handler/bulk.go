package handler

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)


type bulkCreateTagsRequest struct {
	Name string `json:"tag_name"`
}

type bulkCreateTagsResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (h *Handler) BulkCreateTags(c echo.Context) error {
	var req []bulkCreateTagsRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	tagsToCreate := make([]repository.TagInfo, len(req))
	for i, r := range req {
		tagsToCreate[i] = repository.TagInfo{Name: r.Name}
	}

	createdTags, err := h.repo.BulkCreateTags(c.Request().Context(), tagsToCreate)
	if err != nil {
		log.Printf("error in BulkCreateTags repository call: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create tags")
	}

	res := make([]bulkCreateTagsResponse, len(createdTags))
	for i, t := range createdTags {
		res[i] = bulkCreateTagsResponse{ID: t.ID, Name: t.Name}
	}

	return c.JSON(http.StatusCreated, res)
}


type bulkAddStampMetaRequest struct {
	ID          uuid.UUID   `json:"id"`
	TagIDs      []uuid.UUID `json:"tag_ids"`
	Description string      `json:"description"`
}

func (h *Handler) BulkAddStampMeta(c echo.Context) error {
	var req []bulkAddStampMetaRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	additions := make([]repository.StampMetaAddition, len(req))
	for i, r := range req {
		additions[i] = repository.StampMetaAddition{
			StampID:     r.ID,
			TagIDs:      r.TagIDs,
			Description: r.Description,
		}
	}

	err := h.repo.BulkAddStampMeta(c.Request().Context(), additions)
	if err != nil {
		log.Printf("error in BulkAddStampMeta repository call: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to link tags and add descriptions")
	}

	return c.NoContent(http.StatusNoContent)
}

