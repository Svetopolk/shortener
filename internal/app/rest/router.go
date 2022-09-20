package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Svetopolk/shortener/internal/logging"
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
		r.Delete("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
			logging.Enter()
			defer logging.Exit()

			m.batchDelete(w, r)
		})
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			logging.Enter()
			defer logging.Exit()

			m.handlePing(w, r)
		})
		r.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
			m.handleBatch(w, r)
		})
	})

	return r
}
