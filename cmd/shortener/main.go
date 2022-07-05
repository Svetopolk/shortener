package main

import (
	"io"
	"log"
	"net/http"
)

type RequestHandler struct {
	storage Storage
}

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s request %s", r.Method, r.URL)

	if r.Method == http.MethodGet {
		hash := trimFirstRune(r.URL.Path)
		fullUrl := h.storage.get(hash)

		w.Header().Set("Location", fullUrl)
		w.WriteHeader(http.StatusTemporaryRedirect)
		_, err := w.Write([]byte("redirect to " + fullUrl))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hash := h.storage.save(string(body))
		shortUrl := "localhost:8080/" + hash
		_, err = w.Write([]byte(shortUrl))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	m := RequestHandler{}
	http.Handle("/", &m)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
