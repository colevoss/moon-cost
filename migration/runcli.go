package migration

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"

	_ "github.com/tursodatabase/go-libsql"
)

type runCli struct {
	dir        string
	table      string
	dbFilename string
	cli        *MigrationCLI
}

func (r *runCli) init(args []string) error {
	fs := flag.NewFlagSet("run", flag.ExitOnError)

	fs.StringVar(&r.dir, "dir", "migrations", "Directory to find migration files")
	fs.StringVar(&r.table, "table", "migrations", "Database table that migration data is stored in")
	fs.StringVar(&r.dbFilename, "db", "", "SQLite File to run migrations against")
	r.cli.parseUniversalFlags(fs)

	fs.Parse(args)

	if r.dbFilename == "" {
		return fmt.Errorf("Error: db flag required")
	}

	return nil
}

func (r *runCli) run() error {
	dbName := fmt.Sprintf("file:%s", r.dbFilename)

	db, err := sql.Open("libsql", dbName)

	if err != nil {
		return err
	}

	defer db.Close()

	logger := slog.Default()
	ctx := context.Background()

	manager := Manager{
		Dir:    r.dir,
		Table:  r.table,
		DB:     db,
		logger: logger,
	}

	return manager.Run(ctx)
}
