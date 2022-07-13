package main

import (
	"log"
	"net/http"

	"github.com/Svetopolk/shortener/internal/app/rest"
	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/storage"
)

func main() {
	memStorage := storage.NewMemStorage()
	shortService := service.NewShortService(memStorage)
	handler := rest.NewRequestHandler(shortService)
	router := rest.NewRouter(handler)
	log.Fatal(http.ListenAndServe(":8080", router))
}
