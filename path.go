package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// resolvePath expands a leading ~ to the user's home directory, makes
// the path absolute, and replaces the <name> placeholder with the
// application name derived from execPath.
func resolvePath(path, execPath string) (string, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("store: expand home directory: %w", err)
		}
		path = filepath.Join(home, path[1:])
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("store: resolve path %q: %w", path, err)
	}

	return strings.ReplaceAll(path, "<name>", appName(execPath)), nil
}

// appName derives the application name from the executable path.
// It returns "main" for go run / go test binaries (built into the
// go-build cache or named *.test) and when execPath is empty.
func appName(execPath string) string {
	if execPath == "" ||
		strings.Contains(execPath, "go-build") ||
		strings.Contains(execPath, "go_build") ||
		strings.HasSuffix(filepath.Base(execPath), ".test") {
		return "main"
	}
	return filepath.Base(execPath)
}
