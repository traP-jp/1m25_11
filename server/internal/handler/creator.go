package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// クッキーからIDトークンを復元してユーザーIDを取得する関数
func (h *Handler) getUserID(c echo.Context) (uuid.UUID, error) {
	// デバッグのため、リクエストに含まれるCookie名をログに記録（値は記録しない）
	var presentNames []string
	for _, ck := range c.Request().Cookies() {
		presentNames = append(presentNames, ck.Name)
	}
	if len(presentNames) == 0 {
		log.Printf("getUserID: request contained no cookies")
	} else {
		log.Printf("getUserID: request cookie names=%v", presentNames)
	}

	// 分割されたクッキーの数を取得
	countCookie, err := c.Cookie(fmt.Sprintf("%s_count", tokenKey))
	if err != nil {
		log.Printf("getUserID: missing count cookie: %v", err)

		return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized: no auth token").SetInternal(err)
	}
	count, err := strconv.Atoi(countCookie.Value)
	if err != nil {
		log.Printf("getUserID: failed to parse token count (%s): %v", countCookie.Value, err)

		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "failed to parse token count").SetInternal(err)
	}

	log.Printf("getUserID: token count=%d", count)

	// 各クッキーを読み込んで結合
	var idTokenBuilder strings.Builder
	for i := 0; i < count; i++ {
		cookieName := fmt.Sprintf("%s_%d", tokenKey, i)
		cookie, err := c.Cookie(cookieName)
		if err != nil {
			log.Printf("getUserID: missing token chunk %d: %v", i, err)

			return uuid.Nil, echo.NewHTTPError(http.StatusUnauthorized, fmt.Sprintf("unauthorized: missing token part %d", i)).SetInternal(err)
		}
		idTokenBuilder.WriteString(cookie.Value)
	}
	idTokenString := idTokenBuilder.String()

	log.Printf("getUserID: reconstructed id_token length=%d", len(idTokenString))

	// IDトークンのペイロードをデコード
	parts := strings.Split(idTokenString, ".")
	if len(parts) != 3 {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "invalid ID token format")
	}

	payload := parts[1]
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusBadRequest, "failed to decode token payload")
	}

	var payloadMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payloadMap); err != nil {
		return uuid.Nil, echo.NewHTTPError(http.StatusInternalServerError, "failed to parse token payload")
	}

	// ユーザーID（"sub"）を取得してUUIDに変換
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
