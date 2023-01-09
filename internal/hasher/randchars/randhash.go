package randchars

import (
	// "math/rand"
	// "time"
	"encoding/hex"

	"github.com/BillyBones007/my-url-shortener/internal/auth/randbytes"
)

// const characters = "aA1bB2cC3dD4eE5fF6gG7hH8iI9jJ0kKlLmMnNoOpPqQrRsStTuUvVwWxXyYzZ"

// func init() {
// 	rand.Seed(time.Now().UnixNano())
// }

// Тип для работы с шифрованием ссылки
type URLHash struct {
	// LongURL string
	LenHash int
}

// Возвращает короткую ссылку
// func (h URLHash) GetHash(longURL string) string {
// shortURL := make([]byte, 6)
// for i := range shortURL {
// shortURL[i] = characters[rand.Intn(len(characters))]
// }
// s := string(shortURL)
// return "/" + s
// }

// Возвращает короткую ссылку
func (h URLHash) GetHash(lenhash int) string {
	b, _ := randbytes.RandomBytes(lenhash)
	shortURL := hex.EncodeToString(b)
	return "/" + shortURL
}
