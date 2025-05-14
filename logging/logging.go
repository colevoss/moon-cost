package logging

import "log/slog"

func Logger(logger *slog.Logger, args ...any) *slog.Logger {
	loggerInstance := logger

	if loggerInstance == nil {
		loggerInstance = slog.Default()
	}

	return loggerInstance.With(args...)
}
