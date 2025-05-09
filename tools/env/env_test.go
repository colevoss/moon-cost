package env

import (
	"moon-cost/moontest"
	"os"
	"testing"
)

func TestEnvparseLine(t *testing.T) {
	tests := []struct {
		line string
		key  string
		val  string
	}{
		{"foo=", "foo", ""},
		{"foo=bar", "foo", "bar"},
		{"foo = bar", "foo", "bar"},
		{"foo =bar", "foo", "bar"},
		{"foo=bar=baz", "foo", "bar=baz"},
		{"foo= bar=baz", "foo", "bar=baz"},
	}

	for _, test := range tests {
		t.Run(test.line, func(t *testing.T) {
			key, val, err := parseLine(test.line)

			if err != nil {
				t.Fatalf("parseLine(\"%s\") = %s. want nil", test.line, err)
			}

			if key != test.key {
				t.Errorf("parseLine(\"%s\") = %s. want %s", test.line, key, test.key)
			}

			if val != test.val {
				t.Errorf("parseLine(\"%s\") = %s. want %s", test.line, val, test.val)
			}
		})
	}

	t.Run("error", func(t *testing.T) {
	})
}

func TestEnvparseLineErrs(t *testing.T) {
	tests := []string{
		"",
		"key",
		"key val",
		"key - val",
	}

	for _, test := range tests {
		_, _, err := parseLine(test)

		if err == nil {
			t.Error("Got nil. Expected error")
		}
	}
}

func TestEnvAddEnviron(t *testing.T) {
	env := Env{}

	environs := []string{
		"foo=bar",
		"path=./some/path",
		"hello=world:",
	}

	env.AddEnviron(environs)

	expected := map[string]string{
		"foo":   "bar",
		"path":  "./some/path",
		"hello": "world:",
	}

	if len(environs) != len(expected) {
		t.Errorf("len(env) = %d. want %d", len(env), len(environs))
	}

	for k, v := range expected {
		val, _ := env[k]

		if val != v {
			t.Errorf("env var %s = \"%s\". want \"%s\"", k, val, v)
		}
	}
}

func TestEnvReadSuccess(t *testing.T) {
	file := moontest.LoadTestFixture(t, "env.good")

	defer file.Close()

	env := Env{}

	if err := env.Read(file); err != nil {
		t.Errorf("env.Read(./test/env.good) = %s. wanted nil", err)
	}

	data := map[string]string{
		"FOO":          "BAR",
		"MY_VAR":       "VALUE",
		"MY_OTHER_VAR": "OTHER VALUE",
		"EMPTY_VAR":    "",
	}

	// if len(data) != len(env.Data) {
	if len(data) != len(env) {
		t.Errorf("len(env) = %d. want %d", len(env), len(data))
	}

	for k, v := range data {
		envVar, _ := env[k]

		if envVar != v {
			t.Errorf("env[%s] = %s. want %s", k, envVar, v)
		}
	}
}

func TestEnvReadError(t *testing.T) {
	file := moontest.LoadTestFixture(t, "env.bad")

	defer file.Close()

	env := Env{}

	err := env.Read(file)

	if err == nil {
		t.Error("env.Read() = nil. want err")
	}

	expected := "4: expected `=` followed by optional value"

	if err.Error() != expected {
		t.Errorf("env.Read() = %s. want %s", err, expected)
	}
}

func TestEnvLoad(t *testing.T) {
	env, err := Load("./test-fixtures/test-env-load.env")

	if err != nil {
		t.Errorf("Load() = _, %s. want nil error", err)
	}

	tests := []struct {
		name     string
		expected string
	}{
		{"TEST_VAR", "test-var"},
		{"OTHER_VAR", "other-var"},
	}

	for _, test := range tests {
		envVar := os.Getenv(test.name)

		if envVar != test.expected {
			t.Errorf("environment variable %s = %s. want %s", test.name, envVar, test.expected)
		}
	}

	t.Cleanup(func() {
		for k := range env {
			os.Unsetenv(k)
		}
	})
}
