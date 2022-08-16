package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/caarlos0/env/v6"

	"github.com/Svetopolk/shortener/internal/app/rest"
	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/storage"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
}

func main() {
	log.Println("main started")
	defer log.Println("main finished")

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "-a serverAddress")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "-b baseUrl")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "-f fileStoragePath")
	flag.StringVar(&cfg.DatabaseDsn, "d", cfg.DatabaseDsn, "-d DatabaseDsn")
	flag.Parse()

	var store storage.Storage
	if cfg.FileStoragePath != "" {
		log.Println("environment var FILE_STORAGE_PATH is found: " + cfg.FileStoragePath)
		store = storage.NewFileStorage(cfg.FileStoragePath)
	} else {
		store = storage.NewMemStorage()
	}

	dbSource := db.NewDB(cfg.DatabaseDsn)
	defer dbSource.Close()

	if cfg.DatabaseDsn != "" {
		store = storage.NewDBStorage(dbSource)
	}

	shortService := service.NewShortService(store)

	handler := rest.NewRequestHandler(shortService, cfg.BaseURL, dbSource)
	router := rest.NewRouter(handler)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, router))
}
