package migration

import (
	"context"
	"fmt"
)

func (m *Manager) syncMigrations(ctx context.Context, existing, files []Migration) error {
	lastExistingIndex := 0

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

	fmt.Printf("lastExistingIndex: %d, len(files): %d\n", lastExistingIndex, len(files))

	if len(existing) == len(files) {
		m.logger.Info("No migrations to run")
		return nil
	}

	migrationsToRun := files[lastExistingIndex+1:]

	for _, migration := range migrationsToRun {
		fmt.Printf("Will run %+v\n", migration)
	}

	m.createMigrations(ctx, migrationsToRun)

	return nil
}
