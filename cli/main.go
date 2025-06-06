package main

import (
	"context"
	"log"
	"moon-cost/tools/curl"
	"moon-cost/tools/migration"
	"os"
	"os/signal"
)

func run(ctx context.Context) {
	log.SetFlags(0)

	var curl curl.CurlCLI
	curl.Out = os.Stdout
	var migration migration.MigrationCLI

	cli := New()
	cli.Add("curl", &curl)
	cli.Add("migration", &migration)

	args := os.Args[1:]

	if err := cli.Run(ctx, args); err != nil {
		os.Exit(1)
	}
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-c
		cancel()
	}()

	run(ctx)
}
