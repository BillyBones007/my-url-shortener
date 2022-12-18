package db

import (
	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

// Интерфейс для работы с некой базой данных
type DBase interface {
	InsertURL(longURL string, h hasher.URLHasher) error
	SelectLongURL(shortURL string) (model models.Model, err error)
	SelectShortURL(longURL string) (model models.Model, err error)
	URLIsExist(model models.Model) bool
}
