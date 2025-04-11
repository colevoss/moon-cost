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

	curl := curl.CurlCLI{}
	migration := migration.MigrationCLI{}

	cli := New()
	cli.Add("curl", &curl)
	cli.Add("migration", &migration)

	args := os.Args[1:]

	if err := cli.Run(ctx, args); err != nil {
		log.Fatal(err)
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
