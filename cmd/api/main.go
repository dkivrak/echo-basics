package main

import (
	"net/http"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"

	"go.smsk.dev/pkgs/basics/echo-basics/internal/app"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/config"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/db"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/logger"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/logs"
	customMiddleware "go.smsk.dev/pkgs/basics/echo-basics/internal/middleware"
)

func main() {
	cfg := config.MustLoad()
	log := logger.New(cfg.Env)
	database := db.InitDB(cfg)

	appCtx := &app.AppContext{
		DB:     database,
		Logger: log,
		Config: cfg,
	}

	e := echo.New()

	e.Use(echoMiddleware.RequestLogger())
	e.Use(echoMiddleware.Secure())
	e.Use(echoMiddleware.RateLimiter(
		echoMiddleware.NewRateLimiterMemoryStore(cfg.LimitRate),
	))
	e.Use(customMiddleware.InjectAppContext(appCtx))

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello from the other side")
	})

	e.GET("/health", app.HealthCheck)

	api := e.Group("/api", customMiddleware.APIKeyAuth(cfg.APIKey))

	api.POST("/logs", logs.CreateLog)
	api.GET("/logs", logs.FetchLogs)
	api.GET("/logs/id/:id", logs.FetchID)
	api.GET("/logs/timestamp/:timestamp", logs.FetchTimestamp)
	api.GET("/logs/flag/:flag", logs.FetchFlag)
	api.DELETE("/logs/:id", logs.DeleteLog)

	log.Info("starting api server", "port", cfg.Port, "env", cfg.Env)

	if err := e.Start(":" + cfg.Port); err != nil {
		log.Error("failed to start echo application", "error", err)
	}
}
