package main

import (
	"io"
	"net/http"
)

type RequestHandler struct {
	storage Storage
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
	hash := h.storage.save(string(body))
	shortUrl := "http://localhost:8080/" + hash
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortUrl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	hash := trimFirstRune(r.URL.Path)
	fullUrl := h.storage.get(hash)

	w.Header().Set("Location", fullUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err := w.Write([]byte("redirect to " + fullUrl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
