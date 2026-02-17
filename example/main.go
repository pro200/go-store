package main

import (
	"fmt"

	"github.com/pro200/go-store"
)

func main() {
	db, err := store.New(store.Config{
		Path: "/tmp/test.store",
	})
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Set("name", "james")

	var data string
	db.Get("name", &data)
	fmt.Println("->", data)
}
