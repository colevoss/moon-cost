package curl

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

type CurlCLI struct {
	verbose bool
	File    io.Reader
	EnvFile Env

	Args CLIArgs
}

type CLIArgs struct {
	File     string // curl file to load
	Request  string // name of request to run
	Env      string
	Verbose  bool
	Debug    bool
	Suppress bool
	JSON     bool
	Pretty   bool
}

const (
	DefaultRequestName     = "default"
	RequestFlagDescription = "Name of request in file to send"

	VerboseFlagDescription  = "Enables verbose logging"
	DebugFlagDescription    = "Enables debug logging"
	SuppressFlagDescription = "Suppresses all logs"
	JSONFlagDescription     = "Outputs all logs as json"
	PrettyFlagDescription   = "Pretty prints the json (only used in json mode)"
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

	fs.BoolVar(&c.Pretty, "pretty", false, PrettyFlagDescription)
	fs.BoolVar(&c.Pretty, "p", false, PrettyFlagDescription+"(shorthand)")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	return nil
}

func (c *CurlCLI) Init(args []string) error {
	var cliArgs CLIArgs

	if err := cliArgs.Parse(args); err != nil {
		return err
	}

	c.Args = cliArgs
	initLogger(cliArgs)

	return nil
}

func initLogger(args CLIArgs) {
	var handler slog.Handler
	// var out io.Writer
	out := os.Stdout

	level := slog.LevelInfo

	// This should imply that debug takes precedence
	if args.Verbose {
		level = LevelVerbose
	}

	if args.Debug {
		level = slog.LevelDebug
	}

	if args.Suppress {
		handler = slog.DiscardHandler
	} else {
		if args.JSON {
			level = slog.LevelError
			handler = NewJSONHandler(out, level)
		} else {
			handler = NewStandardHandler(out, level)
		}
	}

	slog.SetDefault(slog.New(handler))
}

func (c *CurlCLI) Command(ctx context.Context, args []string) error {
	if err := c.Init(args); err != nil {
		return err
	}

	slog.Debug("Opening curl file", slog.String("file", c.Args.File))

	file, err := os.Open(c.Args.File)
	if err != nil {
		slog.Error(err.Error())
		return err
	}
	defer file.Close()

	var curl Curl
	if err := curl.Read(file); err != nil {
		slog.Error(err.Error(), slog.String("file", c.Args.File))
		return err
	}

	slog.Debug("Parsed curl file", slog.String("file", c.Args.File))

	env, err := c.LoadEnv(ctx, c.Args)

	if err != nil {
		return err
	}

	request, ok := curl.Req(c.Args.Request)

	if !ok {
		slog.Error(ErrRequestNotFound.Error(), "request", c.Args.Request)
		return ErrRequestNotFound
	}

	manager := Manager{
		Curl:    curl,
		Request: request,
		Env:     env,
	}

	slog.Debug("Building request")

	req, err := BuildRequest(ctx, manager)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	c.logRequest(req)
	start := time.Now()
	res, err := http.DefaultClient.Do(req)
	end := time.Now()
	duration := end.Sub(start)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	result, err := NewResult(c.Args.Request, res, request, duration)

	if err != nil {
		slog.Error(err.Error())
		return err
	}

	if err := c.logResult(ctx, result); err != nil {
		slog.Error(err.Error())
		return err
	}

	return CheckResponse(request, res)
}

func (c *CurlCLI) LoadEnv(ctx context.Context, args CLIArgs) (Env, error) {
	env := Env{}

	if args.Env == "" {
		return env, nil
	}

	slog.Debug("Opening env file", slog.String("envFile", args.Env))
	file, err := os.Open(args.Env)

	if err != nil {
		slog.Error(err.Error(), slog.String("envFile", args.Env))
		return env, err
	}

	defer file.Close()

	slog.Debug("Parsing env file", slog.String("envFile", args.Env))

	if err := env.Read(file); err != nil {
		slog.Error(fmt.Sprintf("%s:%s", args.Env, err))
	}

	slog.Log(ctx, LevelVerbose, "Parsed env file", slog.String("envFile", args.Env))

	return env, nil
}

func (c *CurlCLI) logRequest(req *http.Request) {
	url := req.URL.String()
	method := req.Method

	slog.Info(fmt.Sprintf("%s %s", method, url))
}

func (c *CurlCLI) logResult(ctx context.Context, result Result) error {
	// TODO: support pretty printing
	if c.Args.JSON {
		var data []byte
		var err error

		if c.Args.Pretty {
			data, err = json.MarshalIndent(result, "", "  ")
		} else {
			data, err = json.Marshal(result)
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

	b.WriteString(fmt.Sprintf("(%s) %s", result.Duration.Truncate(time.Millisecond).String(), result.Response.Status))

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
			b.WriteString(fmt.Sprintf(". Want %d %s", result.Expected, http.StatusText(result.Expected)))
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
