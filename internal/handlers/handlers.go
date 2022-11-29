package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

func ShortUrlHandler(rw http.ResponseWriter, r *http.Request) {
	dBase := db.DB{}
	hash := hasher.UrlHash{}
	requestHost := r.Host
	if r.URL.Scheme == "" {
		requestHost = "http://" + requestHost
	}

	switch r.Method {
	case "POST":
		body, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		if !urlValid(string(body)) {
			http.Error(rw, "Url incorrected", http.StatusBadRequest)
			return
		}
		if !dBase.UrlIsExist(string(body)) {
			dBase.InsertUrl(string(body), hash)
		}

		sUrl, err := dBase.SelectShortUrl(string(body))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		sUrl = requestHost + sUrl
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusCreated)
		rw.Write([]byte(sUrl))

	case "GET":
		q := r.URL.EscapedPath()
		if q == "/" {
			http.Error(rw, "The query parameter is missing", http.StatusBadRequest)
			return
		}
		lUrl, err := dBase.SelectLongUrl(q)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
		rw.Header().Set("Location", lUrl)
		rw.WriteHeader(http.StatusTemporaryRedirect)
		fmt.Println(rw.Header())

	default:
		rw.WriteHeader(http.StatusBadRequest)
	}
}

// Валидация полученных ссылок
func urlValid(recUrl string) bool {
	pattern := `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)`
	flag := false
	matched, err := regexp.Match(pattern, []byte(recUrl))
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
