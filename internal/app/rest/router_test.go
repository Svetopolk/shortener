package rest

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func TestContentEncodingGzip(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	requestBody := zip("https://ya.ru")
	resp, body := testRequest(t, ts, "POST", "/", requestBody, "Content-Encoding", "gzip")

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, "http://localhost:8080/12345", body)
	closeBody(t, resp)
}

func TestGetAllCookiePresent(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/api/user/urls", strings.NewReader(""))
	reqCookie := http.Cookie{Name: "userID", Value: "3456", Expires: time.Now().Add(time.Hour)}
	req.AddCookie(&reqCookie)
	resp, _ := sendRequest(t, req)

	cookies := resp.Cookies()
	assert.Equal(t, 1, len(cookies))
	cookie := cookies[0]
	assert.Equal(t, "userID", cookie.Name)
	assert.Equal(t, 72, len(cookie.Value))
	closeBody(t, resp)

}

func TestGetAllCookieMissed(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, _ := testRequest(t, ts, "GET", "/api/user/urls", "")

	cookies := resp.Cookies()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, 1, len(cookies))
	cookie := cookies[0]
	assert.Equal(t, "userID", cookie.Name)
	assert.Equal(t, 72, len(cookie.Value))
	closeBody(t, resp)
}

func TestPingDb(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	resp, _ := testRequest(t, ts, "GET", "/ping", "")

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	closeBody(t, resp)
}

func TestUserIDCookiePresent(t *testing.T) {
	ts := getServer()
	defer ts.Close()

	userID := getSignedUserID()
	req, _ := http.NewRequest("GET", ts.URL+"/api/user/urls", strings.NewReader(""))
	reqCookie := http.Cookie{Name: "userID", Value: userID, Expires: time.Now().Add(time.Hour)}
	req.AddCookie(&reqCookie)
	resp, _ := sendRequest(t, req)

	cookies := resp.Cookies()
	assert.Equal(t, 1, len(cookies))
	cookie := cookies[0]
	assert.Equal(t, "userID", cookie.Name)
	assert.Equal(t, userID, cookie.Value)
	closeBody(t, resp)
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body string, headers ...string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, strings.NewReader(body))
	if len(headers) == 2 {
		req.Header.Set(headers[0], headers[1])
	}
	require.NoError(t, err)

	return sendRequest(t, req)
}

func sendRequest(t *testing.T, req *http.Request) (*http.Response, string) {
	client := http.DefaultClient
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

func unzip(original string) string {
	reader := bytes.NewReader([]byte(original))
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		log.Println("error while unzip", err)
		return ""
	}
	output, err := ioutil.ReadAll(gzReader)
	if err != nil {
		log.Println("error while unzip", err)
		return ""
	}
	return string(output)
}

func zip(original string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(original)); err != nil {
		log.Println(err)
	}
	if err := gz.Close(); err != nil {
		log.Println(err)
	}
	return b.String()
}

func closeBody(t *testing.T, resp *http.Response) {
	err := resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func getServer() *httptest.Server {
	r := NewRouter(NewRequestHandler(
		service.NewShortService(storage.NewTestStorage()),
		"http://localhost:8080",
		nil,
	))
	ts := httptest.NewServer(r)
	return ts
}
