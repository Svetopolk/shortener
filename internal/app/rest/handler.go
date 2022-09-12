package rest

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/app/exceptions"
	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/util"
	"github.com/Svetopolk/shortener/internal/logging"
)

type RequestHandler struct {
	service  service.ShortService
	baseURL  string
	dbSource *db.Source
}

func NewRequestHandler(service service.ShortService, baseURL string, db *db.Source) *RequestHandler {
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
	hash, err := h.service.Save(string(body))
	if errors.Is(err, exceptions.ErrURLAlreadyExist) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	shortURL := h.makeShortURL(hash)
	_, err = w.Write([]byte(shortURL))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	logging.Enter()
	defer logging.Exit()

	hash := util.RemoveFirstSymbol(r.URL.Path)
	fullURL, err := h.service.Get(hash)
	if err != nil {
		fullURL = "" // TODO if not found what to do?
	}
	w.Header().Set("Location", fullURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
	_, err = w.Write([]byte("redirect to " + fullURL))
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

	w.Header().Set("Content-Type", "application/json")

	hash, err := h.service.Save(value.URL)
	if errors.Is(err, exceptions.ErrURLAlreadyExist) {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	response := Response{h.makeShortURL(hash)}
	responseString, err := json.Marshal(response)
	if err != nil {
		log.Println("can not marshal response:[", string(resBody), "] ", err)
	}

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

	var batchRequests []BatchRequest
	if err := json.Unmarshal(resBody, &batchRequests); err != nil {
		log.Println("can not unmarshal body:[", string(resBody), "] ", err)
	}

	batchResponses := make([]BatchResponse, 0, len(batchRequests))

	// TODO make it through batch SaveBatch
	for i := range batchRequests {
		hash, err := h.service.Save(batchRequests[i].OriginalURL)
		if err != nil {
			log.Println("unexpected exceptions", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		batchResponse := BatchResponse{batchRequests[i].CorrelationID, h.makeShortURL(hash)}
		batchResponses = append(batchResponses, batchResponse)
	}

	responseString, err := json.Marshal(batchResponses)
	if err != nil {
		log.Println("can not marshal batchResponses:[", string(resBody), "] ", err)
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

	pairs, err3 := h.service.GetAll()

	if err3 != nil {
		log.Println("err when get data from storage", err3)
		w.WriteHeader(http.StatusNoContent)
		return
	}

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
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
