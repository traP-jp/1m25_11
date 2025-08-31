package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler)AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	log.Print("AuthMiddleware called")
	return func(c echo.Context) error {
		if strings.HasSuffix(c.Request().URL.Path, "/login") {
			log.Print("Login request")

			return next(c)
		}
		if strings.Contains(c.Request().URL.Path, "/callback") {
			log.Print("Callback request")
			log.Print(c.Request().URL.Path)

			return next(c)
		}
		log.Printf("AuthMiddleware: %s", c.Request().URL.Path)

		token := getToken(c)
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized4")
		}

		if err := getMe(token); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: Invalid token")
		}

		c.Set("token", token)

		return next(c)
	}
}
func getToken(c echo.Context) string {
	cookieToken, err := c.Cookie(tokenKey)
	 if err == nil && cookieToken != nil {
        return cookieToken.Value
    }


	authHeader := c.Request().Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
		return parts[1]
	}

	return ""
}
//
func getMe(token string) error {
	if token == "valid-token" {
		// 認証成功
		return nil
	}
	// 認証失敗
	return fmt.Errorf("invalid token")
}
