package rest

import (
	"github.com/Svetopolk/shortener/internal/logging"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(m *RequestHandler) chi.Router {
	logging.Enter()
	defer logging.Exit()

	r := chi.NewRouter()
	r.Use(gzipResponseHandle, gzipRequestHandle, userIDCookieHandle)

	r.Route("/", func(r chi.Router) {
		r.Get("/{hash}", func(w http.ResponseWriter, r *http.Request) {
			logging.Enter()
			defer logging.Exit()

			_ = chi.URLParam(r, "hash")
			m.handleGet(w, r)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			logging.Enter()
			defer logging.Exit()

			m.handlePost(w, r)
		})
		r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
			logging.Enter()
			defer logging.Exit()

			m.handleJSONPost(w, r)
		})
		r.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
			logging.Enter()
			defer logging.Exit()

			m.getUserUrls(w, r)
		})
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			logging.Enter()
			defer logging.Exit()

			m.handlePing(w, r)
		})
	})

	return r
}
