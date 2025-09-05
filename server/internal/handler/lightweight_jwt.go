package handler

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// LightweightJWT は軽量化されたJWT構造体です
type LightweightJWT struct {
	Sub               string `json:"sub"`                // ユーザーID
	PreferredUsername string `json:"preferred_username"` // ユーザー名
	Name              string `json:"name"`               // 表示名
	Exp               int64  `json:"exp"`                // 有効期限
	Iat               int64  `json:"iat"`                // 発行時刻
	Iss               string `json:"iss"`                // 発行者
	jwt.RegisteredClaims
}

// TraQUserInfo はtraQのJWTから抽出するユーザー情報です
type TraQUserInfo struct {
	Sub               string      `json:"sub"`
	PreferredUsername string      `json:"preferred_username"`
	Name              string      `json:"name"`
	Exp               interface{} `json:"exp"`
	Iat               interface{} `json:"iat"`
}

// CreateLightweightJWTFromTraQ は元のtraQ JWTから軽量JWTを作成します
func CreateLightweightJWTFromTraQ(originalToken string) (*LightweightJWT, error) {
	// traQのJWTをパースしてペイロード取得（署名検証なし）
	originalClaims, err := ParseTraQJWTPayload(originalToken)
	if err != nil {
		return nil, err
	}

	// 必要な情報を抽出
	userInfo := &TraQUserInfo{}
	claimsBytes, err := json.Marshal(originalClaims)
	if err != nil {
		return nil, err
	}
	
	if err := json.Unmarshal(claimsBytes, userInfo); err != nil {
		return nil, err
	}

	// expをint64に変換
	var exp int64
	switch v := userInfo.Exp.(type) {
	case float64:
		exp = int64(v)
	case int64:
		exp = v
	case int:
		exp = int64(v)
	default:
		return nil, errors.New("invalid exp format")
	}

	// iatをint64に変換
	var iat int64
	switch v := userInfo.Iat.(type) {
	case float64:
		iat = int64(v)
	case int64:
		iat = v
	case int:
		iat = int64(v)
	default:
		iat = time.Now().Unix()
	}

	// 軽量JWT構造体を作成
	lightweightJWT := &LightweightJWT{
		Sub:               userInfo.Sub,
		PreferredUsername: userInfo.PreferredUsername,
		Name:              userInfo.Name,
		Exp:               exp,
		Iat:               iat,
		Iss:               "1m25-app",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userInfo.Sub,
			ExpiresAt: jwt.NewNumericDate(time.Unix(exp, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Unix(iat, 0)),
			Issuer:    "1m25-app",
		},
	}

	return lightweightJWT, nil
}

// GenerateToken は軽量JWTトークン文字列を生成します
func (l *LightweightJWT) GenerateToken(secret string) (string, error) {
	if len(secret) < 32 {
		return "", errors.New("JWT secret must be at least 32 characters")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, l)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GetJWTSecret はJWT署名用の秘密鍵を環境変数から取得します
func GetJWTSecret() (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		return "", errors.New("JWT_SECRET environment variable must be at least 32 characters")
	}
	return secret, nil
}
