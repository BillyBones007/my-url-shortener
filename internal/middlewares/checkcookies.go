package middlewares

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/BillyBones007/my-url-shortener/internal/auth/randbytes"
	"github.com/BillyBones007/my-url-shortener/internal/auth/signcookies"
)

// Тип для ключа используемого при передаче в контексте
type KeyID string

const KeyUUID KeyID = "uuid"

// Содержит рандомный ключ шифрования и вектор инициализации
type Auth struct {
	SecretKey []byte
}

// Мидлварь проверяет наличие у пользователя cookie
func (a *Auth) CheckCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var targetCookie string
		var uuid string
		enc := signcookies.NewEncDecAES(a.SecretKey)
		cookie, err := r.Cookie("uuid")
		if err != nil {
			fmt.Println("Cookie uuid not found")
			encSignUUID, uuid := encryptCookie(enc, a.SecretKey)
			targetCookie = hex.EncodeToString(encSignUUID)
			cookie = &http.Cookie{
				Name:  "uuid",
				Value: targetCookie,
			}
			http.SetCookie(w, cookie)
			r = r.WithContext(context.WithValue(r.Context(), KeyUUID, uuid))
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println("Cookie is exist:")
		fmt.Printf("Name: %s\nValue: %s\n", cookie.Name, cookie.Value)
		valCookie, err := hex.DecodeString(cookie.Value)
		if err != nil {
			log.Printf("ERROR: %s\n", err)
		}
		decCookie := enc.DecryptAES(valCookie)
		if signcookies.CheckSignHMAC(decCookie, a.SecretKey) {
			fmt.Println("INFO: Подпись верна.")
			uuid = hex.EncodeToString(decCookie[:16])
			fmt.Printf("INFO: Дешифрованный uuid: %s\n", uuid)
		} else {
			fmt.Println("INFO: Подпись неверна.")
			fmt.Println("INFO: Установлена новая Cookie.")
			encUUID, uuid := encryptCookie(enc, a.SecretKey)
			targetCookie = hex.EncodeToString(encUUID)
			cookie = &http.Cookie{
				Name:  "uuid",
				Value: targetCookie,
			}
			http.SetCookie(w, cookie)
			r = r.WithContext(context.WithValue(r.Context(), KeyUUID, uuid))
			next.ServeHTTP(w, r)
			return

		}
		r = r.WithContext(context.WithValue(r.Context(), KeyUUID, uuid))
		next.ServeHTTP(w, r)

	})

}

// Вспомогательная функция шифрования cookie
func encryptCookie(enc *signcookies.EncDecAES, key []byte) (encSignUUID []byte, stringUUID string) {
	uuid, err := randbytes.RandomBytes(16)
	if err != nil {
		log.Printf("ERROR: %s\n", err)
	}
	fmt.Printf("INFO: Оригинальный UUID: %x\n", uuid)
	signUUID := signcookies.SignHMAC(uuid, key)
	fmt.Printf("INFO: Подписанный UUID: %x\n", signUUID)
	encSignUUID = enc.EncryptAES(signUUID)
	fmt.Printf("INFO: Шифрованный и подписанный UUID: %x\n", encSignUUID)
	stringUUID = hex.EncodeToString(uuid)
	return encSignUUID, stringUUID
}
