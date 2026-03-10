package logs

import (
	"errors"

	"github.com/labstack/echo/v5"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/app"
	"gorm.io/gorm"
)

func getDB(c *echo.Context) (*gorm.DB, error) {
	appCtx, ok := c.Get("app").(*app.AppContext)
	if !ok || appCtx == nil || appCtx.DB == nil {
		return nil, errors.New("app context missing")
	}
	return appCtx.DB, nil
}
