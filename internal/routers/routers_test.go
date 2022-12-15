package routers

import (
	"bytes"
	"encoding/json"
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

type payloadJSON struct {
	Url string `json:"url"`
}

var rd redirect
var m *maps.MapStorage = maps.NewStorage()
var h randchars.URLHash

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
func testPostJSONRequest(t *testing.T, ts *httptest.Server, endPoint string, body []byte) returnData {
	req, err := http.NewRequest(http.MethodPost, ts.URL+endPoint, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
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

func TestRouter(t *testing.T) {
	r := NewRouter(m, h, "")
	ts := httptest.NewServer(r)
	fmt.Printf("Listener addr: %s\n", ts.Listener.Addr())
	fmt.Printf("URL test server: %s\n", ts.URL)
	defer ts.Close()

	longURL := "https://habr.com/ru/post/702374/"

	retPostData := testPostRequest(t, ts, "/", longURL)
	assert.Equal(t, http.StatusCreated, retPostData.statusCode)
	fmt.Println(string(retPostData.body))
	testGetRequest(t, ts, string(retPostData.body))
	assert.Equal(t, http.StatusTemporaryRedirect, rd.sCodes[0])
	assert.Equal(t, longURL, rd.locations[1])

	payload := payloadJSON{Url: "https://habr.com/ru/post/702373"}
	body, _ := json.Marshal(payload)
	retPostJSONData := testPostJSONRequest(t, ts, "/api/shorten", body)
	fmt.Printf("Return PostJSON: %v\n", string(retPostJSONData.body))
	assert.Equal(t, http.StatusCreated, retPostJSONData.statusCode)
}
