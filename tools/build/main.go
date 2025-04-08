package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()

	slog.SetLogLoggerLevel(slog.LevelDebug)
	file, err := os.Open("./build.json")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	pwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	cfg, err := ConfigFrom(pwd, file)

	if err != nil {
		panic(err)
	}

	if err := cfg.Run(ctx); err != nil {
		fmt.Printf("ERR: %v\n", err)
	}
}
