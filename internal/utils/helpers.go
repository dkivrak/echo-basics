package utils

import (
	"errors"

	"github.com/labstack/echo/v5"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/app"
	"gorm.io/gorm"
)

// GetDB returns the database from the application context stored in Echo context.
// Exported so other packages can obtain the DB instance (was getDB previously).
func GetDB(c *echo.Context) (*gorm.DB, error) {
	appCtx, ok := c.Get("app").(*app.AppContext)
	if !ok || appCtx == nil || appCtx.DB == nil {
		return nil, errors.New("app context missing")
	}
	return appCtx.DB, nil
}
