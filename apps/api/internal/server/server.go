package server

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
	"gorm.io/gorm"

	"pocketpanel/api/internal/config"
	"pocketpanel/api/internal/handlers"
	"pocketpanel/api/internal/middleware"
	validatormiddleware "pocketpanel/api/internal/validator"
)

func New(cfg *config.Config, db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:         "PocketPanel API",
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		BodyLimit:       4 * 1024 * 1024, // 4MB
		StructValidator: validatormiddleware.New(),
		ErrorHandler:    customErrorHandler,
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

	serverHandler := handlers.NewServerHandler(db)
	api.Post("/servers", serverHandler.Create)
	api.Get("/servers", serverHandler.List)
	api.Get("/servers/:id", serverHandler.Get)

	versionsHandler := handlers.NewVersionsHandler(db)
	api.Get("/versions/:type", versionsHandler.GetVersions)

	protected := api.Group("/", middleware.Auth(sessionStore))
	protected.Get("/me", authHandler.Me)

	return app
}

// customErrorHandler formats validation and binding errors into user-friendly responses.
func customErrorHandler(c fiber.Ctx, err error) error {
	// Handle validation errors from go-playground/validator
	if errors, ok := err.(validator.ValidationErrors); ok {
		var messages []string
		for _, e := range errors {
			// Get the human-readable field name from JSON tag
			field := e.Field()
			switch e.Tag() {
			case "required":
				messages = append(messages, field+" is required")
			case "min":
				messages = append(messages, field+" is too short")
			case "max":
				messages = append(messages, field+" is too long")
			case "oneof":
				messages = append(messages, field+" must be one of: "+e.Param())
			default:
				messages = append(messages, field+" is invalid")
			}
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": strings.Join(messages, "; "),
		})
	}

	// Handle binding errors (invalid JSON, missing fields, etc.)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	return nil
}
