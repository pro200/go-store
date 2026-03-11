package store

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/pro200/go-store/lib"
	"github.com/vmihailenco/msgpack/v5"
	"go.etcd.io/bbolt"
)

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrEmptyKey    = errors.New("empty key")
)

var rootBucket = []byte("__root__")

type Store struct {
	db *bbolt.DB
}

func New(path string) (*Store, error) {
	defaultName := "main"
	fullpath, _ := os.Executable()
	if !strings.Contains(fullpath, "go-build") && !strings.Contains(fullpath, "go_build") {
		defaultName = filepath.Base(fullpath)
	}

	// ~/ 홈경로 반영 후 절대경로
	path, err := homedir.Expand(path)
	if err != nil {
		return nil, err
	}
	path, _ = filepath.Abs(path)

	// <name> 을 앱이름으로 치환
	path = strings.ReplaceAll(path, "<name>", defaultName)

	fmt.Println("path->", path)
	// Timeout을 설정하면 해당 시간 내에 lock을 획득하지 못할 경우 "timeout" 에러를 반환하고 즉시 종료됩니다.
	opts := &bbolt.Options{
		Timeout: 1 * time.Nanosecond,
	}

	db, err := bbolt.Open(path, 0600, opts)
	if err != nil {
		return nil, err
	}

	// 내부 고정 루트 버킷 생성
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err2 := tx.CreateBucketIfNotExists(rootBucket)
		return err2
	})
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

/*
 * KV 기능 (dest any 방식)
 */
func (s *Store) Set(key string, value any) error {
	if key == "" {
		return ErrEmptyKey
	}

	data, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}

	data, err = lib.Encrypt(data)
	if err != nil {
		return err
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(rootBucket)
		return b.Put([]byte(key), data)
	})
}

func (s *Store) Get(key string, dest any) error {
	if key == "" {
		return ErrEmptyKey
	}

	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(rootBucket)
		data := b.Get([]byte(key))
		if data == nil {
			return ErrKeyNotFound
		}

		var err error
		data, err = lib.Decrypt(data)
		if err != nil {
			return err
		}

		return msgpack.Unmarshal(data, dest)
	})
}

func (s *Store) Delete(key string) error {
	if key == "" {
		return ErrEmptyKey
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(rootBucket)
		return b.Delete([]byte(key))
	})
}

func (s *Store) Keys() ([]string, error) {
	var keys []string

	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(rootBucket)
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, string(k))
		}
		return nil
	})

	return keys, err
}
