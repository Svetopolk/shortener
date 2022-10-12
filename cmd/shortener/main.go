package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/Svetopolk/shortener/internal/server"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := server.NewByConfig()

	server.Start(ctx)

	//<-ctx.Done()
	//log.Println("shutting down server gracefully start")
	//server.Shutdown()
	//log.Println("shutting down server gracefully finish")
}
