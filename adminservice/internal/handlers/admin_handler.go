package handlers

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func TestHandler(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}
