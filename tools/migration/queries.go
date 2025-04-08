package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

const ensureMigrationsTableQuery = `
CREATE TABLE IF NOT EXISTS %s (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  filename TEXT NOT NULL,
  created INTEGER NOT NULL,
  instruction TEXT NOT NULL
)
`

func (m *Manager) ensureMigrationsTable(ctx context.Context) error {
	m.logger.Debug("Ensuring migrations table exists", "table", m.Table)

	_, err := m.DB.ExecContext(ctx, m.formatQuery(ensureMigrationsTableQuery))

	if err != nil {
		return fmt.Errorf("Error ensuring migrations table %s: %w", m.Table, err)
	}

	return nil
}

const getAllMigrationsQuery = `
SELECT
  id,
  name,
  filename,
  created,
  instruction
FROM %s
ORDER BY created ASC;
`

func (m *Manager) getAllMigrations(ctx context.Context) ([]Migration, error) {
	rows, err := m.DB.QueryContext(ctx, m.formatQuery(getAllMigrationsQuery))

	if err != nil {
		return nil, fmt.Errorf("Error querying for migrations: %w", err)
	}

	defer rows.Close()

	var migrations []Migration

	for rows.Next() {
		var migration Migration
		var createdInt int64

		rows.Scan(
			&migration.Id,
			&migration.Name,
			&migration.Filename,
			&createdInt,
			&migration.Instruction,
		)

		migration.Created = time.UnixMilli(createdInt)
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

const createMigrationQuery = `
INSERT INTO %s (
  name,
  filename,
  created,
  instruction
)
VALUES (?, ?, ?, ?);
`

func (m *Manager) createMigrations(ctx context.Context, migrations []Migration) error {
	tx, err := m.DB.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(m.formatQuery(createMigrationQuery))
	defer stmt.Close()

	if err != nil {
		return err
	}

	for _, migration := range migrations {
		m.logger.Debug("Running migration", "name", migration.Name, "file", migration.Filename, "query", migration.Instruction)

		if err := m.runMigration(ctx, tx, migration); err != nil {
			if rbError := tx.Rollback(); rbError != nil {
				m.logger.Error("error running rolling back migration", "error", err)
				return rbError
			}

			return fmt.Errorf("Error running migration %w", err)
		}

		m.logger.Info("Migration ran successfully", "file", migration.Filename)

		m.logger.Debug("Creating migration", "name", migration.Name, "created", migration.Created)

		res, err := stmt.ExecContext(
			ctx,
			migration.Name,
			migration.Filename,
			migration.Created.UnixMilli(),
			migration.Instruction,
		)

		if err != nil {
			if rbError := tx.Rollback(); rbError != nil {
				return rbError
			}

			return err
		}

		affected, _ := res.RowsAffected()
		id, _ := res.LastInsertId()

		m.logger.Debug("created migration in transaction", "affected", affected, "id", id)
	}

	if err := tx.Commit(); err != nil {
		m.logger.Error("Error committing migration creation transaction", "error", err)
		return err
	}

	return nil
}

func (m *Manager) runMigration(ctx context.Context, tx *sql.Tx, migration Migration) error {
	queries := splitQueries(migration.Instruction)

	for _, q := range queries {
		_, err := tx.ExecContext(ctx, q)

		if err != nil {
			return err
		}
	}

	return nil
}

func splitQueries(queries string) []string {
	split := strings.SplitAfter(queries, ";")
	var trimmedQueries []string

	for _, query := range split {
		trimmed := strings.TrimSpace(query)

		if trimmed == "" {
			continue
		}

		trimmedQueries = append(trimmedQueries, trimmed)
	}

	return trimmedQueries
}
