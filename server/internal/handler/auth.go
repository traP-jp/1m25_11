package handler

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/labstack/echo/v4"
	"github.com/traP-jp/1m25_11/server/pkg/config"
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
var topPageURL = os.Getenv("TOP_PAGE_URL")

func (h *Handler) login(c echo.Context) error {
	redirectURI, codeVerifier, state, err := h.getTraqAuthCode(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError).SetInternal(err)
	}
	cookie := &http.Cookie{
		Name:     h.codeVerifierKey(state),
		Value:    codeVerifier,
		MaxAge:   60 * 60, // 3600秒 = 1時間
		Secure:   config.GetCookieSecure(),
		HttpOnly: true,
		Path:     "/",
		SameSite: config.GetCookieSameSite(),
	}
	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, redirectURI)

}
func (h *Handler) callback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	if code == "" || state == "" {
		return c.Redirect(http.StatusFound, topPageURL)
	}

	resCodeVerifier, err := c.Cookie(h.codeVerifierKey(state))
	if err != nil {
		return c.Redirect(http.StatusFound, topPageURL)
	}
	codeVerifier := resCodeVerifier.Value
	tokenRes, err := h.sendTraqAuthToken(code, codeVerifier)
	if err != nil {
		log.Printf("failed to request token endpoint: %v", err)

		return c.Redirect(http.StatusFound, topPageURL)
	}
	defer tokenRes.Body.Close()

	// Check HTTP status code from token endpoint
	if tokenRes.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(tokenRes.Body)
		log.Printf("token endpoint returned status %d: %s", tokenRes.StatusCode, string(body))

		return c.Redirect(http.StatusFound, topPageURL)
	}

	var tokenData TokenData
	if err := json.NewDecoder(tokenRes.Body).Decode(&tokenData); err != nil {
		log.Printf("failed to decode token response: %v", err)

		return c.Redirect(http.StatusFound, topPageURL)
	}
	// Verify ID token signature and claims using OIDC provider
	idToken := tokenData.IDToken
	if idToken == "" {
		log.Printf("no id_token in token response")

		return c.Redirect(http.StatusFound, topPageURL)
	}
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://q.trap.jp")
	if err != nil {
		log.Printf("failed to create oidc provider: %v", err)

		return c.Redirect(http.StatusFound, topPageURL)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	if _, err := verifier.Verify(ctx, idToken); err != nil {
		log.Printf("id_token verification failed: %v", err)

		return c.Redirect(http.StatusFound, topPageURL)
	}
	log.Printf("callback: id_token verified for state=%s", state)
	// id_token verified, proceed to set cookies and delete code verifier
	deleteCookie := &http.Cookie{
		Name:     h.codeVerifierKey(state),
		Secure:   config.GetCookieSecure(),
		HttpOnly: true,
		Path:     "/",
		SameSite: config.GetCookieSameSite(),
		Value:    "",
		MaxAge:   -1,
	}
	c.SetCookie(deleteCookie)
	log.Printf("callback: set cookie %s (Path=%s, Secure=%t, SameSite=%v, MaxAge=%d)", deleteCookie.Name, deleteCookie.Path, deleteCookie.Secure, deleteCookie.SameSite, deleteCookie.MaxAge)

	var chunks []string
	const maxCookieSize = 3500 // クッキーの最大サイズ（バイト）

	for i := 0; i < len(idToken); i += maxCookieSize {
		end := i + maxCookieSize
		if end > len(idToken) {
			end = len(idToken)
		}
		chunks = append(chunks, idToken[i:end])
	}

	// 分割したチャンクをそれぞれクッキーとして設定
	for i, chunk := range chunks {
		cookieName := fmt.Sprintf("%s_%d", tokenKey, i)
		cookie := &http.Cookie{
			Name:     cookieName,
			Value:    chunk,
			Secure:   config.GetCookieSecure(),
			HttpOnly: true,
			Path:     "/",
			SameSite: config.GetCookieSameSite(),
			MaxAge:   tokenData.ExpiresIn,
		}
		c.SetCookie(cookie)
		log.Printf("callback: set cookie %s (Path=%s, Secure=%t, SameSite=%v, MaxAge=%d)", cookie.Name, cookie.Path, cookie.Secure, cookie.SameSite, cookie.MaxAge)
	}

	// クッキーの総数を記録するためのクッキーも設定
	countCookie := &http.Cookie{
		Name:     fmt.Sprintf("%s_count", tokenKey),
		Value:    strconv.Itoa(len(chunks)),
		Secure:   config.GetCookieSecure(),
		HttpOnly: true,
		Path:     "/",
		SameSite: config.GetCookieSameSite(),
		MaxAge:   tokenData.ExpiresIn,
	}
	c.SetCookie(countCookie)
	log.Printf("callback: set cookie %s (Path=%s, Secure=%t, SameSite=%v, MaxAge=%d)", countCookie.Name, countCookie.Path, countCookie.Secure, countCookie.SameSite, countCookie.MaxAge)

	log.Printf("callback: set %d token cookies (count=%s)", len(chunks), countCookie.Value)

	return c.Redirect(http.StatusFound, topPageURL)
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
