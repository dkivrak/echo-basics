package logs

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

func CreateLog(c *echo.Context) error {
	db, err := getDB(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	var payload struct {
		Flag    FlagEnum `json:"flag"`
		Message string   `json:"message"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	payload.Message = strings.TrimSpace(payload.Message)
	if payload.Message == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "message is required",
		})
	}

	flagStr := strings.ToLower(strings.TrimSpace(string(payload.Flag)))
	if flagStr == "" {
		flagStr = string(InfoFlag)
	}

	allowed := map[string]struct{}{
		string(InfoFlag):  {},
		string(WarnFlag):  {},
		string(ErrorFlag): {},
	}

	if _, ok := allowed[flagStr]; !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid flag",
		})
	}

	log := Log{
		Flag:    FlagEnum(flagStr),
		Message: payload.Message,
	}

	if err := db.WithContext(c.Request().Context()).Create(&log).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not create log",
		})
	}

	return c.JSON(http.StatusCreated, log)
}
