# go-store

`go-store`는 **bbolt 기반의 경량 디스크 저장 Key-Value 스토리지**입니다.  
값은 `msgpack`으로 직렬화한 뒤 암호화하여 저장합니다.

- 저장소: bbolt (파일 기반)
- 직렬화: msgpack
- 암호화: 내부 `lib.Encrypt / lib.Decrypt`

## Features

- 파일 기반 영속 저장소
- 임의 타입 저장 가능 (`any`)
- msgpack 직렬화
- 암호화 저장

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

type User struct {
	Name string
	Age  int
	At   time.Time
}

func main() {
	db, err := store.New("/tmp/test.store")
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
	err = db.Get("user:1", &user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User:", user)

	// 키 목록
	keys, _ := db.Keys()
	fmt.Println("Keys:", keys)

	// 삭제
	_ = db.Delete("user:1")

	err = db.Get("user:1", &user)
	fmt.Println("After delete:", err)
}
```

