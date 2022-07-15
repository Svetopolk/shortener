package rest

import (
	"io"
	"net/http"

	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/util"
)

type RequestHandler struct {
	storage service.ShortService
}

func NewRequestHandler(storage service.ShortService) *RequestHandler {
	return &RequestHandler{storage}
}

func (h *RequestHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash := h.storage.Save(string(body))
	shortURL := "http://localhost:8080/" + hash
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	hash := util.RemoveFirstSymbol(r.URL.Path)
	fullURL := h.storage.Get(hash)

	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err := w.Write([]byte("redirect to " + fullURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
