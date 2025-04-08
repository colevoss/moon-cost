package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os/exec"
	"path/filepath"
	"slices"
)

type Config struct {
	base   string
	logOut io.Writer
	Src    []string `json:"src"`
	Out    string   `json:"out"`
	Ignore []string `json:"ignore"`
}

type buildSrc struct {
	name     string
	fullPath string
	relPath  string
	srcDir   string
}

func ConfigFrom(base string, r io.Reader, out io.Writer) (Config, error) {
	var config Config

	if err := json.NewDecoder(r).Decode(&config); err != nil {
		return config, err
	}

	config.base = base
	config.logOut = out

	return config, nil
}

func (c *Config) Run(ctx context.Context) error {
	srcs, err := c.srcs()

	if err != nil {
		return err
	}

	var errs error

	for _, src := range srcs {
		if slices.Contains(c.Ignore, src.name) {
			slog.Debug(
				"Skipping build target",
				"module", src.name,
				"path", src.relPath,
				"srcDir", src.srcDir,
			)

			continue
		}

		err := c.build(ctx, src)

		if err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}

func (c *Config) build(ctx context.Context, build buildSrc) error {
	output := c.buildOut(build)
	buildCommand := goBuild(output, build.fullPath)

	cmd := buildCommand.exec(ctx)

	slog.Debug(
		"Building module",
		"src", build.relPath,
		"out", output,
		"cmd", buildCommand,
	)

	out, err := cmd.StdoutPipe()

	defer out.Close()

	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	slog.Debug(
		"Started build command",
		"module", build.name,
		"pid", cmd.Process.Pid,
	)

	fmt.Fprintf(c.logOut, "Building %s\n", build.name)

	slog.Info(
		"Building",
		"name", build.name,
		"src", build.relPath,
		"out", output,
	)

	events := ReadEvents(out)

	if err = cmd.Wait(); err != nil {
		if err, ok := err.(*exec.ExitError); !ok {
			fmt.Printf("err: %v\n", err)
			return err
		}
	}

	if cmd.ProcessState.Success() {
		fmt.Fprintf(c.logOut, "Success %s -> %s\n", build.relPath, output)

		slog.Info(
			"Success",
			"name", build.name,
			"src", build.relPath,
			"out", output,
		)

		return nil
	}

	slog.Warn(
		"Error running build command",
		"pid", cmd.Process.Pid,
		"err", err,
	)

	var errs error
	for _, event := range events {
		eventStr := event.String()

		if eventStr == "" {
			continue
		}

		fmt.Printf("\t[%s] %s\n", event.ImportPath, event)
		errs = errors.Join(errs, event)
	}

	return errs
}

func (c *Config) buildOut(build buildSrc) string {
	return filepath.Join(c.Out, build.name)
}

func (c *Config) srcs() ([]buildSrc, error) {
	srcs := []buildSrc{}

	for _, src := range c.Src {
		dirSrcs, err := c.srcFromDir(src)

		if err != nil {
			return srcs, err
		}

		srcs = slices.Concat(srcs, dirSrcs)
	}

	return srcs, nil
}

func (c *Config) srcFromDir(srcDir string) ([]buildSrc, error) {
	srcPaths := []buildSrc{}

	slog.Debug("Inspecting src directory", "src", srcDir)

	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() || path == srcDir {
			return nil
		}

		name := filepath.Base(path)
		abs, err := filepath.Abs(path)

		if err != nil {
			return err
		}

		rel, err := filepath.Rel(c.base, abs)

		if err != nil {
			return err
		}

		slog.Debug("Found module", "src", srcDir, "path", rel, "name", name)

		bs := buildSrc{
			name:     name,
			fullPath: abs,
			relPath:  rel,
			srcDir:   srcDir,
		}

		srcPaths = append(srcPaths, bs)

		return fs.SkipDir
	})

	return srcPaths, err
}
