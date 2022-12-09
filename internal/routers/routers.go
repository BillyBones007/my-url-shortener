package routers

import (
	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/handlers"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(storage db.DBase, hasher hasher.URLHasher) chi.Router {
	handler := handlers.Handler{Storage: storage, Hasher: hasher}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Get("/{shortUrl}", handler.GetLongURLHandler)
		r.Post("/", handler.CreateShortURLHandler)
		r.Put("/", handler.BadRequestHandler)
		r.Delete("/", handler.BadRequestHandler)
		r.Head("/", handler.BadRequestHandler)
	})

	return r
}
