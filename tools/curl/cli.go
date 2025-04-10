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
	EnvFile Env
}

type CLIArgs struct {
	File    string // curl file to load
	Request string // name of request to run
	Verbose bool
	Env     string
}

const (
	DefaultRequestName     = "default"
	RequestFlagDescription = "Name of request in file to send"

	VerboseFlagDescription = "Enables verbose logging"
	EnvFlagDescription     = "Path to env file to load"
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

	fs.StringVar(&c.Env, "env", "", EnvFlagDescription)
	fs.StringVar(&c.Env, "e", "", EnvFlagDescription)

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
	defer file.Close()

	var curl CurlFile
	if err := curl.Read(file); err != nil {
		return err
	}

	env, err := c.LoadEnv(cliArgs)

	if err != nil {
		return err
	}

	manager := Manager{
		Curl: curl,
		Env:  env,
	}

	return manager.Request(ctx, cliArgs.Request)
}

func (c *CurlCLI) LoadEnv(args CLIArgs) (Env, error) {
	env := NewEnv()

	if args.Env == "" {
		return env, nil
	}

	file, err := os.Open(args.Env)

	if err != nil {
		return env, err
	}

	defer file.Close()

	env.Read(file)

	return env, nil
}
