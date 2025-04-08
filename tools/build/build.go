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
	Name     string
	FullPath string
	RelPath  string
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
	targets, err := c.srcs()

	if err != nil {
		return err
	}

	for _, target := range targets {
		if slices.Contains(c.Ignore, target.Name) {
			slog.Debug("Skipping build target", "target", target.Name, "path", target.FullPath)
			continue
		}

		if err := c.build(ctx, target); err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) build(ctx context.Context, build buildSrc) error {
	output := c.buildOut(build)
	buildCommand := goBuild(output, build.FullPath)

	fmt.Printf("Building %s -> %s\n", build.RelPath, output)

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
	return filepath.Join(c.Out, build.Name)
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

	slog.Debug("Walking src directory", "src", srcDir)

	err := filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() || path == srcDir {
			return nil
		}

		name := filepath.Base(path)
		abs, err := filepath.Abs(path)

		if err != nil {
			return err
		}

		fullPath, err := filepath.Rel(c.base, abs)

		if err != nil {
			return err
		}

		bs := buildSrc{
			Name:     name,
			FullPath: abs,
			RelPath:  fullPath,
		}

		srcPaths = append(srcPaths, bs)

		return fs.SkipDir
	})

	return srcPaths, err
}
