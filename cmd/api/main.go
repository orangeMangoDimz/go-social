package main

import (
	"log"

	"github.com/orangeMangoDimz/go-social/internal/db"
	"github.com/orangeMangoDimz/go-social/internal/env"
	"github.com/orangeMangoDimz/go-social/store"
)

const VERSION = "0.0.1"

//	@title			Gopher Social API
//	@description	API for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description

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
			addr:   env.GetString("ADDR", ":8000"),
			db:     cfg,
			env:    env.GetString("ENV", "development"),
			apiURL: env.GetString("EXTERNAL_URL", "localhost:8000"),
		},
		store: store.NewStore(db),
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
