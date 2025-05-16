package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"moon-cost/api"
	"moon-cost/logging"
	"moon-cost/moon"
	"moon-cost/services/auth"
	"net/http"
	"os"

	_ "github.com/tursodatabase/go-libsql"
)

func createAuth(db *sql.DB) *auth.Service {
	repo := auth.SQLiteRepo{
		DB: db,
		ID: moon.DefaultUUIDGenerator,
	}
	// repo := auth.NoopRepo{}
	return auth.NewService(&repo)
}

func run() int {
	prettyHandler := logging.NewPrettyHandler(os.Stderr, slog.LevelDebug)
	logger := slog.New(prettyHandler)
	slog.SetDefault(logger)

	// slog.SetLogLoggerLevel(slog.LevelDebug)

	cfg := api.Config{
		Port: 8080,
	}

	db, err := sql.Open("libsql", "file:./local.db")

	if err != nil {
		slog.Error("error opening db", "err", err)
		return 1
	}

	restApi := api.New(cfg)

	authSvc := createAuth(db)

	authController := api.AuthController{
		Auth: authSvc,
	}

	authController.Init(restApi)

	if err := http.ListenAndServe(restApi.Port(), restApi.Server.Mux); err != nil {
		fmt.Printf("ERR")
		return 1
	}

	return 0
}

func main() {
	os.Exit(run())
}
