package rest

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPositive(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/12345", "")

	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "https://ya.ru", resp.Header.Get("Location"))
	assert.Equal(t, "redirect to https://ya.ru", body)
	closeBody(t, resp)
}

func TestGetEmpty(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "GET", "/98765", "")

	assert.Equal(t, http.StatusTemporaryRedirect, resp.StatusCode)
	assert.Equal(t, "", resp.Header.Get("Location"))
	assert.Equal(t, "redirect to ", body)
	closeBody(t, resp)
}

func TestPostPositive(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/", "https://ya.ru")

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/12345", body)
	closeBody(t, resp)
}

func TestPostBadRequest(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/", "")

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "", body)
	closeBody(t, resp)
}

func TestRouterNotFound(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/1/2", "https://ya.ru")

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "404 page not found\n", body)
	closeBody(t, resp)
}

func TestPostApi(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/api/shorten", `{"url":"https://ya.ru"}`)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, `{"result":"http://localhost:8080/12345"}`, body)
	closeBody(t, resp)
}

func TestAcceptEncodingGzip(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, body := testRequest(t, ts, "POST", "/", "https://ya.ru", "Accept-Encoding", "gzip")

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/12345", unzip(body))
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))
	closeBody(t, resp)
}

//TODO fix test
func _TestContentEncodingGzip(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	requestBody := zip("https://ya.ru")
	//requestBody := "https://ya.ru"
	resp, body := testRequest(t, ts, "POST", "/", requestBody, "Content-Encoding", "gzip")

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/12345", body)
	closeBody(t, resp)

	//resp, body := testRequest(t, ts, "POST", "/", zip("https://ya.ru"), "Content-Encoding", "gzip")

}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string, headers ...string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	if len(headers) == 2 {
		req.Header.Set(headers[0], headers[1])
	}
	require.NoError(t, err)

	client := http.DefaultClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(resp.Body)

	return resp, string(respBody)
}

func unzip(original string) string {
	reader := bytes.NewReader([]byte(original))
	gzReader, e := gzip.NewReader(reader)
	if e != nil {
		log.Fatal(e)
	}
	output, e := ioutil.ReadAll(gzReader)
	if e != nil {
		log.Fatal(e)
	}
	return string(output)
}

func zip(original string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(original)); err != nil {
		log.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		log.Fatal(err)
	}
	return string(b.Bytes())
}

func closeBody(t *testing.T, resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func getServer() *httptest.Server {
	r := NewRouter(NewRequestHandler(service.NewShortService(storage.NewTestStorage()), "http://localhost:8080"))
	ts := httptest.NewServer(r)
	return ts
}
