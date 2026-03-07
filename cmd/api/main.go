package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/app"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/config"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/db"
	"go.smsk.dev/pkgs/basics/echo-basics/internal/logs"
)

// Welcome to my in-person golang tutorial
// including gorm, migrations, echo and more.
//
// This specific tutorial is built to create
// a semi-useful project. My students and I
// decided to go with a simple remote logging
// rest~ish api with using latest versions of
// all dependencies.
//
// Have fun, if you have some ideas (i have used
// some bad practices across the code), create a
// pr to practice both your pr and your coding skills

// dkivrak was here heheheh

func main() {
	config.LoadEnv()

	limitRateStr := os.Getenv("LIMIT_RATE")

	if limitRateStr == "" {
		panic("LIMIT_RATE not set")
	}

	limitRate, err := strconv.ParseFloat(os.Getenv("LIMIT_RATE"), 64)
	if err != nil {
		panic(err)
	}

	db := db.InitDB()

	// Create new echo instance
	e := echo.New()
	// tell echo to use default logger middleware.
	e.Use(middleware.RequestLogger())
	// Load default security practices like,
	// xss, content type sniffing, clickjacking and more.
	e.Use(middleware.Secure())
	// Simple rate limitter setup
	// remove if you don't need it :=)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(limitRate)))

	// Inject AppContext into every request to be able
	// to access database or other utilities easily.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set("app", &logs.AppContext{DB: db})
			return next(c)
		}
	})

	e.GET("/", func(c *echo.Context) error {
		// Sample db usage
		// db := c.Get("app").(*AppContext).DB
		// _ = db.Exec("SELECT 1;")
		return c.String(http.StatusOK, "Hellow from the other side")
	})

	e.GET("/health", app.HealthCheck)

	// /api/*
	api := e.Group("/api") //api routes

	// route structure
	// - list logs (get)
	// - create logs
	// - fetch logs
	//   - fetch based on id (get)
	//   - fetch based on timestamp (get)
	//   - fetch based on flag (get)
	// - delete logs
	//   - delete based on id (only allowed when the flag is lower i.e. 4) (delete)

	api.POST("/create", logs.CreateLog) // create log

	// task: create paginated responses over list returns
	api.GET("/list", logs.FetchLogs)                    // returns list of logs
	api.GET("/fetch/i/:id", logs.FetchID)               // fetch based on id (returns exactly one log instance)
	api.GET("/fetch/t/:timestamp", logs.FetchTimestamp) // fetch based on timestamp (returns the latest given log (possibly only one))
	api.GET("/fetch/f/:flag", logs.FetchFlag)           // fetch based on flag (returns a list)

	api.DELETE("/delete/:id", logs.DeleteLog) // delete log

	if err := e.Start(":" + os.Getenv("PORT")); err != nil { //
		e.Logger.Error("Failed to start echo application", "error", err)
	}
}
