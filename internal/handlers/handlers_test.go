package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShortUrlHandlerPost(t *testing.T) {
	type want struct {
		sCode int
	}

	tests := []struct {
		name    string
		longUrl string
		want    want
	}{
		{
			name:    "PostOK",
			longUrl: "http://habr.com/blablabla",
			want: want{
				sCode: 201,
			},
		},
		{
			name:    "PostBadUrl",
			longUrl: "blablabla",
			want: want{
				sCode: 400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.longUrl))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(ShortUrlHandler)
			h.ServeHTTP(w, request)
			resp := w.Result()

			assert.Equal(t, tt.want.sCode, resp.StatusCode, fmt.Sprintf("Статус код должен быть %d", tt.want.sCode))

		})

	}
}

func TestShortUrlHandlerGet(t *testing.T) {
	type want struct {
		sCode    int
		location string
	}

	tests := []struct {
		name    string
		longUrl string
		want    want
	}{
		{
			name:    "Get_OK",
			longUrl: "http://habr.com/blablabla",
			want: want{
				sCode:    307,
				location: "http://habr.com/blablabla",
			},
		},
		{
			name:    "GetBadUrl",
			longUrl: "blablabla",
			want: want{
				sCode:    400,
				location: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// POST запрос и сохранение полученного короткого url в sUrl
			reqPost := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.longUrl))
			w := httptest.NewRecorder()
			h := http.HandlerFunc(ShortUrlHandler)
			h.ServeHTTP(w, reqPost)
			resp := w.Result()
			sUrl := string(w.Body.String())
			// Для себя посмотреть, что получили
			fmt.Printf("\nResponse POST: %v\n", resp)
			fmt.Printf("ShortUrl: %s\n\n", sUrl)

			// При некорректном заданном url в теле ответа вернется сообщение об ошибке
			// "Url incorrected" и нельзя будет передать sUrl в GET запрос
			// Поэтому проверяем статус код и либо передаем полученный в POST запросе
			// sUrl как есть, либо конструируем GET запрос с добавлением хоста
			if resp.StatusCode != 400 {
				reqGet := httptest.NewRequest(http.MethodGet, sUrl, nil)
				w = httptest.NewRecorder()
				h.ServeHTTP(w, reqGet)
				resp = w.Result()
				loc := resp.Header.Get("Location")
				fmt.Printf("\nResponse GET: %v\n", resp)
				fmt.Printf("Location: %v\n\n", loc)
				assert.Equal(t, tt.want.location, loc, fmt.Sprintf("Location не должен быть nil: %s", loc))
				assert.Equal(t, tt.want.sCode, resp.StatusCode, fmt.Sprintf("Статус код должен быть %d", tt.want.sCode))
			} else {
				reqGet := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:8080/%s", tt.longUrl), nil)
				w = httptest.NewRecorder()
				h.ServeHTTP(w, reqGet)
				resp = w.Result()
				loc := resp.Header.Get("Location")
				assert.Equal(t, tt.want.location, loc, fmt.Sprintf("Location не должен быть пуст: %s", loc))
				assert.Equal(t, tt.want.sCode, resp.StatusCode, fmt.Sprintf("Статус код должен быть %d", tt.want.sCode))
			}

		})

	}
}
