package main

import (
	"log"
	"time"

	"pocketpanel/api/internal/config"
	"pocketpanel/api/internal/database"
	"pocketpanel/api/internal/models"
	"pocketpanel/api/internal/server"
	"pocketpanel/api/internal/sync"

	"github.com/go-co-op/gocron/v2"
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

	versionsSyncer := sync.NewSyncer(db,
		sync.NewMojangFetcher(),
		sync.NewFabricFetcher(),
	)

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	if _, err := scheduler.NewJob(
		gocron.DurationJob(6*time.Hour),
		gocron.NewTask(func() {
			log.Println("Syncing server versions...")
			if err := versionsSyncer.SyncAll(); err != nil {
				log.Printf("Failed to sync versions: %v", err)
			}
		}),
	); err != nil {
		log.Fatalf("Failed to schedule version sync: %v", err)
	}

	scheduler.Start()

	var versionCount int64
	db.Model(&models.Version{}).Count(&versionCount)
	if versionCount == 0 {
		log.Println("Syncing server versions (initial)...")
		if err := versionsSyncer.SyncAll(); err != nil {
			log.Printf("Initial version sync failed: %v", err)
		}
	} else {
		log.Printf("Versions already synced (%d entries), skipping initial sync", versionCount)
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
