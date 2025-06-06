package migration

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"moon-cost/common"
)

type MigrationCLI struct {
	verbose  bool
	suppress bool
	logger   *slog.Logger
	now      common.Now
}

func (mcli *MigrationCLI) init() {
	level := slog.LevelInfo

	if mcli.verbose {
		level = slog.LevelDebug
	}

	if mcli.suppress {
		level = slog.LevelError
	}

	slog.SetLogLoggerLevel(level)

	mcli.logger = slog.Default()

	if mcli.now == nil {
		mcli.now = common.TimeNow{}
	}
}

func (mcli *MigrationCLI) Command(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No command provided.")
	}

	command := args[0]
	flags := args[1:]

	switch command {
	case "create":
		return mcli.Create(flags)

	case "run":
		return mcli.Run(ctx, flags)

	default:
		return fmt.Errorf("Invalid command: %s", command)
	}
}

func (mcli *MigrationCLI) Create(args []string) error {
	createCli := createCli{cli: mcli}

	if err := createCli.init(args); err != nil {
		return err
	}

	mcli.init()

	return createCli.run()
}

func (mcli *MigrationCLI) Run(ctx context.Context, args []string) error {
	runCli := runCli{cli: mcli}

	if err := runCli.Init(args); err != nil {
		return err
	}

	mcli.init()

	return runCli.Command(ctx)
}

func (mcli *MigrationCLI) parseUniversalFlags(fs *flag.FlagSet) {
	fs.BoolVar(&mcli.verbose, "v", false, "Verbose")

	// allows user to set suppress directly set suppress on struct
	if !mcli.suppress {
		fs.BoolVar(&mcli.suppress, "s", false, "Suppress")
	}
}
