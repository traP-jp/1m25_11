package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

const userIDContextKey = "userID"

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := h.getUserID(c)
		if err != nil {
			return err
		}

		c.Set(userIDContextKey, id)

		return next(c)
	}
}

func (h *Handler) ProxySecretMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		secret := os.Getenv("PROXY_SECRET")
		if secret == "" {
			if os.Getenv("APP_ENV") != "development" {
				return echo.NewHTTPError(http.StatusForbidden, "forbidden")
			}

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

		auth := c.Request().Header.Get("Authorization")
		scheme, token, ok := strings.Cut(auth, " ")
		if !ok || !strings.EqualFold(scheme, "bearer") || token != botToken {
			return echo.NewHTTPError(http.StatusUnauthorized, "unauthorized")
		}

		return next(c)
	}
}
