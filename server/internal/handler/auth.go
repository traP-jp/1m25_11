package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

var requestURL = "https://q.trap.jp/api/v3/oauth2/authorize"
var tokenURL = "https://q.trap.jp/api/v3/oauth2/token"
var tokenKey = "traq-auth-token"

type TokenData struct {
	TokenType   string `json:"token_type"`
	AccessToken string `json:"access_token"`	
	IDToken     string `json:"id_token"`
	ExpiresIn   int    `json:"expires_in"`
}

var clientID = os.Getenv("CLIENT_ID")

func (h *Handler) login(c echo.Context) error {
	redirectURI, codeVerifier, state, err := h.getTraqAuthCode(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
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
	
	resCodeVerifier, err := c.Cookie(h.codeVerifierKey(state))
	if err != nil {
		return c.Redirect(http.StatusFound, "/")
	}
	codeVerifier := resCodeVerifier.Value
	tokenRes, err := h.sendTraqAuthToken(code, codeVerifier)
	if err != nil {
		return c.Redirect(http.StatusFound, "/")
	}
	defer tokenRes.Body.Close()
	var tokenData TokenData
	err = json.NewDecoder(tokenRes.Body).Decode(&tokenData)
	if err != nil {
		return c.Redirect(http.StatusFound, "/")
	}
	idToken := tokenData.IDToken
	deleteCookie := &http.Cookie{
    Name:     h.codeVerifierKey(state),
    Secure:   true,
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,
    Value:    "",    
    MaxAge:   -1,    
}
	c.SetCookie(deleteCookie)

	cookie := &http.Cookie{
		Name:     tokenKey,
		Value:    idToken,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode, // 'None'に対応
		MaxAge:   tokenData.ExpiresIn,
	}

	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/")
}

func (h *Handler) getTraqAuthCode(c echo.Context) (string, string, string, error) {
	state, err := h.randomString(10)
	if err != nil {
		return "", "", "", err
	}
	codeVerifier, err := h.randomString(43)
	if err != nil {
		return "", "", "", err
	}
	codeChallenge := h.getCodeChallenge(codeVerifier)
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", clientID)
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")
	u, err := url.Parse(requestURL)
	if err != nil {
		return "", "", "", err
	}
	u.RawQuery = params.Encode()
	redirectURI := u.String()

	return redirectURI, codeVerifier, state, nil
}

func (h *Handler) codeVerifierKey(state string) string {
	return "traq-auth-code-verifier-" + state
}
func (h *Handler) sendTraqAuthToken(code string, codeVerifier string) (*http.Response, error) {
	const baseURL = "https://q.trap.jp/api/v3/oauth2"

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("code", code)
	data.Set("code_verifier", codeVerifier)

	return http.PostForm(baseURL+"/token", data)
}

func (h *Handler) randomString(length int) (string, error) {
	const characters = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	var builder strings.Builder
	builder.Grow(length)

	max := big.NewInt(int64(len(characters)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		builder.WriteByte(characters[n.Int64()])
	}

	return builder.String(), nil
}

func (h *Handler) getCodeChallenge(codeVerifier string) string {
	hasher := sha256.New()
	hasher.Write([]byte(codeVerifier))
	shaSum := hasher.Sum(nil)
	encoder := base64.URLEncoding.WithPadding(base64.NoPadding)

	return encoder.EncodeToString(shaSum)
}
