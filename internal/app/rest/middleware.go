package rest

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
)

var _ http.ResponseWriter = gzipWriter{}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipResponseHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncodingHeader := r.Header.Get("Accept-Encoding")
		if !strings.Contains(acceptEncodingHeader, "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

var _ io.ReadCloser = gzipRequestBody{}

type gzipRequestBody struct {
	Reader io.Reader
	Closer io.Closer
}

func (g gzipRequestBody) Read(p []byte) (n int, err error) {
	return g.Reader.Read(p)
}

func (g gzipRequestBody) Close() error {
	return g.Closer.Close()
}

func gzipRequestHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncodingHeader := r.Header.Get("Content-Encoding")
		if !strings.Contains(contentEncodingHeader, "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		log.Println("Encoded request is received")

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer gz.Close()
		r.Body = gzipRequestBody{Reader: gz, Closer: gz}
		next.ServeHTTP(w, r)
	})
}
