package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Migration struct {
	Id          int
	Created     time.Time
	Name        string
	Filename    string
	Instruction string
}

type MigrationByCreated []Migration

func (m MigrationByCreated) Len() int      { return len(m) }
func (m MigrationByCreated) Swap(i, j int) { m[i], m[j] = m[j], m[i] }
func (m MigrationByCreated) Less(i, j int) bool {
	return m[i].Created.Before(m[j].Created)
}

type parsedMigrationFilename struct {
	Timestamp time.Time
	Name      string
}

var InvalidFilenameError = errors.New("Invalid file name")
var InvalidFileTypeError = errors.New("Invalid file type")

func parseMigrationName(filename string) (parsedMigrationFilename, error) {
	var parsed parsedMigrationFilename

	ext := filepath.Ext(filename)

	if ext != ".sql" {
		return parsed, fmt.Errorf("%w. %s should be .sql file", InvalidFileTypeError, filename)
	}

	base := strings.TrimSuffix(filename, ext)
	parts := strings.Split(base, ".")

	if len(parts) != 2 {
		return parsed, InvalidFilenameError
	}

	timestampPart := parts[0]
	namePart := parts[1]

	timestampInt, err := strconv.Atoi(timestampPart)

	if err != nil {
		return parsed, fmt.Errorf("%w. Invalid timestamp section in %s: %s", InvalidFilenameError, filename, timestampPart)
	}

	timestamp := time.UnixMilli(int64(timestampInt))

	parsed.Timestamp = timestamp
	parsed.Name = namePart

	return parsed, nil
}

func makeMigrationFileName(timestamp time.Time, name string) string {
	return fmt.Sprintf("%d.%s.sql", timestamp.UnixMilli(), name)
}
