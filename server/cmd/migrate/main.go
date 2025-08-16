package main

import (
	"context"
	"log"
	"match-me/config"
	"match-me/ent"
	"match-me/ent/migrate"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	client, err := ent.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer client.Close()

	ctx := context.Background()

	// Create migration schema
	if err := client.Schema.Create(ctx,
		migrate.WithDropColumn(true),
		migrate.WithDropIndex(true),
	); err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}

	log.Println("✅ Migration completed successfully - username column dropped")
}