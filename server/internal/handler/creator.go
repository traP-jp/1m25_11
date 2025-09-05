package handler

import (
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"time"
)

type IDToken struct {
	Audience          string `json:"aud"`
	ExpiresAt         int64  `json:"exp"`
	IssuedAt          int64  `json:"iat"`
	Issuer            string `json:"iss"`
	Name              string `json:"name"`
	Picture           string `json:"picture"`
	PreferredUsername string `json:"preferred_username"`
	Subject           string `json:"sub"`
	Traq              Traq   `json:"traq"`
	UpdatedAt         int64  `json:"updated_at"`
}
type Traq struct {
	Bio         string    `json:"bio"`
	Bot         bool      `json:"bot"`
	DisplayName string    `json:"display_name"`
	Groups      []string  `json:"groups"`
	HomeChannel string    `json:"home_channel"`
	IconFileID  string    `json:"icon_file_id"`
	LastOnline  time.Time `json:"last_online"`
	Permissions []string  `json:"permissions"`
	State       int       `json:"state"`
	Tags        []string  `json:"tags"`
	TwitterID   string    `json:"twitter_id"`
}

func (h *Handler) getUserID(c echo.Context) (uuid.UUID, error) {
	cookie, err := c.Cookie(tokenKey)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: no auth token").SetInternal(err)
	}
	idTokenString := cookie.Value
	parts := strings.Split(idTokenString, ".")
	payload := parts[1]
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "failed to decode token payload")
	}
	var idToken IDToken

	if err := json.Unmarshal(payloadBytes, &idToken); err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to parse token payload")
	}

	creatorID, err := uuid.Parse(idToken.Name)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}

	return creatorID, nil
}
