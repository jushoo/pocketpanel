package handlers

import (
	"errors"
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
			errMsg := err.Error()
			if strings.Contains(errMsg, "port") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": "A server is already using this port",
				})
			}
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

// List handles GET /api/v1/servers
func (h *ServerHandler) List(c fiber.Ctx) error {
	var servers []models.Server
	if err := h.db.Find(&servers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch servers",
		})
	}

	return c.JSON(servers)
}

// Get handles GET /api/v1/servers/:id
func (h *ServerHandler) Get(c fiber.Ctx) error {
	id := c.Params("id")

	var server models.Server
	if err := h.db.First(&server, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Server not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch server",
		})
	}

	return c.JSON(server)
}
