package main

import (
	"context"
	"os"
)

func main() {
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

	if err := manager.Request(ctx, ""); err != nil {
		panic(err)
	}
}
