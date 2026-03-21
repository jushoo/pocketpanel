package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"pocketpanel/api/internal/models"
)

// ServerHandler handles server-related HTTP requests.
type ServerHandler struct {
	db *gorm.DB
}

func NewServerHandler(db *gorm.DB) *ServerHandler {
	return &ServerHandler{db: db}
}

// CreateServerRequest represents the request body for creating a server.
type CreateServerRequest struct {
	Name    string            `json:"name" validate:"required,min=3,max=50"`
	Type    models.ServerType `json:"type" validate:"required,oneof=vanilla fabric"`
	Version string            `json:"version" validate:"required"`
	MinMem  uint              `json:"min_mem" validate:"required,min=1,max=128"`
	MaxMem  uint              `json:"max_mem" validate:"required,min=1,max=128"`
	Port    uint              `json:"port" validate:"required,min=25565,max=65535"`
}

// Create handles POST /api/v1/servers
func (h *ServerHandler) Create(c fiber.Ctx) error {
	var req CreateServerRequest
	if err := c.Bind().Body(&req); err != nil {
		return err // Error handled by customErrorHandler
	}

	if req.MinMem > req.MaxMem {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "min_mem cannot be greater than max_mem",
		})
	}

	server := models.Server{
		Name:    req.Name,
		Type:    req.Type,
		Version: req.Version,
		MinMem:  req.MinMem,
		MaxMem:  req.MaxMem,
		Port:    req.Port,
	}

	if err := h.db.Create(&server).Error; err != nil {
		if isUniqueViolation(err) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Server with this name already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create server",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(server)
}

// isUniqueViolation checks if the error is a unique constraint violation.
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()
	return strings.Contains(msg, "UNIQUE constraint failed") ||
		strings.Contains(msg, "duplicate key value") ||
		strings.Contains(msg, "Error 1062: Duplicate entry")
}
