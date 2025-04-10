package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	cli := CurlCLI{}
	ctx := context.Background()

	if err := cli.Command(ctx, os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func _main() {
	ctx := context.Background()

	file, err := os.Open("./curl.json")

	if err != nil {
		panic(err)
	}

	defer file.Close()

	var curl CurlFile

	if err := curl.Read(file); err != nil {
		panic(err)
	}

	env := map[string]string{
		"base": "localhost:8080",
	}

	manager := Manager{
		Curl: curl,
		Env:  env,
	}

	if err := manager.Request(ctx, "error"); err != nil {
		panic(err)
	}
}
