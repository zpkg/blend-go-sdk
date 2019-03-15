package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"github.com/blend/go-sdk/db"
)

func main() {

	conn := db.MustNewFromEnv()
	conn.Open()

	_, err := conn.Connection().Query("select * from foo")

	fmt.Printf("error: %#v\n", err)
	fmt.Printf("parsed: %#v\n", db.Error(err))
}
