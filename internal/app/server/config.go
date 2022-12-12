package server

import (
	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

type Config struct {
	ServerAddress string           `env:"SERVER_ADDRESS"` // фдрес запуска HTTP сервера (env $SERVER_ADDRESS)
	BaseURL       string           `env:"BASE_URL"`       // базовый адресс результирующего сокращенного url (env $BASE_URL)
	Storage       db.DBase         // объект хранилища, реализующий интерфейс db.DBase
	Hash          hasher.URLHasher // объект хэшера, реализующий интерфейс hasher.URLHasher
}
