// Package config loads the WireGopher configuration from YAML.
package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrNotFound = errors.New("config file not found")
	ErrRead     = errors.New("config file unreadable")
	ErrParse    = errors.New("config parse error")
)

// Config mirrors the YAML. Server type lives in internal/http.
type Config struct {
	Server   Server   `yaml:"server"`
	Firewall Firewall `yaml:"firewall"`
}

// Load reads and parses the YAML at path.
// It only checks: file exists, is readable, and is valid YAML.
// Never returns (nil, nil).
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("%w: empty path (use --config or WIREGOPHER_CONFIG)", ErrNotFound)
	}

	raw, err := os.ReadFile(path)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			return nil, fmt.Errorf(
				"%w: %q\n"+
					"  hint: mount it, e.g. -v ./configs/wiregopher.yaml:/etc/wiregopher/wiregopher.yaml:ro\n"+
					"  hint: or set --config / WIREGOPHER_CONFIG",
				ErrNotFound, path,
			)
		case errors.Is(err, os.ErrPermission):
			return nil, fmt.Errorf("%w: %q (check container user / file mode)", ErrRead, path)
		default:
			return nil, fmt.Errorf("%w: %q: %v", ErrRead, path, err)
		}
	}

	var cfg Config
	if err := yaml.Unmarshal(raw, &cfg); err != nil {
		return nil, fmt.Errorf("%w: %q: %v", ErrParse, path, err)
	}
	return &cfg, nil
}

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
