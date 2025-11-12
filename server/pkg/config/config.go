package config

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func getEnv(key, defaultValue string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}

	return v
}

func AppAddr() string {
	return getEnv("APP_ADDR", ":8080")
}

func MySQL() *mysql.Config {
	c := mysql.NewConfig()

	c.User = getEnv("NS_MARIADB_USER", "root")
	c.Passwd = getEnv("NS_MARIADB_PASSWORD", "pass")
	c.Net = getEnv("DB_NET", "tcp")
	c.Addr = fmt.Sprintf(
		"%s:%s",
		getEnv("NS_MARIADB_HOSTNAME", "localhost"),
		getEnv("NS_MARIADB_PORT", "3306"),
	)
	c.DBName = getEnv("NS_MARIADB_DATABASE", "app")
	c.Collation = "utf8mb4_bin"
	c.AllowNativePasswords = true
	c.ParseTime = true

	return c
}

// AllowedOrigins はCORSで許可されるオリジンのリストを返す
// ALLOWED_ORIGINS環境変数でカンマ区切りで指定
func AllowedOrigins() []string {
	raw := getEnv("ALLOWED_ORIGINS", "https://stampedia.trap.show,https://1m25-11.trap.show,http://localhost:3000")
	parts := strings.Split(raw, ",")
	origins := make([]string, 0, len(parts))
	for _, p := range parts {
		o := strings.TrimSpace(p)
		if o == "" {
			continue
		}
		origins = append(origins, o)
	}

	return origins
}

// 環境変数APP_ENVを確認して、開発モードで実行されているかを IsDevelopment に
func IsDevelopment() bool {
	// APP_ENV変数で明示的に環境を判定。デフォルトは "development"
	env := getEnv("APP_ENV", "development")

	return env == "development"
}

// GetCookieSameSite は環境に基づいて適切なSameSite設定を返す
func GetCookieSameSite() http.SameSite {
	if IsDevelopment() {
		return http.SameSiteLaxMode
	}

	return http.SameSiteNoneMode
}

// GetCookieSecure は環境に基づいて適切なSecure設定を返す
func GetCookieSecure() bool {
	// 開発環境ではHTTPも許可、本番環境ではHTTPS必須
	return !IsDevelopment()
}
