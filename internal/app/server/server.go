package server

import (
	"net/http"

	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
	"github.com/BillyBones007/my-url-shortener/internal/routers"
)

// Возвращает ссылку на готовый экземпляр сервера. Получает в качестве параметров
// адрес, на котором будет доступен сервер, инициализированное хранилище, реализующее
// интерфейс db.DBase, и инициализированный шифровщик, реализующий интерфейс hasher.URLHasher
func NewServer(address string, storage db.DBase, hasher hasher.URLHasher) *http.Server {
	r := routers.NewRouter(storage, hasher)
	server := &http.Server{
		Addr:    address,
		Handler: r,
	}
	return server
}
