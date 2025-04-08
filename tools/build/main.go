package main

import (
	"context"
	"log"
	"log/slog"
	"os"
)

// TODO: Ensure build directory exists

func main() {
	defer log.SetFlags(log.Flags())
	log.SetFlags(0)

	slog.SetLogLoggerLevel(slog.LevelWarn)
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	ctx := context.Background()
	file, err := os.Open("./build.json")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	pwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	cfg, err := ConfigFrom(pwd, file, os.Stdout)

	if err != nil {
		panic(err)
	}

	if err := cfg.Run(ctx); err != nil {
		// TODO: update this to work with json output
		log.Printf("Build errors:\n")
		log.Fatal(err)
	}
}
