package curl

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"moon-cost/assert"
	"moon-cost/logging"
	"moon-cost/tools/env"
	"net/http"
	"strings"
	"time"
)

type CurlCLI struct {
	verbose bool
	File    io.Reader
	EnvFile env.Env

	Args CLIArgs
	Out  io.Writer
}

type CLIArgs struct {
	File     string // curl file to load
	Request  string // name of request to run
	Env      string
	Verbose  bool
	Debug    bool
	Suppress bool
	JSON     bool
	Raw      bool
}

const (
	DefaultRequestName     = "default"
	RequestFlagDescription = "Name of request in file to send"

	VerboseFlagDescription  = "Enables verbose logging"
	DebugFlagDescription    = "Enables debug logging"
	SuppressFlagDescription = "Suppresses all logs"
	JSONFlagDescription     = "Outputs all logs as json"
	RawFlagDescription      = "Prints JSON without indentation (only used in json mode)"
	EnvFlagDescription      = "Path to env file to load"
)

// Its assumed that args is os.Args[2:]
// args[0] should be the filename
// args[1:] should be flags
func (c *CLIArgs) Parse(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Not enough arguments")
	}

	file := args[0]
	c.File = file

	fs := flag.NewFlagSet("curl", flag.ExitOnError)

	fs.StringVar(&c.Request, "request", DefaultRequestName, RequestFlagDescription)
	fs.StringVar(&c.Request, "r", DefaultRequestName, RequestFlagDescription+" (shorthand)")

	fs.StringVar(&c.Env, "env", "", EnvFlagDescription)
	fs.StringVar(&c.Env, "e", "", EnvFlagDescription)

	fs.BoolVar(&c.Verbose, "verbose", false, VerboseFlagDescription)
	fs.BoolVar(&c.Verbose, "v", false, VerboseFlagDescription+" (shorthand)")

	fs.BoolVar(&c.Debug, "debug", false, DebugFlagDescription)
	fs.BoolVar(&c.Debug, "d", false, DebugFlagDescription+" (shorthand)")

	fs.BoolVar(&c.Suppress, "suppress", false, SuppressFlagDescription)
	fs.BoolVar(&c.Suppress, "s", false, SuppressFlagDescription+" (shorthand)")

	fs.BoolVar(&c.JSON, "json", false, JSONFlagDescription)
	fs.BoolVar(&c.JSON, "j", false, JSONFlagDescription+" (shorthand)")

	fs.BoolVar(&c.Raw, "raw", false, RawFlagDescription)

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	return nil
}

func (c *CurlCLI) Init(args []string) error {
	assert.Ensure(c.Out, "CurlCLI out required")

	var cliArgs CLIArgs

	if err := cliArgs.Parse(args); err != nil {
		return err
	}

	c.Args = cliArgs
	initLogger(cliArgs, c.Out)

	return nil
}

func initLogger(args CLIArgs, out io.Writer) {
	var handler slog.Handler

	level := slog.LevelInfo

	if args.JSON {
		level = slog.LevelError
	}

	// This should imply that debug takes precedence
	if args.Verbose {
		level = logging.LevelVerbose
	}

	if args.Debug {
		level = slog.LevelDebug
	}

	if args.Suppress {
		handler = slog.DiscardHandler
	} else {
		if args.JSON {
			handler = slog.NewJSONHandler(out, nil)
		} else {
			handler = logging.NewPrettyHandler(out, level)
		}
	}

	slog.SetDefault(slog.New(handler))
}

func (c *CurlCLI) Command(ctx context.Context, args []string) error {
	if err := c.Init(args); err != nil {
		return err
	}

	var client Client

	slog.Debug("loading curl file", slog.String("file", c.Args.File))

	if err := client.LoadCurl(c.Args.File); err != nil {
		slog.Error(err.Error())
		return err
	}

	slog.Log(ctx, logging.LevelVerbose, "curl file loaded", slog.String("file", c.Args.File))

	if c.Args.Env != "" {
		slog.Debug("loading env file", slog.String("file", c.Args.Env))

		if err := client.LoadEnv(c.Args.Env); err != nil {
			slog.Error(err.Error())
			return err
		}

		slog.Log(ctx, logging.LevelVerbose, "env file loaded", slog.String("file", c.Args.Env))
	}

	if err := client.Use(c.Args.Request); err != nil {
		slog.Error(err.Error())
		return err
	}

	result, err := client.Execute(ctx)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if err := c.logResult(ctx, result); err != nil {
		slog.Error(err.Error())
		return err
	}

	return nil
}

func (c *CurlCLI) logResult(ctx context.Context, result Result) error {
	// TODO: support pretty printing
	if c.Args.JSON {
		var data []byte
		var err error

		if c.Args.Raw {
			data, err = json.Marshal(result)
		} else {
			data, err = json.MarshalIndent(result, "", "  ")
		}

		if err != nil {
			return err
		}

		fmt.Println(string(data))
		return nil
	}

	return c.logResponse(ctx, result)
}

func (c *CurlCLI) logResponse(ctx context.Context, result Result) error {
	level := slog.LevelInfo

	attrs := []any{
		slog.String("request", result.Name),
		slog.Int("status", result.Status),
		slog.String("url", result.URL),
		slog.String("method", result.Method),
		slog.Int64("durationMs", result.DurationMS),
	}

	body, err := result.BodyData()

	if err != nil {
		return err
	}

	attrs = append(attrs, slog.Any("body", body))

	var b strings.Builder

	b.WriteString(fmt.Sprintf("(%s) %s", result.Duration.Truncate(time.Millisecond*1).String(), result.Response.Status))

	if result.Expected == 0 {
		if result.Response.StatusCode >= 400 {
			level = slog.LevelWarn
		}
	} else {
		expected := result.Expected
		pass := result.Response.StatusCode == expected
		attrs = append(attrs, slog.Bool("pass", pass))
		attrs = append(attrs, slog.Int("expected", expected))

		if !pass {
			level = slog.LevelWarn
			b.WriteString(fmt.Sprintf(". want %d %s", result.Expected, http.StatusText(result.Expected)))
		}
	}

	slog.Log(
		ctx,
		level,
		b.String(),
		attrs...,
	)

	return nil
}
