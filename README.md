# go-store

`go-store`는 **bbolt 기반의 경량 디스크 저장 Key-Value 스토리지**입니다.
값은 `msgpack`으로 직렬화한 뒤 AES-256으로 암호화하여 저장합니다.

- 저장소: bbolt (파일 기반)
- 직렬화: msgpack
- 암호화: AES-256-CBC, 레코드별 랜덤 IV

## Features

- 파일 기반 영속 저장소 (부모 디렉토리 자동 생성)
- 임의 타입 저장 가능 (`any`)
- msgpack 직렬화 + 암호화 저장 (레코드별 랜덤 IV)
- **머신 바인딩 암호화**: 키가 머신 식별자에서 유도되어 다른 컴퓨터에서는 데이터를 읽을 수 없음
- 머신 식별자 폴백 체인 (인터넷 연결 불필요):
  1. OS machine id (macOS: IOPlatformUUID / Linux: /etc/machine-id / Windows: MachineGuid)
  2. 첫 번째 활성 네트워크 인터페이스의 MAC 주소
  3. hostname
  - path 편의 기능: `~/` 홈경로 확장, `/tmp` OS 임시 디렉토리 확장, `<name>` 실행 파일 이름으로 치환
- 파일 락 타임아웃 설정 가능 (`WithTimeout`, 기본 1초)

## Installation

```bash
go get github.com/pro200/go-store
```

## Quick Start

```go
package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/pro200/go-store"
)

type User struct {
  Name string
  Age  int
  At   time.Time
}

func main() {
  // "/tmp/<name>.store" -> /var/folders/7w/2tbd3m4s.../T/tmp/main.store (실행 파일 이름으로 치환)
  db, err := store.New("/tmp/<name>.store", store.WithTimeout(3*time.Second))
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()
  
  // 저장
  err = db.Set("user:1", User{
    Name: "Kim",
    Age:  30,
    At:   time.Now(),
  })
  if err != nil {
    log.Fatal(err)
  }
  
  // 조회 (dest 방식)
  var user User
  if err := db.Get("user:1", &user); err != nil {
    log.Fatal(err)
  }
  fmt.Println("User:", user)
  
  // 키 목록
  keys, err := db.Keys()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Keys:", keys)
  
  // 삭제
  if err := db.Delete("user:1"); err != nil {
    log.Fatal(err)
  }
  
  // 삭제 확인
  err = db.Get("user:1", &user)
  if errors.Is(err, store.ErrKeyNotFound) {
    fmt.Println("After delete: key not found (expected)")
  }
}
```

## API

| 함수 | 설명 |
|------|------|
| `New(path string, opts ...Option) (*Store, error)` | 스토어 열기/생성 |
| `(*Store) Set(key string, value any) error` | 저장 |
| `(*Store) Get(key string, dest any) error` | 조회 (dest는 포인터) |
| `(*Store) Delete(key string) error` | 삭제 (없는 키는 no-op) |
| `(*Store) Keys() ([]string, error)` | 전체 키 목록 (사전순) |
| `(*Store) Close() error` | 닫기 |

### Options

| 옵션 | 설명 |
|------|------|
| `WithTimeout(d time.Duration)` | 파일 락 획득 대기 시간 (기본 1초) |

### Errors

`errors.Is`로 비교할 수 있습니다.

- `store.ErrKeyNotFound` — 키가 존재하지 않음
- `store.ErrEmptyKey` — 빈 키 사용

## 주의 사항

- **머신 바인딩**: 암호화 키가 머신 식별자에서 유도됩니다. `.store` 파일을 다른
  컴퓨터로 복사하면 데이터를 복호화할 수 없습니다. (의도된 동작)
- **하위 호환**: 저장 형식과 키 유도 방식이 변경되어 이전 버전(v1.0.x)으로
  작성된 `.store` 파일은 읽을 수 없습니다. 기존 파일을 삭제 후 사용하세요.
- **무결성**: AES-256-CBC는 변조 감지(인증)를 제공하지 않습니다. 변조된 데이터는
  복호화 에러 또는 디코딩 에러로 나타납니다. 향후 AES-GCM 전환을 검토 중입니다.
