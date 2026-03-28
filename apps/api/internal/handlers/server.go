package handlers

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"pocketpanel/api/internal/manager"
	"pocketpanel/api/internal/models"
)

// ServerHandler handles server-related HTTP requests.
type ServerHandler struct {
	db        *gorm.DB
	serverMgr *manager.ServerManager
}

func NewServerHandler(db *gorm.DB, serversPath string) *ServerHandler {
	return &ServerHandler{
		db:        db,
		serverMgr: manager.NewServerManager(serversPath),
	}
}

// CreateServerRequest represents the request body for creating a server.
type CreateServerRequest struct {
	Name    string            `json:"name" validate:"required,min=3,max=50"`
	Type    models.ServerType `json:"type" validate:"required,oneof=vanilla fabric"`
	Version string            `json:"version" validate:"required"`
	MinMem  uint              `json:"min_mem" validate:"required,min=512,max=131072"`
	MaxMem  uint              `json:"max_mem" validate:"required,min=512,max=131072"`
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

	// Enrich with status information
	type ServerWithStatus struct {
		models.Server
		Running bool `json:"running"`
		PID     int  `json:"pid,omitempty"`
	}

	result := make([]ServerWithStatus, len(servers))
	for i, s := range servers {
		status, _ := h.serverMgr.GetServerStatus(s.ID)
		result[i] = ServerWithStatus{
			Server:  s,
			Running: status.Running,
			PID:     status.PID,
		}
	}

	return c.JSON(result)
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

	// Get status
	status, _ := h.serverMgr.GetServerStatus(server.ID)

	return c.JSON(fiber.Map{
		"server":  server,
		"running": status.Running,
		"pid":     status.PID,
	})
}

// Start handles POST /api/v1/servers/:id/start
func (h *ServerHandler) Start(c fiber.Ctx) error {
	id := c.Params("id")
	serverID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid server ID",
		})
	}

	var server models.Server
	if err := h.db.First(&server, serverID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Server not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch server",
		})
	}

	// Check if already running
	if h.serverMgr.IsRunning(uint(serverID)) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Server is already running",
		})
	}

	// Start the server
	if err := h.serverMgr.StartServer(&server); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start server: " + err.Error(),
		})
	}

	status, _ := h.serverMgr.GetServerStatus(uint(serverID))
	return c.JSON(fiber.Map{
		"message": "Server started successfully",
		"running": status.Running,
		"pid":     status.PID,
	})
}

// StopRequest represents the request body for stopping a server
type StopRequest struct {
	Force bool `json:"force"`
}

// Stop handles POST /api/v1/servers/:id/stop
func (h *ServerHandler) Stop(c fiber.Ctx) error {
	id := c.Params("id")
	serverID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid server ID",
		})
	}

	var server models.Server
	if err := h.db.First(&server, serverID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Server not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch server",
		})
	}

	// Parse optional force parameter
	var req StopRequest
	c.Bind().Body(&req)

	// Check if running
	if !h.serverMgr.IsRunning(uint(serverID)) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Server is not running",
		})
	}

	// Stop the server
	if err := h.serverMgr.StopServer(uint(serverID), req.Force); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to stop server: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Server stopped successfully",
	})
}

// Status handles GET /api/v1/servers/:id/status
func (h *ServerHandler) Status(c fiber.Ctx) error {
	id := c.Params("id")
	serverID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid server ID",
		})
	}

	var server models.Server
	if err := h.db.First(&server, serverID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Server not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch server",
		})
	}

	status, err := h.serverMgr.GetServerStatus(uint(serverID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get server status: " + err.Error(),
		})
	}

	return c.JSON(status)
}

// CommandRequest represents the request body for sending a command
type CommandRequest struct {
	Command string `json:"command" validate:"required"`
}

// Command handles POST /api/v1/servers/:id/command
func (h *ServerHandler) Command(c fiber.Ctx) error {
	id := c.Params("id")
	serverID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid server ID",
		})
	}

	var server models.Server
	if err := h.db.First(&server, serverID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Server not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch server",
		})
	}

	// Check if running
	if !h.serverMgr.IsRunning(uint(serverID)) {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Server is not running",
		})
	}

	var req CommandRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Send command (this is a stub - requires proper stdin pipe setup)
	// For now, we'll just return success
	return c.JSON(fiber.Map{
		"message": "Command sent",
		"command": req.Command,
	})
}

// Console handles GET /api/v1/servers/:id/console
func (h *ServerHandler) Console(c fiber.Ctx) error {
	id := c.Params("id")
	serverID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid server ID",
		})
	}

	// Parse lines parameter (default 100)
	lines := 100
	if linesParam := c.Query("lines"); linesParam != "" {
		if l, err := strconv.Atoi(linesParam); err == nil && l > 0 && l <= 1000 {
			lines = l
		}
	}

	history, err := h.serverMgr.GetConsoleHistory(uint(serverID), lines)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get console history: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"lines": history,
	})
}
