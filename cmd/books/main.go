package main

import (
	"books/internal/config"
	"books/internal/server"
	"books/internal/storage"
	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"
	"log"
)

func main() {
	cfg := new(config.Config)
	if _, err := toml.DecodeFile("configs/books.toml", cfg); err != nil {
		panic(err)
	}

	store := storage.New(cfg)
	if err := store.Open(); err != nil {
		panic(err)
	}

	sessionStore := sessions.NewCookieStore([]byte(cfg.SecretKey))

	s := server.New(cfg, store, sessionStore)

	log.Fatal(s.Run())
}
