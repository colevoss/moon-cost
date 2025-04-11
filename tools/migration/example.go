package migration

import (
	"context"
	"database/sql"
	"log/slog"

	_ "github.com/tursodatabase/go-libsql"
)

func example() {
	url := "file:./local.db"
	db, err := sql.Open("libsql", url)

	if err != nil {
		panic(err)
	}

	mig := Manager{
		Dir: "./migrations",
		DB:  db,
	}

	slog.SetLogLoggerLevel(slog.LevelDebug)
	// slog.SetLogLoggerLevel(slog.LevelInfo)
	logger := slog.Default()

	mig.Init(WithLogger(logger))
	if err = mig.Run(context.Background()); err != nil {
		panic(err)
	}
	// _, err = mig.ReadMigrationsFromDir()
}
