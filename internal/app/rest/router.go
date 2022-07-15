package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(m *RequestHandler) chi.Router {
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
