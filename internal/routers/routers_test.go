package routers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type returnData struct {
	statusCode int
	body       string
	location   string
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
		req, err := http.NewRequest(http.MethodGet, rd.body, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		defer resp.Body.Close()
		require.NoError(t, err)

		rd.statusCode = resp.StatusCode
		rd.location = resp.Header.Get("Location")
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	ts := httptest.NewServer(r)
	defer ts.Close()

	longURL := "https://habr.com/ru/post/702374/"

	testRequest(t, ts, "POST", "", longURL, &rd)
	assert.Equal(t, http.StatusCreated, rd.statusCode)

	testRequest(t, ts, "GET", rd.body, "", &rd)
	assert.Equal(t, http.StatusTemporaryRedirect, rd.statusCode)
	assert.Equal(t, longURL, rd.location)

}
