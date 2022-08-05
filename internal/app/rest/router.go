package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(m *RequestHandler) chi.Router {
	r := chi.NewRouter()
	r.Use(gzipResponseHandle, gzipRequestHandle, userIdCookieHandle)

	r.Route("/", func(r chi.Router) {
		r.Get("/{hash}", func(w http.ResponseWriter, r *http.Request) {
			_ = chi.URLParam(r, "hash")
			m.handleGet(w, r)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			m.handlePost(w, r)
		})
		r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
			m.handleJSONPost(w, r)
		})
		r.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
			m.getUserUrls(w, r)
		})
	})

	return r
}
