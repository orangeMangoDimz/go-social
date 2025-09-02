package main

import (
	"log"

	"github.com/orangeMangoDimz/go-social/internal/env"
	"github.com/orangeMangoDimz/go-social/store"
)

func main() {
	app := application{
		config: config{
			addr: env.GetString("ADDR", ":8000"),
		},
		store: store.NewStore(nil),
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
