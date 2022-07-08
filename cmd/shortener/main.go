package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	m := RequestHandler{NewMemStorage()}
	r := NewRouter(m)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func NewRouter(m RequestHandler) chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Get("/{hash}", func(w http.ResponseWriter, r *http.Request) {
			_ = chi.URLParam(r, "hash")
			m.handleGet(w, r)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			m.handlePost(w, r)
		})
	})
	return r
}
