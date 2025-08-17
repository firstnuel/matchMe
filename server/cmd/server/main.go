package main

import (
	"context"
	"log"
	"match-me/api"
	"match-me/config"
	"match-me/internal/repositories"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

// handleServerLifecycle manages the lifecycle of the HTTP server
func handleServerLifecycle(srv *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	serverErr := make(chan error, 1)
	go func() {
		log.Printf("HTTP server listening on addr %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		log.Fatalf("Server error: %v", err)

	case <-stop:
		log.Println("Shutdown signal received. Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("HTTP server shutdown error: %v", err)
		} else {
			log.Printf("HTTP server shutdown cleanly")
		}
	}

	log.Printf("Graceful shutdown complete.")
}
