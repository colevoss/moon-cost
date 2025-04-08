package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
)

type BuildEvent struct {
	ImportPath string
	Action     string
	Output     string
}

func (b BuildEvent) String() string {
	return strings.TrimSpace(b.Output)
}

func (b BuildEvent) Error() string {
	return b.String()
}

func ReadEvents(r io.Reader) []BuildEvent {
	dec := json.NewDecoder(r)

	events := []BuildEvent{}
	var errs error

	for dec.More() {
		var event BuildEvent

		if err := dec.Decode(&event); err != nil {
			errs = errors.Join(errs, err)
			continue
		}

		events = append(events, event)
	}

	if errs != nil {
		slog.Error(
			"Error reading build event output",
		)

		fmt.Println(errs)
	}

	return events
}
