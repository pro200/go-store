package store

import (
	"errors"
	"path/filepath"
	"slices"
	"testing"
	"time"
)

type testUser struct {
	Name string
	Age  int
	At   time.Time
}

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(filepath.Join(t.TempDir(), "test.store"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestSetGetRoundtrip(t *testing.T) {
	s := newTestStore(t)

	want := testUser{Name: "Kim", Age: 30, At: time.Now()}
	if err := s.Set("user:1", want); err != nil {
		t.Fatalf("Set: %v", err)
	}

	var got testUser
	if err := s.Get("user:1", &got); err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != want.Name || got.Age != want.Age || !got.At.Equal(want.At) {
		t.Fatalf("roundtrip mismatch: got %+v, want %+v", got, want)
	}
}

func TestGetMissingKey(t *testing.T) {
	s := newTestStore(t)

	var dest string
	err := s.Get("missing", &dest)
	if !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("Get missing key: got %v, want ErrKeyNotFound", err)
	}
}

func TestEmptyKey(t *testing.T) {
	s := newTestStore(t)

	if err := s.Set("", "v"); !errors.Is(err, ErrEmptyKey) {
		t.Errorf("Set empty key: got %v, want ErrEmptyKey", err)
	}
	var dest string
	if err := s.Get("", &dest); !errors.Is(err, ErrEmptyKey) {
		t.Errorf("Get empty key: got %v, want ErrEmptyKey", err)
	}
	if err := s.Delete(""); !errors.Is(err, ErrEmptyKey) {
		t.Errorf("Delete empty key: got %v, want ErrEmptyKey", err)
	}
}

func TestDelete(t *testing.T) {
	s := newTestStore(t)

	if err := s.Set("k", "v"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.Delete("k"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	var dest string
	if err := s.Get("k", &dest); !errors.Is(err, ErrKeyNotFound) {
		t.Fatalf("Get after delete: got %v, want ErrKeyNotFound", err)
	}

	// Deleting a missing key is a no-op.
	if err := s.Delete("missing"); err != nil {
		t.Fatalf("Delete missing key: %v", err)
	}
}

func TestKeys(t *testing.T) {
	s := newTestStore(t)

	keys, err := s.Keys()
	if err != nil {
		t.Fatalf("Keys on empty store: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("Keys on empty store = %v, want empty", keys)
	}

	for _, k := range []string{"b", "a", "c"} {
		if err := s.Set(k, k); err != nil {
			t.Fatalf("Set %q: %v", k, err)
		}
	}

	keys, err = s.Keys()
	if err != nil {
		t.Fatalf("Keys: %v", err)
	}
	if !slices.Equal(keys, []string{"a", "b", "c"}) {
		t.Fatalf("Keys = %v, want [a b c]", keys)
	}
}

func TestPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "persist.store")

	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := s.Set("k", "persisted"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	s, err = New(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	defer s.Close()

	var got string
	if err := s.Get("k", &got); err != nil {
		t.Fatalf("Get after reopen: %v", err)
	}
	if got != "persisted" {
		t.Fatalf("Get after reopen = %q, want %q", got, "persisted")
	}
}

func TestCreatesParentDirs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "c.store")

	s, err := New(path)
	if err != nil {
		t.Fatalf("New with missing parent dirs: %v", err)
	}
	defer s.Close()

	if err := s.Set("k", "v"); err != nil {
		t.Fatalf("Set: %v", err)
	}
}

func TestLockTimeout(t *testing.T) {
	path := filepath.Join(t.TempDir(), "locked.store")

	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer s.Close()

	timeout := 100 * time.Millisecond
	start := time.Now()
	_, err = New(path, WithTimeout(timeout))
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("second New on locked file should fail")
	}
	if elapsed > 5*timeout {
		t.Fatalf("lock failure took %v, want ~%v", elapsed, timeout)
	}
}
