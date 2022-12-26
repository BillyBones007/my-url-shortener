package main

import (
	"log"

	"github.com/BillyBones007/my-url-shortener/internal/app/server"
)

var cfg server.Config
var fV server.FlagVars

func main() {
	server.ServerConfigurator(&cfg, &fV)
	server := server.NewServer(&cfg)
	log.Fatal(server.ListenAndServe())
}
