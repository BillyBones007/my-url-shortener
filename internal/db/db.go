package db

import (
	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

// Интерфейс для работы с некой базой данных
type DBase interface {
	InsertURL(m *models.MainModel, h hasher.URLHasher) error
	SelectLongURL(m *models.Model) (*models.Model, error)
	// SelectShortURL(m *models.Model) (*models.Model, error)
	// URLIsExist(m *models.Model) bool // утратила актуальность
	UUIDIsExist(uuid string) bool
	SelectAllForUUID(uuid string) ([]models.Model, error)
}
