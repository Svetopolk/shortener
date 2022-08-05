package rest

import (
	"compress/gzip"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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
	ReadCloser io.ReadCloser
}

func (g gzipRequestBody) Read(p []byte) (n int, err error) {
	return g.ReadCloser.Read(p)
}

func (g gzipRequestBody) Close() error {
	return g.ReadCloser.Close()
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
		r.Body = gzipRequestBody{ReadCloser: gz}
		next.ServeHTTP(w, r)
	})
}

const userIdCookieName = "userId"

func userIdCookieHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userId, err := r.Cookie(userIdCookieName)
		if err != nil {
			log.Print("userId cookie not found")
			expiration := time.Now().Add(365 * 24 * time.Hour)
			value := "11111"
			cookie := http.Cookie{Name: userIdCookieName, Value: value, Expires: expiration}
			http.SetCookie(w, &cookie)
		}
		log.Printf("userId cookie %v", userId)

		next.ServeHTTP(w, r)
	})
}
