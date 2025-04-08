package main

import (
	"context"
	"os/exec"
	"strings"
)

type buildCommand []string

func goBuild(out, src string) buildCommand {
	return []string{"go", "build", "-json", "-o", out, src}
}

func (bc buildCommand) String() string {
	return strings.Join(bc, " ")
}

func (bc buildCommand) exec(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, bc[0], bc[1:]...)
}
