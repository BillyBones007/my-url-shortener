package handlers

import (
	"fmt"
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
		recUrl := r.FormValue("Url")
		if !urlValid(recUrl) {
			http.Error(rw, "Url incorrected", http.StatusBadRequest)
			return
		}
		if !dBase.UrlIsExist(recUrl) {
			dBase.InsertUrl(recUrl, hash)
			sUrl, _ := dBase.SelectShortUrl(recUrl)
			rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
			rw.WriteHeader(http.StatusCreated)
			fmt.Fprintf(rw, sUrl)
		} else {
			sUrl, err := dBase.SelectShortUrl(recUrl)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
			rw.WriteHeader(http.StatusCreated)
			fmt.Fprintf(rw, sUrl)
		}

	case "GET":
		q := r.URL.EscapedPath()
		if q == "/" {
			http.Error(rw, "The query parameter is missing", http.StatusBadRequest)
			return
		}
		lUrl, err := dBase.SelectLongUrl("http://localhost:8080" + q)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
		rw.Header().Set("Location", lUrl)
		rw.WriteHeader(http.StatusTemporaryRedirect)
	default:
		rw.WriteHeader(http.StatusBadRequest)
	}
}

// Валидация полученных ссылок
func urlValid(recUrl string) bool {
	flag := false
	parsedUrl, err := url.Parse(recUrl)
	if err == nil && parsedUrl.Scheme != "" && parsedUrl.Host != "" {
		flag = true
	}
	return flag
}
