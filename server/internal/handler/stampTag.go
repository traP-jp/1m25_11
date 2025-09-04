package handler

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/internal/repository"
)

func (h *Handler) createStampTags(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stampId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "スタンプIDの形式が不正です。", Code: "invalid_stamp_id"})
	}
	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "タグIDの形式が不正です。", Code: "invalid_tag_id"})
	}

	creatorID := uuid.Nil

	err = h.repo.CreateStampTags(c.Request().Context(), repository.CreateStampTagParams{
		StampID:   stampID,
		TagID:     tagID,
		CreatorID: creatorID,
	})
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{Message: "指定されたスタンプまたはタグが見つかりません。", Code: "not_found"})
		}
		if errors.Is(err, repository.ErrAlreadyExists) {
			return echo.NewHTTPError(http.StatusConflict, &ErrorResponse{Message: "このタグは既に追加されています。", Code: "tag_already_added"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "タグ付け処理中にエラーが発生しました。", Code: "internal_server_error"})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *Handler) deleteStampTags(c echo.Context) error {
	stampID, err := uuid.Parse(c.Param("stampId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "スタンプIDの形式が不正です。", Code: "invalid_stamp_id"})
	}
	tagID, err := uuid.Parse(c.Param("tagId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, &ErrorResponse{Message: "タグIDの形式が不正です。", Code: "invalid_tag_id"})
	}

	err = h.repo.DeleteStampTags(c.Request().Context(), stampID, tagID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, &ErrorResponse{Message: "指定されたスタンプまたはタグの紐付けが見つかりません。", Code: "not_found"})
		}
		if errors.Is(err, repository.ErrForbidden) {
			return echo.NewHTTPError(http.StatusForbidden, &ErrorResponse{Message: "このタグを削除する権限がありません。", Code: "forbidden"})
		}
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, &ErrorResponse{Message: "タグの削除処理中にエラーが発生しました。", Code: "internal_server_error"})
	}

	return c.NoContent(http.StatusNoContent)
}
