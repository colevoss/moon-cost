package main

import (
	"context"
	"database/sql"
	"log/slog"
	"moon-cost/migration"

	_ "github.com/tursodatabase/go-libsql"
)

func main() {
	url := "file:./local.db"
	db, err := sql.Open("libsql", url)

	if err != nil {
		panic(err)
	}

	mig := migration.Manager{
		Dir: "./migrations",
		DB:  db,
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)
	// slog.SetLogLoggerLevel(slog.LevelInfo)
	logger := slog.Default()

	mig.Init(migration.WithLogger(logger))
	if err = mig.Run(context.Background()); err != nil {
		panic(err)
	}
	// _, err = mig.ReadMigrationsFromDir()
}
