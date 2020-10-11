package main

import (
	"books/internal/config"
	// "books/internal/seed"
	"books/internal/storage"
	"github.com/BurntSushi/toml"
)

func main() {
	config := new(config.Config)
	if _, err := toml.DecodeFile("configs/books.toml", config); err != nil {
		panic(err)
	}

	storage := storage.New(config)
	storage.MustOpen()

}
