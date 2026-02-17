package store

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/pro200/go-store/lib"
	"github.com/vmihailenco/msgpack/v5"
	"go.etcd.io/bbolt"
)

type Config struct {
	Bucket string
	Path   string
}

type Store struct {
	bucket string
	db     *bbolt.DB
	mu     sync.RWMutex
}

// New 함수는 주어진 설정(Config) 또는 기본 설정으로 새로운 Store 인스턴스를 생성합니다.
// 버킷이 존재하지 않을 경우 생성하며, 데이터베이스 파일은 지정된 경로에 저장됩니다.
func New(cfg ...Config) (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	c := Config{
		Bucket: "default",
		Path:   home + "/.default.store",
	}

	if len(cfg) > 0 {
		tmp := cfg[0]
		if tmp.Bucket != "" {
			c.Bucket = tmp.Bucket
		}
		if tmp.Path != "" {
			c.Path = tmp.Path
		}
	}

	// DB 오픈
	db, err := bbolt.Open(c.Path, 0600, nil)
	if err != nil {
		return nil, err
	}

	// 버킷 생성
	err = db.Update(func(tx *bbolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte(c.Bucket))
		return e
	})
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return &Store{bucket: c.Bucket, db: db}, nil
}

// Bucket은 현재 사용 중인 버킷 이름을 반환하거나, 지정된 이름으로 버킷을 생성하고 이를 기본 버킷으로 설정합니다.
func (s *Store) Bucket(name ...string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(name) == 0 {
		return s.bucket
	}

	err := s.db.Update(func(tx *bbolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists([]byte(name[0]))
		return e
	})
	if err != nil {
		return ""
	}

	s.bucket = name[0]
	return s.bucket
}

// DeleteBucket은 지정된 이름의 버킷을 데이터베이스에서 삭제합니다.
// 존재하지 않는 버킷 이름을 전달하면 오류를 반환합니다.
func (s *Store) DeleteBucket(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(name))
	})
}

// Set 메서드는 주어진 키-값 페어를 현재 사용 중인 버킷에 저장합니다.
// 키는 문자열로, 값은 직렬화 가능한 어떤 데이터 타입도 가능합니다.
// 내부적으로 데이터는 msgpack으로 직렬화되고 암호화됩니다.
func (s *Store) Set(key string, value any) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := msgpack.Marshal(value)
	if err != nil {
		return fmt.Errorf("msgpack marshal failed: %w", err)
	}

	encryptedData, _ := lib.Encrypt(string(data))

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		return b.Put([]byte(key), []byte(encryptedData))
	})
}

// Get 메서드는 주어진 키에 해당하는 값을 복호화하여 dest에 언마샬링합니다.
// 키가 존재하지 않을 경우 "key not found" 오류를 반환합니다.
func (s *Store) Get(key string, dest any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(s.bucket))
		data := bucket.Get([]byte(key))
		if data == nil {
			return errors.New("key not found")
		}

		decryptedData, _ := lib.Decrypt(string(data))
		return msgpack.Unmarshal([]byte(decryptedData), dest)
	})
}

// Keys 메서드는 현재 설정된 버킷 내의 모든 키 목록을 문자열 슬라이스로 반환합니다.
func (s *Store) Keys() ([]string, error) {
	var keys []string
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			keys = append(keys, string(k))
		}
		return nil
	})
	return keys, err
}

// Delete 메서드는 현재 선택된 버킷에서 주어진 키에 해당하는 데이터를 삭제합니다.
// 키가 존재하지 않을 경우 오류를 반환하지 않고 성공으로 간주됩니다.
// 오류 발생 시 데이터베이스 트랜잭션은 롤백됩니다.
func (s *Store) Delete(key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(s.bucket))
		return b.Delete([]byte(key))
	})
}

// Close 메서드는 데이터베이스 연결을 안전하게 닫습니다.
// 연결이 성공적으로 닫히면 nil을 반환하며, 실패 시 오류를 반환합니다.
func (s *Store) Close() error {
	return s.db.Close()
}
