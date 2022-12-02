package handlers

import (
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

// Обработчик POST запросов. Принимает длинный URL, проверяет его на валидность,
// помещает его в базу данных и выдает в теле ответа короткий URL
func CreateShortURLHandler(rw http.ResponseWriter, r *http.Request) {
	dBase := db.DB{}
	hash := hasher.URLHash{}
	requestHost := r.Host
	if r.URL.Scheme == "" {
		requestHost = "http://" + requestHost
	}

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
	if !dBase.URLIsExist(string(body)) {
		dBase.InsertURL(string(body), hash)
	}

	sURL, err := dBase.SelectShortURL(string(body))
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	sURL = requestHost + sURL
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(sURL))
}

// Обработчик GET запросов. Проверяет полученный короткий URL в базе,
// достает на основе его оригинальный URL и делает по нему редирект
func GetLongURLHandler(rw http.ResponseWriter, r *http.Request) {
	dBase := db.DB{}
	q := r.URL.EscapedPath()
	if q == "/" {
		http.Error(rw, "The query parameter is missing", http.StatusBadRequest)
		return
	}
	lURL, err := dBase.SelectLongURL(q)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	rw.Header().Set("Location", lURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

// Обработчик не подлежащих обработке запросов
func BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}

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
