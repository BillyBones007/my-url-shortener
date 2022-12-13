package routers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BillyBones007/my-url-shortener/internal/db/maps"
	"github.com/BillyBones007/my-url-shortener/internal/hasher/randchars"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type returnData struct {
	statusCode int
	body       []byte
}

type redirect struct {
	sCodes    []int
	locations []string
}

var rd redirect
var m *maps.MapStorage = maps.NewStorage()
var h randchars.URLHash
var bu string = "http://localhost:8080"

type LogRedirects struct {
	Transport http.RoundTripper
}

func (l LogRedirects) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	RData := &rd
	t := l.Transport
	if t == nil {
		t = http.DefaultTransport
	}
	resp, err = t.RoundTrip(req)
	if err != nil {
		return
	}
	RData.sCodes = append(RData.sCodes, resp.StatusCode)
	RData.locations = append(RData.locations, req.URL.String())
	fmt.Printf("INFO redirect to: %v, status code: %d, location: %v\n", req.URL, resp.StatusCode, resp.Header.Get("Location"))
	return resp, err
}

// Вспомогательная функция тестирования POST запроса (длинный url в теле запроса в виде строки)
func testPostRequest(t *testing.T, ts *httptest.Server, endPoint string, body string) returnData {
	req, err := http.NewRequest(http.MethodPost, ts.URL+endPoint, strings.NewReader(body))
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	bodyResp, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	var retPostData = returnData{statusCode: resp.StatusCode, body: bodyResp}

	return retPostData
}

// Вспомогательная функция тестирования POST запроса (длинный url в теле запроса в виде json объекта)
func testPostJSONRequest(t *testing.T, ts *httptest.Server, endPoint string, body string) returnData {
	client := &http.Client{}
	// req, err := http.NewRequest(http.MethodPost, ts.URL+endPoint, strings.NewReader(body))
	// require.NoError(t, err)
	resp, err := client.Post(ts.URL+endPoint, "application/json", strings.NewReader(body))
	// resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	bodyResp, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	var retPostData = returnData{statusCode: resp.StatusCode, body: bodyResp}

	return retPostData
}

func testGetRequest(t *testing.T, ts *httptest.Server, sURL string) {
	client := &http.Client{Transport: LogRedirects{}}
	req, err := http.NewRequest(http.MethodGet, sURL, nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
}

// func testRequest(t *testing.T, ts *httptest.Server, method, endPoint string, body string, rd *returnData) {
// 	switch method {
// 	case "POST":
// 		req, err := http.NewRequest(http.MethodPost, ts.URL+endPoint, strings.NewReader(body))
// 		require.NoError(t, err)

// 		resp, err := http.DefaultClient.Do(req)
// 		require.NoError(t, err)

// 		bodyResp, err := io.ReadAll(resp.Body)
// 		require.NoError(t, err)

// 		defer resp.Body.Close()

// 		rd.statusCode = resp.StatusCode
// 		rd.body = bodyResp

// 	case "GET":
// 		client := &http.Client{Transport: LogRedirects{}}
// 		req, err := http.NewRequest(http.MethodGet, string(rd.body), nil)
// 		require.NoError(t, err)

// 		resp, err := client.Do(req)
// 		defer resp.Body.Close()
// 		require.NoError(t, err)

// 	}
// }

func TestRouter(t *testing.T) {
	r := NewRouter(m, h, bu)
	ts := httptest.NewServer(r)
	defer ts.Close()

	longURL := "https://habr.com/ru/post/702374/"

	retPostData := testPostRequest(t, ts, "/", longURL)
	assert.Equal(t, http.StatusCreated, retPostData.statusCode)

	testGetRequest(t, ts, string(retPostData.body))
	assert.Equal(t, http.StatusTemporaryRedirect, rd.sCodes[0])
	assert.Equal(t, longURL, rd.locations[1])

	retPostJsonData := testPostJSONRequest(t, ts, "/api/shorten", "{'url': 'https://habr.com/ru/post/702373'}")
	assert.Equal(t, http.StatusCreated, retPostJsonData.statusCode)
}
