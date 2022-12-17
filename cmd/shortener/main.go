package main

import (
	"log"

	"github.com/BillyBones007/my-url-shortener/internal/app/server"
	"github.com/BillyBones007/my-url-shortener/internal/db/files"
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
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ":8080"
	}
	if cfg.StoragePATH == "" {
		cfg.Storage = maps.NewStorage()
		cfg.Hash = randchars.URLHash{}
		log.Println("INFO: $FILE_STORAGE_PATH is empty, map storage is used")
		return
	}
	cfg.Storage, err = files.NewStorage(cfg.StoragePATH)
	if err != nil {
		log.Printf("ERROR: file storage error %s\n", err)
		log.Println("INFO: map storage is used")
		cfg.Storage = maps.NewStorage()
		cfg.Hash = randchars.URLHash{}
		return
	}
	cfg.Hash = randchars.URLHash{}
}

func main() {
	server := server.NewServer(&cfg)
	log.Fatal(server.ListenAndServe())
}
