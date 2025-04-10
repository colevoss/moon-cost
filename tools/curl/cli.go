package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
)

type CurlCLI struct {
	verbose bool
	logger  *slog.Logger
	File    io.Reader
}

type CLIArgs struct {
	File    string // curl file to load
	Request string // name of request to run
	Verbose bool
}

const (
	DefaultRequestName     = "default"
	RequestFlagDescription = "Name of request in file to send"

	VerboseFlagDescription = "Enables verbose logging"
)

// Its assumed that args is os.Args
// args[1] should be the filename
// args[2:] should be flags
func (c *CLIArgs) Parse(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Not enough arguments")
	}

	file := args[1]
	c.File = file

	fs := flag.NewFlagSet("curl", flag.ExitOnError)

	fs.StringVar(&c.Request, "request", DefaultRequestName, RequestFlagDescription)
	fs.StringVar(&c.Request, "r", DefaultRequestName, RequestFlagDescription+" (shorthand)")

	fs.BoolVar(&c.Verbose, "verbose", false, VerboseFlagDescription)
	fs.BoolVar(&c.Verbose, "v", false, VerboseFlagDescription+" (shorthand)")

	if err := fs.Parse(args[2:]); err != nil {
		return err
	}

	return nil
}

func (c *CurlCLI) Command(ctx context.Context, args []string) error {
	var cliArgs CLIArgs

	if err := cliArgs.Parse(args); err != nil {
		return err
	}

	file, err := os.Open(cliArgs.File)

	if err != nil {
		return err
	}

	var curl CurlFile
	if err := curl.Read(file); err != nil {
		return err
	}

	env := map[string]string{
		"base": "localhost:8080",
	}

	manager := Manager{
		Curl: curl,
		Env:  env,
	}

	return manager.Request(ctx, cliArgs.Request)
}
