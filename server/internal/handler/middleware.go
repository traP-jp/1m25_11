package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// bulk エンドポイントは認証不要（Bot による一括インポート用）
		if strings.Contains(c.Request().URL.Path, "/bulk") {
			return next(c)
		}

		if _, err := h.getUserID(c); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		return next(c)
	}
}
