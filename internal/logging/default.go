package logging

import (
	"io"
	"log/slog"
)

func SetDefault(component string, level slog.Level, additionalWriters ...io.Writer) {
	writers := append([]io.Writer{GetConsoleWriter()}, additionalWriters...)
	slog.SetDefault(NewComponentLogger(component, level, writers...))
}

func init() {
	// Set a useful process-wide logger before application initialization starts.
	// This may be overridden by the application.
	SetDefault(defaultComponent, slog.LevelInfo)
}
