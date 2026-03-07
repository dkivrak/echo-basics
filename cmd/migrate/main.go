package main

import (
	"fmt"

	"go.smsk.dev/pkgs/basics/echo-basics/internal/config"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/db"
	"go.smsk.dev/pkgs/basics/echo-basics/migrations"
)

func main() {
	env := config.LoadEnv()
	fmt.Println("Running migrations in env:", env)

	database := db.InitDB()

	if err := migrations.Run(database); err != nil {
		panic(fmt.Sprintf("migration failed: %v", err))
	}

	fmt.Println("Migrations completed successfully")
}
