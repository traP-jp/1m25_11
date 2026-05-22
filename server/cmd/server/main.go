package main

import (
	"context"
	"log"
	"os"

	"github.com/go-co-op/gocron/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/traP-jp/1m25_11/server/cmd/server/server"
	"github.com/traP-jp/1m25_11/server/pkg/config"
	"github.com/traP-jp/1m25_11/server/pkg/database"
)

func checkEnv() {
	appEnv := os.Getenv("APP_ENV")
	switch appEnv {
	case "production":
		log.Println("[OK]   APP_ENV=production")
	case "development":
		log.Println("[WARN] APP_ENV=development（DEV_USER フォールバックが有効）")
	case "":
		log.Fatal("[FAIL] APP_ENV: 未設定（\"production\" または \"development\" を明示してください）")
	default:
		log.Fatalf("[FAIL] APP_ENV=%q は不正な値です（\"production\" または \"development\" を設定してください）", appEnv)
	}

	isDev := config.IsDevelopment()

	if os.Getenv("BOT_TOKEN_KEY") == "" {
		if isDev {
			log.Println("[WARN] BOT_TOKEN_KEY: 未設定（UserCache が空になる）")
		} else {
			log.Fatal("[FAIL] BOT_TOKEN_KEY: 本番環境では必須")
		}
	} else {
		log.Println("[OK]   BOT_TOKEN_KEY")
	}

	if os.Getenv("PROXY_SECRET") == "" {
		if appEnv == "production" {
			log.Fatal("[FAIL] PROXY_SECRET: 本番環境では必須")
		}
		log.Println("[WARN] PROXY_SECRET: 未設定（本番では必須）")
	} else {
		log.Println("[OK]   PROXY_SECRET")
	}

	if appEnv == "production" && os.Getenv("DEV_USER") != "" {
		log.Println("[WARN] DEV_USER: production では無効だが設定されている")
	}

	if os.Getenv("ALLOWED_ORIGINS") == "" {
		log.Println("[WARN] ALLOWED_ORIGINS: 未設定（デフォルト値を使用）")
	} else {
		log.Println("[OK]   ALLOWED_ORIGINS")
	}
}

func main() {
	checkEnv()

	e := echo.New()

	// middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Dynamic CORS (credentials allowed) based on ALLOWED_ORIGINS env
	allowed := config.AllowedOrigins()
	e.Logger.Infof("CORS allowed origins: %v", allowed)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowed,
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH, echo.OPTIONS, echo.HEAD},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           600,
	}))

	// connect to and migrate database
	db, err := database.Setup(config.MySQL())
	if err != nil {
		e.Logger.Fatal(err)
	}
	defer db.Close()

	s := server.Inject(db)

	// Setup existing v1 API routes
	v1API := e.Group("/api/v1")
	s.SetupRoutes(v1API)

	//gocron
	ss, er := gocron.NewScheduler()
	if er != nil {
		log.Fatal(er)
	}

	_, err = ss.NewJob(
		gocron.CronJob("0 19 * * *", false),
		gocron.NewTask(s.Handler.CronJobTask, context.Background()),
	)
	if err != nil {
		log.Fatal(err)
	}

	// UserCacheを毎日午前3時に更新
	_, err = ss.NewJob(
		gocron.CronJob("0 3 * * *", false),
		gocron.NewTask(s.Handler.RefreshUserCache),
	)
	if err != nil {
		log.Fatal(err)
	}

	ss.Start()

	e.Logger.Fatal(e.Start(config.AppAddr()))

}
