package handler

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if _, err := h.getUserID(c); err != nil {
			return err
		}

		return next(c)
	}
}

func (h *Handler) ProxySecretMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		secret := os.Getenv("PROXY_SECRET")
		if secret == "" {
			return next(c)
		}

		if c.Request().Header.Get("X-Proxy-Secret") != secret {
			return echo.NewHTTPError(http.StatusForbidden, "forbidden")
		}

		return next(c)
	}
}

func (h *Handler) BulkAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		botToken := os.Getenv("BOT_TOKEN_KEY")
		if botToken == "" {
			return next(c)
		}

		if c.Request().Header.Get("Authorization") != "Bearer "+botToken {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		return next(c)
	}
}
