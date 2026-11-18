// Package config loads the WireGopher configuration from YAML.
package config

import (
	"errors"
	"fmt"
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
	expandedData := os.ExpandEnv(string(raw))
	if err := yaml.Unmarshal([]byte(expandedData), &cfg); err != nil {
		return nil, fmt.Errorf("%w: %q: %v", ErrParse, path, err)
	}
	return &cfg, nil
}
