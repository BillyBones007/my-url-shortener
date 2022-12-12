package main

import (
	"log"

	"github.com/BillyBones007/my-url-shortener/internal/app/server"
	"github.com/BillyBones007/my-url-shortener/internal/db/maps"
	"github.com/BillyBones007/my-url-shortener/internal/hasher/randchars"
	"github.com/caarlos0/env/v6"
)

var cfg server.Config

func init() {
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.Storage = maps.NewStorage()
	cfg.Hash = randchars.URLHash{}
}

func main() {
	server := server.NewServer(&cfg)
	log.Fatal(server.ListenAndServe())
}
