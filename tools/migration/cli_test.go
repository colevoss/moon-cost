package migration

import (
	"context"
	"moon-cost/common"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateCLI(t *testing.T) {
	ctx := context.Background()

	time := time.Now()
	now := common.TestNow{Time: time}
	dir := t.TempDir()

	cli := MigrationCLI{now: now, suppress: true}

	tests := []struct {
		args          []string
		migrationName string
	}{
		{args: []string{"create", "--dir", dir, "--name", "test-migration"}, migrationName: "test-migration"},
		{args: []string{"create", "--dir", dir, "--name", "test migration with spaces"}, migrationName: "test-migration-with-spaces"},
	}

	for _, test := range tests {
		err := cli.Command(ctx, test.args)

		if err != nil {
			t.Errorf("cli.Command() = %s. want nil", err)
		}

		fileName := makeMigrationFileName(time, test.migrationName)
		path := filepath.Join(dir, fileName)

		stat, err := os.Stat(path)

		if err != nil {
			t.Errorf("os.State = _, %s. want nil", err)
		}

		statName := stat.Name()

		if statName != fileName {
			t.Errorf("created filename is %s. want %s", statName, fileName)
		}
	}
}

func TestCreateCliErrors(t *testing.T) {
	ctx := context.Background()
	f, _ := os.CreateTemp("", "test-file")
	dir := t.TempDir()

	defer os.Remove(f.Name())

	cli := MigrationCLI{suppress: true}

	tests := []struct {
		args []string
		test string
	}{
		{test: "dir flag is file not dir", args: []string{"create", "--dir", f.Name(), "--name", "test-name"}},
		{test: "dir flag value does not exist", args: []string{"create", "--dir", "non-existent", "--name", "test-name"}},
		{test: "no name flag", args: []string{"create", "--dir", dir}},
		{test: "no dir flag", args: []string{"create", "--name", "name"}},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			err := cli.Command(ctx, test.args)

			if err == nil {
				t.Errorf("Expected command to error. Got nil")
			}
		})
	}
}
