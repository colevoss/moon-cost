package main

import (
	"fmt"
	"moon-cost/moontest"
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

		moontest.AssertNilError(t, err)

		moontest.Assert(
			t,
			parsed.Name == test.name,
			"Expected parsed name to be %s. Got %s",
			test.name,
			parsed.Name,
		)

		moontest.Assert(
			t,
			parsed.Timestamp.Equal(test.timestamp),
			"Expected parsed timestamp to be %s. Got %s",
			test.timestamp,
			parsed.Timestamp,
		)
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

		moontest.AssertErrorIs(t, err, InvalidFileTypeError)
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

		moontest.AssertErrorIs(t, err, InvalidFilenameError)
	}
}
