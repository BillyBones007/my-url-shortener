package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/BillyBones007/my-url-shortener/internal/db/models"
)

// Обработчик POST запросов. Принимает длинный URL в теле запроса, проверяет его на валидность,
// помещает его в базу данных и выдает в теле ответа короткий URL
func (h *Handler) CreateShortURLHandler(rw http.ResponseWriter, r *http.Request) {
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
	if !h.Storage.URLIsExist(&model) {
		h.Storage.InsertURL(&model, h.Hasher)
	}

	mURL, err := h.Storage.SelectShortURL(&model)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// Если переменная окружения $BASE_URL не была установлена,
	// то конструируем конечный url из текущего адреса сервера
	bURL := h.BaseURL
	if bURL == "" {
		bURL = "http://" + r.Host
	}
	mURL.ShortURL = bURL + mURL.ShortURL
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(mURL.ShortURL))
}

// Обработчик POST-запроса. Принимает в теле запроса JSON объект {"url": "<some_url>"}
// Возвращает в ответ объект {"result": "<shorten_url>"}
func (h *Handler) CreateInJSONShortURLHandler(rw http.ResponseWriter, r *http.Request) {
	type InObj struct {
		URL string `json:"url"`
	}

	type OutObj struct {
		Result string `json:"result"`
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	var inObj InObj
	err = json.Unmarshal(body, &inObj)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	if !urlValid(inObj.URL) {
		http.Error(rw, "Url incorrected", http.StatusBadRequest)
		return
	}
	model := models.Model{LongURL: inObj.URL}
	if !h.Storage.URLIsExist(&model) {
		h.Storage.InsertURL(&model, h.Hasher)
	}

	mURL, err := h.Storage.SelectShortURL(&model)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	// Если переменная окружения $BASE_URL не была установлена,
	// то конструируем конечный url из текущего адреса сервера
	bURL := h.BaseURL
	if bURL == "" {
		bURL = "http://" + r.Host
	}
	outObj := OutObj{Result: bURL + mURL.ShortURL}
	b, err := json.Marshal(outObj)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusCreated)
	rw.Write(b)
}
