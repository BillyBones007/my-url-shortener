package db

import (
	"errors"

	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

// Вместо базы данных пока используется мапа,
// где ключи - это короткий url,
// значения - длинной url.
func init() {
	DataBase = make(map[string]string)
}

var DataBase map[string]string

// Интерфейс для работы с некой базой данных
type DBase interface {
	InsertUrl(longUrl string, h hasher.UrlHasher) error
	SelectLongUrl(shortUrl string) (longUrl string, err error)
	SelectShortUrl(longUrl string) (shortUrl string, err error)
	UrlIsExist(longUrl string) bool
}

// Тип для работы с мапой
type DB struct{}

// Проверяет, существует ли длинный url в базе
func (d DB) UrlIsExist(lurl string) bool {
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
func (d DB) InsertUrl(lurl string, h hasher.UrlHasher) error {
	// На случай если GetHash выдаст shortUrl, который уже есть
	// в мапе, используем цикл до тех пор, пока не получим уникальное значение
	for {
		shortUrl := h.GetHash(lurl)
		if DataBase[shortUrl] == "" {
			DataBase[shortUrl] = lurl
			break
		}
	}
	return nil
}

// Возвращает длинный url из мапы на основе короткого url
func (d DB) SelectLongUrl(shortUrl string) (longUrl string, err error) {
	shortUrl = "http://localhost:8080" + shortUrl
	if DataBase[shortUrl] == "" {
		longUrl = ""
		err = errors.New("URL not found")
		return longUrl, err
	}
	longUrl = DataBase[shortUrl]
	return longUrl, nil
}

// Возвращает короткий url из мапы на основе длинного url
func (d DB) SelectShortUrl(longUrl string) (shortUrl string, err error) {
	flag := false
	for k, v := range DataBase {
		if v == longUrl {
			shortUrl = k
			flag = true
			break
		}
	}
	if !flag {
		err = errors.New("URL not found")
		shortUrl = ""
		return shortUrl, nil
	}
	return shortUrl, err
}
