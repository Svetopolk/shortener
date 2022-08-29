package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"

	"github.com/Svetopolk/shortener/internal/app/db"
	"github.com/Svetopolk/shortener/internal/app/rest"
	"github.com/Svetopolk/shortener/internal/app/service"
	"github.com/Svetopolk/shortener/internal/app/storage"
	"github.com/Svetopolk/shortener/internal/logging"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
}

func main() {
	logging.Enter()
	defer logging.Exit()

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Println("error while reading config file:", err)
	}

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "-a serverAddress")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "-b baseUrl")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "-f fileStoragePath")
	flag.StringVar(&cfg.DatabaseDsn, "d", cfg.DatabaseDsn, "-d DatabaseDsn")
	flag.Parse()

	log.Printf("config: %+v", cfg)

	var (
		store    storage.Storage
		dbSource *db.Source
	)
	switch {
	case cfg.DatabaseDsn != "":
		log.Println("init store as database based")
		dbSource, err = db.NewDB(cfg.DatabaseDsn)
		if err != nil {
			log.Fatal("failed to init dbSource: " + err.Error())
		}
		defer dbSource.Close()
		store = storage.NewDBStorage(dbSource)
	case cfg.FileStoragePath != "":
		log.Println("environment var FILE_STORAGE_PATH is found: " + cfg.FileStoragePath)
		log.Println("init store as file store based")
		store = storage.NewFileStorage(cfg.FileStoragePath)
	default:
		log.Println("init store as memory store based")
		store = storage.NewMemStorage()
	}

	shortService := service.NewShortService(store)

	handler := rest.NewRequestHandler(shortService, cfg.BaseURL, dbSource)
	router := rest.NewRouter(handler)

	if err = http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Println("listen and serve failed: " + err.Error())
	}
}
