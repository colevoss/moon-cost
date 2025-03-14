package main

import (
	"fmt"
	"moon-cost/migration"
	"os"
)

func main() {
	cli := migration.MigrationCLI{}

	if err := cli.Command(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
