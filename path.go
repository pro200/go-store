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

	if path == "/tmp" || strings.HasPrefix(path, "/tmp/") {
		// OS 마다 보안(샌드박스) 체계와 다중 사용자 보호 정책이 다르기 때문에, /tmp 경로는 OS가 제공하는 임시 디렉토리로 대체
		// macOS 예: /var/folders/7w/2tbd3m4s0hx0qncd_m78p2580000gn/T
		tmp := os.TempDir()
		path = filepath.Join(tmp, path[1:])
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
