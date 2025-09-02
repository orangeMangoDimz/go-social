package main

import (
	"log"

	"github.com/orangeMangoDimz/go-social/internal/db"
	"github.com/orangeMangoDimz/go-social/internal/env"
	"github.com/orangeMangoDimz/go-social/store"
)

func main() {
	cfg := dbConfig{
		addr:         env.GetString("DB_ADDR", "postgres://postgres:root@localhost/social?sslmode=disable"),
		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	db, err := db.New(cfg.addr, cfg.maxOpenConns, cfg.maxIdleConns, cfg.maxIdleTime)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()
	log.Printf("DB Database connection pool established")

	app := application{
		config: config{
			addr: env.GetString("ADDR", ":8000"),
			db:   cfg,
		},
		store: store.NewStore(db),
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
