package handler

import (
	"encoding/base64"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)



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
	var payloadMap map[string]interface{}

	if err := json.Unmarshal(payloadBytes, &payloadMap); err != nil {
        return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to parse token payload")
    }

	creatorIDStr, ok := payloadMap["sub"].(string)
	if !ok {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid token payload")
	}
	creatorID, err := uuid.Parse(creatorIDStr)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest).SetInternal(err)
	}



	return creatorID, nil
}
