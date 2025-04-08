package main

import (
	"fmt"
	"os"
)

func main() {
	cli := MigrationCLI{}

	if err := cli.Command(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
