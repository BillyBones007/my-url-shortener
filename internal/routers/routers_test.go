package routers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BillyBones007/my-url-shortener/internal/db/maps"
	"github.com/BillyBones007/my-url-shortener/internal/hasher/rand"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var m *maps.MapStorage = maps.NewStorage()
var h rand.URLHash

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

type returnData struct {
	sCodes     []int
	locations  []string
	statusCode int
	body       string
}

var rd returnData

func testRequest(t *testing.T, ts *httptest.Server, method, endPoint string, body string, rd *returnData) {
	switch method {
	case "POST":
		req, err := http.NewRequest(http.MethodPost, ts.URL+endPoint, strings.NewReader(body))
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		bodyResp, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		defer resp.Body.Close()

		rd.statusCode = resp.StatusCode
		rd.body = string(bodyResp)

	case "GET":
		client := &http.Client{Transport: LogRedirects{}}
		req, err := http.NewRequest(http.MethodGet, rd.body, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		defer resp.Body.Close()
		require.NoError(t, err)

	}
}

func TestRouter(t *testing.T) {
	r := NewRouter(m, h)
	ts := httptest.NewServer(r)
	defer ts.Close()

	longURL := "https://habr.com/ru/post/702374/"

	testRequest(t, ts, "POST", "", longURL, &rd)
	assert.Equal(t, http.StatusCreated, rd.statusCode)

	testRequest(t, ts, "GET", rd.body, "", &rd)
	assert.Equal(t, http.StatusTemporaryRedirect, rd.sCodes[0])
	assert.Equal(t, longURL, rd.locations[1])

}
