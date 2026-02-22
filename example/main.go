package main

import (
	"fmt"
	"log"
	"time"

	"github.com/pro200/go-store"
)

// =========================
// main
// =========================

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
