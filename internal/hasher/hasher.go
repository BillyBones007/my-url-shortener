package hasher

import (
	"math/rand"
	"time"
)

const characters = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

// Интерфейс для работы с неким шифровщиком ссылок
type UrlHasher interface {
	GetHash(longUrl string) string
}

// Тип для работы с шифрование ссылки
type UrlHash struct {
	LongUrl string
}

// Возвращает короткую ссылку
func (h UrlHash) GetHash(longUrl string) string {
	rand.Seed(time.Now().UnixNano())
	shortUrl := make([]byte, 6)
	for i, _ := range shortUrl {
		shortUrl[i] = characters[rand.Intn(len(characters))]
	}
	s := string(shortUrl)
	return "/" + s
}
