package server

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gorm.io/gorm"

	"pocketpanel/api/internal/config"
	"pocketpanel/api/internal/handlers"
	"pocketpanel/api/internal/middleware"
)

func New(cfg *config.Config, db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "PocketPanel API",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		BodyLimit:    4 * 1024 * 1024, // 4MB
	})

	// Session store
	sessionStore := session.NewStore(session.Config{
		AbsoluteTimeout: 24 * time.Hour,
		CookieHTTPOnly:  true,
		CookieSecure:    cfg.Environment == "production",
		CookieSameSite:  "Lax",
	})

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	api := app.Group("/api/v1")

	authHandler := handlers.NewAuthHandler(db, sessionStore)
	api.Post("/auth/register", authHandler.Register)
	api.Post("/auth/login", authHandler.Login)
	api.Post("/auth/logout", authHandler.Logout)

	protected := api.Group("/", middleware.Auth(sessionStore))
	protected.Get("/me", authHandler.Me)

	return app
}
