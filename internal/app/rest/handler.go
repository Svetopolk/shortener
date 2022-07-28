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
	baseUrl string
}

func NewRequestHandler(service *service.ShortService, baseUrl string) *RequestHandler {
	return &RequestHandler{
		service: *service,
		baseUrl: baseUrl,
	}
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
	shortURL := h.makeShortURL(hash)
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func (h *RequestHandler) handleJSONPost(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	value := Request{}
	if err := json.Unmarshal(resBody, &value); err != nil {
		log.Fatal("can not unmarshal body:[", string(resBody), "] ", err)
	}
	hash := h.service.Save(value.URL)
	response := Response{h.makeShortURL(hash)}
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

func (h *RequestHandler) makeShortURL(hash string) string {
	return h.baseUrl + "/" + hash
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}
