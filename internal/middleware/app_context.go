package middleware

import (
	"github.com/labstack/echo/v5"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/app"
)

func InjectAppContext(appCtx *app.AppContext) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("app", appCtx)
			return next(c)
		}
	}
}
