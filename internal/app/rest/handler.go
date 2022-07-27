package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/util"
)

type RequestHandler struct {
	service service.ShortService
}

const localAddress = "http://localhost:8080/"

func NewRequestHandler(storage *service.ShortService) *RequestHandler {
	return &RequestHandler{*storage}
}

func (h *RequestHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	hash := h.service.Save(string(body))
	shortURL := makeShortUrl(hash)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeShortUrl(hash string) string {
	return localAddress + hash
}

func (h *RequestHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	hash := util.RemoveFirstSymbol(r.URL.Path)
	fullURL := h.service.Get(hash)

	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err := w.Write([]byte("redirect to " + fullURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) handleJsonPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	value := Request{}
	if err := json.Unmarshal(resBody, &value); err != nil {
		log.Fatal("can not unmarshal body:[", string(resBody), "] ", err)
	}
	hash := h.service.Save(value.Url)
	response := Response{makeShortUrl(hash)}
	responseString, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(responseString)
	if err != nil {
		panic(err)
	}
}

type Request struct {
	Url string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}
