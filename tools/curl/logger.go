package curl

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	// "moon-cost/logging"
)

type DiscardLogger struct{}

func (d DiscardLogger) Write(p []byte) (int, error) {
	return 0, nil
}

type JSONLogger struct {
	slog.Handler
	logger *log.Logger
}

func (jl *JSONLogger) Handle(ctx context.Context, record slog.Record) error {
	// level := logging.Red(record.Level.String())

	fields := make(map[string]any, record.NumAttrs()+1)
	if record.Message != "" {
		fields["msg"] = record.Message
	}

	record.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	data, err := json.Marshal(fields)

	if err != nil {
		return err
	}

	jl.logger.Println(string(data))

	// jl.logger.Printf("[%s] %s: %s", level, record.Message, string(data))

	return nil
}

func NewJSONLogger(out io.Writer, level slog.Level) *JSONLogger {
	return &JSONLogger{
		logger: log.New(out, "", 0),
		Handler: slog.NewJSONHandler(out, &slog.HandlerOptions{
			Level: level,
		}),
	}
}
