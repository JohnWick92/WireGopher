package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeTemp writes a temp file and returns its path.
func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoad_NotFound(t *testing.T) {
	_, err := Load("/path/that/does/not/exist.yaml")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
	if !strings.Contains(err.Error(), "hint:") {
		t.Fatalf("error should contain actionable hint, got %q", err.Error())
	}
}

func TestLoad_EmptyPath(t *testing.T) {
	_, err := Load("")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestLoad_Unreadable(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("running as root: permission bits are bypassed")
	}
	p := writeTemp(t, "locked.yaml", "server: {}")
	if err := os.Chmod(p, 0o000); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(p, 0o644) })

	_, err := Load(p)
	if !errors.Is(err, ErrRead) {
		t.Fatalf("want ErrRead, got %v", err)
	}
}

func TestLoad_ParseError(t *testing.T) {
	p := writeTemp(t, "bad.yaml", "server: [oops")
	_, err := Load(p)
	if !errors.Is(err, ErrParse) {
		t.Fatalf("want ErrParse, got %v", err)
	}
}

func TestLoad_Valid(t *testing.T) {
	p := writeTemp(t, "ok.yaml", "server: {}\n")
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load returned (nil, nil) — must never happen")
	}
}

func TestLoad_NeverReturnsNilWithoutError(t *testing.T) {
	p := writeTemp(t, "ok.yaml", "server: {}\n")
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load returned (nil, nil) — must never happen")
	}
}
