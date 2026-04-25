// Package logging initialises the global structured logger using Go 1.21+
// stdlib slog. JSON output is chosen so Cloud Logging can parse log entries
// automatically without a custom parser.
package logging

import (
	"log/slog"
	"os"
)

// Init sets the default slog logger to a JSON handler writing to stdout.
// Call once at program start, before any log output.
func Init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))
}
