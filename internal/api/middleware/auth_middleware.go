package middleware

import (
	"log"

	"github.com/gofiber/fiber/v2"
	config "github.com/maheshrc27/postflow/configs"
	"github.com/maheshrc27/postflow/pkg/utils"
)

func AuthMiddleware(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Cookies(cfg.CookieName)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication cookie missing",
			})
		}

		claims, err := utils.ValidateToken(cfg.SecretKey, tokenString)
		if err != nil {
			c.Cookie(&fiber.Cookie{
				Name:   cfg.CookieName,
				Value:  "",
				Path:   "/",
				MaxAge: -1, // Delete cookie
			})

			log.Printf("Token validation failed: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals("user_id", claims.UserID)
		return c.Next()
	}
}
