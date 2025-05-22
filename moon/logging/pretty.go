package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"
)

type PrettyHandler struct {
	slog.Handler
	logger *log.Logger
}

func NewPrettyHandler(out io.Writer, level slog.Level) *PrettyHandler {
	return &PrettyHandler{
		logger: log.New(out, "", 0),
		Handler: slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: level,
		}),
	}
}

func (sl *PrettyHandler) Handle(ctx context.Context, record slog.Record) error {
	fields := make(map[string]any, record.NumAttrs())

	record.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	color := ColorWhite

	switch record.Level {
	case slog.LevelDebug:
		color = ColorCyan
	case LevelVerbose:
		color = ColorLightMagenta
	case slog.LevelInfo:
		color = ColorGreen
	case slog.LevelWarn:
		color = ColorYellow
	case slog.LevelError:
		color = ColorRed
	}

	var levelBase string

	if record.Level == LevelVerbose {
		levelBase = "VERBOSE"
	} else {
		levelBase = record.Level.String()
	}

	level := fmt.Sprintf("[%s]", levelBase)

	var b strings.Builder

	b.WriteString(fmt.Sprintf(
		"%s %s",
		Color(color, level),
		record.Message,
	))

	if len(fields) > 0 {
		data, err := json.MarshalIndent(fields, "", "  ")

		if err != nil {
			return err
		}

		b.WriteRune(' ')
		b.WriteString(string(data))
	}

	sl.logger.Println(b.String())

	return nil
}
