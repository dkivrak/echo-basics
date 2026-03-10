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

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	return page, limit, offset
}

// FetchLogs returns paginated logs.
func FetchLogs(c *echo.Context) error {
	db, err := getDB(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	page, limit, offset := getPaginationParams(c)

	var total int64
	if err := db.WithContext(c.Request().Context()).
		Model(&Log{}).
		Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	var logs []Log
	if err := db.WithContext(c.Request().Context()).
		Order("timestamp desc").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, PaginatedResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Data:  logs,
	})
}

// FetchID returns a single Log model by UUID.
func FetchID(c *echo.Context) error {
	db, err := getDB(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
	}

	var log Log
	res := db.WithContext(c.Request().Context()).
		First(&log, "id = ?", id)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "log not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": res.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, log)
}

// FetchTimestamp returns the latest log at or before the provided timestamp.
func FetchTimestamp(c *echo.Context) error {
	db, err := getDB(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	tsStr := strings.TrimSpace(c.Param("timestamp"))
	if tsStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "timestamp is required",
		})
	}

	ts, err := time.Parse(time.RFC3339Nano, tsStr)
	if err != nil {
		ts, err = time.Parse(time.RFC3339, tsStr)
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid timestamp format",
		})
	}

	var log Log
	res := db.WithContext(c.Request().Context()).
		Where("timestamp <= ?", ts).
		Order("timestamp desc").
		First(&log)

	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "log not found",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": res.Error.Error(),
		})
	}

	return c.JSON(http.StatusOK, log)
}

// FetchFlag returns paginated logs filtered by flag.
func FetchFlag(c *echo.Context) error {
	db, err := getDB(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	flagParam := strings.ToLower(strings.TrimSpace(c.Param("flag")))
	if flagParam == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "flag is required",
		})
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
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid flag",
		})
	}

	page, limit, offset := getPaginationParams(c)

	var total int64
	if err := db.WithContext(c.Request().Context()).
		Model(&Log{}).
		Where("flag = ?", flagParam).
		Count(&total).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	var logs []Log
	if err := db.WithContext(c.Request().Context()).
		Where("flag = ?", flagParam).
		Order("timestamp desc").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, PaginatedResponse{
		Page:  page,
		Limit: limit,
		Total: total,
		Data:  logs,
	})
}
