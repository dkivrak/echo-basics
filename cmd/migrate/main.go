package main

import (
	"fmt"

	"go.smsk.dev/pkgs/basics/echo-basics/internal/config"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/db"
	"go.smsk.dev/pkgs/basics/echo-basics/migrations"
)

func main() {
	cfg := config.MustLoad()
	database := db.InitDB(cfg)

	fmt.Println("Running migrations in env:", cfg.Env)

	if err := migrations.Run(database); err != nil {
		panic(fmt.Sprintf("migration failed: %v", err))
	}

	fmt.Println("Migrations completed successfully")
}
