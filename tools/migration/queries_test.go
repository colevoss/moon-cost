package migration

import (
	"testing"
)

func TestSplitMigrationCommands(t *testing.T) {

	testSplitQuery := `
CREATE TABLE IF NOT EXIST test (
  id INTEGER PRIMARY KEY,
  name TEXT
);

CREATE TABLE IF NOT EXIST another (
  id INTEGER PRIMARY KEY,
  name TEXT
);

ALTER TABLE test
RENAME TO actually_test;

ALTER TABLE another
ADD COLUMN myColumn TEXT;
`

	split := splitQueries(testSplitQuery)

	if len(split) != 4 {
		t.Errorf("Expected there to be 4 queries. Got %d", len(split))
	}
}
