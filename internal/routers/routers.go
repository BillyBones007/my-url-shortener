package routers

import (
	"github.com/BillyBones007/my-url-shortener/internal/auth/randbytes"
	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/handlers"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
	"github.com/BillyBones007/my-url-shortener/internal/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(storage db.DBase, hasher hasher.URLHasher, baseURL string) chi.Router {
	handler := handlers.Handler{Storage: storage, Hasher: hasher, BaseURL: baseURL}
	secretKey, _ := randbytes.RandomBytes(16)
	auth := middlewares.Auth{SecretKey: secretKey}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/", func(r chi.Router) {
		r.Use(auth.CheckCookies)
		r.Use(middlewares.GzipDecompress)
		r.Use(middlewares.GzipCompress)
		r.Get("/{shortUrl}", handler.GetLongURLHandler)
		r.Get("/api/user/urls", handler.GetAllURLsHandler)
		r.Get("/ping", handler.GetPing)
		r.Post("/", handler.CreateShortURLHandler)
		r.Post("/api/shorten", handler.CreateInJSONShortURLHandler)
		r.Put("/", handler.BadRequestHandler)
		r.Delete("/", handler.BadRequestHandler)
		r.Head("/", handler.BadRequestHandler)
	})

	return r
}
