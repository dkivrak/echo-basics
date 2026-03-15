package routes

import (
	"github.com/labstack/echo/v5"

	"go.smsk.dev/pkgs/basics/echo-basics/internal/app"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/config"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/handlers"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/middleware"
)

func Register(e *echo.Echo, cfg config.Config) {
	e.GET("/health", app.HealthCheck)

	api := e.Group("/api", middleware.APIKeyAuth(cfg.APIKey))

	api.POST("/logs", handlers.CreateLog)
	api.GET("/logs", handlers.FetchLogs)
	api.GET("/logs/id/:id", handlers.FetchID)
	api.GET("/logs/timestamp/:timestamp", handlers.FetchTimestamp)
	api.GET("/logs/flag/:flag", handlers.FetchFlag)
	api.DELETE("/logs/:id", handlers.DeleteLog)
}
