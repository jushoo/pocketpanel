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

	// Create default admin account if it doesn't exist
	password, err := database.CreateDefaultAdmin(db)
	if err != nil {
		log.Fatalf("Failed to create default admin: %v", err)
	}
	if password != "" {
		log.Printf("=================================================")
		log.Printf("Default admin account created!")
		log.Printf("Username: admin")
		log.Printf("Password: %s", password)
		log.Printf("=================================================")
	} else {
		log.Println("Admin account already exists")
	}

	app := server.New(cfg, db)

	log.Printf("Server starting on http://localhost%s", cfg.Port)
	if err := app.Listen(cfg.Port); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
