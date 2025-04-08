package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"path/filepath"
	"slices"
)

type BuildEvent struct {
	ImportPath string
	Action     string
	Output     string
}

type Config struct {
	base   string
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

func ConfigFrom(base string, r io.Reader) (Config, error) {
	var config Config

	if err := json.NewDecoder(r).Decode(&config); err != nil {
		return config, err
	}

	config.base = base

	return config, nil
}

func (c *Config) Run(ctx context.Context) error {
	srcs, err := c.srcs()

	if err != nil {
		return err
	}

	for _, src := range srcs {
		if slices.Contains(c.Ignore, src.name) {
			slog.Debug("Skipping build target", "module", src.name, "path", src.relPath, "srcDir", src.srcDir)
			continue
		}

		if err := c.build(ctx, src); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) build(ctx context.Context, build buildSrc) error {
	output := c.buildOut(build)
	buildCommand := goBuild(output, build.fullPath)

	slog.Info("Building module", "src", build.relPath, "out", output)

	cmd := buildCommand.exec(ctx)
	out, err := cmd.StdoutPipe()

	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func(r io.Reader) {
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			text := scanner.Text()

			var event BuildEvent

			if err := json.Unmarshal([]byte(text), &event); err != nil {
				fmt.Printf("unmarshall err: %v\n", err)
				return
			}

			fmt.Printf("%+v\n", event)
		}
	}(out)

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
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
