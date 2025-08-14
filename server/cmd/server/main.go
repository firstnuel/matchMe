package main

import (
	"context"
	"log"
	"match-me/api"
	"match-me/config"
	"match-me/internal/repositories"
	"time"
)

const (
	ShutdownTimeout = 60 * time.Second
	InitTimeout     = 30 * time.Second
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	initCtx, cancelInit := context.WithTimeout(context.Background(), InitTimeout)
	defer cancelInit()

	// Initialize database client
	client := repositories.NewEntClient(initCtx, cfg)
	defer client.Close()
	log.Printf("Database client initialized")

	// Initialize and start HTTP server
	srv := api.NewHTTPServer(client, cfg)
	handleServerLifecycle(srv)
}
