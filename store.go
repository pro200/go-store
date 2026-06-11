// Package store is a lightweight, file-based key-value store built on
// bbolt. Values are serialized with msgpack and encrypted with
// AES-256-CBC using a key derived from a machine identifier, so data
// written on one machine cannot be read on another.
//
// The path passed to New supports two conveniences: a leading ~ is
// expanded to the user's home directory, and the <name> placeholder is
// replaced with the executable name (or "main" under go run / go test).
package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pro200/go-store/internal/machineid"
	"github.com/vmihailenco/msgpack/v5"
	"go.etcd.io/bbolt"
)

var (
	// ErrKeyNotFound is returned by Get when the key does not exist.
	ErrKeyNotFound = errors.New("key not found")
	// ErrEmptyKey is returned when an empty key is passed to Set, Get or Delete.
	ErrEmptyKey = errors.New("empty key")
)

var rootBucket = []byte("__root__")

// Store is an encrypted key-value store backed by a single bbolt file.
// It is safe for concurrent use by multiple goroutines.
type Store struct {
	db  *bbolt.DB
	key []byte
}

// New opens (or creates) the store file at path. Parent directories
// are created as needed. It fails if the file lock cannot be acquired
// within the configured timeout (default 1s, see WithTimeout) or if no
// machine identifier is available for key derivation.
func New(path string, opts ...Option) (*Store, error) {
	cfg := config{timeout: defaultTimeout}
	for _, opt := range opts {
		opt(&cfg)
	}

	execPath, _ := os.Executable()
	path, err := resolvePath(path, execPath)
	if err != nil {
		return nil, err
	}

	fmt.Println("path: ", path)

	machineID, err := machineid.ID()
	if err != nil {
		return nil, fmt.Errorf("store: derive encryption key: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("store: create directory: %w", err)
	}

	db, err := bbolt.Open(path, 0o600, &bbolt.Options{Timeout: cfg.timeout})
	if err != nil {
		return nil, fmt.Errorf("store: open %s: %w", path, err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(rootBucket)
		return err
	})
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("store: create root bucket: %w", err)
	}

	return &Store{db: db, key: deriveKey(machineID)}, nil
}

// Close releases the database file and its lock.
func (s *Store) Close() error {
	return s.db.Close()
}

// Set serializes value with msgpack, encrypts it and stores it under key.
func (s *Store) Set(key string, value any) error {
	if key == "" {
		return ErrEmptyKey
	}

	data, err := msgpack.Marshal(value)
	if err != nil {
		return fmt.Errorf("store: encode value for key %q: %w", key, err)
	}

	data, err = encrypt(s.key, data)
	if err != nil {
		return fmt.Errorf("store: encrypt value for key %q: %w", key, err)
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(rootBucket).Put([]byte(key), data)
	})
}

// Get decrypts the value stored under key and unmarshals it into dest,
// which must be a pointer. It returns ErrKeyNotFound if the key does
// not exist.
func (s *Store) Get(key string, dest any) error {
	if key == "" {
		return ErrEmptyKey
	}

	return s.db.View(func(tx *bbolt.Tx) error {
		data := tx.Bucket(rootBucket).Get([]byte(key))
		if data == nil {
			return ErrKeyNotFound
		}

		data, err := decrypt(s.key, data)
		if err != nil {
			return fmt.Errorf("store: decrypt value for key %q: %w", key, err)
		}

		if err := msgpack.Unmarshal(data, dest); err != nil {
			return fmt.Errorf("store: decode value for key %q: %w", key, err)
		}
		return nil
	})
}

// Delete removes key from the store. Deleting a missing key is a no-op.
func (s *Store) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(rootBucket).Delete([]byte(key))
	})
}

// Keys returns all keys in the store in lexicographic order.
func (s *Store) Keys() ([]string, error) {
	var keys []string

	err := s.db.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket(rootBucket).Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, string(k))
		}
		return nil
	})

	return keys, err
}
