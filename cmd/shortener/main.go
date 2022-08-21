package main

import (
	"flag"
	"github.com/Svetopolk/shortener/internal/logging"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	var store storage.Storage
	if cfg.FileStoragePath != "" {
		log.Println("environment var FILE_STORAGE_PATH is found: " + cfg.FileStoragePath)
		store = storage.NewFileStorage(cfg.FileStoragePath)
	} else {
		store = storage.NewMemStorage()
	}

	dbSource := db.NewDB(cfg.DatabaseDsn)
	//defer dbSource.Close()

	if cfg.DatabaseDsn != "" {
		store = storage.NewDBStorage(dbSource)
	}

	shortService := service.NewShortService(store)

	handler := rest.NewRequestHandler(shortService, cfg.BaseURL, dbSource)
	router := rest.NewRouter(handler)

	server := &http.Server{Addr: cfg.ServerAddress, Handler: router}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := server.ListenAndServe(); err != nil {
			log.Println("listen and serve failed: " + err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		ch := make(chan os.Signal)
		signal.Notify(ch,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)

		s := <-ch

		log.Println("received signal: " + s.String())
		if err := server.Close(); err != nil {
			log.Println("close failed: " + err.Error())
		}
		return
	}()

	wg.Wait()
}
