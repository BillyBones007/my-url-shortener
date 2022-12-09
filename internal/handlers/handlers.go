package handlers

import (
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

type Handler struct {
	Storage db.DBase
	Hasher  hasher.URLHasher
}

// Обработчик POST запросов. Принимает длинный URL, проверяет его на валидность,
// помещает его в базу данных и выдает в теле ответа короткий URL
func (h *Handler) CreateShortURLHandler(rw http.ResponseWriter, r *http.Request) {
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
	model := models.Model{LongURL: string(body)}
	if !h.Storage.URLIsExist(model) {
		h.Storage.InsertURL(model.LongURL, h.Hasher)
	}

	mURL, err := h.Storage.SelectShortURL(model.LongURL)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	mURL.ShortURL = requestHost + mURL.ShortURL
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(mURL.ShortURL))
}

// Обработчик GET запросов. Проверяет полученный короткий URL в базе,
// достает на основе его оригинальный URL и делает по нему редирект
func (h *Handler) GetLongURLHandler(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.EscapedPath()
	if q == "/" {
		http.Error(rw, "The query parameter is missing", http.StatusBadRequest)
		return
	}
	mURL, err := h.Storage.SelectLongURL(q)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	rw.Header().Set("Location", mURL.LongURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}

// Обработчик не подлежащих обработке запросов
func (h *Handler) BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
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
