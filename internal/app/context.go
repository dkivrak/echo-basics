package app

import (
	"log/slog"

	"go.smsk.dev/pkgs/basics/echo-basics/internal/config"
	"gorm.io/gorm"
)

type AppContext struct {
	DB     *gorm.DB
	Logger *slog.Logger
	Config config.Config
}
