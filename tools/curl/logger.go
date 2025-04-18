package curl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"moon-cost/logging"
	"strings"
	// "moon-cost/logging"
)

type DiscardLogger struct{}

func (d DiscardLogger) Write(p []byte) (int, error) {
	return 0, nil
}

type JSONHandler struct {
	slog.Handler
	enabled bool
	logger  *log.Logger
}

func NewJSONHandler(out io.Writer, level slog.Level, enabled bool) *JSONHandler {
	return &JSONHandler{
		logger:  log.New(out, "", 0),
		enabled: enabled,
		Handler: slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: level,
		}),
	}
}

// func (jl *JSONHandler) Handle(ctx context.Context, record slog.Record) error {
// 	// level := logging.Red(record.Level.String())
//
// 	fields := make(map[string]any, record.NumAttrs()+1)
// 	if record.Message != "" {
// 		fields["msg"] = record.Message
// 	}
//
// 	record.Attrs(func(a slog.Attr) bool {
// 		fields[a.Key] = a.Value.Any()
//
// 		return true
// 	})
//
// 	data, err := json.Marshal(fields)
//
// 	if err != nil {
// 		return err
// 	}
//
// 	jl.logger.Println(string(data))
//
// 	return nil
// }

func (jl *JSONHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if !jl.enabled {
		return false
	}

	return jl.Handler.Enabled(ctx, level)
}

type StandardHandler struct {
	slog.Handler
	logger  *log.Logger
	enabled bool
}

func NewStandardHandler(out io.Writer, level slog.Level, enabled bool) *StandardHandler {
	return &StandardHandler{
		logger:  log.New(out, "", 0),
		enabled: enabled,
		Handler: slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: level,
		}),
	}
}

func (sl *StandardHandler) Enabled(ctx context.Context, level slog.Level) bool {
	if !sl.enabled {
		return false
	}

	return sl.Handler.Enabled(ctx, level)
}

var maxLevelLength = 5 + 2 // 5 for ERROR 2 for  []

func (sl *StandardHandler) Handle(ctx context.Context, record slog.Record) error {
	fields := make(map[string]any, record.NumAttrs())

	record.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	color := logging.ColorWhite

	switch record.Level {
	case slog.LevelDebug:
		color = logging.ColorCyan
	case slog.LevelInfo:
		color = logging.ColorGreen
	case slog.LevelWarn:
		color = logging.ColorYellow
	case slog.LevelError:
		color = logging.ColorRed
	}

	level := fmt.Sprintf("[%s]", record.Level.String())

	var b strings.Builder

	b.WriteString(fmt.Sprintf(
		"%s %s",
		logging.Color(color, level),
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
