package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"match-me/api"
	"match-me/config"
	"match-me/internal/pkg/cloudinary"
	"match-me/internal/pkg/seed"
	"match-me/internal/repositories"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	ShutdownTimeout = 60 * time.Second
	InitTimeout     = 30 * time.Second
)

func main() {
	// Parse command line flags
	var populateFlag = flag.String("p", "", "Populate database with n users (e.g., -p 100)")
	var resetFlag = flag.Bool("r", false, "Reset database (delete all data)")
	var resetPopulateFlag = flag.String("rp", "", "Reset database and populate with n users (e.g., -rp 100)")
	var helpFlag = flag.Bool("h", false, "Show help and usage information")
	flag.Parse()

	// Handle help flag
	if *helpFlag {
		printUsage()
		return
	}

	// Load configuration
	cfg := config.LoadConfig()

	initCtx, cancelInit := context.WithTimeout(context.Background(), InitTimeout)
	defer cancelInit()

	// Initialize database client
	client := repositories.NewEntClient(initCtx, cfg)
	defer client.Close()
	log.Printf("Database client initialized")

	// Handle database operations flags
	seeder := seed.NewSeeder(client)

	// Handle reset flag
	if *resetFlag {
		if err := seeder.ResetDatabase(initCtx); err != nil {
			log.Fatalf("Failed to reset database: %v", err)
		}
		return
	}

	// Handle reset and populate flag
	if *resetPopulateFlag != "" {
		n, err := strconv.Atoi(*resetPopulateFlag)
		if err != nil || n <= 0 || n > 100 {
			log.Fatalf("Invalid number for reset-populate flag: %s. Please provide a positive integer between (0-100).", *resetPopulateFlag)
		}

		if err := seeder.ResetDatabase(initCtx); err != nil {
			log.Fatalf("Failed to reset database: %v", err)
		}

		if err := seeder.PopulateUsers(initCtx, n); err != nil {
			log.Fatalf("Failed to populate database: %v", err)
		}
		log.Printf("Database reset and population completed successfully")
		return
	}

	// Handle population flag
	if *populateFlag != "" {
		n, err := strconv.Atoi(*populateFlag)
		if err != nil || n <= 0 || n > 100 {
			log.Fatalf("Invalid number for population flag: %s. Please provide a positive integer between (0-100).", *populateFlag)
		}

		if err := seeder.PopulateUsers(initCtx, n); err != nil {
			log.Fatalf("Failed to populate database: %v", err)
		}
		log.Printf("Database population completed successfully")
		return
	}

	// set up media storage
	cld := cloudinary.NewCloudinary()

	// Initialize and start HTTP server
	srv := api.NewHTTPServer(client, cfg, cld)
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

// printUsage displays help information for command line flags
func printUsage() {
	fmt.Println("Match-Me Server - Database Management and Server Operations")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  server [FLAGS]")
	fmt.Println()
	fmt.Println("FLAGS:")
	fmt.Println("  -h          Show help and usage information")
	fmt.Println("  -r          Reset database (delete all data)")
	fmt.Println("  -p n        Populate database with n users (1-100)")
	fmt.Println("  -rp n       Reset database and populate with n users (1-100)")
	fmt.Println()
	fmt.Println("NOTE: Database operations (-r, -p, -rp) will exit after completion.")
	fmt.Println("      Only use -p or -rp when you want to add test data.")
}
