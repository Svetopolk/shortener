package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/app/util"

	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/storage"
)

func TestDbDeleteBatchApi(t *testing.T) {
	ts := getDBServer(t)
	defer ts.Close()

	randomUrl := generateRandomUrl()
	log.Println("randomUrl ===== " + randomUrl)

	resp, body := testRequest(t, ts, "POST", "/api/shorten", getRequest(randomUrl))
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	closeBody(t, resp)

	hash := grabHash(t, body)
	log.Println("hash ===== " + hash)

	// get url
	resp2, _ := testRequest(t, ts, "GET", "/"+hash, "")
	assert.Equal(t, http.StatusTemporaryRedirect, resp2.StatusCode)
	assert.Equal(t, randomUrl, resp2.Header.Get("Location"))
	closeBody(t, resp2)

	// delete
	deleteBody := `["` + hash + `", "3456"]`
	resp2, responseBody2 := testRequest(t, ts, "DELETE", "/api/user/urls", deleteBody)
	assert.Equal(t, http.StatusAccepted, resp2.StatusCode)
	assert.Equal(t, ``, responseBody2)
	closeBody(t, resp2)

	// get url after delete
	resp3, _ := testRequest(t, ts, "GET", "/"+hash, "")
	assert.Equal(t, http.StatusGone, resp3.StatusCode)
	assert.Equal(t, "", resp3.Header.Get("Location"))
	closeBody(t, resp3)
}

func generateRandomUrl() string {
	return "https://" + util.RandomString(6) + ".ru"
}

func getRequest(url string) string {
	requestBody, _ := json.Marshal(Request{url})
	requestBodyString := string(requestBody)
	return requestBodyString
}

func grabHash(t *testing.T, body string) string {
	var response Response
	err := json.Unmarshal([]byte(body), &response)
	if err != nil {
		t.Error("can't unmarshal response body ", body)
	}
	hash := util.GrabHashFromUrl(response.Result)
	return hash
}

func getDBServer(t *testing.T) *httptest.Server {
	db, err := db.NewDB("postgres://shortener:pass@localhost:5432/shortener")
	if err != nil {
		t.Skip("no db connection")
	}
	err = db.Ping()

	if err != nil {
		log.Println("exceptions while ping DB:", err)
		t.Skip("no db connection")
	}

	r := NewRouter(NewRequestHandler(
		service.NewShortService(storage.NewDBStorage(db)),
		"http://localhost:8080",
		nil,
	))
	return httptest.NewServer(r)
}
