package server

import (
	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

// Конфигурация сервера
type Config struct {
	// Адрес запуска HTTP сервера (env $SERVER_ADDRESS)
	ServerAddress string `env:"SERVER_ADDRESS"`
	// Базовый адресс результирующего сокращенного url (env $BASE_URL)
	BaseURL string `env:"BASE_URL"`
	// Путь к файлу хранилища (env $FILE_STORAGE_PATH)
	StoragePATH string `env:"FILE_STORAGE_PATH"`
	// Объект хранилища, реализующий интерфейс db.DBase (env $FILE_STORAGE_PATH)
	Storage db.DBase
	// Объект хэшера, реализующий интерфейс hasher.URLHasher
	Hash hasher.URLHasher
}
