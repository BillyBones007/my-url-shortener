package handlers

import (
	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

// Структура для работы с хранилищем
type Handler struct {
	Storage db.DBase         // объект хранилища, реализующий интерфейс db.DBase
	Hasher  hasher.URLHasher // объект хэшера, реализующий интерфейс hasher.URLHasher
	BaseURL string           // базовый адрес результирующего сокращенного url
}
