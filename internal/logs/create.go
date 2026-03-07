package logs

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
)

func CreateLog(c *echo.Context) error {
	db := c.Get("app").(*AppContext).DB
	if db == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "app context missing"})
	}

	var payload struct {
		Flag    FlagEnum `json:"flag"`
		Message string   `json:"message" validate:"required"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "bad request, I had better expectations from you.",
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

	// transaction başlat
	tx := db.Begin()
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not start transaction",
		})
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	log := Log{
		Flag:    FlagEnum(flagStr),
		Message: payload.Message,
	}

	if err := tx.Create(&log).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not create log",
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not commit transaction",
		})
	}

	return c.JSON(http.StatusCreated, log)
}
