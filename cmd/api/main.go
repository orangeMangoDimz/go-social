package main

import (
	"log"

	"github.com/orangeMangoDimz/go-social/internal/env"
)

func main() {
	app := application{
		config: config{
			addr: env.GetString("ADDR", ":8000"),
		},
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
