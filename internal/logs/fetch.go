package logs

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

type PaginatedResponse struct {
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}

func getPaginationParams(c *echo.Context) (int, int, int) {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// çok saçma büyüklükte response dönmemek için üst sınır
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	return page, limit, offset
}

// FetchLogs returns paginated logs.
func FetchLogs(c *echo.Context) error {
	db := c.Get("app").(*AppContext).DB
	if db == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "app context missing"})
	}

	page, limit, offset := getPaginationParams(c)

	var total int64
	if err := db.WithContext(c.Request().Context()).Model(&Log{}).Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var logs []Log
	if err := db.WithContext(c.Request().Context()).
		Order("timestamp desc").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, PaginatedResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Data:  logs,
	})
}

// FetchID returns a single Log model by UUID.
// Returns 400 for bad id, 404 if not found, 500 for DB errors.
func FetchID(c *echo.Context) error {
	db := c.Get("app").(*AppContext).DB
	if db == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "app context missing"})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "bad request, i had better expectations from you.",
		})
	}

	var log Log
	res := db.WithContext(c.Request().Context()).First(&log, "id = ?", id)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "ummm, wait a bit... Nope, its not here. Maybe you didn't send me that log?",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": res.Error.Error()})
	}

	return c.JSON(http.StatusOK, log)
}

// FetchTimestamp returns the latest log at or before the provided timestamp.
// Route param must be an RFC3339 timestamp string.
func FetchTimestamp(c *echo.Context) error {
	db := c.Get("app").(*AppContext).DB
	if db == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "app context missing"})
	}

	tsStr := strings.TrimSpace(c.Param("timestamp"))
	if tsStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "bad request, i had better expectations from you."})
	}

	var ts time.Time
	var err error
	if ts, err = time.Parse(time.RFC3339, tsStr); err != nil {
		ts, err = time.Parse(time.RFC3339Nano, tsStr)
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "bad request, i had better expectations from you."})
	}

	tx := db.WithContext(c.Request().Context()).Begin()
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not start transaction"})
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var log Log
	res := tx.
		Where("timestamp <= ?", ts).
		Order("timestamp desc").
		First(&log)

	if res.Error != nil {
		tx.Rollback()

		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": "ummm, wait a bit... Nope, its not here. Maybe you didn't send me that log?",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": res.Error.Error()})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not commit transaction"})
	}

	return c.JSON(http.StatusOK, log)
}

// FetchFlag returns paginated logs filtered by flag (case-insensitive).
func FetchFlag(c *echo.Context) error {
	db := c.Get("app").(*AppContext).DB
	if db == nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "app context missing"})
	}

	flagParam := strings.ToLower(strings.TrimSpace(c.Param("flag")))
	if flagParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "bad request, i had better expectations from you."})
	}

	allowed := map[string]struct{}{
		string(LogFlag):   {},
		string(DebugFlag): {},
		string(InfoFlag):  {},
		string(WarnFlag):  {},
		string(ErrorFlag): {},
		string(TraceFlag): {},
	}

	if _, ok := allowed[flagParam]; !ok {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "WHAT'S THAT FLAG? I HAVE NEVER SEEN THAT!!!"})
	}

	page, limit, offset := getPaginationParams(c)

	tx := db.WithContext(c.Request().Context()).Begin()
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not start transaction"})
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	var total int64
	if err := tx.Model(&Log{}).Where("flag = ?", flagParam).Count(&total).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var logs []Log
	if err := tx.
		Where("flag = ?", flagParam).
		Order("timestamp desc").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "could not commit transaction"})
	}

	return c.JSON(http.StatusOK, PaginatedResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Data:  logs,
	})
}
