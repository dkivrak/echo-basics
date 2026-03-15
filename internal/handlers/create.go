package handlers

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v5"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/models"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/utils"
)

func CreateLog(c *echo.Context) error {
	db, err := utils.GetDB(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	var payload struct {
		Flag    utils.FlagEnum `json:"flag"`
		Message string         `json:"message"`
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
		flagStr = string(utils.InfoFlag)
	}

	allowed := map[string]struct{}{
		string(utils.InfoFlag):  {},
		string(utils.WarnFlag):  {},
		string(utils.ErrorFlag): {},
	}

	if _, ok := allowed[flagStr]; !ok {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid flag",
		})
	}

	log := models.Log{
		Flag:    utils.FlagEnum(flagStr),
		Message: payload.Message,
	}

	if err := db.WithContext(c.Request().Context()).Create(&log).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not create log",
		})
	}

	return c.JSON(http.StatusCreated, log)
}
