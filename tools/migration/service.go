package main

type MigrationService interface {
	ensureMigrationsTable() error
}
