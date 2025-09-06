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
	c.Collation = "utf8mb4_general_ci"
	c.AllowNativePasswords = true
	c.ParseTime = true

	return c
}

// AllowedOrigins returns the list of CORS allowed origins.
// Comma separated in ALLOWED_ORIGINS environment variable.
// Defaults cover production & local development frontends.
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

// IsDevelopment determines if the application is running in development mode
// by checking if localhost is present in ALLOWED_ORIGINS
func IsDevelopment() bool {
	origins := AllowedOrigins()
	for _, origin := range origins {
		if strings.Contains(origin, "localhost") {
			return true
		}
	}

	return false
}

// GetCookieSameSite returns the appropriate SameSite setting based on environment
func GetCookieSameSite() http.SameSite {
	if IsDevelopment() {
		return http.SameSiteLaxMode
	}

	return http.SameSiteNoneMode
}

// GetCookieSecure returns the appropriate Secure setting based on environment
func GetCookieSecure() bool {
	// 開発環境ではHTTPも許可、本番環境ではHTTPS必須
	return !IsDevelopment()
}
