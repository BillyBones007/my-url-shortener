package db

import (
	"errors"

	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

// Вместо базы данных пока используется мапа,
// где ключи - это короткий url,
// значения - длинный url.
func init() {
	DataBase = make(map[string]string)
}

var DataBase map[string]string

// Интерфейс для работы с некой базой данных
type DBase interface {
	InsertURL(longURL string, h hasher.URLHasher) error
	SelectLongURL(shortURL string) (longURL string, err error)
	SelectShortURL(longURL string) (shortURL string, err error)
	URLIsExist(longURL string) bool
}

// Тип для работы с мапой
type DB struct{}

// Проверяет, существует ли длинный url в базе
func (d DB) URLIsExist(lurl string) bool {
	flag := false
	for _, long := range DataBase {
		if long == lurl {
			flag = true
			break
		}
	}
	return flag
}

// Заполняет мапу. Получает длинный и короткий url
func (d DB) InsertURL(lurl string, h hasher.URLHasher) error {
	// На случай если GetHash выдаст shortUrl, который уже есть
	// в мапе, используем цикл до тех пор, пока не получим уникальное значение
	for {
		shortURL := h.GetHash(lurl)
		if DataBase[shortURL] == "" {
			DataBase[shortURL] = lurl
			break
		}
	}
	return nil
}

// Возвращает длинный url из мапы на основе короткого url
func (d DB) SelectLongURL(shortURL string) (longURL string, err error) {
	if DataBase[shortURL] == "" {
		longURL = ""
		err = errors.New("URL not found")
		return longURL, err
	}
	longURL = DataBase[shortURL]
	return longURL, nil
}

// Возвращает короткий url из мапы на основе длинного url
func (d DB) SelectShortURL(longURL string) (shortURL string, err error) {
	flag := false
	for k, v := range DataBase {
		if v == longURL {
			shortURL = k
			flag = true
			break
		}
	}
	if !flag {
		err = errors.New("URL not found")
		shortURL = ""
		return shortURL, err
	}
	return shortURL, nil
}
