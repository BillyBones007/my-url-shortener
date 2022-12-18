package randchars

import (
	"math/rand"
	"time"
)

const characters = "aA1bB2cC3dD4eE5fF6gG7hH8iI9jJ0kKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Тип для работы с шифрованием ссылки
type URLHash struct {
	LongURL string
}

// Возвращает короткую ссылку
func (h URLHash) GetHash(longURL string) string {
	shortURL := make([]byte, 6)
	for i := range shortURL {
		shortURL[i] = characters[rand.Intn(len(characters))]
	}
	s := string(shortURL)
	return "/" + s
}
