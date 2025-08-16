package main

import (
	"context"
	"database/sql"
	"log"
	"match-me/config"

	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect directly to database using database/sql
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Check if username column exists
	var exists bool
	checkQuery := `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'users' 
			AND column_name = 'username'
			AND table_schema = 'public'
		);`

	err = db.QueryRowContext(ctx, checkQuery).Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check if username column exists: %v", err)
	}

	if !exists {
		log.Println("✅ Username column does not exist - no migration needed")
		return
	}

	// Drop the username column
	dropQuery := `ALTER TABLE users DROP COLUMN username;`
	
	_, err = db.ExecContext(ctx, dropQuery)
	if err != nil {
		log.Fatalf("Failed to drop username column: %v", err)
	}

	log.Println("✅ Successfully dropped username column from users table")
	log.Println("🎉 Registration should now work properly")
}