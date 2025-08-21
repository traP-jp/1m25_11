package handler

import (
	"net/http"

	vd "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

type DescriptionParams struct {
	StampID     uuid.UUID `param:"stamp_id"`
	CreatorID   uuid.UUID `json:"creator_id"`
	Description string    `json:"description"`
}

func (p *DescriptionParams) Validate(requireCreatorID, requireDescription bool) error {
	fields := []*vd.FieldRules{
		vd.Field(p.StampID, vd.Required),
	}
	if requireCreatorID {
		fields = append(fields, vd.Field(&p.CreatorID, vd.Required))
	}
	if requireDescription {
		fields = append(fields, vd.Field(p.Description, vd.Required))
	}

	return vd.ValidateStruct(p, fields...)
}

func (h *Handler) createDescriptions(c echo.Context) error {
	p := new(DescriptionParams)
	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if err := p.Validate(true, true); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if err := h.repo.CreateDescriptions(c.Request().Context(), repository.CreateDescriptionParams{
		StampID:     p.StampID,
		Description: p.Description,
		CreatorID:   p.CreatorID,
	}); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) getDescriptions(c echo.Context) error {
	p := new(DescriptionParams)
	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if err := p.Validate(false, false); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	descriptions, err := h.repo.GetDescriptionsByStampID(c.Request().Context(), p.StampID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.JSON(http.StatusOK, descriptions)
}

func (h *Handler) updateDescriptions(c echo.Context) error {
	p := new(DescriptionParams)
	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if err := p.Validate(true, true); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if err := h.repo.UpdateDescriptions(c.Request().Context(), p.StampID, p.CreatorID, p.Description); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) deleteDescriptions(c echo.Context) error {
	p := new(DescriptionParams)
	if err := c.Bind(p); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if err := p.Validate(true, false); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}
	if err := h.repo.DeleteDescriptions(c.Request().Context(), p.StampID, p.CreatorID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}
