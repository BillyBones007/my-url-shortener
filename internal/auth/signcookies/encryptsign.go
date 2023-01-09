package signcookies

import (
	"crypto/aes"
	"crypto/cipher"
	"log"
)

// ----------------------------------------------------------
// Тип для работы с шифрованием и расшифровкой cookie
type EncDecAES struct {
	AesBlock cipher.Block
	AesGCM   cipher.AEAD
	Nonce    []byte
}

// Конструктор типа EncDecAES
func NewEncDecAES(key []byte) *EncDecAES {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return nil
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
		return nil
	}
	nonce := key[len(key)-aesgcm.NonceSize():]
	return &EncDecAES{AesBlock: aesblock, AesGCM: aesgcm, Nonce: nonce}
}

// Шифрует uuid. Возвращает шифрованное значение.
func (e *EncDecAES) EncryptAES(cookie []byte) []byte {
	dst := e.AesGCM.Seal(nil, e.Nonce, cookie, nil)
	return dst
}

// Дешифрует cookie. Возвращает оригинальный cookie.
func (e *EncDecAES) DecryptAES(encCookie []byte) (src []byte) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("ERROR: %s\n", err)
			src = nil
			return
		}
	}()
	src, err := e.AesGCM.Open(nil, e.Nonce, encCookie, nil)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}
	return src
}

//----------------------------------------------------------
