package migration

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"moon-cost/moon/assert"
)

const DEFAULT_TABLE_NAME = "migrations"

type Manager struct {
	Dir   string
	Table string
	DB    *sql.DB

	logger *slog.Logger
}

type MigrationOption func(m *Manager)

func (m *Manager) Init(options ...MigrationOption) {
	for _, option := range options {
		option(m)
	}

	m.init()
}

func (m *Manager) Run(ctx context.Context) error {
	assert.Ensure(m.logger, "Manager logger is nil")
	assert.Ensure(m.DB, "Manager db is nil")
	assert.Ok(m.Table != "", "Manager Table is blank")

	if err := m.ensureMigrationsTable(ctx); err != nil {
		return err
	}

	existingMigrations, err := m.getAllMigrations(ctx)

	if err != nil {
		return err
	}

	migrationFiles, err := m.inspectDir()

	if err != nil {
		return err
	}

	if err := m.syncMigrations(ctx, existingMigrations, migrationFiles); err != nil {
		return err
	}

	return nil
}

func WithLogger(logger *slog.Logger) MigrationOption {
	return func(m *Manager) {
		m.logger = logger
	}
}

func (m *Manager) init() {
	m.ensureLogger()
	m.ensureTableName()
}

func (m *Manager) ensureLogger() {
	if m.logger != nil {
		return
	}

	m.logger = slog.Default()
}

func (m *Manager) ensureTableName() {
	if m.Table != "" {
		return
	}

	m.Table = DEFAULT_TABLE_NAME
}

func (m *Manager) formatQuery(query string) string {
	return fmt.Sprintf(query, m.Table)
}
