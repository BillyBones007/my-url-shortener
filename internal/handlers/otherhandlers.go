package handlers

import "net/http"

// Обработчик не подлежащих обработке запросов
func (h *Handler) BadRequestHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
}
