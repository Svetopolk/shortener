package main

import (
	"log"
	"net/http"

	"github.com/Svetopolk/shortener/internal/app/rest"
	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

func main() {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	memStorage := storage.NewMemStorage()
	shortService := service.NewShortService(memStorage)
	handler := rest.NewRequestHandler(shortService, cfg.BaseURL)
	router := rest.NewRouter(handler)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}
