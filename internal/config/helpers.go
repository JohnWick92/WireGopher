package config

import (
	"errors"
	"log/slog"
	"os"
)

// LogAndExit logs the load error with the right message and exits non-zero.
// Call it only when err != nil.
func LogAndExit(err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		slog.Error("config not found", "error", err)
	case errors.Is(err, ErrRead):
		slog.Error("config unreadable", "error", err)
	case errors.Is(err, ErrParse):
		slog.Error("config syntax error", "error", err)
	default:
		slog.Error("config load failed", "error", err)
	}
	os.Exit(1)
}
