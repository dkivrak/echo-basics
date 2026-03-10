package middleware

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

func APIKeyAuth(expectedKey string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			givenKey := c.Request().Header.Get("X-API-Key")

			if givenKey == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "missing api key",
				})
			}

			if givenKey != expectedKey {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid api key",
				})
			}

			return next(c)
		}
	}
}
