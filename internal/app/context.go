package app

import (
	"log/slog"

	"github.com/labstack/echo/v5"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/config"
	"gorm.io/gorm"
)

type AppContext struct {
	DB     *gorm.DB
	Logger *slog.Logger
	Config config.Config
}

func InjectAppContext(appCtx *AppContext) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("app", appCtx)
			return next(c)
		}
	}
}
