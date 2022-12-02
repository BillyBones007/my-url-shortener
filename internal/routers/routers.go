package routers

import (
	"github.com/BillyBones007/my-url-shortener/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{shortUrl}", handlers.GetLongURLHandler)
		r.Post("/", handlers.CreateShortURLHandler)
		r.Put("/", handlers.BadRequestHandler)
		r.Delete("/", handlers.BadRequestHandler)
		r.Head("/", handlers.BadRequestHandler)
	})

	return r
}
