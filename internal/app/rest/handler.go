package rest

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Svetopolk/shortener/internal/logging"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/util"
)

type RequestHandler struct {
	service  *service.ShortService
	baseURL  string
	dbSource *db.Source
}

func NewRequestHandler(service *service.ShortService, baseURL string, db *db.Source) *RequestHandler {
	logging.Enter()
	defer logging.Exit()

	return &RequestHandler{
		service:  service,
		baseURL:  baseURL,
		dbSource: db,
	}
}

func (h *RequestHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	logging.Enter()
	defer logging.Exit()

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
	logging.Enter()
	defer logging.Exit()

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
	logging.Enter()
	defer logging.Exit()

	defer r.Body.Close()
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	value := Request{}
	if err := json.Unmarshal(resBody, &value); err != nil {
		log.Println("can not unmarshal body:[", string(resBody), "] ", err)
	}
	hash := h.service.Save(value.URL)
	response := Response{h.makeShortURL(hash)}
	responseString, err := json.Marshal(response)
	if err != nil {
		log.Println("can not marshal response:[", string(resBody), "] ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(responseString)
	if err != nil {
		panic(err)
	}
}

func (h *RequestHandler) handleBatch(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	resBody, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	value := Request{}
	if err := json.Unmarshal(resBody, &value); err != nil {
		log.Println("can not unmarshal body:[", string(resBody), "] ", err)
	}
	hash := h.service.Save(value.URL)
	response := Response{h.makeShortURL(hash)}
	responseString, err := json.Marshal(response)
	if err != nil {
		log.Println("can not marshal response:[", string(resBody), "] ", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(responseString)
	if err != nil {
		panic(err)
	}
}

func (h *RequestHandler) getUserUrls(w http.ResponseWriter, r *http.Request) {
	logging.Enter()
	defer logging.Exit()

	userIDCookie, err := r.Cookie(userIDCookieName)
	if err != nil {
		log.Println("err when get userID cookie:", err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err2 := decodeID(userIDCookie.Value)
	if err2 != nil {
		log.Println("err when decode userID cookie:", err2)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	log.Println(strconv.Itoa(int(userID)) + " userID found")

	pairs := h.service.GetAll()
	var list []ListResponse

	for hash, value := range pairs {
		listResponse := ListResponse{h.makeShortURL(hash), value}
		list = append(list, listResponse)
	}
	responseString, err := json.Marshal(list)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseString)
	if err != nil {
		panic(err)
	}
}

func (h *RequestHandler) handlePing(w http.ResponseWriter, r *http.Request) {
	logging.Enter()
	defer logging.Exit()

	if h.dbSource == nil {
		log.Println("db ping error, db is not initialized")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := h.dbSource.Ping()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("db ping error:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *RequestHandler) makeShortURL(hash string) string {
	return h.baseURL + "/" + hash
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type ListResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type BatchRequest struct {
	CorrelationId string `json:"correlation_id"`
	OriginalUrl   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationId string `json:"correlation_id"`
	ShortUrl      string `json:"short_url"`
}
