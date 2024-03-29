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

const userIDCookieName = "userID"

func userIDCookieHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIDCookie, err := r.Cookie(userIDCookieName)
		newUserIDCookie := userIDCookie
		if err != nil {
			newUserIDCookie = generateNewCookie()
		} else {
			_, err2 := decodeID(userIDCookie.Value)
			if err2 != nil {
				newUserIDCookie = generateNewCookie()
			}
		}
		http.SetCookie(w, newUserIDCookie)
		next.ServeHTTP(w, r)
	})
}

func generateNewCookie() *http.Cookie {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	value := getSignedUserID()
	cookie := http.Cookie{Name: userIDCookieName, Value: value, Expires: expiration}
	return &cookie
}
