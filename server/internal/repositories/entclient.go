package repositories

import (
	"context"
	"database/sql"
	"log"
	"match-me/config"
	"match-me/ent"
	"match-me/ent/schema"
	"match-me/internal/repositories/hooks"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/lib/pq"
)

func NewEntClient(ctx context.Context, cfg *config.Config) *ent.Client {

	// 1. Open a standard database connection.
	db, err := sql.Open(cfg.DbName, cfg.DbURL)
	if err != nil {
		log.Fatalf("failed opening connection to %s: %v", cfg.DbName, err)
	}

	// 2. Try to enable PostGIS extension (optional for development)
	postgisExt := schema.PostGISExtension{}
	if _, err := db.ExecContext(ctx, postgisExt.SQL()); err != nil {
		log.Printf("Warning: PostGIS extension not available: %v", err)
		log.Println("Spatial queries will not work. Install PostGIS with: brew install postgis")
	} else {
		log.Println("PostGIS extension enabled")
	}

	// 3. Create an Ent driver that wraps our existing connection.
	drv := entsql.OpenDB(dialect.Postgres, db)

	// 4. Create the Ent client with the custom driver.
	client := ent.NewClient(ent.Driver(drv))

	// Register hooks
	client.User.Use(hooks.ProfileCompletionHook())
	client.UserPhoto.Use(hooks.PhotoCompletionHook())

	// Run auto migration with the fully configured client.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	log.Println("Ent client connected and schema created")
	return client
}
