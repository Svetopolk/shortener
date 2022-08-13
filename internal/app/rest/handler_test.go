package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/stretchr/testify/assert"
)

func TestStatusHandler(t *testing.T) {
	h := NewRequestHandler(
		service.NewShortService(storage.NewTestStorage()),
		"http://localhost:8080",
		db.NewDB(""),
	)
	type want struct {
		code        int
		response    string
		contentType string
	}

	type request struct {
		method string
		path   string
		body   string
	}

	type test struct {
		name    string
		request request
		want    want
	}

	tests := []test{
		{
			name: "GET missed value",
			request: request{
				method: http.MethodGet,
				path:   "/",
				body:   "",
			},
			want: want{
				code:     307,
				response: `redirect to `,
			},
		},
		{
			name: "GET value",
			request: request{
				method: http.MethodGet,
				path:   "/12345",
				body:   "",
			},
			want: want{
				code:     307,
				response: `redirect to https://ya.ru`,
			},
		},
		{
			name: "POST url",
			request: request{
				method: http.MethodPost,
				path:   "/",
				body:   "https://ya.ru",
			},
			want: want{
				code:     201,
				response: `http://localhost:8080/12345`,
			},
		},
		{
			name: "POST url",
			request: request{
				method: http.MethodPost,
				path:   "/api/shorten",
				body:   `{"url":"https://ya.ru"}`,
			},
			want: want{
				code:     201,
				response: `{"result":"http://localhost:8080/12345"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.request.method, tt.request.path, strings.NewReader(tt.request.body))

			w := httptest.NewRecorder()
			switch tt.request.method {
			case http.MethodGet:
				h.handleGet(w, request)
			case http.MethodPost:
				if tt.request.path == "/" {
					h.handlePost(w, request)
				} else if tt.request.path == "/api/shorten" {
					h.handleJSONPost(w, request)
				} else {
					panic("unexpected path:" + tt.request.path)
				}
			}

			res := w.Result()
			assert.Equal(t, tt.want.code, res.StatusCode)
			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			resBody := readBody(t, res)
			if resBody != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}

		})
	}
}

func readBody(t *testing.T, res *http.Response) string {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.Fatal(err)
		}
	}(res.Body)
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return string(resBody)
}
