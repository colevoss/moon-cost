package migration

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type createCli struct {
	name string
	dir  string
	cli  *MigrationCLI
}

func (c *createCli) init(args []string) error {
	fs := flag.NewFlagSet("create", flag.ExitOnError)

	fs.StringVar(&c.name, "name", "", "")
	fs.StringVar(&c.dir, "dir", "", "")
	c.cli.parseUniversalFlags(fs)

	fs.Parse(args)

	if c.name == "" {
		return fmt.Errorf("Error: name flag required")
	}

	if c.dir == "" {
		return fmt.Errorf("Error: dir flag required")
	}

	return nil
}

func (c *createCli) run() error {
	c.cli.logger.Debug("Checking for dir", "dir", c.dir)

	stat, err := os.Stat(c.dir)

	if err != nil {
		return err
	}

	if !stat.IsDir() {
		return fmt.Errorf("%s is not a dir", c.dir)
	}

	migrationName := strings.ReplaceAll(c.name, " ", "-")
	timestamp := c.cli.now.Now()

	filename := makeMigrationFileName(timestamp, migrationName)

	path := filepath.Join(c.dir, filename)

	c.cli.logger.Debug("Creating migration file", "path", path)

	_, err = os.Create(path)

	if err != nil {
		return fmt.Errorf("Error creating file %w", err)
	}

	c.cli.logger.Info("Created migration file", "path", path)

	return nil
}
