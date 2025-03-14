package migration

import (
	"os"
	"path/filepath"
	"sort"
)

func (m *Manager) inspectDir() ([]Migration, error) {
	m.logger.Debug("Inspecting dir for migration files", "dir", m.Dir)

	entries, err := os.ReadDir(m.Dir)

	if err != nil {
		return nil, err
	}

	var migrations []Migration

	for _, entry := range entries {
		migration, ok, err := m.inspectFile(entry)

		if err != nil {
			return migrations, err
		}

		if !ok {
			continue
		}

		migrations = append(migrations, migration)
	}

	sort.Sort(MigrationByCreated(migrations))

	return migrations, nil
}

func (m *Manager) inspectFile(entry os.DirEntry) (Migration, bool, error) {
	var migration Migration
	filename := entry.Name()

	if entry.IsDir() {
		m.logger.Warn("Ignoring directory", "dir", filename)
		return migration, false, nil
	}

	parsed, err := parseMigrationName(filename)

	if err != nil {
		return migration, false, err
	}

	migration.Created = parsed.Timestamp
	migration.Name = parsed.Name
	migration.Filename = filename

	m.logger.Debug("Parsed migration details from filename", "timestamp", migration.Created, "name", migration.Name)

	contents, err := m.readFile(entry)

	if err != nil {
		return migration, false, err
	}

	migration.Instruction = contents

	return migration, true, nil
}

func (m *Manager) readFile(entry os.DirEntry) (string, error) {
	filename := entry.Name()
	file, err := os.Open(filepath.Join(m.Dir, entry.Name()))

	if err != nil {
		return "", err
	}

	info, err := entry.Info()

	if err != nil {
		return "", err
	}

	size := info.Size()
	bytes := make([]byte, size)

	m.logger.Debug("Reading migration file", "file", filename)

	n, err := file.Read(bytes)

	if err != nil {
		return "", err
	}

	m.logger.Debug("Read migration file", "file", filename, "bytes", n)

	return string(bytes), nil
}
