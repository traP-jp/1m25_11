package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// ParseTraQJWTPayload はtraQのJWTを解析してペイロードを取得します（署名検証なし）
func ParseTraQJWTPayload(tokenString string) (map[string]interface{}, error) {
	// JWTの構造を分割 (header.payload.signature)
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	// ペイロード部分をデコード
	payloadPart := parts[1]
	// Base64URL デコーディングのためのパディング調整
	if padLength := len(payloadPart) % 4; padLength != 0 {
		payloadPart += strings.Repeat("=", 4-padLength)
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(payloadPart)
	if err != nil {
		return nil, err
	}

	// JSONとして解析
	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, err
	}

	return claims, nil
}

// SignJWT はJWTトークンを署名します
func SignJWT(claims *LightweightJWT, secret string) (string, error) {
	return claims.GenerateToken(secret)
}

// VerifyJWT はJWTトークンを検証し、クレームを返します
func VerifyJWT(tokenString string, secret string) (*LightweightJWT, error) {
	if len(secret) < 32 {
		return nil, errors.New("JWT secret must be at least 32 characters")
	}

	// JWTトークンをパースして検証
	token, err := jwt.ParseWithClaims(tokenString, &LightweightJWT{}, func(token *jwt.Token) (interface{}, error) {
		// 署名方式がHMAC-SHA256であることを確認
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	// クレームを取得
	claims, ok := token.Claims.(*LightweightJWT)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token or claims")
	}

	return claims, nil
}

// ValidateJWTSecret はJWT秘密鍵の形式をチェックします
func ValidateJWTSecret(secret string) error {
	if len(secret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters")
	}
	return nil
}
