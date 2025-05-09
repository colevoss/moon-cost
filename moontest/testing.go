package moontest

import (
	"os"
	"path/filepath"
	"testing"
)

func LoadTestFixture(t *testing.T, file string) *os.File {
	t.Helper()

	path := filepath.Join(".", "test-fixtures", file)

	testFixture, err := os.Open(path)

	if err != nil {
		t.Fatalf("Could not load test fixture: %s", err)
	}

	return testFixture
}
