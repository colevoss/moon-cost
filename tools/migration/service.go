package migration

type MigrationService interface {
	ensureMigrationsTable() error
}
