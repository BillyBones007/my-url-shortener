package main

import (
	"log"

	"github.com/BillyBones007/my-url-shortener/internal/app/server"
	"github.com/BillyBones007/my-url-shortener/internal/db/maps"
	"github.com/BillyBones007/my-url-shortener/internal/hasher/randchars"
)

func main() {
	addr := "localhost:8080"
	db := maps.NewStorage()
	hash := randchars.URLHash{}
	server := server.NewServer(addr, db, hash)
	log.Fatal(server.ListenAndServe())
}
