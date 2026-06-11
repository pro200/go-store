package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppName(t *testing.T) {
	tests := []struct {
		execPath string
		want     string
	}{
		{"", "main"},
		{"/tmp/go-build123/b001/exe/main", "main"},
		{"/tmp/go_build_example/exe", "main"},
		{"/home/user/project/store.test", "main"},
		{"/usr/local/bin/myapp", "myapp"},
	}

	for _, tt := range tests {
		if got := appName(tt.execPath); got != tt.want {
			t.Errorf("appName(%q) = %q, want %q", tt.execPath, got, tt.want)
		}
	}
}

func TestResolvePathHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skipf("no home directory: %v", err)
	}

	got, err := resolvePath("~/x/<name>.store", "/opt/myapp")
	if err != nil {
		t.Fatalf("resolvePath: %v", err)
	}
	want := filepath.Join(home, "x", "myapp.store")
	if got != want {
		t.Fatalf("resolvePath = %q, want %q", got, want)
	}
}

func TestResolvePathRelative(t *testing.T) {
	got, err := resolvePath("data/<name>.store", "")
	if err != nil {
		t.Fatalf("resolvePath: %v", err)
	}
	if !filepath.IsAbs(got) {
		t.Fatalf("resolvePath returned relative path %q", got)
	}
	if filepath.Base(got) != "main.store" {
		t.Fatalf("placeholder not replaced: %q", got)
	}
}
