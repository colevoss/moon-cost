package main

import (
	"context"
	"fmt"
)

func (m *Manager) syncMigrations(ctx context.Context, existing, files []Migration) error {
	lastExistingIndex := -1

	if len(existing) > len(files) {
		return fmt.Errorf("Existing migrations count exceeds migration files")
	}

	for i, e := range existing {
		f := files[i]

		if !e.Created.Equal(f.Created) {
			return fmt.Errorf("Existing migration date %s does not match migration file %s", e.Created, f.Created)
		}

		if e.Name != f.Name {
			return fmt.Errorf("Existing migration name %s does not match migration file name %s", e.Name, f.Name)
		}

		lastExistingIndex = i
	}

	migrationsToRun := files[lastExistingIndex+1:]

	if len(migrationsToRun) == 0 {
		m.logger.Info("No new migrations to run")
		return nil
	}

	return m.createMigrations(ctx, migrationsToRun)
}
