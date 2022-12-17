package handlers

import (
	"net/http"

	"github.com/BillyBones007/my-url-shortener/internal/db/models"
)

// Обработчик GET запросов. Проверяет полученный короткий URL в базе,
// достает на основе его оригинальный URL и делает по нему редирект
func (h *Handler) GetLongURLHandler(rw http.ResponseWriter, r *http.Request) {
	q := r.URL.EscapedPath()
	if q == "/" {
		http.Error(rw, "The query parameter is missing", http.StatusBadRequest)
		return
	}
	model := models.Model{ShortURL: q}
	mURL, err := h.Storage.SelectLongURL(&model)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	rw.Header().Set("Location", mURL.LongURL)
	rw.WriteHeader(http.StatusTemporaryRedirect)
}
