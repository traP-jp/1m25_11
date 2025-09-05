package handler

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *Handler) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if strings.HasSuffix(c.Request().URL.Path, "/login") {
			return next(c)
		}
		if strings.Contains(c.Request().URL.Path, "/callback") {
			return next(c)
		}

		token := getToken(c)
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: No token")
		}

		// 軽量JWT検証
		claims, err := verifyLightweightJWT(token)
		if err != nil {
			log.Printf("JWT verification failed: %v", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: Invalid token")
		}

		// ユーザー情報をコンテキストに設定
		c.Set("user_id", claims.Sub)
		c.Set("username", claims.PreferredUsername)
		c.Set("display_name", claims.Name)
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

// verifyLightweightJWT は軽量JWTを検証します
func verifyLightweightJWT(tokenString string) (*LightweightJWT, error) {
	secret, err := GetJWTSecret()
	if err != nil {
		return nil, err
	}

	return VerifyJWT(tokenString, secret)
}
