package signcookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"log"
)

// Подписывает uuid ключом key
func SignHMAC(uuid []byte, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(uuid)
	dst := h.Sum(uuid)
	return dst
}

// Проверяет подпись на подлинность
func CheckSignHMAC(data []byte, key []byte) bool {
	h := hmac.New(sha256.New, key)
	defer func() {
		if err := recover(); err != nil {
			log.Printf("ERROR: %s\n", err)
		}
	}()
	h.Write(data[:16])
	sign := h.Sum(nil)
	return hmac.Equal(sign, data[16:])
}
