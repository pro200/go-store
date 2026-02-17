# go-store

`go-store`는 **bbolt 기반의 경량 디스크 저장 Key-Value 스토리지**입니다.  
값은 `msgpack`으로 직렬화한 뒤 암호화하여 저장합니다.

- 저장소: bbolt (파일 기반)
- 직렬화: msgpack
- 암호화: 내부 `lib.Encrypt / lib.Decrypt`
- 멀티 고루틴 안전 (RWMutex 사용)

## Features

- 파일 기반 영속 저장소
- Bucket 지원 (네임스페이스 개념)
- 임의 타입 저장 가능 (`any`)
- msgpack 직렬화
- 암호화 저장
- Thread-safe

## Installation

```bash
go get github.com/pro200/go-store
```

## Quick Start
```go
package main

import (
	"fmt"
	"github.com/pro200/go-store"
)

func main() {
	db, err := store.New()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// 값 저장
	err = db.Set("name", "pro200")
	if err != nil {
		panic(err)
	}

	// 값 조회
	var name string
	err = db.Get("name", &name)
	if err != nil {
		panic(err)
	}
	
	// Bucket 변경
	db.Bucket("user")

	fmt.Println(name)
}
```

## Configuration

New() 함수는 기본 설정 또는 사용자 정의 설정을 사용할 수 있습니다.

기본 설정
- Bucket: default
- Path: $HOME/.default.store
```go
store, err := store.New()
```
사용자 설정
```go
db, err := store.New(store.Config{
	Bucket: "config",
	Path:   "/tmp/mydb.store",
})
```

## API

### New
```go
func New(cfg ...Config) (*Store, error)
```
새로운 Store 인스턴스를 생성합니다.
DB 파일이 없으면 자동 생성되며, Bucket도 자동 생성됩니다.


### Bucket
```go
func (s *Store) Bucket(name ...string) string
```
- 인자 없음 → 현재 Bucket 이름 반환
- 인자 있음 → Bucket이 없으면 생성 후 기본 Bucket으로 설정


### DeleteBucket
```go
func (s *Store) DeleteBucket(name string) error
```
지정된 Bucket을 삭제합니다.


### Set
```go
func (s *Store) Set(key string, value any) error
```
- msgpack 직렬화
- Encrypt
- bbolt에 저장
```go
db.Set("age", 30)
```


### Get
```go
func (s *Store) Get(key string, dest any) error
```
- 복호화
- msgpack 역직렬화
- 존재하지 않으면 "key not found" 오류 반환
```go
var age int
err := db.Get("age", &age)
```


### Keys
```go
func (s *Store) Keys() ([]string, error)
```
현재 Bucket의 모든 키 목록 반환
```go
keys, _ := db.Keys()
```


### Delete
```go
func (s *Store) Delete(key string) error
```
지정된 키 삭제
존재하지 않아도 에러 없음


### Close
```go
func (s *Store) Close() error
```
DB 연결 종료
