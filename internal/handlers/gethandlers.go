package handlers

import (
	"fmt"
	"net/http"

	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/middlewares"
)

// Обработчик GET запросов. Проверяет полученный короткий URL в базе,
// достает на основе его оригинальный URL и делает по нему редирект
func (h *Handler) GetLongURLHandler(rw http.ResponseWriter, r *http.Request) {
	uuid := r.Context().Value(middlewares.KeyUUID)
	fmt.Printf("INFO: value from context in GetLongURLHandler: %s\n", uuid)
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

// Возвращает все когда-либо сокращенные пользователем url в json формате
func (h *Handler) GetAllURLsHandler(rw http.ResponseWriter, r *http.Request) {
	uuid := r.Context().Value(middlewares.KeyUUID)
	fmt.Printf("INFO: value from context in GetAllURLsHandler: %s\n", uuid)

}
