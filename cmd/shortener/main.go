package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Svetopolk/shortener/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := server.NewByConfig()

	server.Run()

	<-ctx.Done()
	log.Println("shutting down server gracefully start")
	server.Shutdown()
	log.Println("shutting down server gracefully finish")
}
