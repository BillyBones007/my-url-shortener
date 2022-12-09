package maps

import (
	"errors"

	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

// Тип для работы с мапой
type MapStorage struct {
	DataBase map[string]string
}

// Конструктор хранилища. Возвращает указатель на MapStorage
func NewStorage() *MapStorage {
	return &MapStorage{DataBase: make(map[string]string)}
}

// Проверяет, существует ли длинный url в базе
func (m *MapStorage) URLIsExist(model models.Model) bool {
	flag := false
	for _, long := range m.DataBase {
		if long == model.LongURL {
			flag = true
			break
		}
	}
	return flag
}

// Заполняет мапу. Получает длинный и короткий url
func (m *MapStorage) InsertURL(lurl string, h hasher.URLHasher) error {
	// На случай если GetHash выдаст shortUrl, который уже есть
	// в мапе, используем цикл до тех пор, пока не получим уникальное значение
	for {
		shortURL := h.GetHash(lurl)
		if m.DataBase[shortURL] == "" {
			m.DataBase[shortURL] = lurl
			break
		}
	}
	return nil
}

// Возвращает длинный url из мапы на основе короткого url
func (m *MapStorage) SelectLongURL(shortURL string) (model models.Model, err error) {
	if m.DataBase[shortURL] == "" {
		model.LongURL = ""
		model.ShortURL = shortURL
		err = errors.New("URL not found")
		return model, err
	}
	model.LongURL = m.DataBase[shortURL]
	model.ShortURL = shortURL
	return model, nil
}

// Возвращает короткий url из мапы на основе длинного url
func (m *MapStorage) SelectShortURL(longURL string) (model models.Model, err error) {
	flag := false
	for k, v := range m.DataBase {
		if v == longURL {
			model.ShortURL = k
			model.LongURL = longURL
			flag = true
			break
		}
	}
	if !flag {
		err = errors.New("URL not found")
		model.ShortURL = ""
		model.LongURL = longURL
		return model, err
	}
	return model, nil
}
