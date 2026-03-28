package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"pocketpanel/api/internal/models"
)

type AuthHandler struct {
	db    *gorm.DB
	store *session.Store
}

func NewAuthHandler(db *gorm.DB, store *session.Store) *AuthHandler {
	return &AuthHandler{db: db, store: store}
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return err // Error handled by customErrorHandler
	}

	var user models.User
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error: " + err.Error(),
		})
	}

	sess.Set("user_id", user.ID)
	sess.Set("username", user.Username)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save session",
		})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (h *AuthHandler) Logout(c fiber.Ctx) error {
	sess, err := h.store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Session error",
		})
	}

	if err := sess.Destroy(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to logout",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) Me(c fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}

	username, ok := c.Locals("username").(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user context",
		})
	}

	return c.JSON(fiber.Map{
		"id":       userID,
		"username": username,
	})
}
