package handler

import (
	"encoding/json"
	"net/http"
	"math/rand"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/ml/deletecalendar"
	"github.com/labstack/echo/v4"
	"golang.org/x/tools/godoc/redirect"
)

var requestURL = "https://q.trap.jp/api/v3/oauth2/authorize"
var tokenURL = "https://q.trap.jp/api/v3/oauth2/token"
var tokenKey = "traq-auth-token"

type TokenData struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}



func (h *Handler) login(c echo.Context) error {
	redirectURI, codeVerifier, state := h.getTraqAuthCode(c)
	cookie := &http.Cookie{
		Name:     h.codeVerifierKey(state),
		Value:    codeVerifier,
		MaxAge:   60 * 60, // 3600秒 = 1時間
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, redirectURI)

}
func (h *Handler) callback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	if code == "" || state == "" {
		return c.Redirect(http.StatusFound, "/")
	}
	codeVerifier, err := c.Cookie(h.codeVerifierKey(state))
	if err != nil {
		return c.Redirect(http.StatusFound, "/")
	}
	tokenRes := h.sendTraqAuthToken(code, codeVerifier)
	var tokenData TokenData
	err = json.NewDecoder(tokenRes.Body).Decode(&tokenData)
	if err != nil {
		return c.Redirect(http.StatusFound, "/")
	}
	token := tokenData.AccessToken
	deleteCookie := &http.Cookie{
		Name:   h.codeVerifierKey(state),
		Value:  "",
		MaxAge: -1,
	}
	c.SetCookie(deleteCookie)

	cookie := &http.Cookie{
        Name:     tokenKey,
        Value:    token,
        Secure:   true,
        HttpOnly: true,
        SameSite: http.SameSiteNoneMode, // 'None'に対応
        MaxAge:   tokenData.ExpiresIn,
    }

	c.SetCookie(cookie)
	
	return nil
}

func (h *Handler) getTraqAuthCode(c echo.Context) (string, string, string) {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return "redirectURI", "codeVerifier", "state"
}
func (h *Handler) codeVerifierKey(state string) string {
	return "traq-auth-code-verifier-" + state
}
func (h *Handler) sendTraqAuthToken(code string, codeVerifier string) {

}
