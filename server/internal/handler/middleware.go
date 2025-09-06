package handler

import (
	// "context"
	"net/http"
	"strings"
	"log"

	// "github.com/coreos/go-oidc/v3/oidc"
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
		log.Printf("AuthMiddleware: Checking token")
		_, err := h.getUserID(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized: Invalid token")
		}
		

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

// func getMe(c echo.Context, rawToken string) error {
// 	config := &oidc.Config{
// 		ClientID: clientID,
// 	}
// 	provider, err := oidc.NewProvider(c.Request().Context(), "https://q.trap.jp")
// 	if err != nil {
// 		return err
// 	}

// 	ctx := context.Background()
// 	verifier := provider.Verifier(config)

// 	_, err = verifier.Verify(ctx, rawToken)
// 	if err != nil {

// 		return err
// 	}

// 	return nil
// }
