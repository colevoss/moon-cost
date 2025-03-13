package db

import (
	"database/sql"
	_ "github.com/tursodatabase/go-libsql"
)

type SQLite struct {
	URL string
	db  *sql.DB
}

func (s *SQLite) DB() *sql.DB {
	return s.db
}

func (s *SQLite) Open() error {
	db, err := sql.Open("libsql", s.URL)

	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *SQLite) Close() error {
	return s.db.Close()
}
