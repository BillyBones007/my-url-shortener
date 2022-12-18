package main

import (
	"log"
	"net/http"

	"github.com/BillyBones007/my-url-shortener/internal/routers"
)

func main() {
	r := routers.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", r))
}
