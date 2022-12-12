package server

import (
	"net/http"

	"github.com/BillyBones007/my-url-shortener/internal/routers"
)

// Возвращает ссылку на готовый экземпляр сервера.
// Получает в качестве параметра указатель на структуру конфига.
func NewServer(cfg *Config) *http.Server {
	r := routers.NewRouter(cfg.Storage, cfg.Hash, cfg.BaseURL)
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}
	return server
}
