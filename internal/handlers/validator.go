package handlers

import (
	"log"
	"regexp"
)

// Валидация полученных ссылок
func urlValid(recURL string) bool {
	pattern := `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
	flag := false
	matched, err := regexp.Match(pattern, []byte(recURL))
	if err != nil {
		log.Printf("Ошибка возникла при вызове regexp.Match: %s", err)
		return flag
	}
	if !matched {
		log.Print("Ошибка валидации: Url не подходит под паттерн.")
		return flag
	}
	flag = true
	return flag
}
