package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/BillyBones007/my-url-shortener/internal/db"
	"github.com/BillyBones007/my-url-shortener/internal/hasher"
)

func ShortUrlHandler(rw http.ResponseWriter, r *http.Request) {
	dBase := db.DB{}
	hash := hasher.UrlHash{}
	switch r.Method {
	case "POST":
		body, err := io.ReadAll(r.Body)
		fmt.Println(body)
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
			sUrl, _ := dBase.SelectShortUrl(string(body))
			rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
			rw.WriteHeader(http.StatusCreated)
			rw.Write([]byte(sUrl))
		} else {
			sUrl, err := dBase.SelectShortUrl(string(body))
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
			rw.WriteHeader(http.StatusCreated)
			rw.Write([]byte(sUrl))
		}

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
	flag := false
	parsedUrl, err := url.Parse(recUrl)
	fmt.Println(parsedUrl.Scheme)
	fmt.Println(parsedUrl.Host)
	if err == nil && parsedUrl.Scheme == "localhost" {
		flag = true
		return flag
	}
	if err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != "" {
		flag = true
	}
	return flag
}
