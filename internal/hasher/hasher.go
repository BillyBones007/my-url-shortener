package hasher

import (
	"math/rand"
	"time"
)

const characters = "aAbBcCdDeEfFgGhHiIjJkKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

// Интерфейс для работы с неким шифровщиком ссылок
type URLHasher interface {
	GetHash(longURL string) string
}

// Тип для работы с шифрование ссылки
type URLHash struct {
	LongURL string
}

// Возвращает короткую ссылку
func (h URLHash) GetHash(longURL string) string {
	rand.Seed(time.Now().UnixNano())
	shortURL := make([]byte, 6)
	for i := range shortURL {
		shortURL[i] = characters[rand.Intn(len(characters))]
	}
	s := string(shortURL)
	return "/" + s
}
