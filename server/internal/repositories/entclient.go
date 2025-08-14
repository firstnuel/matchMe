package repositories

import (
	"context"
	"log"
	"match-me/config"
	"match-me/ent"

	_ "github.com/lib/pq"
)

func NewEntClient(ctx context.Context, cfg *config.Config) *ent.Client {

	client, err := ent.Open(cfg.DbName, cfg.DbURL)
	if err != nil {
		log.Fatalf("failed opening connection to %s: %v", cfg.DbName, err)
	}

	// Run auto migration
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("Ent client connected and schema created")
	return client
}
