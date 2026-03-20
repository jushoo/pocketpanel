package main

import (
	"log"

	"pocketpanel/api/internal/config"
	"pocketpanel/api/internal/database"
	"pocketpanel/api/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.Connect(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	app := server.New(cfg, db)

	log.Printf("Server starting on http://localhost%s", cfg.Port)
	if err := app.Listen(cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
