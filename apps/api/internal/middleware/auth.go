package middleware

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
)

func Auth(store *session.Store) fiber.Handler {
	return func(c fiber.Ctx) error {
		sess, err := store.Get(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Session error",
			})
		}

		userID := sess.Get("user_id")
		username := sess.Get("username")

		if userID == nil || username == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		c.Locals("userID", userID)
		c.Locals("username", username)

		return c.Next()
	}
}
