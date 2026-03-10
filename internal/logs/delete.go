package logs

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func DeleteLog(c *echo.Context) error {
	db, err := getDB(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if db == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "app context missing",
		})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "bad request, i had better expectations from you.",
		})
	}

	tx := db.WithContext(c.Request().Context()).Begin()
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

	var log Log
	if err := tx.First(&log, "id = ?", id).Error; err != nil {
		tx.Rollback()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "umm, wait a bit... Nope, its not here. Maybe you didn't send me that log?",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	if GetLogLevel(log.Flag) >= 4 {
		tx.Rollback()
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "EEEEYYYY! You can't do that!",
		})
	}

	if err := tx.Delete(&log).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not delete log",
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not commit transaction",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "That record is long gone now. Don't worry, our secret is now safe.",
	})
}
