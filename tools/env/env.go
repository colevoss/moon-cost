package env

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	InvalidEnvLineErr = errors.New("expected `=` followed by optional value")
)

type Env map[string]string

func (e Env) Read(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	l := 0
	for scanner.Scan() {
		l += 1

		line := scanner.Bytes()

		if len(line) == 0 {
			continue
		}

		if line[0] == '#' {
			continue
		}

		key, val, err := parseLine(string(line))

		e[key] = val

		if err != nil {
			return fmt.Errorf("%d: %w", l, err)
		}
	}

	return nil
}

func (e Env) AddEnviron(environ []string) error {
	for _, envVar := range environ {
		key, val, err := parseLine(envVar)

		if err != nil {
			return err
		}

		e[key] = val
	}

	return nil
}

func Load(path string) (Env, error) {
	env := Env{}

	if path == "" {
		path = "./.env"
	}

	file, err := os.Open(path)

	if err != nil {
		return env, err
	}

	defer file.Close()

	if err := env.Read(file); err != nil {
		return env, err
	}

	var errs error

	for k, v := range env {
		err := os.Setenv(k, v)
		errs = errors.Join(errs, err)
	}

	return env, errs
}

func parseLine(line string) (string, string, error) {
	parts := strings.SplitN(line, "=", 2)

	if len(parts) < 2 {
		return "", "", InvalidEnvLineErr
	}

	key := strings.TrimSpace(parts[0])
	val := strings.TrimSpace(parts[1])

	return key, val, nil
}
