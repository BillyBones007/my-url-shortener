package main

import (
	"log"
	"net/http"

	"github.com/BillyBones007/my-url-shortener/internal/handlers"
)

func main() {

	http.HandleFunc("/", handlers.ShortUrlHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
