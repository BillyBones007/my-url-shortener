package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BillyBones007/my-url-shortener/internal/db/models"
	"github.com/BillyBones007/my-url-shortener/internal/middlewares"
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
	fmt.Println(model)
	mURL, err := h.Storage.SelectLongURL(&model)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	} else {
		fmt.Println(mURL)
		rw.Header().Set("Location", mURL.LongURL)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	}
}

// Возвращает все когда-либо сокращенные пользователем url в json формате
func (h *Handler) GetAllURLsHandler(rw http.ResponseWriter, r *http.Request) {
	type OutObj struct {
		Result []models.Model
	}
	id := r.Context().Value(middlewares.KeyUUID)
	uuid := fmt.Sprintf("%v", id)
	fmt.Printf("INFO: value from context in GetAllURLsHandler: %s\n", uuid)
	list, _ := h.Storage.SelectAllForUUID(uuid)
	if list == nil {
		rw.WriteHeader(http.StatusNoContent)
	} else {
		out := OutObj{Result: list}
		b, err := json.Marshal(out)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}

		rw.Header().Set("Content-Type", "application/json; charset=utf-8")
		rw.WriteHeader(http.StatusCreated)
		rw.Write(b)
	}

}
