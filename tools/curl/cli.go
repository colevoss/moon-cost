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
	logger  *slog.Logger
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
}

const (
	DefaultRequestName     = "default"
	RequestFlagDescription = "Name of request in file to send"

	VerboseFlagDescription  = "Enables verbose logging"
	DebugFlagDescription    = "Enables debug logging"
	SuppressFlagDescription = "Suppresses all logs"
	JSONFlagDescription     = "Outputs all logs as json"
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
	c.logger = initLogger(cliArgs)

	return nil
}

func initLogger(args CLIArgs) *slog.Logger {
	var handler slog.Handler
	var out io.Writer

	level := slog.LevelError
	enabled := true

	if args.Suppress {
		enabled = false
	}

	// This should imply that debug takes precedence
	if args.Verbose {
		level = slog.LevelInfo
	}
	if args.Debug {
		level = slog.LevelDebug
	}

	out = os.Stdout

	if args.JSON {
		handler = NewJSONHandler(out, level, enabled)
	} else {
		handler = NewStandardHandler(out, level, enabled)
	}

	return slog.New(handler)
}

func (c *CurlCLI) Command(ctx context.Context, args []string) error {
	if err := c.Init(args); err != nil {
		return err
	}

	c.logger.Debug("Opening curl file", slog.String("file", c.Args.File))

	file, err := os.Open(c.Args.File)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	defer file.Close()

	var curl Curl
	if err := curl.Read(file); err != nil {
		c.logger.Error(err.Error(), slog.String("file", c.Args.File))
		return err
	}

	c.logger.Debug("Parsed curl file", slog.String("file", c.Args.File))

	env, err := c.LoadEnv(c.Args)

	if err != nil {
		return err
	}

	request, ok := curl.Req(c.Args.Request)

	if !ok {
		c.logger.Error(ErrRequestNotFound.Error(), "request", c.Args.Request)
		return ErrRequestNotFound
	}

	manager := Manager{
		Curl:    curl,
		Request: request,
		Env:     env,
	}

	c.logger.Debug("Building request")

	req, err := BuildRequest(ctx, manager)

	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	c.logRequest(req)
	start := time.Now()
	res, err := http.DefaultClient.Do(req)
	end := time.Now()
	duration := end.Sub(start)

	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	if err := c.logResponse(ctx, res, request, duration); err != nil {
		c.logger.Error(err.Error())
		return err
	}

	return CheckResponse(request, res)
}

func (c *CurlCLI) LoadEnv(args CLIArgs) (Env, error) {
	env := Env{}

	if args.Env == "" {
		return env, nil
	}

	c.logger.Debug("Opening env file", slog.String("envFile", args.Env))
	file, err := os.Open(args.Env)

	if err != nil {
		c.logger.Error(err.Error(), slog.String("envFile", args.Env))
		return env, err
	}

	defer file.Close()

	c.logger.Debug("Parsing env file", slog.String("envFile", args.Env))

	if err := env.Read(file); err != nil {
		c.logger.Error(fmt.Sprintf("%s:%s", args.Env, err))
	}

	c.logger.Debug("Parsed env file", slog.String("envFile", args.Env))

	return env, nil
}

func (c *CurlCLI) logRequest(req *http.Request) {
	url := req.URL.String()
	method := req.Method

	c.logger.Info(fmt.Sprintf("%s %s", method, url))
}

func (c *CurlCLI) logResponse(ctx context.Context, res *http.Response, request Request, duration time.Duration) error {
	level := slog.LevelInfo
	url := res.Request.URL.String()
	method := res.Request.Method
	// headers := res.Header
	status := res.StatusCode

	attrs := []any{
		slog.String("request", c.Args.Request),
		slog.Bool("response", true),
		slog.Int("status", status),
		slog.String("url", url),
		slog.String("method", method),
		slog.Int64("durationMs", duration.Milliseconds()),
	}

	// if len(headers) > 0 {
	// 	attrs = append(attrs, slog.Any("headers", headers))
	// }

	body, err := c.getBodyData(res)

	if err != nil {
		return err
	}

	attrs = append(attrs, slog.Any("body", body))

	var b strings.Builder

	b.WriteString(fmt.Sprintf("(%s) %s", duration.Truncate(time.Millisecond).String(), res.Status))

	if request.Expect.Status != 0 {
		expected := request.Expect.Status
		pass := res.StatusCode == expected
		attrs = append(attrs, slog.Bool("pass", pass))
		attrs = append(attrs, slog.Int("expected", expected))

		if !pass {
			level = slog.LevelWarn
			b.WriteString(fmt.Sprintf(". Want %d %s", request.Expect.Status, http.StatusText(request.Expect.Status)))
		}
	}

	c.logger.Log(
		ctx,
		level,
		b.String(),
		attrs...,
	)

	return nil
}

func (c *CurlCLI) getBodyData(res *http.Response) (any, error) {
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	if res.Header.Get("Content-Type") != "application/json" {
		return string(data), nil
	}

	var body map[string]any

	err = json.Unmarshal(data, &body)

	if err != nil {
		return "", err
	}

	return body, nil
}
