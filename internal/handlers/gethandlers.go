package handlers

import "net/http"

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
