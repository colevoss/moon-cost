package moon

import (
	"context"
	"database/sql"
	"moon-cost/moon/assert"
)

const (
	DBTrue  = 1
	DBFalse = 0
)

type Queryer interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Execer interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type QueryerExecer interface {
	Queryer
	Execer
}

type Query struct {
	SQLite string
}

type Queries map[string]Query

func (q Queries) SQLite(name string) string {
	query, ok := q[name]

	assert.Ok(ok, "Expected query %s to exist", name)

	return query.SQLite
}
