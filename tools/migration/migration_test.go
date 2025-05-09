package migration

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

const testTimestampInt = 1741883710097

var timestampStr = fmt.Sprintf("%d", testTimestampInt)
var testTimestamp = time.UnixMilli(testTimestampInt)

func TestParseMigrationNameSuccessfully(t *testing.T) {
	names := []struct {
		name      string
		timestamp time.Time
	}{
		{"test", testTimestamp},
		{"test-with-hyphen", testTimestamp},
		{"test_underscore", testTimestamp},
		{"test-with-123-number", testTimestamp},
	}

	for _, test := range names {
		filename := makeMigrationFileName(test.timestamp, test.name)
		parsed, err := parseMigrationName(filename)

		if err != nil {
			t.Errorf("parseMigrationName(%s) = _, %s. want nil err", filename, err)
		}

		if parsed.Name != test.name {
			t.Errorf("parseMigrationName(%s).Name = %s. want %s", filename, parsed.Name, test.name)
		}

		if !parsed.Timestamp.Equal(test.timestamp) {
			t.Errorf("parseMigrationName(%s).Timestamp = %s. want %s", filename, parsed.Timestamp, test.timestamp)
		}
	}
}

func TestParseMigrationNameRequiresSqlFile(t *testing.T) {
	filenames := []string{
		timestampStr + ".test.pdf",
		timestampStr + ".test.sq",
		timestampStr + ".test.ql",
		timestampStr + ".test.txt",
		timestampStr + ".test.somethingelse",
	}

	for _, test := range filenames {
		_, err := parseMigrationName(test)

		if !errors.Is(err, InvalidFileTypeError) {
			t.Errorf("parseMigrationName() = _, %s. want %s", err, InvalidFileTypeError)
		}
	}
}

func TestParsedMalformedFileName(t *testing.T) {
	filenames := []string{
		"justone.sql",
		"3.file.parts.sql",
		"not-timestamp.name.sql",
	}

	for _, test := range filenames {
		_, err := parseMigrationName(test)

		if !errors.Is(err, InvalidFilenameError) {
			t.Errorf("parseMigrationName() = _, %s. want %s", err, InvalidFilenameError)
		}
	}
}
