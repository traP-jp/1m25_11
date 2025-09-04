package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

type descriptionPayload struct {
	Description string `json:"description"`
}

func (h *Handler) createDescriptions(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stampId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "スタンプIDの形式が不正です。", Code: "invalid_stamp_id"})
	}

	creatorID := uuid.Nil 

	payload := new(descriptionPayload)
	if err = c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "リクエストボディが不正です。", Code: "invalid_request_body"})
	}
	if payload.Description == "" {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "説明文は空にできません。", Code: "empty_description"})
	}

	err = h.repo.CreateDescriptions(c.Request().Context(), repository.CreateDescriptionParams{
		StampID:     stampID,
		Description: payload.Description,
		CreatorID:   creatorID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{Message: "指定されたスタンプが見つかりません。", Code: "stamp_not_found"})
		}
		if errors.Is(err, repository.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, &ErrorResponse{Message: "既にこのスタンプへの説明文を投稿済みです。", Code: "description_already_exists"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "説明文の作成中にエラーが発生しました。", Code: "internal_server_error"})
	}

	return c.NoContent(http.StatusCreated)
}

func (h *Handler) getDescriptions(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stampId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "スタンプIDの形式が不正です。", Code: "invalid_stamp_id"})
	}

	descriptions, err := h.repo.GetDescriptionsByStampID(c.Request().Context(), stampID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{Message: "指定されたスタンプが見つかりません。", Code: "stamp_not_found"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "説明文の取得中にエラーが発生しました。", Code: "internal_server_error"})
	}

	return c.JSON(http.StatusOK, descriptions)
}

func (h *Handler) updateDescriptions(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stampId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "スタンプIDの形式が不正です。", Code: "invalid_stamp_id"})
	}

	creatorID := uuid.Nil 

	payload := new(descriptionPayload)
	if err = c.Bind(payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "リクエストボディが不正です。", Code: "invalid_request_body"})
	}
	if payload.Description == "" {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "説明文は空にできません。", Code: "empty_description"})
	}

	if err = h.repo.UpdateDescriptions(c.Request().Context(), stampID, creatorID, payload.Description); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{Message: "編集対象の説明文が見つかりません。", Code: "description_not_found"})
		}
		if errors.Is(err, repository.ErrForbidden) {
			return echo.NewHTTPError(http.StatusForbidden, &ErrorResponse{Message: "この説明文を編集する権限がありません。", Code: "forbidden"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "説明文の更新中にエラーが発生しました。", Code: "internal_server_error"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) deleteDescriptions(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stampId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "スタンプIDの形式が不正です。", Code: "invalid_stamp_id"})
	}

	creatorID := uuid.Nil 

	if err = h.repo.DeleteDescriptions(c.Request().Context(), stampID, creatorID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{Message: "削除対象の説明文が見つかりません。", Code: "description_not_found"})
		}
		if errors.Is(err, repository.ErrForbidden) {
			return echo.NewHTTPError(http.StatusForbidden, &ErrorResponse{Message: "この説明文を削除する権限がありません。", Code: "forbidden"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "説明文の削除中にエラーが発生しました。", Code: "internal_server_error"})
	}

	return c.NoContent(http.StatusNoContent)
}

