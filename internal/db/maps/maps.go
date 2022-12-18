package maps

import (
	"errors"

	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

// Тип для работы с мапой в роли основного хранилища
type MapStorage struct {
	DataBase map[string]string
}

// Конструктор хранилища. Возвращает указатель на MapStorage
func NewStorage() *MapStorage {
	return &MapStorage{DataBase: make(map[string]string)}
}

// Проверяет, существует ли длинный url в базе
func (m *MapStorage) URLIsExist(model *models.Model) bool {
	flag := false
	for _, long := range m.DataBase {
		if long == model.LongURL {
			flag = true
			break
		}
	}
	return flag
}

// Заполняет мапу. Получает models.Model и хэшер
func (m *MapStorage) InsertURL(model *models.Model, h hasher.URLHasher) error {
	// На случай если GetHash выдаст shortUrl, который уже есть
	// в мапе, используем цикл до тех пор, пока не получим уникальное значение
	for {
		model.ShortURL = h.GetHash(model.LongURL)
		if m.DataBase[model.ShortURL] == "" {
			m.DataBase[model.ShortURL] = model.LongURL
			break
		}
	}
	return nil
}

// Возвращает длинный url из мапы на основе короткого url
func (m *MapStorage) SelectLongURL(model *models.Model) (*models.Model, error) {
	if m.DataBase[model.ShortURL] == "" {
		err := errors.New("URL not found")
		return model, err
	}
	model.LongURL = m.DataBase[model.ShortURL]
	return model, nil
}

// Возвращает короткий url из мапы на основе длинного url
func (m *MapStorage) SelectShortURL(model *models.Model) (*models.Model, error) {
	flag := false
	for k, v := range m.DataBase {
		if v == model.LongURL {
			model.ShortURL = k
			flag = true
			break
		}
	}
	if !flag {
		err := errors.New("URL not found")
		return model, err
	}
	return model, nil
}
