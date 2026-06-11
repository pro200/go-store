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
	// "~/<name>.store" -> /Users/me/main.store (실행 파일 이름으로 치환)
	db, err := store.New("/tmp/test.store", store.WithTimeout(3*time.Second))
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
