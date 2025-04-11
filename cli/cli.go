package main

import (
	"context"
	"errors"
	"fmt"
)

var (
	InsufficientArgsErr = errors.New("Insufficient args")
)

type Command interface {
	Command(context.Context, []string) error
}

type CLI struct {
	Commands map[string]Command
}

func New() *CLI {
	return &CLI{
		Commands: make(map[string]Command),
	}
}

func (c *CLI) Add(name string, command Command) {
	c.Commands[name] = command
}

func (c *CLI) Command(name string) Command {
	command, ok := c.Commands[name]

	if !ok {
		return nil
	}

	return command
}

// args assumes os.Args[1:]
func (c *CLI) Run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return InsufficientArgsErr
	}

	commandName := args[0]

	command, ok := c.Commands[commandName]

	if !ok {
		return fmt.Errorf("Command %s not found", commandName)
	}

	commandArgs := args[1:]

	return command.Command(ctx, commandArgs)
}
