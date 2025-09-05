package main

import (
	"log"

	"github.com/orangeMangoDimz/go-social/internal/db"
	"github.com/orangeMangoDimz/go-social/internal/env"
	"github.com/orangeMangoDimz/go-social/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://postgres:root@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 30, 30, "15m")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStore(conn)
	db.Seed(store, conn)
}
